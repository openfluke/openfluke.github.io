package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/universal"
)

func main() {
	g, err := universal.LoadUniversal("/path/to/snapshot")
	fmt.Println(g != nil, err)
}
