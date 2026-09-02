package gb

type MBC interface {
	Read(addr uint16) byte
	Write(addr uint16, val byte)
	Size() int
}

type ROM struct {
	rom 		[]byte
}

func InitROM(rom []byte) *ROM {
	return &ROM{
		rom: rom,
	} 
}

func (mbc *ROM) Read(addr uint16) byte {
	return mbc.rom[addr] 
}

func (mbc *ROM) Write(addr uint16, val byte) {
	// ROM only, ignore write to ROM space
}

func (mbc *ROM) Size() int {
	return len(mbc.rom)
}

type MBC1 struct {
	rom 		[]byte
	romBank 	byte
	romBanking  bool

	ram 		[]byte
	ramBank 	byte
	ramEnable 	bool
}

func InitMBC1(rom []byte) *MBC1 {
	return &MBC1{
		rom: rom,
		romBank: 1,
		ram: make([]byte, 0x8000),
	}
}

func (mbc *MBC1) Read(addr uint16) byte {
	switch {
		case addr < 0x4000:
			return mbc.rom[addr]
		case addr < 0x8000:
			bank := int(mbc.romBank) * 0x4000
			offset := int(addr - 0x4000)
			return mbc.rom[bank + offset]
		default:
			if !mbc.ramEnable {
				return 0xff
			}
			bank := int(mbc.ramBank) * 0x2000
			return mbc.ram[bank + int(addr - 0xa000)]
	}
}

func (mbc *MBC1) Write(addr uint16, val byte) {
	switch {
		case addr < 0x2000:
			// ram enable
			switch val & 0xf {
				case 0xa:
					mbc.ramEnable = true
				case 0x0:
					mbc.ramEnable = false
			}
		case addr < 0x4000:
			mbc.romBank = (mbc.romBank & 0xe0) | (val & 0x1f)
			mbc.UpdateROMBankIfUnusable()
		case addr < 0x6000:
			if mbc.romBanking {
				mbc.romBank = (mbc.romBank & 0x1f) | (val & 0xe0)
				mbc.UpdateROMBankIfUnusable()
			} else {
				mbc.ramBank = val & 0x3
			}
		case addr < 0x8000:
			mbc.romBanking = (val & 0x1 == 0x00)
			if mbc.romBanking {
				mbc.ramBank = 0
			} else {
				mbc.romBank = mbc.romBank & 0x1f
			}
	}
}

func (mbc *MBC1) Size() int {
	return len(mbc.rom)
}

func (mbc *MBC1) UpdateROMBankIfUnusable() {
	if mbc.romBank == 0x00 || mbc.romBank == 0x20 || mbc.romBank == 0x40 || mbc.romBank == 0x60 {
		mbc.romBank++
	} 
}
