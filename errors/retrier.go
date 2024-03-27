//go:build go1.18
// +build go1.18

package errorsext

import (
	"context"
	"time"

	. "github.com/go-playground/pkg/v5/values/result"
)

// MaxAttemptsMode is used to set the mode for the maximum number of attempts.
//
// eg. Should the max attempts apply to all errors, just ones not determined to be retryable, reset on retryable errors, etc.
type MaxAttemptsMode uint8

const (
	// MaxAttemptsNonRetryableReset will apply the max attempts to all errors not determined to be retryable, but will
	// reset the attempts if a retryable error is encountered after a non-retryable error.
	MaxAttemptsNonRetryableReset MaxAttemptsMode = iota

	// MaxAttemptsNonRetryable will apply the max attempts to all errors not determined to be retryable.
	MaxAttemptsNonRetryable

	// MaxAttempts will apply the max attempts to all errors, even those determined to be retryable.
	MaxAttempts

	// MaxAttemptsUnlimited will not apply a maximum number of attempts.
	MaxAttemptsUnlimited
)

// BackoffFn is a function used to apply a backoff strategy to the retryable function.
//
// It accepts `E` in cases where the amount of time to backoff is dynamic, for example when and http request fails
// with a 429 status code, the `Retry-After` header can be used to determine how long to backoff. It is not required
// to use or handle `E` and can be ignored if desired.
type BackoffFn[E any] func(ctx context.Context, attempt int, e E)

// IsRetryableFn2 is called to determine if the type E is retryable.
type IsRetryableFn2[E any] func(ctx context.Context, e E) (isRetryable bool)

// EarlyReturnFn is the function that can be used to bypass all retry logic, no matter the MaxAttemptsMode, for when the
// type of `E` will never succeed and should not be retried.
//
// eg. If retrying an HTTP request and getting 400 Bad Request, it's unlikely to ever succeed and should not be retried.
type EarlyReturnFn[E any] func(ctx context.Context, e E) (earlyReturn bool)

// Retryer is used to retry any fallible operation.
type Retryer[T, E any] struct {
	isRetryableFn   IsRetryableFn2[E]
	isEarlyReturnFn EarlyReturnFn[E]
	maxAttemptsMode MaxAttemptsMode
	maxAttempts     uint8
	bo              BackoffFn[E]
	timeout         time.Duration
}

// NewRetryer returns a new `Retryer` with sane default values.
//
// The default values are:
// - `MaxAttemptsMode` is `MaxAttemptsNonRetryableReset`.
// - `MaxAttempts` is 5.
// - `Timeout` is 0 no context timeout.
// - `IsRetryableFn` will always return false as `E` is unknown until defined.
// - `BackoffFn` will sleep for 200ms. It's recommended to use exponential backoff for production.
// - `EarlyReturnFn` will be None.
func NewRetryer[T, E any]() Retryer[T, E] {
	return Retryer[T, E]{
		isRetryableFn:   func(_ context.Context, _ E) bool { return false },
		maxAttemptsMode: MaxAttemptsNonRetryableReset,
		maxAttempts:     5,
		bo: func(ctx context.Context, attempt int, _ E) {
			t := time.NewTimer(time.Millisecond * 200)
			defer t.Stop()
			select {
			case <-ctx.Done():
			case <-t.C:
			}
		},
	}
}

// IsRetryableFn sets the `IsRetryableFn` for the `Retryer`.
func (r Retryer[T, E]) IsRetryableFn(fn IsRetryableFn2[E]) Retryer[T, E] {
	if fn == nil {
		fn = func(_ context.Context, _ E) bool { return false }
	}
	r.isRetryableFn = fn
	return r
}

// IsEarlyReturnFn sets the `EarlyReturnFn` for the `Retryer`.
//
// NOTE: If the `EarlyReturnFn` and `IsRetryableFn` are both set and a conflicting `IsRetryableFn` will take precedence.
func (r Retryer[T, E]) IsEarlyReturnFn(fn EarlyReturnFn[E]) Retryer[T, E] {
	r.isEarlyReturnFn = fn
	return r
}

// MaxAttempts sets the maximum number of attempts for the `Retryer`.
//
// NOTE: Max attempts is optional and if not set will retry indefinitely on retryable errors.
func (r Retryer[T, E]) MaxAttempts(mode MaxAttemptsMode, maxAttempts uint8) Retryer[T, E] {
	r.maxAttemptsMode, r.maxAttempts = mode, maxAttempts
	return r
}

// Backoff sets the backoff function for the `Retryer`.
func (r Retryer[T, E]) Backoff(fn BackoffFn[E]) Retryer[T, E] {
	if fn == nil {
		fn = func(_ context.Context, _ int, _ E) {}
	}
	r.bo = fn
	return r
}

// Timeout sets the timeout for the `Retryer`. This is the timeout per `RetyableFn` attempt and not the entirety
// of the `Retryer` execution.
//
// A timeout of 0 will disable the timeout and is the default.
func (r Retryer[T, E]) Timeout(timeout time.Duration) Retryer[T, E] {
	r.timeout = timeout
	return r
}

// Do will execute the provided functions code and automatically retry using the provided retry function.
func (r Retryer[T, E]) Do(ctx context.Context, fn RetryableFn[T, E]) Result[T, E] {
	var attempt int
	remaining := r.maxAttempts
	for {
		var result Result[T, E]
		if r.timeout == 0 {
			result = fn(ctx)
		} else {
			ctx, cancel := context.WithTimeout(ctx, r.timeout)
			result = fn(ctx)
			cancel()
		}
		if result.IsErr() {
			err := result.Err()
			isRetryable := r.isRetryableFn(ctx, err)
			if !isRetryable && r.isEarlyReturnFn != nil && r.isEarlyReturnFn(ctx, err) {
				return result
			}

			switch r.maxAttemptsMode {
			case MaxAttemptsUnlimited:
				goto RETRY
			case MaxAttemptsNonRetryableReset:
				if isRetryable {
					remaining = r.maxAttempts
					goto RETRY
				} else if remaining > 0 {
					remaining--
				}
			case MaxAttemptsNonRetryable:
				if isRetryable {
					goto RETRY
				} else if remaining > 0 {
					remaining--
				}
			case MaxAttempts:
				if remaining > 0 {
					remaining--
				}
			}
			if remaining == 0 {
				return result
			}

		RETRY:
			r.bo(ctx, attempt, err)
			attempt++
			continue
		}
		return result
	}
}
