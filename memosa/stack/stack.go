package stack

type Stack[T any] struct {
	data []T
}

func (s *Stack[T]) Push(values ...T) {
	s.data = append(s.data, values...)
}

func (s *Stack[T]) Peek() (out T, ok bool) {
	if len(s.data) == 0 {
		return out, false
	}

	return s.data[len(s.data)-1], true
}

func (s *Stack[T]) Pop() (out T, ok bool) {
	if len(s.data) == 0 {
		return out, false
	}

	value := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return value, true
}

func (s *Stack[T]) MustPop() T {
	if len(s.data) == 0 {
		panic("empty stack")
	}
	value := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return value
}
