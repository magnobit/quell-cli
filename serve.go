// Copyright 2026 Magnobit. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/magnobit/quell/compile"
	"github.com/spf13/cobra"
)

const maxRequestBytes = 1 << 20 // 1 MB

func newServeCmd() *cobra.Command {
	var port string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP compile server",
		Example: `  quell serve
  quell serve --port 9000`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("port") {
				if p := os.Getenv("PORT"); p != "" {
					port = p
				}
			}
			return serve(port)
		},
	}

	cmd.Flags().StringVar(&port, "port", "8081", "port to listen on (env PORT)")
	return cmd
}

func serve(port string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"service": "quell-compiler",
			"version": version,
		})
	})

	mux.HandleFunc("OPTIONS /compile", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("POST /compile", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		// Panic recovery — a compiler bug must not crash the server
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("compile panic: %v\n%s", rec, debug.Stack())
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": fmt.Sprintf("internal compiler error: %v", rec),
				})
			}
		}()

		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)

		var req struct {
			Code     string `json:"code"`
			Target   string `json:"target"`
			Optimize *bool  `json:"optimize"` // defaults to true when omitted
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			status := http.StatusBadRequest
			msg := "invalid request body"
			if err.Error() == "http: request body too large" {
				status = http.StatusRequestEntityTooLarge
				msg = fmt.Sprintf("request body exceeds %d bytes", maxRequestBytes)
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": msg})
			return
		}
		if req.Code == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "code is required"})
			return
		}

		target := compile.Target(req.Target)
		if target == "" {
			target = compile.OpenQASM
		}

		optimize := true
		if req.Optimize != nil {
			optimize = *req.Optimize
		}

		result, err := compile.CompileWithWarnings(req.Code, target, optimize)
		if err != nil {
			errType := "compile"
			errMsg := err.Error()
			if strings.Contains(errMsg, "line ") || strings.Contains(errMsg, "empty circuit") {
				errType = "parse"
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(map[string]string{
				"error":     errMsg,
				"errorType": errType,
			})
			return
		}

		lang := "python"
		if target == compile.OpenQASM {
			lang = "openqasm"
		}

		json.NewEncoder(w).Encode(map[string]any{
			"result":         result.Code,
			"target":         string(target),
			"language":       lang,
			"optimizerNotes": result.OptimizerNotes,
		})
	})

	fmt.Printf("Quell compile server v%s listening on :%s\n", version, port)
	return http.ListenAndServe(":"+port, mux)
}
