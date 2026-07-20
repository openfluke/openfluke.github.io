#!/usr/bin/env bash
# Back-compat shim — use ./release.sh
exec "$(dirname "$0")/release.sh" "$@"
