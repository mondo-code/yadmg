package gb

const divFreq = 0xff

type Screen interface {
	Render(framebuffer *[160][144][3]uint8)
	IsRunning() bool
	DoInput() byte
}

type Speaker interface {
	EmitAudio(sampleBuffer []float32)
}

type Gameboy struct {
	cpu           *CPU
	ppu           *PPU
	memoryBus     *MemoryBus
	apu           *APU
	emuCycles     uint64
	timerCycles   uint16
	divCycles     uint16
	screen        Screen
	speaker       Speaker
	joypad        *Joypad
	framebuffer   *[ScreenWidth][ScreenHeight][3]uint8
	paused        bool
	LogOpcodes    bool
	screenCleared bool
}

func InitGameboy(romPath string, screen Screen, speaker Speaker) (*Gameboy, error) {
	gb := &Gameboy{}
	gb.screen = screen
	gb.speaker = speaker
	gb.cpu = initCPU(gb)
	mb := initMemoryBus(gb)
	gb.memoryBus = mb
	gb.ppu = initPPU(gb)
	gb.framebuffer = &[ScreenWidth][ScreenHeight][3]uint8{}
	gb.apu = initAPU()
	gb.joypad = initJoypad()
	gb.clearScreen()
	_, err := mb.loadCartridge(romPath)
	if err != nil {
		return nil, err
	}

	return gb, nil
}

func (gb *Gameboy) handleInterrupts() (cycles uint16) {
	if gb.cpu.interruptEnablePending {
		gb.cpu.interruptsEnabled = true
		gb.cpu.interruptEnablePending = false
		return 0
	}

	if !gb.cpu.interruptsEnabled && !gb.cpu.halted {
		return 0
	}

	mask := gb.memoryBus.regIE & gb.memoryBus.regIF

	if mask != 0 {
		gb.cpu.halted = false
	}

	if !gb.cpu.interruptsEnabled {
		return 0
	}

	// there are much nicer ways to write this, will probably refactor later
	// handle requested interrupt
	if mask&VBlankFlag != 0 {
		gb.memoryBus.regIF &= ^VBlankFlag
		gb.cpu.jumpToISRAddress(VBlankAddr)
		return 20
	}

	if mask&LCDStatFlag != 0 {
		gb.memoryBus.regIF &= ^LCDStatFlag
		gb.cpu.jumpToISRAddress(LCDStatAddr)
		return 20
	}

	if mask&TimerFlag != 0 {
		gb.memoryBus.regIF &= ^TimerFlag
		gb.cpu.jumpToISRAddress(TimerAddr)
		return 20
	}

	if mask&SerialFlag != 0 {
		gb.memoryBus.regIF &= ^SerialFlag
		gb.cpu.jumpToISRAddress(SerialAddr)
		return 20
	}

	if mask&JoypadFlag != 0 {
		gb.memoryBus.regIF &= ^JoypadFlag
		gb.cpu.jumpToISRAddress(JoypadAddr)
		return 20
	}

	return 0
}

func (gb *Gameboy) stepDivider(cycles uint16) {
	gb.divCycles += cycles
	if gb.divCycles >= divFreq {
		gb.divCycles -= divFreq
		gb.memoryBus.memory[divAddress]++
	}
}

func (gb *Gameboy) isTimerEnabled() bool {
	tac := gb.memoryBus.readAddress(tacAddress)
	return isBitSet(tac, 2)
}

func (gb *Gameboy) getTimerFreq() uint16 {
	tac := gb.memoryBus.readAddress(tacAddress) & 0b11
	switch tac {
	case 0b00:
		return 1024
	case 0b11:
		return 256
	case 0b10:
		return 64
	default:
		return 16
	}
}

func (gb *Gameboy) stepTimer(cycles uint16) {
	if !gb.isTimerEnabled() {
		return
	}

	gb.timerCycles += cycles
	freq := gb.getTimerFreq()
	for gb.timerCycles >= freq {
		gb.timerCycles -= freq
		tima := gb.memoryBus.readAddress(timaAddress)
		if tima == 0xff {
			tma := gb.memoryBus.readAddress(tmaAddress)
			gb.memoryBus.writeToAddress(timaAddress, tma)
			gb.memoryBus.requestInterrupt(TimerFlag)
		} else {
			gb.memoryBus.writeToAddress(timaAddress, tima+1)
		}
	}
}

func (gb *Gameboy) clearScreen() {
	for x := range ScreenWidth {
		for y := range ScreenHeight {
			gb.framebuffer[x][y][0] = 255
			gb.framebuffer[x][y][1] = 255
			gb.framebuffer[x][y][2] = 255
		}
	}
	gb.screenCleared = true
}

func (gb *Gameboy) Step(cycleBudget uint64) error {
	for gb.emuCycles < cycleBudget {
		cycles, err := gb.cpu.step()
		if err != nil {
			return err
		}
		gb.memoryBus.step(cycles)
		gb.apu.step(cycles)
		gb.emuCycles += uint64(cycles)
	}
	gb.emuCycles -= cycleBudget
	gb.speaker.EmitAudio(gb.apu.getSamples())
	return nil
}
