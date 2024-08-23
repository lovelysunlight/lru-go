package simplelru

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	_, err := New[int, int](0)
	assert.Error(t, err)

	l, err := New[int, int](128)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 256; i++ {
		l.Push(i, i)
	}
	if l.Len() != 128 {
		t.Fatalf("bad len: %v", l.Len())
	}
	if l.Cap() != 128 {
		t.Fatalf("expect %d, but %d", 128, l.Cap())
	}

	for i := 0; i < 128; i++ {
		if _, ok := l.Get(i); ok {
			t.Fatalf("should be evicted")
		}
	}
	for i := 128; i < 256; i++ {
		if _, ok := l.Get(i); !ok {
			t.Fatalf("should not be evicted")
		}
	}
	for i := 128; i < 192; i++ {
		if _, ok := l.Remove(i); !ok {
			t.Fatalf("should be contained")
		}
		if _, ok := l.Remove(i); ok {
			t.Fatalf("should not be contained")
		}
		if _, ok := l.Get(i); ok {
			t.Fatalf("should be deleted")
		}
	}

	l.Get(192) // expect 192 to be last key in l.Keys()

	l.Clear()
	if l.Len() != 0 {
		t.Fatalf("bad len: %v", l.Len())
	}
	if _, ok := l.Get(200); ok {
		t.Fatalf("should contain nothing")
	}
}

func TestCache_Replace_Push(t *testing.T) {
	l, err := New[int, int](2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Push(1, 1)
	l.Push(2, 2)
	l.Push(1, 101)

	var (
		k, v int
		ok   bool
	)
	k, v, ok = l.RemoveOldest()
	assert.True(t, ok)
	assert.EqualValues(t, 2, k)
	assert.EqualValues(t, 2, v)

	k, v, ok = l.RemoveOldest()
	assert.True(t, ok)
	assert.EqualValues(t, 1, k)
	assert.EqualValues(t, 101, v)

	_, _, ok = l.RemoveOldest()
	assert.False(t, ok)
}

func TestCache_GetOldest_RemoveOldest(t *testing.T) {
	l, err := New[int, int](128)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for i := 0; i < 256; i++ {
		l.Put(i, i)
	}
	var (
		k  int
		ok bool
	)

	k, _, ok = l.RemoveOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != 128 {
		t.Fatalf("bad: %v", k)
	}

	k, _, ok = l.RemoveOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != 129 {
		t.Fatalf("bad: %v", k)
	}
}

// Test that Add returns true/false if an eviction occurred
func TestCache_Put(t *testing.T) {
	l, err := New[int, int](1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	var ok bool
	_, ok = l.Put(1, 1)
	assert.False(t, ok)

	_, ok = l.Put(2, 2)
	assert.False(t, ok)

	t.Run("check address", func(t *testing.T) {
		type testCase struct {
			key string
		}
		l, _ := New[int, *testCase](1)
		insert := &testCase{"a"}
		l.Put(1, insert)
		got, _ := l.Peek(1)
		assert.EqualValues(t, &testCase{"a"}, got)

		evict, _ := l.Put(1, &testCase{"b"})
		assert.EqualValues(t, unsafe.Pointer(insert), unsafe.Pointer(evict))
		assert.EqualValues(t, &testCase{"a"}, evict)

		got, _ = l.Peek(1)
		assert.EqualValues(t, &testCase{"b"}, got)
	})
}

// Test that Peek doesn't update recent-ness
func TestCache_Peek(t *testing.T) {
	l, err := New[int, int](2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Put(1, 1)
	l.Put(2, 2)
	if v, ok := l.Peek(1); !ok || v != 1 {
		t.Errorf("1 should be set to 1: %v, %v", v, ok)
	}

	l.Put(3, 3)
	if _, ok := l.Peek(1); ok {
		t.Errorf("should not have updated recent-ness of 1")
	}
}

func TestCache_PeekOldest(t *testing.T) {
	l, err := New[int, int](3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	var (
		k, v int
		ok   bool
	)
	l.Put(1, 1)
	l.Put(2, 2)
	l.Put(3, 3)

	for i := 0; i < 2; i++ {
		k, v, ok = l.PeekOldest()
		_ = k
		assert.True(t, ok)
		assert.EqualValues(t, 1, v)
	}
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
}
