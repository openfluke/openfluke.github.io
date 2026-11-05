#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
for chapter in 110-apple-mlx 111-apple-mlx-caps 112-apple-mlx-forward 113-apple-mlx-results; do
  echo "=== $chapter ==="
  (cd "$chapter" && go run .)
done
