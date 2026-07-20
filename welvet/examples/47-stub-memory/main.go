package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/stub/memory"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	fp := memory.FromGrid(g)
	fmt.Printf("%+v\n", fp)
}
