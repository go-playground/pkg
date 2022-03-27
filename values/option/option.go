//go:build go1.18
// +build go1.18

package option

// Option represents a values that represents a values existence.
//
// nil is usually used on Go however this has two problems:
// 1. Checking if the return values is nil is NOT enforced and can lead to panics.
// 2. Using nil is not good enough when nil itself is a valid values.
//
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
	if o.IsNone() {
		panic("Option.Unwrap: option is None")
	}
	return o.value
}

// Some creates a new Option with the given values.
func Some[T any](value T) Option[T] {
	return Option[T]{value, true}
}

// None creates an empty Option that represents no values.
func None[T any]() Option[T] {
	return Option[T]{}
}
