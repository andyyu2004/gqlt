package slice

func Map[T, U any](xs []T, f func(T) U) []U {
	ys := make([]U, len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

func Filter[T any](xs []T, f func(T) bool) []T {
	ys := []T{}
	for _, x := range xs {
		if f(x) {
			ys = append(ys, x)
		}
	}
	return ys
}

func FilterMap[T, U any](xs []T, f func(T) (U, bool)) []U {
	ys := []U{}
	for _, x := range xs {
		y, ok := f(x)
		if ok {
			ys = append(ys, y)
		}
	}
	return ys
}
