//go:build go1.18
// +build go1.18

package timeext

import (
	"time"
)

var base = time.Now()

// NanoTime returns a monotonically increasing time in nanoseconds.
func NanoTime() int64 {
	return int64(time.Since(base))
}
