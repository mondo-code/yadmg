package gb

// relevant hardware registers
const (
	WaveRamBegin = 0xff30
	WaveRamEnd   = 0xff3f

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
	0x80, // NR10
	0x3f, // NR11
	0x00, // NR12
	0xff, // NR13
	0xbf, // NR14
	0xff, // (ff15)
	0x3f, // NR21
	0x00, // NR22
	0xff, // NR23
	0xbf, // NR24
	0x7f, // NR30
	0xff, // NR31
	0x9f, // NR32
	0xff, // NR33
	0xbf, // NR34
	0xff, // (ff1f)
	0xff, // NR41
	0x00, // NR42
	0x00, // NR43
	0xbf, // NR44
	0x00, // NR50
	0x00, // NR51
	0x70, // NR52
	0xff, // (ff27)
	0xff, // (ff28)
	0xff, // (ff29)
	0xff, // (ff2a)
	0xff, // (ff2b)
	0xff, // (ff2c)
	0xff, // (ff2d)
	0xff, // (ff2e)
	0xff, // (ff2f)
}

const (
	SC1LengthMax uint16 = 64
	SC2LengthMax uint16 = 64
	SC3LengthMax uint16 = 256
	SC4LengthMax uint16 = 64

	FrameSequencerPeriod = 8192
)

type Sweep struct {
	enabled  bool
	decrease bool
	shift    byte
	pace     byte
	timer    byte
	shadow   uint16
}

func (s *Sweep) clock() bool {
	if s.enabled {
		s.timer--
		if s.timer == 0 {
			s.timer = 8
			return true
		}
	}
	return false
}

func (s *Sweep) calculate() uint16 {
	change := s.shadow >> uint16(s.shift)
	var freq uint16
	if s.decrease {
		freq = s.shadow - change
	} else {
		freq = s.shadow + change
	}
	return freq
}

type Envelope struct {
	enabled    bool
	initVolume byte
	volume     byte
	decrease   bool
	timer      byte
	pace       byte
}

func (e *Envelope) clock() {
	if e.enabled && e.timer == 0 {
		e.timer--
		var volume byte
		if e.decrease {
			volume = e.volume - 1
		} else {
			volume = e.volume + 1
		}
		if 0 <= volume && volume <= 15 {
			e.volume = volume
		} else {
			e.enabled = false
		}
	}
}

type LengthTimer struct {
	enabled bool // NRx4 bit 6
	counter uint16
}

func (l *LengthTimer) clock() bool {
	if !l.enabled || l.counter == 0 {
		return false
	}
	l.counter--
	return l.counter == 0
}

type SoundChannel1 struct {
	enabled     bool
	dacEnabled  bool
	env         Envelope
	sweep       Sweep
	counter     uint32
	length      LengthTimer
	waveDuty    byte
	waveDutyPos byte
	period      uint16 // 11-bit period value
	panLeft     bool
	panRight    bool
}

func (sc1 *SoundChannel1) Trigger() {
	sc1.enabled = sc1.dacEnabled
	if sc1.length.counter == 0 {
		sc1.length.counter = SC1LengthMax
	}
	sc1.counter = uint32(2048-sc1.period) * 4
	sc1.env.volume = sc1.env.initVolume
	sc1.env.timer = sc1.env.pace
	if sc1.env.timer == 0 {
		sc1.env.timer = 8
	}
	sc1.sweep.shadow = sc1.period
	sc1.sweep.timer = sc1.sweep.pace
	if sc1.sweep.timer == 0 {
		sc1.sweep.timer = 8
	}

	if sc1.sweep.pace != 0 || sc1.sweep.shift != 0 {
		sc1.sweep.enabled = true
	} else {
		sc1.sweep.enabled = false
	}

	if sc1.sweep.shift != 0 {
		sc1.clockSweep()
	}
}

func (sc1 *SoundChannel1) clockLength() {
	if sc1.length.clock() {
		sc1.length.counter = SC1LengthMax
		sc1.enabled = false
	}
}

func (sc1 *SoundChannel1) clockSweep() {
	freq := sc1.sweep.calculate()
	if freq > 0x7ff {
		// result overflows, so disable channel and throw new period away
		sc1.enabled = false
	}
}

func (sc1 *SoundChannel1) Sample() (leftOut, rightOut uint16) {
	if !sc1.enabled {
		return 0, 0
	}

	// return waveform samples

	return leftOut, rightOut
}

func (sc1 *SoundChannel1) Step(cycles uint16) {}

type SoundChannel2 struct {
	enabled     bool
	dacEnabled  bool
	env         Envelope
	counter     uint32
	length      LengthTimer
	waveDuty    byte
	waveDutyPos byte
	period      uint16 // 11-bit period value
	panLeft     bool
	panRight    bool
}

func (sc2 *SoundChannel2) Trigger() {
	sc2.enabled = sc2.dacEnabled
	if sc2.length.counter == 0 {
		sc2.length.counter = SC2LengthMax
	}
	sc2.counter = uint32(2048-sc2.period) * 4
	sc2.env.timer = sc2.env.pace
	if sc2.env.timer == 0 {
		sc2.env.timer = 8
	}
	sc2.env.volume = sc2.env.initVolume
}

func (sc2 *SoundChannel2) clockLength() {
	if sc2.length.clock() {
		sc2.length.counter = SC2LengthMax
		sc2.enabled = false
	}
}

func (sc2 *SoundChannel2) Sample()            {}
func (sc2 *SoundChannel2) Step(cycles uint16) {}

type SoundChannel3 struct {
	enabled     bool
	dacEnabled  bool
	length      LengthTimer
	outputLevel byte   // NR32 bits 5 and 6
	period      uint16 // 11-bit period value
	timer       uint16
	samplePos   byte
	panLeft     bool
	panRight    bool
}

func (sc3 *SoundChannel3) Trigger() {
	sc3.enabled = sc3.dacEnabled
	if sc3.length.counter == 0 {
		sc3.length.counter = SC3LengthMax
	}
	sc3.samplePos = 0
	sc3.timer = (2048 - sc3.period) * 2
}

func (sc3 *SoundChannel3) clockLength() {
	if sc3.length.clock() {
		sc3.length.counter = SC3LengthMax
		sc3.enabled = false
	}
}

func (sc3 *SoundChannel3) Sample()            {}
func (sc3 *SoundChannel3) Step(cycles uint16) {}

type SoundChannel4 struct {
	enabled      bool
	dacEnabled   bool
	env          Envelope
	length       LengthTimer
	clockDivider byte // NR43 bits 0-2
	LFSRWidth    bool // NR43 bit 3
	clockShift   byte // NR43 bits 4-7
	lfsr         uint16
	timer        uint32
	panLeft      bool
	panRight     bool
}

var noiseDivisors = [8]uint32{8, 16, 32, 48, 64, 80, 96, 112}

func (sc4 *SoundChannel4) reloadTimer() {
	sc4.timer = noiseDivisors[sc4.clockDivider] << sc4.clockShift
}

func (sc4 *SoundChannel4) Trigger() {
	sc4.enabled = sc4.dacEnabled
	if sc4.length.counter == 0 {
		sc4.length.counter = SC4LengthMax
	}
	sc4.env.volume = sc4.env.initVolume
	sc4.env.timer = sc4.env.pace
	if sc4.env.timer == 0 {
		sc4.env.timer = 8
	}
	sc4.reloadTimer()
	sc4.lfsr = 0x7fff
}

func (sc4 *SoundChannel4) clockLength() {
	if sc4.length.clock() {
		sc4.length.counter = SC4LengthMax
		sc4.enabled = false
	}
}

func (sc4 *SoundChannel4) Sample()            {}
func (sc4 *SoundChannel4) Step(cycles uint16) {}

type APU struct {
	enabled  bool   // maps to NR52 register bit 7
	fsCycles uint16 // frame sequencer cycle count
	fsPhase  byte   // next phase of the frame sequencer (0-7)

	sc1 *SoundChannel1
	sc2 *SoundChannel2
	sc3 *SoundChannel3
	sc4 *SoundChannel4

	waveRAM     [16]byte
	volumeLeft  byte
	volumeRight byte
	vinLeft     bool
	vinRight    bool
}

func InitAPU() *APU {
	apu := &APU{enabled: true}

	apu.sc1 = &SoundChannel1{enabled: true, env: Envelope{}, sweep: Sweep{}}
	apu.sc2 = &SoundChannel2{enabled: true}
	apu.sc3 = &SoundChannel3{enabled: true}
	apu.sc4 = &SoundChannel4{enabled: true, env: Envelope{}}

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
		val |= (boolToUint8(apu.sc1.sweep.decrease) << 3)
		val |= apu.sc1.sweep.shift
	case addr == NR11Addr:
		// channel 1 length timer and duty cycle
		val |= (apu.sc1.waveDuty << 6)
		// initial length timer is write-only, so we bits 0-5 should read as 1s
	case addr == NR12Addr:
		// channel 1 volume and envelope
		val = apu.sc1.env.initVolume<<4 | boolToUint8(apu.sc1.env.decrease)<<3 | apu.sc1.env.pace
	// 0xff13 is write-only
	case addr == NR14Addr:
		// channel 1 period high and control
		val |= (boolToUint8(apu.sc1.length.enabled) << 6)
		// trigger and period are write-only, set them with unused bits
	case addr == NR21Addr:
		// channel 2 length timer and duty cycle
		val |= (apu.sc2.waveDuty << 6)
	case addr == NR22Addr:
		// channel 2 volume and envelope
		val = apu.sc2.env.initVolume<<4 | boolToUint8(apu.sc2.env.decrease)<<3 | apu.sc2.env.pace
	// 0xff18 is write-only because it's the channel 2 analog to 0xff13
	case addr == NR24Addr:
		// channel 2 period high and control
		val |= (boolToUint8(apu.sc2.length.enabled) << 6)
	case addr == NR30Addr:
		val |= (boolToUint8(apu.sc3.dacEnabled) << 7)
	// 0xff1b is write-only
	case addr == NR32Addr:
		val = apu.sc3.outputLevel
	// 0xff1d is write-only
	case addr == NR34Addr:
		val |= boolToUint8(apu.sc3.length.enabled) << 6
	// 0xff20 is write-only
	case addr == NR42Addr:
		// channel 4 volume and envelope
		val = apu.sc4.env.initVolume<<4 | boolToUint8(apu.sc4.env.decrease)<<3 | apu.sc4.env.pace
	case addr == NR43Addr:
		// channel 4 frequency and randomness
		val |= (apu.sc4.clockShift << 4)
		val |= (boolToUint8(apu.sc4.LFSRWidth) << 3)
		val |= apu.sc4.clockDivider
	case addr == NR44Addr:
		// channel 4 control
		val |= boolToUint8(apu.sc4.length.enabled) << 6
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
		sweepPace := (val >> 4) & 0b111
		sweepDir := val & 0b1000
		sweepShift := val & 0b0111
		apu.sc1.sweep.pace = sweepPace
		apu.sc1.sweep.decrease = (sweepDir == 0b0000_1000)
		apu.sc1.sweep.shift = sweepShift
	case addr == NR11Addr:
		// channel 1 length timer, duty cycle
		initialLengthTimer := val & 0b0011_1111
		waveDuty := (val & 0b1100_0000) >> 6
		apu.sc1.waveDuty = waveDuty
		apu.sc1.length.counter = SC1LengthMax - uint16(initialLengthTimer)
	case addr == NR12Addr:
		// channel 1 volume and envelope
		// if bits 3-7 are 0, turn off the DAC
		apu.sc1.dacEnabled = ((val & 0xf8) != 0)
		if !apu.sc1.dacEnabled {
			apu.sc1.enabled = false
		}
		apu.sc1.env.pace = val & 0b0000_0111
		apu.sc1.env.decrease = (val & 0b0000_1000) == 0b0000_1000
		apu.sc1.env.initVolume = (val & 0b1111_0000) >> 4
	case addr == NR13Addr:
		// channel 1 period low
		apu.sc1.period = (apu.sc1.period & 0x700) | uint16(val)
	case addr == NR14Addr:
		// channel 1 period high and control
		periodHi := uint16(val&0b0000_0111) << 8
		apu.sc1.period = (apu.sc1.period & 0xff) | periodHi

		sc1Length := &apu.sc1.length
		wasLengthEnabled := apu.sc1.length.enabled
		sc1Length.enabled = bitEnabled(val, 6)
		isLengthPhase := (apu.fsPhase&1 == 0)
		doTrigger := bitEnabled(val, 7)

		// enabling in the first half of length period clocks length
		if !isLengthPhase && !wasLengthEnabled && sc1Length.enabled && sc1Length.counter > 0 {
			sc1Length.counter--
			if sc1Length.counter == 0 && !doTrigger {
				apu.sc1.enabled = false
			}
		}

		if doTrigger {
			reloaded := (sc1Length.counter == 0)
			apu.sc1.Trigger()
			if reloaded && sc1Length.enabled && !isLengthPhase {
				sc1Length.counter = SC1LengthMax - 1
			}
		}
	case addr == NR21Addr:
		// channel 2 length timer, duty cycle
		initialLengthTimer := val & 0b0011_1111
		waveDuty := (val & 0b1100_0000) >> 6
		apu.sc2.waveDuty = waveDuty
		apu.sc2.length.counter = SC2LengthMax - uint16(initialLengthTimer)
	case addr == NR22Addr:
		// channel 2 volume and envelope
		// if bits 3-7 are 0, turn off the DAC
		apu.sc2.dacEnabled = ((val & 0xf8) != 0)
		if !apu.sc2.dacEnabled {
			apu.sc2.enabled = false
		}
		apu.sc2.env.pace = val & 0b0000_0111
		apu.sc2.env.decrease = (val & 0b0000_1000) == 0b0000_1000
		apu.sc2.env.initVolume = (val & 0b1111_0000) >> 4
	case addr == NR23Addr:
		// channel 2 period low
		apu.sc2.period = (apu.sc2.period & 0x700) | uint16(val)
	case addr == NR24Addr:
		// channel 2 period high and control
		periodHi := uint16(val&0b0000_0111) << 8
		apu.sc2.period = (apu.sc2.period & 0xff) | periodHi

		sc2Length := &apu.sc2.length
		wasLengthEnabled := apu.sc2.length.enabled
		sc2Length.enabled = bitEnabled(val, 6)
		isLengthPhase := (apu.fsPhase&1 == 0)
		doTrigger := bitEnabled(val, 7)

		if !isLengthPhase && !wasLengthEnabled && sc2Length.enabled && sc2Length.counter > 0 {
			sc2Length.counter--
			if sc2Length.counter == 0 && !doTrigger {
				apu.sc2.enabled = false
			}
		}

		if doTrigger {
			reloaded := (sc2Length.counter == 0)
			apu.sc2.Trigger()
			if reloaded && sc2Length.enabled && !isLengthPhase {
				sc2Length.counter = SC2LengthMax - 1
			}
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
		apu.sc3.length.counter = SC3LengthMax - uint16(val)
	case addr == NR32Addr:
		// channel 3 output level
		apu.sc3.outputLevel = (val & 0b0110_0000)
	case addr == NR33Addr:
		// channel 3 period low
		apu.sc3.period = (apu.sc3.period & 0x700) | uint16(val)
	case addr == NR34Addr:
		// channel 3 period high and control
		periodHi := uint16(val&0b0000_0111) << 8
		apu.sc3.period = (apu.sc3.period & 0xff) | periodHi

		sc3Length := &apu.sc3.length
		wasLengthEnabled := apu.sc3.length.enabled
		sc3Length.enabled = bitEnabled(val, 6)
		isLengthPhase := (apu.fsPhase&1 == 0)
		doTrigger := bitEnabled(val, 7)

		if !isLengthPhase && !wasLengthEnabled && sc3Length.enabled && sc3Length.counter > 0 {
			sc3Length.counter--
			if sc3Length.counter == 0 && !doTrigger {
				apu.sc3.enabled = false
			}
		}

		if doTrigger {
			reloaded := (sc3Length.counter == 0)
			apu.sc3.Trigger()
			if reloaded && sc3Length.enabled && !isLengthPhase {
				sc3Length.counter = SC3LengthMax - 1
			}
		}
	case addr == NR41Addr:
		// channel 4 length timer
		apu.sc4.length.counter = SC4LengthMax - uint16((val & 0x3f))
	case addr == NR42Addr:
		// channel 4 volume and envelope
		// if bits 3-7 are 0, turn off the DAC
		apu.sc4.dacEnabled = ((val & 0xf8) != 0)
		if !apu.sc4.dacEnabled {
			apu.sc4.enabled = false
		}
		apu.sc4.env.pace = val & 0b0000_0111
		apu.sc4.env.decrease = bitEnabled(val, 3)
		apu.sc4.env.initVolume = (val & 0b1111_0000) >> 4
	case addr == NR43Addr:
		// channel 4 frequency and randomness
		apu.sc4.clockShift = val >> 4
		apu.sc4.LFSRWidth = bitEnabled(val, 3)
		apu.sc4.clockDivider = val & 0x07
	case addr == NR44Addr:
		// channel 4 control
		sc4Length := &apu.sc4.length
		isLengthPhase := (apu.fsPhase&1 == 0)
		wasLengthEnabled := apu.sc4.length.enabled
		sc4Length.enabled = bitEnabled(val, 6)
		doTrigger := bitEnabled(val, 7)

		if !isLengthPhase && !wasLengthEnabled && sc4Length.enabled && sc4Length.counter > 0 {
			sc4Length.counter--
			if sc4Length.counter == 0 && !doTrigger {
				apu.sc4.enabled = false
			}
		}

		if doTrigger {
			reloaded := (sc4Length.counter == 0)
			apu.sc4.Trigger()
			if reloaded && sc4Length.enabled && !isLengthPhase {
				sc4Length.counter = SC4LengthMax - 1
			}
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

func (apu *APU) clockChannelLengths() {
	apu.sc1.clockLength()
	apu.sc2.clockLength()
	apu.sc3.clockLength()
	apu.sc4.clockLength()
}

func (apu *APU) Step(cycles uint16) {
	if !apu.enabled {
		return
	}

	// frame sequencer stuff
	apu.fsCycles += cycles
	for apu.fsCycles >= FrameSequencerPeriod {
		apu.fsCycles -= FrameSequencerPeriod
		switch apu.fsPhase {
		case 0, 4:
			apu.clockChannelLengths()
		case 2, 6:
			apu.clockChannelLengths()
			apu.sc1.clockSweep()
		case 7:
			apu.sc1.env.clock()
			apu.sc2.env.clock()
			apu.sc4.env.clock()
		}
		apu.fsPhase = (apu.fsPhase + 1) & 7
	}

	apu.sc1.Step(cycles)
	apu.sc2.Step(cycles)
	apu.sc3.Step(cycles)
	apu.sc4.Step(cycles)
}
