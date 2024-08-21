package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap_Set_And_Get(t *testing.T) {
	m := New[string, string]()
	m.Set("a", "aa")
	m.Set("b", "bb")
	m.Set("c", "cc")
	assert.EqualValues(t, 3, m.Len())

	got, exists := m.Get("a")
	assert.True(t, exists)
	assert.EqualValues(t, "aa", got)

	got, exists = m.Get("b")
	assert.True(t, exists)
	assert.EqualValues(t, "bb", got)

	got, exists = m.Get("ff")
	assert.False(t, exists)
	assert.EqualValues(t, "", got)
}

func TestMap_Remove(t *testing.T) {
	m := New[string, string]()
	m.Set("a", "aa")
	m.Set("b", "bb")

	got, exists := m.Get("a")
	assert.True(t, exists)
	assert.EqualValues(t, "aa", got)

	got, exists = m.Get("b")
	assert.True(t, exists)
	assert.EqualValues(t, "bb", got)

	got, exists = m.Remove("a")
	assert.True(t, exists)
	assert.EqualValues(t, "aa", got)
	got, exists = m.Get("a")
	assert.False(t, exists)
	assert.EqualValues(t, "", got)

	got, exists = m.Remove("b")
	assert.True(t, exists)
	assert.EqualValues(t, "bb", got)
	got, exists = m.Get("b")
	assert.False(t, exists)
	assert.EqualValues(t, "", got)

	got, exists = m.Remove("c")
	assert.False(t, exists)
	assert.EqualValues(t, "", got)
}
