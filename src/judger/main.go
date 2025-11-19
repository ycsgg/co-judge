package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"judger/logisim"
	"judger/verilog"
)

func main() {
	mode := flag.String("mode", "", "mode: logisim,verilog")
	flag.Parse()

	args := flag.Args()
	if *mode == "logisim" {
		if len(args) < 4 {
			fmt.Fprintln(os.Stderr, "Usage: judger -mode logisim <jar_path> <circ_path> <hex_path> <output_path>")
			os.Exit(2)
		}
		jar := args[0]
		circ := args[1]
		hex := args[2]
		out := args[3]

		res, err := logisim.JudgeLogisim(jar, circ, hex)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Judge error:", err)
			os.Exit(1)
		}
		// Write result lines to output
		if err := os.WriteFile(out, []byte(strings.Join(res.Diffs, "\n")+"\n"), 0644); err != nil {
			fmt.Fprintln(os.Stderr, "write output error:", err)
			os.Exit(1)
		}
		if res.OK {
			fmt.Println("All lines OK")
			os.Exit(0)
		}
		// On mismatch, also write mipsim details to detail.log next to output file
		detailPath := filepath.Join(filepath.Dir(out), "detail.log")
		var b strings.Builder
		for i, ml := range res.MipsLines {
			fmt.Fprintf(&b, "line %d: Instr=0x%08x PC=0x%08x RegWrite=%v RegDest=%d RegData=0x%08x MemWrite=%v MemAddr=0x%08x MemData=0x%08x\n",
				i+1, ml.Instr, ml.PC, ml.RegWrite, ml.RegDest, ml.RegData, ml.MemWrite, ml.MemAddr, ml.MemData)
		}
		if err := os.WriteFile(detailPath, []byte(b.String()), 0644); err != nil {
			fmt.Fprintln(os.Stderr, "write detail.log error:", err)
		}
		fmt.Println("Mismatch found. See output and detail.log for details.")
		os.Exit(1)
	} else if *mode == "verilog" {
		if len(args) < 4 {
			fmt.Fprintln(os.Stderr, "Usage: judger -mode verilog <ise_path> <verilog_path> <hex_path> <output_path>")
			os.Exit(2)
		}
		ise := args[0]
		vfile := args[1]
		hex := args[2]
		out := args[3]

		res, err := verilog.JudgeVerilog(ise, vfile, hex)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Judge error:", err)
			os.Exit(1)
		}
		if err := os.WriteFile(out, []byte(strings.Join(res.Diffs, "\n")+"\n"), 0644); err != nil {
			fmt.Fprintln(os.Stderr, "write output error:", err)
			os.Exit(1)
		}
		if res.OK {
			fmt.Println("All lines OK")
			os.Exit(0)
		}
		fmt.Println("Mismatch found. See output for details.")
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Unknown or missing -mode. Use -mode logisim or -mode verilog")
	os.Exit(2)
}
