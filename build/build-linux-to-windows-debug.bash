#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

echo "Building Windows app (debug)..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -gcflags "all=-N -l" -o "$SCRIPT_DIR/../app_debug.exe" .
echo "Debug build complete: ./app_debug.exe"
echo "Press Enter to exit..."
read
