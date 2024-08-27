package lru

import (
	"slices"
	"sync"

	"github.com/lovelysunlight/lru-go/internal/list"
	"github.com/lovelysunlight/lru-go/simplelfu"
	"github.com/lovelysunlight/lru-go/simplelru"
)

// Cache is a thread-safe fixed size LRU cache.
type Cache[K comparable, V any] struct {
	deepCopyExt[K, V]

	mux   sync.RWMutex
	lru   simplelru.LRUCache[K, V]
	visit simplelfu.LFUCache[K, V]
	fifo  *list.FIFO[K, V]

	visitThreshold uint64
}

// Values returns the size of LRU cache.
//
//	cache, _ := lru.New[string, string](3)
//	fmt.Println(lru.Cap())
func (c *Cache[K, V]) Cap() int {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.lru.Cap()
}

// Len returns the number of items in the LRU cache.
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

	if value, ok = c.lru.Get(key); ok {
		return c.OptionalCopyValue(value), ok
	}
	if c.IsUpgradeToLRUK() {
		if value, ok = c.visit.Get(key); ok {
			visits, _ := c.visit.PeekVisits(key)
			if c.isExpectedVisits(visits) {
				value, _ = c.visit.Remove(key)
				c.lru.Put(key, value)
				return c.OptionalCopyValue(value), true
			}
		}
	} else if c.IsUpgradeTo2Q() {
		if e, ok := c.fifo.Remove(key); ok {
			c.lru.Put(key, e.Value)
			return c.OptionalCopyValue(e.Value), true
		}
	}
	return
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

	if value, ok = c.lru.Peek(key); ok {
		return c.OptionalCopyValue(value), ok
	}
	if c.IsUpgradeToLRUK() {
		if value, ok = c.visit.Peek(key); ok {
			return c.OptionalCopyValue(value), ok
		}
	}
	return
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

	return c.lru.Contains(key) || c.visit.Contains(key)
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

	if key, value, ok = c.lru.PeekOldest(); ok {
		return c.OptionalCopyKey(key), c.OptionalCopyValue(value), ok
	}
	return
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

	if value, ok = c.lru.Remove(key); ok {
		return value, true
	}
	if c.IsUpgradeToLRUK() {
		return c.visit.Remove(key)
	}
	return
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

	if !c.lru.Contains(key) {
		if c.IsUpgradeToLRUK() {
			oldKey, oldValue, ok = c.visit.Push(key, value)
			visits, _ := c.visit.PeekVisits(key)
			if c.isExpectedVisits(visits) {
				c.moveToLru(key)
			}
			return
		} else if c.IsUpgradeTo2Q() {
			if _, ok := c.fifo.Get(key); !ok {
				c.fifo.Push(key, value)
			}
			return oldKey, oldValue, false
		}
	}

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

	if !c.lru.Contains(key) {
		if c.IsUpgradeToLRUK() {
			oldValue, ok = c.visit.Put(key, value)
			visits, _ := c.visit.PeekVisits(key)
			if c.isExpectedVisits(visits) {
				c.moveToLru(key)
			}
			return
		} else if c.IsUpgradeTo2Q() {
			if _, ok := c.fifo.Get(key); !ok {
				c.fifo.Push(key, value)
			}
			return oldValue, false
		}
	}

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

	return c.OptionalCopyKeyN(
		slices.Concat(c.lru.Keys(), c.visit.Keys()),
	)
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

	return c.OptionalCopyValueN(
		slices.Concat(c.lru.Values(), c.visit.Values()),
	)
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
	if c.IsUpgradeToLRUK() {
		c.visit.Clear()
	}
	c.mux.Unlock()
}

// Resize changes the cache size.
//
//	cache, _ := lru.New[string, string](3)
//	cache.Put("apple", "red")
//	cache.Put("banana", "yellow")
//	cache.Put("orange", "orange")
//	fmt.Println(cache.Resize(2), cache.Cap())
func (c *Cache[K, V]) Resize(size int) (evicted int) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.lru.Resize(size)
}

var _ simplelru.LRUCache[any, any] = (*Cache[any, any])(nil)

// VisitCacheCap returns the size of LFU cache.
func (c *Cache[K, V]) VisitCacheCap() int {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.visit.Cap()
}

// VisitCacheLen returns the number of items in the LFU cache.
func (c *Cache[K, V]) VisitCacheLen() int {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.visit.Len()
}

// Resize changes the LFU cache size.
func (c *Cache[K, V]) VisitCacheResize(size int) int {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.visit.Resize(size)
}

// whether or not enable LRU-K algorithm
func (c *Cache[K, V]) IsUpgradeToLRUK() bool {
	return c.visitThreshold > 1
}

// whether or not enable 2Q algorithm
func (c *Cache[K, V]) IsUpgradeTo2Q() bool {
	return c.fifo != nil && c.fifo.Size() > 0
}

func (c *Cache[K, V]) isExpectedVisits(visits uint64) bool {
	return visits >= c.visitThreshold
}

func (c *Cache[K, V]) moveToLru(key K) {
	if value, ok := c.visit.Remove(key); ok {
		c.lru.Push(key, value)
	}
}

type cacheOptionFunc[K comparable, V any] func(*Cache[K, V])

// Enable to return a deep copy of the value in `Get`, `Peek`,  `PeekOldest`, `Keys` and `Values`.
//
//	cache, _ := lru.New[string, string](3, lru.EnableDeepCopy())
//	cache.Put("foo", []int{1,2})
//	v1, _ := cache.Get("apple")
//	v1[0] = 100
//	v2, _ := cache.Get("apple")
//	fmt.Println(v1, v2)
func EnableDeepCopy[K comparable, V any]() cacheOptionFunc[K, V] {
	return func(c *Cache[K, V]) {
		c.deepCopyExt = deepCopyExt[K, V]{
			copy: true,
		}
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
func DisableDeepCopy[K comparable, V any]() cacheOptionFunc[K, V] {
	return func(c *Cache[K, V]) {
		c.deepCopyExt = deepCopyExt[K, V]{
			copy: false,
		}
	}
}

// Enable LRU-K algorithm
func EnableLRUK[K comparable, V any](threshold uint64) cacheOptionFunc[K, V] {
	return func(c *Cache[K, V]) {
		c.visitThreshold = threshold
	}
}

// Enable 2Q algorithm
func Enable2Q[K comparable, V any](size int) cacheOptionFunc[K, V] {
	return func(c *Cache[K, V]) {
		c.fifo = list.NewFIFOList[K, V](size)
	}
}

// Resize the size of visit cache.
func WithVisitCacheSize[K comparable, V any](size int) cacheOptionFunc[K, V] {
	return func(c *Cache[K, V]) {
		if size > 0 {
			c.visit.Resize(size)
		}
	}
}

// New creates an LRU of the given size.
func New[K comparable, V any](size int, opts ...cacheOptionFunc[K, V]) (c *Cache[K, V], err error) {
	c = &Cache[K, V]{
		deepCopyExt: deepCopyExt[K, V]{
			copy: false,
		},
		visitThreshold: 0,
	}
	c.lru, err = simplelru.New[K, V](size)
	if err != nil {
		return
	}
	c.visit, _ = simplelfu.New[K, V](size)
	for _, f := range opts {
		f(c)
	}
	return
}
