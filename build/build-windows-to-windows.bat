@echo off
setlocal enabledelayedexpansion

set "OUTPUT_EXE=app.exe"

echo Building Windows app...
go build -ldflags="-s -w -H=windowsgui" -o "%OUTPUT_EXE%" .

if %ERRORLEVEL% EQU 0 (
    echo ---------------------------------------
    echo SUCCESS! Build complete.
) else (
    echo ---------------------------------------
    echo [ERROR] Build failed. Check the logs above.
)

pause
