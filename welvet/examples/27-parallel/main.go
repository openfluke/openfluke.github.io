package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/parallel"
)

func main() {
	l, err := parallel.New(parallel.Config{
		Dim: 16, Branches: 2, OutFeat: 16, Combine: parallel.CombineConcat,
	})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 16)
	_, y, err := parallel.Forward(l, x)
	fmt.Println(y.Shape, err)
}
