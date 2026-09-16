package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	LCDWidth int32 = 160
	LCDHeight int32 = 144
	initialScale int32 = 3
)

// default binds for dpad will be vim keys because i'm a maniac
// indices correspond to the bit in the input byte returned by Input()
// 0 = A, 1 = B, 2 = select, 3 = start, 4 = right, 5 = left, 6 = up, 7 = down
var keybinds = [8]int32{rl.KeyZ, rl.KeyX, rl.KeyRightShift, rl.KeyEnter, rl.KeyL, rl.KeyH, rl.KeyK, rl.KeyJ}

type LCD struct {
	texture 		rl.Texture2D
}

func (lcd *LCD) Start(name string) {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(LCDWidth * initialScale, LCDHeight * initialScale, name)

	img := rl.GenImageColor(int(LCDWidth), int(LCDHeight), rl.Black)
	lcd.texture = rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
}

func (lcd *LCD) IsRunning() bool {
	return !rl.WindowShouldClose()
}

func (lcd *LCD) Render(pixels *[160][144][3]uint8) {
	data := make([]byte, LCDWidth * LCDHeight * 4)
	i := 0
	for y := range LCDHeight {
		for x := range LCDWidth {
			col := pixels[x][y]
			data[i], data[i+1], data[i+2], data[i+3] = col[0], col[1], col[2], 0xff
			i += 4
		}
	}
	rl.UpdateTexture(lcd.texture, data)
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	winWidth, winHeight := float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
	scale := min(winWidth / float32(LCDWidth), winHeight / float32(LCDHeight))
	destWidth, destHeight := float32(LCDWidth)*scale, float32(LCDHeight)*scale
	offsetX, offsetY := (winWidth-destWidth) / 2, (winHeight-destHeight) / 2
	src := rl.Rectangle{Width: float32(LCDWidth), Height: float32(LCDHeight)}
	dest := rl.Rectangle{X: offsetX, Y: offsetY, Width: destWidth, Height: destHeight}
	rl.DrawTexturePro(lcd.texture, src, dest, rl.Vector2{}, 0, rl.Color{R: 0xe0, G: 0xf8, B: 0xd0, A: 0xff})

	rl.EndDrawing()
}

func (lcd *LCD) DoInput() byte {
	var input byte

	for i, n := range keybinds {
		if rl.IsKeyDown(n) {
			input |= (1 << i)
		}
	}

	return input
}

func (lcd *LCD) Close() {
	rl.CloseWindow()
}
