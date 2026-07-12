// Copyright 2026 Magnobit. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import "os"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
