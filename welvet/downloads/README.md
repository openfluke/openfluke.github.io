# PDF downloads

The feature book PDF is **not stored in git** (it's ~1 MB and regenerated on release).

- **Build:** `python3 _gen_welvet_book.py --run --pdf` → writes `dist/welvet-feature-book-v0.XX.pdf`
- **Publish:** upload that file to [GitHub Releases](https://github.com/openfluke/openfluke.github.io/releases)
- **Version:** filename tracks the Welvet scorecard in `welvet/README.md`

The site title page links to the latest release download URL.
