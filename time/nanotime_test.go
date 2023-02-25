//go:build go1.18
// +build go1.18

package timeext

import (
	"testing"
)

func TestNanoTime(t *testing.T) {
	t1 := NanoTime()
	t2 := NanoTime()
	if t1 >= t2 {
		t.Fatalf("nanotime failed to monotonically increase, t1: %d t2: %d", t1, t2)
	}
}
