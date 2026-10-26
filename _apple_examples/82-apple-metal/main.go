// Mount one layer, execute its forward pass on Apple Metal, and release it.
package main

import (
	"flag"
	"fmt"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
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
	library := flag.String("library", "", "optional libwelvet_metal.dylib path")
	flag.Parse()
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		fmt.Println("SKIP: Apple Metal needs Apple silicon macOS and cgo")
		return nil
	}
	l, e := dense.NewConfigured(2, 2, core.ActivationLinear, core.DTypeFloat32, quant.FormatNone, []float32{1, 0, 0, 1})
	if e != nil {
		return e
	}
	session, e := nativemetal.Mount(l, nativemetal.Options{LibraryPath: *library})
	if e != nil {
		return e
	}
	defer session.Close()
	x := &core.Tensor[float32]{Shape: []int{1, 2}, Data: []float32{.5, -.25}}
	post, e := session.Forward(x)
	if e != nil {
		return e
	}
	fmt.Printf("backend=metal version=%s\ndevice=%s\ninput=%v output=%v shape=%v\n", session.Version(), session.Identity(), x.Data, post.Data, post.Shape)
	return nil
}
