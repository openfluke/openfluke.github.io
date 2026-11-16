#!/usr/bin/env bash
# Run Part XVII Python C-ABI examples (114–117).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")" && pwd)"
# Prefer editable/source tree next to openfluke.github.io; fall back to installed wheel.
W2A_SRC="$(cd "$ROOT/../.." && pwd)/w2a/python/src"
if [[ -d "$W2A_SRC/welvet" ]]; then
  export PYTHONPATH="$W2A_SRC${PYTHONPATH:+:$PYTHONPATH}"
fi

for slug in 114-python-install 115-python-backends 116-python-train-entity 117-python-layers; do
  echo "=== $slug ==="
  (cd "$ROOT/$slug" && python3 main.py) || echo "FAIL $slug"
done
