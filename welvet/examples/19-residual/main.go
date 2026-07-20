package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/residual"
)

func main() {
	l, err := residual.New(residual.Config{Dim: 16, Depth: 2})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 16)
	_, y, err := residual.Forward(l, x)
	fmt.Println(len(y.Data), err)
}
