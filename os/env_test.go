//go:build go1.18
// +build go1.18

package osext

import (
	"fmt"
	"math/rand"
	"testing"
)

type (
	CustomInt     int
	CustomInt8    int8
	CustomInt16   int16
	CustomInt32   int32
	CustomInt64   int64
	CustomUint    uint
	CustomUint8   uint8
	CustomUint16  uint16
	CustomUint32  uint32
	CustomUint64  uint64
	CustomUintptr uintptr
	CustomFloat32 float32
	CustomFloat64 float64
	CustomString  string
)

func TestEnv(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() error
	}{
		// Integer types
		{"int", func() error { return SetAndUnsetTest(t, "24", 24, 42) }},
		{"CustomInt", func() error { return SetAndUnsetTest(t, "24", CustomInt(24), CustomInt(42)) }},
		{"int8", func() error { return SetAndUnsetTest(t, "24", int8(24), int8(42)) }},
		{"CustomInt8", func() error { return SetAndUnsetTest(t, "24", CustomInt8(24), CustomInt8(42)) }},
		{"int16", func() error { return SetAndUnsetTest(t, "24", int16(24), int16(42)) }},
		{"CustomInt16", func() error { return SetAndUnsetTest(t, "24", CustomInt16(24), CustomInt16(42)) }},
		{"int32", func() error { return SetAndUnsetTest(t, "24", int32(24), int32(42)) }},
		{"CustomInt32", func() error { return SetAndUnsetTest(t, "24", CustomInt32(24), CustomInt32(42)) }},
		{"int64", func() error { return SetAndUnsetTest(t, "24", int64(24), int64(42)) }},
		{"CustomInt64", func() error { return SetAndUnsetTest(t, "24", CustomInt64(24), CustomInt64(42)) }},
		{"uint", func() error { return SetAndUnsetTest(t, "24", uint(24), uint(42)) }},
		{"CustomUint", func() error { return SetAndUnsetTest(t, "24", CustomUint(24), CustomUint(42)) }},
		{"uint8", func() error { return SetAndUnsetTest(t, "24", uint8(24), uint8(42)) }},
		{"CustomUint8", func() error { return SetAndUnsetTest(t, "24", CustomUint8(24), CustomUint8(42)) }},
		{"uint16", func() error { return SetAndUnsetTest(t, "24", uint16(24), uint16(42)) }},
		{"CustomUint16", func() error { return SetAndUnsetTest(t, "24", CustomUint16(24), CustomUint16(42)) }},
		{"uint32", func() error { return SetAndUnsetTest(t, "24", uint32(24), uint32(42)) }},
		{"CustomUint32", func() error { return SetAndUnsetTest(t, "24", CustomUint32(24), CustomUint32(42)) }},
		{"uint64", func() error { return SetAndUnsetTest(t, "24", uint64(24), uint64(42)) }},
		{"CustomUint64", func() error { return SetAndUnsetTest(t, "24", CustomUint64(24), CustomUint64(42)) }},
		{"uintptr", func() error { return SetAndUnsetTest(t, "24", uintptr(24), uintptr(42)) }},
		{"CustomUintptr", func() error { return SetAndUnsetTest(t, "24", CustomUintptr(24), CustomUintptr(42)) }},

		// Float types
		{"float32", func() error { return SetAndUnsetTest(t, "24.5", float32(24.5), float32(42.5)) }},
		{"CustomFloat32", func() error { return SetAndUnsetTest(t, "24.5", CustomFloat32(24.5), CustomFloat32(42.5)) }},
		{"float64", func() error { return SetAndUnsetTest(t, "24.5", float64(24.5), float64(42.5)) }},
		{"CustomFloat64", func() error { return SetAndUnsetTest(t, "24.5", CustomFloat64(24.5), CustomFloat64(42.5)) }},

		// String types
		{"string", func() error { return SetAndUnsetTest(t, "hello", "hello", "world") }},
		{"CustomString", func() error { return SetAndUnsetTest(t, "hello", CustomString("hello"), CustomString("world")) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.testFunc(); err != nil {
				t.Errorf("Test %s failed: %v", tt.name, err)
			}
		})
	}
}

func SetAndUnsetTest[T EnvDefaults](t *testing.T, envValueSet string, expectedSet T, expectedDefault T) error {
	constTestEnvKey := fmt.Sprintf("TEST_ENV_%d", rand.Intn(1000000))

	v := EnvOrDefault(constTestEnvKey, expectedDefault)
	if v != expectedDefault {
		return fmt.Errorf("default value mismatch: got %v, want %v", v, expectedDefault)
	}

	t.Setenv(constTestEnvKey, envValueSet)

	v = EnvOrDefault(constTestEnvKey, expectedDefault)
	if v != expectedSet {
		return fmt.Errorf("set value mismatch: got %v, want %v", v, expectedSet)
	}
	return nil
}
