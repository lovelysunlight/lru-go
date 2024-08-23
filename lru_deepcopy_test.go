package lru

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache_disable_deepcopy_Peek(t *testing.T) {
	t.Run("map", func(t *testing.T) {
		cache, _ := New[string, map[string]string](2, DisableDeepCopy())
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
		cache, _ := New[string, []int](2, DisableDeepCopy())
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
		cache, _ := New[string, *TestCase](2, DisableDeepCopy())
		cache.Put("a", &TestCase{Name: "a"})

		v, _ := cache.Peek("a")
		assert.EqualValues(t, &TestCase{Name: "a"}, v)
		v.Name = "b"

		v, _ = cache.Peek("a")
		assert.EqualValues(t, &TestCase{Name: "b"}, v)
	})
}

func TestCache_enable_deepcopy_Peek(t *testing.T) {
	t.Run("map", func(t *testing.T) {
		cache, _ := New[string, map[string]string](2, EnableDeepCopy())
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
		cache, _ := New[string, []int](2, EnableDeepCopy())
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
		cache, _ := New[string, *TestCase](2, EnableDeepCopy())
		cache.Put("a", &TestCase{Name: "a"})

		v, _ := cache.Peek("a")
		assert.EqualValues(t, &TestCase{Name: "a"}, v)
		v.Name = "b"

		v, _ = cache.Peek("a")
		assert.EqualValues(t, &TestCase{Name: "a"}, v)
	})
}

func TestCache_enable_deepcopy_Get(t *testing.T) {
	t.Run("map", func(t *testing.T) {
		cache, _ := New[string, map[string]string](2, EnableDeepCopy())
		cache.Put("a", map[string]string{
			"a": "a",
		})

		v, _ := cache.Get("a")
		assert.EqualValues(t, map[string]string{
			"a": "a",
		}, v)
		v["a"] = "b"

		v, _ = cache.Get("a")
		assert.EqualValues(t, map[string]string{
			"a": "a",
		}, v)
	})
	t.Run("slice", func(t *testing.T) {
		cache, _ := New[string, []int](2, EnableDeepCopy())
		cache.Put("a", []int{1, 2, 3})

		v, _ := cache.Get("a")
		assert.EqualValues(t, []int{1, 2, 3}, v)
		v[0] = 4

		v, _ = cache.Get("a")
		assert.EqualValues(t, []int{1, 2, 3}, v)
	})
	t.Run("object", func(t *testing.T) {
		type TestCase struct {
			Name string
		}
		cache, _ := New[string, *TestCase](2, EnableDeepCopy())
		cache.Put("a", &TestCase{Name: "a"})

		v, _ := cache.Peek("a")
		assert.EqualValues(t, &TestCase{Name: "a"}, v)
		v.Name = "b"

		v, _ = cache.Peek("a")
		assert.EqualValues(t, &TestCase{Name: "a"}, v)
	})
}

func TestCache_disable_deepcopy_Get(t *testing.T) {
	t.Run("map", func(t *testing.T) {
		cache, _ := New[string, map[string]string](2, DisableDeepCopy())
		cache.Put("a", map[string]string{
			"a": "a",
		})

		v, _ := cache.Get("a")
		assert.EqualValues(t, map[string]string{
			"a": "a",
		}, v)
		v["a"] = "b"

		v, _ = cache.Get("a")
		assert.EqualValues(t, map[string]string{
			"a": "b",
		}, v)
	})
	t.Run("slice", func(t *testing.T) {
		cache, _ := New[string, []int](2, DisableDeepCopy())
		cache.Put("a", []int{1, 2, 3})

		v, _ := cache.Get("a")
		assert.EqualValues(t, []int{1, 2, 3}, v)
		v[0] = 4

		v, _ = cache.Get("a")
		assert.EqualValues(t, []int{4, 2, 3}, v)
	})
	t.Run("object", func(t *testing.T) {
		type TestCase struct {
			Name string
		}
		cache, _ := New[string, *TestCase](2, DisableDeepCopy())
		cache.Put("a", &TestCase{Name: "a"})

		v, _ := cache.Get("a")
		assert.EqualValues(t, &TestCase{Name: "a"}, v)
		v.Name = "b"

		v, _ = cache.Get("a")
		assert.EqualValues(t, &TestCase{Name: "b"}, v)
	})
}

func TestCache_enable_deepcopy_PeekOldest(t *testing.T) {
	t.Run("map", func(t *testing.T) {
		cache, _ := New[string, map[string]string](2, EnableDeepCopy())
		cache.Put("a", map[string]string{
			"a": "a",
		})

		_, v, _ := cache.PeekOldest()
		assert.EqualValues(t, map[string]string{
			"a": "a",
		}, v)
		v["a"] = "b"

		_, v, _ = cache.PeekOldest()
		assert.EqualValues(t, map[string]string{
			"a": "a",
		}, v)
	})
	t.Run("slice", func(t *testing.T) {
		cache, _ := New[string, []int](2, EnableDeepCopy())
		cache.Put("a", []int{1, 2, 3})

		_, v, _ := cache.PeekOldest()
		assert.EqualValues(t, []int{1, 2, 3}, v)
		v[0] = 4

		_, v, _ = cache.PeekOldest()
		assert.EqualValues(t, []int{1, 2, 3}, v)
	})
	t.Run("object", func(t *testing.T) {
		type TestCase struct {
			Name string
		}
		cache, _ := New[string, *TestCase](2, EnableDeepCopy())
		cache.Put("a", &TestCase{Name: "a"})

		_, v, _ := cache.PeekOldest()
		assert.EqualValues(t, &TestCase{Name: "a"}, v)
		v.Name = "b"

		_, v, _ = cache.PeekOldest()
		assert.EqualValues(t, &TestCase{Name: "a"}, v)
	})
}

func TestCache_disable_deepcopy_PeekOldest(t *testing.T) {
	t.Run("map", func(t *testing.T) {
		cache, _ := New[string, map[string]string](2, DisableDeepCopy())
		cache.Put("a", map[string]string{
			"a": "a",
		})

		_, v, _ := cache.PeekOldest()
		assert.EqualValues(t, map[string]string{
			"a": "a",
		}, v)
		v["a"] = "b"

		_, v, _ = cache.PeekOldest()
		assert.EqualValues(t, map[string]string{
			"a": "b",
		}, v)
	})
	t.Run("slice", func(t *testing.T) {
		cache, _ := New[string, []int](2, DisableDeepCopy())
		cache.Put("a", []int{1, 2, 3})

		_, v, _ := cache.PeekOldest()
		assert.EqualValues(t, []int{1, 2, 3}, v)
		v[0] = 4

		_, v, _ = cache.PeekOldest()
		assert.EqualValues(t, []int{4, 2, 3}, v)
	})
	t.Run("object", func(t *testing.T) {
		type TestCase struct {
			Name string
		}
		cache, _ := New[string, *TestCase](2, DisableDeepCopy())
		cache.Put("a", &TestCase{Name: "a"})

		_, v, _ := cache.PeekOldest()
		assert.EqualValues(t, &TestCase{Name: "a"}, v)
		v.Name = "b"

		_, v, _ = cache.PeekOldest()
		assert.EqualValues(t, &TestCase{Name: "b"}, v)
	})
}

func TestCache_enable_deepcopy_Values(t *testing.T) {
	t.Run("map", func(t *testing.T) {
		cache, _ := New[string, map[string]string](2, EnableDeepCopy())
		cache.Put("a", map[string]string{
			"a": "a",
		})

		v := cache.Values()
		assert.EqualValues(t, []map[string]string{
			{"a": "a"},
		}, v)
		v[0] = map[string]string{"a": "b"}

		v = cache.Values()
		assert.EqualValues(t, []map[string]string{
			{"a": "a"},
		}, v)
	})
	t.Run("slice", func(t *testing.T) {
		cache, _ := New[string, []int](2, EnableDeepCopy())
		cache.Put("a", []int{1, 2, 3})

		v := cache.Values()
		assert.EqualValues(t, [][]int{{1, 2, 3}}, v)
		v[0][0] = 4

		v = cache.Values()
		assert.EqualValues(t, [][]int{{1, 2, 3}}, v)
	})
	t.Run("object", func(t *testing.T) {
		type TestCase struct {
			Name string
		}
		cache, _ := New[string, *TestCase](2, EnableDeepCopy())
		cache.Put("a", &TestCase{Name: "a"})

		v := cache.Values()
		assert.EqualValues(t, []*TestCase{{Name: "a"}}, v)
		v[0].Name = "b"

		v = cache.Values()
		assert.EqualValues(t, []*TestCase{{Name: "a"}}, v)
	})
}

func TestCache_disable_deepcopy_Values(t *testing.T) {
	t.Run("map", func(t *testing.T) {
		cache, _ := New[string, map[string]string](2, DisableDeepCopy())
		cache.Put("a", map[string]string{
			"a": "a",
		})

		v := cache.Values()
		assert.EqualValues(t, []map[string]string{
			{"a": "a"},
		}, v)
		v[0]["a"] = "b"

		v = cache.Values()
		assert.EqualValues(t, []map[string]string{
			{"a": "b"},
		}, v)
		v[0] = map[string]string{"c": "c"}

		v = cache.Values()
		assert.EqualValues(t, []map[string]string{
			{"a": "b"},
		}, v)
	})
	t.Run("slice", func(t *testing.T) {
		cache, _ := New[string, []int](2, DisableDeepCopy())
		cache.Put("a", []int{1, 2, 3})

		v := cache.Values()
		assert.EqualValues(t, [][]int{{1, 2, 3}}, v)
		v[0][0] = 4

		v = cache.Values()
		assert.EqualValues(t, [][]int{{4, 2, 3}}, v)
		v[0] = []int{1}

		v = cache.Values()
		assert.EqualValues(t, [][]int{{4, 2, 3}}, v)
	})
	t.Run("object", func(t *testing.T) {
		type TestCase struct {
			Name string
		}
		cache, _ := New[string, *TestCase](2, DisableDeepCopy())
		cache.Put("a", &TestCase{Name: "a"})

		v := cache.Values()
		assert.EqualValues(t, []*TestCase{{Name: "a"}}, v)
		v[0].Name = "b"

		v = cache.Values()
		assert.EqualValues(t, []*TestCase{{Name: "b"}}, v)
		v[0] = &TestCase{Name: "c"}

		v = cache.Values()
		assert.EqualValues(t, []*TestCase{{Name: "b"}}, v)
	})
}

func TestCache_enable_deepcopy_Keys(t *testing.T) {
	t.Run("object", func(t *testing.T) {
		type TestCase struct {
			Name string
		}
		cache, _ := New[*TestCase, *TestCase](2, EnableDeepCopy())
		cache.Put(&TestCase{Name: "b"}, &TestCase{Name: "a"})

		v := cache.Keys()
		assert.EqualValues(t, []*TestCase{{Name: "b"}}, v)
		v[0].Name = "c"

		v = cache.Keys()
		assert.EqualValues(t, []*TestCase{{Name: "b"}}, v)
	})
}

func TestCache_disable_deepcopy_Keys(t *testing.T) {
	t.Run("object", func(t *testing.T) {
		type TestCase struct {
			Name string
		}
		cache, _ := New[*TestCase, *TestCase](2, DisableDeepCopy())
		cache.Put(&TestCase{Name: "b"}, &TestCase{Name: "a"})

		v := cache.Keys()
		assert.EqualValues(t, []*TestCase{{Name: "b"}}, v)
		v[0].Name = "c"

		v = cache.Keys()
		assert.EqualValues(t, []*TestCase{{Name: "c"}}, v)
	})
}
