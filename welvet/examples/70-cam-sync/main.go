package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/layers/parallel"
	"github.com/openfluke/welvet/quant"
	"github.com/openfluke/welvet/stub/seed"
	"github.com/openfluke/welvet/weights"
)

func main() {
	s, err := parallel.Bicameral(8, 16, 1, core.ActivationLeakyReLU,
		core.DTypeFloat32, quant.FormatNone)
	if err != nil {
		panic(err)
	}
	hemi := s.Children[1].(*parallel.Layer)
	left := hemi.Branches[0].(*dense.Layer).Weights
	right := hemi.Branches[1].(*dense.Layer).Weights
	_ = seed.InitStoreHe(left, 16, seed.From("cam0"))
	_ = seed.InitStoreHe(right, 16, seed.From("cam1"))

	hemi.SetCamSync(parallel.CamSyncConfig{
		Enabled: true,
		Alpha:   1.0,
		When:    parallel.SyncManual,
	})
	before, _ := weights.StoreCosine(left, right)
	if err := hemi.SyncNow(); err != nil {
		panic(err)
	}
	after, _ := weights.StoreCosine(left, right)
	fmt.Printf("cosine before=%.4f after=%.4f\n", before, after)
}
