# openfluke.github.io

**[Welvet feature book](welvet/)** — why/what for every engine package with runnable Go examples.

```bash
python3 _gen_welvet_book.py --run          # HTML + examples + captured outputs
python3 _gen_welvet_book.py --run --pdf    # + PDF in dist/ (upload to GitHub Release)
./release-book.sh                          # release helper
cd welvet/examples/01-welvet && source ../env.sh && go run .
```

PDF is **not** committed — version **v0.78** from `welvet/README.md`, upload `dist/welvet-feature-book-v0.78.pdf` to [GitHub Releases](https://github.com/openfluke/openfluke.github.io/releases).
