package main

import (
	"fmt"

	"github.com/openfluke/welvet/layers/seqmix"
)

func main() {
	c := seqmix.Contract{Kind: seqmix.KindAttention, DModel: 64, MaxT: 2048}
	fmt.Println(c.Kind.String()) // attention
}
