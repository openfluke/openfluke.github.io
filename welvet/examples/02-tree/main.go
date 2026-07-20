package main

import (
	"fmt"
	"os"
)

func main() {
	// Mental model only — list engine roots you depend on:
	roots := []string{
		"core", "weights", "quant", "simd", "webgpu", "tiling",
		"architecture", "fusedgpu", "layers", "runtime", "systems", "model",
	}
	for _, r := range roots {
		fmt.Println("import github.com/openfluke/welvet/" + r + "/…")
	}
	_ = os.Getenv // apps/stub/w2a are separate concerns
}
