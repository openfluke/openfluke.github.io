package main

import (
	"fmt"

	"github.com/openfluke/welvet/model/hf"
)

func main() {
	info, err := hf.InspectSnapshot("/path/to/hf-snapshot")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(info)
}
