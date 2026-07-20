package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/dense"
)

func main() {
	g := architecture.NewGrid(1, 1, 2, 1) // depth, rows, cols, layersPerCell
	l, _ := dense.New(4, 4, core.ActivationLinear, core.DTypeFloat32)
	if err := dense.Place(g, 0, 0, 0, 0, l); err != nil {
		panic(err)
	}
	_ = g.SetRemoteLink(0, 0, 1, 0, 0, 0, 0, 0) // spatial feedback hop
	fmt.Println("cells", len(g.HopOrder()))
}
