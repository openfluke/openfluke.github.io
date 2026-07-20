package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/kmeans"
)

func main() {
	l, err := kmeans.New(kmeans.Config{NumClusters: 4, FeatureDim: 8})
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 8)
	_, y, err := kmeans.Forward(l, x)
	fmt.Println(y.Shape, err)
}
