// Copyright 2026 Magnobit. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/magnobit/quell/compile"
	"github.com/magnobit/quell/simulate"
	"github.com/spf13/cobra"
)

func newSimulateCmd() *cobra.Command {
	var shots int
	var noiseFlags []string

	cmd := &cobra.Command{
		Use:     "simulate <file.quell>",
		Short:   "Simulate a circuit locally (no backend, no credentials, no network)",
		Example: "  quell simulate bell.quell --shots 2000\n  quell simulate bell.quell --noise depolarizing=0.01",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			compiled, err := compile.CompileFileWithWarnings(args[0], compile.OpenQASM, true)
			if err != nil {
				return fmt.Errorf("parse/compile error: %w", err)
			}
			fmt.Printf("Qubits  : %d\n", compiled.NumQubits)
			fmt.Printf("Gates   : %d\n", compiled.NumInstructions)

			noise, err := mergeNoiseFlags(noiseFlags)
			if err != nil {
				return err
			}
			if noise.Active() {
				fmt.Printf("Noise   : depolarizing=%g amplitude_damping=%g phase_damping=%g bit_flip=%g readout=%g\n",
					noise.Depolarizing, noise.AmplitudeDamping, noise.PhaseDamping, noise.BitFlip, noise.ReadoutError)
			}

			result, err := simulate.RunFileOpts(args[0], simulate.Options{Shots: shots, Noise: noise})
			if err != nil {
				return fmt.Errorf("simulate error: %w", err)
			}
			result.Print()
			return nil
		},
	}

	cmd.Flags().IntVar(&shots, "shots", 1000, "number of measurement samples")
	cmd.Flags().StringArrayVar(&noiseFlags, "noise", nil, "noise model: depolarizing|amplitude_damping|phase_damping|bit_flip|readout=<p> (repeatable)")
	return cmd
}

func mergeNoiseFlags(flags []string) (simulate.NoiseModel, error) {
	var n simulate.NoiseModel
	for _, f := range flags {
		part, err := simulate.ParseNoiseFlag(f)
		if err != nil {
			return n, err
		}
		n = simulate.MergeNoise(n, part)
	}
	return n, n.Validate()
}
