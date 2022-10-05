//go:build go1.18
// +build go1.18

package optionext

import (
	"encoding/json"
	"testing"
	"time"

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

func TestOptionJSON(t *testing.T) {
	type s struct {
		Timestamp Option[time.Time] `json:"ts"`
	}
	now := time.Now().UTC().Truncate(time.Minute)
	tv := s{Timestamp: Some(now)}

	b, err := json.Marshal(tv)
	Equal(t, nil, err)
	Equal(t, `{"ts":"`+now.Format(time.RFC3339)+`"}`, string(b))

	tv = s{}
	b, err = json.Marshal(tv)
	Equal(t, nil, err)
	Equal(t, `{"ts":null}`, string(b))
}

func TestOptionJSONOmitempty(t *testing.T) {
	type s struct {
		Timestamp Option[time.Time] `json:"ts,omitempty"`
	}
	now := time.Now().UTC().Truncate(time.Minute)
	tv := s{Timestamp: Some(now)}

	b, err := json.Marshal(tv)
	Equal(t, nil, err)
	Equal(t, `{"ts":"`+now.Format(time.RFC3339)+`"}`, string(b))

	type s2 struct {
		Timestamp *Option[time.Time] `json:"ts,omitempty"`
	}
	tv2 := &s2{}
	b, err = json.Marshal(tv2)
	Equal(t, nil, err)
	Equal(t, `{}`, string(b))
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

func BenchmarkOption(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opt := returnTypedSomeOption()
		if opt.IsSome() {
			_ = opt.Unwrap()
		}
	}
}

func BenchmarkOptionPtr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opt := returnTypedSomeOptionPtr()
		if opt.IsSome() {
			_ = opt.Unwrap()
		}
	}
}

func BenchmarkNoOptionPtr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := returnTypedNoOption()
		if result != nil {
			_ = result
		}
	}
}

func BenchmarkOptionNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opt := returnTypedSomeOptionNil()
		if opt.IsSome() {
			_ = opt.Unwrap()
		}
	}
}

func BenchmarkNoOptionNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result, found := returnNoOptionNil()
		if found {
			_ = result
		}
	}
}

func returnTypedSomeOptionNil() Option[any] {
	return Some[any](nil)
}

func returnTypedNoOption() *myStruct {
	return new(myStruct)
}

func returnNoOptionNil() (any, bool) {
	return nil, true
}
