//go:build go1.18
// +build go1.18

package errorsext

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	. "github.com/go-playground/assert/v2"
	. "github.com/go-playground/pkg/v5/values/result"
)

// TODO: Add IsRetryable and Retryable to helper functions.

func TestRetrierMaxAttempts(t *testing.T) {
	var i, j int
	result := NewRetryer[int, error]().Backoff(func(ctx context.Context, attempt int, _ error) {
		j++
	}).MaxAttempts(MaxAttempts, 3).Do(context.Background(), func(ctx context.Context) Result[int, error] {
		i++
		if i > 50 {
			panic("infinite loop")
		}
		return Err[int, error](io.EOF)
	})
	Equal(t, result.IsErr(), true)
	Equal(t, result.Err(), io.EOF)
	Equal(t, i, 3)
	Equal(t, j, 2)
}

func TestRetrierMaxAttemptsNonRetryable(t *testing.T) {
	var i, j int
	returnErr := io.ErrUnexpectedEOF
	result := NewRetryer[int, error]().IsRetryableFn(func(_ context.Context, e error) (isRetryable bool) {
		if returnErr == io.EOF {
			return false
		} else {
			return true
		}
	}).Backoff(func(ctx context.Context, attempt int, _ error) {
		j++
		if j == 10 {
			returnErr = io.EOF
		}
	}).MaxAttempts(MaxAttemptsNonRetryable, 3).Do(context.Background(), func(ctx context.Context) Result[int, error] {
		i++
		if i > 50 {
			panic("infinite loop")
		}
		return Err[int, error](returnErr)
	})
	Equal(t, result.IsErr(), true)
	Equal(t, result.Err(), io.EOF)
	Equal(t, i, 13)
	Equal(t, j, 12)
}

func TestRetrierMaxAttemptsNonRetryableReset(t *testing.T) {
	var i, j int
	returnErr := io.EOF
	result := NewRetryer[int, error]().IsRetryableFn(func(_ context.Context, e error) (isRetryable bool) {
		if returnErr == io.EOF {
			return false
		} else {
			return true
		}
	}).Backoff(func(ctx context.Context, attempt int, _ error) {
		j++
		if j == 2 {
			returnErr = io.ErrUnexpectedEOF
		} else if j == 10 {
			returnErr = io.EOF
		}
	}).MaxAttempts(MaxAttemptsNonRetryableReset, 3).Do(context.Background(), func(ctx context.Context) Result[int, error] {
		i++
		if i > 50 {
			panic("infinite loop")
		}
		return Err[int, error](returnErr)
	})
	Equal(t, result.IsErr(), true)
	Equal(t, result.Err(), io.EOF)
	Equal(t, i, 13)
	Equal(t, j, 12)
}

func TestRetrierMaxAttemptsUnlimited(t *testing.T) {
	var i, j int
	r := NewRetryer[int, error]().Backoff(func(ctx context.Context, attempt int, _ error) {
		j++
	}).MaxAttempts(MaxAttemptsUnlimited, 0)

	PanicMatches(t, func() {
		r.Do(context.Background(), func(ctx context.Context) Result[int, error] {
			i++
			if i > 50 {
				panic("infinite loop")
			}
			return Err[int, error](io.EOF)
		})
	}, "infinite loop")
}

func TestRetrierMaxAttemptsTimeout(t *testing.T) {
	result := NewRetryer[int, error]().Backoff(func(ctx context.Context, attempt int, _ error) {
	}).MaxAttempts(MaxAttempts, 1).Timeout(time.Second).
		Do(context.Background(), func(ctx context.Context) Result[int, error] {
			select {
			case <-ctx.Done():
				return Err[int, error](ctx.Err())
			case <-time.After(time.Second * 3):
				return Err[int, error](io.EOF)
			}
		})
	Equal(t, result.IsErr(), true)
	Equal(t, result.Err(), context.DeadlineExceeded)
}

func TestRetrierEarlyReturn(t *testing.T) {
	var earlyReturnCount int

	r := NewRetryer[int, error]().Backoff(func(ctx context.Context, attempt int, _ error) {
	}).MaxAttempts(MaxAttempts, 5).Timeout(time.Second).
		IsEarlyReturnFn(func(ctx context.Context, err error) bool {
			earlyReturnCount++
			return errors.Is(err, io.EOF)
		}).Backoff(nil)

	result := r.Do(context.Background(), func(ctx context.Context) Result[int, error] {
		return Err[int, error](io.EOF)
	})
	Equal(t, result.IsErr(), true)
	Equal(t, result.Err(), io.EOF)
	Equal(t, earlyReturnCount, 1)

	// now let try with retryable overriding early return TL;DR retryable should take precedence over early return
	earlyReturnCount = 0
	isRetryableCount := 0
	result = r.IsRetryableFn(func(ctx context.Context, err error) (isRetryable bool) {
		isRetryableCount++
		return errors.Is(err, io.EOF)
	}).Do(context.Background(), func(ctx context.Context) Result[int, error] {
		return Err[int, error](io.EOF)
	})
	Equal(t, result.IsErr(), true)
	Equal(t, result.Err(), io.EOF)
	Equal(t, earlyReturnCount, 0)
	Equal(t, isRetryableCount, 5)

	// while here let's check the first test case again, `Retrier` should be a copy and original still intact.
	isRetryableCount = 0
	result = r.Do(context.Background(), func(ctx context.Context) Result[int, error] {
		return Err[int, error](io.EOF)
	})
	Equal(t, result.IsErr(), true)
	Equal(t, result.Err(), io.EOF)
	Equal(t, earlyReturnCount, 1)
	Equal(t, isRetryableCount, 0)
}
