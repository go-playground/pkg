//go:build go1.18
// +build go1.18

package httpext

import (
	"context"
	"fmt"
	bytesext "github.com/go-playground/pkg/v5/bytes"
	errorsext "github.com/go-playground/pkg/v5/errors"
	ioext "github.com/go-playground/pkg/v5/io"
	resultext "github.com/go-playground/pkg/v5/values/result"
	"io"
	"net/http"
	"strconv"
)

var (
	// retryableStatusCodes defines the common HTTP response codes that are considered retryable.
	retryableStatusCodes = map[int]bool{
		http.StatusServiceUnavailable: true,
		http.StatusTooManyRequests:    true,
		http.StatusBadGateway:         true,
		http.StatusGatewayTimeout:     true,
		http.StatusRequestTimeout:     true,

		// 524 is a Cloudflare specific error which indicates it connected to the origin server but did not receive
		// response within 100 seconds and so times out.
		// https://support.cloudflare.com/hc/en-us/articles/115003011431-Error-524-A-timeout-occurred#524error
		524: true,
	}
)

// ErrRetryableStatusCode can be used to indicate a retryable HTTP status code was encountered as an error.
type ErrRetryableStatusCode struct {
	StatusCode int
}

func (e ErrRetryableStatusCode) Error() string {
	return fmt.Sprintf("retryable HTTP status code encountered: %d", e.StatusCode)
}

// IsRetryableStatusCode returns if the provided status code is considered retryable.
func IsRetryableStatusCode(code int) bool {
	return retryableStatusCodes[code]
}

// BuildRequestFn is a function used to rebuild an HTTP request for use in retryable code.
type BuildRequestFn func(ctx context.Context) (*http.Request, error)

// IsRetryableStatusCodeFn is a function used to determine if the provided status code is considered retryable.
//
// this is also an optional place to put backoff strategies, logs and metrics emitting.
type IsRetryableStatusCodeFn func(code int) bool

// DoRetryableResponse will execute the provided functions code and automatically retry before returning the *http.Response.
func DoRetryableResponse(ctx context.Context, client *http.Client, onRetryFn errorsext.OnRetryFn[error], buildFn BuildRequestFn) resultext.Result[*http.Response, error] {
	if client == nil {
		client = http.DefaultClient
	}
	var attempt int
	for {
		req, err := buildFn(ctx)
		if err != nil {
			return resultext.Err[*http.Response, error](err)
		}

		resp, err := client.Do(req)
		if err != nil {
			if retryReason, isRetryable := errorsext.IsRetryableHTTP(err); isRetryable {
				opt := onRetryFn(ctx, err, retryReason, attempt)
				if opt.IsSome() {
					return resultext.Err[*http.Response, error](opt.Unwrap())
				}
				attempt++
				continue
			}
			return resultext.Err[*http.Response, error](err)
		}

		if IsRetryableStatusCode(resp.StatusCode) {
			opt := onRetryFn(ctx, ErrRetryableStatusCode{StatusCode: resp.StatusCode}, strconv.Itoa(resp.StatusCode), attempt)
			if opt.IsSome() {
				return resultext.Err[*http.Response, error](opt.Unwrap())
			}
			attempt++
			continue
		}
		return resultext.Ok[*http.Response, error](resp)
	}
}

// DoRetryable will execute the provided functions code and automatically retry before returning the result.
//
// This function currently supports decoding the following automatically based on the response Content-Type with
// Gzip supported:
// - JSON
// - XML
func DoRetryable[T any](ctx context.Context, client *http.Client, expectedResponseCode int, maxMemory bytesext.Bytes, isRetryableFn errorsext.IsRetryableFn[error], onRetryFn errorsext.OnRetryFn[error], buildFn BuildRequestFn) resultext.Result[T, error] {

	return errorsext.DoRetryable(ctx, isRetryableFn, onRetryFn, func(ctx context.Context) resultext.Result[T, error] {

		result := DoRetryableResponse(ctx, client, onRetryFn, buildFn)
		if result.IsErr() {
			return resultext.Err[T, error](result.Err())
		}
		resp := result.Unwrap()
		defer resp.Body.Close()

		if resp.StatusCode != expectedResponseCode {
			b, _ := io.ReadAll(ioext.LimitReader(resp.Body, maxMemory))
			err := fmt.Errorf("invalid response status code: %d body: %s", resp.StatusCode, string(b))
			return resultext.Err[T, error](err)
		}

		data, err := DecodeResponse[T](resp, maxMemory)
		if err != nil {
			return resultext.Err[T, error](err)
		}
		return resultext.Ok[T, error](data)
	})
}
