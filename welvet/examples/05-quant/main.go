package main

import (
	"fmt"

	"github.com/openfluke/welvet/quant"
)

func main() {
	w := []float32{0.1, -0.2, 0.3, 0.4, -0.5, 0.6, 0.05, -0.15}
	b, err := quant.Pack(quant.FormatQ4_0, w, 2, 4)
	if err != nil {
		panic(err)
	}
	y := make([]float32, 2)
	if err := quant.MatVec(b, []float32{1, 1, 1, 1}, y); err != nil {
		panic(err)
	}
	fmt.Println("rows", b.Rows, "cols", b.Cols, "y", y)
}
