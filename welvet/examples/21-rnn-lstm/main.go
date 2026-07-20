package main

import (
	"fmt"

	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/lstm"
	"github.com/openfluke/welvet/layers/rnn"
)

func main() {
	r, _ := rnn.New(rnn.Config{InputSize: 8, HiddenSize: 16, SeqLen: 4})
	l, _ := lstm.New(lstm.Config{InputSize: 8, HiddenSize: 16, SeqLen: 4})
	x := core.NewTensor[float32](1, 4, 8)
	_, yr, _ := rnn.Forward(r, x)
	_, yl, _ := lstm.Forward(l, x)
	fmt.Println(len(yr.Data), len(yl.Data))
}
