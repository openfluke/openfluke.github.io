package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/stub/evaluation"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	in := []*core.Tensor[float32]{core.NewTensor[float32](1, 4)}
	rep, err := evaluation.EvaluateNetwork(g, in, []float64{0, 0, 0, 0})
	fmt.Println(rep, err)
}
