package mathext

import (
	"math"
)

// Min returns the smaller value.
//
// NOTE: this function does not check for difference in floats of 0/zero vs -0/negative zero using Signbit.
func Min[N number](x, y N) N {
	// special case for float
	// IEEE 754 says that only NaNs satisfy f != f.
	if x != x || y != y {
		return N(math.NaN())
	}

	if x < y {
		return x
	}
	return y
}
