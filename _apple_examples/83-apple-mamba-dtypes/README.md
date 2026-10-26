# 83. Apple Metal — Mamba parameter dtypes

**Part:** XII · Apple acceleration (welvet.metal)  
**Package:** `nativemetal + mamba`    
**Status:** ok — verified macOS/arm64 on Apple M5

## When

Choose native parameter storage independently of activation tensors.

## Where

From `example/welvet/83-apple-mamba-dtypes`:

```sh
go run .
go run . -dtype float16
```

## Why

Exercise all 34 supported native encodings without running the full stress matrix.

## What

Creates and mounts a small seeded Mamba model per dtype, runs forward, and prints shape/checksum. Tensor activations stay Float32.

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
float64    shape=[1 3 8] checksum=-0.016678239
float32    shape=[1 3 8] checksum=-0.016678239
float16    shape=[1 3 8] checksum=-0.016669162
… (remaining native dtypes) …
Mamba dtypes=34; activations=float32; parameters=FormatNone native dtype
```
