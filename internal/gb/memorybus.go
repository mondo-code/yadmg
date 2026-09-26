package gb

import "fmt"

// memory map definitions
const (
	BootRomBegin uint16 = 0x0000
	BootRomEnd   uint16 = 0x00ff

	RomBank0Begin uint16 = 0x0000
	RomBank0End   uint16 = 0x3fff

	RomBankNStart uint16 = 0x4000
	RomBankNEnd   uint16 = 0x7fff

	TileRamBegin uint16 = 0x8000
	TileRamEnd   uint16 = 0x97ff

	BackgroundMapBegin uint16 = 0x9800
	BackgroundMapEnd   uint16 = 0x9fff

	CartridgeRamBegin uint16 = 0xa000
	CartridgeRamEnd   uint16 = 0xbfff

	WorkingRamBegin uint16 = 0xc000
	WorkingRamEnd   uint16 = 0xdfff

	EchoRamBegin uint16 = 0xe000
	EchoRamEnd   uint16 = 0xfdff

	OAMBegin uint16 = 0xfe00
	OAMEnd   uint16 = 0xfe9f
	OAMSize  uint16 = OAMEnd - OAMBegin + 1

	UnusedBegin uint16 = 0xfea0
	UnusedEnd   uint16 = 0xfeff

	IORegistersBegin uint16 = 0xff00
	IORegistersEnd   uint16 = 0xff7f

	HighRAMBegin uint16 = 0xff80
	HighRAMEnd   uint16 = 0xfffe

	// relevant hardware addresses
	JoypadAddress     uint16 = 0xff00
	DIVAddress        uint16 = 0xff04
	TIMAAddress       uint16 = 0xff05
	TMAAddress        uint16 = 0xff06
	TACAddress        uint16 = 0xff07
	IFAddress         uint16 = 0xff0f
	LCDControlAddress uint16 = 0xff40
	STATAddress       uint16 = 0xff41
	SCYAddress        uint16 = 0xff42
	SCXAddress        uint16 = 0xff43
	LYAddress         uint16 = 0xff44
	LYCAddress        uint16 = 0xff45
	OAMDMAAddress     uint16 = 0xff46
	BGPAddress        uint16 = 0xff47
	WYAddress         uint16 = 0xff4a
	WXAddress         uint16 = 0xff4b
	IEAddress         uint16 = 0xffff
)

type MemoryBus struct {
	gb     *Gameboy
	cart   *Cartridge
	memory [0x10000]byte
	IF     byte
	IE     byte
}

func InitMemoryBus(gb *Gameboy) *MemoryBus {
	mb := &MemoryBus{
		gb: gb,
	}
	return mb
}

func (mb *MemoryBus) Step(cycles uint16) {
	// have GPU step and communicate VBlank and LCD status, trigger flags for
	// each if it comes back positive
	req := mb.gb.PPU.Step(cycles)
	switch req {
	case VBlankRequest:
		mb.SetVBlank()
		screen := mb.gb.screen
		screen.Render(mb.gb.Framebuffer)
		if mb.gb.joypad.Update(screen.DoInput()) {
			mb.RequestInterrupt(JoypadFlag)
		}
	case LCDStatRequest:
		mb.SetLCD()
	case BothRequest:
		mb.SetVBlank()
		mb.SetLCD()
		screen := mb.gb.screen
		screen.Render(mb.gb.Framebuffer)
		if mb.gb.joypad.Update(screen.DoInput()) {
			mb.RequestInterrupt(JoypadFlag)
		}
	}
}

func (mb *MemoryBus) DMATransfer(val byte) {
	addr := uint16(val) << 8

	for i := range uint16(0xa0) {
		oamData := mb.ReadAddress(addr + i)
		mb.memory[OAMBegin+i] = oamData
	}
}

func (mb *MemoryBus) LoadCartridge(romPath string) (int, error) {
	cart, err := InitCartFromFile(romPath, mb)
	if err != nil {
		return 0, err
	}
	mb.cart = cart
	return cart.mbc.Size(), nil
}

func (mb *MemoryBus) ReadAddress(addr uint16) byte {
	if int(addr) >= len(mb.memory) {
		panic(fmt.Sprintf("out of bounds read: 0x%04X\n", addr))
	}

	switch {
	case addr <= RomBankNEnd, addr >= CartridgeRamBegin && addr <= CartridgeRamEnd:
		return mb.cart.Read(addr)
	case addr == JoypadAddress:
		return mb.gb.joypad.Read()
	case addr == IEAddress:
		return mb.IE
	case addr == IFAddress:
		return mb.IF
	case addr == LYAddress:
		return mb.gb.PPU.line
	case addr == LYCAddress:
		return mb.gb.PPU.lyc
	case addr == SCXAddress:
		return mb.gb.PPU.scx
	case addr == SCYAddress:
		return mb.gb.PPU.scy
	case addr == BGPAddress:
		return mb.gb.PPU.bgp
	case addr == LCDControlAddress:
		return mb.gb.PPU.lcdc
	case addr == STATAddress:
		return mb.gb.PPU.stat
	case addr == WYAddress:
		return mb.gb.PPU.winY
	case addr == WXAddress:
		return mb.gb.PPU.winX + 7
	case addr >= VRAMBegin && addr <= VRAMEnd:
		return mb.gb.PPU.ReadVRAM(addr)
	}
	return mb.memory[addr]
}

func (mb *MemoryBus) WriteToAddress(addr uint16, val byte) {
	switch {
	case addr <= RomBankNEnd, addr >= CartridgeRamBegin && addr <= CartridgeRamEnd:
		mb.cart.Write(addr, val)
	case addr == OAMDMAAddress:
		mb.DMATransfer(val)
	case addr == JoypadAddress:
		mb.gb.joypad.Write(val)
	case addr == DIVAddress:
		// any write to this resets div to $00, and it's reset when STOP is executed
		mb.memory[DIVAddress] = 0
	case addr == TACAddress:
		// reset timer cycle counter on frequency change
		oldFreq := mb.gb.GetTimerFreq()
		mb.memory[TACAddress] = val
		newFreq := mb.gb.GetTimerFreq()
		if oldFreq != newFreq {
			mb.gb.TimerCycles = 0
		}
	case addr == IEAddress:
		mb.IE = val
	case addr == IFAddress:
		mb.IF = val
	case addr == LYAddress:
		mb.gb.PPU.line = 0
	case addr == LYCAddress:
		mb.gb.PPU.WriteLYC(val)
	case addr == SCXAddress:
		mb.gb.PPU.scx = val
	case addr == SCYAddress:
		mb.gb.PPU.scy = val
	case addr == LCDControlAddress:
		mb.gb.PPU.lcdc = val
	case addr == STATAddress:
		// bits 3 through 6 of STAT are read/write, the rest are either unused
		// or read only
		ppu := mb.gb.PPU
		ppu.stat = (ppu.stat & 0x07) | (val & 0x78)
	case addr == BGPAddress:
		mb.gb.PPU.bgp = val
	case addr == WXAddress:
		// WX is window X position + 7, store with the offset upfront for ergonomics
		// while also preventing underflow
		mb.gb.PPU.winX = byte(int16(val) - 7)
	case addr == WYAddress:
		mb.gb.PPU.winY = val
	case addr >= VRAMBegin && addr <= VRAMEnd:
		mb.gb.PPU.WriteVRAM(addr, val)
	default:
		mb.memory[addr] = val
	}
}
