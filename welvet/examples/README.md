# Welvet book examples

One `main.go` per runnable chapter (90 chapters total).

Apple GPU: chapters 82–88; run `./run_apple.sh` on Apple silicon macOS with both native libraries installed. See `82-apple-metal/README.md` for setup and `apple_stress_test/README.md` for benchmarks.

NVIDIA GPU: chapters 89–95; run `./run_nvidia.sh` on Linux/Windows amd64 with `libwelvet_cuda` + NVIDIA runtime. Mount works on both OSes; JSON/mixed routing is Linux amd64. Soft-SKIP without a GPU.

```bash
cd welvet/examples/01-welvet
source ../env.sh   # shared GOCACHE — required for webgpu/CGO examples
go run .
```

Regenerate HTML + capture all outputs:

```bash
cd openfluke.github.io && python3 _gen_welvet_book.py --run
```
