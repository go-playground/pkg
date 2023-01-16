package mapext

// Entry represents a single Map entry.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Retain retains only the elements specified by the function and removes others.
func Retain[K comparable, V any](m map[K]V, fn func(entry Entry[K, V]) bool) {
	for k, v := range m {
		if fn(Entry[K, V]{Key: k, Value: v}) {
			continue
		}
		delete(m, k)
	}
}
