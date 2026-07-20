package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/quant"
)

func main() {
	init := make([]float32, 4*8)
	l, err := dense.NewConfigured(8, 4, core.ActivationReLU, core.DTypeFloat32, quant.FormatNone, init)
	if err != nil {
		panic(err)
	}
	l.Exec.Backend = core.BackendCPUTiled
	x := core.NewTensor[float32](1, 8)
	pre, post, err := dense.Forward(l, x)
	if err != nil {
		panic(err)
	}
	gIn, gW, err := dense.Backward(l, post, x, pre)
	_ = dense.ApplyGradSGD(l, gW, 1e-3)
	fmt.Println(len(gIn.Data), len(gW.Data))
}
