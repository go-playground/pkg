package httpext

import (
	"net/url"

	"github.com/go-playground/form/v4"
)

// FormDecoder is the type used for decoding a form for use
type FormDecoder interface {
	Decode(interface{}, url.Values) error
}

var (
	// DefaultFormDecoder of this package, which is configurable
	DefaultFormDecoder FormDecoder = form.NewDecoder()
)
