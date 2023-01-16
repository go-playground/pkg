package mapext

import (
	. "github.com/go-playground/assert/v2"
	"testing"
)

func TestRetain(t *testing.T) {
	m := map[string]int{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
	}
	Retain(m, func(key string, value int) bool {
		return value < 1 || value > 2
	})
	Equal(t, len(m), 2)
	Equal(t, m["0"], 0)
	Equal(t, m["3"], 3)
}

func TestMap(t *testing.T) {
	// Test Map to slice
	m := map[string]int{
		"0": 0,
		"1": 1,
	}
	slice := Map(m, make([]int, 0, len(m)), func(accum []int, key string, value int) []int {
		return append(accum, value)
	})
	Equal(t, len(slice), 2)
	Equal(t, slice[0], 0)
	Equal(t, slice[1], 1)

	// Test Map to Map of different type
	inverted := Map(m, make(map[int]string, len(m)), func(accum map[int]string, key string, value int) map[int]string {
		accum[value] = key
		return accum
	})
	Equal(t, len(inverted), 2)
	Equal(t, inverted[0], "0")
	Equal(t, inverted[1], "1")
}
