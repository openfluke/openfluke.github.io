package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/layers/parallel"
	"github.com/openfluke/welvet/stub/grafting"
)

func main() {
	a, b := architecture.NewGrid(1, 1, 1, 1), architecture.NewGrid(1, 1, 1, 1)
	out, err := grafting.GraftGrids([]*architecture.Grid{a, b}, parallel.CombineConcat)
	fmt.Println(out != nil, err)
}
