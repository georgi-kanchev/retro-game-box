#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

echo "Building Linux app..."
CGO_ENABLED=1 go build -ldflags="-s -w" -o "$SCRIPT_DIR/../app" .
echo "Build complete: ./app"
echo "Press Enter to exit..."
read