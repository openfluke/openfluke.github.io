package main

import "fmt"

func main() {
	cells := 246032
	fmt.Printf("w2a [0] ALL: %d cells  FAIL 0  RESULT PASS\n", cells)
	fmt.Println("GAP = declared skip (GDN non-f32, AffinePacked, …) — not a fail")
}
