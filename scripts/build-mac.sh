#!/usr/bin/env bash
# build-mac.sh — Costruisce labnexus.app per macOS Apple Silicon (FR-11).
#
# Layout del bundle (pattern idiomatic macOS per CLI-in-.app):
#   labnexus.app/Contents/
#     MacOS/
#       labnexus      ← wrapper bash (CFBundleExecutable)
#       labnexus-bin  ← binario Go reale (Mach-O arm64)
#     Info.plist
#
# Il wrapper apre Terminal via osascript e lancia il binario reale al suo
# interno, così il doppio click dal Finder mostra effettivamente la TUI.
# Fix del bug: .pipeline/bugs/app-bundle-doppio-click-no-output.md
set -euo pipefail

APP_NAME="labnexus"
BIN_NAME="labnexus-bin"
BUNDLE="${APP_NAME}.app"
APP_BIN_DIR="${BUNDLE}/Contents/MacOS"
APP_PLIST="${BUNDLE}/Contents/Info.plist"

rm -rf "${BUNDLE}"
mkdir -p "${APP_BIN_DIR}"

echo "[build-mac] cross-compile darwin/arm64 → ${BIN_NAME}..."
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 \
  go build -trimpath -ldflags "-s -w" \
  -o "${APP_BIN_DIR}/${BIN_NAME}" ./cmd/labnexus

# Wrapper bash: CFBundleExecutable. Lanciato da Finder (doppio click o drag&drop)
# o da CLI. Apre Terminal via osascript e lancia il binario reale al suo interno.
echo "[build-mac] scrittura wrapper bash → ${APP_NAME}..."
cat > "${APP_BIN_DIR}/${APP_NAME}" <<'WRAPPER'
#!/usr/bin/env bash
# CFBundleExecutable wrapper di labnexus.app.
# Lancia il binario reale (labnexus-bin) dentro una finestra Terminal,
# così la TUI può girare in modalità interattiva (FR-10).
# Se invocato con un argomento posizionale = cartella esistente (drag&drop),
# la passa al binario come input pre-selezionato.

set -eu
DIR="$(cd "$(dirname "$0")" && pwd)"
BIN="${DIR}/labnexus-bin"

# Funzione: lancia il binario in Terminal via osascript.
# $1 = comando shell singolarmente quotato da eseguire in Terminal.
launch_in_terminal() {
  /usr/bin/osascript <<EOF
tell application "Terminal"
  activate
  do script "$1"
end tell
EOF
}

# Se è stato passato un argomento posizionale che è una cartella esistente
# (caso drag&drop dal Finder), passala come input al binario.
if [ "$#" -ge 1 ] && [ -d "$1" ]; then
  # Escape singole virgolette nel path per osascript (rare ma possibili)
  INPUT_PATH=$(printf %s "$1" | sed "s/'/'\\\\''/g")
  launch_in_terminal "'${BIN}' '${INPUT_PATH}'"
else
  # Doppio click semplice: TUI parte dalla selezione profilo
  launch_in_terminal "'${BIN}'"
fi
WRAPPER

chmod +x "${APP_BIN_DIR}/${APP_NAME}"

# Info.plist
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
file "${APP_BIN_DIR}/${BIN_NAME}"
echo "[build-mac] wrapper: $(wc -l < "${APP_BIN_DIR}/${APP_NAME}") righe"
