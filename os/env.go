//go:build go1.18
// +build go1.18

package osext

import "os"

// Env retrieves the value of the environment variable named by the key.
// If the variable is present in the environment, its value (which may be empty) is returned.
// Otherwise, the provided defaultValue is returned.
// The function is generic and can return any type T, but the caller must ensure that
// the type assertion is valid for the expected type of the environment variable.
func Env[T any](key string, defaultValue T) T {
	if v, ok := os.LookupEnv(key); ok {
		return any(v).(T)
	}
	return defaultValue
}
