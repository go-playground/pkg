package errorsext

import (
	"context"

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
}

// NewRetryer returns a new `Retryer` with sane default values.
func NewRetryer[T, E any]() Retryer[T, E] {
	return Retryer[T, E]{
		isRetryableFn:   func(_ context.Context, _ E) bool { return false },
		maxAttemptsMode: MaxAttemptsNonRetryableReset,
		maxAttempts:     5,
		bo:              func(ctx context.Context, attempt int) {},
	}
}

// IsRetryableFn sets the `IsRetryableFn` for the `Retryer`.
func (r Retryer[T, E]) IsRetryableFn(fn IsRetryableFn2[E]) Retryer[T, E] {
	r.isRetryableFn = fn
	return r
}

// MaxAttempts sets the maximum number of attempts for the `Retryer`.
//
// NOTE: Max attempts is optional and if not set will retry indefinitely on retryable errors.
func (r Retryer[T, E]) MaxAttempts(mode MaxAttemptsMode, maxAttempts uint8) Retryer[T, E] {
	r.maxAttemptsMode = mode
	r.maxAttempts = maxAttempts
	return r
}

// Backoff sets the backoff function for the `Retryer`.
func (r Retryer[T, E]) Backoff(fn BackoffFn) Retryer[T, E] {
	r.bo = fn
	return r
}

// Do will execute the provided functions code and automatically retry using the provided retry function.
func (r Retryer[T, E]) Do(ctx context.Context, fn RetryableFn[T, E]) Result[T, E] {
	var attempt int
	remaining := r.maxAttempts
	for {
		result := fn(ctx)
		if result.IsErr() {
			isRetryable := r.isRetryableFn(ctx, result.Err())

			switch r.maxAttemptsMode {
			case MaxAttemptsUnlimited:
				goto END
			case MaxAttemptsNonRetryableReset:
				if isRetryable {
					remaining = r.maxAttempts
				} else {
					remaining = decrement(remaining)
				}
			case MaxAttemptsNonRetryable:
				if !isRetryable {
					remaining = decrement(remaining)
				}
			case MaxAttempts:
				remaining = decrement(remaining)
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

func decrement(i uint8) uint8 {
	if i == 0 {
		return 0
	}
	return i - 1
}
