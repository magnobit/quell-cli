// Copyright 2026 Magnobit. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/magnobit/quell/compile"
	"github.com/magnobit/quell/simulate"
	"github.com/spf13/cobra"
)

// newSimulateCmd is an explicit, discoverable entry point for local
// simulation — equivalent to `quell run --backend local`, which does the
// same thing, but named for anyone specifically looking for a "just run it
// locally, no config file needed" command rather than the backend-selection
// mental model `run` uses.
func newSimulateCmd() *cobra.Command {
	var shots int

	cmd := &cobra.Command{
		Use:     "simulate <file.quell>",
		Short:   "Simulate a circuit locally (no backend, no credentials, no network)",
		Example: "  quell simulate bell.quell --shots 2000",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// compile.CompileFileWithWarnings gives us NumQubits/
			// NumInstructions for the header — simulate.RunFile alone
			// doesn't expose the gate count, only the simulation result.
			compiled, err := compile.CompileFileWithWarnings(args[0], compile.OpenQASM, true)
			if err != nil {
				return fmt.Errorf("parse/compile error: %w", err)
			}
			fmt.Printf("Qubits  : %d\n", compiled.NumQubits)
			fmt.Printf("Gates   : %d\n", compiled.NumInstructions)

			result, err := simulate.RunFile(args[0], shots)
			if err != nil {
				return fmt.Errorf("simulate error: %w", err)
			}
			result.Print()
			return nil
		},
	}

	cmd.Flags().IntVar(&shots, "shots", 1000, "number of measurement samples")
	return cmd
}
