package simplelru

import (
	"errors"

	"github.com/lovelysunlight/lru-go/internal/hashmap"
)

type LRU[K comparable, V any] struct {
	size      int
	items     *hashmap.Map[K, *Entry[K, V]]
	evictList *LruList[K, V]
}

// Get implements LRUCache.
func (c *LRU[K, V]) Get(key K) (value V, ok bool) {
	node, ok := c.items.Get(key)
	if ok {
		c.evictList.MoveToFront(node)
		value = node.Value
	}

	return value, ok
}

// Peek implements LRUCache.
func (c *LRU[K, V]) Peek(key K) (value V, ok bool) {
	node, ok := c.items.Get(key)
	if ok {
		value = node.Value
	}

	return value, ok
}

// Pop implements LRUCache.
func (c *LRU[K, V]) Pop(key K) (value V, ok bool) {
	oldNode, ok := c.items.Remove(key)
	if ok {
		c.evictList.Remove(oldNode)
		value = oldNode.Value
	}

	return value, ok
}

// Push implements LRUCache.
func (c *LRU[K, V]) Push(key K, value V) (oldKey K, oldValue V, ok bool) {
	oldKey, oldValue, ok = c.capturingPut(key, value, true)
	return oldKey, oldValue, ok
}

// Put implements LRUCache.
func (c *LRU[K, V]) Put(key K, value V) (oldValue V, ok bool) {
	_, oldValue, ok = c.capturingPut(key, value, false)
	return oldValue, ok
}

// NewLRU constructs an LRU of the given size
func NewLRU[K comparable, V any](size int) (*LRU[K, V], error) {
	if size <= 0 {
		return nil, errors.New("must provide a positive size")
	}

	c := &LRU[K, V]{
		size:      size,
		evictList: NewList[K, V](),
		items:     hashmap.New[K, *Entry[K, V]](),
	}
	return c, nil
}

// Len returns the number of items in the cache.
func (c *LRU[K, V]) Len() int {
	return c.evictList.Length()
}

// Cap returns the capacity of the cache
func (c *LRU[K, V]) Cap() int {
	return c.size
}

// RemoveOldest removes the oldest item from the cache.
func (c *LRU[K, V]) RemoveOldest() (key K, value V, ok bool) {
	if ent := c.evictList.Back(); ent != nil {
		c.removeElement(ent)
		return ent.Key, ent.Value, true
	}
	return
}

// Clears all cache entries.
func (c *LRU[K, V]) Clear() {
	c.evictList.Init()
	c.items.Clear()
}

// removeElement is used to remove a given list element from the cache
func (c *LRU[K, V]) removeElement(e *Entry[K, V]) {
	c.items.Remove(e.Key)
	c.evictList.Remove(e)
}

func (c *LRU[K, V]) capturingPut(key K, val V, capture bool) (K, V, bool) {
	var (
		oldKey K
		oldVal V
	)

	if node, ok := c.items.Get(key); ok {
		c.items.Set(key, c.evictList.PushFront(key, val))
		return node.Key, node.Value, true
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

var _ LRUCache[any, any] = (*LRU[any, any])(nil)
