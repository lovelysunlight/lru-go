package simplelru

// LRUCache is the interface for simple LRU cache.
type LRUCache[K comparable, V any] interface {
	// Adds a value to the cache, returns evicted Value and true if an eviction occurred and
	// updates the "recently used"-ness of the key.
	Put(key K, value V) (oldValue V, ok bool)

	// Adds a value to the cache, returns evicted Key-Value and true if an eviction occurred and
	// updates the "recently used"-ness of the key.
	Push(key K, value V) (oldKey K, oldValue V, ok bool)

	// Returns key's value from the cache and
	// updates the "recently used"-ness of the key. #value, isFound
	Get(key K) (value V, ok bool)

	// Returns key's value without updating the "recently used"-ness of the key.
	Peek(key K) (value V, ok bool)

	// Removes a key from the cache.
	Pop(key K) (value V, ok bool)

	// Removes the oldest entry from cache.
	RemoveOldest() (K, V, bool)

	// Clears all cache entries.
	Clear()

	// Returns the number of items in the cache.
	Len() int

	// Returns the capacity of the cache.
	Cap() int
}
