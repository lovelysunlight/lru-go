package lru

import (
	"sync"

	"github.com/JimChenWYU/lru-go/internal/deepcopy"
	"github.com/JimChenWYU/lru-go/internal/hashmap"
	"github.com/JimChenWYU/lru-go/internal/option"
)

type Key[T comparable] struct {
	K T
}

type Value[T any] struct {
	V T
}

func (v Key[T]) DeepCopy() Key[T] {
	return Key[T]{
		K: deepcopy.Copy(v.K),
	}
}

func (v Value[T]) DeepCopy() Value[T] {
	return Value[T]{
		V: deepcopy.Copy(v.V),
	}
}

func (v Value[T]) GetDeepCopyV() T {
	return deepcopy.Copy(v.V)
}

type tupleKV[K comparable, V any] struct {
	Key Key[K]
	Val Value[V]
}

func (v *tupleKV[K, V]) DeepCopy() *tupleKV[K, V] {
	return &tupleKV[K, V]{
		Key: v.Key.DeepCopy(),
		Val: v.Val.DeepCopy(),
	}
}

type lruEntry[K comparable, V any] struct {
	Data *tupleKV[K, V]
	Next *lruEntry[K, V]
	Prev *lruEntry[K, V]
}

func newLRUEntry[K comparable, V any](key K, val V) *lruEntry[K, V] {
	return &lruEntry[K, V]{
		Data: &tupleKV[K, V]{
			Key: Key[K]{K: key},
			Val: Value[V]{V: val},
		},
		Next: nil,
		Prev: nil,
	}
}

func newLRUEntrySigil[K comparable, V any]() *lruEntry[K, V] {
	return &lruEntry[K, V]{
		Next: nil,
		Prev: nil,
	}
}

type lruCache[K comparable, V any] struct {
	mux   sync.RWMutex
	index hashmap.Map[Key[K], *lruEntry[K, V]]
	cap   int

	head *lruEntry[K, V]
	tail *lruEntry[K, V]
}

func New[K comparable, V any](cap int) *lruCache[K, V] {

	cache := &lruCache[K, V]{
		mux:   sync.RWMutex{},
		index: hashmap.New[Key[K], *lruEntry[K, V]](),
		cap:   cap,
		head:  newLRUEntrySigil[K, V](),
		tail:  newLRUEntrySigil[K, V](),
	}

	cache.head.Next = cache.tail
	cache.tail.Prev = cache.head

	return cache
}

func (c *lruCache[K, V]) Len() int {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.len()
}

// internal use only
func (c *lruCache[K, V]) len() int {
	return c.index.Len()
}

func (c *lruCache[K, V]) Cap() int {
	return c.cap
}

func (c *lruCache[K, V]) Put(key K, val V) option.Option[V] {
	c.mux.Lock()
	defer c.mux.Unlock()

	data := c.capturingPut(key, val, false)
	if data.IsSome() {
		return option.Some(data.Unwrap().Val.GetDeepCopyV())
	}
	return option.None[V]()
}

func (c *lruCache[K, V]) Push(key K, val V) option.Option[tupleKV[K, V]] {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.capturingPut(key, val, true)
}

func (c *lruCache[K, V]) Peek(key K) option.Option[V] {
	c.mux.RLock()
	defer c.mux.RUnlock()

	node, ok := c.index.Get(Key[K]{K: key})
	if !ok {
		return option.None[V]()
	}

	return option.Some(node.Data.Val.GetDeepCopyV())
}

func (c *lruCache[K, V]) Get(key K) option.Option[V] {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if node, ok := c.index.Get(Key[K]{K: key}); ok {
		c.detach(node)
		c.attach(node)

		return option.Some(node.Data.Val.GetDeepCopyV())
	}

	return option.None[V]()
}

func (c *lruCache[K, V]) Pop(key K) option.Option[V] {
	c.mux.Lock()
	defer c.mux.Unlock()

	oldNode, ok := c.index.Remove(Key[K]{K: key})
	if !ok {
		return option.None[V]()
	}

	c.detach(oldNode)

	return option.Some(oldNode.Data.Val.V)
}

func (c *lruCache[K, V]) capturingPut(key K, val V, capture bool) option.Option[tupleKV[K, V]] {
	if node, ok := c.index.Get(Key[K]{K: key}); ok {
		oldNodeData := node.Data.DeepCopy()
		node.Data.Val = Value[V]{V: val}

		c.detach(node)
		c.attach(node)

		return option.Some(*oldNodeData)
	}

	replaced, node := c.replaceOrCreateNode(key, val)
	c.attach(node)
	c.index.Set(node.Data.Key, node)

	return replaced.Filter(func(_ tupleKV[K, V]) bool {
		return capture
	})
}

func (c *lruCache[K, V]) replaceOrCreateNode(key K, val V) (option.Option[tupleKV[K, V]], *lruEntry[K, V]) {
	if c.len() == c.Cap() {
		oldKey := c.tail.Prev.Data.Key
		oldNode, _ := c.index.Remove(oldKey)
		oldKey, oldVal := replace(oldNode, key, val)
		c.detach(oldNode)

		return option.Some(tupleKV[K, V]{
			Key: oldKey,
			Val: oldVal,
		}), oldNode
	}

	return option.None[tupleKV[K, V]](), newLRUEntry(key, val)
}

func (c *lruCache[K, V]) detach(node *lruEntry[K, V]) {
	if node != nil {
		node.Prev.Next = node.Next
		node.Next.Prev = node.Prev
	}
}

func (c *lruCache[K, V]) attach(node *lruEntry[K, V]) {
	if node != nil {
		node.Next = c.head.Next
		node.Prev = c.head
		c.head.Next = node
		node.Next.Prev = node
	}
}

func replace[K comparable, V any](node *lruEntry[K, V], newKey K, newVal V) (Key[K], Value[V]) {
	oldKey, oldVal := node.Data.Key.DeepCopy(), node.Data.Val.DeepCopy()
	node.Data.Key = Key[K]{K: newKey}
	node.Data.Val = Value[V]{V: newVal}

	return oldKey, oldVal
}
