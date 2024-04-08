package chip8

import (
	"fmt"
	"math/rand"

	"github.com/ministergoose/chip8-emu-go/chip8/hardware"
)

func getParametersEX(i uint16) (nnn uint16, kk byte, n byte, x byte, y byte, op byte) {
	nnn = i & 0x0fff
	kk = (byte)(i & 0x00ff)
	n = (byte)(i & 0x000f)
	x = (byte)((i & 0x0f00) >> 8)
	y = (byte)((i & 0x00f0) >> 4)
	op = (byte)((i & 0xf000) >> 12)
	return nnn, kk, n, x, y, op
}

func (cpu *Cpu) RunInstFast(inst uint16) (string, error) {
	nnn, kk, n, x, y, op := getParametersEX(inst)

	switch op {
	case 0x0:
		{
			if nnn == 0x0e0 {
				cpu.display.Cls()
			} else if nnn == 0x0ee {
				addr, _ := cpu.stack.Pop()
				cpu.cnt = addr
			} else {
				// TODO 0nnn - SYS addr
				fmt.Println("0nnn - SYS addr")
			}
		}
	case 0x1:
		{
			cpu.cnt = nnn
		}
	case 0x2:
		{
			cpu.stack.Push(cpu.cnt)
			cpu.cnt = nnn
		}
	case 0x3:
		{
			if cpu.v[x] == kk {
				cpu.cnt += 2
			}
		}
	case 0x4:
		{
			if cpu.v[x] != kk {
				cpu.cnt += 2
			}
		}
	case 0x5:
		{
			if n == 0 {
				if cpu.v[x] == cpu.v[y] {
					cpu.cnt += 2
				}
			}
		}
	case 0x6:
		{
			cpu.v[x] = kk
		}
	case 0x7:
		{
			cpu.v[x] += kk
		}
	case 0x8:
		{
			if n == 0 {
				cpu.v[x] = cpu.v[y]
			} else if n == 1 {
				cpu.v[x] |= cpu.v[y]
				if x == 0x0f {
					cpu.v[15] >>= 7
				} else {
					cpu.v[15] = 0
				}
			} else if n == 2 {
				cpu.v[x] &= cpu.v[y]
				if x == 0x0f {
					cpu.v[15] >>= 7
				} else {
					cpu.v[15] = 0
				}
			} else if n == 3 {
				cpu.v[x] ^= cpu.v[y]
				if x == 0x0f {
					cpu.v[15] >>= 7
				} else {
					cpu.v[15] = 0
				}
			} else if n == 4 {
				res := uint16(cpu.v[x])
				res += uint16(cpu.v[y])

				cpu.v[x] = byte(res & 0xff)
				cpu.v[15] = 0
				if res > 255 {
					cpu.v[15] = 1
				}
			} else if n == 5 {
				res := int16(cpu.v[x])
				res -= int16(cpu.v[y])

				cpu.v[x] = byte(res & 0xff)
				cpu.v[15] = 0
				if res >= 0 {
					cpu.v[15] = 1
				}
			} else if n == 6 {
				cpu.v[x] = cpu.v[y] // only original COSMAC VIP

				carry := byte(0)
				if (cpu.v[x] & 1) == 1 {
					carry = 1
				}

				cpu.v[x] >>= 1
				cpu.v[15] = carry
			} else if n == 7 {
				res := int16(cpu.v[y])
				res -= int16(cpu.v[x])

				cpu.v[x] = byte(res & 0xff)
				cpu.v[15] = 0
				if res >= 0 {
					cpu.v[15] = 1
				}
			} else if n == 0xe {
				cpu.v[x] = cpu.v[y] // only original COSMAC VIP

				carry := byte(0)
				if (cpu.v[x] & 0x80) == 0x80 {
					carry = 1
				}

				cpu.v[x] <<= 1
				cpu.v[15] = carry
			}
		}
	case 0x9:
		{
			if n == 0 {
				if cpu.v[x] != cpu.v[y] {
					cpu.cnt += 2
				}
			}
		}
	case 0xa:
		{
			cpu.i = nnn
		}
	case 0xb:
		{
			cpu.cnt = nnn + uint16(cpu.v[0])
		}
	case 0xc:
		{
			cpu.v[x] = byte(rand.Intn(255)) & kk
		}
	case 0xd:
		{
			x = cpu.v[x] % hardware.DISPLAY_WIDTH
			y = cpu.v[y] % hardware.DISPLAY_HEIGHT

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
		}
	case 0xe:
		{
			if kk == 0x9e {
				if (cpu.kbrd & uint16(1<<cpu.v[x])) != 0 {
					cpu.cnt += 2
				}
			} else if kk == 0xa1 {
				if (cpu.kbrd & uint16(1<<cpu.v[x])) == 0 {
					cpu.cnt += 2
				}
			}
		}
	case 0xf:
		{
			if kk == 0x07 {
				cpu.v[x] = cpu.timerDelay
			} else if kk == 0x0a {
				res := cpu.keyboard.WaitKey()

				if res != 0x80 {
					cpu.v[x] = res
				} else {
					cpu.cnt -= 2
				}
			} else if kk == 0x15 {
				cpu.timerDelay = cpu.v[x]
			} else if kk == 0x18 {
				cpu.timerSound = cpu.v[x]
			} else if kk == 0x1e {
				cpu.i += uint16(cpu.v[x])
			} else if kk == 0x29 {
				cpu.i = SPRITE_ADDR + uint16(((cpu.v[x] & 0x0f) * 5))
			} else if kk == 0x33 {
				b := []byte(fmt.Sprintf("%03d", cpu.v[x]))
				cpu.memory[cpu.i] = b[0] - 48
				cpu.memory[cpu.i+1] = b[1] - 48
				cpu.memory[cpu.i+2] = b[2] - 48
			} else if kk == 0x55 {
				for i := uint16(0); i <= uint16(x); i++ {
					cpu.memory[cpu.i] = cpu.v[i]
					cpu.i += 1
				}
			} else if kk == 0x65 {
				for i := uint16(0); i <= uint16(x); i++ {
					cpu.v[i] = cpu.memory[cpu.i]
					cpu.i += 1
				}
			}
		}
	default:
		{
			return "", fmt.Errorf("unknown opcode: 0x%04X", inst)
		}
	}

	return "", nil
}
