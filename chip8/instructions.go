package chip8

import (
	"fmt"
	"math/rand"

	"github.com/ministergoose/chip8-emu-go/chip8/hardware"
)

func getParameters(i uint16) (nnn uint16, kk byte, n byte, x byte, y byte) {
	nnn = i & 0x0fff
	kk = (byte)(i & 0x00ff)
	n = (byte)(i & 0x000f)
	x = (byte)((i & 0x0f00) >> 8)
	y = (byte)((i & 0x00f0) >> 4)
	return nnn, kk, n, x, y
}

func (cpu *Cpu) InstructionsInit() {
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x00e0, Mask: 0xffff, f: cpu.ins00e0})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x00ee, Mask: 0xffff, f: cpu.ins00ee})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x0000, Mask: 0xf000, f: cpu.ins0nnn})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x1000, Mask: 0xf000, f: cpu.ins1nnn})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x2000, Mask: 0xf000, f: cpu.ins2nnn})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x3000, Mask: 0xf000, f: cpu.ins3xkk})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x4000, Mask: 0xf000, f: cpu.ins4xkk})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x5000, Mask: 0xf00f, f: cpu.ins5xy0})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x6000, Mask: 0xf000, f: cpu.ins6xkk})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x7000, Mask: 0xf000, f: cpu.ins7xkk})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x8000, Mask: 0xf00f, f: cpu.ins8xy0})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x8001, Mask: 0xf00f, f: cpu.ins8xy1})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x8002, Mask: 0xf00f, f: cpu.ins8xy2})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x8003, Mask: 0xf00f, f: cpu.ins8xy3})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x8004, Mask: 0xf00f, f: cpu.ins8xy4})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x8005, Mask: 0xf00f, f: cpu.ins8xy5})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x8006, Mask: 0xf00f, f: cpu.ins8xy6})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x8007, Mask: 0xf00f, f: cpu.ins8xy7})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x800e, Mask: 0xf00f, f: cpu.ins8xyE})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0x9000, Mask: 0xf00f, f: cpu.ins9xy0})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0xa000, Mask: 0xf000, f: cpu.insAnnn})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0xb000, Mask: 0xf000, f: cpu.insBnnn})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0xc000, Mask: 0xf000, f: cpu.insCxkk})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0xd000, Mask: 0xf000, f: cpu.insDxyn})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0xe09e, Mask: 0xf0ff, f: cpu.insEx9E})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0xe0a1, Mask: 0xf0ff, f: cpu.insExA1})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0xf007, Mask: 0xf0ff, f: cpu.insFx07})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0xf00a, Mask: 0xf0ff, f: cpu.insFx0A})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0xf015, Mask: 0xf0ff, f: cpu.insFx15})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0xf018, Mask: 0xf0ff, f: cpu.insFx18})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0xf01e, Mask: 0xf0ff, f: cpu.insFx1E})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0xf029, Mask: 0xf0ff, f: cpu.insFx29})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0xf033, Mask: 0xf0ff, f: cpu.insFx33})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0xf055, Mask: 0xf0ff, f: cpu.insFx55})
	cpu.Opcodes = append(cpu.Opcodes, Opcode{Result: 0xf065, Mask: 0xf0ff, f: cpu.insFx65})
}

/*
   00E0 - CLS
   Clear the display.
*/

func (cpu *Cpu) ins00e0(op uint16) (string, error) {
	cpu.display.Cls()

	return "CLS\t\t; Clear the display", nil
}

/*
   00EE - RET
   Return from a subroutine.
   The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.
*/

func (cpu *Cpu) ins00ee(op uint16) (string, error) {
	addr, err := cpu.stack.Pop()
	if err != nil {
		return "", err
	}
	cpu.cnt = addr

	return "RET\t\t; Return from a subroutine", nil
}

/*
   0nnn - SYS addr
   Jump to a machine code routine at nnn.
   This instruction is only used on the old computers on which Chip-8 was originally implemented. It is ignored by modern interpreters.
*/

func (cpu *Cpu) ins0nnn(op uint16) (string, error) {
	nnn, _, _, _, _ := getParameters(op)
	// TODO
	return fmt.Sprintf("SYS 0x%03x", nnn), nil
}

/*
   1nnn - JP addr
   Jump to location nnn.
   The interpreter sets the program counter to nnn.
*/

func (cpu *Cpu) ins1nnn(op uint16) (string, error) {
	nnn, _, _, _, _ := getParameters(op)

	err := cpu.checkAddr(nnn)
	if err != nil {
		return "", err
	}
	cpu.cnt = nnn

	return fmt.Sprintf("JP 0x%03x\t\t; Jump to location nnn", nnn), nil
}

/*
   2nnn - CALL addr
   Call subroutine at nnn.
   The interpreter increments the stack pointer, then puts the current PC on the top of the stack. The PC is then set to nnn.
*/

func (cpu *Cpu) ins2nnn(op uint16) (string, error) {
	nnn, _, _, _, _ := getParameters(op)

	err := cpu.stack.Push(cpu.cnt)
	if err != nil {
		return "", err
	}
	cpu.cnt = nnn

	return fmt.Sprintf("CALL 0x%03x\t; Call subroutine at nnn", nnn), nil
}

/*
   3xkk - SE Vx, byte
   Skip next instruction if Vx = kk.
   The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
*/

func (cpu *Cpu) ins3xkk(op uint16) (string, error) {
	_, kk, _, x, _ := getParameters(op)

	if cpu.v[x] == kk {
		cpu.cnt += 2
	}

	return fmt.Sprintf("SE V%x, %02d\t\t; Skip next instruction if Vx = kk", x, kk), nil
}

/*
   4xkk - SNE Vx, byte
   Skip next instruction if Vx != kk.
   The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 2.
*/

func (cpu *Cpu) ins4xkk(op uint16) (string, error) {
	_, kk, _, x, _ := getParameters(op)

	if cpu.v[x] != kk {
		cpu.cnt += 2
	}

	return fmt.Sprintf("SNE V%x, %02d\t; Skip next instruction if Vx != kk", x, kk), nil
}

/*
   5xy0 - SE Vx, Vy
   Skip next instruction if Vx = Vy.
   The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.
*/

func (cpu *Cpu) ins5xy0(op uint16) (string, error) {
	_, _, _, x, y := getParameters(op)

	if cpu.v[x] == cpu.v[y] {
		cpu.cnt += 2
	}

	return fmt.Sprintf("SE V%x, V%x\t; Skip next instruction if Vx = Vy", x, y), nil
}

/*
   6xkk - LD Vx, byte
   Set Vx = kk.
   The interpreter puts the value kk into register Vx.
*/

func (cpu *Cpu) ins6xkk(op uint16) (string, error) {
	_, kk, _, x, _ := getParameters(op)

	cpu.v[x] = kk

	return fmt.Sprintf("LD V%x, %02d\t\t; Set Vx = kk", x, kk), nil
}

/*
   7xkk - ADD Vx, byte
   Set Vx = Vx + kk.
   Adds the value kk to the value of register Vx, then stores the result in Vx.
*/

func (cpu *Cpu) ins7xkk(op uint16) (string, error) {
	_, kk, _, x, _ := getParameters(op)

	cpu.v[x] += kk

	return fmt.Sprintf("ADD V%x, %02d\t; Set Vx = Vx + kk", x, kk), nil
}

/*
   8xy0 - LD Vx, Vy
   Set Vx = Vy.
   Stores the value of register Vy in register Vx.
*/

func (cpu *Cpu) ins8xy0(op uint16) (string, error) {
	_, _, _, x, y := getParameters(op)

	cpu.v[x] = cpu.v[y]

	return fmt.Sprintf("LD V%x, V%x\t\t; Set Vx = Vy", x, y), nil
}

/*
   8xy1 - OR Vx, Vy
   Set Vx = Vx OR Vy.
   Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx. A bitwise OR compares the corrseponding bits from two values, and if either bit is 1, then the same bit in the result is also 1. Otherwise, it is 0.
*/

func (cpu *Cpu) ins8xy1(op uint16) (string, error) {
	_, _, _, x, y := getParameters(op)

	cpu.v[x] |= cpu.v[y]
	if x == 0x0f {
		cpu.v[15] >>= 7
	} else {
		cpu.v[15] = 0
	}

	return fmt.Sprintf("OR V%x, V%x\t; Set Vx = Vx OR Vy", x, y), nil
}

/*
   8xy2 - AND Vx, Vy
   Set Vx = Vx AND Vy.
   Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx. A bitwise AND compares the corrseponding bits from two values, and if both bits are 1, then the same bit in the result is also 1. Otherwise, it is 0.
*/

func (cpu *Cpu) ins8xy2(op uint16) (string, error) {
	_, _, _, x, y := getParameters(op)

	cpu.v[x] &= cpu.v[y]
	if x == 0x0f {
		cpu.v[15] >>= 7
	} else {
		cpu.v[15] = 0
	}

	return fmt.Sprintf("AND V%x, V%x\t; Set Vx = Vx AND Vy", x, y), nil
}

/*
   8xy3 - XOR Vx, Vy
   Set Vx = Vx XOR Vy.
   Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx. An exclusive OR compares the corrseponding bits from two values, and if the bits are not both the same, then the corresponding bit in the result is set to 1. Otherwise, it is 0.
*/

func (cpu *Cpu) ins8xy3(op uint16) (string, error) {
	_, _, _, x, y := getParameters(op)

	cpu.v[x] ^= cpu.v[y]
	if x == 0x0f {
		cpu.v[15] >>= 7
	} else {
		cpu.v[15] = 0
	}

	return fmt.Sprintf("XOR V%x, V%x\t; Set Vx = Vx XOR Vy", x, y), nil
}

/*
   8xy4 - ADD Vx, Vy
   Set Vx = Vx + Vy, set VF = carry.
   The values of Vx and Vy are added together. If the result is greater than 8 bits (i.e., > 255,) VF is set to 1, otherwise 0. Only the lowest 8 bits of the result are kept, and stored in Vx.
*/

func (cpu *Cpu) ins8xy4(op uint16) (string, error) {
	_, _, _, x, y := getParameters(op)

	res := uint16(cpu.v[x])
	res += uint16(cpu.v[y])

	cpu.v[x] = byte(res & 0xff)
	cpu.v[15] = 0

	if res > 255 {
		cpu.v[15] = 1
	}

	return fmt.Sprintf("ADD V%x, V%x\t; Set Vx = Vx + Vy, set VF = carry", x, y), nil
}

/*
   8xy5 - SUB Vx, Vy
   Set Vx = Vx - Vy, set VF = NOT borrow.
   If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.
*/

func (cpu *Cpu) ins8xy5(op uint16) (string, error) {
	_, _, _, x, y := getParameters(op)

	res := int16(cpu.v[x])
	res -= int16(cpu.v[y])
	cpu.v[x] = byte(res & 0xff)

	cpu.v[15] = 0
	if res >= 0 {
		cpu.v[15] = 1
	}

	return fmt.Sprintf("SUB V%x, V%x\t; Set Vx = Vx - Vy, set VF = NOT borrow", x, y), nil
}

/*
   8xy6 - SHR Vx {, Vy}
   Set Vx = Vx SHR 1.
   If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
*/

func (cpu *Cpu) ins8xy6(op uint16) (string, error) {
	_, _, _, x, y := getParameters(op)

	cpu.v[x] = cpu.v[y] // only original COSMAC VIP

	carry := byte(0)
	if (cpu.v[x] & 1) == 1 {
		carry = 1
	}

	cpu.v[x] >>= 1
	cpu.v[15] = carry

	return fmt.Sprintf("SHR V%x, V%x\t; Vx = Vx SHR 1", x, y), nil
}

/*
   8xy7 - SUBN Vx, Vy
   Set Vx = Vy - Vx, set VF = NOT borrow.
   If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
*/

func (cpu *Cpu) ins8xy7(op uint16) (string, error) {
	_, _, _, x, y := getParameters(op)

	res := int16(cpu.v[y])
	res -= int16(cpu.v[x])
	cpu.v[x] = byte(res & 0xff)

	cpu.v[15] = 0
	if res >= 0 {
		cpu.v[15] = 1
	}

	return fmt.Sprintf("SUBN V%x, V%x\t; Set Vx = Vy - Vx, set VF = NOT borrow", x, y), nil
}

/*
   8xyE - SHL Vx {, Vy}
   Set Vx = Vx SHL 1.
   If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.
*/

func (cpu *Cpu) ins8xyE(op uint16) (string, error) {
	_, _, _, x, y := getParameters(op)

	cpu.v[x] = cpu.v[y] // only original COSMAC VIP

	carry := byte(0)
	if (cpu.v[x] & 0x80) == 0x80 {
		carry = 1
	}

	cpu.v[x] <<= 1
	cpu.v[15] = carry

	return fmt.Sprintf("SHL V%x, V%x\t; Set Vx = Vx SHL 1", x, y), nil
}

/*
   9xy0 - SNE Vx, Vy
   Skip next instruction if Vx != Vy.
   The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.
*/

func (cpu *Cpu) ins9xy0(op uint16) (string, error) {
	_, _, _, x, y := getParameters(op)

	if cpu.v[x] != cpu.v[y] {
		cpu.cnt += 2
	}

	return fmt.Sprintf("SNE V%x, V%x\t; Skip next instruction if Vx != Vy", x, y), nil
}

/*
   Annn - LD I, addr
   Set I = nnn.
   The value of register I is set to nnn.
*/

func (cpu *Cpu) insAnnn(op uint16) (string, error) {
	nnn, _, _, _, _ := getParameters(op)

	cpu.i = nnn

	return fmt.Sprintf("LD I, 0x%03x\t; Set I = nnn", nnn), nil
}

/*
   Bnnn - JP V0, addr
   Jump to location nnn + V0.
   The program counter is set to nnn plus the value of V0.
*/

func (cpu *Cpu) insBnnn(op uint16) (string, error) {
	nnn, _, _, _, _ := getParameters(op)

	res := nnn + uint16(cpu.v[0])
	err := cpu.checkAddr(res)
	if err != nil {
		return "", err
	}
	cpu.cnt = res

	return fmt.Sprintf("JP V0, 0x%03x\t; Jump to location nnn + V0", nnn), nil
}

/*
   Cxkk - RND Vx, byte
   Set Vx = random byte AND kk.
   The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk. The results are stored in Vx. See instruction 8xy2 for more information on AND.
*/

func (cpu *Cpu) insCxkk(op uint16) (string, error) {
	_, kk, _, x, _ := getParameters(op)

	cpu.v[x] = byte(rand.Intn(255)) & kk

	return fmt.Sprintf("RND V%x, %02d\t; Set Vx = random byte AND kk", x, kk), nil
}

/*
   Dxyn - DRW Vx, Vy, nibble
   Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
   The interpreter reads n bytes from memory, starting at the address stored in I. These bytes are then displayed as sprites on screen at coordinates (Vx, Vy). Sprites are XORed onto the existing screen. If this causes any pixels to be erased, VF is set to 1, otherwise it is set to 0. If the sprite is positioned so part of it is outside the coordinates of the display, it wraps around to the opposite side of the screen. See instruction 8xy3 for more information on XOR, and section 2.4, Display, for more information on the Chip-8 screen and sprites.
*/

func (cpu *Cpu) insDxyn(op uint16) (string, error) {
	_, _, n, vx, vy := getParameters(op)

	// TODO so ugly

	x := cpu.v[vx] % hardware.DISPLAY_WIDTH
	y := cpu.v[vy] % hardware.DISPLAY_HEIGHT

	var dy byte = 0

	cpu.v[15] = 0
	for i := byte(0); i < n; i++ {
		line := cpu.memory[cpu.i+uint16(i)]
		for dx := byte(0); dx < 8; dx++ {
			if (line & (1 << (7 - dx))) != 0 {
				if cpu.display.PutPixel(x+dx, y+dy) {
					cpu.v[15] = 1
				}
			}
		}
		dy++
	}

	return fmt.Sprintf("DRW V%x, V%x, %02d\t; Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision", vx, vy, n), nil
}

/*
   Ex9E - SKP Vx
   Skip next instruction if key with the value of Vx is pressed.
   Checks the keyboard, and if the key corresponding to the value of Vx is currently in the down position, PC is increased by 2.
*/

func (cpu *Cpu) insEx9E(op uint16) (string, error) {
	_, _, _, x, _ := getParameters(op)

	if (cpu.kbrd & uint16(1<<cpu.v[x])) != 0 {
		cpu.cnt += 2
	}

	return fmt.Sprintf("SKP V%x\t; Skip next instruction if key with the value of Vx is pressed", x), nil
}

/*
   ExA1 - SKPN Vx
   Skip next instruction if key with the value of Vx is not pressed.
   Checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position, PC is increased by 2.
*/

func (cpu *Cpu) insExA1(op uint16) (string, error) {
	_, _, _, x, _ := getParameters(op)

	if (cpu.kbrd & uint16(1<<cpu.v[x])) == 0 {
		cpu.cnt += 2
	}

	return fmt.Sprintf("SKPN V%x\t\t; Skip next instruction if key with the value of Vx is not pressed", x), nil
}

/*
   Fx07 - LD Vx, DT
   Set Vx = delay timer value.
   The value of DT is placed into Vx.
*/

func (cpu *Cpu) insFx07(op uint16) (string, error) {
	_, _, _, x, _ := getParameters(op)

	cpu.v[x] = cpu.timerDelay

	return fmt.Sprintf("LD V%x, DT\t\t; Set Vx = delay timer value", x), nil
}

/*
   Fx0A - LD Vx, K
   Wait for a key press, store the value of the key in Vx.
   All execution stops until a key is pressed, then the value of that key is stored in Vx.
*/

func (cpu *Cpu) insFx0A(op uint16) (string, error) {
	_, _, _, x, _ := getParameters(op)

	res := cpu.keyboard.WaitKey()

	if res != 0x80 {
		cpu.v[x] = res
	} else {
		cpu.cnt -= 2
	}

	return fmt.Sprintf("LD V%x, K\t\t; Wait for a key press, store the value of the key in Vx", x), nil
}

/*
   Fx15 - LD DT, Vx
   Set delay timer = Vx.
   DT is set equal to the value of Vx.
*/

func (cpu *Cpu) insFx15(op uint16) (string, error) {
	_, _, _, x, _ := getParameters(op)

	cpu.timerDelay = cpu.v[x]

	return fmt.Sprintf("LD DT, V%x\t\t; Set delay timer = Vx", x), nil
}

/*
   Fx18 - LD ST, Vx
   Set sound timer = Vx.
   ST is set equal to the value of Vx.
*/

func (cpu *Cpu) insFx18(op uint16) (string, error) {
	_, _, _, x, _ := getParameters(op)

	cpu.timerSound = cpu.v[x]

	return fmt.Sprintf("LD ST, V%x\t\t; Set sound timer = Vx", x), nil
}

/*
   Fx1E - ADD I, Vx
   Set I = I + Vx.
   The values of I and Vx are added, and the results are stored in I.
*/

func (cpu *Cpu) insFx1E(op uint16) (string, error) {
	_, _, _, x, _ := getParameters(op)

	cpu.i += uint16(cpu.v[x])

	return fmt.Sprintf("ADD I, V%x\t\t; Set I = I + Vx", x), nil
}

/*
   Fx29 - LD F, Vx
   Set I = location of sprite for digit Vx.
   The value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx. See section 2.4, Display, for more information on the Chip-8 hexadecimal font.
*/

func (cpu *Cpu) insFx29(op uint16) (string, error) {
	_, _, _, x, _ := getParameters(op)

	cpu.i = SPRITE_ADDR + uint16(((cpu.v[x] & 0x0f) * 5))

	return fmt.Sprintf("LD F, V%x\t\t; Set I = location of sprite for digit Vx", x), nil
}

/*
   Fx33 - LD B, Vx
   Store BCD representation of Vx in memory locations I, I+1, and I+2.
   The interpreter takes the decimal value of Vx, and places the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.
*/

func (cpu *Cpu) insFx33(op uint16) (string, error) {
	_, _, _, x, _ := getParameters(op)

	b := []byte(fmt.Sprintf("%03d", cpu.v[x]))
	cpu.memory[cpu.i] = b[0] - 48
	cpu.memory[cpu.i+1] = b[1] - 48
	cpu.memory[cpu.i+2] = b[2] - 48

	return fmt.Sprintf("LD B, V%x\t\t; Store BCD representation of Vx in memory locations I, I+1, and I+2", x), nil
}

/*
   Fx55 - LD [I], Vx
   Store registers V0 through Vx in memory starting at location I.
   The interpreter copies the values of registers V0 through Vx into memory, starting at the address in I.
*/

func (cpu *Cpu) insFx55(op uint16) (string, error) {
	_, _, _, x, _ := getParameters(op)

	for i := uint16(0); i <= uint16(x); i++ {
		cpu.memory[cpu.i] = cpu.v[i]
		cpu.i += 1
	}

	return fmt.Sprintf("LD [I], V%x\t; Store registers V0 through Vx in memory starting at location I", x), nil
}

/*
   Fx65 - LD Vx, [I]
   Read registers V0 through Vx from memory starting at location I.
   The interpreter reads values from memory starting at location I into registers V0 through Vx.
*/

func (cpu *Cpu) insFx65(op uint16) (string, error) {
	_, _, _, x, _ := getParameters(op)

	for i := uint16(0); i <= uint16(x); i++ {
		cpu.v[i] = cpu.memory[cpu.i]
		cpu.i += 1
	}

	return fmt.Sprintf("LD V%x, [I]\t; Read registers V0 through Vx from memory starting at location I", x), nil
}
