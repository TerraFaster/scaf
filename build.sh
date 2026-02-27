#!/bin/bash

set -e

mkdir -p dist

echo "Building for Windows x64..."
GOOS=windows GOARCH=amd64 go build -o dist/scaf-windows-amd64.exe .
echo "Done."

echo "Building for Windows arm64..."
GOOS=windows GOARCH=arm64 go build -o dist/scaf-windows-arm64.exe .
echo "Done."

echo "Building for Linux x64..."
GOOS=linux GOARCH=amd64 go build -o dist/scaf-linux-amd64 .
echo "Done."

echo "Building for Linux arm64..."
GOOS=linux GOARCH=arm64 go build -o dist/scaf-linux-arm64 .
echo "Done."

echo "Building for macOS x64..."
GOOS=darwin GOARCH=amd64 go build -o dist/scaf-darwin-amd64 .
echo "Done."

echo "Building for macOS arm64..."
GOOS=darwin GOARCH=arm64 go build -o dist/scaf-darwin-arm64 .
echo "Done."

echo "Build complete"