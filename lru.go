package lru

import (
	"sync"

	"github.com/lovelysunlight/lru-go/internal/hashmap"
	"github.com/lovelysunlight/lru-go/internal/option"
	"github.com/lovelysunlight/lru-go/internal/pack"
)

type lruCache[K comparable, V any] struct {
	mux   sync.RWMutex
	index hashmap.Map[pack.Key[K], *lruEntry[K, V]]
	cap   int

	head *lruEntry[K, V]
	tail *lruEntry[K, V]
}

func New[K comparable, V any](cap int) *lruCache[K, V] {

	cache := &lruCache[K, V]{
		mux:   sync.RWMutex{},
		index: hashmap.New[pack.Key[K], *lruEntry[K, V]](),
		cap:   cap,
		head:  newLRUEntrySigil[K, V](),
		tail:  newLRUEntrySigil[K, V](),
	}

	cache.head.next = cache.tail
	cache.tail.prev = cache.head

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
		return option.Some(data.Unwrap().GetVal().Unpack())
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

	node, ok := c.index.Get(*newKey(key))
	if !ok {
		return option.None[V]()
	}

	return option.Some(node.GetData().GetVal().DeepCopy().Unpack())
}

func (c *lruCache[K, V]) Get(key K) option.Option[V] {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if node, ok := c.index.Get(*newKey(key)); ok {
		c.detach(node)
		c.attach(node)

		return option.Some(node.GetData().GetVal().DeepCopy().Unpack())
	}

	return option.None[V]()
}

func (c *lruCache[K, V]) Pop(key K) option.Option[V] {
	c.mux.Lock()
	defer c.mux.Unlock()

	oldNode, ok := c.index.Remove(*newKey(key))
	if !ok {
		return option.None[V]()
	}

	c.detach(oldNode)

	return option.Some(oldNode.GetData().GetVal().Unpack())
}

func (c *lruCache[K, V]) capturingPut(key K, val V, capture bool) option.Option[tupleKV[K, V]] {
	if node, ok := c.index.Get(*newKey(key)); ok {
		oldNodeData := node.data.DeepCopy()
		node.data.val = newVal(val)

		c.detach(node)
		c.attach(node)

		return option.Some(oldNodeData)
	}

	replaced, node := c.replaceOrCreateNode(key, val)
	c.attach(node)
	c.index.Set(*node.data.key, node)

	return replaced.Filter(func(_ tupleKV[K, V]) bool {
		return capture
	})
}

func (c *lruCache[K, V]) replaceOrCreateNode(key K, val V) (option.Option[tupleKV[K, V]], *lruEntry[K, V]) {
	if c.len() == c.Cap() {
		oldKey := c.tail.Prev().GetData().GetKey()
		oldNode, _ := c.index.Remove(*oldKey)
		oldTupleKV := oldNode.Replace(key, val)
		c.detach(oldNode)

		return option.Some(oldTupleKV), oldNode
	}

	return option.None[tupleKV[K, V]](), newLRUEntry(key, val)
}

func (c *lruCache[K, V]) detach(node *lruEntry[K, V]) {
	node.Prev().PushBack(node.Next())
	node.Next().PushFront(node.Prev())
}

func (c *lruCache[K, V]) attach(node *lruEntry[K, V]) {
	node.PushBack(c.head.Next())
	node.PushFront(c.head)
	c.head.PushBack(node)
	node.Next().PushFront(node)
}

func newKey[T comparable](v T) *pack.Key[T] {
	return pack.Pack[T, *pack.Key[T]](v)
}

func newVal[T any](v T) *pack.Value[T] {
	return pack.Pack[T, *pack.Value[T]](v)
}
