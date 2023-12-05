//go:build go1.18
// +build go1.18

package optionext

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

var (
	scanType      = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
	byteSliceType = reflect.TypeOf(([]byte)(nil))
	valuerType    = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
	timeType      = reflect.TypeOf((*time.Time)(nil)).Elem()
	stringType    = reflect.TypeOf((*string)(nil)).Elem()
	int64Type     = reflect.TypeOf((*int64)(nil)).Elem()
	float64Type   = reflect.TypeOf((*float64)(nil)).Elem()
	boolType      = reflect.TypeOf((*bool)(nil)).Elem()
)

// Option represents a values that represents a values existence.
//
// nil is usually used on Go however this has two problems:
// 1. Checking if the return values is nil is NOT enforced and can lead to panics.
// 2. Using nil is not good enough when nil itself is a valid value.
//
// This implements the sql.Scanner interface and can be used as a sql value for reading and writing. It supports:
// - String
// - Bool
// - Uint8
// - Float64
// - Int16
// - Int32
// - Int64
// - interface{}/any
// - time.Time
// - Struct - when type is convertable to []byte and assumes JSON.
// - Slice - when type is convertable to []byte and assumes JSON.
// - Map types - when type is convertable to []byte and assumes JSON.
//
// This also implements the `json.Marshaler` and `json.Unmarshaler` interfaces. The only caveat is a None value will result
// in a JSON `null` value. there is no way to hook into the std library to make `omitempty` not produce any value at
// this time.
type Option[T any] struct {
	value  T
	isSome bool
}

// IsSome returns true if the option is not empty.
func (o Option[T]) IsSome() bool {
	return o.isSome
}

// IsNone returns true if the option is empty.
func (o Option[T]) IsNone() bool {
	return !o.isSome
}

// Unwrap returns the values if the option is not empty or panics.
func (o Option[T]) Unwrap() T {
	if o.isSome {
		return o.value
	}
	panic("Option.Unwrap: option is None")
}

// UnwrapOr returns the contained `Some` value or provided default value.
//
// Arguments passed to `UnwrapOr` are eagerly evaluated; if you are passing the result of a function call,
// look to use `UnwrapOrElse`, which can be lazily evaluated.
func (o Option[T]) UnwrapOr(value T) T {
	if o.isSome {
		return o.value
	}
	return value
}

// UnwrapOrElse returns the contained `Some` value or computes it from a provided function.
func (o Option[T]) UnwrapOrElse(fn func() T) T {
	if o.isSome {
		return o.value
	}
	return fn()
}

// UnwrapOrDefault returns the contained `Some` value or the default value of the type T.
func (o Option[T]) UnwrapOrDefault() T {
	return o.value
}

// And calls the provided function with the contained value if the option is Some, return None otherwise.
func (o Option[T]) And(fn func(T) T) Option[T] {
	if o.isSome {
		o.value = fn(o.value)
	}
	return o
}

// AndThen calls the provided function with the contained value if the option is Some, return None otherwise.
//
// This differs from `And` in that the provided function returns an Option[T] allowing changing of the Option itself if
// manipulation of a Same value can be changes to a None.
func (o Option[T]) AndThen(fn func(T) Option[T]) Option[T] {
	if o.isSome {
		return fn(o.value)
	}
	return o
}

// Some creates a new Option with the given values.
func Some[T any](value T) Option[T] {
	return Option[T]{value, true}
}

// None creates an empty Option that represents no values.
func None[T any]() Option[T] {
	return Option[T]{}
}

// MarshalJSON implements the `json.Marshaler` interface.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.isSome {
		return json.Marshal(o.value)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements the `json.Unmarshaler` interface.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && string(data[:4]) == "null" {
		*o = None[T]()
		return nil
	}
	var v T
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	*o = Some(v)
	return nil
}

// Value implements the driver.Valuer interface.
//
// This honours the `driver.Valuer` interface if the value implements it.
// It also supports custom types of the std types and treats all else as []byte/
func (o Option[T]) Value() (driver.Value, error) {
	if o.IsNone() {
		return nil, nil
	}
	value := o.Unwrap()
	val := reflect.ValueOf(value)

	if val.Type().Implements(valuerType) {
		return val.Interface().(driver.Valuer).Value()
	}
	switch val.Kind() {
	case reflect.String:
		return val.Convert(stringType).Interface(), nil
	case reflect.Bool:
		return val.Convert(boolType).Interface(), nil
	case reflect.Int64:
		return val.Convert(int64Type).Interface(), nil
	case reflect.Float64:
		return val.Convert(float64Type).Interface(), nil
	case reflect.Slice, reflect.Array:
		if val.Type().ConvertibleTo(byteSliceType) {
			return val.Convert(byteSliceType).Interface(), nil
		}
		return json.Marshal(val.Interface())
	case reflect.Struct:
		if val.CanConvert(timeType) {
			return val.Convert(timeType).Interface(), nil
		}
		return json.Marshal(val.Interface())
	case reflect.Map:
		return json.Marshal(val.Interface())
	default:
		return val.Interface(), nil
	}
}

// Scan implements the sql.Scanner interface.
func (o *Option[T]) Scan(value any) error {

	if value == nil {
		*o = None[T]()
		return nil
	}

	val := reflect.ValueOf(&o.value)

	if val.Type().Implements(scanType) {
		err := val.Interface().(sql.Scanner).Scan(value)
		if err != nil {
			return err
		}
		o.isSome = true
		return nil
	}

	val = val.Elem()

	switch val.Kind() {
	case reflect.String:
		var v sql.NullString
		if err := v.Scan(value); err != nil {
			return err
		}
		*o = Some(reflect.ValueOf(v.String).Convert(val.Type()).Interface().(T))
	case reflect.Bool:
		var v sql.NullBool
		if err := v.Scan(value); err != nil {
			return err
		}
		*o = Some(reflect.ValueOf(v.Bool).Convert(val.Type()).Interface().(T))
	case reflect.Uint8:
		var v sql.NullByte
		if err := v.Scan(value); err != nil {
			return err
		}
		*o = Some(reflect.ValueOf(v.Byte).Convert(val.Type()).Interface().(T))
	case reflect.Float64:
		var v sql.NullFloat64
		if err := v.Scan(value); err != nil {
			return err
		}
		*o = Some(reflect.ValueOf(v.Float64).Convert(val.Type()).Interface().(T))
	case reflect.Int16:
		var v sql.NullInt16
		if err := v.Scan(value); err != nil {
			return err
		}
		*o = Some(reflect.ValueOf(v.Int16).Convert(val.Type()).Interface().(T))
	case reflect.Int32:
		var v sql.NullInt32
		if err := v.Scan(value); err != nil {
			return err
		}
		*o = Some(reflect.ValueOf(v.Int32).Convert(val.Type()).Interface().(T))
	case reflect.Int64:
		var v sql.NullInt64
		if err := v.Scan(value); err != nil {
			return err
		}
		*o = Some(reflect.ValueOf(v.Int64).Convert(val.Type()).Interface().(T))
	case reflect.Interface:
		*o = Some(reflect.ValueOf(value).Convert(val.Type()).Interface().(T))
	case reflect.Struct:
		if val.CanConvert(timeType) {
			switch t := value.(type) {
			case string:
				tm, err := time.Parse(time.RFC3339Nano, t)
				if err != nil {
					return err
				}
				*o = Some(reflect.ValueOf(tm).Convert(val.Type()).Interface().(T))

			case []byte:
				tm, err := time.Parse(time.RFC3339Nano, string(t))
				if err != nil {
					return err
				}
				*o = Some(reflect.ValueOf(tm).Convert(val.Type()).Interface().(T))

			default:
				var v sql.NullTime
				if err := v.Scan(value); err != nil {
					return err
				}
				*o = Some(reflect.ValueOf(v.Time).Convert(val.Type()).Interface().(T))
			}
			return nil
		}
		fallthrough

	default:
		switch val.Kind() {
		case reflect.Struct, reflect.Slice, reflect.Map:
			v := reflect.ValueOf(value)

			if v.Type().ConvertibleTo(byteSliceType) {
				if err := json.Unmarshal(v.Convert(byteSliceType).Interface().([]byte), &o.value); err != nil {
					return err
				}
				o.isSome = true
				return nil
			}
		}
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", value, o.value)
	}
	return nil
}
