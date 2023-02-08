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

// Option represents a values that represents a values existence.
//
// nil is usually used on Go however this has two problems:
// 1. Checking if the return values is nil is NOT enforced and can lead to panics.
// 2. Using nil is not good enough when nil itself is a valid value.
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

// Some creates a new Option with the given values.
func Some[T any](value T) Option[T] {
	return Option[T]{value, true}
}

// None creates an empty Option that represents no values.
func None[T any]() Option[T] {
	return Option[T]{}
}

// MarshalJSON implements the json.Marshaler interface.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.isSome {
		return json.Marshal(o.value)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
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
func (o Option[T]) Value() (driver.Value, error) {
	if o.isSome {
		return o.Unwrap(), nil
	}
	return nil, nil
}

// Scan implements the sql.Scanner interface.
func (o *Option[T]) Scan(value any) error {
	val := reflect.ValueOf(o.value)
	switch val.Kind() {
	case reflect.String:
		var v sql.NullString
		if err := v.Scan(value); err != nil {
			return err
		}
		if !v.Valid {
			*o = None[T]()
		} else {
			*o = Some(reflect.ValueOf(v.String).Interface().(T))
		}
	case reflect.Bool:
		var v sql.NullBool
		if err := v.Scan(value); err != nil {
			return err
		}
		if !v.Valid {
			*o = None[T]()
		} else {
			*o = Some(reflect.ValueOf(v.Bool).Interface().(T))
		}
	case reflect.Uint8:
		var v sql.NullByte
		if err := v.Scan(value); err != nil {
			return err
		}
		if !v.Valid {
			*o = None[T]()
		} else {
			*o = Some(reflect.ValueOf(v.Byte).Interface().(T))
		}
	case reflect.Float64:
		var v sql.NullFloat64
		if err := v.Scan(value); err != nil {
			return err
		}
		if !v.Valid {
			*o = None[T]()
		} else {
			*o = Some(reflect.ValueOf(v.Float64).Interface().(T))
		}
	case reflect.Int16:
		var v sql.NullInt16
		if err := v.Scan(value); err != nil {
			return err
		}
		if !v.Valid {
			*o = None[T]()
		} else {
			*o = Some(reflect.ValueOf(v.Int16).Interface().(T))
		}
	case reflect.Int32:
		var v sql.NullInt32
		if err := v.Scan(value); err != nil {
			return err
		}
		if !v.Valid {
			*o = None[T]()
		} else {
			*o = Some(reflect.ValueOf(v.Int32).Interface().(T))
		}
	case reflect.Int64:
		var v sql.NullInt64
		if err := v.Scan(value); err != nil {
			return err
		}
		if !v.Valid {
			*o = None[T]()
		} else {
			*o = Some(reflect.ValueOf(v.Int64).Interface().(T))
		}
	case reflect.Struct:
		if val.Type() == reflect.TypeOf(time.Time{}) {
			var v sql.NullTime
			if err := v.Scan(value); err != nil {
				return err
			}
			if !v.Valid {
				*o = None[T]()
			} else {
				*o = Some(reflect.ValueOf(v.Time).Interface().(T))
			}
			return nil
		}
		fallthrough
	default:
		return fmt.Errorf("unsupported Scan, storing driver.value type %T into type %T", value, o.value)
	}
	return nil
}
