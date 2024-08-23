package lru

import (
	"sync"

	"github.com/lovelysunlight/lru-go/internal/deepcopy"
	"github.com/lovelysunlight/lru-go/internal/simplelru"
)

// Cache is a thread-safe fixed size LRU cache.
type Cache[K comparable, V any] struct {
	mux      sync.RWMutex
	lru      simplelru.LRUCache[K, V]
	deepCopy bool
}

// Values returns a slice of the values in the cache, from oldest to newest.
//
//	cache, _ := lru.New[string, string](3)
//	fmt.Println(lru.Cap())
func (c *Cache[K, V]) Cap() int {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.lru.Cap()
}

// Len returns the number of items in the cache.
//
//	cache, _ := lru.New[string, string](3)
//	cache.Put("apple", "red")
//	cache.Put("banana", "yellow")
//	fmt.Println(lru.Len())
func (c *Cache[K, V]) Len() int {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.lru.Len()
}

// Get looks up a key's value from the cache with updating
// the "recently used"-ness of the key.
//
//	cache, _ := lru.New[string, string](3)
//	cache.Put("apple", "red")
//	cache.Put("banana", "yellow")
//	v, ok := cache.Get("banana")
//	fmt.Println(ok, v)
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	value, ok = c.lru.Get(key)
	if ok && c.deepCopy {
		return deepcopy.Copy(value), ok
	}
	return value, ok
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
//
//	cache, _ := lru.New[string, string](3)
//	cache.Put("apple", "red")
//	cache.Put("banana", "yellow")
//	v, ok := cache.Peek("banana")
//	fmt.Println(ok, v)
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
//
//	cache, _ := lru.New[string, string](3)
//	cache.Put("apple", "red")
//	cache.Put("banana", "yellow")
//	fmt.Println(cache.Contains("banana"))
func (c *Cache[K, V]) Contains(key K) (ok bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.lru.Contains(key)
}

// Returns the oldest entry without updating the "recently used"-ness of the key.
//
//	cache, _ := lru.New[string, string](3)
//	cache.Put("apple", "red")
//	cache.Put("banana", "yellow")
//	k, v, ok := cache.PeekOldest()
//	fmt.Println(ok, k, v)
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
//
//	cache, _ := lru.New[string, string](3)
//	cache.Put("apple", "red")
//	cache.Put("banana", "yellow")
//	v, ok := cache.Remove("apple")
//	fmt.Println(ok, v)
func (c *Cache[K, V]) Remove(key K) (value V, ok bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.lru.Remove(key)
}

// Pushes a key-value pair into the cache. If an entry with key `key` already exists in
// the cache or another cache entry is removed (due to the lru's capacity),
// then it returns the old entry's key-value pair or not.
//
//	cache, _ := lru.New[string, string](3)
//	cache.Push("apple", "red")
func (c *Cache[K, V]) Push(key K, value V) (oldKey K, oldValue V, ok bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.lru.Push(key, value)
}

// Puts a key-value pair into cache. If the key already exists in the cache, then it updates
// the key's value and returns the old value or not.
//
//	cache, _ := lru.New[string, string](3)
//	cache.Put("apple", "red")
//	fmt.Println(ok, k, v)
func (c *Cache[K, V]) Put(key K, value V) (oldValue V, ok bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.lru.Put(key, value)
}

// RemoveOldest removes the oldest item from the cache.
//
//	cache, _ := lru.New[string, string](3)
//	k, v, ok := cache.RemoveOldest()
func (c *Cache[K, V]) RemoveOldest() (key K, value V, ok bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.lru.RemoveOldest()
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
//
//	cache, _ := lru.New[string, string](3)
//	cache.Put("apple", "red")
//	cache.Put("banana", "yellow")
//	fmt.Printf("%+v", cache.Keys())
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
//
//	cache, _ := lru.New[string, string](3)
//	cache.Put("apple", "red")
//	cache.Put("banana", "yellow")
//	fmt.Printf("%+v", cache.Values())
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
//
//	cache, _ := lru.New[string, string](3)
//	cache.Put("apple", "red")
//	cache.Put("banana", "yellow")
//	cache.Clear()
//	fmt.Println(cache.Len())
func (c *Cache[K, V]) Clear() {
	c.mux.Lock()
	c.lru.Clear()
	c.mux.Unlock()
}

type cacheOpts struct {
	DeepCopy bool
}

// Enable to return a deep copy of the value in `Get`, `Peek`,  `PeekOldest`, `Keys` and `Values`.
//
//	cache, _ := lru.New[string, string](3, lru.EnableDeepCopy())
//	cache.Put("foo", []int{1,2})
//	v1, _ := cache.Get("apple")
//	v1[0] = 100
//	v2, _ := cache.Get("apple")
//	fmt.Println(v1, v2)
func EnableDeepCopy() func(*cacheOpts) {
	return func(c *cacheOpts) {
		c.DeepCopy = true
	}
}

// Disable to return a deep copy of the value in `Get`, `Peek`,  `PeekOldest`, `Keys` and `Values`.
//
//	cache, _ := lru.New[string, string](3, lru.DisableDeepCopy())
//	cache.Put("foo", []int{1,2})
//	v1, _ := cache.Get("apple")
//	v1[0] = 100
//	v2, _ := cache.Get("apple")
//	fmt.Println(v1, v2)
func DisableDeepCopy() func(*cacheOpts) {
	return func(c *cacheOpts) {
		c.DeepCopy = false
	}
}

// New creates an LRU of the given size.
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
