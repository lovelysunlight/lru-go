package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFIFO(t *testing.T) {
	list := NewFIFOList[int, int](3)
	assert.EqualValues(t, 3, list.Size())
	assert.EqualValues(t, 0, list.Len())
	assert.Equal(t, &list.root, list.Root())

	list.root.next = nil
	list.root.prev = nil
	list.lazyInit()
}

func TestFIFO_Resize(t *testing.T) {
	list := NewFIFOList[int, int](10)
	assert.EqualValues(t, 10, list.Size())
	assert.EqualValues(t, 0, list.Len())
	assert.EqualValues(t, 0, list.Resize(4))
	assert.EqualValues(t, 4, list.Size())

	list.Push(1, 1)
	list.Push(2, 2)
	list.Push(3, 3)
	list.Push(4, 4)

	assert.True(t, list.Contains(1))
	assert.EqualValues(t, 1, list.Resize(3))
	assert.EqualValues(t, 3, list.Size())
	assert.EqualValues(t, 3, list.Len())
	assert.False(t, list.Contains(1))
}

func TestFIFO_Push(t *testing.T) {
	list := NewFIFOList[int, int](3)
	list.Push(1, 1)
	assert.EqualValues(t, 1, list.Len())
	list.Push(2, 2)
	assert.EqualValues(t, 2, list.Len())
	list.Push(3, 3)
	assert.EqualValues(t, 3, list.Len())
	list.Push(4, 4)
	assert.EqualValues(t, 3, list.Len())
}

func TestFIFO_Pop(t *testing.T) {
	list := NewFIFOList[int, int](3)
	list.Push(1, 1)
	list.Push(2, 2)
	list.Push(3, 3)

	e := list.Pop()
	assert.EqualValues(t, 1, e.Key)
	assert.EqualValues(t, 1, e.Value)

	e = list.Pop()
	assert.EqualValues(t, 2, e.Key)
	assert.EqualValues(t, 2, e.Value)

	e = list.Pop()
	assert.EqualValues(t, 3, e.Key)
	assert.EqualValues(t, 3, e.Value)
}

func TestFIFO_Get(t *testing.T) {
	list := NewFIFOList[int, int](3)
	list.Push(1, 1)
	list.Push(2, 2)
	list.Push(3, 3)

	e, ok := list.Get(2)
	assert.True(t, ok)
	assert.EqualValues(t, 2, e.Key)
	assert.EqualValues(t, 2, e.Value)

	e, ok = list.Get(4)
	assert.False(t, ok)
	assert.Nil(t, e)
}

func TestFIFO_Remove(t *testing.T) {
	list := NewFIFOList[int, int](3)
	list.Push(1, 1)
	list.Push(2, 2)
	list.Push(3, 3)

	e, ok := list.Remove(2)
	assert.True(t, ok)
	assert.EqualValues(t, 2, e.Key)
	assert.EqualValues(t, 2, e.Value)

	e, ok = list.Remove(4)
	assert.False(t, ok)
	assert.Nil(t, e)

	assert.EqualValues(t, 2, list.Len())
}

func TestFIFO_RemoveElement(t *testing.T) {
	list := NewFIFOList[int, int](3)
	list.Push(1, 1)
	list.Push(2, 2)
	list.Push(3, 3)
	var (
		e *Entry[int, int]
	)
	e, _ = list.Get(1)
	list.RemoveElement(e)
	assert.EqualValues(t, 2, list.Len())

	e, _ = list.Get(4)
	list.RemoveElement(e)
	assert.EqualValues(t, 2, list.Len())
}

func TestFIFO_Back(t *testing.T) {
	var (
		e *Entry[int, int]
	)
	list := NewFIFOList[int, int](3)
	e = list.Back()
	assert.Nil(t, e)

	list.Push(1, 1)
	e = list.Back()
	assert.EqualValues(t, 1, e.Key)
	assert.EqualValues(t, 1, e.Value)

	list.Push(2, 2)
	e = list.Back()
	assert.EqualValues(t, 1, e.Key)
	assert.EqualValues(t, 1, e.Value)
}
