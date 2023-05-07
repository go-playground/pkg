//go:build go1.18
// +build go1.18

package httpext

import (
	"errors"
	bytesext "github.com/go-playground/pkg/v5/bytes"
	"net/http"
	"strings"
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
