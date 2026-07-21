package main

import "fmt"

func main() {
	// Recompute when a board row flips ✅/🚧/⬜ in welvet/README.md
	earned := 81.0
	fmt.Printf("v0.%02.0f\n", earned) // round(81) → v0.81 until earned==100 → v1.0
}
