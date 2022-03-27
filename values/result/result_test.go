//go:build go1.18
// +build go1.18

package result

import (
	"errors"
	"testing"

	. "github.com/go-playground/assert/v2"
)

type myStruct struct{}

func TestResult(t *testing.T) {
	result := returnOk()
	Equal(t, true, result.IsOk())
	Equal(t, false, result.IsErr())
	Equal(t, true, result.Err() == nil)
	Equal(t, myStruct{}, result.Unwrap())

	result = returnErr()
	Equal(t, false, result.IsOk())
	Equal(t, true, result.IsErr())
	Equal(t, false, result.Err() == nil)
	PanicMatches(t, func() {
		result.Unwrap()
	}, "Result.Unwrap(): result is Err")
}

func returnOk() Result[myStruct] {
	return Ok(myStruct{})
}

func returnErr() Result[myStruct] {
	return Err[myStruct](errors.New("bad"))
}
