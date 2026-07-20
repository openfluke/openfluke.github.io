package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/fountain"
)

func main() {
	src := [][]byte{make([]byte, 64), make([]byte, 64), make([]byte, 64)}
	enc, err := fountain.NewLTEncoder(src, 42)
	fmt.Printf("%T %v\n", enc, err)
}
