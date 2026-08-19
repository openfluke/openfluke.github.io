package main

import (
	"fmt"

	"github.com/openfluke/welvet/layers/parallel"
)

func main() {
	named := parallel.AllNamedTrainModes()
	line := 0
	for _, m := range named {
		if m.IsLineStep() {
			line++
		}
	}
	fp, err := parallel.ParseTrainMode("stepfastproxy")
	fmt.Println("named", len(named), "linestep", line)
	fmt.Println("stepfastproxy", fp.Short(), err)
	fmt.Println(parallel.ShortTrainModeLegend)
}
