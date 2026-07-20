package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/stub/introspection"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	methods, err := introspection.GetMethods(g)
	fmt.Println(len(methods), err)
}
