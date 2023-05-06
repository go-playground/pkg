package errorsext

import (
	"context"
	optionext "github.com/go-playground/pkg/v5/values/option"
	resultext "github.com/go-playground/pkg/v5/values/result"
)

// RetryableFn is a function that can be retried.
type RetryableFn[T, E any] func(ctx context.Context) resultext.Result[T, E]

// IsRetryableFn is called to determine if the error is retryable and optionally returns the reason for logging and metrics.
type IsRetryableFn[E any] func(err E) (reason string, isRetryable bool)

// OnRetryFn is called after IsRetryableFn returns true and before the retry is attempted.
//
// this allows for interception, short-circuiting and adding of backoff strategies.
type OnRetryFn[E any] func(ctx context.Context, reason string, attempt int) optionext.Option[E]

// DoRetryable will execute the provided functions code and automatically retry using the provided retry function.
func DoRetryable[T, E any](ctx context.Context, isRetryFn IsRetryableFn[E], onRetryFn OnRetryFn[E], fn RetryableFn[T, E]) resultext.Result[T, E] {
	var attempt int
	for {
		result := fn(ctx)
		if result.IsErr() {
			if reason, isRetryable := isRetryFn(result.Err()); isRetryable {
				if opt := onRetryFn(ctx, reason, attempt); opt.IsSome() {
					return resultext.Err[T, E](opt.Unwrap())
				}
				attempt++
				continue
			}
		}
		return result
	}
}
