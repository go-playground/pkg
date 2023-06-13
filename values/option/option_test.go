//go:build go1.18
// +build go1.18

package optionext

import (
	"encoding/json"
	"testing"
	"time"

	. "github.com/go-playground/assert/v2"
)

type customScanner struct {
	S string
}

func (c *customScanner) Scan(src interface{}) error {
	c.S = src.(string)
	return nil
}

func TestSQL(t *testing.T) {
	value := int64(123)
	var optionI64 Option[int64]
	var optionI32 Option[int32]
	var optionI16 Option[int16]
	var optionString Option[string]
	var optionBool Option[bool]
	var optionF64 Option[float64]
	var optionByte Option[byte]
	var optionTime Option[time.Time]

	err := optionI64.Scan(value)
	Equal(t, err, nil)
	Equal(t, optionI64, Some(value))

	err = optionI32.Scan(value)
	Equal(t, err, nil)
	Equal(t, optionI32, Some(int32(value)))

	err = optionI16.Scan(value)
	Equal(t, err, nil)
	Equal(t, optionI16, Some(int16(value)))

	err = optionBool.Scan(1)
	Equal(t, err, nil)
	Equal(t, optionBool, Some(true))

	err = optionString.Scan(value)
	Equal(t, err, nil)
	Equal(t, optionString, Some("123"))

	err = optionF64.Scan(2.0)
	Equal(t, err, nil)
	Equal(t, optionF64, Some(2.0))

	err = optionByte.Scan(uint8('1'))
	Equal(t, err, nil)
	Equal(t, optionByte, Some(uint8('1')))

	err = optionTime.Scan("2023-06-13T06:34:32Z")
	Equal(t, err, nil)
	Equal(t, optionTime, Some(time.Date(2023, 6, 13, 6, 34, 32, 0, time.UTC)))

	err = optionTime.Scan([]byte("2023-06-13T06:34:32Z"))
	Equal(t, err, nil)
	Equal(t, optionTime, Some(time.Date(2023, 6, 13, 6, 34, 32, 0, time.UTC)))

	err = optionTime.Scan(time.Date(2023, 6, 13, 6, 34, 32, 0, time.UTC))
	Equal(t, err, nil)
	Equal(t, optionTime, Some(time.Date(2023, 6, 13, 6, 34, 32, 0, time.UTC)))

	// Test nil
	var nullableOption Option[int64]
	err = nullableOption.Scan(nil)
	Equal(t, err, nil)
	Equal(t, nullableOption, None[int64]())

	// custom scanner
	var custom Option[customScanner]
	err = custom.Scan("GOT HERE")
	Equal(t, err, nil)
	Equal(t, custom, Some(customScanner{S: "GOT HERE"}))

	// test unmarshal to struct
	type myStruct struct {
		Name string `json:"name"`
	}

	var optionMyStruct Option[myStruct]
	err = optionMyStruct.Scan([]byte(`{"name":"test"}`))
	Equal(t, err, nil)
	Equal(t, optionMyStruct, Some(myStruct{Name: "test"}))

	err = optionMyStruct.Scan(json.RawMessage(`{"name":"test2"}`))
	Equal(t, err, nil)
	Equal(t, optionMyStruct, Some(myStruct{Name: "test2"}))

	var optionArrayOfMyStruct Option[[]myStruct]
	err = optionArrayOfMyStruct.Scan([]byte(`[{"name":"test"}]`))
	Equal(t, err, nil)
	Equal(t, optionArrayOfMyStruct, Some([]myStruct{{Name: "test"}}))

	var optionMap Option[map[string]any]
	err = optionMap.Scan([]byte(`{"name":"test"}`))
	Equal(t, err, nil)
	Equal(t, optionMap, Some(map[string]any{"name": "test"}))
}

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
