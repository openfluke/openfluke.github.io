package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/stub/serialization"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	raw, err := serialization.SerializeGrid(g)
	fmt.Println(len(raw), err)
}
