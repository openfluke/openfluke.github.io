package main

import (
	"fmt"

	"github.com/openfluke/welvet/lucy"
)

func main() {
	a := lucy.SoftAccOne(0.72, 0.80) // sine scale 0.10
	p := lucy.SoftAccProb(0.91, 1.0) // class scale 1.0
	var snap lucy.Snapshot
	snap.AvgAccuracy = 80
	snap.SoftAcc = a
	snap.InferMs = 8
	snap.TrainMs = 2
	snap.Throughput = 12000
	lucy.Finalize(&snap, lucy.Options{AdaptWindows: 10})
	fmt.Printf("soft=%.1f class=%.1f avail=%.1f score=%.0f\n",
		a, p, snap.Availability, snap.Score)
}
