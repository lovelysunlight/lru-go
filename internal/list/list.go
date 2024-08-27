package list

import (
	"fmt"
	"slices"
	"strings"
)

// DoublyLinkedList represents a doubly linked list.
// The zero Value for DoublyLinkedList is an empty list ready to use.
type DoublyLinkedList[K comparable, V any] struct {
	root Entry[K, V] // sentinel list element, only &root, root.prev, and root.next are used
	len  int         // current list Length excluding (this) sentinel element
}

// NewDoublyLinkedList returns an initialized list.
func NewDoublyLinkedList[K comparable, V any]() *DoublyLinkedList[K, V] {
	return new(DoublyLinkedList[K, V]).Init()
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *DoublyLinkedList[K, V]) Len() int { return l.len }

// lazyInit lazily initializes a zero List Value.
func (l *DoublyLinkedList[K, V]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// Init initializes or clears list l.
func (l *DoublyLinkedList[K, V]) Init() *DoublyLinkedList[K, V] {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// insert inserts e after at, increments l.len, and returns e.
func (l *DoublyLinkedList[K, V]) insert(e, at *Entry[K, V]) *Entry[K, V] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Entry{Value: v}, at).
func (l *DoublyLinkedList[K, V]) insertValue(k K, v V, at *Entry[K, V]) *Entry[K, V] {
	return l.insert(&Entry[K, V]{Value: v, Key: k}, at)
}

// Remove removes e from its list, decrements l.len
func (l *DoublyLinkedList[K, V]) Remove(e *Entry[K, V]) V {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	l.len--

	return e.Value
}

// Root returns the root of the list.
func (l *DoublyLinkedList[K, V]) Root() *Entry[K, V] {
	return &l.root
}

// move moves e to next to at.
func (l *DoublyLinkedList[K, V]) move(e, at *Entry[K, V]) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *DoublyLinkedList[K, V]) PushFront(k K, v V) *Entry[K, V] {
	l.lazyInit()
	return l.insertValue(k, v, &l.root)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *DoublyLinkedList[K, V]) MoveToFront(e *Entry[K, V]) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in List.Remove about initialization of l
	l.move(e, &l.root)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *DoublyLinkedList[K, V]) PushBack(k K, v V) *Entry[K, V] {
	l.lazyInit()
	return l.insertValue(k, v, l.root.prev)
}

// MoveToAt moves element e to the position after at.
func (l *DoublyLinkedList[K, V]) MoveToAt(e *Entry[K, V], at *Entry[K, V]) (ok bool) {
	if at == nil {
		at = &l.root
	}
	if e.list != l || at.next == e {
		return false
	}
	l.move(e, at)
	return true
}

// Back returns the last element of list l or nil if the list is empty.
func (l *DoublyLinkedList[K, V]) Back() *Entry[K, V] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// Root returns the root of the list.
func (l *DoublyLinkedList[K, V]) Debug() []string {
	asc := []string{"root"}
	for ent := l.root.next; ent.list == l; ent = ent.next {
		asc = append(asc, fmt.Sprintf("%v", ent.Key))
	}
	asc = append(asc, "root")
	desc := []string{"root"}
	for ent := l.root.prev; ent.list == l; ent = ent.prev {
		desc = append(desc, fmt.Sprintf("%v", ent.Key))
	}
	desc = append(desc, "root")
	slices.Reverse(desc)

	return []string{
		strings.Join(asc, " -> "),
		strings.Join(desc, " <- "),
	}
}
