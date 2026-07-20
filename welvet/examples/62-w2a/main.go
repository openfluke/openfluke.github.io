package main

import "fmt"

func main() {
	fmt.Println("w2a is a separate module — engine packages never contain tests.")
	fmt.Println("")
	fmt.Println("  cd w2a")
	fmt.Println("  go run .                    # interactive menu")
	fmt.Println("  go test ./tests/dense -v    # timed FormatNone matrix")
	fmt.Println("  go test ./tests/mha -v      # MHA coverage")
}
