package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	LCDWidth int32 = 160
	LCDHeight int32 = 144
	initialScale int32 = 3
)

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

func (lcd *LCD) Close() {
	rl.CloseWindow()
}
