package logisim

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"mipsim/cpu"
)

type LogisimLine struct {
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

// JudgeLogisim injects hex program into circ, runs Logisim and mipsim, and compares traces.
// logisimJarPath: path to logisim.jar
// circPath:       path to the .circ template
// hexPath:        path to hex instruction file (one word per line; supports 0x prefix)
func JudgeLogisim(logisimJarPath, circPath, hexPath string) (JudgeResult, error) {
	circToRun := filepath.Join(filepath.Dir(circPath), "circToRun.circ")
	if err := injectHexIntoCirc(circPath, hexPath, circToRun); err != nil {
		return JudgeResult{}, fmt.Errorf("inject circ failed: %w", err)
	}

	logisimOut, err := runLogisim(logisimJarPath, circToRun)
	if err != nil {
		return JudgeResult{}, fmt.Errorf("run logisim failed: %w", err)
	}
	logiLines := parseLogisimOutput(logisimOut)
	if len(logiLines) == 0 {
		return JudgeResult{}, errors.New("no valid logisim trace lines parsed (expect 167-bit lines)")
	}

	mipsLines, err := runMipsimLocal(hexPath)
	if err != nil {
		return JudgeResult{}, fmt.Errorf("run mipsim(cpu) failed: %w", err)
	}

	diffs := compareTraces(logiLines, mipsLines)
	// Determine OK: check if any mismatch or length error exists
	hasError := false
	for _, d := range diffs {
		if strings.Contains(d, "mismatch") || strings.Contains(d, "length error") {
			hasError = true
			break
		}
	}
	return JudgeResult{OK: !hasError, Diffs: diffs, MipsLines: mipsLines}, nil
}

func injectHexIntoCirc(circPath, hexPath, outPath string) error {
	circBytes, err := os.ReadFile(circPath)
	if err != nil {
		return err
	}
	hexBytes, err := os.ReadFile(hexPath)
	if err != nil {
		return err
	}
	pattern := regexp.MustCompile(`addr/data: 12 32([\s\S]*?)</a>`)
	repl := []byte("addr/data: 12 32\n" + string(hexBytes) + "</a>")
	res := pattern.ReplaceAll(circBytes, repl)
	return os.WriteFile(outPath, res, 0644)
}

func runLogisim(jarPath, circToRun string) (string, error) {
	cmd := exec.Command("java", "-jar", jarPath, circToRun, "-tty", "table")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseLogisimOutput(out string) []LogisimLine {
	var res []LogisimLine
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		bits := make([]rune, 0, 200)
		for _, ch := range line {
			if ch == '0' || ch == '1' {
				bits = append(bits, ch)
			}
		}
		if len(bits) < 167 {
			continue
		}
		bits = bits[:167]
		get := func(a, b int) string { return string(bits[a:b]) }
		parseU := func(s string) uint32 {
			v, _ := strconv.ParseUint(s, 2, 32)
			return uint32(v)
		}
		ll := LogisimLine{
			Instr:    parseU(get(0, 32)),
			PC:       parseU(get(32, 64)),
			RegWrite: get(64, 65) == "1",
			RegDest:  parseU(get(65, 70)),
			RegData:  parseU(get(70, 102)),
			MemWrite: get(102, 103) == "1",
			MemAddr:  parseU(get(103, 135)),
			MemData:  parseU(get(135, 167)),
		}
		res = append(res, ll)
	}
	return res
}

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

func compareTraces(logis []LogisimLine, mips []MipsLine) []string {
	var diffs []string
	// Compare only up to mipsim length per requirement
	minN := len(mips)
	if len(logis) < minN {
		minN = len(logis)
	}
	for i := 0; i < minN; i++ {
		l := logis[i]
		m := mips[i]
		stepDiffs := []string{}
		if l.Instr != m.Instr {
			stepDiffs = append(stepDiffs, fmt.Sprintf("line %d Instr mismatch: logisim=0x%08x mipsim=0x%08x", i+1, l.Instr, m.Instr))
		}
		if l.PC != m.PC {
			stepDiffs = append(stepDiffs, fmt.Sprintf("line %d PC mismatch: logisim=0x%08x mipsim=0x%08x", i+1, l.PC, m.PC))
		}
		// 特判对于zero写入无意义
		if (l.RegDest == 0) && (m.RegDest == 0) {
			// 忽略数据差异
		} else {
			if l.RegWrite != m.RegWrite {
				stepDiffs = append(stepDiffs, fmt.Sprintf("line %d RegWrite mismatch: logisim=%v mipsim=%v", i+1, l.RegWrite, m.RegWrite))
			}
			if l.RegWrite && m.RegWrite {
				// 特判: 如果目标寄存器为0, 忽略数据差异
				if l.RegDest != m.RegDest {
					stepDiffs = append(stepDiffs, fmt.Sprintf("line %d RegDest mismatch: logisim=%d mipsim=%d", i+1, l.RegDest, m.RegDest))
				} else if l.RegDest != 0 { // 只有非 $zero 才比较数据
					if l.RegData != m.RegData {
						stepDiffs = append(stepDiffs, fmt.Sprintf("line %d RegData mismatch: logisim=0x%08x mipsim=0x%08x", i+1, l.RegData, m.RegData))
					}
				}
			}
		}
		if l.MemWrite != m.MemWrite {
			stepDiffs = append(stepDiffs, fmt.Sprintf("line %d MemWrite mismatch: logisim=%v mipsim=%v", i+1, l.MemWrite, m.MemWrite))
		}
		if l.MemWrite && m.MemWrite {
			if l.MemAddr != m.MemAddr {
				stepDiffs = append(stepDiffs, fmt.Sprintf("line %d MemAddr mismatch: logisim=0x%08x mipsim=0x%08x", i+1, l.MemAddr, m.MemAddr))
			}
			if l.MemData != m.MemData { // 没有忽略规则, 直接比较
				stepDiffs = append(stepDiffs, fmt.Sprintf("line %d MemData mismatch: logisim=0x%08x mipsim=0x%08x", i+1, l.MemData, m.MemData))
			}
		}
		if len(stepDiffs) == 0 {
			diffs = append(diffs, fmt.Sprintf("line %d OK", i+1))
		} else {
			diffs = append(diffs, stepDiffs...)
		}
	}
	// Length handling per requirement: follow mipsim length; if logisim has extra non-NOP lines, report length error
	if len(logis) > len(mips) {
		extraAllNOP := true
		for i := len(mips); i < len(logis); i++ {
			if logis[i].Instr != 0x00000000 {
				extraAllNOP = false
				break
			}
		}
		if !extraAllNOP {
			diffs = append(diffs, fmt.Sprintf("length error: logisim=%d mipsim=%d", len(logis), len(mips)))
		}
	}
	// If logisim is shorter than mipsim, we do not report length mismatch per instruction.
	return diffs
}
