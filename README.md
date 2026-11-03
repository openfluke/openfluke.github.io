# openfluke.github.io

**[Welvet feature book](welvet/)** — why/what for every engine package with runnable Go examples.

```bash
python3 _gen_welvet_book.py --run          # HTML + examples + captured outputs
python3 _gen_welvet_book.py --run --pdf    # + PDF in dist/ (upload to GitHub Release)
./release-book.sh                          # release helper
cd welvet/examples/01-welvet && source ../env.sh && go run .
```

PDF is **not** committed — version **v1.6.0** from sibling `../welvet/README.md`, upload `dist/welvet-feature-book-v1.6.0.pdf` to [GitHub Releases](https://github.com/openfluke/openfluke.github.io/releases).

**v1.6.0** adds Vulkan chapters (§103–109) for sibling [`welvet.vulkan`](https://github.com/openfluke/welvet.vulkan) (`nativevulkan`).
