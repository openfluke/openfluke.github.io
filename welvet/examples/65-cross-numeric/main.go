package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/quant"
	"github.com/openfluke/welvet/runtime/training"
)

func main() {
	// Weight storage: int8 FormatNone. Activations/grads: Tensor[int8].
	g := architecture.NewGrid(1, 1, 1, 1)
	init := make([]float32, 4*4)
	for i := range init {
		init[i] = 0.1
	}
	l, err := dense.NewConfigured(4, 4, core.ActivationLinear, core.DTypeInt8, quant.FormatNone, init)
	if err != nil {
		panic(err)
	}
	if err := dense.Place(g, 0, 0, 0, 0, l); err != nil {
		panic(err)
	}
	x := core.NewTensor[int8](1, 4)
	y := core.NewTensor[int8](1, 4)
	for i := 0; i < 4; i++ {
		x.Data[i] = int8(i + 1)
		y.Data[i] = int8(i)
	}
	loss, _, err := training.StepMesh(g, x, y, 1, 0.05)
	fmt.Println(loss, err, l.Weights.RetainsF32Master())
}
