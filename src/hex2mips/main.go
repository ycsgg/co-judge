package main

import (
	"bufio"
	"flag"
	"fmt"
	"hex2mips/disassembler"
	"os"
	"strconv"
	"strings"
)

func parseWord(token string) (uint32, error) {
	tok := strings.TrimSpace(token)
	if tok == "" {
		return 0, fmt.Errorf("empty token")
	}

	// Binary
	if len(tok) == 32 && strings.Trim(tok, "01") == "" {
		v, err := strconv.ParseUint(tok, 2, 32)
		return uint32(v), err
	}

	// Hex
	if strings.HasPrefix(tok, "0x") || strings.HasPrefix(tok, "0X") {
		v, err := strconv.ParseUint(tok[2:], 16, 32)
		return uint32(v), err
	}
	if len(tok) <= 8 && strings.Trim(tok, "0123456789abcdefABCDEF") == "" {
		v, err := strconv.ParseUint(tok, 16, 32)
		return uint32(v), err
	}

	// Decimal
	v, err := strconv.ParseInt(tok, 10, 32)
	if err == nil {
		return uint32(v), nil
	}

	return 0, fmt.Errorf("cannot parse token: %s", tok)
}

func main() {
	baseAddr := flag.String("base", "0x00400000", "Base address for PC (in hex)")
	flag.Parse()

	base, err := strconv.ParseUint(*baseAddr, 0, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing base address: %v\n", err)
		os.Exit(1)
	}
	pc := uint32(base)

	var scanner *bufio.Scanner
	args := flag.Args()
	if len(args) > 0 {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", args[0], err)
			os.Exit(1)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	} else {
		if fi, _ := os.Stdin.Stat(); (fi.Mode() & os.ModeCharDevice) != 0 {
			fmt.Println("Reading machine code from stdin, one instruction per line (hex/bin/dec). Press Ctrl-D to end.")
		}
		scanner = bufio.NewScanner(os.Stdin)
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			if line != "" {
				fmt.Println(line)
			}
			continue
		}

		tokens := strings.Fields(line)
		for _, tok := range tokens {
			word, err := parseWord(tok)
			if err != nil {
				fmt.Printf("0x%08x: ERROR parsing token '%s': %v\n", pc, tok, err)
				pc += 4
				continue
			}
			asm := disassembler.DecodeWord(word, pc)
			fmt.Printf("0x%08x: 0x%08x\t%s\n", pc, word, asm)
			pc += 4
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}
