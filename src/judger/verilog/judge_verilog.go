package verilog

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"mipsim/cpu"
)

type MipsLine struct {
	Instr    uint32
	PC       uint32
	RegWrite bool
	RegDest  uint32
	RegData  uint32
	MemWrite bool
	MemAddr  uint32
	MemData  uint32
}

type JudgeResult struct {
	OK        bool
	Diffs     []string
	MipsLines []MipsLine
}

func JudgeVerilog(isePath, verilogPath, hexPath string) (JudgeResult, error) {
	verilogOut, err := runVerilog(isePath, verilogPath)
	if err != nil {
		return JudgeResult{}, fmt.Errorf("run verilog failed: %w", err)
	}
	veriLines := parseVerilogOutput(verilogOut)
	if len(veriLines) == 0 {
		return JudgeResult{}, errors.New("no valid verilog trace lines parsed")
	}

	mipsLines, err := runMipsimLocal(hexPath)
	if err != nil {
		return JudgeResult{}, fmt.Errorf("run mipsim(cpu) failed: %w", err)
	}

	diffs := compareTrace(veriLines, mipsLines)

	hasError := false
	for _, d := range diffs {
		if strings.Contains(d, "mismatch") {
			hasError = true
			break
		}
	}

	return JudgeResult{
		OK:        !hasError,
		Diffs:     diffs,
		MipsLines: mipsLines,
	}, nil
}

func runVerilog(isePath, verilogPath string) (string, error) {
	var res strings.Builder
	fmt.Printf("Input verilogOutput (-1 as Ended)\n")
	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		line := reader.Text()
		if line == "-1" {
			break
		}
		res.WriteString(line + "\n")
	}
	return res.String(), nil
}

func parseVerilogOutput(out string) []MipsLine {
	var res []MipsLine
	scanner := bufio.NewScanner(strings.NewReader(out))
	/***
	always @(posedge clk) begin
		if (RegWrite && reset != 1) begin
			$display("@%h: $%d <= %h", pc, rd, grfdin);
		end
		if (MemWrite && reset != 1) begin
			$display("@%h: *%h <= %h", pc, aluout, dout2);
		end
	end
	***/
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "@") {
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				continue
			}
			pcStr := strings.TrimPrefix(parts[0], "@")
			pc, err := strconv.ParseUint(pcStr, 16, 32)
			if err != nil {
				continue
			}
			var ml MipsLine
			if strings.Contains(parts[1], "<=") {
				subParts := strings.Split(parts[1], "<=")
				if len(subParts) != 2 {
					continue
				}
				left := strings.TrimSpace(subParts[0])
				right := strings.TrimSpace(subParts[1])
				if strings.HasPrefix(left, "$") {
					regDestStr := strings.TrimPrefix(strings.TrimPrefix(left, "$"), " ")
					regDest, err := strconv.ParseUint(regDestStr, 10, 32)
					if err != nil {
						continue
					}
					regData, err := strconv.ParseUint(right, 16, 32)
					if err != nil {
						continue
					}
					ml = MipsLine{
						PC:       uint32(pc),
						RegWrite: true,
						RegDest:  uint32(regDest),
						RegData:  uint32(regData),
					}
				} else if strings.HasPrefix(left, "*") {
					memAddrStr := strings.TrimPrefix(left, "*")
					memAddr, err := strconv.ParseUint(memAddrStr, 16, 32)
					if err != nil {
						continue
					}
					memData, err := strconv.ParseUint(right, 16, 32)
					if err != nil {
						continue
					}
					ml = MipsLine{
						PC:       uint32(pc),
						MemWrite: true,
						MemAddr:  uint32(memAddr),
						MemData:  uint32(memData),
					}
				}
			}
			res = append(res, ml)
		}
	}
	return res
}

func runMipsimLocal(hexPath string) ([]MipsLine, error) {
	// read hex words
	f, err := os.Open(hexPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var instrs []uint32
	rd := bufio.NewReader(f)
	for {
		ln, err := rd.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
		ln = strings.TrimSpace(ln)
		if ln != "" {
			w64, perr := strconv.ParseUint(strings.TrimPrefix(ln, "0x"), 16, 32)
			if perr == nil {
				instrs = append(instrs, uint32(w64))
			}
		}
		if errors.Is(err, io.EOF) {
			break
		}
	}

	c := cpu.New()
	base := c.PC
	for i, w := range instrs {
		c.Mem[base+uint32(i*4)] = w
	}
	var out []MipsLine
	for {
		idx := (c.PC - base) / 4
		if idx >= uint32(len(instrs)) {
			break
		}
		word := c.Mem[c.PC]
		c.NextPC = c.PC + 4
		res := c.Execute(word)
		ml := MipsLine{
			Instr:    word,
			PC:       c.PC,
			RegWrite: res.RegWrite,
			RegDest:  res.RegDest,
			RegData:  res.RegWriteData,
			MemWrite: res.MemWrite,
			MemAddr:  res.MemDest,
			MemData:  res.MemWriteData,
		}
		out = append(out, ml)
		c.PC = c.NextPC
	}
	return out, nil
}

func compareTrace(verilogTrace []MipsLine, mipsimTrace []MipsLine) []string {
	var diffs []string
	var verilogTraceIdx, mipsimTraceIdx int

	for v := range verilogTrace {
		fmt.Printf("Verilog Trace Line %d: PC=0x%08x RegWrite=%v RegDest=$%d RegData=0x%08x MemWrite=%v MemAddr=0x%08x MemData=0x%08x\n",
			v+1,
			verilogTrace[v].PC,
			verilogTrace[v].RegWrite,
			verilogTrace[v].RegDest,
			verilogTrace[v].RegData,
			verilogTrace[v].MemWrite,
			verilogTrace[v].MemAddr,
			verilogTrace[v].MemData,
		)
	}

	fmt.Printf("<==============================>\n")

	for m := range mipsimTrace {
		fmt.Printf("MIPSIM Trace Line %d: PC=0x%08x RegWrite=%v RegDest=$%d RegData=0x%08x MemWrite=%v MemAddr=0x%08x MemData=0x%08x\n",
			m+1,
			mipsimTrace[m].PC,
			mipsimTrace[m].RegWrite,
			mipsimTrace[m].RegDest,
			mipsimTrace[m].RegData,
			mipsimTrace[m].MemWrite,
			mipsimTrace[m].MemAddr,
			mipsimTrace[m].MemData,
		)
	}

	for verilogTraceIdx < len(verilogTrace) && mipsimTraceIdx < len(mipsimTrace) {
		v := verilogTrace[verilogTraceIdx]
		m := mipsimTrace[mipsimTraceIdx]

		if (!m.RegWrite) && (!m.MemWrite) {
			fmt.Printf("Skipping MIPSIM line with no RegWrite and no MemWrite in PC 0x%08x\n", m.PC)
			mipsimTraceIdx++
			continue
		}

		if v.PC != m.PC {
			diffs = append(diffs, fmt.Sprintf("line %d PC mismatch: verilog=0x%08x mipsim=0x%08x", verilogTraceIdx+1, v.PC, m.PC))
		}

		if v.RegWrite != m.RegWrite {
			diffs = append(diffs, fmt.Sprintf("line %d RegWrite mismatch: verilog=%v mipsim=%v", verilogTraceIdx+1, v.RegWrite, m.RegWrite))
		} else if v.RegWrite && m.RegWrite {
			if v.RegDest != m.RegDest {
				diffs = append(diffs, fmt.Sprintf("line %d RegDest mismatch: verilog=$%d mipsim=$%d", verilogTraceIdx+1, v.RegDest, m.RegDest))
			}
			if v.RegData != m.RegData {
				diffs = append(diffs, fmt.Sprintf("line %d RegData mismatch: verilog=0x%08x mipsim=0x%08x", verilogTraceIdx+1, v.RegData, m.RegData))
			}
		}

		if v.MemWrite != m.MemWrite {
			diffs = append(diffs, fmt.Sprintf("line %d MemWrite mismatch: verilog=%v mipsim=%v", verilogTraceIdx+1, v.MemWrite, m.MemWrite))
		} else if v.MemWrite && m.MemWrite {
			if v.MemAddr != m.MemAddr {
				diffs = append(diffs, fmt.Sprintf("line %d MemAddr mismatch: verilog=0x%08x mipsim=0x%08x", verilogTraceIdx+1, v.MemAddr, m.MemAddr))
			}
			if v.MemData != m.MemData {
				diffs = append(diffs, fmt.Sprintf("line %d MemData mismatch: verilog=0x%08x mipsim=0x%08x", verilogTraceIdx+1, v.MemData, m.MemData))
			}
		}

		verilogTraceIdx++
		mipsimTraceIdx++
	}
	return diffs
}
