package pack

import (
	"github.com/JimChenWYU/lru-go/internal/deepcopy"
)

type Packable[T any] interface {
	Unpack() T
	Pack(T) Packable[T]
	DeepCopy() Packable[T]
}

type Key[T comparable] struct {
	inner T
}

var _ Packable[any] = (*Key[any])(nil)

func (v *Key[T]) Pack(d T) Packable[T] {
	return &Key[T]{
		inner: d,
	}
}

func (v *Key[T]) Unpack() T {
	return v.inner
}

func (v *Key[T]) DeepCopy() Packable[T] {
	return deepcopy.Copy(v)
}

// -------------------------------------------------------

type Value[T any] struct {
	inner T
}

var _ Packable[any] = (*Value[any])(nil)

func (v *Value[T]) Pack(d T) Packable[T] {
	return &Value[T]{
		inner: d,
	}
}

func (v *Value[T]) Unpack() T {
	return v.inner
}

func (v *Value[T]) DeepCopy() Packable[T] {
	return deepcopy.Copy(v)
}

func Unpack[T any, V Packable[T]](v V) T {
	return v.Unpack()
}

func Pack[T any, P Packable[T]](v T) P {
	var p P
	return p.Pack(v).(P)
}

func DeepCopy[T any, P Packable[T]](v P) P {
	return v.DeepCopy().(P)
}
