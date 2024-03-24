package httpext

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/go-playground/assert/v2"
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
	fmt.Println(result.Err())
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
