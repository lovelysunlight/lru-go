package simplelfu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleLFU(t *testing.T) {
	var (
		k  int
		v  int
		ok bool
	)

	t.Run("error", func(t *testing.T) {
		_, err := New[int, int](0)
		assert.Error(t, err)
	})

	t.Run("clear", func(t *testing.T) {
		cache, _ := New[int, int](3)
		cache.Push(1, 11)
		cache.Push(2, 22)
		cache.Push(3, 33)
		assert.EqualValues(t, 3, cache.Len())
		assert.EqualValues(t, 3, cache.Cap())

		cache.Clear()
		assert.EqualValues(t, 0, cache.Len())
		assert.EqualValues(t, 3, cache.Cap())
	})

	t.Run("keys/values", func(t *testing.T) {
		cache, _ := New[int, int](3)
		cache.Push(1, 11)
		cache.Push(2, 22)
		cache.Push(3, 33)
		assert.EqualValues(t, []int{1, 2, 3}, cache.Keys())
		assert.EqualValues(t, []int{11, 22, 33}, cache.Values())
	})

	t.Run("contains", func(t *testing.T) {
		cache, _ := New[int, int](3)
		cache.Push(1, 11)
		cache.Push(2, 22)
		cache.Push(3, 33)
		assert.True(t, cache.Contains(1))
		assert.False(t, cache.Contains(4))
	})

	t.Run("peek/peekLeast", func(t *testing.T) {
		cache, _ := New[int, int](3)
		cache.Push(1, 11)
		cache.Push(1, 11)
		cache.Push(2, 22)
		cache.Push(3, 33)

		v, ok = cache.Peek(1)
		assert.True(t, ok)
		assert.EqualValues(t, 11, v)

		k, v, ok = cache.PeekLeast()
		assert.True(t, ok)
		assert.EqualValues(t, 2, k)
		assert.EqualValues(t, 22, v)

		_, ok = cache.Peek(4)
		assert.False(t, ok)

		cache.Clear()
		_, _, ok = cache.PeekLeast()
		assert.False(t, ok)
	})

	t.Run("remove/removeLeast", func(t *testing.T) {
		cache, _ := New[int, int](3)
		cache.Push(1, 11)
		cache.Push(2, 22)
		cache.Push(3, 33)

		v, ok = cache.Remove(1)
		assert.True(t, ok)
		assert.EqualValues(t, 11, v)

		_, ok = cache.Remove(1)
		assert.False(t, ok)

		k, v, ok = cache.RemoveLeast()
		assert.True(t, ok)
		assert.EqualValues(t, 2, k)
		assert.EqualValues(t, 22, v)

		cache.Clear()

		_, _, ok = cache.RemoveLeast()
		assert.False(t, ok)
	})

	t.Run("logic", func(t *testing.T) {
		cache, _ := New[int, int](3)
		cache.Push(1, 1)
		cache.Push(2, 2)
		cache.Push(3, 3)

		assert.EqualValues(t, []string{
			"root -> 3 -> 2 -> 1 -> root",
			"root <- 3 <- 2 <- 1 <- root",
		}, cache.evictList.Debug())
		k, v, ok = cache.PeekLeast()
		assert.True(t, ok)
		assert.Equal(t, 1, k)
		assert.Equal(t, 1, v)

		cache.Push(4, 4)
		assert.EqualValues(t, []string{
			"root -> 4 -> 3 -> 2 -> root",
			"root <- 4 <- 3 <- 2 <- root",
		}, cache.evictList.Debug())
		k, v, ok = cache.PeekLeast()
		assert.True(t, ok)
		assert.Equal(t, 2, k)
		assert.Equal(t, 2, v)

		cache.Push(2, 2)
		assert.EqualValues(t, []string{
			"root -> 2 -> 4 -> 3 -> root",
			"root <- 2 <- 4 <- 3 <- root",
		}, cache.evictList.Debug())
		k, v, ok = cache.PeekLeast()
		assert.True(t, ok)
		assert.Equal(t, 3, k)
		assert.Equal(t, 3, v)

		cache.Get(4)
		assert.EqualValues(t, []string{
			"root -> 4 -> 2 -> 3 -> root",
			"root <- 4 <- 2 <- 3 <- root",
		}, cache.evictList.Debug())
		k, v, ok = cache.PeekLeast()
		assert.True(t, ok)
		assert.Equal(t, 3, k)
		assert.Equal(t, 3, v)

		cache.Get(4)
		cache.Get(3)
		assert.EqualValues(t, []string{
			"root -> 4 -> 3 -> 2 -> root",
			"root <- 4 <- 3 <- 2 <- root",
		}, cache.evictList.Debug())
		k, v, ok = cache.PeekLeast()
		assert.True(t, ok)
		assert.Equal(t, 2, k)
		assert.Equal(t, 2, v)
	})
}

func TestSimpleLFU_Resize(t *testing.T) {
	cache, _ := New[int, int](3)
	assert.EqualValues(t, 3, cache.Cap())

	cache.Push(1, 11)
	cache.Push(2, 22)
	cache.Push(3, 33)

	evicted := cache.Resize(2)
	assert.EqualValues(t, 1, evicted)
	assert.EqualValues(t, 2, cache.Cap())

	evicted = cache.Resize(10)
	assert.EqualValues(t, 0, evicted)
	assert.EqualValues(t, 10, cache.Cap())
}

func TestSimpleLFU_Get(t *testing.T) {
	cache, _ := New[int, int](3)
	cache.Push(1, 11)
	cache.Push(2, 22)
	cache.Push(3, 33)

	var (
		v  int
		ok bool
	)

	_, v, _ = cache.PeekLeast()
	assert.Equal(t, 11, v)

	v, ok = cache.Get(1)
	assert.True(t, ok)
	assert.Equal(t, 11, v)

	_, v, _ = cache.PeekLeast()
	assert.Equal(t, 22, v)

	_, ok = cache.Get(44)
	assert.False(t, ok)

	_, v, _ = cache.PeekLeast()
	assert.Equal(t, 22, v)
}

func TestSimpleLFU_Push(t *testing.T) {
	cache, _ := New[int, int](2)
	var (
		k, v int
		ok   bool
	)
	_, _, ok = cache.Push(1, 11)
	assert.False(t, ok)

	k, v, ok = cache.Push(1, 111)
	assert.True(t, ok)
	assert.EqualValues(t, 1, k)
	assert.EqualValues(t, 11, v)

	k, v, ok = cache.Push(1, 11)
	assert.True(t, ok)
	assert.EqualValues(t, 1, k)
	assert.EqualValues(t, 111, v)

	_, _, ok = cache.Push(2, 22)
	assert.False(t, ok)

	k, v, ok = cache.Push(3, 33)
	assert.True(t, ok)
	assert.EqualValues(t, 2, k)
	assert.EqualValues(t, 22, v)
}

func TestSimpleLFU_Put(t *testing.T) {
	cache, _ := New[int, int](2)
	var (
		v  int
		ok bool
	)
	_, ok = cache.Put(1, 11)
	assert.False(t, ok)

	v, ok = cache.Put(1, 111)
	assert.True(t, ok)
	assert.EqualValues(t, 11, v)

	v, ok = cache.Put(1, 11)
	assert.True(t, ok)
	assert.EqualValues(t, 111, v)

	_, ok = cache.Put(2, 22)
	assert.False(t, ok)

	_, ok = cache.Put(3, 33)
	assert.False(t, ok)
}

func TestSimpleLFU_PeekUsed(t *testing.T) {
	var (
		v  uint64
		ok bool
	)
	cache, _ := New[int, int](2)
	cache.Put(1, 1)

	v, _ = cache.PeekUsed(1)
	assert.EqualValues(t, 1, v)

	cache.Put(1, 1)
	v, _ = cache.PeekUsed(1)
	assert.EqualValues(t, 2, v)

	cache.Put(1, 11)
	v, _ = cache.PeekUsed(1)
	assert.EqualValues(t, 3, v)

	cache.Get(1)
	v, _ = cache.PeekUsed(1)
	assert.EqualValues(t, 4, v)

	v, ok = cache.PeekUsed(2)
	assert.False(t, ok)
	assert.EqualValues(t, 0, v)
}
