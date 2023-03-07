//go:build go1.18
// +build go1.18

package syncext

import (
	"sync"

	resultext "github.com/go-playground/pkg/v5/values/result"
)

type MutexGuard[T any, M interface{ Unlock() }] struct {
	m M
	T T
}

func (g MutexGuard[T, M]) Unlock() {
	g.m.Unlock()
}

// NewMutex2 creates a new Mutex for use.
func NewMutex2[T any, M *sync.Mutex](value T) Mutex2[T, M] {
	return Mutex2[T, M]{
		m:     new(sync.Mutex),
		value: value,
	}
}

// Mutex2 creates a type safe mutex wrapper ensuring one cannot access the values of a locked values
// without first gaining a lock.
type Mutex2[T any, M *sync.Mutex] struct {
	m     *sync.Mutex
	value T
}

// Lock locks the Mutex and returns value for use, safe for mutation if
//
// If the lock is already in use, the calling goroutine blocks until the mutex is available.
func (m Mutex2[T, M]) Lock() MutexGuard[T, *sync.Mutex] {
	m.m.Lock()
	return MutexGuard[T, *sync.Mutex]{
		m: m.m,
		T: m.value,
	}
}

//// Unlock unlocks the Mutex accepting a value to set as the new or mutated value.
//// It is optional to pass a new value to be set but NOT required for there reasons:
//// 1. If the internal value is already mutable then no need to set as changes apply as they happen.
//// 2. If there's a failure working with the locked value you may NOT want to set it, but still unlock.
//// 3. Supports locked values that are not mutable.
////
//// It is a run-time error if the Mutex is not locked on entry to Unlock.
//func (m Mutex2[T]) Unlock() {
//	m.m.Unlock()
//}

// PerformMut safely locks and unlocks the Mutex values and performs the provided function returning its error if one
// otherwise setting the returned value as the new mutex value.
func (m Mutex2[T, M]) PerformMut(f func(T)) {
	guard := m.Lock()
	f(guard.T)
	guard.Unlock()
}

// TryLock tries to lock Mutex and reports whether it succeeded.
// If it does the value is returned for use in the Ok result otherwise Err with empty value.
func (m Mutex2[T, M]) TryLock() resultext.Result[MutexGuard[T, *sync.Mutex], struct{}] {
	if m.m.TryLock() {
		return resultext.Ok[MutexGuard[T, *sync.Mutex], struct{}](MutexGuard[T, *sync.Mutex]{
			m: m.m,
			T: m.value,
		})
	} else {
		return resultext.Err[MutexGuard[T, *sync.Mutex], struct{}](struct{}{})
	}
}

type RMutexGuard[T any] struct {
	rw *sync.RWMutex
	T  T
}

func (g RMutexGuard[T]) RUnlock() {
	g.rw.RUnlock()
}

// NewRWMutex2 creates a new RWMutex for use.
func NewRWMutex2[T any](value T) RWMutex2[T] {
	return RWMutex2[T]{
		rw:    new(sync.RWMutex),
		value: value,
	}
}

// RWMutex2 creates a type safe RWMutex wrapper ensuring one cannot access the values of a locked values
// without first gaining a lock.
type RWMutex2[T any] struct {
	rw    *sync.RWMutex
	value T
}

// Lock locks the Mutex and returns value for use, safe for mutation if
//
// If the lock is already in use, the calling goroutine blocks until the mutex is available.
func (m RWMutex2[T]) Lock() MutexGuard[T, *sync.RWMutex] {
	m.rw.Lock()
	return MutexGuard[T, *sync.RWMutex]{
		m: m.rw,
		T: m.value,
	}
}

// PerformMut safely locks and unlocks the RWMutex mutable values and performs the provided function.
func (m RWMutex2[T]) PerformMut(f func(T)) {
	guard := m.Lock()
	f(guard.T)
	guard.Unlock()
}

// TryLock tries to lock RWMutex and returns the value in the Ok result if successful.
// If it does the value is returned for use in the Ok result otherwise Err with empty value.
func (m RWMutex2[T]) TryLock() resultext.Result[MutexGuard[T, *sync.RWMutex], struct{}] {
	if m.rw.TryLock() {
		return resultext.Ok[MutexGuard[T, *sync.RWMutex], struct{}](
			MutexGuard[T, *sync.RWMutex]{
				m: m.rw,
				T: m.value,
			})
	} else {
		return resultext.Err[MutexGuard[T, *sync.RWMutex], struct{}](struct{}{})
	}
}

// Perform safely locks and unlocks the RWMutex read-only values and performs the provided function.
func (m RWMutex2[T]) Perform(f func(T)) {
	guard := m.RLock()
	f(guard.T)
	guard.RUnlock()
}

// RLock locks the RWMutex for reading and returns the value for read-only use.
// It should not be used for recursive read locking; a blocked Lock call excludes new readers from acquiring the lock
func (m RWMutex2[T]) RLock() RMutexGuard[T] {
	m.rw.RLock()
	return RMutexGuard[T]{
		rw: m.rw,
		T:  m.value,
	}
}

// TryRLock tries to lock RWMutex for reading and returns the value in the Ok result if successful.
// If it does the value is returned for use in the Ok result otherwise Err with empty value.
func (m RWMutex2[T]) TryRLock() resultext.Result[RMutexGuard[T], struct{}] {
	if m.rw.TryRLock() {
		return resultext.Ok[RMutexGuard[T], struct{}](
			RMutexGuard[T]{
				rw: m.rw,
				T:  m.value,
			},
		)
	} else {
		return resultext.Err[RMutexGuard[T], struct{}](struct{}{})
	}
}
