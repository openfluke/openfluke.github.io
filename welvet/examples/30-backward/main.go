package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
	"github.com/openfluke/welvet/runtime/backward"
	"github.com/openfluke/welvet/runtime/forward"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	l, _ := dense.New(4, 4, core.ActivationLinear, core.DTypeFloat32)
	_ = dense.Place(g, 0, 0, 0, 0, l)
	fwd, _ := forward.Forward(g, core.NewTensor[float32](1, 4))
	bwd, err := backward.Backward(fwd, core.NewTensor[float32](1, 4))
	fmt.Println(bwd != nil, err)
}
