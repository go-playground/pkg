//go:build go1.18

package syncext

import (
	"testing"

	. "github.com/go-playground/assert/v2"
)

func TestMutex(t *testing.T) {
	m := NewMutex(make(map[string]int))
	m.Lock()["foo"] = 1
	m.Unlock()

	m.PerformMut(func(m map[string]int) {
		m["boo"] = 1
	})

	myMap := m.Lock()
	Equal(t, 2, len(myMap))
	m.Unlock()
}

func TestRWMutex(t *testing.T) {
	m := NewRWMutex(make(map[string]int))
	m.Lock()["foo"] = 1
	Equal(t, m.TryLock().IsOk(), false)
	Equal(t, m.TryRLock().IsOk(), false)
	m.Unlock()

	m.PerformMut(func(m map[string]int) {
		m["boo"] = 2
	})

	myMap := m.RLock()
	Equal(t, len(myMap), 2)
	Equal(t, m.TryRLock().IsOk(), true)
	m.RUnlock()

	m.Perform(func(m map[string]int) {
		Equal(t, 1, m["foo"])
		Equal(t, 2, m["boo"])
	})
}
