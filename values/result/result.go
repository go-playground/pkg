//go:build go1.18
// +build go1.18

package result

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
	return !r.IsOk()
}

// Unwrap returns the values of the result. It panics if there is no result due to not checking for errors.
func (r Result[T, E]) Unwrap() T {
	if !r.isOk {
		panic("Result.Unwrap(): result is Err")
	}
	return r.ok
}

// Err returns the error of the result. To be used after calling IsOK()
func (r Result[T, E]) Err() E {
	return r.err
}
