package errorsext

import (
	"context"
	"io"
	"testing"

	. "github.com/go-playground/assert/v2"
	. "github.com/go-playground/pkg/v5/values/result"
)

// TODO: Add IsRetryable and Retryable to helper functions.

func TestRetrierMaxAttempts(t *testing.T) {
	var i, j int
	result := NewRetryer[int, error]().Backoff(func(ctx context.Context, attempt int) {
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
	result := NewRetryer[int, error]().IsRetryableFn(func(e error) (isRetryable bool) {
		if returnErr == io.EOF {
			return false
		} else {
			return true
		}
	}).Backoff(func(ctx context.Context, attempt int) {
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
	result := NewRetryer[int, error]().IsRetryableFn(func(e error) (isRetryable bool) {
		if returnErr == io.EOF {
			return false
		} else {
			return true
		}
	}).Backoff(func(ctx context.Context, attempt int) {
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
	r := NewRetryer[int, error]().Backoff(func(ctx context.Context, attempt int) {
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
