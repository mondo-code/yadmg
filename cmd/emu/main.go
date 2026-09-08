package main 

import (
	"flag"
	"gbemu/internal/gb"
	"log"
)

var (
	debugMode = flag.Bool("debug", false, "debug output to track PC, SP, and register values for each opcode") 
)

func main() {
	flag.Parse()
	romPath := flag.Arg(0)
	if romPath == "" {
		log.Fatal("no ROM file provided")
	}

	lcd := &LCD{}
	gameboy, err := gb.InitGameboy(romPath, lcd)
	if err != nil {
		log.Fatal(err)
	}

	if *debugMode {
		gameboy.LogOpcodes = true
	}

	// raylib window init
	lcd.Start("yadmg")
	defer lcd.Close()

	for lcd.IsRunning() {
		if err := gameboy.Step(); err != nil {
			log.Fatalf("error during gameboy step: %v\n", err.Error())
		}
	}
}
