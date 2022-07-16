package dsa

type Stack[T any] struct {
	Items		[]T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Count() int {
	return len(s.Items)
}

func (s *Stack[T]) Pop() T {
	n := s.Count() - 1
	item := s.Items[n]

	s.Items = s.Items[:n]
	return item
}

func (s *Stack[T]) Push(item T) {
	s.Items = append(s.Items, item)
}

func (s *Stack[T]) Peek() T {
	n := s.Count() - 1
	item := s.Items[n]

	return item
}