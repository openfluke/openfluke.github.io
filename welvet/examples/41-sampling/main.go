package main

import (
	"fmt"

	"github.com/openfluke/welvet/model/sampling"
)

func main() {
	logits := []float32{0.1, 2.4, 0.3, -1}
	fmt.Println(sampling.ArgMax(logits))
	fmt.Println(sampling.SampleTopK(logits, 2, 0.8, true))
}
