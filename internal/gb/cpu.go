package gb 

import (
	"fmt"
)

type Registers struct {
	A byte
	B byte
	C byte
	D byte
	E byte
	H byte
	L byte
}

type CPU struct {
	gb 				*Gameboy
	regs 			Registers
	flags   		FlagsRegister
	PC      		uint16
	SP 				uint16
	instructions 	[0x100]Instruction
	cbInstructions 	[0x100]Instruction
	Halted 			bool
	IME 			bool
	IMEPending		bool
}

func (cpu *CPU) getBC() uint16 {
	return bytesToWord(cpu.regs.B, cpu.regs.C)
}

func (cpu *CPU) setBC(val uint16) {
	cpu.regs.B, cpu.regs.C = wordToBytes(val)
}

func (cpu *CPU) getDE() uint16 {
	return bytesToWord(cpu.regs.D, cpu.regs.E)
}

func (cpu *CPU) setDE(val uint16) {
	cpu.regs.D, cpu.regs.E = wordToBytes(val)
}

func (cpu *CPU) getHL() uint16 {
	return bytesToWord(cpu.regs.H, cpu.regs.L)
}

func (cpu *CPU) setHL(val uint16) {
	cpu.regs.H, cpu.regs.L = wordToBytes(val)
}

func (cpu *CPU) getAF() uint16 {
	return bytesToWord(cpu.regs.A, cpu.flags.convertToByte())
}

func (cpu *CPU) setAF(val uint16) {
	aByte, fByte := wordToBytes(val)
	cpu.regs.A = aByte
	cpu.flags.convertFromByte(fByte)
}

func (cpu *CPU) getSP() uint16 {
	return cpu.SP
}

func (cpu *CPU) setSP(val uint16) {
	cpu.SP = val
}

func (cpu *CPU) Step() (cycles uint16, err error) {
	cycles += cpu.gb.HandleInterrupts()

	var opCycles uint16 = 4
	var nextPC uint16 = cpu.PC
	if !cpu.Halted {
		opcode := cpu.gb.MemoryBus.ReadAddress(cpu.PC)
		instruction := cpu.instructions[opcode]

		if instruction == nil {
			return 0, fmt.Errorf("nil opcode: 0x%02X", opcode)
		}

		operands := []byte{
			cpu.gb.MemoryBus.ReadAddress(cpu.PC+1),
			cpu.gb.MemoryBus.ReadAddress(cpu.PC+2),
		}

		if cpu.gb.LogOpcodes {
			cpu.LogInstruction()
		}
		nextPC, opCycles = instruction(cpu, operands)
	}

	cycles += opCycles
	cpu.PC = nextPC
	cpu.gb.StepDivider(opCycles)
	cpu.gb.StepTimer(opCycles)

	return cycles, nil
}

// this is representable as a byte but DMG-01 justifies this as being less error prone 
// would like to refactor this into a byte later for performance
type FlagsRegister struct {
	Zero bool
	Subtract bool
	HalfCarry bool
	Carry bool
}

const (
	ZeroFlagLocation 		= 7
	SubtractFlagLocation 	= 6
	HalfCarryFlagLoction 	= 5
	CarryFlagLocation 		= 4
)

func (fr *FlagsRegister) convertToByte() byte {
	return byte((boolToUint8(fr.Zero) << ZeroFlagLocation) | 
	(boolToUint8(fr.Subtract) << SubtractFlagLocation) | 
	(boolToUint8(fr.HalfCarry) << HalfCarryFlagLoction) | 
	(boolToUint8(fr.Carry) << CarryFlagLocation))
}

func (fr *FlagsRegister) convertFromByte(val byte) {
	fr.Zero = ((val >> ZeroFlagLocation) & 0b1) != 0
	fr.Subtract = ((val >> SubtractFlagLocation) & 0b1) != 0
	fr.HalfCarry = ((val >> HalfCarryFlagLoction) & 0b1) != 0
	fr.Carry = ((val >> CarryFlagLocation) & 0b1) != 0
}

func InitCPU(gb *Gameboy) *CPU {
	cpu := &CPU{
		gb: gb,
		// initial values for registers
		regs: Registers{
			A: 0x01, 
			B: 0x00,
			C: 0x13,
			D: 0x00,
			E: 0xd8,
			H: 0x01,
			L: 0x4d,
		},
		flags: FlagsRegister{
			Zero: true,
			Subtract: false,
			HalfCarry: true,
			Carry: true,
		},
		PC: 0x0100,  // initialized at cartridge entry point
		SP: 0xfffe,
		Halted: false,
		IME: false,
	}

	cpu.initInstructions()
	cpu.initInstructionsCB()

	return cpu
}

func (cpu *CPU) nop() uint16 {
	return cpu.PC + 1
}

// jumps
func (cpu *CPU) jrCond(cond bool, offset int8) uint16 {
	next := cpu.PC + 2
	if cond {
		return uint16(int32(next) + int32(offset))
	}
	return next
}

func (cpu *CPU) jpCond(cond bool, target uint16) uint16 {
	if cond {
		return target
	}
	return cpu.PC + 3
}

func (cpu *CPU) call(doJump bool, target uint16) uint16 {
	nextPC := cpu.PC + 3
	if doJump {
		cpu.push(nextPC)
		return target
	}
	return nextPC
}

func (cpu *CPU) ret(doJump bool) uint16 {
	if doJump {
		return cpu.pop()
	}
	return cpu.PC + 1
}

func (cpu *CPU) push(val uint16) {
	cpu.SP--
	cpu.gb.MemoryBus.WriteToAddress(cpu.SP, getUpperByte(val))
	cpu.SP--
	cpu.gb.MemoryBus.WriteToAddress(cpu.SP, getLowerByte(val))
}

func (cpu *CPU) pop() (res uint16) {
	sp := cpu.getSP()
	byte1 := cpu.gb.MemoryBus.ReadAddress(sp)
	res = setLowerByte(res, byte1)
	byte2 := cpu.gb.MemoryBus.ReadAddress(sp + 1)
	res = setUpperByte(res, byte2)
	cpu.setSP(sp + 2)
	return res 
}

func (cpu *CPU) rst(addr uint16) uint16 {
	cpu.SP--
	cpu.gb.MemoryBus.WriteToAddress(cpu.SP, getUpperByte(cpu.PC + 1))
	cpu.SP--
	cpu.gb.MemoryBus.WriteToAddress(cpu.SP, getLowerByte(cpu.PC + 1))
	return addr
}

func (cpu *CPU) add(val byte, addCarry bool) {
	var carryIn byte
	if addCarry && cpu.flags.Carry {
		carryIn = 1
	}

	sum, sumCarry := sumBytes(val, cpu.regs.A)
	total, doCarry := sumBytes(sum, carryIn)

	cpu.flags.Subtract = false
	cpu.flags.Carry = sumCarry || doCarry 
	cpu.flags.HalfCarry = (cpu.regs.A & 0xf) + (val & 0xf) + carryIn > 0xf   

	cpu.regs.A = total 
	cpu.flags.Zero = (cpu.regs.A == 0)
}

func (cpu *CPU) add_hl(val uint16) {
	hlVal := bytesToWord(cpu.regs.H, cpu.regs.L) 
	sum, overflow := sumUint16(val, uint16(hlVal))
	cpu.flags.Subtract = false
	cpu.flags.Carry = overflow
	cpu.flags.HalfCarry = (hlVal & 0x0fff) + (val & 0x0fff) > 0x0fff   
	cpu.regs.H, cpu.regs.L = wordToBytes(sum)
}

// subtract the value in target register from the value in the A register
func (cpu *CPU) sub(val byte, addCarry bool) {
	var carryIn byte
	if addCarry && cpu.flags.Carry {
		carryIn = 1
	}

	diff, underflow := subBytes(cpu.regs.A, val)
	res, doUnderflow := subBytes(diff, carryIn)

	cpu.flags.Subtract = true 
	cpu.flags.Carry = underflow || doUnderflow
	cpu.flags.HalfCarry = (cpu.regs.A & 0xf) < (val & 0xf) + carryIn 

	cpu.regs.A = res 
	cpu.flags.Zero = (cpu.regs.A == 0)
}

// bitwise ops
func (cpu *CPU) and(val byte) {
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = true 
	cpu.flags.Carry = false
	cpu.regs.A &= val 
	cpu.flags.Zero = (cpu.regs.A == 0)
}

func (cpu *CPU) or(val byte) {
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false
	cpu.flags.Carry = false
	cpu.regs.A |= val 
	cpu.flags.Zero = (cpu.regs.A == 0)
}

func (cpu *CPU) xor(val byte) {
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false
	cpu.flags.Carry = false
	cpu.regs.A ^= val 
	cpu.flags.Zero = (cpu.regs.A == 0)
}

// compare A with val 
func (cpu *CPU) cp_n(val byte) {
	cpu.flags.Subtract = true
	if (cpu.regs.A & 0xf) < (val & 0xf) {
		cpu.flags.HalfCarry = true
	} else {
		cpu.flags.HalfCarry = false 
	}

	if cpu.regs.A < val {
		cpu.flags.Carry = true
	} else {
		cpu.flags.Carry = false
	}

	cpu.flags.Zero = (cpu.regs.A - val) == 0
}

// increment 8 bit target
func (cpu *CPU) inc_n(target *byte) {
	cpu.flags.Subtract = false 
	cpu.flags.HalfCarry = (*target & 0x0f) == 0x0f
	*target += 1
	cpu.flags.Zero = (*target == 0)
}

// increment 16 bit target
func (cpu *CPU) inc_nn(getReg func() uint16, setReg func(uint16)) {
	val := getReg()
	setReg(val + 1)
}

// decrement 8 bit target
func (cpu *CPU) dec_n(target *byte) {
	cpu.flags.Subtract = true 
	cpu.flags.HalfCarry = (*target & 0x0f) == 0x00
	*target -= 1
	cpu.flags.Zero = (*target == 0)
}

// decrement 16 bit target 
func (cpu *CPU) dec_nn(getReg func() uint16, setReg func(uint16)) {
	val := getReg()
	setReg(val - 1)
}

// toggle carry flag
func (cpu *CPU) ccf() {
	cpu.flags.Carry = !cpu.flags.Carry
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false
}

func (cpu *CPU) scf() {
	cpu.flags.Carry = true
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false
}

func (cpu *CPU) rra() {
	newCarry := false
	regVal := cpu.regs.A

	// if the rightmost bit is a 1, it will "fall off" into the carry flag on rotate
	if regVal & 0x01 == 0x01 {
		newCarry = true 
	} 
	regVal >>= 1 

	// leftmost bit needs to be 1 if carry is set
	if cpu.flags.Carry {
		regVal |= 0x80
	}

	cpu.flags.Carry = newCarry
	cpu.regs.A = regVal
	cpu.flags.Zero = false
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false
}

func (cpu *CPU) rla() {
	newCarry := false
	regVal := cpu.regs.A

	if regVal & 0x80 == 0x80 {
		newCarry = true
	}
	regVal <<= 1

	if cpu.flags.Carry {
		regVal |= 0x01
	}

	cpu.flags.Carry = newCarry
	cpu.regs.A = regVal
	cpu.flags.Zero = false
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false
}

// right rotate A register, not through carry flag
func (cpu *CPU) rrca() {
	regVal := cpu.regs.A
	bit0 := regVal & 0x01 
	regVal >>= 1
	
	if bit0 != 0 {
		regVal |= 0x80
	}

	cpu.regs.A = regVal
	cpu.flags.Carry = (bit0 != 0)
	cpu.flags.Zero = false
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false
}

func (cpu *CPU) rlca() {
	regVal := cpu.regs.A
	bit7 := regVal & 0x80 
	regVal <<= 1
	
	if bit7 != 0 {
		regVal |= 0x01 
	}

	cpu.regs.A = regVal
	cpu.flags.Carry = (bit7 != 0)
	cpu.flags.Zero = false
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false
}

func (cpu *CPU) cpl() {
	 cpu.regs.A = ^cpu.regs.A
}

// test bit in target register
func (cpu *CPU) bit(bit byte, target *byte) {
	cpu.flags.Zero = ((*target & (1 << bit)) == 0)
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = true 
}

// reset particular bit of a register to 0 
func (cpu *CPU) reset(bit byte, target *byte) {
	*target &= ^(1 << bit)
}

func (cpu *CPU) set(bit byte, target *byte) {
	*target |= (1 << bit)
}

func (cpu *CPU) srl(target *byte) {
	cpu.flags.Carry = ((*target & 1) == 1)
	*target >>= 1
	cpu.flags.Zero = (*target == 0)
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false
}

// rotate bit right through carry flag
func (cpu *CPU) rr(val *byte) {
	rotated := *val >> 1
	if cpu.flags.Carry {
		rotated |= 0x80
	}

	cpu.flags.Carry = ((*val & 1) == 1)
	cpu.flags.Zero = (rotated == 0)
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false
	*val = rotated
}

func (cpu *CPU) rl(val *byte) {
	rotated := *val << 1
	if cpu.flags.Carry {
		rotated |= 0x01 
	}

	cpu.flags.Carry = ((*val >> 7) == 0x01)
	cpu.flags.Zero = (rotated == 0)
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false
	*val = rotated
}

// rotate, not through carry flag
func (cpu *CPU) rrc(val *byte) {
	bit0 := *val & 0x01
	rotated := *val >> 1
	if bit0 != 0 {
		rotated |= 0x80 
	}

	cpu.flags.Carry = (bit0 != 0)
	cpu.flags.Zero = (rotated == 0)
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false
	*val = rotated
}

func (cpu *CPU) rlc(val *byte) {
	bit7 := *val & 0x80 
	rotated := *val << 1
	if bit7 != 0 {
		rotated |= 0x01 
	}

	cpu.flags.Carry = (bit7 != 0)
	cpu.flags.Zero = (rotated == 0)
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false
	*val = rotated
}

// arithmetic shift preserves the most significant bit because it was made to compensate for signed numbers
// only relevant for right shift
func (cpu *CPU) sra(val *byte) {
	signBit := *val >> 7
	rotated := *val >> 1 
	if signBit != 0 {
		rotated |= 0x80 
	}

	cpu.flags.Zero = (rotated == 0)
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false

	if (*val & 0x01) == 0x01{
		cpu.flags.Carry = true
	} else {
		cpu.flags.Carry = false
	}

	*val = rotated
}

func (cpu *CPU) sla(val *byte) {
	rotated := *val << 1 
	cpu.flags.Zero = (rotated == 0)
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false

	if (*val & 0x80) == 0x80 {
		cpu.flags.Carry = true
	} else {
		cpu.flags.Carry = false
	}

	*val = rotated
}

// swap nibbles in given byte
func (cpu *CPU) swap(val *byte) {
	upperNibble := (*val & 0xf0) >> 4
	lowerNibble := (*val & 0x0f) << 4
	swapped := upperNibble ^ lowerNibble

	cpu.flags.Zero = (swapped == 0)
	cpu.flags.Subtract = false
	cpu.flags.HalfCarry = false
	cpu.flags.Carry = false

	*val = swapped
}

func (cpu *CPU) inc_sp() {
	cpu.SP = (cpu.SP + 1) & 0xffff
}

func (cpu *CPU) dec_sp() {
	cpu.SP = (cpu.SP - 1) & 0xffff
}

func (cpu *CPU) getDAA() (byte, bool) {
	var offset byte = 0
	setCarry := false
	aVal := cpu.regs.A
	halfCarry := cpu.flags.HalfCarry 
	carry := cpu.flags.Carry
	subtract := cpu.flags.Subtract

	if (!subtract && aVal & 0xf > 0x09) || halfCarry {
		offset |= 0x06
	}

	if (!subtract && aVal > 0x99) || carry {
		offset |= 0x60
		setCarry = true
	}

	if subtract {
		aVal -= offset
	} else {
		aVal += offset
	}

	return aVal, setCarry 
}
