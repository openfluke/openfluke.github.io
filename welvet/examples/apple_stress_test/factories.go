package main

import (
	"fmt"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/mamba"
	"github.com/openfluke/welvet/quant"
	"github.com/openfluke/welvet/runtime/dispatch"
	"github.com/openfluke/welvet/systems/dna"
	"math"
	"reflect"
)

var layerNames = []string{"mamba"}

type geometry struct{ Dim, Batch, Seq int }

func profile(name string) (geometry, error) {
	switch name {
	case "smoke":
		return geometry{8, 2, 3}, nil
	case "medium":
		return geometry{128, 4, 16}, nil
	case "stress":
		return geometry{512, 8, 32}, nil
	case "large":
		return geometry{1024, 16, 64}, nil
	}
	return geometry{}, fmt.Errorf("unknown profile %q", name)
}
func pattern(n int, seed uint64, scale float64) []float32 {
	v := make([]float32, n)
	for i := range v {
		seed ^= seed << 13
		seed ^= seed >> 7
		seed ^= seed << 17
		v[i] = float32((float64(seed%20001)/10000 - 1) * scale)
	}
	return v
}
func tensor(shape []int, seed uint64) *core.Tensor[float32] {
	t := core.NewTensor[float32](shape...)
	copy(t.Data, pattern(t.Len(), seed, .2))
	return t
}
func makeLayer(name string, dt core.DType, g geometry, seed uint64) (any, []int, error) {
	if name != "mamba" {
		return nil, nil, fmt.Errorf("this stress test only runs Mamba")
	}
	d := g.Dim
	w := func(n, fan int) []float32 { return pattern(n, seed, .5/math.Sqrt(float64(fan))) }
	op, err := mamba.NewConfigured(mamba.Config{DModel: d, DState: 8, Expand: 2, SeqLen: g.Seq}, dt, quant.FormatNone, w(4*d*d, d), pattern(2*d*d, seed+1, .5/math.Sqrt(float64(2*d))), nil, nil)
	return op, []int{g.Batch, g.Seq, d}, err
}
func coreMeta(op any) core.Layer { return op.(*mamba.Layer).Core }
func parameters(ops []any) ([]float32, error) {
	var out []float32
	for _, op := range ops {
		v, e := dna.FlattenOp(op)
		if e != nil {
			return nil, e
		}
		out = append(out, v...)
		l := op.(*mamba.Layer)
		out = append(out, l.ALog...)
		out = append(out, l.DSkip...)
	}
	return out, nil
}
func setBackend(op any, b core.Backend) {
	exec := core.ExecConfig{Backend: b, BackendPinned: true, MultiCore: true, TileSize: 32}
	reflect.ValueOf(op).Elem().FieldByName("Exec").Set(reflect.ValueOf(exec))
	dispatch.SyncExec(op, exec)
}
