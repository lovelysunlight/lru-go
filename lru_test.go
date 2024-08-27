package lru

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assert_opt_eq[V any](t *testing.T, ok bool, got, v V) {
	assert.True(t, ok)
	assert.EqualValues(t, v, got)
}

func TestLruCache_Error(t *testing.T) {
	_, err := New[int, int](0)
	assert.Error(t, err)
}

func TestLruCache_Resize(t *testing.T) {
	cache, _ := New[int, int](3)
	assert.EqualValues(t, 3, cache.Cap())

	cache.Push(1, 1)
	cache.Push(2, 2)
	cache.Push(3, 3)
	assert.EqualValues(t, 1, cache.Resize(2))
	assert.EqualValues(t, 2, cache.Cap())
	assert.EqualValues(t, 2, cache.Len())
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

func TestLruCache_Remove(t *testing.T) {
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

	v, ok = cache.Remove("apple")
	assert_opt_eq(t, ok, v, "red")
	assert.EqualValues(t, 1, cache.Len())
	_, ok = cache.Peek("apple")
	assert.False(t, ok)

	v, ok = cache.Peek("banana")
	assert_opt_eq(t, ok, v, "yellow")

	v, ok = cache.Remove("banana")
	assert_opt_eq(t, ok, v, "yellow")
	assert.EqualValues(t, 0, cache.Len())
	_, ok = cache.Peek("banana")
	assert.False(t, ok)

	_, ok = cache.Remove("orange")
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

func TestCache_PeekOldest(t *testing.T) {
	var (
		k, v int
		ok   bool
	)

	t.Run("peek oldest", func(t *testing.T) {
		l, err := New[int, int](3)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		l.Put(1, 1)
		l.Put(2, 2)
		l.Put(3, 3)

		for i := 0; i < 3; i++ {
			k, v, ok = l.PeekOldest()
			_ = k
			assert.True(t, ok)
			assert.EqualValues(t, 1, v)
		}
	})

	t.Run("peek oldest no value", func(t *testing.T) {
		l, _ := New[int, int](3)
		_, _, ok = l.PeekOldest()
		assert.False(t, ok)
	})
}

func TestCache_Keys_Values(t *testing.T) {
	l, err := New[int, int](3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	l.Put(1, 10)
	l.Put(2, 20)
	l.Put(3, 30)

	assert.EqualValues(t, []int{1, 2, 3}, l.Keys())
	assert.EqualValues(t, []int{10, 20, 30}, l.Values())
}

func TestCache_Contains(t *testing.T) {
	tests := []struct {
		name     string
		initData [][2]string
		key      string
		want     bool
	}{
		{
			name: "contains",
			initData: [][2]string{
				{"foo", "foo"},
				{"zoo", "zoo"},
			},
			key:  "foo",
			want: true,
		},
		{
			name: "not contains",
			initData: [][2]string{
				{"foo", "foo"},
				{"zoo", "zoo"},
			},
			key:  "bar",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := New[string, string](3)
			require.NoError(t, err)

			for _, data := range tt.initData {
				l.Push(data[0], data[1])
			}
			assert.EqualValues(t, tt.want, l.Contains(tt.key))
		})
	}

	t.Run("contains without updating the recent-ness.", func(t *testing.T) {
		l, _ := New[string, string](3)
		l.Put("foo", "foo")
		l.Put("zoo", "zoo")

		k, v, ok := l.PeekOldest()
		assert.True(t, ok)
		assert.EqualValues(t, "foo", k)
		assert.EqualValues(t, "foo", v)

		assert.True(t, l.Contains("foo"))

		k, v, ok = l.PeekOldest()
		assert.True(t, ok)
		assert.EqualValues(t, "foo", k)
		assert.EqualValues(t, "foo", v)
	})
}

func TestCacheUpgradeToLRUK_Push(t *testing.T) {
	var (
		// k, v int
		ok bool
	)
	cache, _ := New(3,
		WithVisitCacheSize[int, int](4),
		EnableLRUK[int, int](3),
	)
	assert.EqualValues(t, 3, cache.Cap())
	assert.EqualValues(t, 4, cache.VisitsCap())
	assert.EqualValues(t, 3, cache.visitThreshold)

	t.Run("push to lfu and evict", func(t *testing.T) {
		cache.VisitsResize(3)
		_, _, ok = cache.Push(1, 1)
		assert.False(t, ok)
		_, _, ok = cache.Push(2, 2)
		assert.False(t, ok)
		_, _, ok = cache.Push(3, 3)
		assert.False(t, ok)
		_, _, ok = cache.Push(4, 4)
		assert.False(t, ok)
	})

	t.Run("incr visits", func(t *testing.T) {
		cache, _ := New(3,
			WithVisitCacheSize[int, int](4),
			EnableLRUK[int, int](3),
		)

		cache.VisitsResize(3)
		cache.Push(1, 1)
		v, _ := cache.visit.PeekVisits(1)
		assert.EqualValues(t, 1, v)

		cache.Push(1, 1)
		v, _ = cache.visit.PeekVisits(1)
		assert.EqualValues(t, 2, v)

		cache.Push(1, 1)
		v, _ = cache.visit.PeekVisits(1)
		assert.EqualValues(t, 3, v)
	})
}

func TestCacheUpgradeToLRUK_Put(t *testing.T) {
	var (
		// v  int
		ok bool
	)
	cache, _ := New(3,
		WithVisitCacheSize[int, int](4),
		EnableLRUK[int, int](3),
	)
	assert.EqualValues(t, 3, cache.Cap())
	assert.EqualValues(t, 4, cache.VisitsCap())
	assert.EqualValues(t, 3, cache.visitThreshold)

	t.Run("put to lfu and evict", func(t *testing.T) {
		cache.Clear()
		cache.VisitsResize(3)
		_, ok = cache.Put(1, 1)
		assert.False(t, ok)
		_, ok = cache.Put(2, 2)
		assert.False(t, ok)
		_, ok = cache.Put(3, 3)
		assert.False(t, ok)
		_, ok = cache.Put(4, 4)
		assert.False(t, ok)
	})

	t.Run("incr visits", func(t *testing.T) {
		cache, _ := New(3,
			WithVisitCacheSize[int, int](4),
			EnableLRUK[int, int](3),
		)

		cache.VisitsResize(3)
		cache.Put(1, 1)
		v, _ := cache.visit.PeekVisits(1)
		assert.EqualValues(t, 1, v)

		cache.Put(1, 1)
		v, _ = cache.visit.PeekVisits(1)
		assert.EqualValues(t, 2, v)

		cache.Put(1, 1)
		v, _ = cache.visit.PeekVisits(1)
		assert.EqualValues(t, 3, v)
	})
}

func TestCacheUpgradeToLRUK_Get(t *testing.T) {
	var (
		v  int
		ok bool
	)
	cache, _ := New(3,
		WithVisitCacheSize[int, int](3),
		EnableLRUK[int, int](3),
	)
	cache.Push(1, 1)
	cache.Push(2, 2)
	cache.Push(2, 2)

	v, ok = cache.Get(1)
	assert.True(t, ok)
	assert.EqualValues(t, 1, v)
	assert.EqualValues(t, 0, cache.Len())
	assert.EqualValues(t, 2, cache.VisitsLen())

	v, ok = cache.Get(1) // move to lru
	assert.True(t, ok)
	assert.EqualValues(t, 1, v)
	assert.EqualValues(t, 1, cache.Len())
	assert.EqualValues(t, 1, cache.VisitsLen())

	v, ok = cache.Get(2) // move to lru
	assert.True(t, ok)
	assert.EqualValues(t, 2, v)
	assert.EqualValues(t, 2, cache.Len())
	assert.EqualValues(t, 0, cache.VisitsLen())
}

func TestCacheUpgradeTo2Q_FIFOResize(t *testing.T) {
	cache, _ := New(4,
		WithVisitCacheSize[int, int](4),
		Enable2Q[int, int](2),
	)
	assert.EqualValues(t, 2, cache.FIFOCap())
	cache.FIFOResize(3)
	assert.EqualValues(t, 3, cache.FIFOCap())
}

func TestCacheUpgradeTo2Q_Get(t *testing.T) {
	cache, _ := New(4,
		WithVisitCacheSize[int, int](4),
		Enable2Q[int, int](2),
	)

	var (
		v  int
		ok bool
	)
	cache.Push(1, 1)
	cache.Push(2, 2)
	assert.EqualValues(t, 0, cache.Len())
	assert.EqualValues(t, 2, cache.FIFOLen())

	_, ok = cache.Peek(1)
	assert.False(t, ok)

	_, ok = cache.Peek(2)
	assert.False(t, ok)

	v, ok = cache.Get(1)
	assert.True(t, ok)
	assert.EqualValues(t, 1, v)
	assert.EqualValues(t, 1, cache.Len())
	assert.EqualValues(t, 1, cache.FIFOLen())

	v, ok = cache.Get(2)
	assert.True(t, ok)
	assert.EqualValues(t, 2, v)
	assert.EqualValues(t, 2, cache.Len())
	assert.EqualValues(t, 0, cache.FIFOLen())
}

func TestCacheUpgradeTo2Q_Push_Put(t *testing.T) {
	t.Run("push", func(t *testing.T) {
		cache, _ := New(4,
			WithVisitCacheSize[int, int](4),
			Enable2Q[int, int](2),
		)
		var (
			v  int
			ok bool
		)
		cache.Push(1, 1)
		cache.Push(2, 2)
		assert.EqualValues(t, 0, cache.Len())
		assert.EqualValues(t, 2, cache.FIFOLen())

		_, ok = cache.Peek(1)
		assert.False(t, ok)

		_, ok = cache.Peek(2)
		assert.False(t, ok)

		v, ok = cache.Get(1)
		assert.True(t, ok)
		assert.EqualValues(t, 1, v)
		assert.EqualValues(t, 1, cache.Len())
		assert.EqualValues(t, 1, cache.FIFOLen())

		v, ok = cache.Get(2)
		assert.True(t, ok)
		assert.EqualValues(t, 2, v)
		assert.EqualValues(t, 2, cache.Len())
		assert.EqualValues(t, 0, cache.FIFOLen())
	})

	t.Run("put", func(t *testing.T) {
		cache, _ := New(4,
			WithVisitCacheSize[int, int](4),
			Enable2Q[int, int](2),
		)
		var (
			v  int
			ok bool
		)
		cache.Put(1, 1)
		cache.Put(2, 2)
		assert.EqualValues(t, 0, cache.Len())
		assert.EqualValues(t, 2, cache.FIFOLen())

		_, ok = cache.Peek(1)
		assert.False(t, ok)

		_, ok = cache.Peek(2)
		assert.False(t, ok)

		v, ok = cache.Get(1)
		assert.True(t, ok)
		assert.EqualValues(t, 1, v)
		assert.EqualValues(t, 1, cache.Len())
		assert.EqualValues(t, 1, cache.FIFOLen())

		v, ok = cache.Get(2)
		assert.True(t, ok)
		assert.EqualValues(t, 2, v)
		assert.EqualValues(t, 2, cache.Len())
		assert.EqualValues(t, 0, cache.FIFOLen())
	})
}
