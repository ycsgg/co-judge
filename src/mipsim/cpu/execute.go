package cpu

func (c *CPU) Execute(instrHex uint32) ExecResult {
	opcode := instrHex >> 26
	rs := (instrHex >> 21) & 0x1F
	rt := (instrHex >> 16) & 0x1F
	rd := (instrHex >> 11) & 0x1F
	shamt := (instrHex >> 6) & 0x1F
	funct := instrHex & 0x3F
	imm := instrHex & 0xFFFF
	immSe := signExtend16(imm)
	addr := instrHex & 0x03FFFFFF

	res := ExecResult{}

	defer func() { c.Regs[0] = 0 }()

	switch opcode {
	case 0x00: // R-type
		switch funct {
		case 0x20: // add
			c.Regs[rd] = uint32(int32(c.Regs[rs]) + int32(c.Regs[rt]))
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x21: // addu
			c.Regs[rd] = c.Regs[rs] + c.Regs[rt]
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x22: // sub
			c.Regs[rd] = uint32(int32(c.Regs[rs]) - int32(c.Regs[rt]))
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x23: // subu
			c.Regs[rd] = c.Regs[rs] - c.Regs[rt]
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x24: // and
			c.Regs[rd] = c.Regs[rs] & c.Regs[rt]
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x25: // or
			c.Regs[rd] = c.Regs[rs] | c.Regs[rt]
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x26: // xor
			c.Regs[rd] = c.Regs[rs] ^ c.Regs[rt]
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x27: // nor
			c.Regs[rd] = ^(c.Regs[rs] | c.Regs[rt])
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x00: // sll
			c.Regs[rd] = c.Regs[rt] << shamt
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x02: // srl
			c.Regs[rd] = c.Regs[rt] >> shamt
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x03: // sra
			c.Regs[rd] = uint32(int32(c.Regs[rt]) >> shamt)
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x04: // sllv
			c.Regs[rd] = c.Regs[rt] << (c.Regs[rs] & 0x1F)
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x06: // srlv
			c.Regs[rd] = c.Regs[rt] >> (c.Regs[rs] & 0x1F)
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x07: // srav
			c.Regs[rd] = uint32(int32(c.Regs[rt]) >> (c.Regs[rs] & 0x1F))
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x08: // jr
			c.NextPC = c.Regs[rs]
		case 0x09: // jalr
			c.Regs[rd] = c.PC + 8
			c.NextPC = c.Regs[rs]
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x10: // mfhi
			c.Regs[rd] = c.Hi
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x11: // mthi
			c.Hi = c.Regs[rs]
		case 0x12: // mflo
			c.Regs[rd] = c.Lo
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x13: // mtlo
			c.Lo = c.Regs[rs]
		case 0x18: // mult
			val := int64(int32(c.Regs[rs])) * int64(int32(c.Regs[rt]))
			c.Hi = uint32(val >> 32)
			c.Lo = uint32(val & 0xFFFFFFFF)
		case 0x19: // multu
			val := uint64(c.Regs[rs]) * uint64(c.Regs[rt])
			c.Hi = uint32(val >> 32)
			c.Lo = uint32(val & 0xFFFFFFFF)
		case 0x1A: // div
			if c.Regs[rt] != 0 {
				c.Lo = uint32(int32(c.Regs[rs]) / int32(c.Regs[rt]))
				c.Hi = uint32(int32(c.Regs[rs]) % int32(c.Regs[rt]))
			}
		case 0x1B: // divu
			if c.Regs[rt] != 0 {
				c.Lo = c.Regs[rs] / c.Regs[rt]
				c.Hi = c.Regs[rs] % c.Regs[rt]
			}
		case 0x2a: // slt
			if int32(c.Regs[rs]) < int32(c.Regs[rt]) {
				c.Regs[rd] = 1
			} else {
				c.Regs[rd] = 0
			}
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		case 0x2b: // sltu
			if c.Regs[rs] < c.Regs[rt] {
				c.Regs[rd] = 1
			} else {
				c.Regs[rd] = 0
			}
			if rd != 0 {
				res.RegWrite = true
				res.RegDest = rd
				res.RegWriteData = c.Regs[rd]
			}
		}
	case 0x02: // j
		c.NextPC = (c.PC & 0xF0000000) | (addr << 2)
	case 0x03: // jal
		c.Regs[31] = c.PC + 4
		c.NextPC = (c.PC & 0xF0000000) | (addr << 2)
		res.RegWrite = true
		res.RegDest = 31
		res.RegWriteData = c.Regs[31]
	case 0x08: // addi
		c.Regs[rt] = uint32(int32(c.Regs[rs]) + int32(immSe))
		if rt != 0 {
			res.RegWrite = true
			res.RegDest = rt
			res.RegWriteData = c.Regs[rt]
		}
	case 0x09: // addiu
		c.Regs[rt] = c.Regs[rs] + immSe
		if rt != 0 {
			res.RegWrite = true
			res.RegDest = rt
			res.RegWriteData = c.Regs[rt]
		}
	case 0x0C: // andi
		c.Regs[rt] = c.Regs[rs] & imm
		if rt != 0 {
			res.RegWrite = true
			res.RegDest = rt
			res.RegWriteData = c.Regs[rt]
		}
	case 0x0D: // ori
		c.Regs[rt] = c.Regs[rs] | imm
		if rt != 0 {
			res.RegWrite = true
			res.RegDest = rt
			res.RegWriteData = c.Regs[rt]
		}
	case 0x0E: // xori
		c.Regs[rt] = c.Regs[rs] ^ imm
		if rt != 0 {
			res.RegWrite = true
			res.RegDest = rt
			res.RegWriteData = c.Regs[rt]
		}
	case 0x0F: // lui
		c.Regs[rt] = imm << 16
		if rt != 0 {
			res.RegWrite = true
			res.RegDest = rt
			res.RegWriteData = c.Regs[rt]
		}
	case 0x04: // beq
		if c.Regs[rs] == c.Regs[rt] {
			c.NextPC = c.PC + 4 + (immSe << 2)
		}
	case 0x05: // bne
		if c.Regs[rs] != c.Regs[rt] {
			c.NextPC = c.PC + 4 + (immSe << 2)
		}
	case 0x06: // blez
		if int32(c.Regs[rs]) <= 0 {
			c.NextPC = c.PC + 4 + (immSe << 2)
		}
	case 0x07: // bgtz
		if int32(c.Regs[rs]) > 0 {
			c.NextPC = c.PC + 4 + (immSe << 2)
		}
	case 0x01: // REGIMM
		switch rt {
		case 0x00: // bltz
			if int32(c.Regs[rs]) < 0 {
				c.NextPC = c.PC + 4 + (immSe << 2)
			}
		case 0x01: // bgez
			if int32(c.Regs[rs]) >= 0 {
				c.NextPC = c.PC + 4 + (immSe << 2)
			}
		}
	case 0x20: // lb
		addrCalc := c.Regs[rs] + immSe
		val := c.Mem[addrCalc&^3] >> ((3 - (addrCalc % 4)) * 8) & 0xFF
		if val&0x80 != 0 {
			val |= 0xFFFFFF00
		}
		c.Regs[rt] = val
		if rt != 0 {
			res.RegWrite = true
			res.RegDest = rt
			res.RegWriteData = c.Regs[rt]
		}
	case 0x21: // lh
		addrCalc := c.Regs[rs] + immSe
		val := c.Mem[addrCalc&^3] >> ((2 - (addrCalc % 4)) * 8) & 0xFFFF
		if val&0x8000 != 0 {
			val |= 0xFFFF0000
		}
		c.Regs[rt] = val
		if rt != 0 {
			res.RegWrite = true
			res.RegDest = rt
			res.RegWriteData = c.Regs[rt]
		}
	case 0x23: // lw
		addrCalc := c.Regs[rs] + immSe
		c.Regs[rt] = c.Mem[addrCalc]
		if rt != 0 {
			res.RegWrite = true
			res.RegDest = rt
			res.RegWriteData = c.Regs[rt]
		}
	case 0x24: // lbu
		addrCalc := c.Regs[rs] + immSe
		c.Regs[rt] = c.Mem[addrCalc&^3] >> ((3 - (addrCalc % 4)) * 8) & 0xFF
		if rt != 0 {
			res.RegWrite = true
			res.RegDest = rt
			res.RegWriteData = c.Regs[rt]
		}
	case 0x25: // lhu
		addrCalc := c.Regs[rs] + immSe
		c.Regs[rt] = c.Mem[addrCalc&^3] >> ((2 - (addrCalc % 4)) * 8) & 0xFFFF
		if rt != 0 {
			res.RegWrite = true
			res.RegDest = rt
			res.RegWriteData = c.Regs[rt]
		}
	case 0x28: // sb
		addrCalc := c.Regs[rs] + immSe
		shift := (3 - (addrCalc % 4)) * 8
		mask := ^(uint32(0xFF) << shift)
		c.Mem[addrCalc&^3] = (c.Mem[addrCalc&^3] & mask) | ((c.Regs[rt] & 0xFF) << shift)
		res.MemWrite = true
		res.MemDest = addrCalc
		res.MemWriteData = c.Regs[rt] & 0xFF
	case 0x29: // sh
		addrCalc := c.Regs[rs] + immSe
		shift := (2 - (addrCalc % 4)) * 8
		mask := ^(uint32(0xFFFF) << shift)
		c.Mem[addrCalc&^3] = (c.Mem[addrCalc&^3] & mask) | ((c.Regs[rt] & 0xFFFF) << shift)
		res.MemWrite = true
		res.MemDest = addrCalc
		res.MemWriteData = c.Regs[rt] & 0xFFFF
	case 0x2B: // sw
		addrCalc := c.Regs[rs] + immSe
		c.Mem[addrCalc] = c.Regs[rt]
		res.MemWrite = true
		res.MemDest = addrCalc
		res.MemWriteData = c.Regs[rt]
	case 0x0A: // slti
		if int32(c.Regs[rs]) < int32(immSe) {
			c.Regs[rt] = 1
		} else {
			c.Regs[rt] = 0
		}
		if rt != 0 {
			res.RegWrite = true
			res.RegDest = rt
			res.RegWriteData = c.Regs[rt]
		}
	case 0x0B: // sltiu
		if c.Regs[rs] < immSe {
			c.Regs[rt] = 1
		} else {
			c.Regs[rt] = 0
		}
		if rt != 0 {
			res.RegWrite = true
			res.RegDest = rt
			res.RegWriteData = c.Regs[rt]
		}
	}

	return res
}
