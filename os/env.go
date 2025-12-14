//go:build go1.18
// +build go1.18

package osext

import (
	"os"
	"reflect"
	"strconv"

	constraintsext "github.com/go-playground/pkg/v5/constraints"
)

// EnvDefaults interface defines the supported types that can be used for environment variables conversions.
type EnvDefaults interface {
	constraintsext.Number | ~string
}

// Env retrieves the value of the environment variable named by the key.
//
// If the variable is not set, or if the conversion fails due to incorrect value,
// a default value is returned.
func Env[T EnvDefaults](key string, defaultValue T) T {
	if v, ok := os.LookupEnv(key); ok {
		ty := reflect.TypeOf(defaultValue)
		elem := reflect.New(ty).Elem()

		switch ty.Kind() {
		case reflect.String:
			elem.SetString(v)
			return elem.Interface().(T)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if i, err := strconv.ParseInt(v, 10, ty.Bits()); err == nil {
				elem.SetInt(i)
				return elem.Interface().(T)
			}
		case reflect.Uintptr, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if i, err := strconv.ParseUint(v, 10, ty.Bits()); err == nil {
				elem.SetUint(i)
				return elem.Interface().(T)
			}
		case reflect.Float32, reflect.Float64:
			if f, err := strconv.ParseFloat(v, ty.Bits()); err == nil {
				elem.SetFloat(f)
				return elem.Interface().(T)
			}
		}
	}
	return defaultValue
}
