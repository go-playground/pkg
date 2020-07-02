package httpext

import (
	"net/url"

	"github.com/go-playground/form/v4"
)

// FormDecoder is the type used for decoding a form for use
type FormDecoder interface {
	Decode(interface{}, url.Values) error
}

// FormEncoder is the type used for encoding form data
type FormEncoder interface {
	Encode(interface{}) (url.Values, error)
}

var (
	// DefaultFormDecoder of this package, which is configurable
	DefaultFormDecoder FormDecoder = form.NewDecoder()

	// DefaultFormEncoder of this package, which is configurable
	DefaultFormEncoder FormEncoder = form.NewEncoder()
)
