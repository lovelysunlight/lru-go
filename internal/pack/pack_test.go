package pack

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	k := Pack[string, *Key[string]]("foo")
	assert.EqualValues(t, "foo", Unpack(k))
	assert.EqualValues(t, k, DeepCopy(k))

	copy := k
	assert.Equal(t, unsafe.Pointer(k), unsafe.Pointer(copy))
	assert.NotEqual(t, unsafe.Pointer(k), unsafe.Pointer(DeepCopy(k)))
}

func TestVal(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		k := Pack[string, *Value[string]]("foo")
		assert.EqualValues(t, "foo", Unpack(k))
		assert.EqualValues(t, k, DeepCopy(k))

		copy := k
		assert.Equal(t, unsafe.Pointer(k), unsafe.Pointer(copy))
		assert.NotEqual(t, unsafe.Pointer(k), unsafe.Pointer(DeepCopy(k)))
	})

	t.Run("struct", func(t *testing.T) {
		type TestCase struct {
			Name string
		}
		k := Pack[*TestCase, *Value[*TestCase]](&TestCase{
			Name: "foo",
		})
		assert.EqualValues(t, &TestCase{
			Name: "foo",
		}, Unpack(k))
		assert.EqualValues(t, k, DeepCopy(k))

		copy := k
		assert.Equal(t, unsafe.Pointer(k), unsafe.Pointer(copy))
		assert.NotEqual(t, unsafe.Pointer(k), unsafe.Pointer(DeepCopy(k)))
	})
}
