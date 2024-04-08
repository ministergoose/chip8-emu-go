package empty

type StackEmpty struct {
}

func NewStackEmpty() *StackEmpty {
	return &StackEmpty{}
}

func (s *StackEmpty) Reset() {

}

func (s *StackEmpty) Push(addr uint16) error {
	return nil
}

func (s *StackEmpty) Pop() (uint16, error) {
	return 0, nil
}
