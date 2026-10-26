# Welvet book examples

One `main.go` per runnable chapter (83 chapters total).

Apple GPU: chapters 82–88; run `./run_apple.sh` on Apple silicon macOS with both native libraries installed. See `82-apple-metal/README.md` for setup and `apple_stress_test/README.md` for benchmarks.

```bash
cd welvet/examples/01-welvet
source ../env.sh   # shared GOCACHE — required for webgpu/CGO examples
go run .
```

Regenerate HTML + capture all outputs:

```bash
cd openfluke.github.io && python3 _gen_welvet_book.py --run
```
