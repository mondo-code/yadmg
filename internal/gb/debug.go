package gb 

import "fmt"

// this is gameboy-doctor's preferred debug output
func (cpu *CPU) LogInstruction() {
	fmt.Printf("A:%02X F:%02X B:%02X C:%02X D:%02X E:%02X H:%02X L:%02X SP:%04X PC:%04X PCMEM:%02X,%02X,%02X,%02X\n",
				cpu.regs.A, cpu.flags.convertToByte(), cpu.regs.B, cpu.regs.C, cpu.regs.D, cpu.regs.E, cpu.regs.H, cpu.regs.L,
				cpu.SP, cpu.PC, cpu.gb.MemoryBus.ReadAddress(cpu.PC), 
				cpu.gb.MemoryBus.ReadAddress(cpu.PC+1), cpu.gb.MemoryBus.ReadAddress(cpu.PC+2), cpu.gb.MemoryBus.ReadAddress(cpu.PC+3))
}
