//go:build go1.18
// +build go1.18

package httpext

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/go-playground/assert/v2"
	errorsext "github.com/go-playground/pkg/v5/errors"
	. "github.com/go-playground/pkg/v5/values/result"
)

func TestRetryer_SuccessNoRetries(t *testing.T) {
	ctx := context.Background()

	type Test struct {
		Name string
	}
	tst := Test{Name: "test"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = JSON(w, http.StatusOK, tst)
	}))
	defer server.Close()

	retryer := NewRetryer()

	result := retryer.DoResponse(ctx, func(ctx context.Context) Result[*http.Request, error] {
		req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		if err != nil {
			return Err[*http.Request, error](err)
		}
		return Ok[*http.Request, error](req)
	}, http.StatusOK)
	Equal(t, result.IsOk(), true)
	Equal(t, result.Unwrap().StatusCode, http.StatusOK)
	defer result.Unwrap().Body.Close()

	var responseResult Test
	err := retryer.Do(ctx, func(ctx context.Context) Result[*http.Request, error] {
		req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		if err != nil {
			return Err[*http.Request, error](err)
		}
		return Ok[*http.Request, error](req)
	}, &responseResult, http.StatusOK)
	Equal(t, err, nil)
	Equal(t, responseResult, tst)
}

func TestRetryer_SuccessWithRetries(t *testing.T) {
	ctx := context.Background()
	var count int

	type Test struct {
		Name string
	}
	tst := Test{Name: "test"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if count < 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			count++
			return
		}
		_ = JSON(w, http.StatusOK, tst)
	}))
	defer server.Close()

	retryer := NewRetryer().Backoff(nil)

	result := retryer.DoResponse(ctx, func(ctx context.Context) Result[*http.Request, error] {
		req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		if err != nil {
			return Err[*http.Request, error](err)
		}
		return Ok[*http.Request, error](req)
	}, http.StatusOK)
	Equal(t, result.IsOk(), true)
	Equal(t, result.Unwrap().StatusCode, http.StatusOK)
	defer result.Unwrap().Body.Close()

	count = 0 // reset count

	var responseResult Test
	err := retryer.Do(ctx, func(ctx context.Context) Result[*http.Request, error] {
		req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		if err != nil {
			return Err[*http.Request, error](err)
		}
		return Ok[*http.Request, error](req)
	}, &responseResult, http.StatusOK)
	Equal(t, err, nil)
	Equal(t, responseResult, tst)
}

func TestRetryer_FailureMaxRetries(t *testing.T) {
	ctx := context.Background()

	type Test struct {
		Name string
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	retryer := NewRetryer().Backoff(nil).MaxAttempts(errorsext.MaxAttempts, 2)

	result := retryer.DoResponse(ctx, func(ctx context.Context) Result[*http.Request, error] {
		req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		if err != nil {
			return Err[*http.Request, error](err)
		}
		return Ok[*http.Request, error](req)
	}, http.StatusOK)
	Equal(t, result.IsErr(), true)

	var responseResult Test
	err := retryer.Do(ctx, func(ctx context.Context) Result[*http.Request, error] {
		req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		if err != nil {
			return Err[*http.Request, error](err)
		}
		return Ok[*http.Request, error](req)
	}, &responseResult, http.StatusOK)
	NotEqual(t, err, nil)
}

func TestRetryer_ExtractStatusBody(t *testing.T) {
	ctx := context.Background()
	eStr := "nooooooooooooo!"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(eStr))
	}))
	defer server.Close()

	retryer := NewRetryer().MaxAttempts(errorsext.MaxAttempts, 3)

	result := retryer.DoResponse(ctx, func(ctx context.Context) Result[*http.Request, error] {
		req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		if err != nil {
			return Err[*http.Request, error](err)
		}
		return Ok[*http.Request, error](req)
	}, http.StatusOK)
	Equal(t, result.IsErr(), true)
	var esc ErrStatusCode
	Equal(t, errors.As(result.Err(), &esc), true)
	Equal(t, esc.IsRetryableStatusCode, false)
	// check the ultimate failed response body is intact
	Equal(t, string(esc.Body), eStr)
}

func TestRetryer_ExtractStatusBodyEarlyReturn(t *testing.T) {
	ctx := context.Background()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
	}))
	defer server.Close()

	var count int

	retryer := NewRetryer().Backoff(func(_ context.Context, _ int, _ error) {
		count++
	}).MaxAttempts(errorsext.MaxAttempts, 2)

	result := retryer.DoResponse(ctx, func(ctx context.Context) Result[*http.Request, error] {
		req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		if err != nil {
			return Err[*http.Request, error](err)
		}
		return Ok[*http.Request, error](req)
	}, http.StatusOK)
	Equal(t, result.IsErr(), true)
	var esc ErrStatusCode
	Equal(t, errors.As(result.Err(), &esc), true)
	Equal(t, esc.IsRetryableStatusCode, false)
	// check the ultimate failed response body is intact
	Equal(t, string(esc.Body), http.StatusText(http.StatusUnauthorized))
	Equal(t, count, 0)
}
