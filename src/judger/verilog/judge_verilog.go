package verilog

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
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

func JudgeVerilog(isePath, verilogPath, prjPath, tbPath, tclPath, hexPath string) (JudgeResult, error) {
	err := loadCode(verilogPath, hexPath)
	if err != nil {
		return JudgeResult{}, fmt.Errorf("load code failed: %w", err)
	}
	verilogOut, err := runVerilog(isePath, verilogPath, prjPath, tbPath, tclPath)
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

func loadCode(verilogPath, codePath string) error {
	data, err := os.ReadFile(codePath)
	if err != nil {
		return err
	}
	err = os.WriteFile(verilogPath+string(os.PathSeparator)+"code.txt", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func runVerilog(isePath, verilogPath, prjPath, tbPath, tclPath string) (string, error) {
	cmd1 := exec.Command("fuse", "-nodebug", "-prj", prjPath, "-o", "judge.exe", tbPath)
	cmd1.Dir = verilogPath
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd1.Stdout = &errOut
	cmd1.Stderr = &errOut
	err := cmd1.Run()
	if err != nil {
		return errOut.String(), fmt.Errorf("fuse command failed: %w\nOutput:\n%s", err, errOut.String())
	}
	cmd2 := exec.Command("./judge.exe", "-nolog", "-tclbatch", tclPath)
	cmd2.Stdout = &out
	cmd2.Stderr = &out
	cmd2.Dir = verilogPath
	err = cmd2.Run()
	cmd3 := exec.Command("rm", "judge.exe")
	cmd3.Dir = verilogPath
	cmd3.Run()
	if err != nil {
		return out.String(), fmt.Errorf("verilog simulation failed: %w\nOutput:\n%s", err, out.String())
	}
	return out.String(), err
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
	for step := 1; ; step++ {
		if step > c.MaxSteps { // 使用 CPU 默认限制，避免评测卡死
			break
		}
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
