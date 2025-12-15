//go:build go1.22
// +build go1.22

package osext

import (
	"os"
	"reflect"
	"strconv"

	constraintsext "github.com/go-playground/pkg/v5/constraints"
	. "github.com/go-playground/pkg/v5/values/option"
	. "github.com/go-playground/pkg/v5/values/result"
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
	return GetEnv[T](key).UnwrapOr(defaultValue)
}

// GetEnv retrieves the value of the environment variable named by the key.
//
// If the variable is not set, or if the conversion fails it returns `None`, otherwise `Some`.
func GetEnv[T EnvDefaults](key string) Option[T] {
	r := LookupEnv[T](key)
	if r.IsErr() {
		return None[T]()
	}
	return r.Unwrap()
}

// LookupEnv retrieves the value of the environment variable named by the key.
//
// If the variable is not present it returns Ok(None)
// If the variable is present and conversion is successful, the value Ok(Some) is returned
// If the variable is present and conversion fails it returns Err(error)
func LookupEnv[T EnvDefaults](key string) Result[Option[T], error] {
	if v, ok := os.LookupEnv(key); ok {
		ty := reflect.TypeFor[T]()
		elem := reflect.New(ty).Elem()

		switch ty.Kind() {
		case reflect.String:
			elem.SetString(v)
			return Ok[Option[T], error](Some(elem.Interface().(T)))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(v, 10, ty.Bits())
			if err != nil {
				return Err[Option[T], error](err)
			}
			elem.SetInt(i)
			return Ok[Option[T], error](Some(elem.Interface().(T)))

		case reflect.Uintptr, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i, err := strconv.ParseUint(v, 10, ty.Bits())
			if err != nil {
				return Err[Option[T], error](err)
			}
			elem.SetUint(i)
			return Ok[Option[T], error](Some(elem.Interface().(T)))
		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(v, ty.Bits())
			if err != nil {
				return Err[Option[T], error](err)
			}
			elem.SetFloat(f)
			return Ok[Option[T], error](Some(elem.Interface().(T)))
		}
	}
	return Ok[Option[T], error](None[T]())
}
