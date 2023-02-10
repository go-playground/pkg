package mathext

import (
	"math"
	"testing"

	. "github.com/go-playground/assert/v2"
)

func TestMin(t *testing.T) {
	Equal(t, true, math.IsNaN(Min(math.NaN(), 1)))
	Equal(t, true, math.IsNaN(Min(1, math.NaN())))
	Equal(t, math.Inf(-1), Min(math.Inf(0), math.Inf(-1)))
	Equal(t, math.Inf(-1), Min(math.Inf(-1), math.Inf(0)))
	Equal(t, 1.0, Min(1.333, 1.0))
	Equal(t, 1.0, Min(1.0, 1.333))
	Equal(t, 1, Min(3, 1))
	Equal(t, 1, Min(1, 3))
	Equal(t, -0, Min(0, -0))
	Equal(t, -0, Min(-0, 0))
}

func BenchmarkMinInf(b *testing.B) {
	n1 := math.Inf(0)
	n2 := math.Inf(-1)

	for i := 0; i < b.N; i++ {
		_ = Min(n1, n2)
	}
}

func BenchmarkMinNaN(b *testing.B) {
	n1 := math.Inf(0)
	n2 := math.NaN()

	for i := 0; i < b.N; i++ {
		_ = Min(n1, n2)
	}
}

func BenchmarkMinNumber(b *testing.B) {
	n1 := 1
	n2 := 3

	for i := 0; i < b.N; i++ {
		_ = Min(n1, n2)
	}
}
