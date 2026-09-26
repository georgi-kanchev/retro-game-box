#!/bin/bash
set -euo pipefail

echo "Building Windows app..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o app.exe .
echo "Build complete: ./app.exe"
echo "Press Enter to exit..."
read
