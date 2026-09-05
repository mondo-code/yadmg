package gb

type Gameboy struct {
	CPU 		*CPU
	PPU			*PPU
	MemoryBus 	*MemoryBus
	TimerCycles uint16
	DivCycles	uint16
	paused		bool
	LogOpcodes	bool
}

const DivFreq = 0xff

func InitGameboy(romPath string) (*Gameboy, error) {
	gb := &Gameboy{}
	gb.CPU = InitCPU(gb)
	mb := InitMemoryBus(gb)
	gb.MemoryBus = mb
	ppu := InitPPU(gb)
	gb.PPU = ppu
	_, err := mb.LoadCartridge(romPath)
	if err != nil {
		return nil, err
	}

	return gb, nil
}

func (gb *Gameboy) HandleInterrupts() (cycles uint16) {
	if gb.CPU.IMEPending {
		gb.CPU.IME = true
		gb.CPU.IMEPending = false
		return 0
	}

	if !gb.CPU.IME && !gb.CPU.Halted {
		return 0
	}

	mask := gb.MemoryBus.IE & gb.MemoryBus.IF

	if mask != 0 {
		gb.CPU.Halted = false 
	}

	if !gb.CPU.IME {
		return 0
	}

	// there are much nicer ways to write this, will probably refactor later
	// handle requested interrupt 
	if mask & VBlankFlag != 0 {
		gb.MemoryBus.IF &= ^VBlankFlag 
		gb.CPU.jumpToISRAddress(VBlankAddr)
		return 20
	}

	if mask & LCDStatFlag != 0 {
		gb.MemoryBus.IF &= ^LCDStatFlag 
		gb.CPU.jumpToISRAddress(LCDStatAddr)
		return 20
	}

	if mask & TimerFlag != 0 {
		gb.MemoryBus.IF &= ^TimerFlag
		gb.CPU.jumpToISRAddress(TimerAddr)
		return 20
	}

	if mask & SerialFlag != 0 {
		gb.MemoryBus.IF &= ^SerialFlag
		gb.CPU.jumpToISRAddress(SerialAddr)
		return 20
	}

	if mask & JoypadFlag != 0 {
		gb.MemoryBus.IF &= ^JoypadFlag
		gb.CPU.jumpToISRAddress(JoypadAddr)
		return 20
	}

	return 0
}

func (gb *Gameboy) StepDivider(cycles uint16) {
	gb.DivCycles += cycles
	if gb.DivCycles >= DivFreq {
		gb.DivCycles -= DivFreq 
		gb.MemoryBus.memory[DIVAddress]++
	}
}

func (gb *Gameboy) IsTimerEnabled() bool {
	tac := gb.MemoryBus.ReadAddress(TACAddress)
	return bitEnabled(tac, 2)
}

func (gb *Gameboy) GetTimerFreq() uint16 {
	tac := gb.MemoryBus.ReadAddress(TACAddress) & 0b11
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

func (gb *Gameboy) StepTimer(cycles uint16) {
	if !gb.IsTimerEnabled() {
		return
	}

	gb.TimerCycles += cycles
	freq := gb.GetTimerFreq()
	for gb.TimerCycles >= freq {
		gb.TimerCycles -= freq
		tima := gb.MemoryBus.ReadAddress(TIMAAddress)
		if tima == 0xff {
			tma := gb.MemoryBus.ReadAddress(TMAAddress)
			gb.MemoryBus.WriteToAddress(TIMAAddress, tma)
			gb.MemoryBus.RequestInterrupt(TimerFlag)
		} else {
			gb.MemoryBus.WriteToAddress(TIMAAddress, tima + 1)
		}
	}
}

func (gb *Gameboy) Step() error {
	cycles, err := gb.CPU.Step()
	if err != nil {
		return err
	}
	gb.MemoryBus.Step(cycles)
	return nil
}
