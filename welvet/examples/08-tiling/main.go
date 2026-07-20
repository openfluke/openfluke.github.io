package main

import (
	"fmt"

	"github.com/openfluke/welvet/tiling"
)

func main() {
	tile := tiling.CPUTile(0) // default 32
	mc := tiling.PreferMultiCore(8, 256, tile)
	wg := tiling.GPUWorkgroupsX(1024, 0)
	fmt.Println(tile, mc, wg)
}
