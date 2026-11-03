// Explicit mount ownership, weight return to Go, ENTITY save and reload.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/layers/parallel"
	"github.com/openfluke/welvet/model/entity"
	"github.com/openfluke/welvet/nativevulkan"
	"github.com/openfluke/welvet/quant"
)

func main() {
	if e := run(); e != nil {
		log.Fatal(e)
	}
}

func run() error {
	out := flag.String("out", "", "ENTITY destination; default is a new temporary directory")
	flag.Parse()
	if !vulkanPlatformOK() {
		fmt.Println("SKIP: Vulkan needs linux/amd64 and cgo")
		return nil
	}
	l, e := dense.NewConfigured(2, 2, core.ActivationLinear, core.DTypeFloat32, quant.FormatNone, []float32{.3, -.2, .1, .4})
	if e != nil {
		return e
	}
	model, e := parallel.NewStack(l)
	if e != nil {
		return e
	}
	before, e := l.Weights.FlattenF32()
	if e != nil {
		return e
	}
	s, e := nativevulkan.Mount(model)
	if e != nil {
		fmt.Printf("SKIP: Vulkan unavailable (%v)\n", e)
		return nil
	}
	defer s.Close()
	x := &core.Tensor[float32]{Shape: []int{1, 2}, Data: []float32{.5, .25}}
	target := core.NewTensor[float32](1, 2)
	for i := 0; i < 5; i++ {
		if _, e = s.TrainMSE(x, target, parallel.ModeNormalBP, .1); e != nil {
			return e
		}
	}
	unchanged, e := l.Weights.FlattenF32()
	if e != nil {
		return e
	}
	if !reflect.DeepEqual(before, unchanged) {
		return fmt.Errorf("explicit mount unexpectedly changed Go parameters before SyncToGo")
	}
	fmt.Println("Go parameters unchanged before SyncToGo")
	if e = s.SyncToGo(); e != nil {
		return e
	}
	after, e := l.Weights.FlattenF32()
	if e != nil {
		return e
	}
	if reflect.DeepEqual(before, after) {
		return fmt.Errorf("expected trained parameters after SyncToGo")
	}
	if *out == "" {
		dir, e := os.MkdirTemp("", "welvet-vulkan-entity-")
		if e != nil {
			return e
		}
		*out = filepath.Join(dir, "trained.entity")
	}
	if e = entity.WriteCameralFile(*out, model, parallel.ModeNormalBP); e != nil {
		return e
	}
	loaded, mode, e := entity.LoadCameral(*out)
	if e != nil {
		return e
	}
	if mode != parallel.ModeNormalBP {
		return fmt.Errorf("training mode changed during save")
	}
	restored, e := nativevulkan.Mount(loaded)
	if e != nil {
		return e
	}
	defer restored.Close()
	want, e := s.Weights()
	if e != nil {
		return e
	}
	got, e := restored.Weights()
	if e != nil {
		return e
	}
	if !reflect.DeepEqual(want, got) {
		return fmt.Errorf("saved parameters changed on reload")
	}
	fmt.Printf("weights returned to Go; saved=%s; reload=exact\n", *out)
	return nil
}

func vulkanPlatformOK() bool {
	return runtime.GOOS == "linux" && runtime.GOARCH == "amd64"
}
