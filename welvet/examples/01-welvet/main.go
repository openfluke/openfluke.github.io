package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
)

func main() {
	l, err := dense.New(8, 4, core.ActivationReLU, core.DTypeFloat32)
	if err != nil {
		panic(err)
	}
	l.Exec.Backend = core.BackendCPUTiled

	x := core.NewTensor[float32](1, 8)
	for i := range x.Data {
		x.Data[i] = float32(i) * 0.1
	}
	pre, post, err := dense.Forward(l, x)
	if err != nil {
		panic(err)
	}
	fmt.Println("pre", len(pre.Data), "post", len(post.Data))
}
