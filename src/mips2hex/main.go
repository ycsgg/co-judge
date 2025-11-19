package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"mips2hex/assembler"
	"mips2hex/emitter"
	"mips2hex/parser"
)

func main() {
	inPath := flag.String("input", "", "输入 MIPS asm 文件路径")
	outPath := flag.String("output", "", "输出 hex 文件路径")
	baseStr := flag.String("base", "0", ".text 段基址(0x前缀十六进制),用于分支/跳转计算")
	flag.Parse()

	if *inPath == "" || *outPath == "" {
		fmt.Fprintln(os.Stderr, "用法: go run main.go -input code.s -output instr.txt")
		os.Exit(2)
	}

	lines, err := parser.ReadFileLines(*inPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "读取输入文件失败:", err)
		os.Exit(1)
	}

	// 解析基址
	baseVal, err := strconv.ParseUint(*baseStr, 16, 32)
	if err != nil {
		fmt.Fprintln(os.Stderr, "解析 base 失败:", err)
		os.Exit(2)
	}

	items, labels, err := parser.ParseLines(lines)
	if err != nil {
		fmt.Fprintln(os.Stderr, "解析失败:", err)
		os.Exit(1)
	}

	words, err := assembler.Assemble(items, labels, uint32(baseVal))
	if err != nil {
		fmt.Fprintln(os.Stderr, "汇编失败:", err)
		os.Exit(1)
	}

	if err := emitter.WriteHexLines(*outPath, words); err != nil {
		fmt.Fprintln(os.Stderr, "写入失败:", err)
		os.Exit(1)
	}

	fmt.Printf("完成：写入 %d 条指令到 %s\n", len(words), *outPath)
}
