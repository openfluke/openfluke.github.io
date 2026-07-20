package main

import (
	"fmt"

	"github.com/openfluke/welvet/simd"
)

func main() {
	if !simd.SimdEnabled() {
		fmt.Println("SIMD not available on this arch — BackendSIMD must hard-error")
		return
	}
	acc := simd.DotTile([]float32{1, 2, 3, 4}, []float32{1, 1, 1, 1}, 0, 4, 0)
	fmt.Println(acc)
}
