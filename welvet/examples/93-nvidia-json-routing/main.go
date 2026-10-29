// JSON selects backends; the existing Go Stack scheduler moves activations/gradients.
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/layers/parallel"
	"github.com/openfluke/welvet/nativecpp"
	"github.com/openfluke/welvet/nativecuda"
	"github.com/openfluke/welvet/quant"
	"github.com/openfluke/welvet/runtime/placement"
)

//go:embed placement.json
var defaultPlan []byte

func main() {
	if e := run(); e != nil {
		log.Fatal(e)
	}
}

func run() error {
	path := flag.String("config", "", "optional replacement execution JSON")
	flag.Parse()
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		fmt.Println("SKIP: CUDA JSON routing needs Linux amd64 and cgo")
		return nil
	}
	raw := defaultPlan
	if *path != "" {
		var e error
		raw, e = os.ReadFile(*path)
		if e != nil {
			return e
		}
	}
	cfg, e := nativecpp.ParseExecConfigJSON(raw)
	if e != nil {
		return e
	}
	router, e := nativecuda.EnableRouting()
	if e != nil {
		fmt.Printf("SKIP: CUDA routing unavailable (%v)\n", e)
		return nil
	}
	defer router.Close()
	var ops []any
	for i := 0; i < 3; i++ {
		l, e := dense.NewConfigured(2, 2, core.ActivationTanh, core.DTypeFloat32, quant.FormatNone, []float32{.8, .1, -.1, .9})
		if e != nil {
			return e
		}
		ops = append(ops, l)
	}
	model, e := parallel.NewStack(ops...)
	if e != nil {
		return e
	}
	if e = placement.Stack(model, cfg, "main", "0,0,0,0"); e != nil {
		return e
	}
	for i, op := range ops {
		fmt.Printf("layer=%d backend=%s\n", i, op.(*dense.Layer).Exec.Backend)
	}
	x := &core.Tensor[float32]{Shape: []int{1, 2}, Data: []float32{.5, -.3}}
	target := core.NewTensor[float32](1, 2)
	loss, e := parallel.TrainStackMSE(model, x, target, parallel.ModeNormalBP, .03)
	if e != nil {
		fmt.Printf("SKIP: CUDA train unavailable (%v)\n", e)
		return nil
	}
	stats := router.Stats()
	if stats.Backwards == 0 || stats.Updates == 0 {
		return fmt.Errorf("plan did not execute native CUDA backward/SGD")
	}
	fmt.Printf("loss=%.8g cuda routing=%+v\nRouted updates are already in Go stores; Close releases cached mounts.\n", loss, stats)
	return nil
}
