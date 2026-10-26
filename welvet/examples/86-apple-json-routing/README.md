# 86. Apple Metal — JSON backend placement

**Part:** XII · Apple acceleration (welvet.metal)  
**Package:** `nativemetal + runtime/placement`    
**Status:** ok — verified macOS/arm64 on Apple M5

## When

Assign Go and Metal operations using the shared execution JSON.

## Where

From `example/welvet/86-apple-json-routing`:

```sh
go run .
go run . -config placement.json
```

## Why

Use ordinary Forward/TrainStackMSE calls after installing the routing owner.

## What

Embeds placement.json: layer 0 runs Go CPU, layers 1–2 run Metal. EnableRouting retains native mounts and returns updated parameters to Go automatically. Close the router after the workload.

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
layer=0 backend=cpu-tiled
layer=1 backend=metal
layer=2 backend=metal
loss=0.055711714 metal routing={Mounts:2 Forwards:4 Backwards:2 Updates:2 TrainingSteps:0}
Routed updates are already in Go stores; Close releases cached mounts.
```
