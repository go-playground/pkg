//go:build go1.18
// +build go1.18

package syncext

import (
	resultext "github.com/go-playground/pkg/v5/values/result"
	"testing"

	. "github.com/go-playground/assert/v2"
)

func TestMutex2(t *testing.T) {
	m := NewMutex2(make(map[string]int))
	m.Lock()["foo"] = 1
	m.Unlock()

	m.PerformMut(func(m map[string]int) {
		m["boo"] = 1
	})
	myMap := m.Lock()
	Equal(t, 2, len(myMap))
	Equal(t, myMap["foo"], 1)
	Equal(t, myMap["boo"], 1)
	Equal(t, m.TryLock(), resultext.Err[map[string]int](struct{}{}))
	m.Unlock()

	result := m.TryLock()
	Equal(t, result.IsOk(), true)
	m.Unlock()
}

func TestRWMutex2(t *testing.T) {
	m := NewRWMutex2(make(map[string]int))
	m.Lock()["foo"] = 1
	Equal(t, m.TryLock().IsOk(), false)
	Equal(t, m.TryRLock().IsOk(), false)
	m.Unlock()

	m.PerformMut(func(m map[string]int) {
		m["boo"] = 2
	})
	mp := m.Lock()
	Equal(t, mp["foo"], 1)
	Equal(t, mp["boo"], 2)
	m.Unlock()

	myMap := m.RLock()
	Equal(t, len(myMap), 2)
	Equal(t, m.TryRLock().IsOk(), true)
	m.RUnlock()

	m.Perform(func(m map[string]int) {
		Equal(t, 1, m["foo"])
		Equal(t, 2, m["boo"])
	})
	myMap = m.RLock()
	Equal(t, len(myMap), 2)
	Equal(t, myMap["foo"], 1)
	Equal(t, myMap["boo"], 2)
	m.RUnlock()
}
