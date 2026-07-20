package main

import (
	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/systems/tanhi"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	cfg := tanhi.ConfigFromGrid(g)
	tanhi.EmitSweep(cfg, "epoch-0") // non-blocking best-effort
}
