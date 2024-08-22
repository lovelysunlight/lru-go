package option

type Option[T any] struct {
	data T
	ok   bool
}

func Some[T any](data T) Option[T] {
	return Option[T]{
		data: data,
		ok:   true,
	}
}

func None[T any]() Option[T] {
	return Option[T]{
		ok: false,
	}
}

func (o Option[T]) IsSome() bool {
	return o.ok
}

func (o Option[T]) Unwrap() T {
	if !o.ok {
		panic("called `Option.Unwrap()` on a `None` value")
	}

	return o.data
}

func (o Option[T]) Filter(predicate func(T) bool) Option[T] {
	if o.IsSome() {
		if predicate(o.data) {
			return Some(o.data)
		}
	}
	return None[T]()
}
