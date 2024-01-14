//go:build go1.18
// +build go1.18

package resultext

// Result represents the result of an operation that is successful or not.
type Result[T, E any] struct {
	ok   T
	err  E
	isOk bool
}

// Ok returns a Result that contains the given values.
func Ok[T, E any](value T) Result[T, E] {
	return Result[T, E]{ok: value, isOk: true}
}

// Err returns a Result that contains the given error.
func Err[T, E any](err E) Result[T, E] {
	return Result[T, E]{err: err}
}

// IsOk returns true if the result is successful with no error.
func (r Result[T, E]) IsOk() bool {
	return r.isOk
}

// IsErr returns true if the result is not successful and has an error.
func (r Result[T, E]) IsErr() bool {
	return !r.isOk
}

// Unwrap returns the values of the result. It panics if there is no result due to not checking for errors.
func (r Result[T, E]) Unwrap() T {
	if r.isOk {
		return r.ok
	}
	panic("Result.Unwrap(): result is Err")
}

// UnwrapOr returns the contained Ok value or a provided default.
//
// Arguments passed to UnwrapOr are eagerly evaluated; if you are passing the result of a function call,
// look to use `UnwrapOrElse`, which can be lazily evaluated.
func (r Result[T, E]) UnwrapOr(value T) T {
	if r.isOk {
		return r.ok
	}
	return value
}

// UnwrapOrElse returns the contained Ok value or computes it from a provided function.
func (r Result[T, E]) UnwrapOrElse(fn func() T) T {
	if r.isOk {
		return r.ok
	}
	return fn()
}

// UnwrapOrDefault returns the contained Ok value or the default value of the type T.
func (r Result[T, E]) UnwrapOrDefault() T {
	return r.ok
}

// And calls the provided function with the contained value if the result is Ok, returns the Result value otherwise.
func (r Result[T, E]) And(fn func(T) T) Result[T, E] {
	if r.isOk {
		r.ok = fn(r.ok)
	}
	return r
}

// AndThen calls the provided function with the contained value if the result is Ok, returns the Result value otherwise.
//
// This differs from `And` in that the provided function returns a Result[T, E] allowing changing of the Option value
// itself.
func (r Result[T, E]) AndThen(fn func(T) Result[T, E]) Result[T, E] {
	if r.isOk {
		return fn(r.ok)
	}
	return r
}

// Err returns the error of the result. To be used after calling IsOK()
func (r Result[T, E]) Err() E {
	return r.err
}
