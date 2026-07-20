package main

import (
	"fmt"

	"github.com/openfluke/welvet/fusedgpu"
	"github.com/openfluke/welvet/model/transformer"
)

func main() {
	m, err := transformer.LoadEntity("model.entity")
	if err != nil {
		fmt.Println("need an ENTITY file:", err)
		return
	}
	spec, err := m.ExportFusedGPUSpec()
	if err != nil {
		panic(err)
	}
	eng, err := fusedgpu.NewFromSpec(spec)
	if err != nil {
		panic(err)
	}
	defer eng.Close()
	logits, err := eng.AppendTokens([]uint32{1, 2, 3})
	fmt.Println(len(logits), err)
}
