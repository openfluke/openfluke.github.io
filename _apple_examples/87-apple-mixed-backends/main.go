// Keep a normal cgo build: -tags welvet_metal intentionally disables C++.
package main

import (
	_ "embed"
	"fmt"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/layers/parallel"
	"github.com/openfluke/welvet/nativecpp"
	"github.com/openfluke/welvet/nativemetal"
	"github.com/openfluke/welvet/quant"
	"github.com/openfluke/welvet/runtime/placement"
	"log"
	"runtime"
)

//go:embed placement.json
var plan []byte

func main() {
	if e := run(); e != nil {
		log.Fatal(e)
	}
}
func run() error {
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		fmt.Println("SKIP: mixed Apple routing needs Apple silicon macOS and cgo")
		return nil
	}
	if !nativecpp.NativeCPUAvailable() {
		fmt.Println("SKIP: mixed routing needs the C++ artifact and normal build without -tags welvet_metal")
		return nil
	}
	router, e := nativemetal.EnableRouting()
	if e != nil {
		return e
	}
	defer router.Close()
	cfg, e := nativecpp.ParseExecConfigJSON(plan)
	if e != nil {
		return e
	}
	var ops []any
	for i := 0; i < 13; i++ {
		l, e := dense.NewConfigured(2, 2, core.ActivationTanh, core.DTypeFloat32, quant.FormatNone, []float32{.8, .12, -.15, .9})
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
	x := &core.Tensor[float32]{Shape: []int{1, 2}, Data: []float32{.5, -.3}}
	target := core.NewTensor[float32](1, 2)
	before := nativecpp.CompletedCalls()
	var loss float64
	for i := 0; i < 5; i++ {
		loss, e = parallel.TrainStackMSE(model, x, target, parallel.ModeNormalBP, .03)
		if e != nil {
			return e
		}
	}
	calls := nativecpp.CompletedCalls() - before
	stats := router.Stats()
	if calls == 0 || stats.Backwards == 0 || stats.Updates == 0 {
		return fmt.Errorf("expected both C++ and Metal native calls")
	}
	fmt.Println("Go layers 0-3 -> C++ layers 4-8 -> Metal layers 9-12")
	fmt.Printf("loss=%.8g cpp calls=%d metal routing=%+v\n", loss, calls, stats)
	fmt.Println("Float32 activations travel forward; gradients travel back; routed Metal weights return automatically to Go.")
	return nil
}
