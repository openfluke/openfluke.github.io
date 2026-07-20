package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/cnn2"
)

func main() {
	l, err := cnn2.New(cnn2.Config{
		InChannels: 1, Filters: 4, Height: 8, Width: 8, Kernel: 3, Stride: 1,
	})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 1, 8, 8)
	_, y, err := cnn2.Forward(l, x)
	fmt.Println(y.Shape, err)
}
