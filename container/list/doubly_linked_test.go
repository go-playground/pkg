//go:build go1.18
// +build go1.18

package listext

import (
	. "github.com/go-playground/assert/v2"
	"testing"
)

func TestSingleEntryPopBack(t *testing.T) {

	l := NewDoublyLinked[int]()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	// push some data and then re-check
	zeroNode := l.PushFront(0)
	Equal(t, zeroNode.Value(), 0)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 1)
	Equal(t, zeroNode.Prev(), nil)
	Equal(t, zeroNode.Next(), nil)

	// test popping where one node is both head and tail
	back := l.PopBack()
	Equal(t, back.Value(), 0)
	Equal(t, back.Next(), nil)
	Equal(t, back.Prev(), nil)
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	front := l.PopFront()
	Equal(t, front, nil)
}

func TestSingleEntryPopFront(t *testing.T) {

	l := NewDoublyLinked[int]()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	// push some data and then re-check
	zeroNode := l.PushFront(0)
	Equal(t, zeroNode.Value(), 0)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 1)
	Equal(t, zeroNode.Prev(), nil)
	Equal(t, zeroNode.Next(), nil)

	// test popping where one node is both head and tail
	front := l.PopFront()
	Equal(t, front.Value(), 0)
	Equal(t, front.Next(), nil)
	Equal(t, front.Prev(), nil)
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	back := l.PopBack()
	Equal(t, back, nil)

}

func TestDoubleEntryPopBack(t *testing.T) {

	l := NewDoublyLinked[int]()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	// push some data and then re-check
	zeroNode := l.PushFront(0)
	oneNode := l.PushFront(1)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 2)
	Equal(t, zeroNode.Value(), 0)
	Equal(t, oneNode.Value(), 1)
	Equal(t, zeroNode.Prev().Value(), 1)
	Equal(t, zeroNode.Next(), nil)
	Equal(t, oneNode.Prev(), nil)
	Equal(t, oneNode.Next().Value(), 0)

	back := l.PopBack()
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 1)
	Equal(t, back.Value(), 0)
	Equal(t, back.Next(), nil)
	Equal(t, back.Prev(), nil)
	Equal(t, l.Front().Value(), 1)
	Equal(t, l.Back().Value(), 1)

	back2 := l.PopBack()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)
	Equal(t, back2.Value(), 1)
	Equal(t, back2.Next(), nil)
	Equal(t, back2.Prev(), nil)
	Equal(t, l.Front(), nil)
	Equal(t, l.Back(), nil)
}

func TestTripleEntryPopBack(t *testing.T) {

	l := NewDoublyLinked[int]()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	// push some data and then re-check
	zeroNode := l.PushFront(0)
	oneNode := l.PushFront(1)
	twoNode := l.PushFront(2)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 3)
	Equal(t, zeroNode.Value(), 0)
	Equal(t, oneNode.Value(), 1)
	Equal(t, twoNode.Value(), 2)
	Equal(t, zeroNode.Next(), nil)
	Equal(t, zeroNode.Prev().Value(), 1)
	Equal(t, zeroNode.Prev().Prev().Value(), 2)
	Equal(t, zeroNode.Prev().Prev().Prev(), nil)
	Equal(t, oneNode.Next().Value(), 0)
	Equal(t, oneNode.Next().Next(), nil)
	Equal(t, oneNode.Prev().Value(), 2)
	Equal(t, oneNode.Prev().Prev(), nil)
	Equal(t, twoNode.Prev(), nil)
	Equal(t, twoNode.Next().Value(), 1)
	Equal(t, twoNode.Next().Next().Value(), 0)
	Equal(t, twoNode.Next().Next().Next(), nil)

	// remove front
	l.Remove(twoNode)

	// remove back
	l.Remove(zeroNode)
}

func TestLinkedList(t *testing.T) {

	l := NewDoublyLinked[int]()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	// push some data and then re-check
	zeroNode := l.PushFront(0)
	oneNode := l.PushFront(1)
	twoNode := l.PushFront(2)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 3)

	// test next logic
	Equal(t, zeroNode.Value(), 0)
	Equal(t, zeroNode.Next(), nil)
	Equal(t, zeroNode.Prev().Value(), 1)
	Equal(t, zeroNode.Prev().Prev().Value(), 2)
	Equal(t, zeroNode.Prev().Prev().Prev(), nil)
	Equal(t, oneNode.Value(), 1)
	Equal(t, oneNode.Next().Value(), 0)
	Equal(t, oneNode.Next().Next(), nil)
	Equal(t, oneNode.Prev().Value(), 2)
	Equal(t, oneNode.Prev().Prev(), nil)
	Equal(t, twoNode.Value(), 2)
	Equal(t, twoNode.Prev(), nil)
	Equal(t, twoNode.Next().Value(), 1)
	Equal(t, twoNode.Next().Next().Value(), 0)
	Equal(t, twoNode.Next().Next().Next(), nil)

	// remove middle node and test again
	l.Remove(oneNode)
	Equal(t, oneNode.Value(), 1)
	Equal(t, oneNode.Prev(), nil)
	Equal(t, oneNode.Next(), nil)

	// move to front
	l.MoveToFront(zeroNode)
	Equal(t, l.Front().Value(), 0)
	Equal(t, l.Back().Value(), 2)

	// move to back
	l.MoveToBack(zeroNode)
	Equal(t, l.Front().Value(), 2)
	Equal(t, l.Back().Value(), 0)

	// test clearing
	l.Clear()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)
}
