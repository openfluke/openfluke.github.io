package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/convt1"
)

func main() {
	l, err := convt1.New(convt1.Config{InChannels: 4, Filters: 2, SeqLen: 8, Kernel: 3, Stride: 2})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 4, 8)
	_, y, err := convt1.Forward(l, x)
	fmt.Println(y.Shape, err)
}
