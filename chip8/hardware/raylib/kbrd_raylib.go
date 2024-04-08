package raylib

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

/*
| 1 | 2 | 3 | C |
| 4 | 5 | 6 | D |
| 7 | 8 | 9 | E |
| A | 0 | B | F |
*/

var KEYS []int32 = []int32{
	rl.KeyKp0,        // 0 - "num 0"
	rl.KeyKp7,        // 1 - "num 7"
	rl.KeyKp8,        // 2 - "num 8"
	rl.KeyKp9,        // 3 - "num 9"
	rl.KeyKp4,        // 4 - "num 4"
	rl.KeyKp5,        // 5 - "num 5"
	rl.KeyKp6,        // 6 - "num 6"
	rl.KeyKp1,        // 7 - "num 1"
	rl.KeyKp2,        // 8 - "num 2"
	rl.KeyKp3,        // 9 - "num 3"
	rl.KeyKpDecimal,  // A - "num ."
	rl.KeyKpEnter,    // B - "num Enter"
	rl.KeyKpDivide,   // C - "num /"
	rl.KeyKpMultiply, // D - "num *"
	rl.KeyKpSubtract, // E - "num -"
	rl.KeyKpAdd,      // F - "num +"
}

type KeyboardRaylib struct {
	status uint16
}

func NewKeyboardRaylib() *KeyboardRaylib {
	return &KeyboardRaylib{status: 0}
}

func (kbrd *KeyboardRaylib) ReadKeys() uint16 {
	kbrd.status = 0
	for i, k := range KEYS {
		if rl.IsKeyDown(k) {
			kbrd.status |= (1 << i)
		}
	}
	return kbrd.status
}

func (kbrd *KeyboardRaylib) WaitKey() byte {
	for i, k := range KEYS {
		if rl.IsKeyReleased(k) {
			return byte(i)
		}
	}
	return 0x80
}
