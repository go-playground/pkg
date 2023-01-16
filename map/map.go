//go:build go1.18
// +build go1.18

package mapext

// Retain retains only the elements specified by the function and removes others.
func Retain[K comparable, V any](m map[K]V, fn func(key K, value V) bool) {
	for k, v := range m {
		if fn(k, v) {
			continue
		}
		delete(m, k)
	}
}

// Map allows mapping of a map[K]V -> U.
func Map[K comparable, V any, U any](m map[K]V, init U, fn func(accum U, key K, value V) U) U {
	accum := init
	for k, v := range m {
		accum = fn(accum, k, v)
	}
	return accum
}
