package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/metacognition"
)

func main() {
	l, err := metacognition.New(metacognition.Config{
		Dim: 16, Rules: metacognition.DefaultStabilityRules(),
	})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 16)
	_, y, err := metacognition.Forward(l, x)
	fmt.Println(len(y.Data), err)
}
