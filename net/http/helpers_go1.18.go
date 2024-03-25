//go:build go1.18
// +build go1.18

package httpext

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	asciiext "github.com/go-playground/pkg/v5/ascii"
	bytesext "github.com/go-playground/pkg/v5/bytes"
	. "github.com/go-playground/pkg/v5/values/option"
)

// DecodeResponse takes the response and attempts to discover its content type via
// the http headers and then decode the request body into the provided type.
//
// Example if header was "application/json" would decode using
// json.NewDecoder(ioext.LimitReader(r.Body, maxMemory)).Decode(v).
func DecodeResponse[T any](r *http.Response, maxMemory bytesext.Bytes) (result T, err error) {
	typ := r.Header.Get(ContentType)
	if idx := strings.Index(typ, ";"); idx != -1 {
		typ = typ[:idx]
	}
	switch typ {
	case nakedApplicationJSON:
		err = decodeJSON(r.Header, r.Body, NoQueryParams, nil, maxMemory, &result)
	case nakedApplicationXML:
		err = decodeXML(r.Header, r.Body, NoQueryParams, nil, maxMemory, &result)
	default:
		err = errors.New("unsupported content type")
	}
	return
}

// HasRetryAfter parses the Retry-After header and returns the duration if possible.
func HasRetryAfter(headers http.Header) Option[time.Duration] {
	if ra := headers.Get(RetryAfter); ra != "" {
		if asciiext.IsDigit(ra[0]) {
			if n, err := strconv.ParseInt(ra, 10, 64); err == nil {
				return Some(time.Duration(n) * time.Second)
			}
		} else {
			// not a number so must be a date in the future
			if t, err := http.ParseTime(ra); err == nil {
				return Some(time.Until(t))
			}
		}
	}
	return None[time.Duration]()
}
