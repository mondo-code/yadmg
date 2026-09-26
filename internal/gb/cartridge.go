package gb

import (
	"os"
)

type Cartridge struct {
	memoryBus *MemoryBus
	mbc       MBC
	name      string
	filepath  string
}

func (c *Cartridge) Read(addr uint16) byte {
	return c.mbc.Read(addr)
}

func (c *Cartridge) Write(addr uint16, val byte) {
	c.mbc.Write(addr, val)
}

func InitCartFromFile(path string, bus *MemoryBus) (*Cartridge, error) {
	cart := &Cartridge{}
	rom, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// handle cartridge header data
	mbcFlag := rom[0x147]
	ramSizeFlag := rom[0x149]

	var ramSize uint32
	switch ramSizeFlag {
	case 0x00, 0x01:
		ramSize = 0
	case 0x02:
		ramSize = 8192
	case 0x03:
		ramSize = 32768
	case 0x04:
		ramSize = 131072
	case 0x05:
		ramSize = 65536
	}

	switch mbcFlag {
	case 0x00, 0x08, 0x09:
		cart.mbc = InitROM(rom)
	default:
		switch {
		case mbcFlag <= 0x03:
			cart.mbc = InitMBC1(rom, ramSize)
		case mbcFlag <= 0x06:
			cart.mbc = InitMBC2(rom)
		case mbcFlag <= 0x13:
			cart.mbc = InitMBC3(rom, ramSize)
		}
	}

	return cart, nil
}
