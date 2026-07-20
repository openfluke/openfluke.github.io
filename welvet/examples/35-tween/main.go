package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/systems/tween"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	st := tween.NewState[float32](g, tween.DefaultConfig())
	_, err := tween.Forward(g, st, core.NewTensor[float32](1, 4))
	fmt.Println(err)
}
