package simplelfu

import "sync/atomic"

type LFUValue[T any] struct {
	value T
	visit uint64
}

func (v *LFUValue[T]) SetValue(value T) {
	v.value = value
}

func (v *LFUValue[T]) Value() T {
	return v.value
}

func (v *LFUValue[T]) IncrVisit() {
	atomic.AddUint64(&v.visit, 1)
}

func (v *LFUValue[T]) GetVisit() uint64 {
	return atomic.LoadUint64(&v.visit)
}
