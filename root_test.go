// Copyright 2026 Magnobit. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"strings"
	"testing"
)

func TestRootHasAnalysisCommands(t *testing.T) {
	root := newRootCmd()
	want := []string{"draw", "state", "observe", "gradient", "vqe", "inspect", "simulate"}
	for _, name := range want {
		if root.Commands() == nil {
			t.Fatal("no commands")
		}
		found := false
		for _, c := range root.Commands() {
			if c.Name() == name {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing command %s", name)
		}
	}
	help := root.UsageString()
	if !strings.Contains(help, "draw") {
		t.Fatal(help)
	}
}
