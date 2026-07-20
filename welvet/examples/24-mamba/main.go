package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/mamba"
)

func main() {
	l, err := mamba.New(mamba.Config{DModel: 32, DState: 16, Expand: 2, SeqLen: 8})
	if err != nil {
		fmt.Println(err)
		return
	}
	x := core.NewTensor[float32](1, 8, 32)
	_, y, err := mamba.Forward(l, x)
	fmt.Println(y != nil, err)
}
