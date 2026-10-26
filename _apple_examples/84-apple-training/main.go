// Native Metal MSE training: persistent Stack Line state or an actual Grid clock.
package main

import (
	"flag"
	"fmt"
	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/layers/parallel"
	"github.com/openfluke/welvet/nativemetal"
	"github.com/openfluke/welvet/quant"
	"log"
	"runtime"
)

func main() {
	if e := run(); e != nil {
		log.Fatal(e)
	}
}
func run() error {
	selected := flag.String("mode", "NormalBP", "training mode name, or all 29")
	steps := flag.Int("steps", 8, "training calls; retain the session for Step pipeline state")
	flag.Parse()
	if *steps < 1 {
		return fmt.Errorf("steps must be positive")
	}
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		fmt.Println("SKIP: Apple Metal needs Apple silicon macOS and cgo")
		return nil
	}
	if *selected == "all" {
		for mode := parallel.ModeNormalBP; mode <= parallel.ModeStepTweenSplitSparse; mode++ {
			if e := train(mode, *steps); e != nil {
				return e
			}
		}
		return nil
	}
	mode, e := parallel.ParseTrainMode(*selected)
	if e != nil {
		return e
	}
	if mode < parallel.ModeNormalBP || mode > parallel.ModeStepTweenSplitSparse {
		return fmt.Errorf("choose one of the 29 update modes")
	}
	return train(mode, *steps)
}
func train(mode parallel.TrainMode, steps int) error {
	var ops []any
	for layer := 0; layer < 2; layer++ {
		w := make([]float32, 16)
		for i := range w {
			w[i] = float32(i%5-2) * .04
		}
		l, e := dense.NewConfigured(4, 4, core.ActivationTanh, core.DTypeFloat32, quant.FormatNone, w)
		if e != nil {
			return e
		}
		ops = append(ops, l)
	}
	stack, e := parallel.NewStack(ops...)
	if e != nil {
		return e
	}
	var model any = stack
	schedule := "Stack"
	if mode.RequiresGrid() {
		g := architecture.NewGrid(3, 1, 1, 1)
		g.Cells[0].Layer.IsDisabled = true
		for i, op := range ops {
			g.Cells[i+1].Layer = op.(*dense.Layer).Core
			g.Cells[i+1].Op = op
		}
		model = g
		schedule = "Grid mesh clock (4 ticks/call)"
	} else if mode.IsLineStep() {
		schedule = "Stack Line pipeline"
	}
	s, e := nativemetal.Mount(model)
	if e != nil {
		return e
	}
	defer s.Close()
	x := core.NewTensor[float32](2, 4)
	for i := range x.Data {
		x.Data[i] = float32(i%5-2) * .1
	}
	target := core.NewTensor[float32](2, 4)
	var loss float64
	for i := 0; i < steps; i++ {
		loss, e = s.TrainMSE(x, target, mode, .01, 4)
		if e != nil {
			return fmt.Errorf("%s: %w", mode, e)
		}
	}
	if e = s.SyncToGo(); e != nil {
		return e
	}
	fmt.Printf("mode=%s schedule=%s calls=%d loss=%.8g weights=synced-to-Go\n", mode, schedule, steps, loss)
	return nil
}
