package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/mha"
)

func main() {
	l, err := mha.New(mha.DecoderCausal(32, 4, 4))
	if err != nil {
		panic(err)
	}
	x := core.NewTensor[float32](1, 8, 32) // [batch, seq, dim]
	pre, post, err := mha.Forward(l, x)
	fmt.Println(pre != nil, len(post.Data), err)
}
