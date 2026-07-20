package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/swiglu"
)

func main() {
	l, err := swiglu.New(swiglu.DefaultFFN(64))
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 64)
	_, post, err := swiglu.Forward(l, x)
	fmt.Println(len(post.Data), err)
}
