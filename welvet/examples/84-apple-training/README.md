# 84. Apple Metal — training modes

**Part:** XII · Apple acceleration (welvet.metal)  
**Package:** `nativemetal + parallel + architecture`    
**Status:** ok — verified macOS/arm64 on Apple M5

## When

Train a mounted model using native Metal credit and optimizer kernels.

## Where

From `example/welvet/84-apple-training`:

```sh
go run .
go run . -mode StepBP -steps 8
go run . -mode MeshTweenChain
go run . -mode all
```

## Why

Keep one session alive for Step pipeline state; use an actual Grid for Mesh modes.

## What

Runs a two-layer Dense model in NormalBP by default. All 29 update modes are selectable. Mesh modes use a disabled ingress plus two active Grid cells and four clock ticks per call.

## Requirements

Apple silicon macOS, Go with cgo, and the installed independent Metal dylib.
Set `WELVET_METAL_LIBRARY` or install it from the Metal checkout:

```sh
cd /path/to/welvet.metal
./scripts/macos_install_go.sh
```

Use the Go v1.3.0 checkout containing `nativemetal`. Examples 82–86 can also use
`go run -tags welvet_metal .` to exclude the independent C++ binding. Example 87
needs the normal build to keep both native backends. Unsupported operating
systems print SKIP. Missing libraries on a supported host report an error.

These are usage examples with small shapes. For speed measurements, use
[`../apple_stress_test`](../apple_stress_test/README.md). The results are numerical
CPU/GPU comparisons, not model prediction accuracy.

Backend labels in routed examples: `cpu-tiled` is Go CPU, `native_cpu` is
welvet.cpp, and `metal` is the independent Apple GPU backend.

## Captured sample

```text
mode=NormalBP schedule=Stack calls=8 loss=6.7954462e-06 weights=synced-to-Go
```
