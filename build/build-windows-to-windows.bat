@echo off
setlocal enabledelayedexpansion

set "SCRIPT_DIR=%~dp0"
set "OUTPUT_EXE=%SCRIPT_DIR%..\app.exe"

cd /d "%SCRIPT_DIR%.."

echo Building Windows app...
go build -ldflags="-s -w -H=windowsgui" -o "%OUTPUT_EXE%" .

if %ERRORLEVEL% EQU 0 (
    echo ---------------------------------------
    echo SUCCESS! Build complete: %OUTPUT_EXE%
) else (
    echo ---------------------------------------
    echo [ERROR] Build failed. Check the logs above.
)

pause
