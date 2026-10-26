# 85. Apple Metal — weight return and ENTITY save

**Part:** XII · Apple acceleration (welvet.metal)  
**Package:** `nativemetal + entity`    
**Status:** ok — verified macOS/arm64 on Apple M5

## When

Bring trained parameters back to Go for saving and reloading.

## Where

From `example/welvet/85-apple-save-weights`:

```sh
go run .
go run . -out trained.entity
```

## Why

Explicit Mount owns a copy: Go parameters stay unchanged until SyncToGo.

## What

Asserts original Go weights remain unchanged during native training, explicitly pulls them, writes a Go ENTITY file, reloads it, and verifies the native parameter snapshot exactly. With no -out, saves into a new temporary directory and prints the path.

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
Go parameters unchanged before SyncToGo
weights returned to Go; saved=<temporary directory>/trained.entity; reload=exact
```
