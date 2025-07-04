package util

import "errors"

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, error) {
	if len(s.items) == 0 {
		var zeroValue T
		return zeroValue, errors.New("stack is empty")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, nil
}

func (s *Stack[T]) Top() (T, error) {
	if len(s.items) == 0 {
		var zeroValue T
		return zeroValue, errors.New("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}
