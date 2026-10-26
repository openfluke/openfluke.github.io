// apple_stress_test exercises Go -> native C++/Metal bindings in-process.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/openfluke/welvet/core"
	"github.com/openfluke/welvet/layers/parallel"
	"github.com/openfluke/welvet/nativecpp"
	"github.com/openfluke/welvet/nativemetal"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type unsupported struct{ reason string }

func (e unsupported) Error() string { return e.reason }

type config struct {
	Profile, Layers, DTypes, Modes, Phase, Out   string
	Samples, Iterations, Steps, Warmup, MaxCases int
	LR                                           float64
	List, Strict                                 bool
	Context                                      context.Context `json:"-"`
}

func selected(value, name string) bool {
	if strings.EqualFold(value, "all") {
		return true
	}
	for _, s := range strings.Split(value, ",") {
		if strings.EqualFold(strings.TrimSpace(s), name) {
			return true
		}
	}
	return false
}
func modesFor(value string) ([]parallel.TrainMode, error) {
	var out []parallel.TrainMode
	for m := parallel.ModeNormalBP; m <= parallel.ModeStepTweenSplitSparse; m++ {
		if selected(value, m.String()) {
			out = append(out, m)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no modes selected; use -list for the 29 names")
	}
	return out, nil
}
func dtypesFor(value string) ([]core.DType, error) {
	var out []core.DType
	for _, dt := range core.AllDTypes {
		if selected(value, dt.String()) {
			out = append(out, dt)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no dtypes selected; use -list")
	}
	return out, nil
}
func main() { os.Exit(run()) }
func run() int {
	c := config{}
	flag.StringVar(&c.Profile, "profile", "stress", "smoke, medium, stress, large; workload geometry, not a speed guarantee")
	flag.StringVar(&c.Layers, "layers", "mamba", "comma-separated layer names, or all")
	flag.StringVar(&c.DTypes, "dtypes", "all", "comma-separated parameter dtype names, or all (Float32 activations)")
	flag.StringVar(&c.Modes, "modes", "all", "comma-separated training modes, or all 29")
	flag.StringVar(&c.Phase, "phase", "all", "all, forward, train")
	flag.StringVar(&c.Out, "out", "", "report directory; default results/<timestamp>")
	flag.IntVar(&c.Samples, "samples", 3, "independent seeded repeats; median timing and replay check")
	flag.IntVar(&c.Iterations, "iters", 3, "forward calls per sample")
	flag.IntVar(&c.Steps, "steps", 8, "training calls per sample; enough to pop the Line pipeline")
	flag.IntVar(&c.Warmup, "warmup", 1, "untimed warm-up calls on a separate disposable model")
	flag.IntVar(&c.MaxCases, "max-cases", 0, "optional row limit; zero runs the full selected matrix")
	flag.Float64Var(&c.LR, "lr", .0001, "learning rate shared by all backends")
	flag.BoolVar(&c.List, "list", false, "print layers, dtype names, modes and selected case plan without loading native libraries")
	flag.BoolVar(&c.Strict, "strict", false, "return failure for below-industry deviations, as well as execution errors")
	flag.Parse()
	if math.IsNaN(c.LR) || math.IsInf(c.LR, 0) {
		fmt.Fprintln(os.Stderr, "learning rate must be finite")
		return 1
	}
	if c.Samples < 1 || c.Iterations < 1 || c.Steps < 1 || c.Warmup < 0 || c.MaxCases < 0 {
		fmt.Fprintln(os.Stderr, "samples/iters/steps must be positive; warmup/max-cases nonnegative")
		return 1
	}
	if c.Phase != "all" && c.Phase != "forward" && c.Phase != "train" {
		fmt.Fprintln(os.Stderr, "invalid phase")
		return 1
	}
	g, e := profile(c.Profile)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		return 1
	}
	dtypes, e := dtypesFor(c.DTypes)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		return 1
	}
	modes, e := modesFor(c.Modes)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		return 1
	}
	var layers []string
	for _, l := range layerNames {
		if selected(c.Layers, l) {
			layers = append(layers, l)
		}
	}
	if len(layers) == 0 {
		fmt.Fprintln(os.Stderr, "no layers selected")
		return 1
	}
	// Validate spelling rather than silently reducing a requested matrix.
	for value, allowed := range map[string][]string{c.Layers: layerNames, c.DTypes: dtypeNames(), c.Modes: modeNames()} {
		if strings.EqualFold(value, "all") {
			continue
		}
		for _, token := range strings.Split(value, ",") {
			found := false
			for _, a := range allowed {
				if strings.EqualFold(strings.TrimSpace(token), a) {
					found = true
					break
				}
			}
			if !found {
				fmt.Fprintln(os.Stderr, "unknown selection:", token)
				return 1
			}
		}
	}
	count := 0
	if c.Phase == "all" || c.Phase == "forward" {
		count += len(layers) * len(dtypes)
	}
	if c.Phase == "all" || c.Phase == "train" {
		count += len(layers) * len(dtypes) * len(modes)
	}

	if c.List {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{"layers": layerNames, "dtypes": dtypeNames(), "modes": modeNames(), "geometry": g, "selected_rows": count, "activation_dtype": "float32"})
		return 0
	}
	if !nativecpp.NativeCPUAvailable() {
		fmt.Fprintln(os.Stderr, "C++ binding unavailable; use a normal cgo build with its native artifact, without -tags welvet_metal")
		return 1
	}
	probe, _, e := makeLayer("mamba", core.DTypeFloat32, geometry{Dim: 8, Batch: 1, Seq: 2}, 123)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		return 1
	}
	metal, e := nativemetal.Mount(probe)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		return 1
	}
	identity, metalVersion := metal.Identity(), metal.Version()
	metal.Close()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	c.Context = ctx
	if c.Out == "" {
		c.Out = filepath.Join("results", time.Now().Format("20060102-150405"))
	}
	metadata := map[string]any{"started": time.Now().Format(time.RFC3339), "go_version": runtime.Version(), "goos": runtime.GOOS, "goarch": runtime.GOARCH, "cpu_threads": runtime.GOMAXPROCS(0), "cpp_version": nativecpp.Version(), "metal_version": metalVersion, "metal_identity": identity, "config": c, "geometry": g, "planned_rows": count, "metric": "median wall time per public API call; setup and final pull separate; det=exact bits / within tolerance", "cpp_path": "Go scheduling/SGD/proxy with real C ABI forward/backward; host Metacognition rules retained", "metal_path": "official Go nativemetal persistent session; no external engine executables"}
	report, e := openReport(c.Out, metadata)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		return 1
	}
	defer report.close()
	fmt.Printf("Apple stress test: %d selected rows; profile=%s geometry=%+v\nC++ %s | Metal %s %s\nreports: %s\n", count, c.Profile, g, nativecpp.Version(), metalVersion, identity, c.Out)
	fmt.Println("det percentages = actual output exact bits / values within tolerance; negative latency reduction means slower. Setup and final pulls are in CSV/JSONL.")
	limit := func() bool { return ctx.Err() != nil || (c.MaxCases > 0 && report.rows >= c.MaxCases) }
	add := func(phase, name string, dt core.DType, mode parallel.TrainMode) bool {
		if limit() {
			return false
		}
		ms := make([]measurement, 0, 3)
		for _, b := range []string{"go", "cpp", "metal"} {
			ms = append(ms, measure(name, dt, g, b, mode, phase == "train", c))
			runtime.GC()
		}
		r := makeRow(phase, name, dt.String(), "—", c.Profile, ms)
		if phase == "train" {
			r.Mode = mode.String()
		}
		allUnsupported := true
		for _, m := range ms {
			if !m.Unsupported {
				allUnsupported = false
			}
		}
		if allUnsupported {
			r.Status = "unsupported"
		}
		if e := report.add(r); e != nil {
			fmt.Fprintln(os.Stderr, e)
			cancel()
			return false
		}
		return true
	}
	if c.Phase == "all" || c.Phase == "forward" {
		for _, l := range layers {
			for _, dt := range dtypes {
				if !add("forward", l, dt, parallel.ModeNormalBP) {
					break
				}
			}
			if limit() {
				break
			}
		}
	}
	if c.Phase == "all" || c.Phase == "train" {
		for _, l := range layers {
			for _, dt := range dtypes {
				for _, m := range modes {
					if !add("train", l, dt, m) {
						break
					}
				}
				if limit() {
					break
				}
			}
			if limit() {
				break
			}
		}
	}
	if e := report.summary(ctx.Err() != nil); e != nil {
		fmt.Fprintln(os.Stderr, e)
		return 1
	}
	fmt.Printf("Completed %d/%d rows: errors=%d deviations=%d. Reports: %s\n", report.rows, count, report.errors, report.deviations, c.Out)
	if ctx.Err() != nil {
		return 130
	}
	if report.errors > 0 || (c.Strict && report.deviations > 0) {
		return 1
	}
	return 0
}
func dtypeNames() []string {
	var out []string
	for _, d := range core.AllDTypes {
		out = append(out, d.String())
	}
	return out
}
func modeNames() []string {
	var out []string
	for m := parallel.ModeNormalBP; m <= parallel.ModeStepTweenSplitSparse; m++ {
		out = append(out, m.String())
	}
	return out
}

// Assert the error's contract separately from execution failures.
func isUnsupported(e error) bool { var u unsupported; return errors.As(e, &u) }
