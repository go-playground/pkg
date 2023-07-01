//go:build go1.18
// +build go1.18

package timeext

import "time"

// Instant represents a monotonic instant in time.
//
// Instants are opaque types that can only be compared with one another and allows measuring of duration.
type Instant int64

// NewInstant returns a new Instant.
func NewInstant() Instant {
	return Instant(NanoTime())
}

// Elapsed returns the duration since the instant was created.
func (i Instant) Elapsed() time.Duration {
	return time.Duration(Instant(NanoTime()) - i)
}

// Since returns the duration elapsed from another Instant, or zero is that Instant is later than this one.
func (i Instant) Since(instant Instant) time.Duration {
	if instant > i {
		return 0
	}
	return time.Duration(i - instant)
}
