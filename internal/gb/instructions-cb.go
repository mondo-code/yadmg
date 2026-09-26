package gb

func (cpu *CPU) initInstructionsCB() {
	// rlc instructions
	cpu.cbInstructions[0x00] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rlc(&cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x01] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rlc(&cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x02] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rlc(&cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x03] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rlc(&cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x04] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rlc(&cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x05] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rlc(&cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x06] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hlByte := cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		cpu.rlc(&hlByte)
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0x07] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rlc(&cpu.regs.A)
		return cpu.PC + 2, 8
	}

	// rrc instructions
	cpu.cbInstructions[0x08] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rrc(&cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x09] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rrc(&cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x0a] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rrc(&cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x0b] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rrc(&cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x0c] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rrc(&cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x0d] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rrc(&cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x0e] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hlByte := cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		cpu.rrc(&hlByte)
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0x0f] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rrc(&cpu.regs.A)
		return cpu.PC + 2, 8
	}

	// rl instructions
	cpu.cbInstructions[0x10] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rl(&cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x11] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rl(&cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x12] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rl(&cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x13] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rl(&cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x14] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rl(&cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x15] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rl(&cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x16] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hlByte := cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		cpu.rl(&hlByte)
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0x17] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rl(&cpu.regs.A)
		return cpu.PC + 2, 8
	}

	// rr instructions
	cpu.cbInstructions[0x18] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rr(&cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x19] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rr(&cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x1a] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rr(&cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x1b] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rr(&cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x1c] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rr(&cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x1d] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rr(&cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x1e] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hlByte := cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		cpu.rr(&hlByte)
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0x1f] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.rr(&cpu.regs.A)
		return cpu.PC + 2, 8
	}

	// sla instructions
	cpu.cbInstructions[0x20] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.sla(&cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x21] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.sla(&cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x22] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.sla(&cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x23] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.sla(&cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x24] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.sla(&cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x25] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.sla(&cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x26] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hlByte := cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		cpu.sla(&hlByte)
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0x27] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.sla(&cpu.regs.A)
		return cpu.PC + 2, 8
	}

	// sra instructions
	cpu.cbInstructions[0x28] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.sra(&cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x29] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.sra(&cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x2a] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.sra(&cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x2b] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.sra(&cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x2c] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.sra(&cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x2d] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.sra(&cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x2e] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hlByte := cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		cpu.sra(&hlByte)
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0x2f] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.sra(&cpu.regs.A)
		return cpu.PC + 2, 8
	}

	// swap instructions
	cpu.cbInstructions[0x30] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.swap(&cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x31] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.swap(&cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x32] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.swap(&cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x33] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.swap(&cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x34] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.swap(&cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x35] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.swap(&cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x36] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hlByte := cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		cpu.swap(&hlByte)
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0x37] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.swap(&cpu.regs.A)
		return cpu.PC + 2, 8
	}

	// srl instructions
	cpu.cbInstructions[0x38] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.srl(&cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x39] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.srl(&cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x3a] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.srl(&cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x3b] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.srl(&cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x3c] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.srl(&cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x3d] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.srl(&cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x3e] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hlByte := cpu.gb.MemoryBus.ReadAddress(cpu.getHL())
		cpu.srl(&hlByte)
		cpu.gb.MemoryBus.WriteToAddress(cpu.getHL(), hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0x3f] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.srl(&cpu.regs.A)
		return cpu.PC + 2, 8
	}

	// bit test instructions
	cpu.cbInstructions[0x40] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(0, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x41] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(0, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x42] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(0, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x43] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(0, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x44] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(0, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x45] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(0, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x46] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := bytesToWord(cpu.regs.H, cpu.regs.L)
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.bit(0, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 12
	}

	cpu.cbInstructions[0x47] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(0, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x48] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(1, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x49] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(1, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x4a] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(1, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x4b] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(1, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x4c] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(1, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x4d] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(1, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x4e] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := bytesToWord(cpu.regs.H, cpu.regs.L)
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.bit(1, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 12
	}

	cpu.cbInstructions[0x4f] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(1, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x50] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(2, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x51] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(2, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x52] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(2, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x53] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(2, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x54] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(2, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x55] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(2, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x56] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := bytesToWord(cpu.regs.H, cpu.regs.L)
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.bit(2, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 12
	}

	cpu.cbInstructions[0x57] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(2, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x58] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(3, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x59] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(3, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x5a] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(3, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x5b] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(3, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x5c] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(3, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x5d] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(3, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x5e] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := bytesToWord(cpu.regs.H, cpu.regs.L)
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.bit(3, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 12
	}

	cpu.cbInstructions[0x5f] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(3, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x60] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(4, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x61] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(4, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x62] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(4, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x63] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(4, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x64] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(4, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x65] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(4, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x66] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.bit(4, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 12
	}

	cpu.cbInstructions[0x67] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(4, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x68] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(5, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x69] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(5, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x6a] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(5, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x6b] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(5, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x6c] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(5, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x6d] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(5, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x6e] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.bit(5, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 12
	}

	cpu.cbInstructions[0x6f] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(5, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x70] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(6, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x71] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(6, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x72] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(6, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x73] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(6, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x74] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(6, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x75] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(6, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x76] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.bit(6, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 12
	}

	cpu.cbInstructions[0x77] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(6, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x78] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(7, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x79] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(7, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x7a] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(7, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x7b] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(7, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x7c] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(7, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x7d] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(7, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x7e] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.bit(7, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 12
	}

	cpu.cbInstructions[0x7f] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.bit(7, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	// res instructions
	cpu.cbInstructions[0x80] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(0, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x81] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(0, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x82] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(0, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x83] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(0, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x84] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(0, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x85] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(0, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x86] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.reset(0, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0x87] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(0, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x88] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(1, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x89] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(1, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x8a] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(1, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x8b] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(1, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x8c] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(1, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x8d] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(1, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x8e] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.reset(1, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0x8f] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(1, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x90] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(2, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x91] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(2, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x92] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(2, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x93] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(2, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x94] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(2, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x95] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(2, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x96] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.reset(2, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0x97] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(2, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x98] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(3, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x99] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(3, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x9a] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(3, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x9b] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(3, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x9c] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(3, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x9d] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(3, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0x9e] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.reset(3, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0x9f] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(3, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xa0] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(4, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xa1] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(4, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xa2] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(4, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xa3] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(4, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xa4] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(4, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xa5] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(4, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xa6] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.reset(4, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0xa7] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(4, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xa8] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(5, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xa9] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(5, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xaa] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(5, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xab] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(5, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xac] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(5, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xad] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(5, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xae] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.reset(5, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0xaf] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(5, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xb0] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(6, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xb1] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(6, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xb2] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(6, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xb3] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(6, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xb4] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(6, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xb5] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(6, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xb6] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.reset(6, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0xb7] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(6, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xb8] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(7, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xb9] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(7, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xba] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(7, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xbb] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(7, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xbc] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(7, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xbd] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(7, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xbe] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.reset(7, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0xbf] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.reset(7, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	// set instructions
	cpu.cbInstructions[0xc0] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(0, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xc1] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(0, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xc2] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(0, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xc3] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(0, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xc4] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(0, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xc5] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(0, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xc6] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.set(0, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0xc7] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(0, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xc8] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(1, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xc9] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(1, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xca] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(1, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xcb] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(1, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xcc] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(1, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xcd] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(1, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xce] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.set(1, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0xcf] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(1, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xd0] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(2, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xd1] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(2, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xd2] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(2, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xd3] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(2, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xd4] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(2, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xd5] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(2, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xd6] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.set(2, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0xd7] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(2, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xd8] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(3, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xd9] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(3, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xda] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(3, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xdb] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(3, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xdc] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(3, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xdd] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(3, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xde] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.set(3, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0xdf] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(3, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xe0] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(4, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xe1] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(4, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xe2] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(4, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xe3] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(4, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xe4] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(4, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xe5] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(4, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xe6] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.set(4, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0xe7] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(4, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xe8] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(5, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xe9] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(5, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xea] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(5, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xeb] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(5, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xec] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(5, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xed] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(5, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xee] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.set(5, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0xef] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(5, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xf0] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(6, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xf1] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(6, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xf2] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(6, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xf3] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(6, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xf4] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(6, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xf5] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(6, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xf6] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.set(6, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0xf7] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(6, &cpu.regs.A)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xf8] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(7, &cpu.regs.B)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xf9] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(7, &cpu.regs.C)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xfa] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(7, &cpu.regs.D)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xfb] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(7, &cpu.regs.E)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xfc] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(7, &cpu.regs.H)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xfd] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(7, &cpu.regs.L)
		return cpu.PC + 2, 8
	}

	cpu.cbInstructions[0xfe] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		hl := cpu.getHL()
		hlByte := cpu.gb.MemoryBus.ReadAddress(hl)
		cpu.set(7, &hlByte)
		cpu.gb.MemoryBus.WriteToAddress(hl, hlByte)
		return cpu.PC + 2, 16
	}

	cpu.cbInstructions[0xff] = func(cpu *CPU, operands []byte) (length uint16, cycles uint16) {
		cpu.set(7, &cpu.regs.A)
		return cpu.PC + 2, 8
	}
}
