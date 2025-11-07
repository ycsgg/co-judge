package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"mipsim/cpu"
)

func main() {
	fileFlag := flag.String("f", "", "hex instruction file (each line a 32-bit word, optionally 0x prefix)")
	flag.Parse()

	if *fileFlag == "" {
		fmt.Println("Usage: mipsim -f <hex_file>")
		os.Exit(1)
	}

	file, err := os.Open(*fileFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open file error: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var instrs []uint32
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		hexStr := strings.TrimPrefix(line, "0x")
		val, err := strconv.ParseUint(hexStr, 16, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip invalid hex '%s': %v\n", line, err)
			continue
		}
		instrs = append(instrs, uint32(val))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "read file error: %v\n", err)
		os.Exit(1)
	}

	c := cpu.New()
	c.Run(instrs)
}
