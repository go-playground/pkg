//go:build go1.18
// +build go1.18

package osext

import (
	"os"
	"testing"

	. "github.com/go-playground/assert/v2"
)

func TestEnv(t *testing.T) {
	key := "ENV_FOO"
	os.Setenv(key, "FOO")
	defer os.Clearenv()

	Equal(t, "FOO", Env(key, "default_value"))

	os.Unsetenv(key)
	Equal(t, "default_value", Env(key, "default_value"))
}
