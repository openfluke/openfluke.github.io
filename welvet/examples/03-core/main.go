package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
)

func main() {
	t := core.NewTensor[float32](2, 4)
	t.Data[0] = -1
	t.Data[0] = core.Activate(t.Data[0], core.ActivationReLU) // 0

	cfg := core.ExecConfig{
		Backend:   core.BackendSIMD,
		MultiCore: true,
		TileSize:  32,
	}
	fmt.Println(cfg.Backend.String(), core.ParseDType("bfloat16"), t.Len())
}
