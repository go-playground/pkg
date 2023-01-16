//go:build go1.18
// +build go1.18

package syncext

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
	"sync"

	resultext "github.com/go-playground/pkg/v5/values/result"
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

// Lock locks the Mutex and returns value for use, safe for mutation if
//
// If the lock is already in use, the calling goroutine blocks until the mutex is available.
func (m *Mutex[T]) Lock() T {
	m.m.Lock()
	return m.value
}

// Unlock unlocks the Mutex accepting a value to set as the new or mutated value.
// It is optional to pass a new value to be set but NOT required for there reasons:
// 1. If the internal value is already mutable then no need to set as changes apply as they happen.
// 2. If there's a failure working with the locked value you may NOT want to set it, but still unlock.
// 3. Supports locked values that are not mutable.
//
// It is a run-time error if the Mutex is not locked on entry to Unlock.
func (m *Mutex[T]) Unlock(value optionext.Option[T]) {
	if value.IsSome() {
		m.value = value.Unwrap()
	}
	m.m.Unlock()
}

// PerformMut safely locks and unlocks the Mutex values and performs the provided function returning its error if one
// otherwise setting the returned value as the new mutex value.
func (m *Mutex[T]) PerformMut(f func(T) (T, error)) error {
	value := m.Lock()
	result, err := f(value)
	if err != nil {
		m.Unlock(optionext.None[T]())
		return err
	}
	m.Unlock(optionext.Some(result))
	return nil
}

// TryLock tries to lock Mutex and reports whether it succeeded.
// If it does the value is returned for use in the Ok result otherwise Err with empty value.
func (m *Mutex[T]) TryLock() resultext.Result[T, struct{}] {
	if m.m.TryLock() {
		return resultext.Ok[T, struct{}](m.value)
	} else {
		return resultext.Err[T, struct{}](struct{}{})
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

// Lock locks the Mutex and returns value for use, safe for mutation if
//
// If the lock is already in use, the calling goroutine blocks until the mutex is available.
func (m *RWMutex[T]) Lock() T {
	m.rw.Lock()
	return m.value
}

// Unlock unlocks the Mutex accepting a value to set as the new or mutated value.
// It is optional to pass a new value to be set but NOT required for there reasons:
// 1. If the internal value is already mutable then no need to set as changes apply as they happen.
// 2. If there's a failure working with the locked value you may NOT want to set it, but still unlock.
// 3. Supports locked values that are not mutable.
//
// It is a run-time error if the Mutex is not locked on entry to Unlock.
func (m *RWMutex[T]) Unlock(value optionext.Option[T]) {
	if value.IsSome() {
		m.value = value.Unwrap()
	}
	m.rw.Unlock()
}

// PerformMut safely locks and unlocks the RWMutex mutable values and performs the provided function.
func (m *RWMutex[T]) PerformMut(f func(T) (T, error)) error {
	value := m.Lock()
	result, err := f(value)
	if err != nil {
		m.Unlock(optionext.None[T]())
		return err
	}
	m.Unlock(optionext.Some(result))
	return nil
}

// TryLock tries to lock RWMutex and returns the value in the Ok result if successful.
// If it does the value is returned for use in the Ok result otherwise Err with empty value.
func (m *RWMutex[T]) TryLock() resultext.Result[T, struct{}] {
	if m.rw.TryLock() {
		return resultext.Ok[T, struct{}](m.value)
	} else {
		return resultext.Err[T, struct{}](struct{}{})
	}
}

// Perform safely locks and unlocks the RWMutex read-only values and performs the provided function.
func (m *RWMutex[T]) Perform(f func(T) error) error {
	result := m.RLock()
	err := f(result)
	if err != nil {
		m.RUnlock()
		return err
	}
	m.RUnlock()
	return nil
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
func (m *RWMutex[T]) TryRLock() resultext.Result[T, struct{}] {
	if m.rw.TryRLock() {
		return resultext.Ok[T, struct{}](m.value)
	} else {
		return resultext.Err[T, struct{}](struct{}{})
	}
}
