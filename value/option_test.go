//go:build go1.18

package value

import (
	"testing"

	. "github.com/go-playground/assert/v2"
)

func TestNilOption(t *testing.T) {
	value := Some[any](nil)
	Equal(t, false, value.IsNone())
	Equal(t, true, value.IsSome())
	Equal(t, nil, value.Unwrap())

	ret := returnTypedNoneOption()
	Equal(t, true, ret.IsNone())
	Equal(t, false, ret.IsSome())
	PanicMatches(t, func() {
		ret.Unwrap()
	}, "Option.Unwrap: option is None")

	ret = returnTypedSomeOption()
	Equal(t, false, ret.IsNone())
	Equal(t, true, ret.IsSome())
	Equal(t, myStruct{}, ret.Unwrap())

	retPtr := returnTypedNoneOptionPtr()
	Equal(t, true, retPtr.IsNone())
	Equal(t, false, retPtr.IsSome())

	retPtr = returnTypedSomeOptionPtr()
	Equal(t, false, retPtr.IsNone())
	Equal(t, true, retPtr.IsSome())
	Equal(t, new(myStruct), retPtr.Unwrap())
}

type myStruct struct{}

func returnTypedNoneOption() Option[myStruct] {
	return None[myStruct]()
}

func returnTypedSomeOption() Option[myStruct] {
	return Some(myStruct{})
}

func returnTypedNoneOptionPtr() Option[*myStruct] {
	return None[*myStruct]()
}

func returnTypedSomeOptionPtr() Option[*myStruct] {
	return Some(new(myStruct))
}
