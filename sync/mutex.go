//go:build go1.18

package sync

import (
	"sync"
)

// NewMutex creates a new Mutex for use.
func NewMutex[T any](value T) *Mutex[T] {
	return &Mutex[T]{
		value: value,
	}
}

// Mutex creates a type safe mutex wrapper ensuring one cannot access the value of a locked value
// without first gaining a lock.
type Mutex[T any] struct {
	m     sync.Mutex
	value T
}

// PerformMut safely locks and unlocks the Mutex value and performs the provided function.
//
// Too bad Go doesn't support PerformMut[R any](func(T) R) R syntax :(
func (m *Mutex[T]) PerformMut(f func(T)) {
	m.Lock()
	defer m.Unlock()
	f(m.value)
}

// Lock locks Mutex and returns value for mutable use.
func (m *Mutex[T]) Lock() T {
	m.m.Lock()
	return m.value
}

// Unlock unlocks mutable lock for value.
func (m *Mutex[T]) Unlock() {
	m.m.Unlock()
}

// NewRwMutex creates a new RWMutex for use.
func NewRWMutex[T any](value T) *RWMutex[T] {
	return &RWMutex[T]{
		value: value,
	}
}

// RWMutex creates a type safe RWMutex wrapper ensuring one cannot access the value of a locked value
// without first gaining a lock.
type RWMutex[T any] struct {
	rw    sync.RWMutex
	value T
}

// PerformMut safely locks and unlocks the RWMutex mutable value and performs the provided function.
//
// Too bad Go doesn't support PerformMut[R any](func(T) R) R syntax :(
func (m *RWMutex[T]) PerformMut(f func(T)) {
	m.Lock()
	defer m.Unlock()
	f(m.value)
}

// Lock locks mutex and returns value for mutable use.
func (m *RWMutex[T]) Lock() T {
	m.rw.Lock()
	return m.value
}

// Unlock unlocks mutable lock for value.
func (m *RWMutex[T]) Unlock() {
	m.rw.Unlock()
}

// Perform safely locks and unlocks the RWMutex read-only value and performs the provided function.
//
// Too bad Go doesn't support Perform[R any](func(T) R) R syntax :(
func (m *RWMutex[T]) Perform(f func(T)) {
	m.RLock()
	defer m.RUnlock()
	f(m.value)
}

// Lock locks mutex and returns value for read-only use.
func (m *RWMutex[T]) RLock() T {
	m.rw.RLock()
	return m.value
}

// Unlock unlocks read-only lock for value.
func (m *RWMutex[T]) RUnlock() {
	m.rw.RUnlock()
}
