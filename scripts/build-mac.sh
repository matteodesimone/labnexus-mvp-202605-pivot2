#!/usr/bin/env bash
# build-mac.sh — Costruisce labnexus.app per macOS Apple Silicon (FR-11).
# Bundle minimo: Contents/MacOS/labnexus + Contents/Info.plist.
# Binario NON firmato (Sprint 1 — Apple Developer cert in Sprint 2).
set -euo pipefail

APP_NAME="labnexus"
BUNDLE="${APP_NAME}.app"
APP_BIN_DIR="${BUNDLE}/Contents/MacOS"
APP_PLIST="${BUNDLE}/Contents/Info.plist"

rm -rf "${BUNDLE}"
mkdir -p "${APP_BIN_DIR}"

echo "[build-mac] cross-compile darwin/arm64..."
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 \
  go build -trimpath -ldflags "-s -w" \
  -o "${APP_BIN_DIR}/${APP_NAME}" ./cmd/labnexus

cat > "${APP_PLIST}" <<'PLIST'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTD/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleExecutable</key><string>labnexus</string>
  <key>CFBundleIdentifier</key><string>com.labnexus.app</string>
  <key>CFBundleName</key><string>labnexus</string>
  <key>CFBundleVersion</key><string>1.0</string>
  <key>CFBundlePackageType</key><string>APPL</string>
  <key>LSMinimumSystemVersion</key><string>11.0</string>
  <key>LSHandlerRank</key><string>Alternate</string>
  <key>CFBundleDocumentTypes</key>
  <array>
    <dict>
      <key>CFBundleTypeName</key><string>Folder</string>
      <key>LSItemContentTypes</key>
      <array><string>public.folder</string></array>
      <key>CFBundleTypeRole</key><string>Editor</string>
    </dict>
  </array>
</dict>
</plist>
PLIST

echo "[build-mac] bundle pronto: ${BUNDLE}"
file "${APP_BIN_DIR}/${APP_NAME}" || true
