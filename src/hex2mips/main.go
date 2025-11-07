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
	// Handle hex without prefix, but be careful not to misinterpret decimals
	if len(tok) <= 8 && !strings.ContainsAny(tok, "ghijklmnopqrstuvwxyzGHIJKLMNOPQRSTUVWXYZ") {
		v, err := strconv.ParseUint(tok, 16, 32)
		if err == nil {
			// Check if it could also be a valid decimal
			_, decErr := strconv.ParseInt(tok, 10, 32)
			// If it parses as hex and not as decimal, or if it's clearly hex, prefer hex.
			// This logic is tricky. For now, we assume if it looks like hex (no non-hex chars), it is.
			// A more robust solution might require more context.
			// The original python script had a similar ambiguity.
			if decErr != nil || len(tok) > 2 { // Avoid ambiguity for short numbers like "10"
				return uint32(v), nil
			}
		}
	}
	// Fallback to trying hex for any string of 8 hex chars
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
	// Define flags
	baseAddr := flag.String("base", "0x00400000", "Base address for PC (in hex)")
	nonInteractive := flag.Bool("n", false, "Non-interactive mode, no prompts")
	inputHex := flag.String("input_hex", "", "Hex string to disassemble")
	inputHexShort := flag.String("ih", "", "Hex string to disassemble (shorthand)")
	inputFile := flag.String("input", "", "Input file path")
	inputFileShort := flag.String("i", "", "Input file path (shorthand)")

	flag.Parse()

	// Determine effective input values from long or short flags
	hexStr := *inputHex
	if *inputHexShort != "" {
		hexStr = *inputHexShort
	}
	filePath := *inputFile
	if *inputFileShort != "" {
		filePath = *inputFileShort
	}

	// Parse base address
	base, err := strconv.ParseUint(*baseAddr, 0, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing base address: %v\n", err)
		os.Exit(1)
	}
	pc := uint32(base)

	// --- Input Handling Logic ---
	// 1. --input_hex / -ih
	if hexStr != "" {
		word, err := parseWord(hexStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing hex input '%s': %v\n", hexStr, err)
			os.Exit(1)
		}
		asm := disassembler.DecodeWord(word, pc)
		fmt.Printf("0x%08x: 0x%08x\t%s\n", pc, word, asm)
		return
	}

	// 2. --input / -i or positional argument
	if filePath == "" && flag.NArg() > 0 {
		filePath = flag.Arg(0)
	}

	var scanner *bufio.Scanner
	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", filePath, err)
			os.Exit(1)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	} else {
		// 3. Stdin
		if !*nonInteractive {
			if fi, _ := os.Stdin.Stat(); (fi.Mode() & os.ModeCharDevice) != 0 {
				fmt.Println("Reading machine code from stdin, one instruction per line (hex/bin/dec). Press Ctrl-D to end.")
			}
		}
		scanner = bufio.NewScanner(os.Stdin)
	}

	// --- Processing Loop ---
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			if line != "" && !*nonInteractive {
				fmt.Println(line)
			}
			continue
		}

		tokens := strings.Fields(line)
		for _, tok := range tokens {
			word, err := parseWord(tok)
			if err != nil {
				if !*nonInteractive {
					fmt.Printf("0x%08x: ERROR parsing token '%s': %v\n", pc, tok, err)
				}
				pc += 4
				continue
			}
			asm := disassembler.DecodeWord(word, pc)
			fmt.Printf("0x%08x\t%s\n", word, asm)
			pc += 4
		}
	}

	if err := scanner.Err(); err != nil {
		if !*nonInteractive {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		}
		os.Exit(1)
	}
}
