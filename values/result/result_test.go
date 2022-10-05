//go:build go1.18
// +build go1.18

package resultext

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

func returnOk() Result[myStruct, error] {
	return Ok[myStruct, error](myStruct{})
}

func returnErr() Result[myStruct, error] {
	return Err[myStruct, error](errors.New("bad"))
}

func BenchmarkResultOk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := returnOk()
		if res.IsOk() {
			_ = res.Unwrap()
		}
	}
}

func BenchmarkResultErr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := returnErr()
		if res.IsOk() {
			_ = res.Unwrap()
		}
	}
}

func BenchmarkNoResultOk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res, err := returnOkNoResult()
		if err != nil {
			_ = res
		}
	}
}

func BenchmarkNoResultErr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res, err := returnErrNoResult()
		if err != nil {
			_ = res
		}
	}
}

func returnOkNoResult() (myStruct, error) {
	return myStruct{}, nil
}

func returnErrNoResult() (myStruct, error) {
	return myStruct{}, errors.New("bad")
}
