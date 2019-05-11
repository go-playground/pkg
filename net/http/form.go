package httpext

import (
	"net/url"

	"github.com/go-playground/form"
)

// FormDecoder is the type used for decoding a form for use
type FormDecoder interface {
	Decode(interface{}, url.Values) error
}

var (
	// DefaultDecoder of this package, which is configurable
	DefaultDecoder FormDecoder = form.NewDecoder()
)
