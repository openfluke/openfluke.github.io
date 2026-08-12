package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/parallel"
	"github.com/openfluke/welvet/quant"
)

func main() {
	// Dense stem → 2 hemispheres (add) → Dense head
	s, err := parallel.Bicameral(8, 16, 1, core.ActivationLeakyReLU,
		core.DTypeFloat32, quant.FormatNone)
	if err != nil {
		panic(err)
	}
	// Mix: left hemi StepBP, right hemi StepTweenChain
	hemi := s.Children[1].(*parallel.Layer)
	hemi.SetBranchModes(parallel.ModeStepBP, parallel.ModeStepTweenChain)

	x := core.NewTensor[float32](1, 8)
	t := core.NewTensor[float32](1, 1)
	t.Data[0] = 0.5
	loss, err := parallel.TrainStackMSE(s, x, t, parallel.ModeStepBP, 0.01)
	fmt.Println("loss", loss, "err", err)
}
