//go:build go1.21
// +build go1.21

package mathext

import (
	"cmp"
)

// Max returns the larger value.
//
// NOTE: this function does not check for difference in floats of 0/zero vs -0/negative zero using Signbit.
func Max[N cmp.Ordered](x, y N) N {
	return max(x, y)
}
