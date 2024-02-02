package iterator

type Iterator[T any] func() (T, bool)

func FilterMap[T, U any](it Iterator[T], f func(T) (U, bool)) Iterator[U] {
	return func() (out U, ok bool) {
		for {
			value, ok := it()
			if !ok {
				return out, false
			}

			y, ok := f(value)
			if ok {
				return y, true
			}
		}
	}
}

func (it Iterator[T]) Collect() []T {
	var values []T
	for {
		value, ok := it()
		if !ok {
			return values
		}

		values = append(values, value)
	}
}
