#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
for chapter in 96-amd-rocm 97-amd-mamba-dtypes 98-amd-training 99-amd-save-weights 100-amd-json-routing 101-amd-mixed-backends 102-amd-rocm-results; do
  echo "=== $chapter ==="
  (cd "$chapter" && go run .)
done
