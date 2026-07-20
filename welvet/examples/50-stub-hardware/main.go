package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/hardware"
)

func main() {
	a := hardware.Audit()
	fmt.Printf("%+v\n", a.CPU)
	fmt.Println(hardware.Description())
}
