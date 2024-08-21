package lru

import (
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

type lruEntry[K comparable, V any] struct {
	Key  Key[K]
	Val  Value[V]
	Next *lruEntry[K, V]
	Prev *lruEntry[K, V]
}

func newLRUEntry[K comparable, V any](key K, val V) *lruEntry[K, V] {
	return &lruEntry[K, V]{
		Key:  Key[K]{K: key},
		Val:  Value[V]{V: val},
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
	index hashmap.Map[Key[K], *lruEntry[K, V]]
	cap   int

	head *lruEntry[K, V]
	tail *lruEntry[K, V]
}

func New[K comparable, V any](cap int) *lruCache[K, V] {

	cache := &lruCache[K, V]{
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
	return c.index.Len()
}

func (c *lruCache[K, V]) Cap() int {
	return c.cap
}

func (c *lruCache[K, V]) Put(key K, val V) option.Option[V] {
	data := c.capturingPut(key, val, false)
	if data.IsSome() {
		return option.Some(data.Unwrap().Val.GetDeepCopyV())
	}
	return option.None[V]()
}

func (c *lruCache[K, V]) Push(key K, val V) option.Option[tupleKV[K, V]] {
	return c.capturingPut(key, val, true)
}

func (c *lruCache[K, V]) Peek(key K) option.Option[V] {
	node, ok := c.index.Get(Key[K]{K: key})
	if !ok {
		return option.None[V]()
	}

	return option.Some(node.Val.GetDeepCopyV())
}

func (c *lruCache[K, V]) Get(key K) option.Option[V] {
	if node, ok := c.index.Get(Key[K]{K: key}); ok {
		c.detach(node)
		c.attach(node)

		return option.Some(node.Val.GetDeepCopyV())
	}

	return option.None[V]()
}

func (c *lruCache[K, V]) Pop(key K) option.Option[V] {
	oldNode, ok := c.index.Remove(Key[K]{K: key})
	if !ok {
		return option.None[V]()
	}

	c.detach(oldNode)

	return option.Some(oldNode.Val.V)
}

func (c *lruCache[K, V]) capturingPut(key K, val V, capture bool) option.Option[tupleKV[K, V]] {
	if node, ok := c.index.Get(Key[K]{K: key}); ok {
		oldVal := node.Val.DeepCopy()
		node.Val = Value[V]{V: val}

		c.detach(node)
		c.attach(node)

		return option.Some(tupleKV[K, V]{
			Key: node.Key,
			Val: oldVal,
		})
	}

	replaced, node := c.replaceOrCreateNode(key, val)
	c.attach(node)
	c.index.Set(node.Key, node)

	return replaced.Filter(func(_ tupleKV[K, V]) bool {
		return capture
	})
}

func (c *lruCache[K, V]) replaceOrCreateNode(key K, val V) (option.Option[tupleKV[K, V]], *lruEntry[K, V]) {
	if c.Len() == c.Cap() {
		oldKey := c.tail.Prev.Key
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
	oldKey, oldVal := node.Key.DeepCopy(), node.Val.DeepCopy()
	node.Key = Key[K]{K: newKey}
	node.Val = Value[V]{V: newVal}

	return oldKey, oldVal
}
