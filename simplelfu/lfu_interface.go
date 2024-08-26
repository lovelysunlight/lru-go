package simplelfu

// LFUCache is the interface for simple LFU cache.
type LFUCache[K comparable, V any] interface {
	// Adds a value to the cache, returns evicted Value and true if an eviction occurred and
	// updates the "recently used"-ness of the key.
	Put(key K, value V) (oldValue V, ok bool)

	// Adds a value to the cache, returns evicted Key-Value and true if an eviction occurred and
	// updates the "recently used"-ness of the key.
	Push(key K, value V) (oldKey K, oldValue V, ok bool)

	// Returns key's value from the cache and
	// updates the "recently used"-ness of the key. #value, isFound
	Get(key K) (value V, ok bool)

	// Checks if a key exists in cache without updating the recent-ness.
	Contains(key K) (ok bool)

	// Returns key's value without updating the "recently used"-ness of the key.
	Peek(key K) (value V, ok bool)

	// Returns used of key's value
	PeekUsed(key K) (used uint64, ok bool)

	// Removes a key from the cache.
	Remove(key K) (value V, ok bool)

	// Removes the least item from cache.
	RemoveLeast() (K, V, bool)

	// Returns the least entry without updating the "recently used"-ness of the key.
	PeekLeast() (K, V, bool)

	// Returns a slice of the keys in the cache, from oldest to newest.
	Keys() []K

	// Values returns a slice of the values in the cache, from oldest to newest.
	Values() []V

	// Clears all cache entries.
	Clear()

	// Returns the number of items in the cache.
	Len() int

	// Returns the capacity of the cache.
	Cap() int

	// Resizes cache, returning number evicted
	Resize(int) int
}
