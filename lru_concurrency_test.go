package lru

import (
	"fmt"
	"math"
	"testing"

	"github.com/sourcegraph/conc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initData(cache *Cache[int, int], total int) {
	for i := 1; i <= total; i++ {
		cache.Put(i, i)
	}
}

func TestConcurrency_Put(t *testing.T) {
	max := 1 << 13
	cache, _ := New[int, int](max)

	wg := conc.NewWaitGroup()
	for i := 1; i <= max; i++ {
		wg.Go(func() {
			cache.Put(i, i)
		})
	}
	wg.Wait()

	for i := 1; i < max; i++ {
		v, ok := cache.Peek(i)
		assert.True(t, ok)
		assert.EqualValues(t, i, v)
	}

	assert.Len(t, cache.Keys(), max)
	assert.Len(t, cache.Values(), max)
}

func TestConcurrency_Get(t *testing.T) {
	max := 1 << 13
	var i float64 = 2
	tests := make([]int, 0, max)
	for {
		n := int(math.Pow(2, i))
		if n >= max {
			break
		}
		tests = append(tests, n)
		i++
	}

	for _, n := range tests {
		cache, _ := New[int, int](max)
		initData(cache, max)

		keys := cache.Keys()
		require.Len(t, keys, max)
		t.Run(fmt.Sprint("N = ", n), func(t *testing.T) {
			wg := conc.NewWaitGroup()
			for i := 1; i <= n; i++ {
				wg.Go(func() {
					v, ok := cache.Get(i)
					assert.True(t, ok)
					assert.EqualValues(t, i, v)
				})
			}
			wg.Wait()

			k, v, ok := cache.PeekOldest()
			assert.True(t, ok)
			assert.EqualValues(t, n+1, k)
			assert.EqualValues(t, n+1, v)
		})
	}
}

func TestConcurrency_Peek(t *testing.T) {
	max := 1 << 13
	var i float64 = 2
	tests := make([]int, 0, max)
	for {
		n := int(math.Pow(2, i))
		if n >= max {
			break
		}
		tests = append(tests, n)
		i++
	}

	for _, n := range tests {
		cache, _ := New[int, int](max)
		initData(cache, max)

		keys := cache.Keys()
		require.Len(t, keys, max)
		t.Run(fmt.Sprint("N = ", n), func(t *testing.T) {
			wg := conc.NewWaitGroup()
			for i := 1; i <= n; i++ {
				wg.Go(func() {
					v, ok := cache.Peek(i)
					assert.True(t, ok)
					assert.EqualValues(t, i, v)
				})
			}
			wg.Wait()

			k, v, ok := cache.PeekOldest()
			assert.True(t, ok)
			assert.EqualValues(t, 1, k)
			assert.EqualValues(t, 1, v)
		})
	}
}

func TestConcurrency_RemoveOldest(t *testing.T) {
	max := 1 << 13
	var i float64 = 2
	tests := make([]int, 0, max)
	for {
		n := int(math.Pow(2, i))
		if n >= max {
			break
		}
		tests = append(tests, n)
		i++
	}

	for _, n := range tests {
		cache, _ := New[int, int](max)
		initData(cache, max)

		keys := cache.Keys()
		require.Len(t, keys, max)
		t.Run(fmt.Sprint("N = ", n), func(t *testing.T) {
			wg := conc.NewWaitGroup()
			for i := 1; i <= n; i++ {
				wg.Go(func() {
					_, _, ok := cache.RemoveOldest()
					assert.True(t, ok)
				})
			}
			wg.Wait()

			k, v, ok := cache.PeekOldest()
			assert.True(t, ok)
			assert.EqualValues(t, n+1, k)
			assert.EqualValues(t, n+1, v)

			for expected := n + 1; cache.Len() > n; expected++ {
				k, v, ok := cache.RemoveOldest()
				assert.True(t, ok)
				assert.EqualValues(t, expected, k)
				assert.EqualValues(t, expected, v)
			}

			for cache.Len() > n {
				k, v, ok := cache.RemoveOldest()
				assert.True(t, ok)
				assert.LessOrEqual(t, k, n)
				assert.LessOrEqual(t, v, n)
			}
		})
	}
}

func TestConcurrency_Remove(t *testing.T) {
	max := 1 << 13
	var i float64 = 2
	tests := make([]int, 0, max)
	for {
		n := int(math.Pow(2, i))
		if n >= max {
			break
		}
		tests = append(tests, n)
		i++
	}

	for _, n := range tests {
		cache, _ := New[int, int](max)
		initData(cache, max)

		keys := cache.Keys()
		require.Len(t, keys, max)
		t.Run(fmt.Sprint("N = ", n), func(t *testing.T) {
			wg := conc.NewWaitGroup()
			for i := 1; i <= n; i++ {
				wg.Go(func() {
					v, ok := cache.Remove(i)
					assert.True(t, ok)
					assert.EqualValues(t, i, v)
				})
			}
			wg.Wait()

			k, v, ok := cache.PeekOldest()
			assert.True(t, ok)
			assert.EqualValues(t, n+1, k)
			assert.EqualValues(t, n+1, v)
		})
	}
}
