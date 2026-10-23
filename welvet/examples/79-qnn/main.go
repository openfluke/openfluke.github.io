package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/nativeqnn"
	"github.com/openfluke/welvet/quant"
)

func main() {
	fmt.Println("version", nativeqnn.Version())
	fmt.Println("available", nativeqnn.Available())
	fmt.Println("default", core.DefaultBackend())
	if !nativeqnn.Available() {
		fmt.Println("SKIP: need windows/arm64 + libwelvet_qnn (+ QAIRT)")
		return
	}
	fmt.Println("layers", nativeqnn.SupportedLayers())

	init := make([]float32, 4*4)
	for i := range init {
		init[i] = float32(i+1) * 0.05
	}
	l, err := dense.NewConfigured(4, 4, core.ActivationLinear, core.DTypeFloat32, quant.FormatNone, init)
	must(err)
	must(l.SyncToQNN())
	must(l.PrepareQNN(2))
	l.Exec.Backend = core.BackendQNN

	x := core.NewTensor[float32](2, 4)
	for i := range x.Data {
		x.Data[i] = float32(i+1) * 0.1
	}
	_, post, err := dense.Forward(l, x)
	must(err)
	fmt.Println("fwd", len(post.Data), "backend", l.Exec.Backend)
	fmt.Println("qnn done")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
