package hashmap

type Map[K comparable, V any] struct {
	inner map[K]V
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{inner: make(map[K]V)}
}

func (m *Map[K, V]) Get(k K) (V, bool) {
	v, ok := m.inner[k]
	return v, ok
}

func (m *Map[K, V]) Set(k K, v V) {
	m.inner[k] = v
}

func (m *Map[K, V]) Remove(k K) (V, bool) {
	if v, ok := m.Get(k); ok {
		delete(m.inner, k)
		return v, true
	}
	var v V
	return v, false
}

func (m *Map[K, V]) Clear() {
	m.inner = make(map[K]V)
}

func (m *Map[K, V]) Len() int {
	return len(m.inner)
}
