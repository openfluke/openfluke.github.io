package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/quant"
	"github.com/openfluke/welvet/weights"
)

func main() {
	s, err := weights.New[float32](2, 2, []float32{1, 0, 0, 1}, core.DTypeFloat32, quant.FormatNone)
	if err != nil {
		panic(err)
	}
	y := make([]float32, 2)
	if err := weights.MatVec(s, []float32{3, 4}, y); err != nil {
		panic(err)
	}
	fmt.Println(y) // ≈ [3, 4]
}
