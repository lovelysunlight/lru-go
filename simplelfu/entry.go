package simplelfu

import "sync/atomic"

type LFUValue[T any] struct {
	value  T
	visits uint64
}

func (v *LFUValue[T]) SetValue(value T) {
	v.value = value
}

func (v *LFUValue[T]) Value() T {
	return v.value
}

func (v *LFUValue[T]) IncrVisits() {
	atomic.AddUint64(&v.visits, 1)
}

func (v *LFUValue[T]) GetVisits() uint64 {
	return atomic.LoadUint64(&v.visits)
}
