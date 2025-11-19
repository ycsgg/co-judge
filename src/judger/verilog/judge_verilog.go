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
			if strings.Contains(parts[1], "<=") {
				subParts := strings.Split(parts[1], "<=")
				if len(subParts) != 2 {
					continue
				}
				left := strings.TrimSpace(subParts[0])
				right := strings.TrimSpace(subParts[1])
				if strings.HasPrefix(left, "$") {
					regDestStr := strings.TrimPrefix(left, "$")
					regDest, err := strconv.ParseUint(regDestStr, 10, 32)
					if err != nil {
						continue
					}
					regData, err := strconv.ParseUint(right, 16, 32)
					if err != nil {
						continue
					}
					res = append(res, MipsLine{
						PC:       uint32(pc),
						RegWrite: true,
						RegDest:  uint32(regDest),
						RegData:  uint32(regData),
					})
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
					res = append(res, MipsLine{
						PC:       uint32(pc),
						MemWrite: true,
						MemAddr:  uint32(memAddr),
						MemData:  uint32(memData),
					})
				}
			}
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
	var varilogTraceIdx, mipsimTraceIdx int
	for varilogTraceIdx < len(verilogTrace) && mipsimTraceIdx < len(mipsimTrace) {
		v := verilogTrace[varilogTraceIdx]
		m := mipsimTrace[mipsimTraceIdx]

		if !m.RegWrite && !m.MemWrite {
			mipsimTraceIdx++
			continue
		}

		if v.PC != m.PC {
			diffs = append(diffs, "PC mismatch at "+strconv.FormatUint(uint64(v.PC), 16))
		}
		if v.RegWrite != m.RegWrite {
			diffs = append(diffs, "RegWrite mismatch at PC "+strconv.FormatUint(uint64(v.PC), 16))
		} else if v.RegWrite && m.RegWrite {
			if v.RegDest != m.RegDest {
				diffs = append(diffs, "RegDest mismatch at PC "+strconv.FormatUint(uint64(v.PC), 16))
			}
			if v.RegData != m.RegData {
				diffs = append(diffs, "RegData mismatch at PC "+strconv.FormatUint(uint64(v.PC), 16))
			}
		}
		if v.MemWrite != m.MemWrite {
			diffs = append(diffs, "MemWrite mismatch at PC "+strconv.FormatUint(uint64(v.PC), 16))
		} else if v.MemWrite && m.MemWrite {
			if v.MemAddr != m.MemAddr {
				diffs = append(diffs, "MemAddr mismatch at PC "+strconv.FormatUint(uint64(v.PC), 16))
			}
			if v.MemData != m.MemData {
				diffs = append(diffs, "MemData mismatch at PC "+strconv.FormatUint(uint64(v.PC), 16))
			}
		}

		varilogTraceIdx++
		mipsimTraceIdx++
	}
	return diffs
}
