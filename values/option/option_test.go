//go:build go1.18
// +build go1.18

package optionext

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	. "github.com/go-playground/assert/v2"
)

type valueTest struct {
}

func (valueTest) Value() (driver.Value, error) {
	return "value", nil
}

type customStringType string

type testStructType struct {
	Name string
}

func TestAndXXX(t *testing.T) {
	s := Some(1)
	Equal(t, Some(3), s.And(func(i int) int { return 3 }))
	Equal(t, Some(3), s.AndThen(func(i int) Option[int] { return Some(3) }))
	Equal(t, None[int](), s.AndThen(func(i int) Option[int] { return None[int]() }))

	n := None[int]()
	Equal(t, None[int](), n.And(func(i int) int { return 3 }))
	Equal(t, None[int](), n.AndThen(func(i int) Option[int] { return Some(3) }))
	Equal(t, None[int](), n.AndThen(func(i int) Option[int] { return None[int]() }))
	Equal(t, None[int](), s.AndThen(func(i int) Option[int] { return None[int]() }))
}

func TestUnwraps(t *testing.T) {
	none := None[int]()
	PanicMatches(t, func() { none.Unwrap() }, "Option.Unwrap: option is None")

	v := none.UnwrapOr(3)
	Equal(t, 3, v)

	v = none.UnwrapOrElse(func() int { return 2 })
	Equal(t, 2, v)

	v = none.UnwrapOrDefault()
	Equal(t, 0, v)

	// now test with a pointer type.
	type myStruct struct {
		S string
	}

	sNone := None[*myStruct]()
	PanicMatches(t, func() { sNone.Unwrap() }, "Option.Unwrap: option is None")

	v2 := sNone.UnwrapOr(&myStruct{S: "blah"})
	Equal(t, &myStruct{S: "blah"}, v2)

	v2 = sNone.UnwrapOrElse(func() *myStruct { return &myStruct{S: "blah 2"} })
	Equal(t, &myStruct{S: "blah 2"}, v2)

	v2 = sNone.UnwrapOrDefault()
	Equal(t, nil, v2)
}

func TestSQLDriverValue(t *testing.T) {

	var v valueTest
	Equal(t, reflect.TypeOf(v).Implements(valuerType), true)

	// none
	nOpt := None[string]()
	nVal, err := nOpt.Value()
	Equal(t, err, nil)
	Equal(t, nVal, nil)

	// string + convert custom string type
	sOpt := Some("myString")
	sVal, err := sOpt.Value()
	Equal(t, err, nil)

	_, ok := sVal.(string)
	Equal(t, ok, true)
	Equal(t, sVal, "myString")

	sCustOpt := Some(customStringType("string"))
	sCustVal, err := sCustOpt.Value()
	Equal(t, err, nil)
	Equal(t, sCustVal, "string")

	_, ok = sCustVal.(string)
	Equal(t, ok, true)

	// bool
	bOpt := Some(true)
	bVal, err := bOpt.Value()
	Equal(t, err, nil)

	_, ok = bVal.(bool)
	Equal(t, ok, true)
	Equal(t, bVal, true)

	// int64
	iOpt := Some(int64(2))
	iVal, err := iOpt.Value()
	Equal(t, err, nil)

	_, ok = iVal.(int64)
	Equal(t, ok, true)
	Equal(t, iVal, int64(2))

	// float64
	fOpt := Some(1.1)
	fVal, err := fOpt.Value()
	Equal(t, err, nil)

	_, ok = fVal.(float64)
	Equal(t, ok, true)
	Equal(t, fVal, 1.1)

	// time.Time
	dt := time.Now().UTC()
	dtOpt := Some(dt)
	dtVal, err := dtOpt.Value()
	Equal(t, err, nil)

	_, ok = dtVal.(time.Time)
	Equal(t, ok, true)
	Equal(t, dtVal, dt)

	// Slice []byte
	b := []byte("myBytes")
	bytesOpt := Some(b)
	bytesVal, err := bytesOpt.Value()
	Equal(t, err, nil)

	_, ok = bytesVal.([]byte)
	Equal(t, ok, true)
	Equal(t, bytesVal, b)

	// Slice []uint8
	b2 := []uint8("myBytes")
	bytes2Opt := Some(b2)
	bytes2Val, err := bytes2Opt.Value()
	Equal(t, err, nil)

	_, ok = bytes2Val.([]byte)
	Equal(t, ok, true)
	Equal(t, bytes2Val, b2)

	// Array []byte
	a := []byte{'1', '2', '3'}
	arrayOpt := Some(a)
	arrayVal, err := arrayOpt.Value()
	Equal(t, err, nil)

	_, ok = arrayVal.([]byte)
	Equal(t, ok, true)
	Equal(t, arrayVal, a)

	// Slice []byte
	data := []testStructType{{Name: "test"}}
	b, err = json.Marshal(data)
	Equal(t, err, nil)

	dataOpt := Some(data)
	dataVal, err := dataOpt.Value()
	Equal(t, err, nil)

	_, ok = dataVal.([]byte)
	Equal(t, ok, true)
	Equal(t, dataVal, b)

	// Map
	data2 := map[string]int{"test": 1}
	b, err = json.Marshal(data2)
	Equal(t, err, nil)

	data2Opt := Some(data2)
	data2Val, err := data2Opt.Value()
	Equal(t, err, nil)

	_, ok = data2Val.([]byte)
	Equal(t, ok, true)
	Equal(t, data2Val, b)

	// Struct
	data3 := testStructType{Name: "test"}
	b, err = json.Marshal(data3)
	Equal(t, err, nil)

	data3Opt := Some(data3)
	data3Val, err := data3Opt.Value()
	Equal(t, err, nil)

	_, ok = data3Val.([]byte)
	Equal(t, ok, true)
	Equal(t, data3Val, b)
}

type customScanner struct {
	S string
}

func (c *customScanner) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	c.S = src.(string)
	return nil
}

func TestSQLScanner(t *testing.T) {
	value := int64(123)
	var optionI64 Option[int64]
	var optionI32 Option[int32]
	var optionI16 Option[int16]
	var optionString Option[string]
	var optionBool Option[bool]
	var optionF64 Option[float64]
	var optionByte Option[byte]
	var optionTime Option[time.Time]
	var optionInterface Option[any]

	err := optionInterface.Scan(1)
	Equal(t, err, nil)
	Equal(t, optionInterface, Some(any(1)))

	err = optionInterface.Scan("blah")
	Equal(t, err, nil)
	Equal(t, optionInterface, Some(any("blah")))

	err = optionI64.Scan(value)
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

	// custom scanner scan nil
	var customNil Option[customScanner]
	err = customNil.Scan(nil)
	Equal(t, err, nil)
	Equal(t, customNil, None[customScanner]())

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

	// test custom types
	var ct Option[customStringType]
	err = ct.Scan("test")
	Equal(t, err, nil)
	Equal(t, ct, Some(customStringType("test")))
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
