package gb

const (
	VRAMBegin uint16 = 0x8000
	VRAMEnd   uint16 = 0x9fff
	VRAMSize  uint16 = VRAMEnd - VRAMBegin + 1

	TilesSize int = 384

	ScreenHeight = 144
	ScreenWidth  = 160
)

type TilePixelValue int

const (
	Zero = iota
	One
	Two
	Three
)

// relevant hardware register bits
// LCDC
const (
	BGWindowEnableBit byte = iota
	OBJEnableBit
	OBJSizeBit
	BGTilemapBit
	BGWindowTileBit
	WindowEnableBit
	WindowTilemapBit
	LCDEnableBit
)

// STAT
const (
	CoincidenceBit byte = iota + 2
	Mode0IntSelectBit
	Mode1IntSelectBit
	Mode2IntSelectBit
	LYCIntSelectBit
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
	gb      		*Gameboy
	vram    		[VRAMSize]byte
	oam     		[OAMSize]byte
	tileSet 		[TilesSize]Tile
	tileScanline 	[ScreenWidth]byte
	cycles  		uint16
	mode    		PPUMode
	line    		byte // maps to ly
	lyc     		byte
	scy     		byte // scroll offset y
	scx     		byte // scroll offset x
	lcdc    		byte // lcd control register
	bgp     		byte // background palette
	stat 			byte // maps to STAT register
	winX			byte // maps to WX register at 0xff4b
	winY			byte // maps to WY register at 0xff4a

	lyEqualsLYC 	bool  // ly == lyc coincidence
	// pan docs says this is supposed to be here, will use it eventually
	yCondition		bool  // WY == LY condition maintained through frame
}

func InitPPU(gb *Gameboy) *PPU {
	ppu := &PPU{
		gb: gb, 
		mode: OAMAccess,
		lcdc: 0x91,
		stat: 0x85,
	}

	for i := range TilesSize {
		ppu.tileSet[i] = emptyTile()
	}

	return ppu
}

func (ppu *PPU) LYEqualsLYCInterruptEnabled() bool {
	return bitEnabled(ppu.stat, LYCIntSelectBit)
}

func (ppu *PPU) OAMInterruptEnabled() bool {
	return bitEnabled(ppu.stat, Mode2IntSelectBit)
}

func (ppu *PPU) VBlankInterruptEnabled() bool {
	return bitEnabled(ppu.stat, Mode1IntSelectBit)
}

func (ppu *PPU) HBlankInterruptEnabled() bool {
	return bitEnabled(ppu.stat, Mode0IntSelectBit)
}

func (ppu *PPU) setMode(mode PPUMode) {
	ppu.mode = mode
	ppu.stat = (ppu.stat &^ 0b11) | byte(mode)
}

// update graphics by individual frame
func (ppu *PPU) Step(cycles uint16) InterruptRequest {
	if !bitEnabled(ppu.lcdc, LCDEnableBit) {
		if !ppu.gb.screenCleared {
			ppu.gb.clearScreen()
		}
		ppu.cycles = 456
		ppu.line = 0
		ppu.setMode(HBlank)
		return NoRequest
	}
	ppu.gb.screenCleared = false

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
				ppu.setMode(VBlank)
				request.add(VBlankRequest)
				if ppu.VBlankInterruptEnabled() {
					request.add(LCDStatRequest)
				}
			} else {
				ppu.setMode(OAMAccess)
				if ppu.OAMInterruptEnabled() {
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
				ppu.setMode(OAMAccess)
				ppu.line = 0
				if ppu.OAMInterruptEnabled() {
					request.add(LCDStatRequest)
				}
			}
			ppu.setEqualLinesCheck(&request)
		}
	case OAMAccess:
		if ppu.cycles >= 80 {
			ppu.cycles -= 80
			ppu.setMode(VRAMAccess)
		}
	case VRAMAccess:
		if ppu.cycles >= 172 {
			ppu.cycles -= 172
			if ppu.HBlankInterruptEnabled() {
				request.add(LCDStatRequest)
			}
			ppu.DrawScanline()
			ppu.setMode(HBlank)
		}
	}
	return request
}

// helper for DrawFrame to compute index in the tileset cache
func tileSetIndex(tileID byte, unsignedAddressingEnabled bool) int {
	if unsignedAddressingEnabled {
		return int(tileID)
	}
	return 256 + int(int8(tileID))
}

func applyPalette(colorIdx TilePixelValue, bgp byte) byte {
	shade := (bgp >> (uint(colorIdx) * 2)) & 0x03
	shades := [4]byte{0xff, 0xaa, 0x55, 0x00}
	return shades[shade]
}

func (ppu *PPU) setTilePixel(x, y, colorNum byte) {
	shade := applyPalette(TilePixelValue(colorNum), ppu.bgp)
	ppu.gb.Framebuffer[x][y][0] = shade
	ppu.gb.Framebuffer[x][y][1] = shade
	ppu.gb.Framebuffer[x][y][2] = shade
	ppu.tileScanline[x] = colorNum
}

func (ppu *PPU) DrawTiles(scanline byte) {
	var inWindow bool
	if bitEnabled(ppu.lcdc, WindowEnableBit) && ppu.winY <= ppu.line {
		inWindow = true
	}

	var tileData uint16 = 0x8800
	var unsignedBytes bool
	if bitEnabled(ppu.lcdc, BGWindowTileBit) {
		tileData = 0x8000
		unsignedBytes = true
	}

	var tileMapBit byte = BGTilemapBit
	var backgroundMemory uint16 = 0x9800
	if inWindow {
		tileMapBit = WindowTilemapBit
	}
	if bitEnabled(ppu.lcdc, tileMapBit) {
		backgroundMemory = 0x9c00
	}

	var yPos byte
	if !inWindow {
		yPos = ppu.scy + scanline
	} else {
		yPos = scanline - ppu.winY
	}

	var tileRow = uint16(yPos / 8) * 32

	for pixel := range byte(160) {
		x := pixel + ppu.scx
		if inWindow && pixel >= ppu.winX {
			x = pixel - ppu.winX
		}

		tileCol := uint16(x / 8)
		tileAddr := backgroundMemory + tileRow + tileCol

		tileLocation := tileData
		var tileNum int16
		if unsignedBytes {
			tileNum = int16(ppu.ReadVRAM(tileAddr))
			tileLocation += uint16(tileNum * 16)
		} else {
			tileNum = int16(int8(ppu.ReadVRAM(tileAddr)))
			tileLocation = uint16(int32(tileLocation) + int32((tileNum + 128) * 16))
		}

		line := (yPos % 8) * 2
		tileData1 := ppu.ReadVRAM(tileLocation + uint16(line))
		tileData2 := ppu.ReadVRAM(tileLocation + uint16(line) + 1)

		colorBit := 7 - (x % 8)
		colorNum := ((tileData2 >> colorBit) & 1) << 1 | (tileData1 >> colorBit) & 1
		ppu.setTilePixel(pixel, scanline, colorNum)
	}
}

func (ppu *PPU) DrawScanline() {
	if bitEnabled(ppu.lcdc, BGWindowEnableBit) {
		ppu.DrawTiles(ppu.line)
	}
	// TODO: implement sprite rendering
}

func (ppu *PPU) setEqualLinesCheck(request *InterruptRequest) {
	equal := ppu.line == ppu.lyc
	if equal {
		SetBit(&ppu.stat, CoincidenceBit)
		if ppu.LYEqualsLYCInterruptEnabled() {
			request.add(LCDStatRequest)
		}
	} else {
		UnsetBit(&ppu.stat, CoincidenceBit)
	}
	ppu.lyEqualsLYC = equal
}

// update LYC and LY == LYC coincidence
func (ppu *PPU) WriteLYC(val byte) {
	ppu.lyc = val
	equal := ppu.lyc == ppu.line
	if equal {
		SetBit(&ppu.stat, CoincidenceBit)
	} else {
		UnsetBit(&ppu.stat, CoincidenceBit)
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
	tileByte2 := ppu.vram[normalizedAddress+1]

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
