package gb

const (
	RomBankSize uint32 = 0x4000
	RamBankSize uint32 = 0x2000
)

type MBC interface {
	Read(addr uint16) byte
	Write(addr uint16, val byte)
	Size() int
}

type ROM struct {
	rom []byte
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

// MBC1: max 2MB ROM and/or 32KiB RAM
type MBC1 struct {
	rom          []byte
	lowRomBank   uint32
	upperRomBank uint32
	bank1        byte
	bank2        byte
	mode         bool
	ram          []byte
	ramEnabled   bool
	ramOffset    uint32
}

func InitMBC1(rom []byte, ramSize uint32) *MBC1 {
	mbc := &MBC1{
		rom:          rom,
		bank1:        1,
		upperRomBank: RomBankSize,
		bank2:        0,
		ram:          make([]byte, 0x8000),
	}
	mbc.computeRomOffsets()
	mbc.computeRamOffset()
	return mbc
}

func (mbc *MBC1) computeRomOffsets() {
	upper := mbc.bank2 << 5
	lower := mbc.bank1 & 0x1f

	if mbc.mode {
		mbc.lowRomBank = RomBankSize * uint32(upper)
	} else {
		mbc.lowRomBank = 0
	}
	mbc.upperRomBank = RomBankSize * uint32(lower|upper)
}

func (mbc *MBC1) computeRamOffset() {
	var bank uint32
	if mbc.mode {
		bank = uint32(mbc.bank2)
	}
	mbc.ramOffset = RamBankSize * bank
}

func effectiveRomAddr(bank uint32, addr uint16, romLen int) uint32 {
	return (bank | (uint32(addr) & 0x3fff)) & uint32((romLen - 1))
}

func (mbc *MBC1) effectiveRamAddr(addr uint16) uint32 {
	return (mbc.ramOffset | (uint32(addr) & 0x1fff)) & uint32(len(mbc.ram)-1)
}

func (mbc *MBC1) Read(addr uint16) byte {
	switch {
	case addr < 0x4000:
		romAddr := effectiveRomAddr(mbc.lowRomBank, addr, mbc.Size())
		return mbc.rom[romAddr]
	case addr < 0x8000:
		romAddr := effectiveRomAddr(mbc.upperRomBank, addr, mbc.Size())
		return mbc.rom[romAddr]
	case addr >= 0xa000 && addr <= 0xbfff:
		if !mbc.ramEnabled {
			return 0xff
		}
		ramAddr := mbc.effectiveRamAddr(addr)
		return mbc.ram[ramAddr]
	}
	return 0xff
}

func (mbc *MBC1) Write(addr uint16, val byte) {
	switch {
	case addr < 0x2000:
		mbc.ramEnabled = ((val & 0xf) == 0b1010)
	case addr < 0x4000:
		// bank 1 isn't allowed to be 0
		mbc.bank1 = val & 0x1f
		if mbc.bank1 == 0 {
			mbc.bank1++
		}
		mbc.computeRomOffsets()
	case addr < 0x6000:
		mbc.bank2 = val & 0b11
		mbc.computeRomOffsets()
		mbc.computeRamOffset()
	case addr < 0x8000:
		mbc.mode = bitEnabled(val, 0)
		mbc.computeRomOffsets()
		mbc.computeRamOffset()
	case addr >= 0xa000 && addr <= 0xbfff:
		if !mbc.ramEnabled {
			return
		}
		ramAddr := mbc.effectiveRamAddr(addr)
		mbc.ram[ramAddr] = val
	}
}

func (mbc *MBC1) Size() int {
	return len(mbc.rom)
}

// MBC2: max 256 KiB ROM, 512 x 4 bits RAM
type MBC2 struct {
	rom          []byte
	romBank      byte
	upperRomBank uint32
	ram          []byte // these function as half-bytes, upper 4 bits ignored
	ramEnabled   bool
}

func (mbc *MBC2) computeRomOffset() {
	mbc.upperRomBank = RomBankSize * uint32(mbc.romBank)
}

func (mbc *MBC2) effectiveRamAddr(addr uint16) uint32 {
	return uint32(addr) & 0x1ff
}

func InitMBC2(rom []byte) *MBC2 {
	mbc := &MBC2{rom: rom, romBank: 1, ram: make([]byte, 512)}
	mbc.computeRomOffset()
	return mbc
}

func (mbc *MBC2) Read(addr uint16) byte {
	switch {
	case addr < 0x4000:
		return mbc.rom[addr]
	case addr < 0x8000:
		romAddr := effectiveRomAddr(mbc.upperRomBank, addr, mbc.Size())
		return mbc.rom[romAddr]
	case addr >= 0xa000 && addr <= 0xbfff: // built-in RAM plus echoes
		if !mbc.ramEnabled {
			return 0xff
		}
		ramAddr := mbc.effectiveRamAddr(addr)
		return mbc.ram[ramAddr] | 0xf0
	}
	return 0xff
}

func (mbc *MBC2) Write(addr uint16, val byte) {
	switch {
	// RAM enable or ROM bank, depending on bit 8
	case addr < 0x4000:
		if addr&0x100 == 0 {
			mbc.ramEnabled = ((val & 0x0f) == 0xa)
		} else {
			// bank isn't allowed to be 0
			mbc.romBank = val & 0x0f
			if mbc.romBank == 0 {
				mbc.romBank++
			}
			mbc.computeRomOffset()
		}
	// builtin RAM and echoes, only use lower 4 bits
	case addr >= 0xa000 && addr <= 0xbfff:
		if mbc.ramEnabled {
			ramAddr := mbc.effectiveRamAddr(addr)
			mbc.ram[ramAddr] = (val & 0xf)
		}
	}
}

func (mbc *MBC2) Size() int {
	return len(mbc.rom)
}

type MBC3 struct {
	rom          []byte
	romBank      byte
	upperRomBank uint32

	ram       []byte
	ramOffset uint32

	rtc        [5]byte
	latchedRTC [5]byte
	latched    bool

	selectedBank byte

	ramAndTimerEnabled bool
}

func InitMBC3(rom []byte, ramSize uint32) *MBC3 {
	return &MBC3{rom: rom, romBank: 1, ram: make([]byte, ramSize)}
}

func (mbc *MBC3) computeRomOffsets() {
	mbc.upperRomBank = RomBankSize * uint32(mbc.romBank)
}

func (mbc *MBC3) computeRamOffset() {
	mbc.ramOffset = RamBankSize * uint32(mbc.selectedBank&0x07)
}

func (mbc *MBC3) effectiveRamAddr(addr uint16) uint32 {
	return (mbc.ramOffset | (uint32(addr) & 0x1fff)) & uint32(len(mbc.ram)-1)
}

func (mbc *MBC3) Read(addr uint16) byte {
	switch {
	case addr < 0x4000:
		return mbc.rom[addr]
	case addr < 0x8000:
		romAddr := effectiveRomAddr(mbc.upperRomBank, addr, mbc.Size())
		return mbc.rom[romAddr]
	case addr >= 0xa000 && addr <= 0xbfff:
		// RAM bank 00 - 07 or RTC register
		if !mbc.ramAndTimerEnabled {
			return 0xff
		}
		if mbc.selectedBank >= 0x08 && mbc.selectedBank <= 0x0c {
			return mbc.latchedRTC[mbc.selectedBank-0x08]
		}
		ramAddr := mbc.effectiveRamAddr(addr)
		return mbc.ram[ramAddr]
	}
	return 0xff
}

func (mbc *MBC3) Write(addr uint16, val byte) {
	switch {
	case addr < 0x2000:
		mbc.ramAndTimerEnabled = ((val & 0xf) == 0b1010)
	case addr < 0x4000:
		// ROM bank number
		// same as MBC1 but the whole 7 bits of the ROM bank number are written directly to this address
		// as with MBC1, val == 0x0 => select bank 1
		// other values 01 - 7f select the bank number
		mbc.romBank = val & 0x7f
		if mbc.romBank == 0 {
			mbc.romBank++
		}
		mbc.computeRomOffsets()
	case addr < 0x6000:
		// RAM bank number or RTC register select
		if val <= 0x0c {
			mbc.selectedBank = val
			mbc.computeRamOffset()
		}
	case addr < 0x8000:
		// latch clock data
		// if $00 THEN $01 is written, current time becomes latched into RTC registers
		if val == 0x00 {
			mbc.latched = false
		} else if val == 0x01 && mbc.latched == false {
			mbc.latched = true
			mbc.latchedRTC = mbc.rtc
		}
	case addr >= 0xa000 && addr <= 0xbfff:
		// RTC register 08 - 0c or RAM bank, depending on bank number/selected RTC register
		if !mbc.ramAndTimerEnabled {
			return
		}
		if mbc.selectedBank >= 0x08 && mbc.selectedBank <= 0x0c {
			mbc.rtc[mbc.selectedBank-0x08] = val
		} else if mbc.selectedBank <= 0x03 {
			ramAddr := mbc.effectiveRamAddr(addr)
			mbc.ram[ramAddr] = val
		}
	}
}

func (mbc *MBC3) Size() int {
	return len(mbc.rom)
}
