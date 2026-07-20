package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/embedding"
)

func main() {
	l, err := embedding.New(embedding.Config{VocabSize: 32, EmbeddingDim: 16, SeqLen: 4})
	if err != nil {
		panic(err)
	}
	ids := core.NewTensor[float32](1, 4)
	ids.Data[0], ids.Data[1] = 3, 7
	_, y, err := embedding.Forward(l, ids)
	fmt.Println(y.Shape, err)
}
