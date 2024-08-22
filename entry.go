package lru

import "github.com/lovelysunlight/lru-go/internal/pack"

type tupleKV[K comparable, V any] struct {
	key *pack.Key[K]
	val *pack.Value[V]
}

func (v tupleKV[K, V]) GetVal() *pack.Value[V] {
	return v.val
}

func (v tupleKV[K, V]) GetKey() *pack.Key[K] {
	return v.key
}

func (v tupleKV[K, V]) DeepCopy() tupleKV[K, V] {
	return tupleKV[K, V]{
		key: pack.DeepCopy(v.key),
		val: pack.DeepCopy(v.val),
	}
}

type lruEntry[K comparable, V any] struct {
	data *tupleKV[K, V]
	next *lruEntry[K, V]
	prev *lruEntry[K, V]
}

func newLRUEntry[K comparable, V any](key K, val V) *lruEntry[K, V] {
	return &lruEntry[K, V]{
		data: &tupleKV[K, V]{
			key: pack.Pack[K, *pack.Key[K]](key),
			val: pack.Pack[V, *pack.Value[V]](val),
		},
		next: nil,
		prev: nil,
	}
}

func newLRUEntrySigil[K comparable, V any]() *lruEntry[K, V] {
	return &lruEntry[K, V]{
		next: nil,
		prev: nil,
	}
}

func (e *lruEntry[K, V]) Prev() *lruEntry[K, V] {
	if e == nil {
		return nil
	}
	return e.prev
}

func (e *lruEntry[K, V]) Next() *lruEntry[K, V] {
	if e == nil {
		return nil
	}
	return e.next
}

func (e *lruEntry[K, V]) PushFront(prev *lruEntry[K, V]) {
	if e == nil {
		return
	}
	e.prev = prev
}

func (e *lruEntry[K, V]) PushBack(next *lruEntry[K, V]) {
	if e == nil {
		return
	}
	e.next = next
}

func (e *lruEntry[K, V]) GetData() tupleKV[K, V] {
	if e == nil {
		return tupleKV[K, V]{}
	}
	return *e.data
}

func (e *lruEntry[K, V]) Replace(newKey K, newVal V) tupleKV[K, V] {
	res := e.data.DeepCopy()
	e.data.key = pack.Pack[K, *pack.Key[K]](newKey)
	e.data.val = pack.Pack[V, *pack.Value[V]](newVal)

	return res
}
