package main

import "fmt"

func main() {
	// Recompute when a board row flips ✅/🚧/⬜ in welvet/README.md
	earned := 95.0
	fmt.Printf("v0.%02.0f\n", earned) // round(95) → v0.95 until earned==100 → v1.0
}
