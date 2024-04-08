package empty

type KeyboardEmpty struct {
	status uint16
}

func NewKeyboardEmpty() *KeyboardEmpty {
	return &KeyboardEmpty{status: 0}
}

func (kbrd *KeyboardEmpty) ReadKeys() uint16 {
	kbrd.status = 0
	return kbrd.status
}

func (kbrd *KeyboardEmpty) WaitKey() byte {
	return 0
}
