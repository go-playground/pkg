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
	Retain(m, func(entry Entry[string, int]) bool {
		return entry.Value < 1 || entry.Value > 2
	})
	Equal(t, len(m), 2)
	Equal(t, m["0"], 0)
	Equal(t, m["3"], 3)
}
