package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/observer"
)

func main() {
	var obs observer.Observer = &observer.ConsoleObserver{}
	fmt.Printf("%T\n", obs)
}
