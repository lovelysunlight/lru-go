package list

type FIFO[K comparable, V any] struct {
	items map[K]*Entry[K, V]
	root  Entry[K, V] // sentinel list element, only &root, root.prev, and root.next are used
	len   int         // current list Length excluding (this) sentinel element
	size  int
}

// NewFIFOList returns an initialized list.
func NewFIFOList[K comparable, V any](size int) *FIFO[K, V] {
	l := &FIFO[K, V]{size: size}
	return l.Init()
}

// Root returns the root of list l.
func (l *FIFO[K, V]) Root() *Entry[K, V] {
	return &l.root
}

// Back returns the last element of list l or nil if the list is empty.
func (l *FIFO[K, V]) Back() *Entry[K, V] {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *FIFO[K, V]) Len() int { return l.len }

// Size returns the size of list l.
func (l *FIFO[K, V]) Size() int { return l.size }

// lazyInit lazily initializes a zero List Value.
func (l *FIFO[K, V]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// Init initializes or clears list l.
func (l *FIFO[K, V]) Init() *FIFO[K, V] {
	l.items = make(map[K]*Entry[K, V])
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// Get returns the element with the given key from the list.
func (l *FIFO[K, V]) Get(key K) (e *Entry[K, V], ok bool) {
	e, ok = l.items[key]
	return
}

// Push inserts the key-value pair into the list.
func (l *FIFO[K, V]) Push(key K, value V) *Entry[K, V] {
	l.lazyInit()
	if l.size == l.len {
		l.Pop()
	}
	return l.insert(&Entry[K, V]{Key: key, Value: value}, &l.root)
}

// Pop removes the last element of list l and returns it.
func (l *FIFO[K, V]) Pop() (e *Entry[K, V]) {
	if e = l.Back(); e != nil {
		l.delete(e)
	}
	return
}

// Remove removes the element with the given key from the list.
func (l *FIFO[K, V]) Remove(key K) (e *Entry[K, V], ok bool) {
	if e, ok = l.items[key]; ok {
		l.delete(e)
	}
	return
}

// RemoveElement removes the given element from the list.
func (l *FIFO[K, V]) RemoveElement(e *Entry[K, V]) {
	if e != nil {
		l.delete(e)
	}
}

// insert inserts e after at.
func (l *FIFO[K, V]) insert(e, at *Entry[K, V]) *Entry[K, V] {
	l.items[e.Key] = e
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// delete removes e from its list.
func (l *FIFO[K, V]) delete(e *Entry[K, V]) {
	delete(l.items, e.Key)
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	l.len--
}
