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

# Escape un argomento per inclusione in `do script "..."` di AppleScript.
# Doppio livello:
#   1) shell single-quoting: ' → '\''  (per il `do script` che esegue shell)
#   2) AppleScript string literal: \ → \\, " → \"  (per la stringa AppleScript stessa)
# Questo gestisce robustamente path con apostrofi, virgolette, backslash.
escape_for_osascript() {
  printf %s "$1" \
    | sed "s/'/'\\\\''/g" \
    | sed 's/\\/\\\\/g; s/"/\\"/g'
}

# Costruisce il comando shell da passare a `do script`, sempre passando per escape_for_osascript.
build_command() {
  local bin_safe input_safe
  bin_safe=$(escape_for_osascript "$1")
  if [ "$#" -ge 2 ] && [ -n "$2" ]; then
    input_safe=$(escape_for_osascript "$2")
    printf "'%s' '%s'" "$bin_safe" "$input_safe"
  else
    printf "'%s'" "$bin_safe"
  fi
}

# Determina l'argomento posizionale (drag&drop di una cartella esistente)
if [ "$#" -ge 1 ] && [ -d "$1" ]; then
  CMD=$(build_command "${BIN}" "$1")
else
  CMD=$(build_command "${BIN}")
fi

/usr/bin/osascript <<EOF
tell application "Terminal"
  activate
  do script "${CMD}"
end tell
EOF
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
