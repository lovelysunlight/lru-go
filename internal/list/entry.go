package list

type Entry[K comparable, V any] struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Entry[K, V]

	// The list to which this element belongs.
	list *DoublyLinkedList[K, V]

	// The LRU Key of this element.
	Key K

	// The Value stored with this element.
	Value V
}

// PrevEntry returns the previous list element or nil.
// note: can't get root entry
func (e *Entry[K, V]) PrevEntry() *Entry[K, V] {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// NextEntry returns the next list element or nil.
// note: can't get root entry
func (e *Entry[K, V]) NextEntry() *Entry[K, V] {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}
