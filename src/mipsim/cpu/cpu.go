package cpu

import (
	"fmt"
	"hex2mips/disassembler"
)

type CPU struct {
	PC     uint32
	Regs   [32]uint32
	Mem    map[uint32]uint32
	Hi, Lo uint32
	NextPC uint32
}

type ExecResult struct {
	RegWrite     bool
	RegDest      uint32
	RegWriteData uint32
	MemWrite     bool
	MemDest      uint32
	MemWriteData uint32
}

func New() *CPU {
	return &CPU{PC: 0x3000, Regs: [32]uint32{}, Mem: make(map[uint32]uint32)}
}

func signExtend16(x uint32) uint32 {
	if x&0x8000 != 0 {
		return x | 0xFFFF0000
	}
	return x
}

func (c *CPU) Run(instrs []uint32) {
	base := c.PC
	for i, w := range instrs {
		c.Mem[base+uint32(i*4)] = w
	}
	step := 1
	for {
		fmt.Printf("=== Step %d ===\n", step)
		step++
		idx := (c.PC - base) / 4
		if idx >= uint32(len(instrs)) {
			break
		}
		word := c.Mem[c.PC]
		c.NextPC = c.PC + 4
		asm := disassembler.DecodeWord(word, c.PC)
		res := c.Execute(word)
		fmt.Printf("Instr: 0x%08x   %s\n", word, asm)
		fmt.Printf("PC : 0x%08x\n", c.PC)
		if res.RegWrite {
			fmt.Printf("RegWrite : 1\n")
			fmt.Printf("RegDest : %d (%s)\n", res.RegDest, regName(res.RegDest))
			fmt.Printf("RegWriteData : 0x%08x\n", res.RegWriteData)
		} else {
			fmt.Printf("RegWrite : 0\n")
			fmt.Printf("RegDest : 0 ($zero)\n")
			fmt.Printf("RegWriteData : 0x00000000\n")
		}
		if res.MemWrite {
			fmt.Printf("MemWrite : 1\n")
			fmt.Printf("MemDest : 0x%08x\n", res.MemDest)
			fmt.Printf("MemWriteData : 0x%08x\n", res.MemWriteData)
		} else {
			fmt.Printf("MemWrite : 0\n")
			fmt.Printf("MemDest : 0x00000000\n")
			fmt.Printf("MemWriteData : 0x00000000\n")
		}
		c.PC = c.NextPC
	}
}

func regName(n uint32) string {
	names := []string{
		"$zero", "$at", "$v0", "$v1",
		"$a0", "$a1", "$a2", "$a3",
		"$t0", "$t1", "$t2", "$t3", "$t4", "$t5", "$t6", "$t7",
		"$s0", "$s1", "$s2", "$s3", "$s4", "$s5", "$s6", "$s7",
		"$t8", "$t9", "$k0", "$k1", "$gp", "$sp", "$fp", "$ra",
	}
	if n < 32 {
		return names[n]
	}
	return fmt.Sprintf("$r%d", n)
}
