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
type BackoffFn func(ctx context.Context, attempt int)

// IsRetryableFn2 is called to determine if the type E is retryable.
type IsRetryableFn2[E any] func(ctx context.Context, e E) (isRetryable bool)

// Retryer is used to retry any fallible operation.
type Retryer[T, E any] struct {
	isRetryableFn   IsRetryableFn2[E]
	maxAttemptsMode MaxAttemptsMode
	maxAttempts     uint8
	bo              BackoffFn
	timeout         time.Duration
}

//TODO: Add example usages to documentation, reminder these are building blocks and although can be used directly
//      are more likely to be used as part of a higher level function.

// NewRetryer returns a new `Retryer` with sane default values.
//
// The default values are:
// - `IsRetryableFn` will always return false as `E` is unknown until defined.
// - `MaxAttemptsMode` is `MaxAttemptsNonRetryableReset`.
// - `MaxAttempts` is 5.
// - `BackoffFn` will sleep for 200ms. It's recommended to use exponential backoff for production.
// - `Timeout` is 0.
func NewRetryer[T, E any]() Retryer[T, E] {
	return Retryer[T, E]{
		isRetryableFn:   func(_ context.Context, _ E) bool { return false },
		maxAttemptsMode: MaxAttemptsNonRetryableReset,
		maxAttempts:     5,
		bo: func(ctx context.Context, attempt int) {
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

// MaxAttempts sets the maximum number of attempts for the `Retryer`.
//
// NOTE: Max attempts is optional and if not set will retry indefinitely on retryable errors.
func (r Retryer[T, E]) MaxAttempts(mode MaxAttemptsMode, maxAttempts uint8) Retryer[T, E] {
	r.maxAttemptsMode, r.maxAttempts = mode, maxAttempts
	return r
}

// Backoff sets the backoff function for the `Retryer`.
func (r Retryer[T, E]) Backoff(fn BackoffFn) Retryer[T, E] {
	if fn == nil {
		fn = func(_ context.Context, _ int) {}
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
			isRetryable := r.isRetryableFn(ctx, result.Err())

			switch r.maxAttemptsMode {
			case MaxAttemptsUnlimited:
				goto END
			case MaxAttemptsNonRetryableReset:
				if isRetryable {
					remaining = r.maxAttempts
				} else if remaining > 0 {
					remaining--
				}
			case MaxAttemptsNonRetryable:
				if !isRetryable {
					if remaining > 0 {
						remaining--
					}
				} else {
					goto END
				}
			case MaxAttempts:
				if remaining > 0 {
					remaining--
				}
			}
			if remaining == 0 {
				return result
			}
		END:
			r.bo(ctx, attempt)
			attempt++
			continue
		}
		return result
	}
}
