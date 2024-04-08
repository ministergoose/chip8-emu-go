package chip8

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ministergoose/chip8-emu-go/chip8/hardware"
	"github.com/ministergoose/chip8-emu-go/chip8/hardware/empty"
)

const (
	MEMORY_SIZE        = 4096
	SPRITE_ADDR uint16 = 0x00
	START_ADDR  uint16 = 0x200
	IPS         int    = 700 // instr per second
)

var SPRITES []byte = []byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type InsFunc func(op uint16) (string, error)

type Opcode struct {
	Result uint16
	Mask   uint16
	f      InsFunc
}

type Cpu struct {
	v       [16]byte
	i       uint16
	cnt     uint16
	memory  [MEMORY_SIZE]byte
	stack   Stack
	display hardware.Display
	sound   hardware.Sound

	keyboard hardware.Keyboard
	kbrd     uint16

	timerDelay byte
	timerSound byte

	Opcodes []Opcode
	Debug   bool
}

func (c *Cpu) RunInst(inst uint16) (string, error) {
	for _, x := range c.Opcodes {
		if (inst & x.Mask) == x.Result {
			return x.f(inst)
		}
	}

	return fmt.Sprintf("Unk: 0x%04X", inst), fmt.Errorf("unknown opcode: 0x%04X", inst)
}

func (c *Cpu) Load(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	fsize := len(data)

	if fsize > MEMORY_SIZE-int(START_ADDR) {
		return fmt.Errorf("program is too big! Length: %d", fsize)
	}

	if c.Debug {
		fmt.Printf("Load %d bytes\n", fsize)
		fmt.Print(hex.Dump(data))
	}

	c.DMA(START_ADDR, data, uint16(len(data)))

	return nil
}

func (c *Cpu) Run() {
	//frameTime := time.Second / 60
	frameTime := time.Millisecond * 16

	for !c.display.ShouldClose() {
		c.kbrd = c.keyboard.ReadKeys()

		start := time.Now()
		for i := 0; i < (IPS / 60); i++ {
			inst := uint16(uint16(c.memory[c.cnt])<<8) + uint16(c.memory[c.cnt+1])
			c.cnt += 2
			//str, err := c.RunInst(inst)
			str, err := c.RunInstFast(inst)
			if err != nil {
				log.Println(err)
			}
			if c.Debug {
				log.Println(str)
			}
		}
		delayTime := frameTime - time.Since(start)

		if delayTime > 0 {
			time.Sleep(delayTime)
		}

		c.TimersTick()
		c.display.Draw()
	}

}

func (c *Cpu) Reset() {
	for i := START_ADDR; i < MEMORY_SIZE; i++ {
		c.memory[i] = 0
	}

	for i := 0; i <= 15; i++ {
		c.v[i] = 0
	}

	c.display.Cls()
	c.stack.Reset()

	c.i = 0
	c.cnt = START_ADDR
	c.timerDelay = 0
	c.timerSound = 0
}

func (c *Cpu) checkAddr(addr uint16) error {
	if addr >= MEMORY_SIZE {
		return fmt.Errorf("bad address: %03x", addr)
	}

	return nil
}

func (c *Cpu) DMA(destPos uint16, src []byte, length uint16) {
	for i := uint16(0); i < length; i++ {
		c.memory[i+destPos] = src[i]
	}
}

func (c *Cpu) PrintDebug() {
	fmt.Println("   | 0| 1| 2| 3| 4| 5| 6| 7| 8| 9| A| B| C| D| E| F|")
	fmt.Print(" v |")
	for _, v := range c.v {
		fmt.Printf("%02X|", v)
	}
	fmt.Println()
}

func (c *Cpu) TimersTick() {
	if c.timerSound > 0 {
		c.sound.Beep()
		c.timerSound -= 1
	}
	if c.timerDelay > 0 {
		c.timerDelay -= 1
	}
}

func NewCPU(dspl hardware.Display, kbrd hardware.Keyboard, snd hardware.Sound) *Cpu {
	rand.Seed(time.Now().UnixNano())

	c := Cpu{display: dspl, keyboard: kbrd, sound: snd, stack: NewStackStd()}
	c.InstructionsInit()
	c.DMA(SPRITE_ADDR, SPRITES, uint16(len(SPRITES)))
	c.Reset()

	return &c
}

func Disassembler(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(fmt.Sprintf("%s.dis.txt", filePath))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	c := Cpu{display: empty.NewDisplayEmpty(), keyboard: empty.NewKeyboardEmpty(), sound: hardware.NewSoundStd(), stack: empty.NewStackEmpty()}
	c.InstructionsInit()

	fsize := len(data)

	f.WriteString("Addr : [code] Asm\t\t; Desc\n")
	for i := 0; i < fsize; i += 2 {
		inst := uint16(uint16(data[i])<<8) + uint16(data[i+1])
		str, err := c.RunInst(inst)
		f.WriteString(fmt.Sprintf("%04Xh: [%04X] %s\n", (0x200 + i), inst, str))
		if err != nil {
			f.WriteString(fmt.Sprintf("%v\n", err))
		}
	}

}
