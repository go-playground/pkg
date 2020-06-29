package unsafeext

import (
	"reflect"
	"unsafe"
)

// BytesToString converts an array of bytes into a string without allocating.
// The byte slice passed to this function is not to be used after this call as it's unsafe; you have been warned.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes converts an existing string into an []byte without allocating.
// The string passed to this functions is not to be used again after this call as it's unsafe; you have been warned.
func StringToBytes(s string) (b []byte) {
	strHdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sliceHdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sliceHdr.Data = strHdr.Data
	sliceHdr.Cap = strHdr.Len
	sliceHdr.Len = strHdr.Len
	return
}
