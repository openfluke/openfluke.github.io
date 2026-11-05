# openfluke.github.io

**[Welvet feature book](welvet/)** — why/what for every engine package with runnable Go examples.

```bash
python3 _gen_welvet_book.py --run          # HTML + examples + captured outputs
python3 _gen_welvet_book.py --run --pdf    # + PDF in dist/ (upload to GitHub Release)
./release-book.sh                          # release helper
cd welvet/examples/01-welvet && source ../env.sh && go run .
```

PDF is **not** committed — version **v1.7.0** from sibling `../welvet/README.md`, upload `dist/welvet-feature-book-v1.7.0.pdf` to [GitHub Releases](https://github.com/openfluke/openfluke.github.io/releases).

**v1.7.0** adds forward-only Apple MLX via sibling [`welvet.mlx`](https://github.com/openfluke/welvet.mlx) (`nativemlx` / `BackendMLX`). **v1.6.x** shipped Vulkan chapters (§103–109) for [`welvet.vulkan`](https://github.com/openfluke/welvet.vulkan) (`nativevulkan`).
