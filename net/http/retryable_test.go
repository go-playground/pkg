//go:build go1.18
// +build go1.18

package httpext

import (
	"context"
	. "github.com/go-playground/assert/v2"
	bytesext "github.com/go-playground/pkg/v5/bytes"
	errorsext "github.com/go-playground/pkg/v5/errors"
	optionext "github.com/go-playground/pkg/v5/values/option"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDoRetryable(t *testing.T) {

	ctx := context.Background()
	type response struct {
		Name string `json:"name"`
	}
	expected := "Joey Bloggs"
	var requests int

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		Equal(t, JSON(w, http.StatusOK, response{Name: expected}), nil)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	fn := func(ctx context.Context) (*http.Request, error) {
		return http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	}
	retryCount := 0
	dummyOnRetryFn := func(ctx context.Context, origErr error, reason string, attempt int) optionext.Option[error] {
		retryCount++
		return optionext.None[error]()
	}

	result := DoRetryable[response](ctx, nil, http.StatusOK, bytesext.MiB, IsRetryableStatusCode, errorsext.IsRetryableHTTP, dummyOnRetryFn, fn)
	Equal(t, result.IsErr(), false)
	Equal(t, result.Unwrap().Name, expected)
	Equal(t, retryCount, 2)
}
