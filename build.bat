@echo off
setlocal enabledelayedexpansion

REM Create dist directory if it doesn't exist
if not exist dist (
    mkdir dist
)

REM Build for Windows x64
set GOOS=windows
set GOARCH=amd64
echo Building for Windows x64...
go build -o dist\scaf-windows-amd64.exe .
echo Done.

REM Build for Windows arm64
set GOOS=windows
set GOARCH=arm64
echo Building for Windows arm64...
go build -o dist\scaf-windows-arm64.exe .
echo Done.

REM Build for Linux x64
set GOOS=linux
set GOARCH=amd64
echo Building for Linux x64...
go build -o dist\scaf-linux-amd64 .
echo Done.

REM Build for Linux arm64
set GOOS=linux
set GOARCH=arm64
echo Building for Linux arm64...
go build -o dist\scaf-linux-arm64 .
echo Done.

REM Build for macOS x64
set GOOS=darwin
set GOARCH=amd64
echo Building for macOS x64...
go build -o dist\scaf-darwin-amd64 .
echo Done.

REM Build for macOS arm64
set GOOS=darwin
set GOARCH=arm64
echo Building for macOS arm64...
go build -o dist\scaf-darwin-arm64 .
echo Done.

echo Build complete
pause