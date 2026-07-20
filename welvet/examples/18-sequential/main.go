package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/sequential"
)

func main() {
	l, err := sequential.New(sequential.Config{Dim: 16, Depth: 3})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 16)
	_, y, err := sequential.Forward(l, x)
	fmt.Println(len(y.Data), err)
}
