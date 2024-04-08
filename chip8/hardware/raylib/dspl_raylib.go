package raylib

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ministergoose/chip8-emu-go/chip8/hardware"
)

type DisplayRaylib struct {
	buffer [hardware.DISPLAY_HEIGHT][hardware.DISPLAY_WIDTH]byte
	camera rl.Camera2D
	title  string

	bgrColor color.RGBA
	frgColor color.RGBA
}

func NewDisplayRaylib() *DisplayRaylib {
	return &DisplayRaylib{}
}

func (dspl *DisplayRaylib) Init(title string, scale float32) {
	dspl.title = title
	rl.SetTraceLog(rl.LogError)
	rl.InitWindow(int32(hardware.DISPLAY_WIDTH*scale), int32(hardware.DISPLAY_HEIGHT*scale), title)
	rl.SetTargetFPS(60)
	dspl.camera = rl.NewCamera2D(rl.NewVector2(0.0, 0.0), rl.NewVector2(0.0, 0.0), 0.0, scale)

	dspl.bgrColor = rl.Black
	dspl.frgColor = rl.Green
}

func (dspl *DisplayRaylib) PutPixel(x, y byte) bool {
	if x < hardware.DISPLAY_WIDTH && y < hardware.DISPLAY_HEIGHT {
		dspl.buffer[y][x] ^= 1
		return dspl.buffer[y][x] == 0
	}

	return false
}

func (dspl *DisplayRaylib) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(dspl.bgrColor)
	rl.BeginMode2D(dspl.camera)

	for y := 0; y < hardware.DISPLAY_HEIGHT; y++ {
		for x := 0; x < hardware.DISPLAY_WIDTH; x++ {
			if dspl.buffer[y][x] > 0 {
				rl.DrawPixel(int32(x), int32(y), dspl.frgColor)
			}
		}
	}

	rl.EndMode2D()
	rl.EndDrawing()
	rl.SetWindowTitle(fmt.Sprintf("%s [FPS: %.2f]", dspl.title, rl.GetFPS()))
}

func (dspl *DisplayRaylib) Cls() {
	for y := 0; y < hardware.DISPLAY_HEIGHT; y++ {
		for x := 0; x < hardware.DISPLAY_WIDTH; x++ {
			dspl.buffer[y][x] = 0
		}
	}
}

func (dspl *DisplayRaylib) Dump() {
	fmt.Print("  |")
	for x := 0; x < hardware.DISPLAY_WIDTH; x++ {
		fmt.Printf("%02d|", x)
	}

	for y := 0; y < hardware.DISPLAY_HEIGHT; y++ {
		fmt.Printf("\n%02d|", y)
		for x := 0; x < hardware.DISPLAY_WIDTH; x++ {
			fmt.Printf("%2d|", dspl.buffer[y][x])
		}
	}
	fmt.Println()
}

func (dspl *DisplayRaylib) ShouldClose() bool {
	return rl.WindowShouldClose()
}

func (dspl *DisplayRaylib) Close() {
	rl.CloseWindow()
}
