//go:build go1.18

package value

// Result represents the result of an operation that is successful or not.
type Result[T any] struct {
	value T
	err   error
}

// Ok returns a Result that contains the given value.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

// Err returns a Result that contains the given error.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// IsOk returns true if the result is successful with no error.
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr returns true if the result is not successful and has an error.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Unwrap returns the value of the result. It panics if there is no result due to not checking for errors.
func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic("Result.Unwrap(): result is Err")
	}
	return r.value
}

// Err returns the error of the result. To be used after calling IsOK()
func (r Result[T]) Err() error {
	return r.err
}
