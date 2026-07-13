// Copyright 2026 Magnobit. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/magnobit/quell/compile"
	"github.com/spf13/cobra"
)

func newCompileCmd() *cobra.Command {
	var target, outFile string
	var optimize, noOptimize bool

	cmd := &cobra.Command{
		Use:   "compile <file.quell>",
		Short: "Compile to OpenQASM, Qiskit, Cirq, or Braket",
		Example: `  quell compile bell.quell
  quell compile --target qiskit bell.quell
  quell compile --target cirq --no-optimize -o out.py bell.quell`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !strings.HasSuffix(args[0], ".quell") {
				return fmt.Errorf("expected a .quell file, got: %s", args[0])
			}

			finalOptimize := optimize
			if noOptimize {
				finalOptimize = false
			}

			// CompileFile (not CompileWithWarnings on read content) so
			// "import" lines resolve relative to this file's directory, or
			// against an installed package.
			result, err := compile.CompileFileWithWarnings(args[0], compile.Target(target), finalOptimize)
			if err != nil {
				return fmt.Errorf("compile error: %w", err)
			}

			for _, n := range result.OptimizerNotes {
				fmt.Printf("Optimizer: %s\n", n)
			}
			for _, w := range result.Warnings {
				fmt.Printf("Warning: %s\n", w)
			}

			if outFile != "" {
				if err := os.WriteFile(outFile, []byte(result.Code), 0644); err != nil {
					return fmt.Errorf("write error: %w", err)
				}
				fmt.Printf("Written to %s\n", outFile)
				return nil
			}
			fmt.Println(result.Code)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVar(&target, "target", string(compile.OpenQASM), "openqasm|qiskit|cirq|braket")
	f.StringVarP(&outFile, "output", "o", "", "write compiled output to file instead of stdout")
	f.BoolVar(&optimize, "optimize", true, "enable the IR optimizer (default)")
	f.BoolVar(&noOptimize, "no-optimize", false, "disable the IR optimizer")

	return cmd
}
