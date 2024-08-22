package lru

import (
	"sync"

	"github.com/lovelysunlight/lru-go/internal/deepcopy"
	"github.com/lovelysunlight/lru-go/internal/simplelru"
)

type Cache[K comparable, V any] struct {
	mux sync.RWMutex
	lru *simplelru.LRU[K, V]
	// `Get`, `Peek` return value is immutable or not, default true.
	immutable bool
}

// Values returns a slice of the values in the cache, from oldest to newest.
func (c *Cache[K, V]) Cap() int {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.lru.Cap()
}

// Len returns the number of items in the cache.
func (c *Cache[K, V]) Len() int {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.lru.Len()
}

// Get looks up a key's value from the cache.
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	value, ok = c.lru.Get(key)
	if ok && c.immutable {
		return deepcopy.Copy(value), ok
	}
	return value, ok
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *Cache[K, V]) Peek(key K) (value V, ok bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	value, ok = c.lru.Peek(key)
	if ok && c.immutable {
		return deepcopy.Copy(value), ok
	}
	return value, ok
}

// Remove removes the provided key from the cache, returning the value if the
// key was contained.
func (c *Cache[K, V]) Pop(key K) (value V, ok bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.lru.Pop(key)
}

// Pushes a key-value pair into the cache. If an entry with key `key` already exists in
// the cache or another cache entry is removed (due to the lru's capacity),
// then it returns the old entry's key-value pair or not.
func (c *Cache[K, V]) Push(key K, value V) (oldKey K, oldValue V, ok bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.lru.Push(key, value)
}

// Puts a key-value pair into cache. If the key already exists in the cache, then it updates
// the key's value and returns the old value or not.
func (c *Cache[K, V]) Put(key K, value V) (oldValue V, ok bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.lru.Put(key, value)
}

// RemoveOldest removes the oldest item from the cache.
func (c *Cache[K, V]) RemoveOldest() (key K, value V, ok bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.lru.RemoveOldest()
}

// Clears all cache entries.
func (c *Cache[K, V]) Clear() {
	c.mux.Lock()
	c.lru.Clear()
	c.mux.Unlock()
}

func WithMutable[K comparable, V any]() func(*Cache[K, V]) {
	return func(c *Cache[K, V]) {
		c.immutable = false
	}
}

func WithImmutable[K comparable, V any]() func(*Cache[K, V]) {
	return func(c *Cache[K, V]) {
		c.immutable = true
	}
}

func New[K comparable, V any](size int, opts ...func(*Cache[K, V])) (c *Cache[K, V], err error) {
	// create a cache with default settings
	c = &Cache[K, V]{
		immutable: false,
	}
	for _, f := range opts {
		f(c)
	}
	c.lru, err = simplelru.NewLRU[K, V](size)
	return
}

var _ simplelru.LRUCache[any, any] = (*Cache[any, any])(nil)
