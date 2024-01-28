package stringsext

import "strings"

// Join is a wrapper around strings.Join with a more ergonomic interface when you don't already have a slice of strings.
//
// Join concatenates the variadic elements placing the separator string between each element.
func Join(sep string, s ...string) string {
	return strings.Join(s, sep)
}
