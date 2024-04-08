package empty

type DisplayEmpty struct {
}

func NewDisplayEmpty() *DisplayEmpty {
	return &DisplayEmpty{}
}

func (dspl *DisplayEmpty) Init(title string, scale float32) {
}

func (dspl *DisplayEmpty) PutPixel(x, y byte) bool {
	return false
}

func (dspl *DisplayEmpty) Draw() {
}

func (dspl *DisplayEmpty) Cls() {

}

func (dspl *DisplayEmpty) Dump() {

}

func (dspl *DisplayEmpty) ShouldClose() bool {
	return false
}

func (dspl *DisplayEmpty) Close() {

}
