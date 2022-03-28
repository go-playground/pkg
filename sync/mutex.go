//go:build go1.18

package sync

import (
	"sync"

	"github.com/go-playground/pkg/v5/values/result"
)

// NewMutex creates a new Mutex for use.
func NewMutex[T any](value T) *Mutex[T] {
	return &Mutex[T]{
		value: value,
	}
}

// Mutex creates a type safe mutex wrapper ensuring one cannot access the values of a locked values
// without first gaining a lock.
type Mutex[T any] struct {
	m     sync.Mutex
	value T
}

// PerformMut safely locks and unlocks the Mutex values and performs the provided function.
//
// Too bad Go doesn't support PerformMut[R any](func(T) R) R syntax :(
func (m *Mutex[T]) PerformMut(f func(T)) {
	m.Lock()
	defer m.Unlock()
	f(m.value)
}

// Lock locks the Mutex and returns value for mutable use.
// If the lock is already in use, the calling goroutine blocks until the mutex is available.
func (m *Mutex[T]) Lock() T {
	m.m.Lock()
	return m.value
}

// Unlock unlocks the Mutex. It is a run-time error if the Mutex is not locked on entry to Unlock.
func (m *Mutex[T]) Unlock() {
	m.m.Unlock()
}

// TryLock tries to lock Mutex and reports whether it succeeded.
// If it does the value is returned for use in the Ok result otherwise Err with empty value.
func (m *Mutex[T]) TryLock() result.Result[T, struct{}] {
	if m.m.TryLock() {
		return result.Ok[T, struct{}](m.value)
	} else {
		return result.Err[T, struct{}](struct{}{})
	}
}

// NewRWMutex creates a new RWMutex for use.
func NewRWMutex[T any](value T) *RWMutex[T] {
	return &RWMutex[T]{
		value: value,
	}
}

// RWMutex creates a type safe RWMutex wrapper ensuring one cannot access the values of a locked values
// without first gaining a lock.
type RWMutex[T any] struct {
	rw    sync.RWMutex
	value T
}

// PerformMut safely locks and unlocks the RWMutex mutable values and performs the provided function.
//
// Too bad Go doesn't support PerformMut[R any](func(T) R) R syntax :(
func (m *RWMutex[T]) PerformMut(f func(T)) {
	m.Lock()
	defer m.Unlock()
	f(m.value)
}

// Lock locks mutex and returns values for mutable use.
func (m *RWMutex[T]) Lock() T {
	m.rw.Lock()
	return m.value
}

// Unlock unlocks mutable lock for values.
func (m *RWMutex[T]) Unlock() {
	m.rw.Unlock()
}

// TryLock tries to lock RWMutex and returns the value in the Ok result if successful.
// If it does the value is returned for use in the Ok result otherwise Err with empty value.
func (m *RWMutex[T]) TryLock() result.Result[T, struct{}] {
	if m.rw.TryLock() {
		return result.Ok[T, struct{}](m.value)
	} else {
		return result.Err[T, struct{}](struct{}{})
	}
}

// Perform safely locks and unlocks the RWMutex read-only values and performs the provided function.
//
// Too bad Go doesn't support Perform[R any](func(T) R) R syntax :(
func (m *RWMutex[T]) Perform(f func(T)) {
	m.RLock()
	defer m.RUnlock()
	f(m.value)
}

// RLock locks the RWMutex for reading and returns the value for read-only use.
// It should not be used for recursive read locking; a blocked Lock call excludes new readers from acquiring the lock
func (m *RWMutex[T]) RLock() T {
	m.rw.RLock()
	return m.value
}

// RUnlock undoes a single RLock call; it does not affect other simultaneous readers.
// It is a run-time error if rw is not locked for reading on entry to RUnlock.
func (m *RWMutex[T]) RUnlock() {
	m.rw.RUnlock()
}

// TryRLock tries to lock RWMutex for reading and returns the value in the Ok result if successful.
// If it does the value is returned for use in the Ok result otherwise Err with empty value.
func (m *RWMutex[T]) TryRLock() result.Result[T, struct{}] {
	if m.rw.TryRLock() {
		return result.Ok[T, struct{}](m.value)
	} else {
		return result.Err[T, struct{}](struct{}{})
	}
}
