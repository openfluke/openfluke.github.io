package main

import (
	"fmt"

	"github.com/openfluke/welvet/layers/gdn"
)

func main() {
	l, err := gdn.New(gdn.Config{
		HiddenSize: 64, NumKeyHeads: 4, NumValueHeads: 4,
		KeyHeadDim: 16, ValueHeadDim: 16, ConvKernel: 4,
	})
	if err != nil {
		fmt.Println("config rejected:", err)
		return
	}
	l.Reset()
	x := make([]float32, 64)
	y := make([]float32, 64)
	fmt.Println(l.ForwardDecode(x, y))
}
