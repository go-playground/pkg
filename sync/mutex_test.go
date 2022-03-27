//go:build go1.18

package sync

import (
	"testing"
)

func TestMutex(t *testing.T) {
	m := NewMutex(make(map[string]int))
	m.Lock()["foo"] = 1
	m.Unlock()

	myMap := m.Lock()
	if len(myMap) != 1 {
		t.Errorf("Expected map to have 1 element, got %d", len(myMap))
	}
	m.Unlock()
}

func TestRWMutex(t *testing.T) {
	m := NewRWMutex(make(map[string]int))
	m.Lock()["foo"] = 1
	m.Unlock()

	myMap := m.RLock()
	if len(myMap) != 1 {
		t.Errorf("Expected map to have 1 element, got %d", len(myMap))
	}
	m.RUnlock()
}
