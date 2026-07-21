package main

import "fmt"

func main() {
	// Recompute when a board row flips ✅/🚧/⬜ in welvet/README.md
	earned := 77.5
	fmt.Printf("v0.%02.0f\n", earned) // round(77.5) → v0.78 until earned==100 → v1.0
}
