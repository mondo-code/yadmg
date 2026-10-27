package main

import (
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	LCDWidth     int32 = 160
	LCDHeight    int32 = 144
	initialScale int32 = 3

	// audio constants
	sampleRate  uint32 = 44100
	sampleSize  uint32 = 32
	numChannels uint32 = 2
	chunkFrames        = 1024
)

var gameboyScreenTint = rl.Color{R: 0xe0, G: 0xf8, B: 0xd0, A: 0xff}

type Frontend struct {
	texture    rl.Texture2D
	keymap     [numButtons]int32
	stream     rl.AudioStream
	audioQueue []float32
	chunk      []float32
}

func InitFrontend() *Frontend {
	fe := &Frontend{}
	fe.keymap = loadKeymap()
	rl.InitAudioDevice()
	rl.SetAudioStreamBufferSizeDefault(chunkFrames)
	fe.stream = rl.LoadAudioStream(sampleRate, sampleSize, numChannels)
	fe.chunk = make([]float32, chunkFrames*2)
	fe.audioQueue = make([]float32, 0)
	return fe
}

func (fe *Frontend) Start(name string) {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(LCDWidth*initialScale, LCDHeight*initialScale, name)

	img := rl.GenImageColor(int(LCDWidth), int(LCDHeight), rl.Black)
	fe.texture = rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	rl.PlayAudioStream(fe.stream)
}

func (fe *Frontend) IsRunning() bool {
	return !rl.WindowShouldClose()
}

func (fe *Frontend) Render(pixels *[160][144][3]uint8) {
	data := make([]byte, LCDWidth*LCDHeight*4)
	i := 0
	for y := range LCDHeight {
		for x := range LCDWidth {
			col := pixels[x][y]
			data[i], data[i+1], data[i+2], data[i+3] = col[0], col[1], col[2], 0xff
			i += 4
		}
	}
	rl.UpdateTexture(fe.texture, data)
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	winWidth, winHeight := float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
	scale := min(winWidth/float32(LCDWidth), winHeight/float32(LCDHeight))
	destWidth, destHeight := float32(LCDWidth)*scale, float32(LCDHeight)*scale
	offsetX, offsetY := (winWidth-destWidth)/2, (winHeight-destHeight)/2
	src := rl.Rectangle{Width: float32(LCDWidth), Height: float32(LCDHeight)}
	dest := rl.Rectangle{X: offsetX, Y: offsetY, Width: destWidth, Height: destHeight}
	rl.DrawTexturePro(fe.texture, src, dest, rl.Vector2{}, 0, gameboyScreenTint)

	rl.EndDrawing()
}

func (fe *Frontend) DoInput() byte {
	var input byte

	for i, n := range fe.keymap {
		if rl.IsKeyDown(n) {
			input |= (1 << i)
		}
	}

	return input
}

// this just queues audio from the sample buffer that Gameboy gives us
func (fe *Frontend) EmitAudio(sampleBuffer []float32) {
	fe.audioQueue = append(fe.audioQueue, sampleBuffer...)
	if limit := 4 * len(fe.chunk); len(fe.audioQueue) > limit {
		// if emulator runs ahead of audio, drop oldest samples to bound latency
		overrun := len(fe.audioQueue) - limit
		fe.audioQueue = fe.audioQueue[:copy(fe.audioQueue, fe.audioQueue[overrun:])]
	}
}

func (fe *Frontend) Update() {
	for rl.IsAudioStreamProcessed(fe.stream) {
		n := copy(fe.chunk, fe.audioQueue)
		fe.audioQueue = slices.Delete(fe.audioQueue, 0, n)
		// pad rest of the chunk w/ silence
		clear(fe.chunk[n:])
		rl.UpdateAudioStream(fe.stream, fe.chunk[:chunkFrames])
	}
}

func (fe *Frontend) Close() {
	rl.UnloadAudioStream(fe.stream)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}
