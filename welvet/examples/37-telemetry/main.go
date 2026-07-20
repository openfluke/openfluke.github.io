package main

import (
	"fmt"

	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/systems/telemetry"
)

func main() {
	g := architecture.NewGrid(1, 1, 1, 1)
	bp := telemetry.ExtractNetworkBlueprint(g, "demo")
	fmt.Printf("%+v\n", bp)
}
