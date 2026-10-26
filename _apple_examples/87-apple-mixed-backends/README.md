# 87. Apple Metal — mixed Go/C++/Metal

**Part:** XII · Apple acceleration (welvet.metal)  
**Package:** `nativecpp + nativemetal + placement`    
**Status:** ok — verified macOS/arm64 on Apple M5

## When

Run different segments of one training graph on different backends.

## Where

From `example/welvet/87-apple-mixed-backends`:

```sh
go run .
```

## Why

Verify actual native C++ calls and Metal forward/backward/SGD, including gradient handoff.

## What

Embeds a 13-layer plan: Go 0–3, C++ 4–8, Metal 9–12. Runs five NormalBP updates, prints native-call counts and router statistics, and retains current weights in Go. Requires the normal build and C++ artifact.

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
Go layers 0-3 -> C++ layers 4-8 -> Metal layers 9-12
loss=0.0043765833 cpp calls=75 metal routing={Mounts:4 Forwards:40 Backwards:20 Updates:20 TrainingSteps:0}
Float32 activations travel forward; gradients travel back; routed Metal weights return automatically to Go.
```
