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

// EnvOrDefault retrieves the value of the environment variable named by the key.
//
// If the variable is not set, or if the conversion fails due to incorrect value,
// a default value is returned.
func EnvOrDefault[T EnvDefaults](key string, defaultValue T) T {
	if v, ok := os.LookupEnv(key); ok {

		rv := reflect.ValueOf(defaultValue)
		ty := rv.Type()

		switch ty.Kind() {
		case reflect.String:
			return reflect.ValueOf(v).Convert(ty).Interface().(T)
		case reflect.Int:
			i, err := strconv.ParseInt(v, 10, 0)
			if err == nil {
				return reflect.ValueOf(i).Convert(ty).Interface().(T)
			}
		case reflect.Int8:
			i, err := strconv.ParseInt(v, 10, 8)
			if err == nil {
				return reflect.ValueOf(i).Convert(ty).Interface().(T)
			}
		case reflect.Int16:
			i, err := strconv.ParseInt(v, 10, 16)
			if err == nil {
				return reflect.ValueOf(i).Convert(ty).Interface().(T)
			}
		case reflect.Int32:
			i, err := strconv.ParseInt(v, 10, 32)
			if err == nil {
				return reflect.ValueOf(i).Convert(ty).Interface().(T)
			}
		case reflect.Int64:
			i, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				return reflect.ValueOf(i).Convert(ty).Interface().(T)
			}
		case reflect.Uint, reflect.Uintptr:
			i, err := strconv.ParseUint(v, 10, 0)
			if err == nil {
				return reflect.ValueOf(i).Convert(ty).Interface().(T)
			}
		case reflect.Uint8:
			i, err := strconv.ParseUint(v, 10, 8)
			if err == nil {
				return reflect.ValueOf(i).Convert(ty).Interface().(T)
			}
		case reflect.Uint16:
			i, err := strconv.ParseUint(v, 10, 16)
			if err == nil {
				return reflect.ValueOf(i).Convert(ty).Interface().(T)
			}
		case reflect.Uint32:
			i, err := strconv.ParseUint(v, 10, 32)
			if err == nil {
				return reflect.ValueOf(i).Convert(ty).Interface().(T)
			}
		case reflect.Uint64:
			i, err := strconv.ParseUint(v, 10, 64)
			if err == nil {
				return reflect.ValueOf(i).Convert(ty).Interface().(T)
			}
		case reflect.Float32:
			f, err := strconv.ParseFloat(v, 32)
			if err == nil {
				return reflect.ValueOf(f).Convert(ty).Interface().(T)
			}
		case reflect.Float64:
			f, err := strconv.ParseFloat(v, 64)
			if err == nil {
				return reflect.ValueOf(f).Convert(ty).Interface().(T)
			}
		}
	}
	return defaultValue
}
