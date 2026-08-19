package main

import (
	"fmt"

	"github.com/openfluke/welvet/lucy"
)

func main() {
	board := lucy.BuildLPD([]lucy.Sample{
		{ID: "f32", Mode: "sgd", Acc: 90, Thru: 200, Avail: 40, Score: 100, RAMKiB: 1000},
		{ID: "int8", Mode: "sgd", Acc: 82, Thru: 180, Avail: 38, Score: 85, RAMKiB: 180},
		{ID: "bin", Mode: "sgd", Acc: 12, Thru: 400, Avail: 50, Score: 40, RAMKiB: 40},
	})
	top := board.Top[0]
	fmt.Printf("lead=%s band=%s LPD=%.2f trap=%s\n",
		top.ID, top.Band, top.LPD, board.Trap[0].ID)
}
