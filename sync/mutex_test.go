//go:build go1.18
// +build go1.18

package syncext

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
	"testing"

	. "github.com/go-playground/assert/v2"
)

func TestMutex(t *testing.T) {
	m := NewMutex(make(map[string]int))
	m.Lock()["foo"] = 1
	m.Unlock(optionext.None[map[string]int]())

	err := m.PerformMut(func(m map[string]int) (map[string]int, error) {
		m["boo"] = 1
		return m, nil
	})
	Equal(t, err, nil)

	myMap := m.Lock()
	Equal(t, 2, len(myMap))
	m.Unlock(optionext.None[map[string]int]())
}

func TestRWMutex(t *testing.T) {
	m := NewRWMutex(make(map[string]int))
	m.Lock()["foo"] = 1
	Equal(t, m.TryLock().IsOk(), false)
	Equal(t, m.TryRLock().IsOk(), false)
	m.Unlock(optionext.None[map[string]int]())

	err := m.PerformMut(func(m map[string]int) (map[string]int, error) {
		m["boo"] = 2
		return m, nil
	})
	Equal(t, err, nil)

	myMap := m.RLock()
	Equal(t, len(myMap), 2)
	Equal(t, m.TryRLock().IsOk(), true)
	m.RUnlock()

	err = m.Perform(func(m map[string]int) error {
		Equal(t, 1, m["foo"])
		Equal(t, 2, m["boo"])
		return nil
	})
	Equal(t, err, nil)
}
