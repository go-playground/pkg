//go:build go1.18
// +build go1.18

package optionext

import (
	"encoding/json"
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

// And calls the provided function with the contained value if the option is Some, returns the None value otherwise.
func (o Option[T]) And(fn func(T) T) Option[T] {
	if o.isSome {
		o.value = fn(o.value)
	}
	return o
}

// AndThen calls the provided function with the contained value if the option is Some, returns the None value otherwise.
//
// This differs from `And` in that the provided function returns an Option[T] allowing changing of the Option value
// itself.
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
