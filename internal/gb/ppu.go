package gb 

const (
	VRAMBegin uint16 = 0x8000
	VRAMEnd uint16 = 0x9fff
	VRAMSize uint16 = VRAMEnd - VRAMBegin + 1

	TilesSize int = 384
)

type TilePixelValue int
const (
	Zero = iota
	One
	Two
	Three
)

type Tile [8][8]TilePixelValue

func emptyTile() Tile {
	var t Tile
	for i := range 8 {
		for j := range 8 {
			t[i][j] = Zero 
		}
	}
	return t
}

type InterruptRequest int
const (
	NoRequest InterruptRequest = iota
	VBlankRequest
	LCDStatRequest
	BothRequest
)

func (r *InterruptRequest) add(other InterruptRequest) {
	switch {
		case *r == NoRequest:
			*r = other
		case *r == VBlankRequest && other == LCDStatRequest:
			*r = BothRequest
		case *r == LCDStatRequest && other == VBlankRequest:
			*r = BothRequest
	}
}

type Screen interface {
	Render(framebuffer *[160][144][3]uint8)
	IsRunning() bool
}

type PPUMode int 
const (
	HBlank PPUMode = iota	
	VBlank
	OAMAccess
	VRAMAccess
)

type PPU struct {
	gb 		*Gameboy
	vram 	[VRAMSize]byte
	oam 	[OAMSize]byte
	tileSet [TilesSize]Tile
	cycles 	uint16
	line 	byte  // maps to ly at 0xff44
	lyc		byte  // maps to lyc at 0xff45
	mode 	PPUMode

	lcdEnabled 					bool
	lyEqualsLYCInterruptEnabled bool
	oamInterruptEnabled 		bool
	vblankInterruptEnabled 		bool
	hblankInterruptEnabled 		bool
	lyEqualsLYC 				bool
}

func InitPPU(gb *Gameboy) *PPU {
	ppu := &PPU{gb: gb, mode: HBlank}

	for i := range TilesSize {
		ppu.tileSet[i] = emptyTile()
	}

	return ppu
}

// update graphics by individual frame
func (ppu *PPU) Step(cycles uint16) InterruptRequest {
	request := NoRequest
	ppu.cycles += cycles

	// set mode in STAT register depending on cycles of current line
	// OAM 80 cycles, mode3/drawing is 172, hblank is 204
	switch ppu.mode {
		case HBlank:
			if ppu.cycles >= 204 {
				ppu.cycles -= 204
				ppu.line += 1

				if ppu.line >= 144 {
					ppu.mode = VBlank
					request.add(VBlankRequest)
					if ppu.vblankInterruptEnabled {
						request.add(LCDStatRequest)
					}
				} else {
					ppu.mode = OAMAccess
					if ppu.oamInterruptEnabled {
						request.add(LCDStatRequest)
					}
				}
				ppu.setEqualLinesCheck(&request)
			}
		case VBlank:
			if ppu.cycles >= 456 {
				ppu.cycles -= 456
				ppu.line += 1
				if ppu.line == 154 {
					ppu.mode = OAMAccess
					ppu.line = 0
					if ppu.oamInterruptEnabled {
						request.add(LCDStatRequest)
					}
				}
				ppu.setEqualLinesCheck(&request)
			}
		case OAMAccess:
			if ppu.cycles >= 80 {
				ppu.cycles -= 80
				ppu.mode = VRAMAccess
			}
		case VRAMAccess:
			if ppu.cycles >= 172 {
				ppu.cycles -= 172
				if ppu.hblankInterruptEnabled {
					request.add(LCDStatRequest)
				}
				ppu.mode = HBlank
			}
	}
	return request
}

func (ppu *PPU) setEqualLinesCheck(request *InterruptRequest) {
	equal := ppu.line == ppu.lyc
	if equal && ppu.lyEqualsLYCInterruptEnabled {
		request.add(LCDStatRequest)
	}
	ppu.lyEqualsLYC = equal
}

func (ppu *PPU) ReadVRAM(addr uint16) byte {
	return ppu.vram[addr - VRAMBegin]
}

func (ppu *PPU) WriteVRAM(addr uint16, val byte) {
	physicalAddr := addr - VRAMBegin
	ppu.vram[physicalAddr] = val
	// the complicated part of this is caching our tile set in the tileSet field
	// we aren't writing to tileset storage if the address goes past 0x1800
	if physicalAddr >= 0x1800 {
		return
	}

	// tiles' rows are encoded in two bytes with the first byte always on an even address,
	// so doing a bitwise and on the address with 0xfffe gives us the first byte.
	normalizedAddress := physicalAddr & 0xfffe
	tileByte1 := ppu.vram[normalizedAddress]
	tileByte2 := ppu.vram[normalizedAddress + 1]

	// tile is 8 rows tall, since each row is encoded with two bytes a tile is 16 bytes in total
	tileIndex := physicalAddr / 16
	rowIndex := (physicalAddr % 16) / 2

	// get 8 pixels that make up a given row
	for pixelIndex := range 8 {
		// first, find corresponding bit that encodes pixel's value 
		var mask byte = 1 << (7 - pixelIndex)
		lsb := tileByte1 & mask
		msb := tileByte2 & mask

		var value TilePixelValue 
		if lsb != 0 && msb != 0 {
			value = Three
		} else if lsb == 0 && msb != 0 {
			value = Two
		} else if lsb != 0 && msb == 0 {
			value = One
		} else {
			value = Zero
		}
		ppu.tileSet[tileIndex][rowIndex][pixelIndex] = value
	}
}

// TODO: implement
func (ppu *PPU) WriteOAM(addr uint16, val byte) {}
