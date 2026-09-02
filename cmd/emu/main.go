package main

import (
	"flag"
	"syscall"
	"os"
	"os/signal"
	"log"
	"gbemu/internal/gb"
)

var (
	debugMode = flag.Bool("debug", false, "debug output to track PC, SP, and register values for each opcode") 
)

func main() {
	flag.Parse()
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	go start()
	sig := <- cancelChan
	log.Printf("caught signal %v\n", sig)
}

func start() {
	romPath := flag.Arg(0)
	if romPath == "" {
		log.Fatal("no ROM file provided")
	}

	gb, err := gb.InitGameboy(romPath) 
	if err != nil {
		log.Fatal(err)
	}

	if *debugMode {
		gb.LogOpcodes = true
	}

	for {
		if err := gb.Step(); err != nil {
			log.Fatalf("error during gameboy step: %v\n", err.Error())
		}
	}
}
