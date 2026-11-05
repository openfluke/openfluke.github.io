// SyncToMLX one Dense layer, forward on Apple MLX GPU, release nothing explicit
// (shared DefaultSession lives for the process).
package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
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
	fmt.Println("version", nativemlx.Version())
	fmt.Println("available", nativemlx.Available())
	fmt.Println("prefer", core.PreferMLX())
	if !nativemlx.Available() {
		fmt.Println("SKIP: need libwelvet_mlx.dylib + Metal GPU (../welvet.mlx/scripts/macos_install_go.sh)")
		return nil
	}

	l, err := dense.NewConfigured(2, 2, core.ActivationLinear, core.DTypeFloat32, quant.FormatNone, []float32{1, 0, 0, 1})
	if err != nil {
		return err
	}
	if err := l.SyncToMLX(); err != nil {
		return err
	}
	l.Exec.Backend = core.BackendMLX

	x := &core.Tensor[float32]{Shape: []int{1, 2}, Data: []float32{0.5, -0.25}}
	_, post, err := dense.Forward(l, x)
	if err != nil {
		return err
	}
	fmt.Printf("backend=mlx synced=%v\ninput=%v output=%v shape=%v\n", l.MLXSynced(), x.Data, post.Data, post.Shape)
	return nil
}
