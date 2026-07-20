package main

import (
	"fmt"

	"github.com/openfluke/welvet/webgpu"
)

func main() {
	if !webgpu.Available() {
		fmt.Println("no adapter:", webgpu.InitError())
		return
	}
	y := make([]float32, 2)
	err := webgpu.DenseGEMV(
		[]float32{1, 0, 0, 1},
		[]float32{1.5, 2.5},
		y, 1, 2, 2,
	)
	fmt.Println(y, err, webgpu.AdapterName())
}
