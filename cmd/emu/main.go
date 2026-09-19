package main

import (
	"flag"
	"gbemu/internal/gb"
	"log"
	"time"
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

	lcd.Start("yadmg")
	defer lcd.Close()

	ticker := time.NewTicker(4 * time.Millisecond)
	last := time.Now()
	for now := range ticker.C {
		if !lcd.IsRunning() {
			return
		}

		cycleBudget := uint64(now.Sub(last)) * gb.CPUSpeedHz / uint64(time.Second)
		last = now
		if err := gameboy.Step(cycleBudget); err != nil {
			log.Fatalf("error during gameboy step: %v\n", err.Error())
		}
	}
}
