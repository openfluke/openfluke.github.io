#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
for chapter in 89-nvidia-cuda 90-nvidia-mamba-dtypes 91-nvidia-training 92-nvidia-save-weights 93-nvidia-json-routing 94-nvidia-mixed-backends 95-nvidia-cuda-results; do
  echo "=== $chapter ==="
  (cd "$chapter" && go run .)
done
