package disassembler

import (
	"fmt"
)

var regs = []string{
	"$zero", "$at", "$v0", "$v1",
	"$a0", "$a1", "$a2", "$a3",
	"$t0", "$t1", "$t2", "$t3", "$t4", "$t5", "$t6", "$t7",
	"$s0", "$s1", "$s2", "$s3", "$s4", "$s5", "$s6", "$s7",
	"$t8", "$t9", "$k0", "$k1", "$gp", "$sp", "$fp", "$ra",
}

func reg(n uint32) string {
	if n < 32 {
		return regs[n]
	}
	return fmt.Sprintf("$r%d", n)
}

func signExtend16(x uint32) int32 {
	if x&0x8000 != 0 {
		return int32(x) | ^0xFFFF
	}
	return int32(x)
}

func decodeR(opcode, rs, rt, rd, shamt, funct, pc uint32) string {
	rMap := map[uint32]string{
		0x20: "add", 0x21: "addu", 0x22: "sub", 0x23: "subu",
		0x24: "and", 0x25: "or", 0x26: "xor", 0x27: "nor",
		0x00: "sll", 0x02: "srl", 0x03: "sra",
		0x04: "sllv", 0x06: "srlv", 0x07: "srav",
		0x08: "jr", 0x09: "jalr",
		0x0C: "syscall", 0x0D: "break",
		0x10: "mfhi", 0x11: "mthi", 0x12: "mflo", 0x13: "mtlo",
		0x18: "mult", 0x19: "multu", 0x1A: "div", 0x1B: "divu",
		0x2a: "slt", 0x2b: "sltu",
	}

	switch funct {
	case 0x00, 0x02, 0x03: // sll, srl, sra
		if name, ok := rMap[funct]; ok {
			if rd == 0 && rt == 0 && shamt == 0 { // nop
				return "nop"
			}
			return fmt.Sprintf("%s %s, %s, %d", name, reg(rd), reg(rt), shamt)
		}
	case 0x04, 0x06, 0x07: // sllv, srlv, srav
		if name, ok := rMap[funct]; ok {
			return fmt.Sprintf("%s %s, %s, %s", name, reg(rd), reg(rt), reg(rs))
		}
	case 0x08: // jr
		return fmt.Sprintf("jr %s", reg(rs))
	case 0x09: // jalr
		if rd == 31 {
			return fmt.Sprintf("jalr %s", reg(rs))
		}
		return fmt.Sprintf("jalr %s, %s", reg(rd), reg(rs))
	case 0x0C:
		return "syscall"
	case 0x10, 0x12: // mfhi, mflo
		if name, ok := rMap[funct]; ok {
			return fmt.Sprintf("%s %s", name, reg(rd))
		}
	case 0x11, 0x13: // mthi, mtlo
		if name, ok := rMap[funct]; ok {
			return fmt.Sprintf("%s %s", name, reg(rs))
		}
	case 0x18, 0x19, 0x1A, 0x1B: // mult, multu, div, divu
		if name, ok := rMap[funct]; ok {
			return fmt.Sprintf("%s %s, %s", name, reg(rs), reg(rt))
		}
	default:
		if name, ok := rMap[funct]; ok {
			return fmt.Sprintf("%s %s, %s, %s", name, reg(rd), reg(rs), reg(rt))
		}
	}
	return fmt.Sprintf("Rtype_unknown_funct_0x%02x", funct)
}

func decodeI(opcode, rs, rt, imm, pc uint32) string {
	imms := signExtend16(imm)
	iMap := map[uint32]string{
		0x04: "beq", 0x05: "bne", 0x06: "blez", 0x07: "bgtz",
		0x08: "addi", 0x09: "addiu", 0x0A: "slti", 0x0B: "sltiu",
		0x0C: "andi", 0x0D: "ori", 0x0E: "xori", 0x0F: "lui",
		0x20: "lb", 0x21: "lh", 0x23: "lw", 0x24: "lbu", 0x25: "lhu",
		0x28: "sb", 0x29: "sh", 0x2B: "sw",
	}

	switch opcode {
	case 0x01: // REGIMM
		regimmMap := map[uint32]string{0x00: "bltz", 0x01: "bgez", 0x10: "bltzal", 0x11: "bgezal"}
		if name, ok := regimmMap[rt]; ok {
			target := (pc + 4) + (uint32(imms) << 2)
			return fmt.Sprintf("%s %s, 0x%08x", name, reg(rs), target)
		}
	case 0x04, 0x05: // beq, bne
		name := iMap[opcode]
		target := (pc + 4) + (uint32(imms) << 2)
		return fmt.Sprintf("%s %s, %s, 0x%08x", name, reg(rs), reg(rt), target)
	case 0x06, 0x07: // blez, bgtz
		name := iMap[opcode]
		target := (pc + 4) + (uint32(imms) << 2)
		return fmt.Sprintf("%s %s, 0x%08x", name, reg(rs), target)
	case 0x0F: // lui
		return fmt.Sprintf("lui %s, 0x%04x", reg(rt), imm)
	case 0x08, 0x09, 0x0A, 0x0B: // addi, addiu, slti, sltiu
		name := iMap[opcode]
		return fmt.Sprintf("%s %s, %s, %d", name, reg(rt), reg(rs), imms)
	case 0x0C, 0x0D, 0x0E: // andi, ori, xori
		name := iMap[opcode]
		return fmt.Sprintf("%s %s, %s, 0x%04x", name, reg(rt), reg(rs), imm)
	case 0x20, 0x21, 0x23, 0x24, 0x25, 0x28, 0x29, 0x2B: // memory
		name := iMap[opcode]
		return fmt.Sprintf("%s %s, %d(%s)", name, reg(rt), imms, reg(rs))
	}
	return fmt.Sprintf("Itype_unknown_op_0x%02x", opcode)
}

func decodeJ(opcode, addr, pc uint32) string {
	jMap := map[uint32]string{0x02: "j", 0x03: "jal"}
	if name, ok := jMap[opcode]; ok {
		target := ((pc + 4) & 0xF0000000) | (addr << 2)
		return fmt.Sprintf("%s 0x%08x", name, target)
	}
	return fmt.Sprintf("Jtype_unknown_op_0x%02x", opcode)
}

// DecodeWord decodes a single 32-bit MIPS instruction word.
func DecodeWord(word uint32, pc uint32) string {
	opcode := word >> 26
	if word == 0 {
		return "nop"
	}

	switch opcode {
	case 0x00: // R-type
		rs := (word >> 21) & 0x1F
		rt := (word >> 16) & 0x1F
		rd := (word >> 11) & 0x1F
		shamt := (word >> 6) & 0x1F
		funct := word & 0x3F
		return decodeR(opcode, rs, rt, rd, shamt, funct, pc)
	case 0x02, 0x03: // J-type
		addr := word & 0x03FFFFFF
		return decodeJ(opcode, addr, pc)
	default: // I-type
		rs := (word >> 21) & 0x1F
		rt := (word >> 16) & 0x1F
		imm := word & 0xFFFF
		return decodeI(opcode, rs, rt, imm, pc)
	}
}
