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

func TestLinkedListPushFront(t *testing.T) {

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

func TestLinkedListPushBack(t *testing.T) {

	l := NewDoublyLinked[int]()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	// push some data and then re-check
	zeroNode := l.PushBack(0)
	oneNode := l.PushBack(1)
	twoNode := l.PushBack(2)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 3)

	// test next logic
	Equal(t, zeroNode.Value(), 0)
	Equal(t, zeroNode.Next().Value(), 1)
	Equal(t, zeroNode.Next().Next().Value(), 2)
	Equal(t, zeroNode.Next().Next().Next(), nil)
	Equal(t, zeroNode.Prev(), nil)
	Equal(t, oneNode.Value(), 1)
	Equal(t, oneNode.Next().Value(), 2)
	Equal(t, oneNode.Next().Next(), nil)
	Equal(t, oneNode.Prev().Value(), 0)
	Equal(t, oneNode.Prev().Prev(), nil)
	Equal(t, twoNode.Value(), 2)
	Equal(t, twoNode.Prev().Value(), 1)
	Equal(t, twoNode.Prev().Prev().Value(), 0)
	Equal(t, twoNode.Prev().Prev().Prev(), nil)
	Equal(t, twoNode.Next(), nil)

	// remove middle node and test again
	l.Remove(oneNode)
	Equal(t, oneNode.Value(), 1)
	Equal(t, oneNode.Prev(), nil)
	Equal(t, oneNode.Next(), nil)

	// move to front
	l.MoveToBack(zeroNode)
	Equal(t, l.Front().Value(), 2)
	Equal(t, l.Back().Value(), 0)

	// move to back
	l.MoveToFront(zeroNode)
	Equal(t, l.Front().Value(), 0)
	Equal(t, l.Back().Value(), 2)

	// test clearing
	l.Clear()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)
}

func TestLinkedListMoving(t *testing.T) {

	l := NewDoublyLinked[int]()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	// test pushing after with one node
	node1 := l.PushFront(0)
	node2 := l.PushAfter(node1, 1)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 2)
	Equal(t, l.Front().Value(), node1.Value())
	Equal(t, l.Back().Value(), node2.Value())

	// test moving after with two nodes
	l.MoveAfter(node2, node1)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 2)
	Equal(t, l.Front().Value(), node2.Value())
	Equal(t, l.Back().Value(), node1.Value())

	// test clearing
	l.Clear()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	// test pushing before with one node
	node1 = l.PushFront(0)
	node2 = l.PushBefore(node1, 1)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 2)
	Equal(t, l.Front().Value(), node2.Value())
	Equal(t, l.Back().Value(), node1.Value())

	// test moving before with two nodes
	l.MoveBefore(node2, node1)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 2)
	Equal(t, l.Front().Value(), node1.Value())
	Equal(t, l.Back().Value(), node2.Value())

	// test clearing
	l.Clear()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	// testing the same as above BUT with 3 nodes attached
	node1 = l.PushFront(0)
	node2 = l.PushAfter(node1, 1)
	node3 := l.PushAfter(node2, 2)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 3)
	Equal(t, l.Front().Value(), node1.Value())
	Equal(t, l.Front().Next().Value(), node2.Value())
	Equal(t, l.Back().Value(), node3.Value())
	Equal(t, l.Back().Prev().Value(), node2.Value())

	l.MoveBefore(node2, node3)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 3)
	Equal(t, l.Front().Value(), node1.Value())
	Equal(t, l.Front().Next().Value(), node3.Value())
	Equal(t, l.Back().Value(), node2.Value())
	Equal(t, l.Back().Prev().Value(), node3.Value())

	l.MoveAfter(node3, node1)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 3)
	Equal(t, l.Front().Value(), node3.Value())
	Equal(t, l.Front().Next().Value(), node1.Value())
	Equal(t, l.Back().Value(), node2.Value())
	Equal(t, l.Back().Prev().Value(), node1.Value())

	// test clearing
	l.Clear()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)

	// testing the same as above BUT with 4 nodes attached moving the middle nodes back and forth
	node1 = l.PushFront(0)
	node2 = l.PushAfter(node1, 1)
	node3 = l.PushAfter(node2, 2)
	node4 := l.PushAfter(node3, 3)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 4)
	Equal(t, l.Front().Value(), node1.Value())
	Equal(t, l.Front().Next().Value(), node2.Value())
	Equal(t, l.Front().Next().Next().Value(), node3.Value())
	Equal(t, l.Front().Next().Next().Next().Value(), node4.Value())
	Equal(t, l.Front().Next().Next().Next().Next(), nil)
	Equal(t, l.Back().Value(), node4.Value())
	Equal(t, l.Back().Prev().Value(), node3.Value())
	Equal(t, l.Back().Prev().Prev().Value(), node2.Value())
	Equal(t, l.Back().Prev().Prev().Prev().Value(), node1.Value())
	Equal(t, l.Back().Prev().Prev().Prev().Prev(), nil)

	l.MoveAfter(node3, node2)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 4)
	Equal(t, l.Front().Value(), node1.Value())
	Equal(t, l.Front().Next().Value(), node3.Value())
	Equal(t, l.Front().Next().Next().Value(), node2.Value())
	Equal(t, l.Front().Next().Next().Next().Value(), node4.Value())
	Equal(t, l.Front().Next().Next().Next().Next(), nil)
	Equal(t, l.Back().Value(), node4.Value())
	Equal(t, l.Back().Prev().Value(), node2.Value())
	Equal(t, l.Back().Prev().Prev().Value(), node3.Value())
	Equal(t, l.Back().Prev().Prev().Prev().Value(), node1.Value())
	Equal(t, l.Back().Prev().Prev().Prev().Prev(), nil)

	l.MoveAfter(node2, node3)
	Equal(t, l.IsEmpty(), false)
	Equal(t, l.Len(), 4)
	Equal(t, l.Front().Value(), node1.Value())
	Equal(t, l.Front().Next().Value(), node2.Value())
	Equal(t, l.Front().Next().Next().Value(), node3.Value())
	Equal(t, l.Front().Next().Next().Next().Value(), node4.Value())
	Equal(t, l.Front().Next().Next().Next().Next(), nil)
	Equal(t, l.Back().Value(), node4.Value())
	Equal(t, l.Back().Prev().Value(), node3.Value())
	Equal(t, l.Back().Prev().Prev().Value(), node2.Value())
	Equal(t, l.Back().Prev().Prev().Prev().Value(), node1.Value())
	Equal(t, l.Back().Prev().Prev().Prev().Prev(), nil)

	// test clearing
	l.Clear()
	Equal(t, l.IsEmpty(), true)
	Equal(t, l.Len(), 0)
}
