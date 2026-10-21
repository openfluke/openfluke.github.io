package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/nativecpp"
	"github.com/openfluke/welvet/quant"
)

func main() {
	fmt.Println("version", nativecpp.Version())
	fmt.Println("available", nativecpp.NativeCPUAvailable())
	fmt.Println("default", core.DefaultBackend())
	fmt.Println("prefer", core.PreferNativeCPU())
	if !nativecpp.NativeCPUAvailable() {
		fmt.Println("SKIP: need linux/amd64 libwelvet")
		return
	}

	init := make([]float32, 4*4)
	for i := range init {
		init[i] = float32(i+1) * 0.05
	}
	l, err := dense.NewConfigured(4, 4, core.ActivationLinear, core.DTypeFloat32, quant.FormatNone, init)
	must(err)
	l.Exec.Backend = core.BackendNativeCPU

	cfg, err := nativecpp.ParseExecConfigJSON([]byte(`{"default":"native_cpu"}`))
	must(err)
	l.ApplyExecConfig(cfg, "main", "0,0,0,0", "0")
	fmt.Println("exec", l.Exec.Backend)

	x := core.NewTensor[float32](2, 4)
	for i := range x.Data {
		x.Data[i] = float32(i+1) * 0.1
	}
	pre, post, err := dense.Forward(l, x)
	must(err)
	_, gW, err := dense.Backward(l, post, x, pre)
	must(err)
	fmt.Println("fwd/bwd", len(post.Data), len(gW.Data))

	tgt := core.NewTensor[float32](2, 4)
	for i := range tgt.Data {
		tgt.Data[i] = 0.25
	}
	loss, err := dense.TrainNativeMSE(l, x, tgt, 0.05)
	must(err)
	fmt.Println("pending", l.NativeWeightsPending(), "loss", loss)
	must(l.SyncWeightsFromNative())
	fmt.Println("synced", !l.NativeWeightsPending())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
