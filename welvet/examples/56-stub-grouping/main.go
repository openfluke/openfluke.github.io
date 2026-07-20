package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/grouping"
)

func main() {
	tensors := []grouping.DetectedTensor{
		{Name: "model.layers.0.self_attn.q_proj.weight"},
		{Name: "model.layers.0.self_attn.k_proj.weight"},
		{Name: "model.layers.0.self_attn.v_proj.weight"},
		{Name: "model.layers.0.self_attn.o_proj.weight"},
	}
	ok, hint := grouping.DetectMHA("block0", tensors, 64, 4)
	fmt.Println(ok, hint)
}
