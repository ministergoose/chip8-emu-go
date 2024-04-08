package hardware

type Keyboard interface {
	WaitKey() byte
	ReadKeys() uint16
}
