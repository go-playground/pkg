package urlext

import (
	"testing"

	. "github.com/go-playground/assert/v2"
)

func TestEncodeToURLValues(t *testing.T) {
	type Test struct {
		Domain string `form:"domain"`
		Next   string `form:"next"`
	}

	s := Test{Domain: "company.org", Next: "NIDEJ89#(@#NWJK"}
	values, err := EncodeToURLValues(s)
	Equal(t, err, nil)
	Equal(t, len(values), 2)
	Equal(t, values.Encode(), "domain=company.org&next=NIDEJ89%23%28%40%23NWJK")
}
