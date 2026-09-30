package gb

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Cartridge struct {
	mbc      MBC
	name     string // currently dead, will use eventually
	filepath string
}

func (c *Cartridge) Read(addr uint16) byte {
	return c.mbc.Read(addr)
}

func (c *Cartridge) Write(addr uint16, val byte) {
	c.mbc.Write(addr, val)
}

func InitCartFromFile(path string) (*Cartridge, error) {
	cart := &Cartridge{}
	cart.filepath = path
	rom, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(rom) < 0x150 {
		return nil, fmt.Errorf("insufficient ROM size")
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

	// reference: https://gbdev.io/pandocs/The_Cartridge_Header.html
	switch mbcFlag {
	// ROM only
	case 0x00, 0x08:
		cart.mbc = InitROM(rom, false)
	case 0x09:
		cart.mbc = InitROM(rom, true)
	// MBC1
	case 0x01, 0x02:
		cart.mbc = InitMBC1(rom, ramSize, false)
	case 0x03:
		cart.mbc = InitMBC1(rom, ramSize, true)
	// MBC2
	case 0x05:
		cart.mbc = InitMBC2(rom, false)
	case 0x06:
		cart.mbc = InitMBC2(rom, true)
	// MBC3
	case 0x11, 0x12:
		cart.mbc = InitMBC3(rom, ramSize, false)
	case 0x0f, 0x10, 0x13:
		cart.mbc = InitMBC3(rom, ramSize, true)
	default:
		return nil, fmt.Errorf("unrecognized or unsupported cartridge type")
	}

	cart.LoadSave()
	return cart, nil
}

func (c *Cartridge) Save() {
	saveData := c.mbc.SaveData()
	if len(saveData) == 0 {
		// don't want cartridges with no battery to unnecessarily try to write save data
		return
	}
	// write RAM snapshot to a file called {cartridge path}.sav
	savePath := strings.TrimSuffix(c.filepath, filepath.Ext(c.filepath)) + ".sav"
	saveFile, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("failed to open file %s for writing save data\n", savePath)
		return
	}
	defer saveFile.Close()

	if _, err = saveFile.Write(saveData); err != nil {
		log.Printf("failed to write save data to file %s\n", savePath)
	}
}

func (c *Cartridge) LoadSave() {
	// can we find save file for this cartridge?
	savePath := strings.TrimSuffix(c.filepath, filepath.Ext(c.filepath)) + ".sav"
	saveFile, err := os.OpenFile(savePath, os.O_RDONLY, 0644)
	if err != nil {
		// no save data, we don't have to do anything
		return
	}
	defer saveFile.Close()

	saveData, err := os.ReadFile(savePath)
	if err != nil {
		log.Printf("failed to load save data from existing save file %s\n", savePath)
		return
	}

	c.mbc.LoadData(saveData)
}

func (c *Cartridge) Step() {
	// if we have a battery, check if saving to battery is requested
	if c.mbc.HasBattery() && c.mbc.SaveRequested() {
		// check if we need a save because RAM was switched from enable to disable
		c.Save()
		c.mbc.SatisfySaveRequest()
	}
}
