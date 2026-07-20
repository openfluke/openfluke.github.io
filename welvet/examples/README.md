# Welvet book examples

One `main.go` per chapter (64 total).

```bash
cd welvet/examples/01-welvet
source ../env.sh   # shared GOCACHE — required for webgpu/CGO examples
go run .
```

Regenerate HTML + capture all outputs:

```bash
cd openfluke.github.io && python3 _gen_welvet_book.py --run
```
