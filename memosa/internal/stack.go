package internal

type Stack[T any] struct {
	data []T
}

func (s *Stack[T]) Push(value T) {
	s.data = append(s.data, value)
}

func (s *Stack[T]) Peek() *T {
	if len(s.data) == 0 {
		return nil
	}

	return &s.data[len(s.data)-1]
}

func (s *Stack[T]) Pop() T {
	if len(s.data) == 0 {
		panic("empty stack")
	}
	value := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return value
}
