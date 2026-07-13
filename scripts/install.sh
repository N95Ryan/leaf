#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"

if ! command -v go >/dev/null 2>&1; then
  echo "Go is not installed. Install Go 1.25+ from https://go.dev/dl/" >&2
  exit 1
fi

echo "Building and installing leaf..."
go install ./cmd/leaf
echo "Done. Run: leaf"
