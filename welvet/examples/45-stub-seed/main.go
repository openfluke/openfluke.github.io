package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/seed"
)

func main() {
	s := seed.From("demo-net", 0)
	rng := seed.New(s)
	fmt.Println(rng.NormFloat64())
}
