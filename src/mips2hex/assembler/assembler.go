package assembler

import (
	"fmt"
	"strconv"
	"strings"

	"mips2hex/regs"
	"mips2hex/types"
)

var instrTable = map[string]types.Instruction{
	"add":   {Type: types.RType, Opcode: 0x00, Funct: 0x20},
	"addu":  {Type: types.RType, Opcode: 0x00, Funct: 0x21},
	"sub":   {Type: types.RType, Opcode: 0x00, Funct: 0x22},
	"subu":  {Type: types.RType, Opcode: 0x00, Funct: 0x23},
	"and":   {Type: types.RType, Opcode: 0x00, Funct: 0x24},
	"or":    {Type: types.RType, Opcode: 0x00, Funct: 0x25},
	"xor":   {Type: types.RType, Opcode: 0x00, Funct: 0x26},
	"nor":   {Type: types.RType, Opcode: 0x00, Funct: 0x27},
	"slt":   {Type: types.RType, Opcode: 0x00, Funct: 0x2a},
	"sltu":  {Type: types.RType, Opcode: 0x00, Funct: 0x2b},
	"sll":   {Type: types.RType, Opcode: 0x00, Funct: 0x00}, // Special R-type format (rd, rt, shamt)
	"srl":   {Type: types.RType, Opcode: 0x00, Funct: 0x02}, // Special R-type format (rd, rt, shamt)
	"sra":   {Type: types.RType, Opcode: 0x00, Funct: 0x03}, // Special R-type format (rd, rt, shamt)
	"sllv":  {Type: types.RType, Opcode: 0x00, Funct: 0x04}, // Special R-type format (rd, rs, rt)
	"srlv":  {Type: types.RType, Opcode: 0x00, Funct: 0x06}, // Special R-type format (rd, rs, rt)
	"srav":  {Type: types.RType, Opcode: 0x00, Funct: 0x07}, // Special R-type format (rd, rs, rt)
	"jr":    {Type: types.RType, Opcode: 0x00, Funct: 0x08}, // Special R-type format (rs)
	"jalr":  {Type: types.RType, Opcode: 0x00, Funct: 0x09}, // Special R-type format (rd, rs)
	"mult":  {Type: types.RType, Opcode: 0x00, Funct: 0x18},
	"multu": {Type: types.RType, Opcode: 0x00, Funct: 0x19},
	"div":   {Type: types.RType, Opcode: 0x00, Funct: 0x1a},
	"divu":  {Type: types.RType, Opcode: 0x00, Funct: 0x1b},
	"mfhi":  {Type: types.RType, Opcode: 0x00, Funct: 0x10},
	"mflo":  {Type: types.RType, Opcode: 0x00, Funct: 0x12},
	"mthi":  {Type: types.RType, Opcode: 0x00, Funct: 0x11},
	"mtlo":  {Type: types.RType, Opcode: 0x00, Funct: 0x13},

	"addi":  {Type: types.IType, Opcode: 0x08},
	"addiu": {Type: types.IType, Opcode: 0x09},
	"andi":  {Type: types.IType, Opcode: 0x0c},
	"ori":   {Type: types.IType, Opcode: 0x0d},
	"xori":  {Type: types.IType, Opcode: 0x0e},
	"slti":  {Type: types.IType, Opcode: 0x0a},
	"sltiu": {Type: types.IType, Opcode: 0x0b},
	"lui":   {Type: types.IType, Opcode: 0x0f}, // Special I-type format (rt, imm)
	"lw":    {Type: types.IType, Opcode: 0x23}, // Special I-type format rt, offset(base)
	"lh":    {Type: types.IType, Opcode: 0x21},
	"lhu":   {Type: types.IType, Opcode: 0x25},
	"lb":    {Type: types.IType, Opcode: 0x20},
	"lbu":   {Type: types.IType, Opcode: 0x24},
	"sw":    {Type: types.IType, Opcode: 0x2b}, // Special I-type format rt, offset(base)
	"sh":    {Type: types.IType, Opcode: 0x29},
	"sb":    {Type: types.IType, Opcode: 0x28},
	"beq":   {Type: types.IType, Opcode: 0x04}, // Special I-type format rs, rt, offset
	"bne":   {Type: types.IType, Opcode: 0x05}, // Special I-type format rs, rt, offset
	"blez":  {Type: types.IType, Opcode: 0x06},
	"bgtz":  {Type: types.IType, Opcode: 0x07},

	"j":   {Type: types.JType, Opcode: 0x02},
	"jal": {Type: types.JType, Opcode: 0x03},

	"li":  {Type: types.Special},
	"nop": {Type: types.Special},
}

// Assemble: 将解析出来的 items 与 label 表翻译为机器码 uint32 列表
func Assemble(items []types.Item, labels map[string]uint32, base uint32) ([]uint32, error) {
	var out []uint32
	addr := uint32(0)

	for _, it := range items {
		switch it.Kind {
		case types.Word:
			v, err := parseNumber(it.Raw, labels, base)
			if err != nil {
				return nil, fmt.Errorf("line %d: 解析 .word %s 失败: %v", it.LineNo, it.Raw, err)
			}
			out = append(out, v)
			addr += 4

		case types.Instr:
			op := strings.ToLower(it.Tokens[0])
			instr, ok := instrTable[op]
			if !ok {
				return nil, fmt.Errorf("line %d: 不支持的指令: %s", it.LineNo, op)
			}

			var words []uint32
			var err error

			switch instr.Type {
			case types.RType:
				words, err = assembleRType(it, instr)
			case types.IType:
				words, err = assembleIType(it, instr, labels, addr, base)
			case types.JType:
				words, err = assembleJType(it, instr, labels, base)
			case types.Special:
				words, err = assembleSpecial(it, labels, base)
			default:
				err = fmt.Errorf("line %d: 未知指令类型 for %s", it.LineNo, op)
			}

			if err != nil {
				return nil, err
			}
			out = append(out, words...)
			addr += uint32(len(words) * 4)

		default:
			return nil, fmt.Errorf("未知 item 类型")
		}
	}
	return out, nil
}

func assembleRType(it types.Item, instr types.Instruction) ([]uint32, error) {
	toks := it.Tokens
	op := strings.ToLower(toks[0])
	var word uint32

	switch op {
	// Format: op rd, rs, rt
	case "add", "addu", "sub", "subu", "and", "or", "xor", "nor", "slt", "sltu":
		if len(toks) < 4 {
			return nil, fmt.Errorf("line %d: %s 需要 3 个寄存器操作数", it.LineNo, op)
		}
		rd := regs.RegOf(toks[1], it.LineNo)
		rs := regs.RegOf(toks[2], it.LineNo)
		rt := regs.RegOf(toks[3], it.LineNo)
		word = (uint32(rs) << 21) | (uint32(rt) << 16) | (uint32(rd) << 11) | instr.Funct
	// Format: op rd, rt, rs (variable shifts)
	case "sllv", "srlv", "srav":
		if len(toks) < 4 {
			return nil, fmt.Errorf("line %d: %s 需要 3 个寄存器操作数", it.LineNo, op)
		}
		rd := regs.RegOf(toks[1], it.LineNo)
		rt := regs.RegOf(toks[2], it.LineNo)
		rs := regs.RegOf(toks[3], it.LineNo)
		word = (uint32(rs) << 21) | (uint32(rt) << 16) | (uint32(rd) << 11) | instr.Funct
	// Format: op rd, rt, shamt
	case "sll", "srl", "sra":
		if len(toks) < 4 {
			return nil, fmt.Errorf("line %d: %s 需要 rd, rt, shamt", it.LineNo, op)
		}
		rd := regs.RegOf(toks[1], it.LineNo)
		rt := regs.RegOf(toks[2], it.LineNo)
		shamt, err := strconv.Atoi(strings.TrimSpace(toks[3]))
		if err != nil {
			return nil, fmt.Errorf("line %d: shamt 解析失败: %v", it.LineNo, err)
		}
		word = (uint32(rt) << 16) | (uint32(rd) << 11) | (uint32(shamt&0x1f) << 6) | instr.Funct
	// Format: op rs
	case "jr", "mthi", "mtlo":
		if len(toks) < 2 {
			return nil, fmt.Errorf("line %d: %s 需要 1 个寄存器操作数", it.LineNo, op)
		}
		rs := regs.RegOf(toks[1], it.LineNo)
		word = (uint32(rs) << 21) | instr.Funct
	// Format: op rd
	case "mfhi", "mflo":
		if len(toks) < 2 {
			return nil, fmt.Errorf("line %d: %s 需要 1 个寄存器操作数", it.LineNo, op)
		}
		rd := regs.RegOf(toks[1], it.LineNo)
		word = (uint32(rd) << 11) | instr.Funct
	// Format: op rs, rt
	case "mult", "multu", "div", "divu":
		if len(toks) < 3 {
			return nil, fmt.Errorf("line %d: %s 需要 2 个寄存器操作数", it.LineNo, op)
		}
		rs := regs.RegOf(toks[1], it.LineNo)
		rt := regs.RegOf(toks[2], it.LineNo)
		word = (uint32(rs) << 21) | (uint32(rt) << 16) | instr.Funct
	// Format: op rd, rs
	case "jalr":
		if len(toks) < 2 { // 支持 jalr $rs 和 jalr $rd, $rs 两种格式
			return nil, fmt.Errorf("line %d: %s 至少需要一个操作数", it.LineNo, op)
		}
		var rd, rs int
		if len(toks) == 2 { // jalr $rs -> $ra is implicitly $rd (31)
			rd = 31
			rs = regs.RegOf(toks[1], it.LineNo)
		} else { // jalr $rd, $rs
			rd = regs.RegOf(toks[1], it.LineNo)
			rs = regs.RegOf(toks[2], it.LineNo)
		}
		word = (uint32(rs) << 21) | (uint32(rd) << 11) | instr.Funct
	default:
		return nil, fmt.Errorf("line %d: 不支持的R类型指令: %s", it.LineNo, op)
	}
	return []uint32{word}, nil
}

func assembleIType(it types.Item, instr types.Instruction, labels map[string]uint32, addr uint32, base uint32) ([]uint32, error) {
	toks := it.Tokens
	op := strings.ToLower(toks[0])
	var word uint32

	switch op {
	// Format: op rt, rs, imm
	case "addi", "addiu", "andi", "ori", "xori", "slti", "sltiu":
		if len(toks) < 4 {
			return nil, fmt.Errorf("line %d: %s 需要 3 个操作数", it.LineNo, op)
		}
		rt := regs.RegOf(toks[1], it.LineNo)
		rs := regs.RegOf(toks[2], it.LineNo)
		imm, err := parseNumber(toks[3], labels, base)
		if err != nil {
			return nil, fmt.Errorf("line %d: 解析立即数失败: %v", it.LineNo, err)
		}
		word = (instr.Opcode << 26) | (uint32(rs) << 21) | (uint32(rt) << 16) | (uint32(imm) & 0xffff)
	// Format: op rt, imm
	case "lui":
		if len(toks) < 3 {
			return nil, fmt.Errorf("line %d: lui 需要 reg, imm", it.LineNo)
		}
		rt := regs.RegOf(toks[1], it.LineNo)
		imm, err := parseNumber(toks[2], labels, base)
		if err != nil {
			return nil, fmt.Errorf("line %d: lui 立即数解析失败: %v", it.LineNo, err)
		}
		word = (instr.Opcode << 26) | (uint32(rt) << 16) | (imm & 0xffff)
	// Format: op rt, offset(base)
	case "lw", "lh", "lhu", "lb", "lbu", "sw", "sh", "sb":
		if len(toks) < 3 {
			return nil, fmt.Errorf("line %d: %s 需要 2 个操作数", it.LineNo, op)
		}
		rt := regs.RegOf(toks[1], it.LineNo)
		off, baseReg, err := parseOffsetBase(toks[2], labels, base)
		if err != nil {
			return nil, fmt.Errorf("line %d: 解析 offset(base) 失败: %v", it.LineNo, err)
		}
		word = (instr.Opcode << 26) | (uint32(baseReg) << 21) | (uint32(rt) << 16) | (uint32(off) & 0xffff)
	// Format: op rs, rt, label
	case "beq", "bne":
		if len(toks) < 4 {
			return nil, fmt.Errorf("line %d: %s 需要 3 个操作数", it.LineNo, op)
		}
		rs := regs.RegOf(toks[1], it.LineNo)
		rt := regs.RegOf(toks[2], it.LineNo)
		off, err := branchOffset(toks[3], labels, addr, base)
		if err != nil {
			return nil, fmt.Errorf("line %d: 分支目标解析失败: %v", it.LineNo, err)
		}
		word = (instr.Opcode << 26) | (uint32(rs) << 21) | (uint32(rt) << 16) | (uint32(off) & 0xffff)
	// Format: op rs, label
	case "blez", "bgtz":
		if len(toks) < 3 {
			return nil, fmt.Errorf("line %d: %s 需要 2 个操作数", it.LineNo, op)
		}
		rs := regs.RegOf(toks[1], it.LineNo)
		off, err := branchOffset(toks[2], labels, addr, base)
		if err != nil {
			return nil, fmt.Errorf("line %d: 分支目标解析失败: %v", it.LineNo, err)
		}
		word = (instr.Opcode << 26) | (uint32(rs) << 21) | (uint32(off) & 0xffff)
	default:
		return nil, fmt.Errorf("line %d: 不支持的I类型指令: %s", it.LineNo, op)
	}
	return []uint32{word}, nil
}

func assembleJType(it types.Item, instr types.Instruction, labels map[string]uint32, base uint32) ([]uint32, error) {
	toks := it.Tokens
	if len(toks) < 2 {
		return nil, fmt.Errorf("line %d: %s 需要目标标签或地址", it.LineNo, toks[0])
	}
	targetTok := toks[1]
	targetAddr, ok := labels[targetTok]
	var targetAbs uint32
	if ok {
		targetAbs = base + targetAddr
	} else {
		v, err := parseNumber(targetTok, labels, base)
		if err != nil {
			return nil, fmt.Errorf("line %d: 未知跳转目标 %s", it.LineNo, targetTok)
		}
		targetAbs = v
	}
	field := (targetAbs >> 2) & 0x03ffffff
	word := (instr.Opcode << 26) | field
	return []uint32{word}, nil
}

func assembleSpecial(it types.Item, labels map[string]uint32, base uint32) ([]uint32, error) {
	toks := it.Tokens
	op := strings.ToLower(toks[0])

	switch op {
	case "nop":
		return []uint32{0x00000000}, nil
	case "li":
		// li rd, imm  --> lui at, imm>>16 ; ori rd, at, imm&0xFFFF
		if len(toks) < 3 {
			return nil, fmt.Errorf("line %d: li 需要两个操作数", it.LineNo)
		}
		rd := regs.RegOf(toks[1], it.LineNo)
		imm32, err := parseNumber(toks[2], labels, base)
		if err != nil {
			return nil, fmt.Errorf("line %d: li 立即数解析失败: %v", it.LineNo, err)
		}
		hi := (imm32 >> 16) & 0xffff
		lo := imm32 & 0xffff
		// lui at, hi
		lui := (instrTable["lui"].Opcode << 26) | (uint32(1) << 16) | hi
		// ori rd, at, lo
		ori := (instrTable["ori"].Opcode << 26) | (uint32(1) << 21) | (uint32(rd) << 16) | lo
		return []uint32{lui, ori}, nil
	}
	return nil, fmt.Errorf("line %d: 未知特殊指令 %s", it.LineNo, op)
}

func parseNumber(s string, labels map[string]uint32, base uint32) (uint32, error) {
	s = strings.TrimSpace(s)
	if v, ok := labels[s]; ok {
		return base + v, nil
	}
	// hex
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		v, err := strconv.ParseUint(s[2:], 16, 32)
		return uint32(v), err
	}
	// negative?
	if strings.HasPrefix(s, "-") {
		v, err := strconv.ParseInt(s, 10, 32)
		return uint32(v), err
	}
	// decimal
	v, err := strconv.ParseUint(s, 10, 32)
	return uint32(v), err
}

func parseOffsetBase(s string, labels map[string]uint32, base uint32) (int32, int, error) {
	// 期望形如: 4($t0) 或 label($t0)
	s = strings.TrimSpace(s)
	if i := strings.Index(s, "("); i >= 0 {
		j := strings.Index(s, ")")
		if j < 0 {
			return 0, 0, fmt.Errorf("缺少右括号")
		}
		offStr := strings.TrimSpace(s[:i])
		baseStr := strings.TrimSpace(s[i+1 : j])
		baseReg := regs.RegOf(baseStr, 0)
		if offStr == "" {
			return 0, baseReg, nil
		}
		if v, ok := labels[offStr]; ok {
			return int32(base + v), baseReg, nil
		}
		// 支持 0x.. 或 十进制或负数
		if strings.HasPrefix(offStr, "0x") || strings.HasPrefix(offStr, "0X") {
			val, err := strconv.ParseInt(offStr[2:], 16, 32)
			return int32(val), baseReg, err
		}
		val, err := strconv.ParseInt(offStr, 10, 32)
		return int32(val), baseReg, err
	}
	return 0, 0, fmt.Errorf("offset(base) 形式期望，但收到: %s", s)
}

func branchOffset(target string, labels map[string]uint32, curAddr uint32, base uint32) (int32, error) {
	// offset = (labelAddr - (curAddr + 4)) / 4
	target = strings.TrimSpace(target)
	if v, ok := labels[target]; ok {
		curAbs := int32(base + curAddr)
		targetAbs := int32(base + v)
		offset := targetAbs - (curAbs + 4)
		return offset / 4, nil
	}
	// numeric
	if strings.HasPrefix(target, "0x") || strings.HasPrefix(target, "0X") {
		val, err := strconv.ParseInt(target[2:], 16, 16)
		return int32(val), err
	}
	val, err := strconv.ParseInt(target, 10, 32)
	if err == nil {
		return int32(val), nil
	}
	return 0, fmt.Errorf("未知分支目标: %s", target)
}
