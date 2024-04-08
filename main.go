package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ministergoose/chip8-emu-go/chip8"
	"github.com/ministergoose/chip8-emu-go/chip8/hardware"
	"github.com/ministergoose/chip8-emu-go/chip8/hardware/raylib"
)

func parseArgs() (bool, string) {
	filePath := ""
	diss := false

	if (len(os.Args) == 3) && (os.Args[1] == "diss") {
		return true, os.Args[2]
	} else if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [diss] <file path>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	filePath = os.Args[1]

	return diss, filePath
}

func main() {
	diss, filePath := parseArgs()

	if diss {
		chip8.Disassembler(filePath)
		return
	}

	dspl := raylib.NewDisplayRaylib()
	dspl.Init("Chip8 Go", 10.0)
	defer dspl.Close()

	kbrd := raylib.NewKeyboardRaylib()
	snd := hardware.NewSoundStd()

	Cpu := chip8.NewCPU(dspl, kbrd, snd)
	err := Cpu.Load(filePath)
	if err != nil {
		log.Fatal(err)
	}
	Cpu.Run()
}
