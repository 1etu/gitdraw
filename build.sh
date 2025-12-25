#!/bin/bash

# GitDraw Release Build Script
# Builds GUI for macOS (Intel + Apple Silicon) and Windows
# Builds CLI for all platforms

set -e

VERSION="${1:-1.0.0}"
DIST_DIR="./dist"
BUILD_DIR="./build/bin"

echo "GitDraw Build v${VERSION}"
echo "=================================="

rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

WAILS=$(which wails 2>/dev/null || echo "$HOME/go/bin/wails")

echo ""
echo "📦 Building GUI apps..."

echo "  → macOS Apple Silicon..."
$WAILS build -tags gui -platform darwin/arm64 -clean -skipbindings 2>/dev/null || $WAILS build -tags gui -platform darwin/arm64 -clean
mv "$BUILD_DIR/gitdraw.app" "$DIST_DIR/GitDraw-macOS-AppleSilicon.app"

echo "  → macOS Intel..."
$WAILS build -tags gui -platform darwin/amd64 -clean -skipbindings 2>/dev/null || $WAILS build -tags gui -platform darwin/amd64 -clean
mv "$BUILD_DIR/gitdraw.app" "$DIST_DIR/GitDraw-macOS-Intel.app"

echo "  → Windows x64..."
$WAILS build -tags gui -platform windows/amd64 -clean -skipbindings 2>/dev/null || $WAILS build -tags gui -platform windows/amd64 -clean
mv "$BUILD_DIR/gitdraw.exe" "$DIST_DIR/GitDraw-Windows-x64.exe"

echo ""
echo "Building CLI binaries..."

echo "  → CLI macOS Apple Silicon..."
GOOS=darwin GOARCH=arm64 go build -o "$DIST_DIR/gitdraw-cli-macos-arm64" .

echo "  → CLI macOS Intel..."
GOOS=darwin GOARCH=amd64 go build -o "$DIST_DIR/gitdraw-cli-macos-amd64" .

echo "  → CLI Windows x64..."
GOOS=windows GOARCH=amd64 go build -o "$DIST_DIR/gitdraw-cli-windows-x64.exe" .

echo "  → CLI Linux x64..."
GOOS=linux GOARCH=amd64 go build -o "$DIST_DIR/gitdraw-cli-linux-amd64" .

echo ""
echo "Creating archives..."

cd "$DIST_DIR"

zip -r -q "GitDraw-macOS-AppleSilicon-v${VERSION}.zip" "GitDraw-macOS-AppleSilicon.app"
zip -r -q "GitDraw-macOS-Intel-v${VERSION}.zip" "GitDraw-macOS-Intel.app"

zip -q "GitDraw-Windows-x64-v${VERSION}.zip" "GitDraw-Windows-x64.exe"

zip -q "gitdraw-cli-macos-arm64-v${VERSION}.zip" "gitdraw-cli-macos-arm64"
zip -q "gitdraw-cli-macos-amd64-v${VERSION}.zip" "gitdraw-cli-macos-amd64"
zip -q "gitdraw-cli-windows-x64-v${VERSION}.zip" "gitdraw-cli-windows-x64.exe"
zip -q "gitdraw-cli-linux-amd64-v${VERSION}.zip" "gitdraw-cli-linux-amd64"

cd ..

echo ""
echo "Build complete! Files in ./dist:"
echo ""
ls -lh "$DIST_DIR"/*.zip 2>/dev/null || ls -lh "$DIST_DIR"
echo ""
echo "OK"
