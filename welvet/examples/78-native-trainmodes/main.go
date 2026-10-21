package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/layers/parallel"
	"github.com/openfluke/welvet/nativecpp"
	"github.com/openfluke/welvet/quant"
)

func main() {
	if !nativecpp.NativeCPUAvailable() {
		fmt.Println("SKIP: need linux/amd64 libwelvet")
		return
	}
	demoStackFamilies()
	demoParallelModes()
	demoGridMesh()
	demoLineStep()
	fmt.Println("native-trainmodes done")
}

func demoStackFamilies() {
	modes := []parallel.TrainMode{
		parallel.ModeNormalBP, parallel.ModeTween, parallel.ModeTweenChain,
		parallel.ModeTweenSplit, parallel.ModeTweenAlt, parallel.ModeFreeze,
	}
	for _, mode := range modes {
		children := make([]any, 2)
		for i := 0; i < 2; i++ {
			init := make([]float32, 16)
			for j := range init {
				init[j] = 0.04 * float32(i*16+j+1)
			}
			d, err := dense.NewConfigured(4, 4, core.ActivationLinear, core.DTypeFloat32, quant.FormatNone, init)
			must(err)
			children[i] = d
		}
		s, err := parallel.NewStack(children...)
		must(err)
		s.SetBackend(core.BackendNativeCPU)
		x, t := core.NewTensor[float32](2, 4), core.NewTensor[float32](2, 4)
		for i := range x.Data {
			x.Data[i] = 0.08 * float32(i+1)
			t.Data[i] = 0.3
		}
		loss, err := parallel.TrainStackMSE(s, x, t, mode, 0.05)
		must(err)
		must(s.SyncWeightsFromNative())
		fmt.Println("stack", mode, "loss", loss)
	}
}

func demoParallelModes() {
	for _, mode := range []parallel.TrainMode{parallel.ModeTween, parallel.ModeAdversarial} {
		init := make([]float32, 2*4*4)
		for i := range init {
			init[i] = 0.02 * float32(i+1)
		}
		l, err := parallel.NewConfigured(parallel.Config{Dim: 4, OutFeat: 4, Branches: 2, Combine: parallel.CombineAdd},
			core.DTypeFloat32, quant.FormatNone, init, nil)
		must(err)
		l.Exec.Backend = core.BackendNativeCPU
		l.SyncBranchExec()
		x, t := core.NewTensor[float32](2, 4), core.NewTensor[float32](2, 4)
		for i := range x.Data {
			x.Data[i] = 0.07 * float32(i+1)
			t.Data[i] = 0.25
		}
		loss, err := parallel.TrainMSE(l, x, t, mode, 0.08)
		must(err)
		must(l.SyncWeightsFromNative())
		fmt.Println("parallel", mode, "loss", loss)
	}
}

func demoGridMesh() {
	const dim, batch, ticks = 3, 2, 2
	packed := make([]float32, 2*dim*dim)
	for i := range packed {
		packed[i] = 0.05 * float32(i+1)
	}
	wAfter := make([]float32, len(packed))
	x := make([]float32, batch*dim)
	tgt := make([]float32, batch*dim)
	for i := range x {
		x[i] = 0.1 * float32(i+1)
		tgt[i] = 0.3
	}
	loss, err := nativecpp.GridTrainMSE(
		int32(parallel.ModeMeshBP), 1, 1, 2, 1, batch, ticks, dim,
		1, 0, int32(core.ActivationLinear), int32(core.DTypeFloat32), 0,
		packed, x, tgt, 0.1, wAfter,
	)
	must(err)
	fmt.Println("grid MeshBP loss", loss)
}

func demoLineStep() {
	const nChildren, nTicks, batch, dim = 2, 3, 2, 3
	inH := []int32{dim, dim}
	outH := []int32{dim, dim}
	act := []int32{int32(core.ActivationLinear), int32(core.ActivationLinear)}
	packed := make([]float32, nChildren*dim*dim)
	for i := range packed {
		packed[i] = 0.04 * float32(i+1)
	}
	wAfter := make([]float32, len(packed))
	xTicks := make([]float32, nTicks*batch*dim)
	tgtTicks := make([]float32, nTicks*batch*dim)
	for i := range xTicks {
		xTicks[i] = 0.05 * float32(i+1)
		tgtTicks[i] = 0.2
	}
	loss, err := nativecpp.StackLineTrainMSE(
		int32(parallel.ModeStepTween), nChildren, nTicks, batch, 0.1, int32(core.DTypeFloat32),
		inH, outH, act, packed, xTicks, tgtTicks, wAfter,
	)
	must(err)
	fmt.Println("line StepTween loss", loss)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
