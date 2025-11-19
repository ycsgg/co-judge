package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"mips2hex/assembler"
	"mips2hex/parser"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

func TestEndToEndAssembler(t *testing.T) {
	testDir := "test"

	for i := 1; i <= 10; i++ {
		testName := fmt.Sprintf("code%d.s", i)
		t.Run(testName, func(t *testing.T) {
			inputPath := filepath.Join(testDir, fmt.Sprintf("code%d.s", i))
			expectedOutputPath := filepath.Join(testDir, fmt.Sprintf("instr%d.txt", i))

			expectedLines, err := readLines(expectedOutputPath)
			if err != nil {
				t.Fatalf("无法读取期望的结果文件 %s: %v", expectedOutputPath, err)
			}

			lines, err := parser.ReadFileLines(inputPath)
			if err != nil {
				t.Fatalf("无法读取输入文件 %s: %v", inputPath, err)
			}
			items, labels, err := parser.ParseLines(lines)
			if err != nil {
				t.Fatalf("解析文件 %s 失败: %v", inputPath, err)
			}

			machineCode, err := assembler.Assemble(items, labels, 0x00003000)
			if err != nil {
				t.Fatalf("汇编文件 %s 失败: %v", inputPath, err)
			}

			var actualLines []string
			for _, code := range machineCode {
				actualLines = append(actualLines, fmt.Sprintf("%08x", code))
			}

			if len(actualLines) != len(expectedLines) {
				t.Fatalf("输出的行数不匹配: 期望 %d, 实际 %d", len(expectedLines), len(actualLines))
			}

			for lineIdx := 0; lineIdx < len(actualLines); lineIdx++ {
				actualLine := actualLines[lineIdx]
				expectedLine := expectedLines[lineIdx]

				if actualLine != expectedLine {
					t.Errorf("在第 %d 行发现不匹配:\n期望: %s\n实际: %s", lineIdx+1, expectedLine, actualLine)
					t.FailNow()
				}
			}
		})
	}
}
