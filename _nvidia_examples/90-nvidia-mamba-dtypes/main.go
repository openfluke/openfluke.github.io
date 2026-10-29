// Float32 tensors with independently selected native parameter encodings.
package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/mamba"
	"github.com/openfluke/welvet/nativecuda"
	"github.com/openfluke/welvet/quant"
)

func main() {
	if e := run(); e != nil {
		log.Fatal(e)
	}
}

func run() error {
	selected := flag.String("dtype", "all", "native parameter dtype name, or all 34")
	flag.Parse()
	if !cudaPlatformOK() {
		fmt.Println("SKIP: NVIDIA CUDA needs linux/amd64 or windows/amd64 and cgo")
		return nil
	}
	count := 0
	for _, dt := range core.AllDTypes {
		if *selected != "all" && !strings.EqualFold(*selected, dt.String()) {
			continue
		}
		if e := forward(dt); e != nil {
			if count == 0 && *selected == "all" {
				fmt.Printf("SKIP: CUDA unavailable (%v)\n", e)
				return nil
			}
			return fmt.Errorf("%s: %w", dt, e)
		}
		count++
	}
	if count == 0 {
		return fmt.Errorf("unknown dtype %q", *selected)
	}
	fmt.Printf("Mamba dtypes=%d; activations=float32; parameters=FormatNone native dtype\n", count)
	return nil
}

func forward(dt core.DType) error {
	const d, inner = 8, 16
	in, out := make([]float32, 2*inner*d), make([]float32, d*inner)
	for i := range in {
		in[i] = float32(i%17-8) * .01
	}
	for i := range out {
		out[i] = float32(i%13-6) * .01
	}
	model, e := mamba.NewConfigured(mamba.Config{DModel: d, DState: 4, Expand: 2, SeqLen: 3}, dt, quant.FormatNone, in, out, nil, nil)
	if e != nil {
		return e
	}
	s, e := nativecuda.Mount(model)
	if e != nil {
		return e
	}
	defer s.Close()
	x := core.NewTensor[float32](1, 3, d)
	for i := range x.Data {
		x.Data[i] = float32(i%7-3) * .05
	}
	post, e := s.Forward(x)
	if e != nil {
		return e
	}
	var checksum float64
	for _, v := range post.Data {
		checksum += float64(v)
	}
	fmt.Printf("%-10s shape=%v checksum=%.8g\n", dt, post.Shape, checksum)
	return nil
}

func cudaPlatformOK() bool {
	return runtime.GOARCH == "amd64" && (runtime.GOOS == "linux" || runtime.GOOS == "windows")
}
