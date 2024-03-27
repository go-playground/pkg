//go:build go1.18
// +build go1.18

package httpext

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	bytesext "github.com/go-playground/pkg/v5/bytes"
	errorsext "github.com/go-playground/pkg/v5/errors"
	ioext "github.com/go-playground/pkg/v5/io"
	typesext "github.com/go-playground/pkg/v5/types"
	valuesext "github.com/go-playground/pkg/v5/values"
	. "github.com/go-playground/pkg/v5/values/result"
)

// ErrStatusCode can be used to treat/indicate a status code as an error and ability to indicate if it is retryable.
type ErrStatusCode struct {
	// StatusCode is the HTTP response status code that was encountered.
	StatusCode int

	// IsRetryableStatusCode indicates if the status code is considered retryable.
	IsRetryableStatusCode bool

	// Headers contains the headers from the HTTP response.
	Headers http.Header

	// Body is the optional body of the HTTP response.
	Body []byte
}

// Error returns the error message for the status code.
func (e ErrStatusCode) Error() string {
	return "status code encountered: " + strconv.Itoa(e.StatusCode)
}

// IsRetryable returns if the provided status code is considered retryable.
func (e ErrStatusCode) IsRetryable() bool {
	return e.IsRetryableStatusCode
}

// BuildRequestFn2 is a function used to rebuild an HTTP request for use in retryable code.
type BuildRequestFn2 func(ctx context.Context) Result[*http.Request, error]

// DecodeAnyFn is a function used to decode the response body into the desired type.
type DecodeAnyFn func(ctx context.Context, resp *http.Response, maxMemory bytesext.Bytes, v any) error

// IsRetryableStatusCodeFn2 is a function used to determine if the provided status code is considered retryable.
type IsRetryableStatusCodeFn2 func(ctx context.Context, code int) bool

// Retryer is used to retry any fallible operation.
//
// The `Retryer` is designed to be stateless and reusable. Configuration is also copy and so a base `Retryer` can be
// used and changed for one-off requests eg. changing max attempts resulting in a new `Retrier` for that request.
type Retryer struct {
	isRetryableFn           errorsext.IsRetryableFn2[error]
	isRetryableStatusCodeFn IsRetryableStatusCodeFn2
	isEarlyReturnFn         errorsext.EarlyReturnFn[error]
	decodeFn                DecodeAnyFn
	backoffFn               errorsext.BackoffFn[error]
	client                  *http.Client
	timeout                 time.Duration
	maxBytes                bytesext.Bytes
	mode                    errorsext.MaxAttemptsMode
	maxAttempts             uint8
}

// NewRetryer returns a new `Retryer` with sane default values.
//
// The default values are:
//   - `IsRetryableFn` uses the existing `errorsext.IsRetryableHTTP` function.
//   - `MaxAttemptsMode` is `MaxAttemptsNonRetryableReset`.
//   - `MaxAttempts` is 5.
//   - `BackoffFn` will sleep for 200ms or is successful `Retry-After` header can be parsed. It's recommended to use
//     exponential backoff for production with a quick copy-paste-modify of the default function
//   - `Timeout` is 0.
//   - `IsRetryableStatusCodeFn` is set to the existing `IsRetryableStatusCode` function.
//   - `IsEarlyReturnFn` is set to check if the error is an `ErrStatusCode` and if the status code is non-retryable.
//   - `Client` is set to `http.DefaultClient`.
//   - `MaxBytes` is set to 2MiB.
//   - `DecodeAnyFn` is set to the existing `DecodeResponseAny` function that supports JSON and XML.
//
// WARNING: The default functions may receive enhancements or fixes in the future which could change their behavior,
// however every attempt will be made to maintain backwards compatibility or made additive-only if possible.
func NewRetryer() Retryer {
	return Retryer{
		client:      http.DefaultClient,
		maxBytes:    2 * bytesext.MiB,
		mode:        errorsext.MaxAttemptsNonRetryableReset,
		maxAttempts: 5,
		isRetryableFn: func(ctx context.Context, err error) (isRetryable bool) {
			_, isRetryable = errorsext.IsRetryableHTTP(err)
			return
		},
		isRetryableStatusCodeFn: func(_ context.Context, code int) bool { return IsRetryableStatusCode(code) },
		isEarlyReturnFn: func(_ context.Context, err error) bool {
			var sce ErrStatusCode
			if errors.As(err, &sce) {
				return IsNonRetryableStatusCode(sce.StatusCode)
			}
			return false
		},
		decodeFn: func(ctx context.Context, resp *http.Response, maxMemory bytesext.Bytes, v any) error {
			err := DecodeResponseAny(resp, maxMemory, v)
			if err != nil {
				return err
			}
			return nil
		},
		backoffFn: func(ctx context.Context, attempt int, err error) {

			wait := time.Millisecond * 200

			var sce ErrStatusCode
			if errors.As(err, &sce) {
				if sce.Headers != nil && (sce.StatusCode == http.StatusTooManyRequests || sce.StatusCode == http.StatusServiceUnavailable) {
					if ra := HasRetryAfter(sce.Headers); ra.IsSome() {
						wait = ra.Unwrap()
					}
				}
			}

			t := time.NewTimer(wait)
			defer t.Stop()
			select {
			case <-ctx.Done():
			case <-t.C:
			}
		},
	}
}

// Client sets the `http.Client` for the `Retryer`.
func (r Retryer) Client(client *http.Client) Retryer {
	r.client = client
	return r
}

// IsRetryableFn sets the `IsRetryableFn` for the `Retryer`.
func (r Retryer) IsRetryableFn(fn errorsext.IsRetryableFn2[error]) Retryer {
	r.isRetryableFn = fn
	return r
}

// IsRetryableStatusCodeFn is called to determine if the status code is retryable.
func (r Retryer) IsRetryableStatusCodeFn(fn IsRetryableStatusCodeFn2) Retryer {
	if fn == nil {
		fn = func(_ context.Context, _ int) bool { return false }
	}
	r.isRetryableStatusCodeFn = fn
	return r
}

// IsEarlyReturnFn sets the `EarlyReturnFn` for the `Retryer`.
func (r Retryer) IsEarlyReturnFn(fn errorsext.EarlyReturnFn[error]) Retryer {
	r.isEarlyReturnFn = fn
	return r
}

// DecodeFn sets the decode function for the `Retryer`.
func (r Retryer) DecodeFn(fn DecodeAnyFn) Retryer {
	if fn == nil {
		fn = func(_ context.Context, _ *http.Response, _ bytesext.Bytes, _ any) error { return nil }
	}
	r.decodeFn = fn
	return r
}

// MaxAttempts sets the maximum number of attempts for the `Retryer`.
//
// NOTE: Max attempts is optional and if not set will retry indefinitely on retryable errors.
func (r Retryer) MaxAttempts(mode errorsext.MaxAttemptsMode, maxAttempts uint8) Retryer {
	r.mode, r.maxAttempts = mode, maxAttempts
	return r
}

// Backoff sets the backoff function for the `Retryer`.
func (r Retryer) Backoff(fn errorsext.BackoffFn[error]) Retryer {
	r.backoffFn = fn
	return r
}

// MaxBytes sets the maximum memory to use when decoding the response body including:
// - upon unexpected status codes.
// - when decoding the response body.
// - when draining the response body before closing allowing connection re-use.
func (r Retryer) MaxBytes(i bytesext.Bytes) Retryer {
	r.maxBytes = i
	return r

}

// Timeout sets the timeout for the `Retryer`. This is the timeout per `RetyableFn` attempt and not the entirety
// of the `Retryer` execution.
//
// A timeout of 0 will disable the timeout and is the default.
func (r Retryer) Timeout(timeout time.Duration) Retryer {
	r.timeout = timeout
	return r
}

// DoResponse will execute the provided functions code and automatically retry before returning the *http.Response
// based on HTTP status code, if defined, and can be used when processing of the response body may not be necessary
// or something custom is required.
//
// NOTE: it is up to the caller to close the response body if a successful request is made.
func (r Retryer) DoResponse(ctx context.Context, fn BuildRequestFn2, expectedResponseCodes ...int) Result[*http.Response, error] {
	return errorsext.NewRetryer[*http.Response, error]().
		IsRetryableFn(r.isRetryableFn).
		MaxAttempts(r.mode, r.maxAttempts).
		Backoff(r.backoffFn).
		Timeout(r.timeout).
		IsEarlyReturnFn(r.isEarlyReturnFn).
		Do(ctx, func(ctx context.Context) Result[*http.Response, error] {
			req := fn(ctx)
			if req.IsErr() {
				return Err[*http.Response, error](req.Err())
			}

			resp, err := r.client.Do(req.Unwrap())
			if err != nil {
				return Err[*http.Response, error](err)
			}

			if len(expectedResponseCodes) > 0 {
				for _, code := range expectedResponseCodes {
					if resp.StatusCode == code {
						goto RETURN
					}
				}
				b, _ := io.ReadAll(ioext.LimitReader(resp.Body, r.maxBytes))
				_ = resp.Body.Close()
				return Err[*http.Response, error](ErrStatusCode{
					StatusCode:            resp.StatusCode,
					IsRetryableStatusCode: r.isRetryableStatusCodeFn(ctx, resp.StatusCode),
					Headers:               resp.Header,
					Body:                  b,
				})
			}

		RETURN:
			return Ok[*http.Response, error](resp)
		})
}

// Do will execute the provided functions code and automatically retry using the provided retry function decoding
// the response body into the desired type `v`, which must be passed as mutable.
func (r Retryer) Do(ctx context.Context, fn BuildRequestFn2, v any, expectedResponseCodes ...int) error {
	result := errorsext.NewRetryer[typesext.Nothing, error]().
		IsRetryableFn(r.isRetryableFn).
		MaxAttempts(r.mode, r.maxAttempts).
		Backoff(r.backoffFn).
		Timeout(r.timeout).
		IsEarlyReturnFn(r.isEarlyReturnFn).
		Do(ctx, func(ctx context.Context) Result[typesext.Nothing, error] {
			req := fn(ctx)
			if req.IsErr() {
				return Err[typesext.Nothing, error](req.Err())
			}

			resp, err := r.client.Do(req.Unwrap())
			if err != nil {
				return Err[typesext.Nothing, error](err)
			}
			defer func() {
				_, _ = io.Copy(io.Discard, ioext.LimitReader(resp.Body, r.maxBytes))
				_ = resp.Body.Close()
			}()

			if len(expectedResponseCodes) > 0 {
				for _, code := range expectedResponseCodes {
					if resp.StatusCode == code {
						goto DECODE
					}
				}

				b, _ := io.ReadAll(ioext.LimitReader(resp.Body, r.maxBytes))
				return Err[typesext.Nothing, error](ErrStatusCode{
					StatusCode:            resp.StatusCode,
					IsRetryableStatusCode: r.isRetryableStatusCodeFn(ctx, resp.StatusCode),
					Headers:               resp.Header,
					Body:                  b,
				})
			}

		DECODE:
			if err = r.decodeFn(ctx, resp, r.maxBytes, v); err != nil {
				return Err[typesext.Nothing, error](err)
			}
			return Ok[typesext.Nothing, error](valuesext.Nothing)
		})
	if result.IsErr() {
		return result.Err()
	}
	return nil
}
