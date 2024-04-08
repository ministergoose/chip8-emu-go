package chip8

import "fmt"

const (
	STACK_SIZE = 16
)

type Stack interface {
	Reset()
	Push(addr uint16) error
	Pop() (uint16, error)
}

type StackStd struct {
	stack [STACK_SIZE]uint16
	index byte
}

func NewStackStd() *StackStd {
	return &StackStd{index: 0}
}

func (s *StackStd) Reset() {
	for i := 0; i < STACK_SIZE; i++ {
		s.stack[i] = 0
	}
	s.index = 0
}

func (s *StackStd) Push(addr uint16) error {
	if s.index >= STACK_SIZE-1 {
		return fmt.Errorf("stack overflow")
	}

	s.index++
	s.stack[s.index] = addr

	return nil
}

func (s *StackStd) Pop() (uint16, error) {
	if s.index == 0 {
		return 0, fmt.Errorf("stack is empty")
	}

	res := s.stack[s.index]
	s.index--

	return res, nil
}
