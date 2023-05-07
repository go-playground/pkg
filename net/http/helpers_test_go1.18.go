package httpext

import (
	. "github.com/go-playground/assert/v2"
	bytesext "github.com/go-playground/pkg/v5/bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecodeResponse(t *testing.T) {

	type result struct {
		ID int `json:"id" xml:"id"`
	}

	tests := []struct {
		name     string
		handler  http.HandlerFunc
		expected result
	}{
		{
			name: "Test JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				Equal(t, JSON(w, http.StatusOK, result{ID: 3}), nil)
			},
			expected: result{ID: 3},
		},
		{
			name: "Test XML",
			handler: func(w http.ResponseWriter, r *http.Request) {
				Equal(t, XML(w, http.StatusOK, result{ID: 5}), nil)
			},
			expected: result{ID: 5},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/", tc.handler)

			server := httptest.NewServer(mux)
			defer server.Close()

			req, err := http.NewRequest(http.MethodGet, server.URL, nil)
			Equal(t, err, nil)

			resp, err := http.DefaultClient.Do(req)
			Equal(t, err, nil)
			Equal(t, resp.StatusCode, http.StatusOK)

			res, err := DecodeResponse[result](resp, bytesext.MiB)
			Equal(t, err, nil)
			Equal(t, tc.expected.ID, res.ID)
		})
	}
}
