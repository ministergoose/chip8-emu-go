package hardware

const (
	DISPLAY_WIDTH  = 64
	DISPLAY_HEIGHT = 32
)

type Display interface {
	Init(title string, scale float32)
	PutPixel(x, y byte) bool
	Draw()
	Cls()
	Dump()
	ShouldClose() bool
	Close()
}
