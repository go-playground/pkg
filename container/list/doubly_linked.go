//go:build go1.18
// +build go1.18

package listext

// Node is an element of the doubly linked list.
type Node[V any] struct {
	next, prev *Node[V]
	Value      V
}

// Next returns the nodes next Value or nil if it is at the tail.
func (n *Node[V]) Next() *Node[V] {
	return n.next
}

// Prev returns the nodes previous Value or nil if it is at the head.
func (n *Node[V]) Prev() *Node[V] {
	return n.prev
}

// DoublyLinkedList is a doubly linked list
type DoublyLinkedList[V any] struct {
	head, tail *Node[V]
	len        int
}

// NewDoublyLinked creates a DoublyLinkedList for use.
func NewDoublyLinked[V any]() *DoublyLinkedList[V] {
	return new(DoublyLinkedList[V])
}

// PushFront adds an element first in the list.
func (d *DoublyLinkedList[V]) PushFront(v V) *Node[V] {
	node := &Node[V]{
		Value: v,
	}
	d.pushFront(node)
	return d.head
}

func (d *DoublyLinkedList[V]) pushFront(node *Node[V]) {
	node.next = d.head
	node.prev = nil

	if d.head == nil {
		d.tail = node
	} else {
		d.head.prev = node
	}
	d.head = node
	d.len++
}

// PopFront removes the first element and returns it or nil.
func (d *DoublyLinkedList[V]) PopFront() *Node[V] {
	if d.IsEmpty() {
		return nil
	}

	node := d.head
	d.head = node.next
	if d.head == nil {
		d.tail = nil
	} else {
		d.head.prev = nil
	}
	d.len--
	// ensure no leakage
	node.next, node.prev = nil, nil
	return node
}

// PushBack appends an element to the back of a list.
func (d *DoublyLinkedList[V]) PushBack(v V) *Node[V] {
	node := &Node[V]{
		Value: v,
	}
	d.pushBack(node)
	return d.tail
}

func (d *DoublyLinkedList[V]) pushBack(node *Node[V]) {
	node.prev = d.tail
	node.next = nil

	if d.tail == nil {
		d.head = node
	} else {
		d.tail.next = node
	}
	d.tail = node
	d.len++
}

// PushAfter pushes the supplied Value after the supplied node.
//
// The supplied node must be attached to the current list otherwise undefined behaviour could occur.
func (d *DoublyLinkedList[V]) PushAfter(node *Node[V], v V) *Node[V] {
	newNode := &Node[V]{
		Value: v,
	}
	d.moveAfter(node, newNode)
	return newNode
}

// MoveAfter moves the `moving` node after the supplied `node`.
//
// The supplied `node` and `moving` nodes must be attached to the current list otherwise
// undefined behaviour could occur.
func (d *DoublyLinkedList[V]) MoveAfter(node *Node[V], moving *Node[V]) {
	// first detach node were moving after, in case it was already attached somewhere else in the list.
	d.Remove(moving)
	d.moveAfter(node, moving)
}

func (d *DoublyLinkedList[V]) moveAfter(node *Node[V], moving *Node[V]) {
	next := node.next

	// no next means node == d.tail
	if next == nil {
		d.pushBack(moving)
	} else {
		node.next = moving
		moving.prev = node
		moving.next = next
		next.prev = moving
		d.len++
	}
}

// PushBefore pushes the supplied Value before the supplied node.
//
// The supplied node must be attached to the current list otherwise undefined behaviour could occur.
func (d *DoublyLinkedList[V]) PushBefore(node *Node[V], v V) *Node[V] {
	newNode := &Node[V]{
		Value: v,
	}
	d.moveBefore(node, newNode)
	return newNode
}

// InsertBefore inserts the supplied node before the supplied node.
//
// The supplied node must be attached to the current list otherwise undefined behaviour could occur.
func (d *DoublyLinkedList[V]) InsertBefore(node *Node[V], inserting *Node[V]) {
	d.moveBefore(node, inserting)
}

// InsertAfter inserts the supplied node after the supplied node.
//
// The supplied node must be attached to the current list otherwise undefined behaviour could occur.
func (d *DoublyLinkedList[V]) InsertAfter(node *Node[V], inserting *Node[V]) {
	d.moveAfter(node, inserting)
}

// MoveBefore moves the `moving` node before the supplied `node`.
//
// The supplied `node` and `moving` nodes must be attached to the current list otherwise
// undefined behaviour could occur.
func (d *DoublyLinkedList[V]) MoveBefore(node *Node[V], moving *Node[V]) {
	// first detach node were moving after, in case it was already attached somewhere else in the list.
	d.Remove(moving)
	d.moveBefore(node, moving)
}

func (d *DoublyLinkedList[V]) moveBefore(node *Node[V], moving *Node[V]) {
	prev := node.prev

	// no prev means node == d.head
	if prev == nil {
		d.pushFront(moving)
	} else {
		node.prev = moving
		moving.next = node
		moving.prev = prev
		prev.next = moving
		d.len++
	}
}

// PopBack removes the last element from a list and returns it or nil.
func (d *DoublyLinkedList[V]) PopBack() *Node[V] {
	if d.IsEmpty() {
		return nil
	}

	node := d.tail
	d.tail = node.prev

	if d.tail == nil {
		d.head = nil
	} else {
		d.tail.next = nil
	}
	d.len--
	// ensure no leakage
	node.next, node.prev = nil, nil
	return node
}

// Front returns the front/head element for use without removing it or nil list is empty.
func (d *DoublyLinkedList[V]) Front() *Node[V] {
	return d.head
}

// Back returns the end/tail element for use without removing it or nil list is empty.
func (d *DoublyLinkedList[V]) Back() *Node[V] {
	return d.tail
}

// Remove removes the provided element from the Linked List.
//
// The supplied node must be attached to the current list otherwise undefined behaviour could occur.
func (d *DoublyLinkedList[V]) Remove(node *Node[V]) {
	if node.prev == nil {
		// is head node
		_ = d.PopFront()
	} else if node.next == nil {
		// is tail node
		_ = d.PopBack()
	} else {
		// is both head and tail nodes, must remap
		node.next.prev = node.prev
		node.prev.next = node.next
		// ensure no leakage
		node.next, node.prev = nil, nil
		d.len--
	}
}

// MoveToFront moves the provided node to the front/head.
//
// The supplied node must be attached to the current list otherwise undefined behaviour could occur.
func (d *DoublyLinkedList[V]) MoveToFront(node *Node[V]) {
	d.Remove(node)
	d.pushFront(node)
}

// InsertAtFront pushes the provided node to the front/head.
//
// The supplied node must not be attached to any list otherwise undefined behaviour could occur.
func (d *DoublyLinkedList[V]) InsertAtFront(node *Node[V]) {
	d.pushFront(node)
}

// MoveToBack moves the provided node to the end/tail.
//
// The supplied node must be attached to the current list otherwise undefined behaviour could occur.
func (d *DoublyLinkedList[V]) MoveToBack(node *Node[V]) {
	d.Remove(node)
	d.pushBack(node)
}

// InsertAtBack pushes the provided node to the back/tail.
//
// The supplied node must not be attached to any list otherwise undefined behaviour could occur.
func (d *DoublyLinkedList[V]) InsertAtBack(node *Node[V]) {
	d.pushBack(node)
}

// IsEmpty returns true if the list is empty.
func (d *DoublyLinkedList[V]) IsEmpty() bool {
	return d.len == 0
}

// Len returns length of the Linked List.
func (d *DoublyLinkedList[V]) Len() int {
	return d.len
}

// Clear removes all elements from the Linked List.
func (d *DoublyLinkedList[V]) Clear() {
	// must loop and clean up references to each other.
	for {
		if d.PopBack() == nil {
			break
		}
	}
	d.head, d.tail, d.len = nil, nil, 0
}
