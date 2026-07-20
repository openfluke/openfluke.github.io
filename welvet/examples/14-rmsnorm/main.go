package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/rmsnorm"
	"github.com/openfluke/welvet/quant"
)

func main() {
	gamma := []float32{1, 1, 1, 1}
	l, err := rmsnorm.NewConfigured(rmsnorm.Config{Dim: 4}, core.DTypeFloat32, quant.FormatNone, gamma)
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 4)
	_, y, err := rmsnorm.Forward(l, x)
	fmt.Println(y.Data, err)
}
