#!/bin/bash
set -euo pipefail

echo "Building Windows app (debug)..."
GOOS=windows GOARCH=amd64 go build -gcflags "all=-N -l" -o app_debug.exe .
echo "Debug build complete: ./app_debug.exe"
echo "Press Enter to exit..."
read
