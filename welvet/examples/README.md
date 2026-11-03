# Welvet book examples

One `main.go` per runnable chapter (97 chapters total).

Apple GPU: chapters 82–88; run `./run_apple.sh` on Apple silicon macOS with both native libraries installed. See `82-apple-metal/README.md` for setup and `apple_stress_test/README.md` for benchmarks.

NVIDIA GPU: chapters 89–95; run `./run_nvidia.sh` on Linux/Windows amd64 with `libwelvet_cuda` + NVIDIA runtime. Mount works on both OSes; JSON/mixed routing is Linux amd64. Soft-SKIP without a GPU.

AMD GPU: chapters 96–102; run `./run_amd.sh` on Linux/Windows amd64 with `libwelvet_rocm` + HIP runtime (sibling github.com/openfluke/welvet.rocm). Mount works on both OSes; JSON/mixed routing is Linux amd64. Soft-SKIP without a GPU; compute untested until ROCm-capable HSA.

```bash
cd welvet/examples/01-welvet
source ../env.sh   # shared GOCACHE — required for webgpu/CGO examples
go run .
```

Regenerate HTML + capture all outputs:

```bash
cd openfluke.github.io && python3 _gen_welvet_book.py --run
```
