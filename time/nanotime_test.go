//go:build go1.18
// +build go1.18

package timeext

import (
	"testing"
	"time"
)

func TestNanoTime(t *testing.T) {
	t1 := NanoTime()
	time.Sleep(time.Second)
	t2 := NanoTime()
	if t2-t1 < int64(time.Second) {
		t.Fatalf("nanotime failed to monotonically increase, t1: %d t2: %d", t1, t2)
	}
}

func BenchmarkNanoTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NanoTime()
	}
}

func BenchmarkNanoTimeUsingUnixNano(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = time.Now().UnixNano()
	}
}
