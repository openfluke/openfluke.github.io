#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
for chapter in 103-vulkan-mount 104-vulkan-mamba-dtypes 105-vulkan-training 106-vulkan-save-weights 107-vulkan-json-routing 108-vulkan-mixed-backends 109-vulkan-results; do
  echo "=== $chapter ==="
  (cd "$chapter" && go run .)
done
