// Mount one layer, execute its forward pass on NVIDIA CUDA, and release it.
package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/nativecuda"
	"github.com/openfluke/welvet/quant"
)

func main() {
	if e := run(); e != nil {
		log.Fatal(e)
	}
}

func run() error {
	library := flag.String("library", "", "optional libwelvet_cuda.so/.dll path")
	flag.Parse()
	if !cudaPlatformOK() {
		fmt.Println("SKIP: NVIDIA CUDA needs linux/amd64 or windows/amd64 and cgo")
		return nil
	}
	l, e := dense.NewConfigured(2, 2, core.ActivationLinear, core.DTypeFloat32, quant.FormatNone, []float32{1, 0, 0, 1})
	if e != nil {
		return e
	}
	session, e := nativecuda.Mount(l, nativecuda.Options{LibraryPath: *library})
	if e != nil {
		fmt.Printf("SKIP: CUDA unavailable (%v)\n", e)
		return nil
	}
	defer session.Close()
	x := &core.Tensor[float32]{Shape: []int{1, 2}, Data: []float32{.5, -.25}}
	post, e := session.Forward(x)
	if e != nil {
		return e
	}
	fmt.Printf("backend=cuda version=%s\ndevice=%s\ninput=%v output=%v shape=%v\n", session.Version(), session.Identity(), x.Data, post.Data, post.Shape)
	return nil
}

func cudaPlatformOK() bool {
	return runtime.GOARCH == "amd64" && (runtime.GOOS == "linux" || runtime.GOOS == "windows")
}
