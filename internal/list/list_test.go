package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLruList_lazyInit(t *testing.T) {
	l := NewList[int, int]()
	assert.EqualValues(t, 0, l.Length())
	l.root.next = nil
	l.root.prev = nil

	first := l.PushFront(1, 11)
	assert.EqualValues(t, 1, first.Key)
	assert.EqualValues(t, 11, first.Value)
	assert.Equal(t, first, l.root.next)
	assert.EqualValues(t, 1, l.Length())

	second := l.PushFront(2, 22)
	assert.EqualValues(t, 2, second.Key)
	assert.EqualValues(t, 22, second.Value)
	assert.Equal(t, second, l.root.next)
	assert.Equal(t, first, l.root.next.next)
	assert.EqualValues(t, 2, l.Length())
}

func TestLruList_PushFront(t *testing.T) {
	l := NewList[int, int]()
	assert.EqualValues(t, 0, l.Length())

	first := l.PushFront(1, 11)
	assert.EqualValues(t, 1, first.Key)
	assert.EqualValues(t, 11, first.Value)
	assert.Equal(t, first, l.root.next)
	assert.EqualValues(t, 1, l.Length())

	second := l.PushFront(2, 22)
	assert.Nil(t, second.PrevEntry(), nil) // forbid root
	assert.Equal(t, first.PrevEntry(), second)

	assert.EqualValues(t, 2, second.Key)
	assert.EqualValues(t, 22, second.Value)
	assert.Equal(t, second, l.root.next)
	assert.Equal(t, first, l.root.next.next)
	assert.EqualValues(t, 2, l.Length())
}

func TestLruList_Remove(t *testing.T) {
	l := NewList[int, int]()
	first := l.PushFront(1, 11)
	second := l.PushFront(2, 22)

	removed := l.Remove(first)
	assert.EqualValues(t, 11, removed)
	assert.Equal(t, second, l.root.next)

	removed = l.Remove(second)
	assert.EqualValues(t, 22, removed)
	assert.Equal(t, &l.root, l.root.next)
	assert.Equal(t, &l.root, l.root.prev)
}

func TestLruList_Back(t *testing.T) {
	l := NewList[int, int]()
	assert.Nil(t, l.Back())

	first := l.PushFront(1, 11)
	_ = l.PushFront(2, 22)
	assert.Equal(t, first, l.Back())
}

func TestLruList_MoveToFront(t *testing.T) {
	t.Run("move root", func(t *testing.T) {
		l := NewList[int, int]()

		first := l.PushFront(1, 11)
		second := l.PushFront(2, 22)
		l.move(&l.root, &l.root)
		assert.Equal(t, second, l.root.next)
		assert.Equal(t, first, l.root.next.next)
	})
	t.Run("move entry", func(t *testing.T) {
		l := NewList[int, int]()

		first := l.PushFront(1, 11)
		second := l.PushFront(2, 22)
		assert.Equal(t, second, l.root.next)

		l.MoveToFront(second)
		assert.Equal(t, second, l.root.next)

		l.MoveToFront(first)
		assert.Equal(t, first, l.root.next)
	})
}
