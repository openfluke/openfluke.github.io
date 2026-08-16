package main

import "fmt"

func main() {
	fmt.Println("w2a is a separate module — engine packages never contain tests.")
	fmt.Println("")
	fmt.Println("  cd w2a")
	fmt.Println("  go run .                                      # interactive; [0] = ALL")
	fmt.Println("  go test ./tests/dense -v")
	fmt.Println("  go test ./tests/parallel -run Test49AllTrainModesCubes -count=1 -v")
}
