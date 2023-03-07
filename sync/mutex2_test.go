//go:build go1.18
// +build go1.18

package syncext

import (
	resultext "github.com/go-playground/pkg/v5/values/result"
	"sync"
	"testing"

	. "github.com/go-playground/assert/v2"
)

func TestMutex2(t *testing.T) {
	m := NewMutex2(make(map[string]int))
	guard := m.Lock()
	guard.T["foo"] = 1
	guard.Unlock()

	m.PerformMut(func(m map[string]int) {
		m["boo"] = 1
	})
	guard = m.Lock()
	myMap := guard.T
	Equal(t, 2, len(myMap))
	Equal(t, myMap["foo"], 1)
	Equal(t, myMap["boo"], 1)
	Equal(t, m.TryLock(), resultext.Err[MutexGuard[map[string]int, *sync.Mutex]](struct{}{}))
	guard.Unlock()

	result := m.TryLock()
	Equal(t, result.IsOk(), true)
	result.Unwrap().Unlock()
}

func TestRWMutex2(t *testing.T) {
	m := NewRWMutex2(make(map[string]int))
	guard := m.Lock()
	guard.T["foo"] = 1
	Equal(t, m.TryLock().IsOk(), false)
	Equal(t, m.TryRLock().IsOk(), false)
	guard.Unlock()

	m.PerformMut(func(m map[string]int) {
		m["boo"] = 2
	})
	guard = m.Lock()
	mp := guard.T
	Equal(t, mp["foo"], 1)
	Equal(t, mp["boo"], 2)
	guard.Unlock()

	rguard := m.RLock()
	myMap := rguard.T
	Equal(t, len(myMap), 2)
	Equal(t, m.TryRLock().IsOk(), true)
	rguard.RUnlock()

	m.Perform(func(m map[string]int) {
		Equal(t, 1, m["foo"])
		Equal(t, 2, m["boo"])
	})
	rguard = m.RLock()
	myMap = rguard.T
	Equal(t, len(myMap), 2)
	Equal(t, myMap["foo"], 1)
	Equal(t, myMap["boo"], 2)
	rguard.RUnlock()
}
