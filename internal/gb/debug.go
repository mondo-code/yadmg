package gb

import "fmt"

// this is gameboy-doctor's preferred debug output
func (cpu *CPU) logInstruction() {
	fmt.Printf("A:%02X F:%02X B:%02X C:%02X D:%02X E:%02X H:%02X L:%02X SP:%04X PC:%04X PCMEM:%02X,%02X,%02X,%02X\n",
		cpu.regs.A, cpu.flags.convertToByte(), cpu.regs.B, cpu.regs.C, cpu.regs.D, cpu.regs.E, cpu.regs.H, cpu.regs.L,
		cpu.sp, cpu.pc, cpu.gb.memoryBus.readAddress(cpu.pc),
		cpu.gb.memoryBus.readAddress(cpu.pc+1), cpu.gb.memoryBus.readAddress(cpu.pc+2), cpu.gb.memoryBus.readAddress(cpu.pc+3))
}
