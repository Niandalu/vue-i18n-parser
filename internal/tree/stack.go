package tree

import (
	"errors"
)

type StackItem struct {
	k string
	v interface{}
}

type Stack struct {
	data []StackItem
}

func NewStack() *Stack {
	return new(Stack)
}

func (s *Stack) Push(record StackItem) {
	s.data = append(s.data, record)
}

func (s *Stack) Pop() (StackItem, error) {
	l := len(s.data)

	if l == 0 {
		return StackItem{}, errors.New("Empty Stack")
	}

	result := s.data[l-1]
	s.data = s.data[:l-1]
	return result, nil
}
