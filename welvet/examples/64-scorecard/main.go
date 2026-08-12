package main

import "fmt"

func main() {
	// Scorecard earned from welvet/README.md; Version cell may be a patch tag.
	earned := 95.0
	fmt.Println("v0.95.1") // patch tag; scorecard round(earned) → v0.95 until 100 → v1.0
	_ = earned
}
