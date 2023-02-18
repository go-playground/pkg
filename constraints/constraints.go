//go:build go1.18
// +build go1.18

package constraintsext

// Number represents any non-complex number eg. Integer and Float.
type Number interface {
	Integer | Float
}

// Integer represents any integer type both signed and unsigned.
type Integer interface {
	Signed | Unsigned
}

// Float represents any float type.
type Float interface {
	~float32 | ~float64
}

// Signed represents any signed integer.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned represents any unsigned integer.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Complex represents any complex number type.
type Complex interface {
	~complex64 | ~complex128
}
