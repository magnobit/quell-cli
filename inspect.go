// Copyright 2026 Magnobit. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newInspectCmd() *cobra.Command {
	var kind string
	cmd := &cobra.Command{
		Use:   "inspect <file.quell>",
		Short: "Print the parse summary, canonical IR, optimized IR, or QIR subset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := loadAnalysis(args[0])
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			switch strings.ToLower(kind) {
			case "ir":
				_, err = fmt.Fprint(out, m.Canonical())
			case "optimized":
				_, err = fmt.Fprint(out, m.CanonicalOptimized())
			case "qir":
				text, qerr := m.QIR()
				if qerr != nil {
					return qerr
				}
				_, err = fmt.Fprint(out, text)
			default:
				fn, gates := m.Summary()
				_, err = fmt.Fprintf(out, "functions %d\ngates %d\n", fn, gates)
			}
			return err
		},
	}
	cmd.Flags().StringVar(&kind, "kind", "ast", "ast, ir, optimized, or qir")
	return cmd
}
