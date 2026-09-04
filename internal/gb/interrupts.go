package gb 

const (
	// interrupt bits
	VBlankFlag byte = 0x01
	LCDStatFlag byte = 0x02
	TimerFlag byte = 0x04
	SerialFlag byte = 0x08
	JoypadFlag byte = 0x10

	// interrupt vector addresses for jumps 
	VBlankAddr = 0x40
	LCDStatAddr = 0x48
	TimerAddr = 0x50
	SerialAddr = 0x58
	JoypadAddr = 0x60
)

func (mb *MemoryBus) VBlankInterruptEnabled() bool {
	return bitEnabled(mb.IE, 0)
}

func (mb *MemoryBus) LCDStatInterruptEnabled() bool {
	return bitEnabled(mb.IE, 1)
}

func (mb *MemoryBus) TimerInterruptEnabled() bool {
	return bitEnabled(mb.IE, 2)
}

func (mb *MemoryBus) SerialInterruptEnabled() bool {
	return bitEnabled(mb.IE, 3)
}

func (mb *MemoryBus) JoypadInterruptEnabled() bool {
	return bitEnabled(mb.IE, 4)
}

func (mb *MemoryBus) SetVBlank() {
	mb.IF |= 1
}

func (mb *MemoryBus) SetLCD() {
	mb.IF |= (1 << 1)
}

func (mb *MemoryBus) SetTimer() {
	mb.IF |= (1 << 2)
}

func (mb *MemoryBus) SetSerial() {
	mb.IF |= (1 << 3)
}

func (mb *MemoryBus) SetJoypad() {
	mb.IF |= (1 << 4)
}

func (cpu *CPU) jumpToISRAddress(addr uint16) {
	cpu.IME = false
	cpu.push(cpu.PC)
	cpu.PC = addr
}

func (mb *MemoryBus) RequestInterrupt(interrupt byte) {
	req := mb.IF | interrupt
	mb.WriteToAddress(0xff0f, req)
}
