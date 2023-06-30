//go:build go1.18
// +build go1.18

package timeext

import (
	"testing"
	"time"
)

func TestInstant(t *testing.T) {
	i := NewInstant()
	if i.Elapsed() < 0 {
		t.Fatalf("elapsed time should be always be monotonically increasing")
	}
	i2 := NewInstant()
	time.Sleep(time.Second)
	if i2.Since(i) <= 0 {
		t.Fatalf("time since instant should always be after")
	}
	if i.Since(i2) != 0 {
		t.Fatalf("time since instant should be zero")
	}
}
