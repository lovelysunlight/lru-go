package simplelfu

import (
	"errors"

	"github.com/lovelysunlight/lru-go/internal/hashmap"
	"github.com/lovelysunlight/lru-go/internal/list"
)

type Cache[K comparable, V any] struct {
	size      int
	items     *hashmap.Map[K, *list.Entry[K, *LFUValue[V]]]
	evictList *list.DoublyLinkedList[K, *LFUValue[V]]
}

// Cap implements LFUCache.
func (c *Cache[K, V]) Cap() int { return c.size }

// Clear implements LFUCache.
func (c *Cache[K, V]) Clear() {
	c.items.Clear()
	c.evictList.Init()
}

// Contains implements LFUCache.
func (c *Cache[K, V]) Contains(key K) (ok bool) {
	_, ok = c.items.Get(key)
	return ok
}

// Get implements LFUCache.
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	if ent, ok := c.items.Get(key); ok {
		ent.Value.IncrVisit()
		c.moveForward(ent)
		return ent.Value.Value(), true
	}
	return
}

// Keys implements LFUCache.
func (c *Cache[K, V]) Keys() []K {
	keys := make([]K, c.evictList.Len())
	i := 0
	for ent := c.evictList.Back(); ent != nil; ent = ent.PrevEntry() {
		keys[i] = ent.Key
		i++
	}
	return keys
}

// Len implements LFUCache.
func (c *Cache[K, V]) Len() int { return c.evictList.Len() }

// Peek implements LFUCache.
func (c *Cache[K, V]) Peek(key K) (value V, ok bool) {
	if v, ok := c.items.Get(key); ok {
		return v.Value.Value(), true
	}
	return
}

// Returns used of key's value
func (c *Cache[K, V]) PeekUsed(key K) (used uint64, ok bool) {
	if v, ok := c.items.Get(key); ok {
		return v.Value.GetVisit(), true
	}
	return
}

// PeekLeast implements LFUCache.
func (c *Cache[K, V]) PeekLeast() (key K, value V, ok bool) {
	if ent := c.evictList.Back(); ent != nil {
		return ent.Key, ent.Value.Value(), true
	}
	return
}

// Push implements LFUCache.
func (c *Cache[K, V]) Push(key K, value V) (oldKey K, oldValue V, ok bool) {
	return c.capturingPut(key, value, true)
}

// Put implements LFUCache.
func (c *Cache[K, V]) Put(key K, value V) (oldValue V, ok bool) {
	_, v, ok := c.capturingPut(key, value, false)
	return v, ok
}

// Remove implements LFUCache.
func (c *Cache[K, V]) Remove(key K) (value V, ok bool) {
	ent, ok := c.items.Remove(key)
	if ok {
		c.evictList.Remove(ent)
		value = ent.Value.Value()
	}

	return value, ok
}

// RemoveLeast implements LFUCache.
func (c *Cache[K, V]) RemoveLeast() (key K, value V, ok bool) {
	if ent := c.evictList.Back(); ent != nil {
		c.removeElement(ent)
		return ent.Key, ent.Value.Value(), true
	}
	return
}

// Resize implements LFUCache.
func (c *Cache[K, V]) Resize(size int) int {
	diff := c.Len() - size
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < diff; i++ {
		c.RemoveLeast()
	}
	c.size = size
	return diff
}

// Values implements LFUCache.
func (c *Cache[K, V]) Values() []V {
	values := make([]V, c.evictList.Len())
	i := 0
	for ent := c.evictList.Back(); ent != nil; ent = ent.PrevEntry() {
		values[i] = ent.Value.Value()
		i++
	}
	return values
}

// New constructs an LFU of the given size
func New[K comparable, V any](size int) (*Cache[K, V], error) {
	if size <= 0 {
		return nil, errors.New("must provide a positive size")
	}

	c := &Cache[K, V]{
		size:      size,
		evictList: list.NewDoublyLinkedList[K, *LFUValue[V]](),
		items:     hashmap.New[K, *list.Entry[K, *LFUValue[V]]](),
	}
	return c, nil
}

var _ LFUCache[any, any] = (*Cache[any, any])(nil)

func (c *Cache[K, V]) capturingPut(key K, val V, capture bool) (K, V, bool) {
	var (
		oldKey K
		oldVal V
	)

	if ent, ok := c.items.Get(key); ok {
		oldVal := ent.Value.Value()
		ent.Value.SetValue(val)
		ent.Value.IncrVisit()
		c.moveForward(ent)
		return key, oldVal, true
	}

	var ok bool
	if c.Len() == c.Cap() {
		oldKey, oldVal, ok = c.RemoveLeast()
	}
	ent := c.evictList.PushBack(key, &LFUValue[V]{
		value: val, visit: 1,
	})
	c.items.Set(key, ent)
	c.moveForward(ent)

	if !capture {
		var (
			emptyKey K
			emptyVal V
		)
		return emptyKey, emptyVal, false
	}

	return oldKey, oldVal, ok
}

// removeElement is used to remove a given list element from the cache
func (c *Cache[K, V]) removeElement(ent *list.Entry[K, *LFUValue[V]]) {
	c.items.Remove(ent.Key)
	c.evictList.Remove(ent)
}

// moveForward is used to move a given list element to the front of the cache
func (c *Cache[K, V]) moveForward(ent *list.Entry[K, *LFUValue[V]]) {
	// if ent is root entry, ent.PrevEntry() is nil
	for prev := ent.PrevEntry(); prev != nil; prev = ent.PrevEntry() {
		if ent.Value.GetVisit() < prev.Value.GetVisit() {
			break
		}
		// because ent can't get root entry from PrevEntry(),
		// so we don't need to check
		_ = c.evictList.MoveToAt(ent, prev.PrevEntry())
	}
}
