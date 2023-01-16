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
