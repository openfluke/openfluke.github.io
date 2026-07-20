package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/layernorm"
)

func main() {
	l, err := layernorm.New(layernorm.Config{Dim: 8})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 8)
	_, y, err := layernorm.Forward(l, x)
	fmt.Println(len(y.Data), err)
}
