package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	appext "github.com/go-playground/pkg/v5/app"
	errorsext "github.com/go-playground/pkg/v5/errors"
	httpext "github.com/go-playground/pkg/v5/net/http"
	. "github.com/go-playground/pkg/v5/values/result"
)

// customize as desired to meet your needs including custom retryable status codes, errors etc.
var retrier = httpext.NewRetryer()

func main() {
	ctx := appext.Context().Build()

	type Test struct {
		Date time.Time
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = httpext.JSON(w, http.StatusOK, Test{Date: time.Now().UTC()})
	}))
	defer server.Close()

	// fetch response
	fn := func(ctx context.Context) Result[*http.Request, error] {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
		if err != nil {
			return Err[*http.Request, error](err)
		}
		return Ok[*http.Request, error](req)
	}

	var result Test
	err := retrier.Do(ctx, fn, &result, http.StatusOK)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %+v\n", result)

	// `Retrier` configuration is copy and so the base `Retrier` can be used and even customized for one-off requests.
	// eg for this request we only want max attempt of 2 instead of default.
	err = retrier.MaxAttempts(errorsext.MaxAttempts, 2).Do(ctx, fn, &result, http.StatusOK)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %+v\n", result)
}
