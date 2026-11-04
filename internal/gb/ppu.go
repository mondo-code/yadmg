package gb

const (
	vramBegin uint16 = 0x8000
	vramEnd   uint16 = 0x9fff
	vramSize  uint16 = vramEnd - vramBegin + 1

	tilesSize int = 384

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
	bgWindowEnableBit byte = iota
	objEnableBit
	objSizeBit
	bgTilemapBit
	bgWindowTileBit
	windowEnableBit
	windowTilemapBit
	lcdEnableBit
)

// STAT
const (
	coincidenceBit byte = iota + 2
	mode0IntSelectBit
	mode1IntSelectBit
	mode2IntSelectBit
	lycIntSelectBit
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

type interruptRequest int

const (
	noRequest interruptRequest = iota
	vblankRequest
	lcdStatRequest
	bothRequest
)

func (r *interruptRequest) add(other interruptRequest) {
	switch {
	case *r == noRequest:
		*r = other
	case *r == vblankRequest && other == lcdStatRequest:
		*r = bothRequest
	case *r == lcdStatRequest && other == vblankRequest:
		*r = bothRequest
	}
}

type ppuMode int

const (
	hblank ppuMode = iota
	vblank
	oamAccess
	vramAccess
)

type PPU struct {
	gb           *Gameboy
	vram         [vramSize]byte
	tileSet      [tilesSize]Tile
	tileScanline [ScreenWidth]byte
	cycles       uint16
	mode         ppuMode
	line         byte // maps to ly
	lyc          byte
	scy          byte // scroll offset y
	scx          byte // scroll offset x
	lcdc         byte // lcd control register
	bgp          byte // background palette
	stat         byte // maps to STAT register
	winX         byte // maps to WX register at 0xff4b
	winY         byte // maps to WY register at 0xff4a
	windowLine   byte // internal window line counter; advances only on lines the window was drawn

	lyEqualsLYC bool // ly == lyc coincidence
	// pan docs says this is supposed to be here, will (probably) use it eventually
	yCondition bool // WY == LY condition maintained through frame
}

func initPPU(gb *Gameboy) *PPU {
	ppu := &PPU{
		gb:   gb,
		mode: oamAccess,
		lcdc: 0x91,
		stat: 0x85,
	}

	for i := range tilesSize {
		ppu.tileSet[i] = emptyTile()
	}

	return ppu
}

func (ppu *PPU) lyEqualsLYCInterruptEnabled() bool {
	return isBitSet(ppu.stat, lycIntSelectBit)
}

func (ppu *PPU) oamInterruptEnabled() bool {
	return isBitSet(ppu.stat, mode2IntSelectBit)
}

func (ppu *PPU) vblankInterruptEnabled() bool {
	return isBitSet(ppu.stat, mode1IntSelectBit)
}

func (ppu *PPU) hblankInterruptEnabled() bool {
	return isBitSet(ppu.stat, mode0IntSelectBit)
}

func (ppu *PPU) setMode(mode ppuMode) {
	ppu.mode = mode
	ppu.stat = (ppu.stat &^ 0b11) | byte(mode)
}

// update graphics by individual frame
func (ppu *PPU) step(cycles uint16) interruptRequest {
	if !isBitSet(ppu.lcdc, lcdEnableBit) {
		if !ppu.gb.screenCleared {
			ppu.gb.clearScreen()
		}
		ppu.cycles = 456
		ppu.line = 0
		ppu.windowLine = 0
		ppu.setMode(hblank)
		return noRequest
	}
	ppu.gb.screenCleared = false

	request := noRequest
	ppu.cycles += cycles

	// set mode in STAT register depending on cycles of current line
	// OAM 80 cycles, mode3/drawing is 172, hblank is 204
	switch ppu.mode {
	case hblank:
		if ppu.cycles >= 204 {
			ppu.cycles -= 204
			ppu.line += 1

			if ppu.line >= 144 {
				ppu.setMode(vblank)
				request.add(vblankRequest)
				if ppu.vblankInterruptEnabled() {
					request.add(lcdStatRequest)
				}
			} else {
				ppu.setMode(oamAccess)
				if ppu.oamInterruptEnabled() {
					request.add(lcdStatRequest)
				}
			}
			ppu.setEqualLinesCheck(&request)
		}
	case vblank:
		if ppu.cycles >= 456 {
			ppu.cycles -= 456
			ppu.line += 1
			if ppu.line == 154 {
				ppu.setMode(oamAccess)
				ppu.line = 0
				ppu.windowLine = 0
				if ppu.oamInterruptEnabled() {
					request.add(lcdStatRequest)
				}
			}
			ppu.setEqualLinesCheck(&request)
		}
	case oamAccess:
		if ppu.cycles >= 80 {
			ppu.cycles -= 80
			ppu.setMode(vramAccess)
		}
	case vramAccess:
		if ppu.cycles >= 172 {
			ppu.cycles -= 172
			if ppu.hblankInterruptEnabled() {
				request.add(lcdStatRequest)
			}
			ppu.drawScanline()
			ppu.setMode(hblank)
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

func (ppu *PPU) setTilePixel(x, y, colorNum, palette byte) {
	ppu.setPixel(x, y, colorNum, palette, true)
	ppu.tileScanline[x] = colorNum
}

func (ppu *PPU) setPixel(x, y, colorNum, palette byte, priority bool) {
	shade := applyPalette(TilePixelValue(colorNum), palette)
	if ppu.tileScanline[x] == 0 || priority {
		ppu.gb.framebuffer[x][y][0] = shade
		ppu.gb.framebuffer[x][y][1] = shade
		ppu.gb.framebuffer[x][y][2] = shade
	}
}

func (ppu *PPU) tileRowAddress(mapAddr uint16, row byte) uint16 {
	tileID := ppu.readVRAM(mapAddr)
	var base uint16 = 0x8000
	if !isBitSet(ppu.lcdc, bgWindowTileBit) {
		// signed mode: IDs 0..127 -> 0x9000..0x97ff, IDs 128..255 (-128..-1) -> 0x8800..0x8fff.
		// flipping the top bit turns the signed ID into an unsigned offset from 0x8800.
		base = 0x8800
		tileID ^= 0x80
	}
	return base + uint16(tileID)<<4 + uint16(row)<<1
}

func (ppu *PPU) drawTileRow(addr uint16, col, count, x, y byte) {
	lo := ppu.readVRAM(addr)
	hi := ppu.readVRAM(addr + 1)
	for i := range count {
		bit := 7 - (col + i)
		colorNum := (hi>>bit&1)<<1 | lo>>bit&1
		ppu.setTilePixel(x+i, y, colorNum, ppu.bgp)
	}
}

func (ppu *PPU) drawBackground(scanline, endX byte) {
	mapBase := uint16(0x9800)
	if isBitSet(ppu.lcdc, bgTilemapBit) {
		mapBase = 0x9c00
	}

	mapY := scanline + ppu.scy
	rowBase := mapBase + uint16(mapY>>3)<<5
	tileRow := mapY & 7

	var x byte
	for x < endX {
		mapX := x + ppu.scx
		col := mapX & 7
		run := min(8-col, endX-x) // clamp so the last tile can't run past endX

		addr := ppu.tileRowAddress(rowBase+uint16(mapX>>3), tileRow)
		ppu.drawTileRow(addr, col, run, x, scanline)
		x += run
	}
}

func (ppu *PPU) drawWindow(scanline, startX byte) {
	mapBase := uint16(0x9800)
	if isBitSet(ppu.lcdc, windowTilemapBit) {
		mapBase = 0x9c00
	}

	mapY := ppu.windowLine
	rowBase := mapBase + uint16(mapY>>3)<<5
	tileRow := mapY & 7

	x := startX
	for x < ScreenWidth {
		// winX holds WX-7 wrapped to a byte, so when WX < 7 (left edge offscreen)
		// this subtraction wraps to the right positive offset into the window.
		mapX := x - ppu.winX
		col := mapX & 7
		run := min(8-col, ScreenWidth-x)

		addr := ppu.tileRowAddress(rowBase+uint16(mapX>>3), tileRow)
		ppu.drawTileRow(addr, col, run, x, scanline)
		x += run
	}
}

func (ppu *PPU) windowStartX() byte {
	if !isBitSet(ppu.lcdc, windowEnableBit) || ppu.line < ppu.winY {
		return ScreenWidth
	}
	wx := ppu.winX + 7
	switch {
	case wx < 7:
		return 0
	case wx > 166:
		return ScreenWidth
	}
	return wx - 7
}

func (ppu *PPU) drawSprites() {
	mb := ppu.gb.memoryBus
	scanline := int32(ppu.line)
	var ySize int32 = 8
	if isBitSet(ppu.lcdc, objSizeBit) {
		ySize = 16
	}

	objPalette0 := mb.readAddress(0xff48)
	objPalette1 := mb.readAddress(0xff49)

	var minX [ScreenWidth]int32
	const priorityOffset int32 = 100
	lineSprites := 0
	for sprite := range uint16(40) {
		index := 0xfe00 + (sprite * 4)

		yPos := int32(mb.readAddress(index)) - 16
		if scanline < yPos || scanline >= (yPos+ySize) {
			continue
		}

		if lineSprites >= 10 {
			break
		}
		lineSprites++

		xPos := int32(mb.readAddress(index+1)) - 8
		tileIndex := mb.readAddress(index + 2)
		attributes := mb.readAddress(index + 3)

		// in 8x16 mode the LSB of the tile index is ignored
		if ySize == 16 {
			tileIndex &= 0xfe
		}

		xFlip := isBitSet(attributes, 5)
		yFlip := isBitSet(attributes, 6)
		priority := !isBitSet(attributes, 7)

		line := scanline - yPos
		if yFlip {
			line = ySize - line - 1
		}

		dataAddress := uint16(tileIndex)*16 + uint16(line*2)
		data1 := ppu.vram[dataAddress]
		data2 := ppu.vram[dataAddress+1]

		for tilePixel := range byte(8) {
			pix := int16(xPos) + int16(7-tilePixel)
			if pix < 0 || pix >= ScreenWidth {
				continue
			}

			// this pixel is already owned by smaller or equal sprite
			if minX[pix] != 0 && minX[pix] <= xPos+priorityOffset {
				continue
			}

			colorBit := tilePixel
			if xFlip {
				colorBit = 7 - tilePixel
			}

			colorNum := (getBit(data2, colorBit)<<1 | getBit(data1, colorBit))
			// color 0 <=> transparent
			if colorNum == 0 {
				continue
			}

			var palette = objPalette0
			if isBitSet(attributes, 4) {
				palette = objPalette1
			}
			ppu.setPixel(byte(pix), byte(scanline), colorNum, palette, priority)
			minX[pix] = xPos + priorityOffset
		}
	}
}

func (ppu *PPU) drawScanline() {
	if isBitSet(ppu.lcdc, bgWindowEnableBit) {
		winStart := ppu.windowStartX()
		ppu.drawBackground(ppu.line, winStart)
		if winStart < ScreenWidth {
			ppu.drawWindow(ppu.line, winStart)
			ppu.windowLine++
		}
	}

	if isBitSet(ppu.lcdc, objEnableBit) {
		ppu.drawSprites()
	}
}

func (ppu *PPU) setEqualLinesCheck(request *interruptRequest) {
	equal := (ppu.line == ppu.lyc)
	if equal {
		setBit(&ppu.stat, coincidenceBit)
		if ppu.lyEqualsLYCInterruptEnabled() {
			request.add(lcdStatRequest)
		}
	} else {
		unsetBit(&ppu.stat, coincidenceBit)
	}
	ppu.lyEqualsLYC = equal
}

// update LYC and LY == LYC coincidence
func (ppu *PPU) writeLYC(val byte) {
	ppu.lyc = val
	equal := ppu.lyc == ppu.line
	if equal {
		setBit(&ppu.stat, coincidenceBit)
	} else {
		unsetBit(&ppu.stat, coincidenceBit)
	}
	ppu.lyEqualsLYC = equal
}

func (ppu *PPU) readVRAM(addr uint16) byte {
	return ppu.vram[addr-vramBegin]
}

func (ppu *PPU) writeVRAM(addr uint16, val byte) {
	physicalAddr := addr - vramBegin
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
