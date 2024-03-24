package httpext

import (
	"context"
	"net/http"
	"strconv"
	"time"

	bytesext "github.com/go-playground/pkg/v5/bytes"
	errorsext "github.com/go-playground/pkg/v5/errors"
	. "github.com/go-playground/pkg/v5/values/result"
)

// ErrStatusCode can be used to treat/indicate a status code as an error and ability to indicate if it is retryable.
type ErrStatusCode struct {
	StatusCode            int
	IsRetryableStatusCode bool
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

// DecodeFn is a function used to decode the response body into the desired type.
type DecodeFn[T any] func(ctx context.Context, resp *http.Response, maxMemory bytesext.Bytes) Result[T, error]

// Retryer is used to retry any fallible operation.
type Retryer[T any] struct {
	isRetryableFn           errorsext.IsRetryableFn2[error]
	isRetryableStatusCodeFn IsRetryableStatusCodeFn
	decodeFn                DecodeFn[T]
	backoffFn               errorsext.BackoffFn
	client                  *http.Client
	timeout                 time.Duration
	maxMemory               bytesext.Bytes
	mode                    errorsext.MaxAttemptsMode
	maxAttempts             uint8
}

// NewRetryer returns a new `Retryer` with sane default values.
//
// The default values are:
// - `IsRetryableFn` uses the existing `errorsext.IsRetryableHTTP` function.
// - `MaxAttemptsMode` is `MaxAttemptsNonRetryableReset`.
// - `MaxAttempts` is 5.
// - `BackoffFn` will sleep for 200ms. It's recommended to use exponential backoff for production.
// - `Timeout` is 0.
// - `IsRetryableStatusCodeFn` is set to the existing `IsRetryableStatusCode` function.
// - `Client` is set to `http.DefaultClient`.
// - `MaxMemory` is set to 2MiB.
// - `DecodeFn` is set to the existing `DecodeResponse` function that supports JSON and XML.
func NewRetryer[T any]() Retryer[T] {
	return Retryer[T]{
		isRetryableFn: func(ctx context.Context, err error) (isRetryable bool) {
			_, isRetryable = errorsext.IsRetryableHTTP(err)
			return
		},
		isRetryableStatusCodeFn: IsRetryableStatusCode,
		decodeFn: func(ctx context.Context, resp *http.Response, maxMemory bytesext.Bytes) Result[T, error] {
			data, err := DecodeResponse[T](resp, maxMemory)
			if err != nil {
				return Err[T, error](err)
			}
			return Ok[T, error](data)
		},
		client:    http.DefaultClient,
		maxMemory: 2 * bytesext.MiB,
	}
}

// Client sets the `http.Client` for the `Retryer`.
func (r Retryer[T]) Client(client *http.Client) Retryer[T] {
	r.client = client
	return r
}

// IsRetryableFn sets the `IsRetryableFn` for the `Retryer`.
func (r Retryer[T]) IsRetryableFn(fn errorsext.IsRetryableFn2[error]) Retryer[T] {
	r.isRetryableFn = fn
	return r
}

// IsRetryableStatusCodeFn is called to determine if the status code is retryable.
func (r Retryer[T]) IsRetryableStatusCodeFn(fn IsRetryableStatusCodeFn) Retryer[T] {
	r.isRetryableStatusCodeFn = fn
	return r
}

// DecodeFn sets the decode function for the `Retryer`.
func (r Retryer[T]) DecodeFn(fn DecodeFn[T]) Retryer[T] {
	r.decodeFn = fn
	return r
}

// MaxAttempts sets the maximum number of attempts for the `Retryer`.
//
// NOTE: Max attempts is optional and if not set will retry indefinitely on retryable errors.
func (r Retryer[T]) MaxAttempts(mode errorsext.MaxAttemptsMode, maxAttempts uint8) Retryer[T] {
	r.mode, r.maxAttempts = mode, maxAttempts
	return r
}

// Backoff sets the backoff function for the `Retryer`.
func (r Retryer[T]) Backoff(fn errorsext.BackoffFn) Retryer[T] {
	r.backoffFn = fn
	return r
}

// MaxMemory sets the maximum memory to use when decoding the response body.
func (r Retryer[T]) MaxMemory(maxMemory bytesext.Bytes) Retryer[T] {
	r.maxMemory = maxMemory
	return r

}

// Timeout sets the timeout for the `Retryer`. This is the timeout per `RetyableFn` attempt and not the entirety
// of the `Retryer` execution.
//
// A timeout of 0 will disable the timeout and is the default.
func (r Retryer[T]) Timeout(timeout time.Duration) Retryer[T] {
	r.timeout = timeout
	return r
}

// DoResponse will execute the provided functions code and automatically retry before returning the *http.Response
// based on HTTP status code, if defined, and can be used when processing of the response body may not be necessary
// or something custom is required.
//
// NOTE: it is up to the caller to close the response body if a successful request is made.
func (r Retryer[T]) DoResponse(ctx context.Context, fn BuildRequestFn2, expectedResponseCodes ...int) Result[*http.Response, error] {
	return errorsext.NewRetryer[*http.Response, error]().
		IsRetryableFn(r.isRetryableFn).
		MaxAttempts(r.mode, r.maxAttempts).
		Backoff(r.backoffFn).
		Timeout(r.timeout).
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
				return Err[*http.Response, error](ErrStatusCode{StatusCode: resp.StatusCode, IsRetryableStatusCode: r.isRetryableStatusCodeFn(resp.StatusCode)})
			}

		RETURN:
			return Ok[*http.Response, error](resp)
		})
}

// Do will execute the provided functions code and automatically retry using the provided retry function.
func (r Retryer[T]) Do(ctx context.Context, fn BuildRequestFn2, expectedResponseCodes ...int) Result[T, error] {
	return errorsext.NewRetryer[T, error]().
		IsRetryableFn(r.isRetryableFn).
		MaxAttempts(r.mode, r.maxAttempts).
		Backoff(r.backoffFn).
		Timeout(r.timeout).
		Do(ctx, func(ctx context.Context) Result[T, error] {
			req := fn(ctx)
			if req.IsErr() {
				return Err[T, error](req.Err())
			}
			resp, err := r.client.Do(req.Unwrap())
			if err != nil {
				return Err[T, error](err)
			}
			defer resp.Body.Close()

			if len(expectedResponseCodes) > 0 {
				for _, code := range expectedResponseCodes {
					if resp.StatusCode == code {
						goto DECODE
					}
				}
				return Err[T, error](ErrStatusCode{StatusCode: resp.StatusCode, IsRetryableStatusCode: r.isRetryableStatusCodeFn(resp.StatusCode)})
			}

		DECODE:
			return r.decodeFn(ctx, resp, r.maxMemory)
		})
}
