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

	HighRAMBegin uint16 = 0xff80
	HighRAMEnd   uint16 = 0xfffe

	AudioBegin uint16 = 0xff10
	AudioEnd   uint16 = 0xff26

	// relevant hardware addresses
	joypadAddress     uint16 = 0xff00
	divAddress        uint16 = 0xff04
	timaAddress       uint16 = 0xff05
	tmaAddress        uint16 = 0xff06
	tacAddress        uint16 = 0xff07
	ifAddress         uint16 = 0xff0f
	lcdControlAddress uint16 = 0xff40
	statAddress       uint16 = 0xff41
	scyAddress        uint16 = 0xff42
	scxAddress        uint16 = 0xff43
	lyAddress         uint16 = 0xff44
	lycAddress        uint16 = 0xff45
	oamDMAAddress     uint16 = 0xff46
	bgpAddress        uint16 = 0xff47
	wyAddress         uint16 = 0xff4a
	wxAddress         uint16 = 0xff4b
	ieAddress         uint16 = 0xffff
)

type MemoryBus struct {
	gb     *Gameboy
	cart   *Cartridge
	memory [0x10000]byte
	regIF  byte
	regIE  byte
}

func initMemoryBus(gb *Gameboy) *MemoryBus {
	mb := &MemoryBus{
		gb: gb,
	}
	return mb
}

func (mb *MemoryBus) step(cycles uint16) {
	// have GPU step and communicate VBlank and LCD status, trigger flags for
	// each if it comes back positive
	req := mb.gb.ppu.step(cycles)
	switch req {
	case vblankRequest:
		mb.setVBlank()
		screen := mb.gb.screen
		screen.Render(mb.gb.framebuffer)
		if mb.gb.joypad.update(screen.DoInput()) {
			mb.requestInterrupt(JoypadFlag)
		}
	case lcdStatRequest:
		mb.setLCD()
	case bothRequest:
		mb.setVBlank()
		mb.setLCD()
		screen := mb.gb.screen
		screen.Render(mb.gb.framebuffer)
		if mb.gb.joypad.update(screen.DoInput()) {
			mb.requestInterrupt(JoypadFlag)
		}
	}
	mb.cart.step()
}

func (mb *MemoryBus) dmaTransfer(val byte) {
	addr := uint16(val) << 8

	for i := range uint16(0xa0) {
		oamData := mb.readAddress(addr + i)
		mb.memory[OAMBegin+i] = oamData
	}
}

func (mb *MemoryBus) loadCartridge(romPath string) (int, error) {
	cart, err := initCartFromFile(romPath)
	if err != nil {
		return 0, err
	}
	mb.cart = cart
	return cart.mbc.romSize(), nil
}

func (mb *MemoryBus) readAddress(addr uint16) byte {
	if int(addr) >= len(mb.memory) {
		panic(fmt.Sprintf("out of bounds read: 0x%04X\n", addr))
	}

	switch {
	case addr <= RomBankNEnd, addr >= CartridgeRamBegin && addr <= CartridgeRamEnd:
		return mb.cart.read(addr)
	case addr == joypadAddress:
		return mb.gb.joypad.read()
	case addr == ieAddress:
		return mb.regIE | 0xe0
	case addr == ifAddress:
		return mb.regIF | 0xe0
	case addr == lyAddress:
		return mb.gb.ppu.line
	case addr == lycAddress:
		return mb.gb.ppu.lyc
	case addr == scxAddress:
		return mb.gb.ppu.scx
	case addr == scyAddress:
		return mb.gb.ppu.scy
	case addr == bgpAddress:
		return mb.gb.ppu.bgp
	case addr == lcdControlAddress:
		return mb.gb.ppu.lcdc
	case addr == statAddress:
		return mb.gb.ppu.stat
	case addr == wyAddress:
		return mb.gb.ppu.winY
	case addr == wxAddress:
		return mb.gb.ppu.winX + 7
	case addr >= vramBegin && addr <= vramEnd:
		return mb.gb.ppu.readVRAM(addr)
	case addr >= AudioBegin && addr <= WaveRamEnd:
		return mb.gb.apu.read(addr)
	}
	return mb.memory[addr]
}

func (mb *MemoryBus) writeToAddress(addr uint16, val byte) {
	switch {
	case addr <= RomBankNEnd, addr >= CartridgeRamBegin && addr <= CartridgeRamEnd:
		mb.cart.write(addr, val)
	case addr == oamDMAAddress:
		mb.dmaTransfer(val)
	case addr == joypadAddress:
		mb.gb.joypad.write(val)
	case addr == divAddress:
		// any write to this resets div to $00, and it's reset when STOP is executed
		mb.memory[divAddress] = 0
	case addr == tacAddress:
		// reset timer cycle counter on frequency change
		oldFreq := mb.gb.getTimerFreq()
		mb.memory[tacAddress] = val
		newFreq := mb.gb.getTimerFreq()
		if oldFreq != newFreq {
			mb.gb.timerCycles = 0
		}
	case addr == ieAddress:
		mb.regIE = val
	case addr == ifAddress:
		mb.regIF = val
	case addr == lyAddress:
		mb.gb.ppu.line = 0
	case addr == lycAddress:
		mb.gb.ppu.writeLYC(val)
	case addr == scxAddress:
		mb.gb.ppu.scx = val
	case addr == scyAddress:
		mb.gb.ppu.scy = val
	case addr == lcdControlAddress:
		mb.gb.ppu.lcdc = val
	case addr == statAddress:
		// bits 3 through 6 of STAT are read/write, the rest are either unused
		// or read only
		ppu := mb.gb.ppu
		ppu.stat = (ppu.stat & 0x07) | (val & 0x78)
	case addr == bgpAddress:
		mb.gb.ppu.bgp = val
	case addr == wxAddress:
		// WX is window X position + 7, store with the offset upfront for ergonomics
		// while also preventing underflow
		mb.gb.ppu.winX = byte(int16(val) - 7)
	case addr == wyAddress:
		mb.gb.ppu.winY = val
	case addr >= vramBegin && addr <= vramEnd:
		mb.gb.ppu.writeVRAM(addr, val)
	case addr >= AudioBegin && addr <= WaveRamEnd:
		mb.gb.apu.write(addr, val)
	default:
		mb.memory[addr] = val
	}
}
