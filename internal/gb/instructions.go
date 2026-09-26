package gb

type Instruction func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16)

func (cpu *CPU) initInstructions() {
	// nop
	cpu.instructions[0x00] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		return cpu.nop(), 4
	}

	// ld bc,n16
	cpu.instructions[0x01] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.setBC(bytesToWord(operands[1], operands[0]))
		return cpu.PC + 3, 12
	}

	// ld [bc],a
	cpu.instructions[0x02] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.gb.MemoryBus.WriteToAddress(cpu.getBC(), cpu.regs.A)
		return cpu.PC + 1, 8
	}

	// inc bc
	cpu.instructions[0x03] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.inc_nn(cpu.getBC, cpu.setBC)
		return cpu.PC + 1, 8
	}

	// inc b
	cpu.instructions[0x04] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.inc_n(&cpu.regs.B)
		return cpu.PC + 1, 4
	}

	// dec b
	cpu.instructions[0x05] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.dec_n(&cpu.regs.B)
		return cpu.PC + 1, 4
	}

	// ld b,n8
	cpu.instructions[0x06] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.B = operands[0]
		return cpu.PC + 2, 8
	}

	// rlca
	cpu.instructions[0x07] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.rlca()
		return cpu.PC + 1, 4
	}

	// ld [a16] sp
	cpu.instructions[0x08] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		addr := bytesToWord(operands[1], operands[0])
		cpu.gb.MemoryBus.WriteToAddress(addr, getLowerByte(cpu.SP))
		cpu.gb.MemoryBus.WriteToAddress(addr+1, getUpperByte(cpu.SP))
		return cpu.PC + 3, 20
	}

	// add hl,bc
	cpu.instructions[0x09] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		bc := cpu.getBC()
		cpu.add_hl(bc)
		return cpu.PC + 1, 8
	}

	// ld a,[bc]
	cpu.instructions[0x0a] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.A = cpu.gb.MemoryBus.ReadAddress(cpu.getBC())
		return cpu.PC + 1, 8
	}

	// dec bc
	cpu.instructions[0x0b] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.dec_nn(cpu.getBC, cpu.setBC)
		return cpu.PC + 1, 8
	}

	// inc c
	cpu.instructions[0x0c] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.inc_n(&cpu.regs.C)
		return cpu.PC + 1, 4
	}

	// dec c
	cpu.instructions[0x0d] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.dec_n(&cpu.regs.C)
		return cpu.PC + 1, 4
	}

	// ld c, n8
	cpu.instructions[0x0e] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.C = operands[0]
		return cpu.PC + 2, 8
	}

	// rrca
	cpu.instructions[0x0f] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.rrca()
		return cpu.PC + 1, 4
	}

	// stop
	cpu.instructions[0x10] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		// there's a lot of control flow that can be added here (carefully) but for a DMG emulator it's not necessary
		cpu.Halted = true
		return cpu.PC + 2, 4
	}

	// ld de, n16
	cpu.instructions[0x11] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.setDE(bytesToWord(operands[1], operands[0]))
		return cpu.PC + 3, 12
	}

	// ld [de], a
	cpu.instructions[0x12] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.gb.MemoryBus.WriteToAddress(cpu.getDE(), cpu.regs.A)
		return cpu.PC + 1, 8
	}

	// inc de
	cpu.instructions[0x13] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.inc_nn(cpu.getDE, cpu.setDE)
		return cpu.PC + 1, 8
	}

	// inc d
	cpu.instructions[0x14] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.inc_n(&cpu.regs.D)
		return cpu.PC + 1, 4
	}

	// dec d
	cpu.instructions[0x15] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.dec_n(&cpu.regs.D)
		return cpu.PC + 1, 4
	}
	// ld d, n8
	cpu.instructions[0x16] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.D = operands[0]
		return cpu.PC + 2, 8
	}
	// rla
	cpu.instructions[0x17] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.rla()
		return cpu.PC + 1, 4
	}

	// jr e8
	cpu.instructions[0x18] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		return cpu.jrCond(true, int8(operands[0])), 12
	}

	// add hl,de
	cpu.instructions[0x19] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		de := bytesToWord(cpu.regs.D, cpu.regs.E)
		cpu.add_hl(de)
		return cpu.PC + 1, 8
	}

	// ld a,[de]
	cpu.instructions[0x1a] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.A = cpu.gb.MemoryBus.ReadAddress(bytesToWord(cpu.regs.D, cpu.regs.E))
		return cpu.PC + 1, 8
	}

	// dec de
	cpu.instructions[0x1b] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.dec_nn(cpu.getDE, cpu.setDE)
		return cpu.PC + 1, 8
	}

	// inc e
	cpu.instructions[0x1c] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.inc_n(&cpu.regs.E)
		return cpu.PC + 1, 4
	}

	// dec e
	cpu.instructions[0x1d] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.dec_n(&cpu.regs.E)
		return cpu.PC + 1, 4
	}

	// ld e,n8
	cpu.instructions[0x1e] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.E = operands[0]
		return cpu.PC + 2, 8
	}

	// rra
	cpu.instructions[0x1f] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.rra()
		return cpu.PC + 1, 4
	}

	// jr nz,e8
	cpu.instructions[0x20] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		nz := !cpu.flags.Zero
		next := cpu.jrCond(nz, int8(operands[0]))
		if nz {
			return next, 12
		}
		return next, 8
	}

	// ld hl,n16
	cpu.instructions[0x21] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.setHL(bytesToWord(operands[1], operands[0]))
		return cpu.PC + 3, 12
	}

	// ld [hl+],a
	cpu.instructions[0x22] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		hl := cpu.getHL()
		cpu.gb.MemoryBus.WriteToAddress(hl, cpu.regs.A)
		cpu.setHL(hl + 1)
		return cpu.PC + 1, 8
	}

	// inc hl
	cpu.instructions[0x23] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.inc_nn(cpu.getHL, cpu.setHL)
		return cpu.PC + 1, 8
	}

	// inc h
	cpu.instructions[0x24] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.inc_n(&cpu.regs.H)
		return cpu.PC + 1, 4
	}

	// dec h
	cpu.instructions[0x25] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.dec_n(&cpu.regs.H)
		return cpu.PC + 1, 4
	}

	// ld h, n8
	cpu.instructions[0x26] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.H = operands[0]
		return cpu.PC + 2, 8
	}

	// daa
	cpu.instructions[0x27] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		res, setCarry := cpu.getDAA()
		cpu.regs.A = res
		cpu.flags.Zero = (cpu.regs.A == 0)
		cpu.flags.HalfCarry = false
		cpu.flags.Carry = setCarry
		return cpu.PC + 1, 4
	}

	// jr z,e8
	cpu.instructions[0x28] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		zero := cpu.flags.Zero
		next := cpu.jrCond(zero, int8(operands[0]))

		if zero {
			return next, 12
		}
		return next, 8
	}

	// add hl,hl
	cpu.instructions[0x29] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		hl := bytesToWord(cpu.regs.H, cpu.regs.L)
		cpu.add_hl(hl)
		return cpu.PC + 1, 8
	}

	// ld a,[hl+]
	cpu.instructions[0x2a] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		hl := bytesToWord(cpu.regs.H, cpu.regs.L)
		val := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.regs.A = val
		cpu.setHL(hl + 1)
		return cpu.PC + 1, 8
	}

	// dec hl
	cpu.instructions[0x2b] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.dec_nn(cpu.getHL, cpu.setHL)
		return cpu.PC + 1, 8
	}

	// inc l
	cpu.instructions[0x2c] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.inc_n(&cpu.regs.L)
		return cpu.PC + 1, 4
	}

	// dec l
	cpu.instructions[0x2d] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.dec_n(&cpu.regs.L)
		return cpu.PC + 1, 4
	}

	// ld l,n8
	cpu.instructions[0x2e] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.L = operands[0]
		return cpu.PC + 2, 8
	}

	// cpl
	cpu.instructions[0x2f] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.A = ^cpu.regs.A
		cpu.flags.HalfCarry = true
		cpu.flags.Subtract = true
		return cpu.PC + 1, 4
	}

	// jr nc,e8
	cpu.instructions[0x30] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		nc := !cpu.flags.Carry
		next := cpu.jrCond(nc, int8(operands[0]))
		if nc {
			return next, 12
		}
		return next, 8
	}

	// ld sp,n16
	cpu.instructions[0x31] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.SP = bytesToWord(operands[1], operands[0])
		return cpu.PC + 3, 12
	}

	// ld [hl-],a
	cpu.instructions[0x32] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		hl := cpu.getHL()
		cpu.gb.MemoryBus.WriteToAddress(hl, cpu.regs.A)
		cpu.setHL(hl - 1)
		return cpu.PC + 1, 8
	}

	// inc sp
	cpu.instructions[0x33] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.inc_nn(cpu.getSP, cpu.setSP)
		return cpu.PC + 1, 8
	}

	// inc [hl]
	cpu.instructions[0x34] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		hlByte := cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		cpu.inc_n(&hlByte)
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), hlByte)
		return cpu.PC + 1, 12
	}

	// dec [hl]
	cpu.instructions[0x35] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		hl := cpu.getHL()
		val := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.dec_n(&val)
		cpu.gb.MemoryBus.WriteToAddress(hl, val)
		return cpu.PC + 1, 12
	}

	// ld [hl],n8
	cpu.instructions[0x36] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), operands[0])
		return cpu.PC + 2, 12
	}

	// scf
	cpu.instructions[0x37] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.scf()
		return cpu.PC + 1, 4
	}

	// jr c,e8
	cpu.instructions[0x38] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		c := cpu.flags.Carry
		next := cpu.jrCond(c, int8(operands[0]))
		if c {
			return next, 12
		}
		return next, 8
	}

	// add hl,sp
	cpu.instructions[0x39] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add_hl(cpu.SP)
		return cpu.PC + 1, 8
	}

	// ld a,[hl-]
	cpu.instructions[0x3a] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		hlVal := bytesToWord(cpu.regs.H, cpu.regs.L)
		cpu.regs.A = cpu.gb.MemoryBus.ReadAddress(hlVal)
		cpu.setHL(hlVal - 1)
		return cpu.PC + 1, 8
	}

	// dec sp
	cpu.instructions[0x3b] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.dec_nn(cpu.getSP, cpu.setSP)
		return cpu.PC + 1, 8
	}

	// inc a
	cpu.instructions[0x3c] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.inc_n(&cpu.regs.A)
		return cpu.PC + 1, 4
	}

	// dec a
	cpu.instructions[0x3d] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.dec_n(&cpu.regs.A)
		return cpu.PC + 1, 4
	}

	// ld a,n8
	cpu.instructions[0x3e] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.A = operands[0]
		return cpu.PC + 2, 8
	}

	// ccf (complement carry flag)
	cpu.instructions[0x3f] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.ccf()
		return cpu.PC + 1, 4
	}

	// beginning of lots of load instructions
	// ld b,b
	cpu.instructions[0x40] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.B = cpu.regs.B
		return cpu.PC + 1, 4
	}

	// ld b,c
	cpu.instructions[0x41] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.B = cpu.regs.C
		return cpu.PC + 1, 4
	}

	// ld b,d
	cpu.instructions[0x42] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.B = cpu.regs.D
		return cpu.PC + 1, 4
	}

	// ld b,e
	cpu.instructions[0x43] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.B = cpu.regs.E
		return cpu.PC + 1, 4
	}

	// ld b,h
	cpu.instructions[0x44] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.B = cpu.regs.H
		return cpu.PC + 1, 4
	}

	// ld b,l
	cpu.instructions[0x45] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.B = cpu.regs.L
		return cpu.PC + 1, 4
	}

	// ld b,[hl]
	cpu.instructions[0x46] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.B = cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		return cpu.PC + 1, 8
	}

	// ld b,a
	cpu.instructions[0x47] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.B = cpu.regs.A
		return cpu.PC + 1, 4
	}

	// ld c,b
	cpu.instructions[0x48] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.C = cpu.regs.B
		return cpu.PC + 1, 4
	}

	// ld c,c
	cpu.instructions[0x49] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.C = cpu.regs.C
		return cpu.PC + 1, 4
	}

	// ld c,d
	cpu.instructions[0x4a] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.C = cpu.regs.D
		return cpu.PC + 1, 4
	}

	// ld c,e
	cpu.instructions[0x4b] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.C = cpu.regs.E
		return cpu.PC + 1, 4
	}

	// ld c,h
	cpu.instructions[0x4c] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.C = cpu.regs.H
		return cpu.PC + 1, 4
	}

	// ld c,l
	cpu.instructions[0x4d] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.C = cpu.regs.L
		return cpu.PC + 1, 4
	}

	// ld c,[hl]
	cpu.instructions[0x4e] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.C = cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		return cpu.PC + 1, 8
	}

	// ld c,a
	cpu.instructions[0x4f] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.C = cpu.regs.A
		return cpu.PC + 1, 4
	}

	// ld d,b
	cpu.instructions[0x50] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.D = cpu.regs.B
		return cpu.PC + 1, 4
	}

	// ld d,c
	cpu.instructions[0x51] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.D = cpu.regs.C
		return cpu.PC + 1, 4
	}

	// ld d,d
	cpu.instructions[0x52] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.D = cpu.regs.D
		return cpu.PC + 1, 4
	}

	// ld d,e
	cpu.instructions[0x53] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.D = cpu.regs.E
		return cpu.PC + 1, 4
	}

	// ld d,h
	cpu.instructions[0x54] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.D = cpu.regs.H
		return cpu.PC + 1, 4
	}

	// ld d,l
	cpu.instructions[0x55] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.D = cpu.regs.L
		return cpu.PC + 1, 4
	}

	// ld d,[hl]
	cpu.instructions[0x56] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.D = cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		return cpu.PC + 1, 8
	}

	// ld d,a
	cpu.instructions[0x57] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.D = cpu.regs.A
		return cpu.PC + 1, 4
	}

	// ld e,b
	cpu.instructions[0x58] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.E = cpu.regs.B
		return cpu.PC + 1, 4
	}

	// ld e,c
	cpu.instructions[0x59] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.E = cpu.regs.C
		return cpu.PC + 1, 4
	}

	// ld e,d
	cpu.instructions[0x5a] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.E = cpu.regs.D
		return cpu.PC + 1, 4
	}

	// ld e,e
	cpu.instructions[0x5b] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.E = cpu.regs.E
		return cpu.PC + 1, 4
	}

	// ld e,h
	cpu.instructions[0x5c] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.E = cpu.regs.H
		return cpu.PC + 1, 4
	}

	// ld e,l
	cpu.instructions[0x5d] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.E = cpu.regs.L
		return cpu.PC + 1, 4
	}

	// ld e,[hl]
	cpu.instructions[0x5e] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.E = cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		return cpu.PC + 1, 8
	}

	// ld e,a
	cpu.instructions[0x5f] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.E = cpu.regs.A
		return cpu.PC + 1, 4
	}

	// ld h,b
	cpu.instructions[0x60] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.H = cpu.regs.B
		return cpu.PC + 1, 4
	}

	// ld h,c
	cpu.instructions[0x61] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.H = cpu.regs.C
		return cpu.PC + 1, 4
	}

	// ld h,d
	cpu.instructions[0x62] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.H = cpu.regs.D
		return cpu.PC + 1, 4
	}

	// ld h,e
	cpu.instructions[0x63] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.H = cpu.regs.E
		return cpu.PC + 1, 4
	}

	// ld h,h
	cpu.instructions[0x64] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.H = cpu.regs.H
		return cpu.PC + 1, 4
	}

	// ld h,l
	cpu.instructions[0x65] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.H = cpu.regs.L
		return cpu.PC + 1, 4
	}

	// ld h,[hl]
	cpu.instructions[0x66] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.H = cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		return cpu.PC + 1, 8
	}

	// ld h,a
	cpu.instructions[0x67] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.H = cpu.regs.A
		return cpu.PC + 1, 4
	}

	// ld l,b
	cpu.instructions[0x68] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.L = cpu.regs.B
		return cpu.PC + 1, 4
	}

	// ld l,c
	cpu.instructions[0x69] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.L = cpu.regs.C
		return cpu.PC + 1, 4
	}

	// ld l,d
	cpu.instructions[0x6a] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.L = cpu.regs.D
		return cpu.PC + 1, 4
	}

	// ld l,e
	cpu.instructions[0x6b] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.L = cpu.regs.E
		return cpu.PC + 1, 4
	}

	// ld l,h
	cpu.instructions[0x6c] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.L = cpu.regs.H
		return cpu.PC + 1, 4
	}

	// ld l,h
	cpu.instructions[0x6d] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.L = cpu.regs.L
		return cpu.PC + 1, 4
	}

	// ld l,[hl]
	cpu.instructions[0x6e] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.L = cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		return cpu.PC + 1, 8
	}

	// ld l,a
	cpu.instructions[0x6f] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.L = cpu.regs.A
		return cpu.PC + 1, 4
	}

	// ld [hl],b
	cpu.instructions[0x70] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), cpu.regs.B)
		return cpu.PC + 1, 8
	}

	// ld [hl],c
	cpu.instructions[0x71] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), cpu.regs.C)
		return cpu.PC + 1, 8
	}

	// ld [hl],d
	cpu.instructions[0x72] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), cpu.regs.D)
		return cpu.PC + 1, 8
	}

	// ld [hl],e
	cpu.instructions[0x73] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), cpu.regs.E)
		return cpu.PC + 1, 8
	}

	// ld [hl],h
	cpu.instructions[0x74] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), cpu.regs.H)
		return cpu.PC + 1, 8
	}

	// ld [hl],l
	cpu.instructions[0x75] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), cpu.regs.L)
		return cpu.PC + 1, 8
	}

	// halt
	cpu.instructions[0x76] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.Halted = true
		// instr_timing has this at 0 t-cycles?
		return cpu.PC + 1, 4
	}

	// ld [hl],a
	cpu.instructions[0x77] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), cpu.regs.A)
		return cpu.PC + 1, 8
	}

	// ld a,b
	cpu.instructions[0x78] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.A = cpu.regs.B
		return cpu.PC + 1, 4
	}

	// ld a,c
	cpu.instructions[0x79] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.A = cpu.regs.C
		return cpu.PC + 1, 4
	}

	// ld a,d
	cpu.instructions[0x7a] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.A = cpu.regs.D
		return cpu.PC + 1, 4
	}

	// ld a,e
	cpu.instructions[0x7b] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.A = cpu.regs.E
		return cpu.PC + 1, 4
	}

	// ld a,h
	cpu.instructions[0x7c] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.A = cpu.regs.H
		return cpu.PC + 1, 4
	}

	// ld a,l
	cpu.instructions[0x7d] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.A = cpu.regs.L
		return cpu.PC + 1, 4
	}

	// ld a,[hl]
	cpu.instructions[0x7e] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.A = cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		return cpu.PC + 1, 8
	}

	// ld a,a
	cpu.instructions[0x7f] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.regs.A = cpu.regs.A
		return cpu.PC + 1, 4
	}

	// add instructions
	// add a,b
	cpu.instructions[0x80] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(cpu.regs.B, false)
		return cpu.PC + 1, 4
	}

	// add a,c
	cpu.instructions[0x81] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(cpu.regs.C, false)
		return cpu.PC + 1, 4
	}

	// add a,d
	cpu.instructions[0x82] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(cpu.regs.D, false)
		return cpu.PC + 1, 4
	}

	// add a,e
	cpu.instructions[0x83] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(cpu.regs.E, false)
		return cpu.PC + 1, 4
	}

	// add a,h
	cpu.instructions[0x84] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(cpu.regs.H, false)
		return cpu.PC + 1, 4
	}

	// add a,l
	cpu.instructions[0x85] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(cpu.regs.L, false)
		return cpu.PC + 1, 4
	}

	// add a,[hl]
	cpu.instructions[0x86] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		hl := cpu.getHL()
		hlVal := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.add(hlVal, false)
		return cpu.PC + 1, 8
	}

	// add a,a
	cpu.instructions[0x87] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(cpu.regs.A, false)
		return cpu.PC + 1, 4
	}

	// adc a,b
	cpu.instructions[0x88] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(cpu.regs.B, true)
		return cpu.PC + 1, 4
	}

	// adc a,c
	cpu.instructions[0x89] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(cpu.regs.C, true)
		return cpu.PC + 1, 4
	}

	// adc a,d
	cpu.instructions[0x8a] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(cpu.regs.D, true)
		return cpu.PC + 1, 4
	}

	// adc a,e
	cpu.instructions[0x8b] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(cpu.regs.E, true)
		return cpu.PC + 1, 4
	}

	// adc a,h
	cpu.instructions[0x8c] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(cpu.regs.H, true)
		return cpu.PC + 1, 4
	}

	// adc a,l
	cpu.instructions[0x8d] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(cpu.regs.L, true)
		return cpu.PC + 1, 4
	}

	// adc a,[hl]
	cpu.instructions[0x8e] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		hl := cpu.getHL()
		hlVal := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.add(hlVal, true)
		return cpu.PC + 1, 8
	}

	// adc a,a
	cpu.instructions[0x8f] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(cpu.regs.A, true)
		return cpu.PC + 1, 4
	}

	// sub a,b
	cpu.instructions[0x90] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(cpu.regs.B, false)
		return cpu.PC + 1, 4
	}

	// sub a,c
	cpu.instructions[0x91] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(cpu.regs.C, false)
		return cpu.PC + 1, 4
	}

	// sub a,d
	cpu.instructions[0x92] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(cpu.regs.D, false)
		return cpu.PC + 1, 4
	}

	// sub a,e
	cpu.instructions[0x93] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(cpu.regs.E, false)
		return cpu.PC + 1, 4
	}

	// sub a,h
	cpu.instructions[0x94] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(cpu.regs.H, false)
		return cpu.PC + 1, 4
	}

	// sub a,l
	cpu.instructions[0x95] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(cpu.regs.L, false)
		return cpu.PC + 1, 4
	}

	// sub a,[hl]
	cpu.instructions[0x96] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		hl := cpu.getHL()
		hlVal := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.sub(hlVal, false)
		return cpu.PC + 1, 8
	}

	// sub a,a
	cpu.instructions[0x97] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(cpu.regs.A, false)
		return cpu.PC + 1, 4
	}

	// sbc a,b
	cpu.instructions[0x98] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(cpu.regs.B, true)
		return cpu.PC + 1, 4
	}

	// sbc a,c
	cpu.instructions[0x99] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(cpu.regs.C, true)
		return cpu.PC + 1, 4
	}

	// sbc a,d
	cpu.instructions[0x9a] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(cpu.regs.D, true)
		return cpu.PC + 1, 4
	}

	// sbc a,e
	cpu.instructions[0x9b] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(cpu.regs.E, true)
		return cpu.PC + 1, 4
	}

	// sbc a,h
	cpu.instructions[0x9c] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(cpu.regs.H, true)
		return cpu.PC + 1, 4
	}

	// sbc a,l
	cpu.instructions[0x9d] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(cpu.regs.L, true)
		return cpu.PC + 1, 4
	}

	// sbc a,[hl]
	cpu.instructions[0x9e] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		hl := cpu.getHL()
		hlVal := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.sub(hlVal, true)
		return cpu.PC + 1, 8
	}

	// sbc a,a
	cpu.instructions[0x9f] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(cpu.regs.A, true)
		return cpu.PC + 1, 4
	}

	// and a,b
	cpu.instructions[0xa0] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.and(cpu.regs.B)
		return cpu.PC + 1, 4
	}

	// and a,c
	cpu.instructions[0xa1] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.and(cpu.regs.C)
		return cpu.PC + 1, 4
	}

	// and a,d
	cpu.instructions[0xa2] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.and(cpu.regs.D)
		return cpu.PC + 1, 4
	}

	// and a,e
	cpu.instructions[0xa3] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.and(cpu.regs.E)
		return cpu.PC + 1, 4
	}

	// and a,h
	cpu.instructions[0xa4] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.and(cpu.regs.H)
		return cpu.PC + 1, 4
	}

	// and a,l
	cpu.instructions[0xa5] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.and(cpu.regs.L)
		return cpu.PC + 1, 4
	}

	// and a,[hl]
	cpu.instructions[0xa6] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		hl := cpu.getHL()
		hlVal := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.and(hlVal)
		return cpu.PC + 1, 8
	}

	// and a,a
	cpu.instructions[0xa7] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.and(cpu.regs.A)
		return cpu.PC + 1, 4
	}

	// xor a,b
	cpu.instructions[0xa8] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.xor(cpu.regs.B)
		return cpu.PC + 1, 4
	}

	// xor a,c
	cpu.instructions[0xa9] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.xor(cpu.regs.C)
		return cpu.PC + 1, 4
	}

	// xor a,d
	cpu.instructions[0xaa] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.xor(cpu.regs.D)
		return cpu.PC + 1, 4
	}

	// xor a,e
	cpu.instructions[0xab] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.xor(cpu.regs.E)
		return cpu.PC + 1, 4
	}

	// xor a,h
	cpu.instructions[0xac] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.xor(cpu.regs.H)
		return cpu.PC + 1, 4
	}

	// xor a,l
	cpu.instructions[0xad] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.xor(cpu.regs.L)
		return cpu.PC + 1, 4
	}

	// xor a,[hl]
	cpu.instructions[0xae] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		hl := cpu.getHL()
		hlVal := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.xor(hlVal)
		return cpu.PC + 1, 8
	}

	// xor a,a
	cpu.instructions[0xaf] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.xor(cpu.regs.A)
		return cpu.PC + 1, 4
	}

	// or a,b
	cpu.instructions[0xb0] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.or(cpu.regs.B)
		return cpu.PC + 1, 4
	}

	// or a,c
	cpu.instructions[0xb1] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.or(cpu.regs.C)
		return cpu.PC + 1, 4
	}

	// or a,d
	cpu.instructions[0xb2] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.or(cpu.regs.D)
		return cpu.PC + 1, 4
	}

	// or a,e
	cpu.instructions[0xb3] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.or(cpu.regs.E)
		return cpu.PC + 1, 4
	}

	// or a,h
	cpu.instructions[0xb4] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.or(cpu.regs.H)
		return cpu.PC + 1, 4
	}

	// or a,l
	cpu.instructions[0xb5] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.or(cpu.regs.L)
		return cpu.PC + 1, 4
	}

	// or a,[hl]
	cpu.instructions[0xb6] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		hl := cpu.getHL()
		hlVal := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.or(hlVal)
		return cpu.PC + 1, 8
	}

	// or a,a
	cpu.instructions[0xb7] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.or(cpu.regs.A)
		return cpu.PC + 1, 4
	}

	// cp a,b
	cpu.instructions[0xb8] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.cp_n(cpu.regs.B)
		return cpu.PC + 1, 4
	}

	// cp a,c
	cpu.instructions[0xb9] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.cp_n(cpu.regs.C)
		return cpu.PC + 1, 4
	}

	// cp a,d
	cpu.instructions[0xba] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.cp_n(cpu.regs.D)
		return cpu.PC + 1, 4
	}

	// cp a,e
	cpu.instructions[0xbb] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.cp_n(cpu.regs.E)
		return cpu.PC + 1, 4
	}

	// cp a,h
	cpu.instructions[0xbc] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.cp_n(cpu.regs.H)
		return cpu.PC + 1, 4
	}

	// cp a,l
	cpu.instructions[0xbd] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.cp_n(cpu.regs.L)
		return cpu.PC + 1, 4
	}

	// cp a,[hl]
	cpu.instructions[0xbe] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		hl := cpu.getHL()
		hlVal := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.cp_n(hlVal)
		return cpu.PC + 1, 8
	}

	// cp a,a
	cpu.instructions[0xbf] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.cp_n(cpu.regs.A)
		return cpu.PC + 1, 4
	}

	// ret nz
	cpu.instructions[0xc0] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		nz := !cpu.flags.Zero

		nextPC = cpu.ret(nz)
		if nz {
			return nextPC, 20
		}
		return nextPC, 8
	}

	// pop bc
	cpu.instructions[0xc1] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		word := cpu.pop()
		cpu.setBC(word)
		return cpu.PC + 1, 12
	}

	// jp nz,a16
	cpu.instructions[0xc2] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		nz := !cpu.flags.Zero
		offset := bytesToWord(operands[1], operands[0])
		next := cpu.jpCond(nz, offset)
		if nz {
			return next, 16
		}
		return next, 12
	}

	// jp a16
	cpu.instructions[0xc3] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		offset := bytesToWord(operands[1], operands[0])
		return cpu.jpCond(true, offset), 16
	}

	// call nz,a16
	cpu.instructions[0xc4] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		nz := !cpu.flags.Zero
		target := uint16(operands[1])<<8 | uint16(operands[0])
		if nz {
			return cpu.call(nz, target), 24
		}
		return cpu.call(nz, target), 12
	}

	// push bc
	cpu.instructions[0xc5] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.push(bytesToWord(cpu.regs.B, cpu.regs.C))
		return cpu.PC + 1, 16
	}

	// add a,n8
	cpu.instructions[0xc6] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(operands[0], false)
		return cpu.PC + 2, 8
	}

	// rst $00
	cpu.instructions[0xc7] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		return cpu.rst(0x0000), 16
	}

	// ret z
	cpu.instructions[0xc8] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		z := cpu.flags.Zero
		if z {
			return cpu.ret(z), 20
		}
		return cpu.ret(z), 8
	}

	// ret
	cpu.instructions[0xc9] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		return cpu.ret(true), 16
	}

	// jp z,a16
	cpu.instructions[0xca] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		z := cpu.flags.Zero
		offset := bytesToWord(operands[1], operands[0])
		next := cpu.jpCond(z, offset)
		if z {
			return next, 16
		}
		return next, 12
	}

	// prefix
	cpu.instructions[0xcb] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		instr := cpu.cbInstructions[operands[0]]
		return instr(cpu, operands)
	}

	// call z,a16
	cpu.instructions[0xcc] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		z := cpu.flags.Zero
		target := uint16(operands[1])<<8 | uint16(operands[0])
		if z {
			return cpu.call(z, target), 24
		}
		return cpu.call(z, target), 12
	}

	// call a16
	cpu.instructions[0xcd] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		target := uint16(operands[1])<<8 | uint16(operands[0])
		return cpu.call(true, target), 24
	}

	// adc a,n8
	cpu.instructions[0xce] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.add(operands[0], true)
		return cpu.PC + 2, 8
	}

	// rst $08
	cpu.instructions[0xcf] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		return cpu.rst(0x0008), 16
	}

	// ret nc
	cpu.instructions[0xd0] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		nc := !cpu.flags.Carry
		if nc {
			return cpu.ret(nc), 20
		}
		return cpu.ret(nc), 8
	}

	// pop de
	cpu.instructions[0xd1] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		word := cpu.pop()
		cpu.setDE(word)
		return cpu.PC + 1, 12
	}

	// jp nc,a16
	cpu.instructions[0xd2] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		nc := !cpu.flags.Carry
		offset := bytesToWord(operands[1], operands[0])
		next := cpu.jpCond(nc, offset)
		if nc {
			return next, 16
		}
		return next, 12
	}

	// call nc,a16
	cpu.instructions[0xd4] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		nc := !cpu.flags.Carry
		target := uint16(operands[1])<<8 | uint16(operands[0])
		if nc {
			return cpu.call(nc, target), 24
		}
		return cpu.call(nc, target), 12
	}

	// push de
	cpu.instructions[0xd5] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.push(bytesToWord(cpu.regs.D, cpu.regs.E))
		return cpu.PC + 1, 16
	}

	// sub a,n8
	cpu.instructions[0xd6] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(operands[0], false)
		return cpu.PC + 2, 8
	}

	// rst $10
	cpu.instructions[0xd7] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		return cpu.rst(0x0010), 16
	}

	// ret c
	cpu.instructions[0xd8] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		c := cpu.flags.Carry
		if c {
			return cpu.ret(c), 20
		}
		return cpu.ret(c), 8
	}

	// reti
	cpu.instructions[0xd9] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.IMEPending = true
		return cpu.ret(true), 16
	}

	// jp c,a16
	cpu.instructions[0xda] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		c := cpu.flags.Carry
		offset := bytesToWord(operands[1], operands[0])
		next := cpu.jpCond(c, offset)
		if c {
			return next, 16
		}
		return next, 12
	}

	// call c,a16
	cpu.instructions[0xdc] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		c := cpu.flags.Carry
		target := uint16(operands[1])<<8 | uint16(operands[0])
		if c {
			return cpu.call(c, target), 24
		}
		return cpu.call(c, target), 12
	}

	// sbc a,n8
	cpu.instructions[0xde] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.sub(operands[0], true)
		return cpu.PC + 2, 8
	}

	// rst $18
	cpu.instructions[0xdf] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		return cpu.rst(0x0018), 16
	}

	// ldh [a8],a
	cpu.instructions[0xe0] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		a8 := 0xff00 + uint16(operands[0])
		cpu.gb.MemoryBus.WriteToAddress(a8, cpu.regs.A)
		return cpu.PC + 2, 12
	}

	// pop hl
	cpu.instructions[0xe1] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		word := cpu.pop()
		cpu.setHL(word)
		return cpu.PC + 1, 12
	}

	// ldh [c],a
	cpu.instructions[0xe2] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.gb.MemoryBus.WriteToAddress(0xff00+uint16(cpu.regs.C), cpu.regs.A)
		return cpu.PC + 1, 8
	}

	// push hl
	cpu.instructions[0xe5] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.push(bytesToWord(cpu.regs.H, cpu.regs.L))
		return cpu.PC + 1, 16
	}

	// and a,n8
	cpu.instructions[0xe6] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.and(operands[0])
		return cpu.PC + 2, 8
	}

	// rst $20
	cpu.instructions[0xe7] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		return cpu.rst(0x0020), 16
	}

	// add sp,e8
	cpu.instructions[0xe8] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		sp := cpu.getSP()
		cpu.flags.Zero = false
		cpu.flags.Subtract = false

		e8 := int16(int8(operands[0]))
		res := uint16(int32(sp) + int32(e8))

		cpu.flags.HalfCarry = (sp&0x0f)+(uint16(operands[0]&0x0f)) > 0x0f
		cpu.flags.Carry = (sp&0xff)+uint16(operands[0]) > 0xff

		cpu.setSP(res)
		return cpu.PC + 2, 16
	}

	// jp hl
	cpu.instructions[0xe9] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		return bytesToWord(cpu.regs.H, cpu.regs.L), 4
	}

	// ld [a16],a
	cpu.instructions[0xea] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		addr := bytesToWord(operands[1], operands[0])
		cpu.gb.MemoryBus.WriteToAddress(addr, cpu.regs.A)
		return cpu.PC + 3, 16
	}

	// xor a,n8
	cpu.instructions[0xee] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.xor(operands[0])
		return cpu.PC + 2, 8
	}

	// rst $28
	cpu.instructions[0xef] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		return cpu.rst(0x0028), 16
	}

	// ldh a,[a8]
	cpu.instructions[0xf0] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		val := cpu.gb.MemoryBus.ReadAddress(0xff00 + uint16(operands[0]))
		cpu.regs.A = val
		return cpu.PC + 2, 12
	}

	// pop af
	cpu.instructions[0xf1] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		word := cpu.pop()
		cpu.regs.A = getUpperByte(word)
		cpu.flags.convertFromByte(getLowerByte(word))
		return cpu.PC + 1, 12
	}

	// ldh a, [c]
	cpu.instructions[0xf2] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		addr := 0xff00 + uint16(cpu.regs.C)
		cpu.regs.A = cpu.gb.MemoryBus.ReadAddress(addr)
		return cpu.PC + 1, 8
	}

	// DI
	cpu.instructions[0xf3] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.IME = false
		return cpu.PC + 1, 4
	}

	// push af
	cpu.instructions[0xf5] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.push(bytesToWord(cpu.regs.A, cpu.flags.convertToByte()))
		return cpu.PC + 1, 16
	}

	// or a,n8
	cpu.instructions[0xf6] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.or(operands[0])
		return cpu.PC + 2, 8
	}

	// rst $30
	cpu.instructions[0xf7] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		return cpu.rst(0x0030), 16
	}

	// ld hl,sp + e8
	cpu.instructions[0xf8] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		sp := cpu.getSP()
		cpu.flags.Zero = false
		cpu.flags.Subtract = false

		e8 := int16(int8(operands[0]))
		res := uint16(int32(sp) + int32(e8))
		cpu.setHL(res)

		cpu.flags.HalfCarry = (sp&0x0f)+(uint16(operands[0]&0x0f)) > 0x0f
		cpu.flags.Carry = (sp&0xff)+uint16(operands[0]) > 0xff

		return cpu.PC + 2, 12
	}

	// ld sp,hl
	cpu.instructions[0xf9] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.SP = cpu.getHL()
		return cpu.PC + 1, 8
	}

	// ld a,[a16]
	cpu.instructions[0xfa] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		addr := bytesToWord(operands[1], operands[0])
		cpu.regs.A = cpu.gb.MemoryBus.ReadAddress(addr)
		return cpu.PC + 3, 16
	}

	// ei
	cpu.instructions[0xfb] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.IMEPending = true
		return cpu.PC + 1, 4
	}

	// cp a,n8
	cpu.instructions[0xfe] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		cpu.cp_n(operands[0])
		return cpu.PC + 2, 8
	}

	// rst $38
	cpu.instructions[0xff] = func(cpu *CPU, operands []byte) (nextPC uint16, cycles uint16) {
		return cpu.rst(0x0038), 16
	}
}
