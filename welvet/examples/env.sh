#!/usr/bin/env bash
# Shared Go build cache for book examples (avoids /tmp quota exhaustion with webgpu CGO).
_root="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export GOCACHE="${GOCACHE:-$_root/.cache/gocache}"
export GOTMPDIR="${GOTMPDIR:-$_root/.cache/gotmp}"
export TMPDIR="${TMPDIR:-$_root/.cache/tmp}"
mkdir -p "$GOCACHE" "$GOTMPDIR" "$TMPDIR"
