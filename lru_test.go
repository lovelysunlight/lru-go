package lru

import (
	"testing"

	"github.com/JimChenWYU/lru-go/internal/option"
	"github.com/stretchr/testify/assert"
)

func assert_opt_eq[V any](t *testing.T, opt option.Option[V], v V) {
	assert.True(t, opt.IsSome())
	assert.EqualValues(t, v, opt.Unwrap())
}

func TestLruCache_Put_And_Get(t *testing.T) {
	cache := New[string, string](2)
	assert.EqualValues(t, option.None[string](), cache.Put("apple", "red"))
	assert.EqualValues(t, option.None[string](), cache.Put("banana", "yellow"))

	assert.EqualValues(t, 2, cache.Cap())
	assert.EqualValues(t, 2, cache.Len())

	assert_opt_eq(t, cache.Get("apple"), "red")
	assert_opt_eq(t, cache.Get("banana"), "yellow")
	assert.EqualValues(t, option.None[string](), cache.Get("orange"))
}

func TestLruCache_Put_And_Peek(t *testing.T) {
	cache := New[string, string](2)
	assert.EqualValues(t, option.None[string](), cache.Put("apple", "red"))
	assert.EqualValues(t, option.None[string](), cache.Put("banana", "yellow"))
	assert.EqualValues(t, 2, cache.Cap())
	assert.EqualValues(t, 2, cache.Len())

	assert_opt_eq(t, cache.Peek("apple"), "red")
	assert_opt_eq(t, cache.Peek("banana"), "yellow")

	assert.EqualValues(t, option.None[string](), cache.Peek("orange"))

	assert.EqualValues(t, option.Some("yellow"), cache.Put("banana", "foo"))
	assert.EqualValues(t, 2, cache.Cap())
	assert.EqualValues(t, 2, cache.Len())
}

func TestLruCache_Push_And_Peek(t *testing.T) {
	cache := New[string, string](2)
	assert.EqualValues(t, option.None[string](), cache.Put("apple", "red"))
	assert.EqualValues(t, option.None[string](), cache.Put("banana", "yellow"))

	assert.EqualValues(t, 2, cache.Cap())
	assert.EqualValues(t, 2, cache.Len())

	assert.EqualValues(t, option.Some(tupleKV[string, string]{
		Key: Key[string]{K: "banana"},
		Val: Value[string]{V: "yellow"},
	}), cache.Push("banana", "foo"))

	assert.EqualValues(t, option.Some(tupleKV[string, string]{
		Key: Key[string]{K: "apple"},
		Val: Value[string]{V: "red"},
	}), cache.Push("apple", "bar"))

	assert_opt_eq(t, cache.Peek("apple"), "bar")
	assert_opt_eq(t, cache.Peek("banana"), "foo")

	assert.EqualValues(t, option.Some(tupleKV[string, string]{
		Key: Key[string]{K: "banana"},
		Val: Value[string]{V: "foo"},
	}), cache.Push("orange", "orange"))

	assert.EqualValues(t, 2, cache.Cap())
	assert.EqualValues(t, 2, cache.Len())

	assert.EqualValues(t, option.None[string](), cache.Peek("banana"))
}

func TestLruCache_Pop(t *testing.T) {
	cache := New[string, string](2)
	cache.Put("apple", "red")
	cache.Put("banana", "yellow")

	assert.EqualValues(t, 2, cache.Len())
	assert_opt_eq(t, cache.Peek("apple"), "red")
	assert_opt_eq(t, cache.Peek("banana"), "yellow")

	poped := cache.Pop("apple")
	assert.True(t, poped.IsSome())
	assert.EqualValues(t, "red", poped.Unwrap())
	assert.EqualValues(t, 1, cache.Len())

	assert.False(t, cache.Peek("apple").IsSome())
	assert_opt_eq(t, cache.Peek("banana"), "yellow")

	poped = cache.Pop("banana")
	assert.True(t, poped.IsSome())
	assert.EqualValues(t, "yellow", poped.Unwrap())
	assert.EqualValues(t, 0, cache.Len())

	poped = cache.Pop("orange")
	assert.False(t, poped.IsSome())
	assert.EqualValues(t, 0, cache.Len())
}
