package main

import "fmt"

func main() {
	// Recompute when a board row flips ✅/🚧/⬜ in welvet/README.md
	earned := 76.0
	fmt.Printf("v0.%02.0f\n", earned) // v0.76 until earned==100 → v1.0
}
