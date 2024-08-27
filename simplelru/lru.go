package simplelru

import (
	"errors"

	"github.com/lovelysunlight/lru-go/internal/hashmap"
	"github.com/lovelysunlight/lru-go/internal/list"
)

// Cache is a non-thread safe fixed size LRU cache.
type Cache[K comparable, V any] struct {
	size      int
	items     *hashmap.Map[K, *list.Entry[K, V]]
	evictList *list.DoublyLinkedList[K, V]
}

// Returns key's value from the cache and
// updates the "recently used"-ness of the key. #value, isFound
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	node, ok := c.items.Get(key)
	if ok {
		c.evictList.MoveToFront(node)
		value = node.Value
	}

	return value, ok
}

// Returns key's value without updating the "recently used"-ness of the key.
func (c *Cache[K, V]) Peek(key K) (value V, ok bool) {
	node, ok := c.items.Get(key)
	if ok {
		value = node.Value
	}

	return value, ok
}

// Checks if a key exists in cache without updating the recent-ness.
func (c *Cache[K, V]) Contains(key K) (ok bool) {
	_, ok = c.items.Get(key)
	return ok
}

// Returns the oldest entry without updating the "recently used"-ness of the key.
func (c *Cache[K, V]) PeekOldest() (key K, value V, ok bool) {
	if ent := c.evictList.Back(); ent != nil {
		return ent.Key, ent.Value, true
	}
	return
}

// Removes a key from the cache.
func (c *Cache[K, V]) Remove(key K) (value V, ok bool) {
	oldNode, ok := c.items.Remove(key)
	if ok {
		c.evictList.Remove(oldNode)
		value = oldNode.Value
	}

	return value, ok
}

// Adds a value to the cache, returns evicted Key-Value and true if an eviction occurred and
// updates the "recently used"-ness of the key.
func (c *Cache[K, V]) Push(key K, value V) (oldKey K, oldValue V, ok bool) {
	oldKey, oldValue, ok = c.capturingPut(key, value, true)
	return oldKey, oldValue, ok
}

// Adds a value to the cache, returns evicted Value and true if an eviction occurred and
// updates the "recently used"-ness of the key.
func (c *Cache[K, V]) Put(key K, value V) (oldValue V, ok bool) {
	_, oldValue, ok = c.capturingPut(key, value, false)
	return oldValue, ok
}

// New constructs an LRU of the given size
func New[K comparable, V any](size int) (*Cache[K, V], error) {
	if size <= 0 {
		return nil, errors.New("must provide a positive size")
	}

	c := &Cache[K, V]{
		size:      size,
		evictList: list.NewDoublyLinkedList[K, V](),
		items:     hashmap.New[K, *list.Entry[K, V]](),
	}
	return c, nil
}

// Len returns the number of items in the cache.
func (c *Cache[K, V]) Len() int {
	return c.evictList.Len()
}

// Cap returns the capacity of the cache
func (c *Cache[K, V]) Cap() int {
	return c.size
}

// RemoveOldest removes the oldest item from the cache.
func (c *Cache[K, V]) RemoveOldest() (key K, value V, ok bool) {
	if ent := c.evictList.Back(); ent != nil {
		c.removeElement(ent)
		return ent.Key, ent.Value, true
	}
	return
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *Cache[K, V]) Keys() []K {
	keys := make([]K, c.evictList.Len())
	i := 0
	for ent := c.evictList.Back(); ent != nil; ent = ent.PrevEntry() {
		keys[i] = ent.Key
		i++
	}
	return keys
}

// Values returns a slice of the values in the cache, from oldest to newest.
func (c *Cache[K, V]) Values() []V {
	values := make([]V, c.evictList.Len())
	i := 0
	for ent := c.evictList.Back(); ent != nil; ent = ent.PrevEntry() {
		values[i] = ent.Value
		i++
	}
	return values
}

// Clears all cache entries.
func (c *Cache[K, V]) Clear() {
	c.evictList.Init()
	c.items.Clear()
}

// Resize changes the cache size.
func (c *Cache[K, V]) Resize(size int) (evicted int) {
	diff := c.Len() - size
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < diff; i++ {
		c.RemoveOldest()
	}
	c.size = size
	return diff
}

// removeElement is used to remove a given list element from the cache
func (c *Cache[K, V]) removeElement(e *list.Entry[K, V]) {
	c.items.Remove(e.Key)
	c.evictList.Remove(e)
}

func (c *Cache[K, V]) capturingPut(key K, val V, capture bool) (K, V, bool) {
	var (
		oldKey K
		oldVal V
	)

	if node, ok := c.items.Get(key); ok {
		oldVal := node.Value
		node.Value = val
		c.evictList.MoveToFront(node)
		return key, oldVal, true
	}

	var ok bool
	if c.Len() == c.Cap() {
		oldKey, oldVal, ok = c.RemoveOldest()
	}
	c.items.Set(key, c.evictList.PushFront(key, val))

	if !capture {
		var (
			emptyKey K
			emptyVal V
		)
		return emptyKey, emptyVal, false
	}

	return oldKey, oldVal, ok
}

var _ LRUCache[any, any] = (*Cache[any, any])(nil)
