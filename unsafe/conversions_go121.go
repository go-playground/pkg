//go:build go1.21

package unsafeext

import (
	"unsafe"
)

// BytesToString converts an array of bytes into a string without allocating.
// The byte slice passed to this function is not to be used after this call as it's unsafe; you have been warned.
func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// StringToBytes converts an existing string into an []byte without allocating.
// The string passed to these functions is not to be used again after this call as it's unsafe; you have been warned.
func StringToBytes(s string) (b []byte) {
	d := unsafe.StringData(s)
	b = unsafe.Slice(d, len(s))
	return
}
