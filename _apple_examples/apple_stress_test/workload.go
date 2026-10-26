package main

import (
	"fmt"
	"github.com/openfluke/welvet/architecture"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/parallel"
	"github.com/openfluke/welvet/nativecpp"
	"github.com/openfluke/welvet/nativemetal"
	"github.com/openfluke/welvet/runtime/dispatch"
	"github.com/openfluke/welvet/runtime/step"
	"github.com/openfluke/welvet/runtime/training"
	"math"
	"time"
)

type workload struct {
	op            any
	ops           []any
	stack         *parallel.Stack
	grid          *architecture.Grid
	input, target *core.Tensor[float32]
	mode          parallel.TrainMode
	session       *nativemetal.Session
	ticks         int
}

func buildWork(name string, dt core.DType, g geometry, backend string, mode parallel.TrainMode, train bool) (*workload, error) {
	primary, shape, e := makeLayer(name, dt, g, 0x1234567)
	if e != nil {
		return nil, e
	}
	w := &workload{op: primary, ops: []any{primary}, mode: mode, ticks: 4}
	w.input = tensor(shape, 0x7413)
	outShape := append([]int(nil), shape...)
	if train {
		head, _, err := makeLayer(name, dt, g, 0x2345678)
		if err != nil {
			return nil, err
		}
		w.ops = append(w.ops, head)
		w.stack, e = parallel.NewStack(w.ops...)
		if e != nil {
			return nil, e
		}
		w.target = tensor(outShape, 0x555555)
		if mode.RequiresGrid() {
			w.grid = architecture.NewGrid(3, 1, 1, 1)
			w.grid.Cells[0].Layer.IsDisabled = true
			for i, op := range w.ops {
				w.grid.Cells[i+1].Layer = coreMeta(op)
				w.grid.Cells[i+1].Op = op
			}
		}
	}
	b := core.BackendCPUTiled
	if backend == "cpp" {
		b = core.BackendNativeCPU
	}
	for _, op := range w.ops {
		setBackend(op, b)
	}
	// Keep Go's graph scheduler explicit for C++ credit/Line modes: native
	// forward/backward calls are checked, while scheduling/SGD/proxies can be Go.
	if backend == "metal" {
		var model any = primary
		if train {
			model = w.stack
			if w.grid != nil {
				model = w.grid
			}
		}
		w.session, e = nativemetal.Mount(model)
		if e != nil {
			return nil, e
		}
	}
	return w, nil
}
func (w *workload) close() {
	if w.session != nil {
		_ = w.session.Close()
	}
}
func (w *workload) forward() (*core.Tensor[float32], error) {
	if w.session != nil {
		return w.session.Forward(w.input)
	}
	if w.stack != nil {
		_, o, e := parallel.ForwardStack(w.stack, w.input)
		return o, e
	}
	_, o, e := dispatch.Forward(w.op, w.input)
	return o, e
}
func (w *workload) train(lr float64) (float64, error) {
	if w.session != nil {
		return w.session.TrainMSE(w.input, w.target, w.mode, lr, w.ticks)
	}
	if w.grid == nil {
		return parallel.TrainStackMSE(w.stack, w.input, w.target, w.mode, lr)
	}
	return meshTrain(w, lr)
}
func (w *workload) pull() error {
	if w.session != nil {
		return w.session.SyncToGo()
	}
	for _, op := range w.ops {
		if s, ok := op.(interface{ SyncWeightsFromNative() error }); ok {
			if e := s.SyncWeightsFromNative(); e != nil {
				return e
			}
		}
	}
	return nil
}
func meshTrain(w *workload, lr float64) (float64, error) {
	g := w.grid
	s := step.New[float32](g)
	for tick := 0; tick < w.ticks; tick++ {
		s.SetInput(w.input)
		if _, e := step.Forward(g, s, true); e != nil {
			return 0, e
		}
	}
	post := s.LayerData[len(g.Cells)-1]
	if post == nil {
		return 0, fmt.Errorf("mesh produced no output")
	}
	loss, e := training.MSE(post, w.target)
	if e != nil {
		return 0, e
	}
	gy, e := training.MSEGrad(post, w.target)
	if e != nil {
		return 0, e
	}
	if w.mode == parallel.ModeMeshBP || w.mode == parallel.ModeMeshTweenChain {
		_, grads, e := step.Backward(g, s, gy)
		if e != nil {
			return 0, e
		}
		for i, pair := range grads {
			if pair[1] != nil && len(pair[1].Data) > 0 {
				if e = dispatch.ApplyGradSGD(g.Cells[i].Op, pair[1], lr); e != nil {
					return 0, e
				}
			}
		}
		return loss, nil
	}
	local := parallel.ModeTween
	switch w.mode {
	case parallel.ModeMeshTweenSplit:
		local = parallel.ModeTweenSplit
	case parallel.ModeMeshTweenAlt:
		local = parallel.ModeTweenAlt
	case parallel.ModeMeshTweenSplitFastProxy:
		local = parallel.ModeTweenSplitFastProxy
	case parallel.ModeMeshTweenSplitSparse:
		local = parallel.ModeTweenSplitSparse
	}
	last := len(s.HistoryIn) - 1
	for i := len(g.Cells) - 1; i >= 0; i-- {
		c := &g.Cells[i]
		if c.Layer.IsDisabled {
			continue
		}
		x, pre := s.HistoryIn[last][i], s.HistoryPre[last][i]
		if x == nil || pre == nil {
			continue
		}
		wrapper, e := parallel.NewStack(c.Op)
		if e != nil {
			return 0, e
		}
		// Preserve the cell's already assigned backend while using public Go credit scheduling.
		if e = parallel.TrainStack(wrapper, gy, x, pre, local, lr); e != nil {
			return 0, e
		}
		next := core.NewTensor[float32](x.Shape...)
		if len(gy.Data) == len(next.Data) {
			copy(next.Data, gy.Data)
		}
		gy = next
	}
	return loss, nil
}

type measurement struct {
	Backend, Path, Error      string
	Unsupported               bool
	SetupMS, SteadyUS, PullMS float64
	Calls                     uint64
	Output, Weights           []float32
	Loss, ReplayExactPct      float64
	ReplayMeasured            bool
	Params                    int
}

func measure(name string, dt core.DType, g geometry, backend string, mode parallel.TrainMode, train bool, c config) (m measurement) {
	m.ReplayExactPct = 100
	m.Backend = backend
	m.Path = "Go CPU"
	if backend == "cpp" {
		m.Path = "Go -> C ABI forward/backward; Go scheduling/SGD/proxy"
	}
	if backend == "metal" {
		m.Path = "Go -> mounted Metal session"
	}
	var durations, setups, pulls []float64
	for sample := -1; sample < c.Samples; sample++ {
		if c.Context.Err() != nil {
			m.Error = "cancelled"
			return
		}
		start := time.Now()
		w, e := buildWork(name, dt, g, backend, mode, train)
		if e != nil {
			m.Error = e.Error()
			m.Unsupported = isUnsupported(e)
			return
		}
		setup := float64(time.Since(start)) / float64(time.Millisecond)
		n := c.Iterations
		if sample < 0 {
			n = c.Warmup
			if n == 0 {
				w.close()
				continue
			}
		}
		if train {
			n = c.Steps
			if sample < 0 {
				n = max(c.Warmup, 1)
			}
		}
		if sample >= 0 {
			warmStart := time.Now()
			for j := 0; j < c.Warmup; j++ {
				if _, e = w.forward(); e != nil {
					w.close()
					m.Error = e.Error()
					return
				}
			}
			setup += float64(time.Since(warmStart)) / float64(time.Millisecond)
		}
		count := nativecpp.CompletedCalls()
		start = time.Now()
		var post *core.Tensor[float32]
		var loss float64
		for i := 0; i < n; i++ {
			if c.Context.Err() != nil {
				e = c.Context.Err()
				break
			}
			if train {
				loss, e = w.train(c.LR)
			} else {
				post, e = w.forward()
			}
			if e != nil {
				break
			}
		}
		elapsed := time.Since(start)
		calls := nativecpp.CompletedCalls() - count
		if e != nil {
			w.close()
			m.Error = e.Error()
			return
		}
		if backend == "cpp" && calls == 0 {
			w.close()
			m.Error = "no successful C ABI calls: refusing to label Go fallback as C++"
			return
		}
		if sample >= 0 {
			start = time.Now()
			e = w.pull()
			pull := float64(time.Since(start)) / float64(time.Millisecond)
			if e != nil {
				w.close()
				m.Error = e.Error()
				return
			}
			if train {
				post, e = w.forward()
				if e != nil {
					w.close()
					m.Error = e.Error()
					return
				}
			}
			weights, e := parameters(w.ops)
			if e != nil {
				w.close()
				m.Error = e.Error()
				return
			}
			if math.IsNaN(loss) || math.IsInf(loss, 0) {
				w.close()
				m.Error = "non-finite loss"
				return
			}
			for _, values := range [][]float32{post.Data, weights} {
				for _, v := range values {
					if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
						w.close()
						m.Error = "non-finite output or parameter"
						return
					}
				}
			}
			if sample == 0 {
				m.Output = append([]float32(nil), post.Data...)
				m.Weights = weights
				m.Loss = loss
				m.Params = len(weights)
			} else {
				replay := compare(append(append([]float32{}, m.Output...), m.Weights...), append(append([]float32{}, post.Data...), weights...))
				if replay.NonFinite > 0 {
					w.close()
					m.Error = "non-finite replay"
					return
				}
				m.ReplayMeasured = true
				if replay.ExactPct < m.ReplayExactPct {
					m.ReplayExactPct = replay.ExactPct
				}
			}
			durations = append(durations, float64(elapsed)/float64(time.Microsecond)/float64(n))
			setups = append(setups, setup)
			pulls = append(pulls, pull)
			m.Calls += calls
		}
		w.close()
	}
	m.SteadyUS = median(durations)
	m.SetupMS = median(setups)
	m.PullMS = median(pulls)
	return
}
