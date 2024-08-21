package lru

import "github.com/JimChenWYU/lru-go/internal/deepcopy"

type Key[T comparable] struct {
	inner T
}

func (v Key[T]) DeepCopy() Key[T] {
	return Key[T]{
		inner: deepcopy.Copy(v.inner),
	}
}

func (v Key[T]) Get() T {
	return v.inner
}

type Value[T any] struct {
	inner T
}

func (v Value[T]) DeepCopy() Value[T] {
	return Value[T]{
		inner: deepcopy.Copy(v.inner),
	}
}

func (v Value[T]) Get() T {
	return v.inner
}

type tupleKV[K comparable, V any] struct {
	key Key[K]
	val Value[V]
}

func (v tupleKV[K, V]) GetVal() Value[V] {
	return v.val
}

func (v tupleKV[K, V]) GetKey() Key[K] {
	return v.key
}

func (v tupleKV[K, V]) DeepCopy() tupleKV[K, V] {
	return tupleKV[K, V]{
		key: v.key.DeepCopy(),
		val: v.val.DeepCopy(),
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
			key: Key[K]{inner: key},
			val: Value[V]{inner: val},
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
	e.data.key = Key[K]{inner: newKey}
	e.data.val = Value[V]{inner: newVal}

	return res
}
