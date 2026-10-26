#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
for chapter in 82-apple-metal 83-apple-mamba-dtypes 84-apple-training 85-apple-save-weights 86-apple-json-routing 87-apple-mixed-backends 88-apple-mamba-results; do
  echo "=== $chapter ==="
  (cd "$chapter" && go run .)
done
