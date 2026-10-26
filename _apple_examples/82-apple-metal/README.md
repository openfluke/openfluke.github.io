# 82. Apple Metal — mount and forward

**Part:** XII · Apple acceleration (welvet.metal)  
**Package:** `nativemetal`    
**Status:** ok — verified macOS/arm64 on Apple M5

## When

Mount a Go layer and execute its forward pass on Apple GPU.

## Where

From `example/welvet/82-apple-metal`:

```sh
go run .
go run . -library /path/to/libwelvet_metal.dylib
```

## Why

See the backend version and actual Metal device identity; release the session explicitly.

## What

Creates an identity Dense layer, mounts it, forwards a Float32 tensor, and closes the session.

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
backend=metal version=0.68
device=metal:Apple M5:reg=4294968322:ws=19069665280
input=[0.5 -0.25] output=[0.5 -0.25] shape=[1 2]
```
