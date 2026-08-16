package main

import (
	"fmt"

	"github.com/openfluke/welvet/layers/parallel"
)

func main() {
	n := 0
	for _, m := range parallel.AllNamedTrainModes() {
		n++
		_ = m.String()
	}
	fp, err := parallel.ParseTrainMode("fastproxy")
	fmt.Println("named", n, "fastproxy", fp, err)
}
