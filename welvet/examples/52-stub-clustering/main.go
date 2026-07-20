package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/stub/clustering"
)

func main() {
	pts := []*core.Tensor[float32]{
		core.NewTensor[float32](2),
		core.NewTensor[float32](2),
	}
	pts[0].Data[0], pts[0].Data[1] = 0, 0
	pts[1].Data[0], pts[1].Data[1] = 1, 1
	cent, assign := clustering.KMeansCluster(pts, 2, 10, false)
	fmt.Println(len(cent), assign)
}
