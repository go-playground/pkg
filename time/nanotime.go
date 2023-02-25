//go:build go1.18
// +build go1.18

package timeext

import (
	_ "unsafe"
)

//go:noescape
//go:linkname nanotime runtime.nanotime
func nanotime() int64

// NanoTime returns the time from the monotonic clock in nanoseconds.
func NanoTime() int64 {
	return nanotime()
}
