#!/usr/bin/env bash
# Build Welvet feature book + PDF for a release. PDF goes to dist/ (gitignored).
set -euo pipefail
cd "$(dirname "$0")"

python3 _gen_welvet_book.py --run --pdf

pdf="$(ls -1 dist/welvet-feature-book-*.pdf | tail -1)"
ver="$(basename "$pdf" .pdf | sed 's/welvet-feature-book-//')"
echo ""
echo "Built: $pdf"
echo ""
echo "Publish to GitHub Releases:"
echo "  gh release create \"$ver\" \"$pdf\" --title \"Welvet feature book $ver\" --notes \"Generated from welvet scorecard $ver.\""
echo "  # or upload to an existing release:"
echo "  gh release upload \"$ver\" \"$pdf\" --clobber"
