package lru

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assert_opt_eq[V any](t *testing.T, ok bool, got, v V) {
	assert.True(t, ok)
	assert.EqualValues(t, v, got)
}

func TestLruCache_Put_And_Get(t *testing.T) {
	cache, _ := New[string, string](2)
	var (
		v  string
		ok bool
	)
	_, ok = cache.Put("apple", "red")
	assert.False(t, ok)

	_, ok = cache.Put("banana", "yellow")
	assert.False(t, ok)

	assert.EqualValues(t, 2, cache.Cap())
	assert.EqualValues(t, 2, cache.Len())

	v, ok = cache.Get("apple")
	assert_opt_eq(t, ok, v, "red")

	v, ok = cache.Get("banana")
	assert_opt_eq(t, ok, v, "yellow")

	_, ok = cache.Get("orange")
	assert.False(t, ok)
}

func TestLruCache_Put_And_Peek(t *testing.T) {
	cache, _ := New[string, string](2)
	var (
		v  string
		ok bool
	)
	_, ok = cache.Put("apple", "red")
	assert.False(t, ok)

	_, ok = cache.Put("banana", "yellow")
	assert.False(t, ok)
	assert.EqualValues(t, 2, cache.Cap())
	assert.EqualValues(t, 2, cache.Len())

	v, ok = cache.Peek("apple")
	assert_opt_eq(t, ok, v, "red")

	v, ok = cache.Peek("banana")
	assert_opt_eq(t, ok, v, "yellow")

	_, ok = cache.Peek("orange")
	assert.False(t, ok)

	v, ok = cache.Put("banana", "foo")
	assert_opt_eq(t, ok, v, "yellow")

	assert.EqualValues(t, 2, cache.Cap())
	assert.EqualValues(t, 2, cache.Len())
}

func TestLruCache_Push_And_Peek(t *testing.T) {
	cache, _ := New[string, string](2)
	var (
		k, v string
		ok   bool
	)
	_, ok = cache.Put("apple", "red")
	assert.False(t, ok)

	_, ok = cache.Put("banana", "yellow")
	assert.False(t, ok)
	assert.EqualValues(t, 2, cache.Cap())
	assert.EqualValues(t, 2, cache.Len())

	k, v, ok = cache.Push("banana", "foo")
	assert_opt_eq(t, ok, k, "banana")
	assert_opt_eq(t, ok, v, "yellow")

	k, v, ok = cache.Push("apple", "bar")
	assert_opt_eq(t, ok, k, "apple")
	assert_opt_eq(t, ok, v, "red")

	v, ok = cache.Peek("apple")
	assert_opt_eq(t, ok, v, "bar")

	v, ok = cache.Peek("banana")
	assert_opt_eq(t, ok, v, "foo")

	k, v, ok = cache.Push("orange", "orange")
	assert_opt_eq(t, ok, k, "banana")
	assert_opt_eq(t, ok, v, "foo")

	assert.EqualValues(t, 2, cache.Cap())
	assert.EqualValues(t, 2, cache.Len())

	_, ok = cache.Peek("banana")
	assert.False(t, ok)
}

func TestLruCache_Pop(t *testing.T) {
	cache, _ := New[string, string](2)
	cache.Put("apple", "red")
	cache.Put("banana", "yellow")
	assert.EqualValues(t, 2, cache.Len())

	var (
		v  string
		ok bool
	)

	v, ok = cache.Peek("apple")
	assert_opt_eq(t, ok, v, "red")
	v, ok = cache.Peek("banana")
	assert_opt_eq(t, ok, v, "yellow")

	v, ok = cache.Pop("apple")
	assert_opt_eq(t, ok, v, "red")
	assert.EqualValues(t, 1, cache.Len())
	_, ok = cache.Peek("apple")
	assert.False(t, ok)

	v, ok = cache.Peek("banana")
	assert_opt_eq(t, ok, v, "yellow")

	v, ok = cache.Pop("banana")
	assert_opt_eq(t, ok, v, "yellow")
	assert.EqualValues(t, 0, cache.Len())
	_, ok = cache.Peek("banana")
	assert.False(t, ok)

	_, ok = cache.Pop("orange")
	assert.False(t, ok)
	assert.EqualValues(t, 0, cache.Len())
}

func TestCache_RemoveOldest(t *testing.T) {
	cache, _ := New[string, string](2)
	cache.Put("apple", "red")
	cache.Put("banana", "yellow")
	assert.EqualValues(t, 2, cache.Len())

	var (
		k, v string
		ok   bool
	)
	k, v, ok = cache.RemoveOldest()
	assert_opt_eq(t, ok, k, "apple")
	assert_opt_eq(t, ok, v, "red")
	assert.EqualValues(t, 1, cache.Len())

	k, v, ok = cache.RemoveOldest()
	assert_opt_eq(t, ok, k, "banana")
	assert_opt_eq(t, ok, v, "yellow")
	assert.EqualValues(t, 0, cache.Len())

	cache.Put("apple", "red")
	cache.Put("banana", "yellow")
	_, _ = cache.Get("apple")

	k, v, ok = cache.RemoveOldest()
	assert_opt_eq(t, ok, k, "banana")
	assert_opt_eq(t, ok, v, "yellow")

	k, v, ok = cache.RemoveOldest()
	assert_opt_eq(t, ok, k, "apple")
	assert_opt_eq(t, ok, v, "red")
}

func TestCache_Clear(t *testing.T) {
	cache, _ := New[string, string](2)
	cache.Put("apple", "red")
	cache.Put("banana", "yellow")

	assert.EqualValues(t, 2, cache.Len())

	cache.Clear()
	assert.EqualValues(t, 0, cache.Len())
}

func TestCache_mutable(t *testing.T) {
	t.Run("map", func(t *testing.T) {
		cache, _ := New[string, map[string]string](2, WithMutable())
		cache.Put("a", map[string]string{
			"a": "a",
		})

		v, _ := cache.Peek("a")
		assert.EqualValues(t, map[string]string{
			"a": "a",
		}, v)
		v["a"] = "b"

		v, _ = cache.Peek("a")
		assert.EqualValues(t, map[string]string{
			"a": "b",
		}, v)
	})
	t.Run("slice", func(t *testing.T) {
		cache, _ := New[string, []int](2, WithMutable())
		cache.Put("a", []int{1, 2, 3})

		v, _ := cache.Peek("a")
		assert.EqualValues(t, []int{1, 2, 3}, v)
		v[0] = 4

		v, _ = cache.Peek("a")
		assert.EqualValues(t, []int{4, 2, 3}, v)
	})
	t.Run("object", func(t *testing.T) {
		type TestCase struct {
			Name string
		}
		cache, _ := New[string, *TestCase](2, WithMutable())
		cache.Put("a", &TestCase{Name: "a"})

		v, _ := cache.Peek("a")
		assert.EqualValues(t, &TestCase{Name: "a"}, v)
		v.Name = "b"

		v, _ = cache.Peek("a")
		assert.EqualValues(t, &TestCase{Name: "b"}, v)
	})
}

func TestCache_immutable(t *testing.T) {
	t.Run("map", func(t *testing.T) {
		cache, _ := New[string, map[string]string](2, WithImmutable())
		cache.Put("a", map[string]string{
			"a": "a",
		})

		v, _ := cache.Peek("a")
		assert.EqualValues(t, map[string]string{
			"a": "a",
		}, v)
		v["a"] = "b"

		v, _ = cache.Peek("a")
		assert.EqualValues(t, map[string]string{
			"a": "a",
		}, v)
	})
	t.Run("slice", func(t *testing.T) {
		cache, _ := New[string, []int](2, WithImmutable())
		cache.Put("a", []int{1, 2, 3})

		v, _ := cache.Peek("a")
		assert.EqualValues(t, []int{1, 2, 3}, v)
		v[0] = 4

		v, _ = cache.Peek("a")
		assert.EqualValues(t, []int{1, 2, 3}, v)
	})
	t.Run("object", func(t *testing.T) {
		type TestCase struct {
			Name string
		}
		cache, _ := New[string, *TestCase](2, WithImmutable())
		cache.Put("a", &TestCase{Name: "a"})

		v, _ := cache.Peek("a")
		assert.EqualValues(t, &TestCase{Name: "a"}, v)
		v.Name = "b"

		v, _ = cache.Peek("a")
		assert.EqualValues(t, &TestCase{Name: "a"}, v)
	})
}