package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	"go.yaml.in/yaml/v4"
)

// indices correspond to the bit in the input byte returned by Input()
const (
	buttonA = iota
	buttonB
	buttonSelect
	buttonStart
	dpadRight
	dpadLeft
	dpadUp
	dpadDown
	numButtons
)

// indexed according to above button indices
var buttonNames = [numButtons]string{
	"buttonA",
	"buttonB",
	"buttonSelect",
	"buttonStart",
	"dpadRight",
	"dpadLeft",
	"dpadUp",
	"dpadDown",
}

var defaultKeybinds = [numButtons]int32{
	buttonA:      rl.KeyZ,
	buttonB:      rl.KeyX,
	buttonSelect: rl.KeyRightShift,
	buttonStart:  rl.KeyEnter,
	dpadRight:    rl.KeyRight,
	dpadLeft:     rl.KeyLeft,
	dpadUp:       rl.KeyUp,
	dpadDown:     rl.KeyDown,
}

var keyNames = func() map[string]int32 {
	names := map[string]int32{
		"Right":      rl.KeyRight,
		"Left":       rl.KeyLeft,
		"Up":         rl.KeyUp,
		"Down":       rl.KeyDown,
		"RightShift": rl.KeyRightShift,
		"LeftShift":  rl.KeyLeftShift,
		"Enter":      rl.KeyEnter,
	}

	// letter keys
	for c := 'A'; c <= 'Z'; c++ {
		names[string(c)] = rl.KeyA + (c - 'A')
	}

	return names
}()

var configPaths = func() []string {
	paths := []string{"keymap.yml"}
	if dir, err := os.UserConfigDir(); err == nil {
		paths = append(paths, filepath.Join(dir, "yadmg", "keymap.yml"))
	}
	return paths
}()

func readConfigFile() (res []byte, err error) {
	for _, path := range configPaths {
		data, err := os.ReadFile(path)
		if err == nil {
			return data, nil
		}
	}
	return nil, fmt.Errorf("searched %s: %w", strings.Join(configPaths, ", "), fs.ErrNotExist)
}

func parseKeymap(data []byte) ([numButtons]int32, error) {
	var entries map[string]string
	if err := yaml.Unmarshal(data, &entries); err != nil {
		return defaultKeybinds, fmt.Errorf("failed to unmarshal keymap data, reverting to default keymap: %v", err)
	}

	binds := defaultKeybinds
	for button, key := range entries {
		i := slices.Index(buttonNames[:], button)
		if i < 0 {
			log.Printf("found unknown button %v, using default bind\n", button)
			continue
		}
		keyNum, ok := keyNames[key]
		if !ok {
			log.Printf("found unknown key %v, using default bind\n", key)
			continue
		}
		binds[i] = keyNum
	}

	return binds, nil
}

func loadKeymap() [numButtons]int32 {
	data, err := readConfigFile()
	if err != nil {
		log.Printf("failed to load keymap.yml, using default: %v", err)
		return defaultKeybinds
	}

	binds, err := parseKeymap(data)
	if err != nil {
		log.Printf("failed to parse keymap data: %v", err)
		return defaultKeybinds
	}
	return binds
}
