package lru

import (
	"sync"

	"github.com/lovelysunlight/lru-go/internal/deepcopy"
	"github.com/lovelysunlight/lru-go/internal/simplelru"
)

type Cache[K comparable, V any] struct {
	mux sync.RWMutex
	lru *simplelru.LRU[K, V]
	// `Get`, `Peek` and so on will acquire a copy value. default false.
	deepCopy bool
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
	if ok && c.deepCopy {
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
	if ok && c.deepCopy {
		return deepcopy.Copy(value), ok
	}
	return value, ok
}

// Checks if a key exists in cache without updating the recent-ness.
func (c *Cache[K, V]) Contains(key K) (ok bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.lru.Contains(key)
}

// Returns the oldest entry without updating the "recently used"-ness of the key.
func (c *Cache[K, V]) PeekOldest() (key K, value V, ok bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	key, value, ok = c.lru.PeekOldest()
	if ok && c.deepCopy {
		return deepcopy.Copy(key), deepcopy.Copy(value), ok
	}
	return key, value, ok
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

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *Cache[K, V]) Keys() []K {
	c.mux.RLock()
	defer c.mux.RUnlock()

	items := c.lru.Keys()
	if c.deepCopy {
		for i, v := range items {
			items[i] = deepcopy.Copy(v)
		}
	}
	return items
}

// Values returns a slice of the values in the cache, from oldest to newest.
func (c *Cache[K, V]) Values() []V {
	c.mux.RLock()
	defer c.mux.RUnlock()

	items := c.lru.Values()
	if c.deepCopy {
		for i, v := range items {
			items[i] = deepcopy.Copy(v)
		}
	}
	return items
}

// Clears all cache entries.
func (c *Cache[K, V]) Clear() {
	c.mux.Lock()
	c.lru.Clear()
	c.mux.Unlock()
}

type cacheOpts struct {
	DeepCopy bool
}

func EnableDeepCopy() func(*cacheOpts) {
	return func(c *cacheOpts) {
		c.DeepCopy = true
	}
}

func DisableDeepCopy() func(*cacheOpts) {
	return func(c *cacheOpts) {
		c.DeepCopy = false
	}
}

func New[K comparable, V any](size int, options ...func(*cacheOpts)) (c *Cache[K, V], err error) {
	// create a cache with default settings
	config := &cacheOpts{DeepCopy: false}
	for _, f := range options {
		f(config)
	}
	c = &Cache[K, V]{
		deepCopy: config.DeepCopy,
	}

	c.lru, err = simplelru.NewLRU[K, V](size)
	return
}

var _ simplelru.LRUCache[any, any] = (*Cache[any, any])(nil)
