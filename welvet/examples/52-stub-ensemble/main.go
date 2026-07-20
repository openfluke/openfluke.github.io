package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/ensemble"
)

func main() {
	votes := [][]int{{0, 1, 1}, {0, 0, 1}, {0, 1, 0}}
	fmt.Println(ensemble.MajorityVote(votes))
}
