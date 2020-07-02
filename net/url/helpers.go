package urlext

import (
	"net/url"

	httpext "github.com/go-playground/pkg/v5/net/http"
)

// EncodeToURLValues encodes a struct or field into a set of url.Values
func EncodeToURLValues(v interface{}) (url.Values, error) {
	return httpext.DefaultFormEncoder.Encode(v)
}
