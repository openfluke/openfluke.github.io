package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/softmax"
)

func main() {
	l, err := softmax.New(softmax.Config{Dim: 4, Kind: softmax.KindStandard})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 4)
	copy(x.Data, []float32{2, 1, 0.1, -1})
	_, y, err := softmax.Forward(l, x)
	fmt.Println(y.Data, err)
}
