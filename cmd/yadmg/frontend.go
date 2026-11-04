package main

import (
	"github.com/mondo-code/yadmg/internal/gb"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	initialScale int32  = 3
	sampleSize   uint32 = 32

	// audio constants
	chunkFrames = 1024
)

var gameboyScreenTint = rl.Color{R: 0xe0, G: 0xf8, B: 0xd0, A: 0xff}

type frontend struct {
	texture    rl.Texture2D
	keymap     [numButtons]int32
	stream     rl.AudioStream
	audioQueue []float32
	chunk      []float32
	muted      bool
}

func initFrontend(mute bool) *frontend {
	fe := &frontend{}
	fe.muted = mute
	fe.keymap = loadKeymap()
	if !mute {
		rl.InitAudioDevice()
		rl.SetAudioStreamBufferSizeDefault(chunkFrames)
		fe.stream = rl.LoadAudioStream(gb.AudioSampleRate, sampleSize, gb.AudioOutputChannels)
		fe.chunk = make([]float32, chunkFrames*2)
		fe.audioQueue = make([]float32, 0)
	}
	return fe
}

func (fe *frontend) start(name string) {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(gb.ScreenWidth*initialScale, gb.ScreenHeight*initialScale, name)

	img := rl.GenImageColor(int(gb.ScreenWidth), int(gb.ScreenHeight), rl.Black)
	fe.texture = rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	rl.PlayAudioStream(fe.stream)
}

func (fe *frontend) Render(pixels *[160][144][3]uint8) {
	data := make([]byte, gb.ScreenWidth*gb.ScreenHeight*4)
	i := 0
	for y := range gb.ScreenHeight {
		for x := range gb.ScreenWidth {
			col := pixels[x][y]
			data[i], data[i+1], data[i+2], data[i+3] = col[0], col[1], col[2], 0xff
			i += 4
		}
	}
	rl.UpdateTexture(fe.texture, data)
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	winWidth, winHeight := float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
	scale := min(winWidth/float32(gb.ScreenWidth), winHeight/float32(gb.ScreenHeight))
	destWidth, destHeight := float32(gb.ScreenWidth)*scale, float32(gb.ScreenHeight)*scale
	offsetX, offsetY := (winWidth-destWidth)/2, (winHeight-destHeight)/2
	src := rl.Rectangle{Width: float32(gb.ScreenWidth), Height: float32(gb.ScreenHeight)}
	dest := rl.Rectangle{X: offsetX, Y: offsetY, Width: destWidth, Height: destHeight}
	rl.DrawTexturePro(fe.texture, src, dest, rl.Vector2{}, 0, gameboyScreenTint)

	rl.EndDrawing()
}

func (fe *frontend) IsRunning() bool {
	return !rl.WindowShouldClose()
}

func (fe *frontend) DoInput() byte {
	var input byte

	for i, n := range fe.keymap {
		if rl.IsKeyDown(n) {
			input |= (1 << i)
		}
	}

	return input
}

// this just queues audio from the sample buffer that Gameboy gives us
func (fe *frontend) EmitAudio(sampleBuffer []float32) {
	if fe.muted {
		return
	}
	fe.audioQueue = append(fe.audioQueue, sampleBuffer...)
	if limit := 4 * len(fe.chunk); len(fe.audioQueue) > limit {
		// if emulator runs ahead of audio, drop oldest samples to bound latency
		overrun := len(fe.audioQueue) - limit
		fe.audioQueue = fe.audioQueue[:copy(fe.audioQueue, fe.audioQueue[overrun:])]
	}
}

func (fe *frontend) processAudio() {
	if fe.muted {
		return
	}
	for rl.IsAudioStreamProcessed(fe.stream) {
		n := copy(fe.chunk, fe.audioQueue)
		fe.audioQueue = slices.Delete(fe.audioQueue, 0, n)
		// pad rest of the chunk w/ silence
		clear(fe.chunk[n:])
		rl.UpdateAudioStream(fe.stream, fe.chunk[:chunkFrames])
	}
}

func (fe *frontend) close() {
	rl.UnloadAudioStream(fe.stream)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}
