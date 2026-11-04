package gb

const (
	// interrupt bits
	VBlankFlag  byte = 0x01
	LCDStatFlag byte = 0x02
	TimerFlag   byte = 0x04
	SerialFlag  byte = 0x08
	JoypadFlag  byte = 0x10

	// interrupt vector addresses for jumps
	VBlankAddr  = 0x40
	LCDStatAddr = 0x48
	TimerAddr   = 0x50
	SerialAddr  = 0x58
	JoypadAddr  = 0x60
)

func (mb *MemoryBus) setVBlank() {
	setBit(&mb.regIF, 0)
}

func (mb *MemoryBus) setLCD() {
	setBit(&mb.regIF, 1)
}

func (cpu *CPU) jumpToISRAddress(addr uint16) {
	cpu.interruptsEnabled = false
	cpu.push(cpu.pc)
	cpu.pc = addr
}

func (mb *MemoryBus) requestInterrupt(interrupt byte) {
	req := mb.regIF | interrupt
	mb.writeToAddress(0xff0f, req)
}
