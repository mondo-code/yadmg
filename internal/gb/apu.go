package gb

import (
	"math"
)

// relevant hardware registers
const (
	WaveRamBegin = 0xff30
	WaveRamEnd = 0xff3f

	// sound channel 1
	NR10Addr = 0xff10
	NR11Addr = 0xff11
	NR12Addr = 0xff12
	NR13Addr = 0xff13
	NR14Addr = 0xff14

	// sound channel 2
	NR21Addr = 0xff16
	NR22Addr = 0xff17
	NR23Addr = 0xff18
	NR24Addr = 0xff19

	// sound channel 3
	NR30Addr = 0xff1a
	NR31Addr = 0xff1b
	NR32Addr = 0xff1c
	NR33Addr = 0xff1d
	NR34Addr = 0xff1e

	// sound channel 4
	NR41Addr = 0xff20
	NR42Addr = 0xff21
	NR43Addr = 0xff22
	NR44Addr = 0xff23

	// global control registers
	NR50Addr = 0xff24
	NR51Addr = 0xff25
	NR52Addr = 0xff26
)

// bit masks for reading audio registers 
var AudioRegisterMasks = [0x30]byte{
	0x80,	// NR10
	0x3f,	// NR11
	0x00,	// NR12
	0xff,	// NR13
	0xbf,	// NR14
	0xff,	// (ff15)
	0x3f,	// NR21
	0x00,	// NR22
	0xff,	// NR23
	0xbf,	// NR24
	0x7f,	// NR30
	0xff,	// NR31
	0x9f,	// NR32
	0xff,	// NR33
	0xbf,	// NR34
	0xff,	// (ff1f)
	0xff,	// NR41
	0x00,	// NR42
	0x00,	// NR43
	0xbf,	// NR44
	0x00,	// NR50
	0x00,	// NR51
	0x70,	// NR52
	0xff,	// (ff27)
	0xff,	// (ff28)
	0xff,	// (ff29)
	0xff,	// (ff2a)
	0xff,	// (ff2b)
	0xff,	// (ff2c)
	0xff,	// (ff2d)
	0xff,	// (ff2e)
	0xff,	// (ff2f)
}

type Sweep struct {
	enabled bool
	direction bool  // 0 == increasing, 1 == decreasing
	shift byte
	pace byte
	timer byte
	shadow uint16 
}

type Envelope struct {
	initVolume byte
	volume byte
	direction bool  // 0 == decreasing, 1 == increasing
	pace byte
}

type SoundChannel1 struct {
	enabled bool
	dacEnabled bool
	env Envelope
	sweep Sweep
	frequency uint16
	counter uint32
	length	byte
	lengthEnabled bool
	waveDuty byte
	waveDutyPos byte
	period uint16  // 11-bit period value
	timer byte
	panLeft bool
	panRight bool
}

func (sc1 *SoundChannel1) Trigger() {
	sc1.enabled = sc1.dacEnabled 
	if sc1.timer == 0 {
		sc1.timer = 64
	}
	sc1.counter = uint32(2048 - sc1.period * 4)
	sc1.env.pace = 64
	sc1.env.volume = sc1.env.initVolume
	sc1.sweep.shadow = sc1.period
	sc1.sweep.timer = 0x00

	if sc1.sweep.pace != 0 || sc1.sweep.shift != 0 {
		sc1.enabled = true
	} else {
		sc1.enabled = false
	}

	if sc1.sweep.shift != 0 {
		// calculate next period to check for overflow
		delta := sc1.sweep.shadow >> uint16(sc1.sweep.shift)
		var newPeriod uint16
		if sc1.sweep.direction {
			newPeriod = sc1.sweep.shadow - delta 
		} else {
			newPeriod = sc1.sweep.shadow + delta 
		}
		if newPeriod > 0x7ff {
			// result overflows, so disable channel and throw new period away
			sc1.enabled = false
			return
		}
		sc1.period = newPeriod
	}
}

func (sc1 *SoundChannel1) Sample() (leftOut, rightOut uint16) {
	if !sc1.enabled {
		return 0, 0
	}

	// return waveform samples

	return leftOut, rightOut
}

type SoundChannel2 struct {
	enabled bool
	dacEnabled bool
	env Envelope
	frequency uint16
	counter uint32
	length	byte
	lengthEnabled bool
	waveDuty byte
	waveDutyPos byte
	period uint16  // 11-bit period value
	panLeft bool
	panRight bool
}

func (sc2 *SoundChannel2) Trigger() {}
func (sc2 *SoundChannel2) Sample() {}

type SoundChannel3 struct {
	enabled bool
	dacEnabled bool
	length uint16	
	lengthEnabled bool
	outputLevel byte  // NR32 bits 5 and 6
	period uint16  // 11-bit period value
	timer byte
	samplePos byte
	panLeft bool
	panRight bool
}

func (sc3 *SoundChannel3) Trigger() {}
func (sc3 *SoundChannel3) Sample() {}

type SoundChannel4 struct {
	enabled bool
	dacEnabled bool
	env Envelope
	length	byte
	lengthEnabled bool
	clockDivider byte  // NR43 bits 0-2
	LFSRWidth bool  // NR43 bit 3
	clockShift byte  // NR43 bits 4-7
	lfsr uint16
	panLeft bool
	panRight bool
}

func (sc4 *SoundChannel4) Trigger() {}
func (sc4 *SoundChannel4) Sample() {}

type APU struct {
	enabled bool  // maps to NR52 register bit 7
	cycles uint16

	sc1 *SoundChannel1
	sc2 *SoundChannel2
	sc3 *SoundChannel3
	sc4 *SoundChannel4

	waveRAM	[16]byte
	volumeLeft byte
	volumeRight byte
	vinLeft bool
	vinRight bool
}

func InitAPU() *APU {
	apu := &APU{enabled: true}

	apu.sc1 = &SoundChannel1{ enabled: true, env: Envelope{}, sweep: Sweep{} }
	apu.sc2 = &SoundChannel2{ enabled: true }
	apu.sc3 = &SoundChannel3{ enabled: true }
	apu.sc4 = &SoundChannel4{ enabled: true, env: Envelope{} }

	// post-boot register values
	apu.Write(NR10Addr, 0x80)
	apu.Write(NR11Addr, 0xbf)
	apu.Write(NR12Addr, 0xf3)
	apu.Write(NR13Addr, 0xff)
	apu.Write(NR14Addr, 0xbf)
	apu.Write(NR21Addr, 0x3f)
	apu.Write(NR22Addr, 0x00)
	apu.Write(NR23Addr, 0xff)
	apu.Write(NR24Addr, 0xbf)
	apu.Write(NR30Addr, 0x7f)
	apu.Write(NR31Addr, 0xff)
	apu.Write(NR32Addr, 0x9f)
	apu.Write(NR33Addr, 0xff)
	apu.Write(NR34Addr, 0xbf)
	apu.Write(NR41Addr, 0xff)
	apu.Write(NR42Addr, 0x00)
	apu.Write(NR43Addr, 0x00)
	apu.Write(NR44Addr, 0xbf)
	apu.Write(NR50Addr, 0x77)
	apu.Write(NR51Addr, 0xf3)
	apu.Write(NR52Addr, 0xf1)

	return apu
}

func (apu *APU) powerOff() {
	var addr uint16 
	for addr = NR10Addr; addr < NR52Addr; addr++ {
		// length addresses are skipped
		if addr == NR11Addr || addr == NR21Addr || addr == NR31Addr || addr == NR41Addr {
			continue
		}
		apu.Write(addr, 0)
	}
	apu.sc1.waveDuty = 0
	apu.sc2.waveDuty = 0
}

func (apu *APU) Read(addr uint16) byte {
	var val byte
	switch {
	case addr == NR10Addr:
		val |= (apu.sc1.sweep.pace << 4)
		val |= (boolToUint8(apu.sc1.sweep.direction) << 3)
		val |= apu.sc1.sweep.shift
	case addr == NR11Addr:
		// channel 1 length timer and duty cycle
		val |= (apu.sc1.waveDuty << 6)
		// initial length timer is write-only, so we bits 0-5 should read as 1s
	case addr == NR12Addr:
		// channel 1 volume and envelope
		val = apu.sc1.env.initVolume<<4 | boolToUint8(apu.sc1.env.direction)<<3 | apu.sc1.env.pace
	// 0xff13 is write-only
	case addr == NR14Addr:
		// channel 1 period high and control
		val |= (boolToUint8(apu.sc1.lengthEnabled) << 6)
		// trigger and period are write-only, set them with unused bits
	case addr == NR21Addr:
		// channel 2 length timer and duty cycle
		val |= (apu.sc2.waveDuty << 6)
	case addr == NR22Addr:
		// channel 2 volume and envelope
		val = apu.sc2.env.initVolume<<4 | boolToUint8(apu.sc2.env.direction)<<3 | apu.sc2.env.pace
	// 0xff18 is write-only because it's the channel 2 analog to 0xff13
	case addr == NR24Addr:
		// channel 2 period high and control
		val |= (boolToUint8(apu.sc2.lengthEnabled) << 6)
	case addr == NR30Addr:
		val |= (boolToUint8(apu.sc3.dacEnabled) << 7)
	// 0xff1b is write-only
	case addr == NR32Addr:
		val = apu.sc3.outputLevel
	// 0xff1d is write-only
	case addr == NR34Addr:
		val |= boolToUint8(apu.sc3.lengthEnabled) << 6
	// 0xff20 is write-only
	case addr == NR42Addr:
		// channel 4 volume and envelope
		val = apu.sc4.env.initVolume<<4 | boolToUint8(apu.sc4.env.direction)<<3 | apu.sc4.env.pace
	case addr == NR43Addr:
		// channel 4 frequency and randomness
		val |= (apu.sc4.clockShift << 4)
		val |= (boolToUint8(apu.sc4.LFSRWidth) << 3)
		val |= apu.sc4.clockDivider
	case addr == NR44Addr:
		// channel 4 control
		val |= boolToUint8(apu.sc4.lengthEnabled) << 6
	// master control registers
	case addr == NR50Addr:
		// master volume and VIN panning
		val |= apu.volumeRight
		val |= (boolToUint8(apu.vinRight) << 3)
		val |= (apu.volumeLeft << 4)
		val |= (boolToUint8(apu.vinLeft) << 7)
	case addr == NR51Addr:
		// sound panning
		val |= boolToUint8(apu.sc1.panRight)
		val |= (boolToUint8(apu.sc1.panLeft) << 4)

		val |= (boolToUint8(apu.sc2.panRight) << 1)
		val |= (boolToUint8(apu.sc2.panLeft) << 5)

		val |= (boolToUint8(apu.sc3.panRight) << 2)
		val |= (boolToUint8(apu.sc3.panLeft) << 6)

		val |= (boolToUint8(apu.sc4.panRight) << 3)
		val |= (boolToUint8(apu.sc4.panLeft) << 7)
	case addr == NR52Addr:
		// audio master control
		val |= (boolToUint8(apu.enabled) << 7)
		val |= boolToUint8(apu.sc1.enabled)
		val |= (boolToUint8(apu.sc2.enabled) << 1)
		val |= (boolToUint8(apu.sc3.enabled) << 2)
		val |= (boolToUint8(apu.sc4.enabled) << 3)
	case addr >= WaveRamBegin && addr <= WaveRamEnd:
		return apu.waveRAM[addr-WaveRamBegin]
	default:
		return 0xff
	}
	return val | AudioRegisterMasks[addr-0xff10]
}

func (apu *APU) Write(addr uint16, val byte) {
	// channels are read-only when the APU is disabled
	if !apu.enabled && addr != NR52Addr && addr < WaveRamBegin {
		return
	}
	switch {
	case addr == NR10Addr:
		// channel 1 sweep register
		// if pace == 0, iterations are instantly disabled and will be reloaded when something else
		// is written
		pace := (val >> 4) & 0b111
		direction := val & 0b1000
		periodStep := val & 0b0111
		apu.sc1.sweep.pace = pace
		apu.sc1.sweep.direction = (direction == 0b0000_1000)
		apu.sc1.sweep.shift = periodStep
	case addr == NR11Addr:
		// channel 1 length timer, duty cycle
		initialLengthTimer := val & 0b0011_1111 
		waveDuty := (val & 0b1100_0000) >> 6
		apu.sc1.waveDuty = waveDuty
		apu.sc1.length = 64 - initialLengthTimer
	case addr == NR12Addr:
		// channel 1 volume and envelope
		// if bits 3-7 are 0, turn off the DAC
		apu.sc1.dacEnabled = ((val&0xf8) != 0)
		if !apu.sc1.dacEnabled {
			apu.sc1.enabled = false
		}
		apu.sc1.env.pace = val & 0b0000_0111
		apu.sc1.env.direction = (val & 0b0000_1000) == 0b0000_1000
		apu.sc1.env.initVolume = (val & 0b1111_0000) >> 4
	case addr == NR13Addr:
		// channel 1 period low
		apu.sc1.period = (apu.sc1.period & 0x700) | uint16(val) 
	case addr == NR14Addr:
		// channel 1 period high and control
		periodHi := uint16(val & 0b0000_0111) << 8
		apu.sc1.lengthEnabled = (val & 0b0100_0000) == 0b0100_0000
		apu.sc1.period = (apu.sc1.period & 0xff) | periodHi 
		if val & 0b1000_0000 == 0b1000_0000 {
			apu.sc1.Trigger()
		}
	case addr == NR21Addr:
		// channel 2 length timer, duty cycle 
		initialLengthTimer := val & 0b0011_1111 
		waveDuty := (val & 0b1100_0000) >> 6
		apu.sc2.waveDuty = waveDuty
		apu.sc2.length = 64 - initialLengthTimer
	case addr == NR22Addr:
		// channel 2 volume and envelope
		// if bits 3-7 are 0, turn off the DAC
		apu.sc2.dacEnabled = ((val&0xf8) != 0)
		if !apu.sc2.dacEnabled {
			apu.sc2.enabled = false
		}
		apu.sc2.env.pace = val & 0b0000_0111
		apu.sc2.env.direction = (val & 0b0000_1000) == 0b0000_1000
		apu.sc2.env.initVolume = (val & 0b1111_0000) >> 4
	case addr == NR23Addr:
		// channel 2 period low
		apu.sc2.period = (apu.sc2.period & 0x700) | uint16(val) 
	case addr == NR24Addr:
		// channel 2 period high and control
		periodHi := uint16(val & 0b0000_0111) << 8
		apu.sc2.lengthEnabled = (val & 0b0100_0000) == 0b0100_0000
		apu.sc2.period = (apu.sc2.period & 0xff) | periodHi 
		if val & 0b1000_0000 == 0b1000_0000 {
			apu.sc2.Trigger()
		}
	case addr == NR30Addr:
		// channel 3 DAC enable
		// this register can only turn off the DAC
		apu.sc3.dacEnabled = bitEnabled(val, 7)
		if !apu.sc3.dacEnabled {
			apu.sc3.enabled = false
		}
	case addr == NR31Addr:
		// channel 3 length timer
		apu.sc3.length = 256 - uint16(val)
	case addr == NR32Addr:
		// channel 3 output level
		apu.sc3.outputLevel = (val & 0b0110_0000)
	case addr == NR33Addr:
		// channel 3 period low
		apu.sc3.period = (apu.sc3.period & 0x700) | uint16(val)
	case addr == NR34Addr:
		// channel 3 period high and control
		apu.sc3.lengthEnabled = bitEnabled(val, 6)
		periodHi := uint16(val & 0x7) << 8
		apu.sc3.period = (apu.sc3.period & 0xff) | periodHi
		if bitEnabled(val, 7) {
			apu.sc3.Trigger()
		}
	case addr == NR41Addr:
		// channel 4 length timer
		apu.sc4.length = 64 - (val & 0x3f)
	case addr == NR42Addr:
		// channel 4 volume and envelope
		// if bits 3-7 are 0, turn off the DAC
		apu.sc4.dacEnabled = ((val&0xf8) != 0)
		if !apu.sc4.dacEnabled {
			apu.sc4.enabled = false
		}
		apu.sc4.env.pace = val & 0b0000_0111
		apu.sc4.env.direction = bitEnabled(val, 3) 
		apu.sc4.env.initVolume = (val & 0b1111_0000) >> 4
	case addr == NR43Addr:
		// channel 4 frequency and randomness
		apu.sc4.clockShift = val >> 4 
		apu.sc4.LFSRWidth = bitEnabled(val, 3)
		apu.sc4.clockDivider = val & 0x07
	case addr == NR44Addr:
		// channel 4 control
		apu.sc4.lengthEnabled = bitEnabled(val, 6)
		if bitEnabled(val, 7) {
			apu.sc4.Trigger()
		}
	case addr == NR50Addr:
		// master volume and VIN panning
		rightVol := val & 0x07
		leftVol := val & 0x70
		apu.vinLeft = bitEnabled(val, 7)
		apu.vinRight = bitEnabled(val, 3)
		apu.volumeLeft = (leftVol >> 4)
		apu.volumeRight = rightVol
	case addr == NR51Addr:
		// sound panning
		apu.sc1.panRight = bitEnabled(val, 0)
		apu.sc1.panLeft = bitEnabled(val, 4)

		apu.sc2.panRight = bitEnabled(val, 1)
		apu.sc2.panLeft = bitEnabled(val, 5)

		apu.sc3.panRight = bitEnabled(val, 2)
		apu.sc3.panLeft = bitEnabled(val, 6)

		apu.sc4.panRight = bitEnabled(val, 3)
		apu.sc4.panLeft = bitEnabled(val, 7)
	case addr == NR52Addr:
		// audio master control
		on := bitEnabled(val, 7)
		if !on && apu.enabled {
			apu.powerOff()
		}
		apu.enabled = on
	case addr >= WaveRamBegin && addr <= WaveRamEnd:
		apu.waveRAM[addr-WaveRamBegin] = val
	}
}

func (apu *APU) Step(cycles uint16) {
	apu.cycles += cycles

	apu.sc1.Sample()
	apu.sc2.Sample()
	apu.sc3.Sample()
	apu.sc4.Sample()

}

func calcFrequency(clockDivider, clockShift byte) float64 {
	// leaving this detail out for now: shift being 14 or 15 stops the channel from being clocked entirely
	divider := float64(clockDivider)
	shift := float64(clockShift)
	return (262144 / (divider * math.Pow(2, shift)))
}
