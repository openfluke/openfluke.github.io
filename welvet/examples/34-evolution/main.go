package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/systems/evolution"
)

func main() {
	a, b := architecture.NewGrid(1, 1, 1, 1), architecture.NewGrid(1, 1, 1, 1)
	child, err := evolution.SpliceDNA(a, b, evolution.DefaultSpliceConfig())
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = evolution.NEATMutate(child, evolution.DefaultNEATConfig(16))
	fmt.Println(err)
}
