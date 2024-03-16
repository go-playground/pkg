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

	// MaxAttemptsTotal will apply the max attempts to all errors, even those determined to be retryable.
	MaxAttemptsTotal

	// MaxAttemptsInfinite will not apply a maximum number of attempts.
	MaxAttemptsInfinite
)

// BackoffFn is a function used to apply a backoff strategy to the retryable function.
type BackoffFn func(ctx context.Context, attempt int)

// IsRetryableFn2 is called to determine if the type E is retryable.
type IsRetryableFn2[E any] func(e E) (isRetryable bool)

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
		isRetryableFn:   func(_ E) bool { return false },
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
	maxAttempts := r.maxAttempts
	for {
		result := fn(ctx)
		if result.IsErr() {
			if r.maxAttemptsMode != MaxAttemptsInfinite && maxAttempts == 0 {
				return result
			}
			if r.isRetryableFn(result.Err()) {
				if r.maxAttemptsMode == MaxAttemptsNonRetryableReset {
					maxAttempts = r.maxAttempts
				} else if r.maxAttemptsMode != MaxAttemptsInfinite {
					maxAttempts--
				}
				r.bo(ctx, attempt)
				attempt++
				continue
			}
		}
		return result
	}
}
