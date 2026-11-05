// Dense PrepareMLX + Sequential SyncToMLX forward on Apple MLX.
package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/layers/sequential"
	"github.com/openfluke/welvet/nativemlx"
	"github.com/openfluke/welvet/quant"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		fmt.Println("SKIP: Apple MLX needs Apple silicon macOS and cgo")
		return nil
	}
	if !nativemlx.Available() {
		fmt.Println("SKIP: need libwelvet_mlx.dylib + Metal GPU")
		return nil
	}

	if err := demoDensePrepare(); err != nil {
		return err
	}
	return demoSequential()
}

func demoDensePrepare() error {
	init := make([]float32, 4*4)
	for i := range init {
		init[i] = float32(i+1) * 0.05
	}
	l, err := dense.NewConfigured(4, 4, core.ActivationLinear, core.DTypeFloat32, quant.FormatNone, init)
	if err != nil {
		return err
	}
	if err := l.SyncToMLX(); err != nil {
		return err
	}
	if err := l.PrepareMLX(1); err != nil {
		return err
	}
	l.Exec.Backend = core.BackendMLX
	x := core.NewTensor[float32](1, 4)
	for i := range x.Data {
		x.Data[i] = 0.1 * float32(i+1)
	}
	_, post, err := dense.Forward(l, x)
	if err != nil {
		return err
	}
	fmt.Printf("dense prepare+exec out=%v\n", post.Data)
	return nil
}

func demoSequential() error {
	// Two identity Dense children (Depth=2, Dim=3).
	packed := append(identity3(), identity3()...)
	seq, err := sequential.NewConfigured(sequential.Config{Dim: 3, Depth: 2}, core.DTypeFloat32, quant.FormatNone, packed)
	if err != nil {
		return err
	}
	if err := seq.SyncToMLX(); err != nil {
		return err
	}
	seq.Exec.Backend = core.BackendMLX
	x := &core.Tensor[float32]{Shape: []int{1, 3}, Data: []float32{1, 2, 3}}
	_, post, err := sequential.Forward(seq, x)
	if err != nil {
		return err
	}
	fmt.Printf("sequential mlx out=%v synced=%v\n", post.Data, seq.MLXSynced())
	return nil
}

func identity3() []float32 {
	return []float32{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
}
