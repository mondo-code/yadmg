package gb 

import (
	"fmt"
	"os"
)

type Cartridge struct {
	memoryBus 	*MemoryBus
	mbc 		MBC
	name 		string
	filepath 	string
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

	if len(rom) > len(bus.memory) {
		return nil, fmt.Errorf("size of ROM exceeds size of memory") 
	}

	// determine MBC type from cartridge header
	mbcFlag := rom[0x0147]
	switch mbcFlag {
	case 0x00, 0x08, 0x09, 0x0b, 0x0c, 0x0d:
		cart.mbc = InitROM(rom)
	default:
		switch {
			case mbcFlag <= 0x03:
				cart.mbc = InitMBC1(rom)
			// TODO: more MBC types
		}
	}

	return cart, nil
}
