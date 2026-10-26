package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
)

type accuracy struct {
	Comparable    bool    `json:"comparable"`
	Count         int     `json:"count"`
	ExactPct      float64 `json:"exact_pct"`
	WithinPct     float64 `json:"within_tolerance_pct"`
	MaxAbs        float64 `json:"max_abs"`
	MaxNormalized float64 `json:"max_normalized_error"`
	NonFinite     int     `json:"non_finite"`
	Band          string  `json:"band"`
}

// ExactPct is the actual fraction of matching Float32 bit patterns. WithinPct
// counts abs(a-b)<=1e-5+1e-3*max(abs(a),abs(b)); these are not model accuracy.
func compare(a, b []float32) accuracy {
	m := accuracy{Count: len(a), Band: "not_comparable"}
	if len(a) != len(b) {
		return m
	}
	m.Comparable = true
	if len(a) == 0 {
		m.Band = "no_parameters"
		return m
	}
	exact, within := 0, 0
	for i, x := range a {
		av, bv := float64(x), float64(b[i])
		if math.IsNaN(av) || math.IsNaN(bv) || math.IsInf(av, 0) || math.IsInf(bv, 0) {
			m.NonFinite++
			continue
		}
		if math.Float32bits(x) == math.Float32bits(b[i]) {
			exact++
		}
		delta := math.Abs(av - bv)
		scale := math.Max(math.Abs(av), math.Abs(bv))
		if delta <= 1e-5+1e-3*scale {
			within++
		}
		m.MaxAbs = math.Max(m.MaxAbs, delta)
		m.MaxNormalized = math.Max(m.MaxNormalized, delta/(1+scale))
	}
	m.ExactPct = 100 * float64(exact) / float64(len(a))
	m.WithinPct = 100 * float64(within) / float64(len(a))
	m.Band = "below_industry"
	if m.NonFinite > 0 {
		m.Band = "non_finite"
	} else if exact == len(a) {
		m.Band = "exact100"
	} else if m.MaxNormalized <= 1e-5 {
		m.Band = "strict80"
	} else if m.MaxNormalized <= 1e-3 {
		m.Band = "industry60"
	}
	return m
}
func median(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	sort.Float64s(v)
	i := len(v) / 2
	if len(v)%2 == 0 {
		return (v[i-1] + v[i]) / 2
	}
	return v[i]
}

type backendResult struct {
	Backend        string   `json:"backend"`
	Path           string   `json:"path"`
	Error          string   `json:"error,omitempty"`
	SetupMS        float64  `json:"setup_ms"`
	SteadyUS       float64  `json:"steady_us_per_call"`
	PullMS         float64  `json:"final_pull_ms"`
	Calls          uint64   `json:"successful_cpp_calls"`
	Params         int      `json:"parameters"`
	Loss           float64  `json:"loss"`
	ReplayExactPct *float64 `json:"replay_exact_pct"`
}

func result(m measurement) backendResult {
	r := backendResult{Backend: m.Backend, Path: m.Path, Error: m.Error, SetupMS: m.SetupMS, SteadyUS: m.SteadyUS, PullMS: m.PullMS, Calls: m.Calls, Params: m.Params, Loss: m.Loss}
	if m.ReplayMeasured {
		v := m.ReplayExactPct
		r.ReplayExactPct = &v
	}
	return r
}

type pairResult struct {
	A, B                          string
	SpeedupX, LatencyReductionPct float64
	Output, Weights               accuracy
	LossAbsError                  float64
}
type row struct {
	Phase, Layer, DType, Mode, Profile, Status, Note string
	Results                                          []backendResult
	Pairs                                            []pairResult
}

func makeRow(phase, layer, dt, mode, profile string, ms []measurement) row {
	r := row{Phase: phase, Layer: layer, DType: dt, Mode: mode, Profile: profile, Status: "ok"}
	for _, m := range ms {
		r.Results = append(r.Results, result(m))
		if m.Error != "" {
			r.Status = "error"
		}
	}
	for _, ab := range [][2]int{{0, 1}, {0, 2}, {1, 2}} {
		a, b := ms[ab[0]], ms[ab[1]]
		p := pairResult{A: a.Backend, B: b.Backend}
		if a.Error == "" && b.Error == "" {
			if b.SteadyUS > 0 {
				p.SpeedupX = a.SteadyUS / b.SteadyUS
			}
			if a.SteadyUS > 0 {
				p.LatencyReductionPct = 100 * (1 - b.SteadyUS/a.SteadyUS)
			}
			p.Output = compare(a.Output, b.Output)
			p.Weights = compare(a.Weights, b.Weights)
			p.LossAbsError = math.Abs(a.Loss - b.Loss)
			if p.Output.NonFinite+p.Weights.NonFinite > 0 {
				r.Status = "error"
			}
			if p.Output.Band == "below_industry" || p.Weights.Band == "below_industry" {
				if r.Status == "ok" {
					r.Status = "deviation"
				}
			}
		}
		r.Pairs = append(r.Pairs, p)
	}
	return r
}

type writer struct {
	jsonl                    *os.File
	csvFile                  *os.File
	csv                      *csv.Writer
	rows, errors, deviations int
	dir                      string
}

func openReport(dir string, metadata any) (*writer, error) {
	if e := os.MkdirAll(dir, 0755); e != nil {
		return nil, e
	}
	j, e := os.Create(filepath.Join(dir, "results.jsonl"))
	if e != nil {
		return nil, e
	}
	f, e := os.Create(filepath.Join(dir, "results.csv"))
	if e != nil {
		j.Close()
		return nil, e
	}
	w := &writer{jsonl: j, csvFile: f, csv: csv.NewWriter(f), dir: dir}
	if e = json.NewEncoder(j).Encode(map[string]any{"type": "metadata", "data": metadata}); e != nil {
		w.close()
		return nil, e
	}
	e = w.csv.Write([]string{"phase", "layer", "dtype", "mode", "profile", "status", "backend", "path", "error", "setup_ms", "steady_us", "pull_ms", "cpp_calls", "parameters", "loss", "replay_exact_pct", "pair", "speedup_x", "latency_reduction_pct", "output_exact_pct", "output_within_tol_pct", "output_band", "weight_exact_pct", "weight_within_tol_pct", "weight_band", "max_output_abs", "max_weight_abs"})
	w.csv.Flush()
	if e != nil {
		w.close()
		return nil, e
	}
	return w, nil
}
func (w *writer) add(r row) error {
	w.rows++
	if r.Status == "error" {
		w.errors++
	}
	if r.Status == "deviation" {
		w.deviations++
	}
	if e := json.NewEncoder(w.jsonl).Encode(map[string]any{"type": "result", "data": r}); e != nil {
		return e
	}
	for i, b := range r.Results {
		p := r.Pairs[i]
		replay := "not_measured"
		if b.ReplayExactPct != nil {
			replay = fmt.Sprint(*b.ReplayExactPct)
		}
		if e := w.csv.Write([]string{r.Phase, r.Layer, r.DType, r.Mode, r.Profile, r.Status, b.Backend, b.Path, b.Error, fmt.Sprint(b.SetupMS), fmt.Sprint(b.SteadyUS), fmt.Sprint(b.PullMS), fmt.Sprint(b.Calls), fmt.Sprint(b.Params), fmt.Sprint(b.Loss), replay, p.A + "/" + p.B, fmt.Sprint(p.SpeedupX), fmt.Sprint(p.LatencyReductionPct), fmt.Sprint(p.Output.ExactPct), fmt.Sprint(p.Output.WithinPct), p.Output.Band, fmt.Sprint(p.Weights.ExactPct), fmt.Sprint(p.Weights.WithinPct), p.Weights.Band, fmt.Sprint(p.Output.MaxAbs), fmt.Sprint(p.Weights.MaxAbs)}); e != nil {
			return e
		}
	}
	w.csv.Flush()
	if e := w.csv.Error(); e != nil {
		return e
	}
	fmt.Printf("%-7s %-14s %-10s %-28s %s", r.Phase, r.Layer, r.DType, r.Mode, r.Status)
	for _, b := range r.Results {
		if b.Error != "" {
			fmt.Printf(" | %s ERROR: %s", b.Backend, b.Error)
		} else {
			fmt.Printf(" | %s %.1fµs", b.Backend, b.SteadyUS)
		}
	}
	for _, p := range r.Pairs {
		if !p.Output.Comparable {
			fmt.Printf(" | %s/%s n/a", p.A, p.B)
			continue
		}
		fmt.Printf(" | %s/%s %.4gx %+.1f%% det %.1f/%.1f%%", p.A, p.B, p.SpeedupX, p.LatencyReductionPct, p.Output.ExactPct, p.Output.WithinPct)
	}
	fmt.Println()
	return nil
}
func (w *writer) close() {
	if w.csv != nil {
		w.csv.Flush()
	}
	if w.csvFile != nil {
		_ = w.csvFile.Close()
	}
	if w.jsonl != nil {
		_ = w.jsonl.Close()
	}
}
func (w *writer) summary(cancelled bool) error {
	data, e := json.MarshalIndent(map[string]any{"rows": w.rows, "errors": w.errors, "deviations": w.deviations, "cancelled": cancelled}, "", "  ")
	if e != nil {
		return e
	}
	return os.WriteFile(filepath.Join(w.dir, "summary.json"), data, 0644)
}
