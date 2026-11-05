# Welvet book examples

One `main.go` per runnable chapter (108 chapters total).

Apple GPU: chapters 82–88; run `./run_apple.sh` on Apple silicon macOS with both native libraries installed. See `82-apple-metal/README.md` for setup and `apple_stress_test/README.md` for benchmarks.

NVIDIA GPU: chapters 89–95; run `./run_nvidia.sh` on Linux/Windows amd64 with `libwelvet_cuda` + NVIDIA runtime. Mount works on both OSes; JSON/mixed routing is Linux amd64. Soft-SKIP without a GPU.

AMD GPU: chapters 96–102; run `./run_amd.sh` on Linux/Windows amd64 with `libwelvet_rocm` + HIP runtime (sibling github.com/openfluke/welvet.rocm). Mount works on both OSes; JSON/mixed routing is Linux amd64. Soft-SKIP without a GPU; compute untested until ROCm-capable HSA.

Vulkan GPU: chapters 103–109; run `./run_vulkan.sh` on Linux amd64 with `libwelvet_vulkan` + Vulkan loader/runtime (sibling github.com/openfluke/welvet.vulkan). Soft-SKIP without a GPU. No silent host fallback labeled as Vulkan.

Apple MLX (v1.7): chapters 110–113; run `./run_mlx.sh` on Apple silicon macOS with `libwelvet_mlx` installed (sibling github.com/openfluke/welvet.mlx). Forward-only `nativemlx` / `BackendMLX` / `SyncToMLX`; validate with w2a `go_v_mlx`. No TrainMSE on MLX.

```bash
cd welvet/examples/01-welvet
source ../env.sh   # shared GOCACHE — required for webgpu/CGO examples
go run .
```

Regenerate HTML + capture all outputs:

```bash
cd openfluke.github.io && python3 _gen_welvet_book.py --run
```
