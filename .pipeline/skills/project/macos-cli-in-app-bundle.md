# Skill — CLI in macOS .app bundle (wrapper bash + osascript)

Pattern per impacchettare un binario CLI dentro un macOS `.app` bundle in modo che il doppio click dal Finder apra Terminal con la CLI in esecuzione. Emerso dal bug `app-bundle-doppio-click-no-output` (Sprint 1 LabNexus, fix 2026-05-19).

## Quando applicarlo

- CLI Go/Rust/Python che vuoi distribuire come `.app` macOS (no firma Apple Developer, no installer)
- L'utente target NON è uno sviluppatore (vuole doppio click, non Terminal manuale)
- Il binario richiede una TUI interattiva o output line-based che resta visibile

## Problema che risolve

Se `CFBundleExecutable` punta direttamente al binario:
1. Launch Services lo lancia in background senza terminale
2. `stdin`/`stdout`/`stderr` non sono TTY
3. Codice che rileva non-TTY (es. `term.IsTerminal()`) esce con errore
4. L'utente vede **niente** perché lo stderr non ha finestra

## Pattern

### Struttura bundle

```
<app>.app/Contents/
  MacOS/
    <app>      ← wrapper bash (CFBundleExecutable, ~1.5 KB)
    <app>-bin  ← binario reale (Mach-O arm64)
  Info.plist   ← CFBundleExecutable=<app>
```

### Wrapper bash (template)

```bash
#!/usr/bin/env bash
set -eu
DIR="$(cd "$(dirname "$0")" && pwd)"
BIN="${DIR}/<app>-bin"

# Doppio escape: shell single-quoting + AppleScript string literal
escape_for_osascript() {
  printf %s "$1" \
    | sed "s/'/'\\\\''/g" \
    | sed 's/\\/\\\\/g; s/"/\\"/g'
}

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
```

### Info.plist (minimum)

```xml
<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
<dict>
  <key>CFBundleExecutable</key><string><app></string>
  <key>CFBundleIdentifier</key><string>com.example.<app></string>
  <key>CFBundleName</key><string><app></string>
  <key>CFBundleVersion</key><string>1.0</string>
  <key>CFBundlePackageType</key><string>APPL</string>
  <key>LSMinimumSystemVersion</key><string>11.0</string>
  <!-- Per drag&drop di cartelle sull'icona -->
  <key>LSHandlerRank</key><string>Alternate</string>
  <key>CFBundleDocumentTypes</key>
  <array>
    <dict>
      <key>CFBundleTypeName</key><string>Folder</string>
      <key>LSItemContentTypes</key><array><string>public.folder</string></array>
      <key>CFBundleTypeRole</key><string>Editor</string>
    </dict>
  </array>
</dict>
</plist>
```

### Build script (Go esempio)

```bash
#!/usr/bin/env bash
set -euo pipefail
APP_NAME="<app>"
BUNDLE="${APP_NAME}.app"
APP_BIN_DIR="${BUNDLE}/Contents/MacOS"

rm -rf "${BUNDLE}"
mkdir -p "${APP_BIN_DIR}"

GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 \
  go build -trimpath -ldflags "-s -w" \
  -o "${APP_BIN_DIR}/${APP_NAME}-bin" ./cmd/<app>

cat > "${APP_BIN_DIR}/${APP_NAME}" <<'WRAPPER'
[wrapper bash template sopra, NB heredoc 'WRAPPER' impedisce expansion]
WRAPPER
chmod +x "${APP_BIN_DIR}/${APP_NAME}"

cat > "${BUNDLE}/Contents/Info.plist" <<'PLIST'
[Info.plist sopra]
PLIST
```

## Test integration

Pattern `TestMain` + assertion sul filesystem prodotto:

```go
//go:build darwin

package bundle_test

import (
    "os"
    "os/exec"
    "testing"
)

func TestMain(m *testing.M) {
    // Build idempotente: lo script fa `rm -rf <app>.app` all'inizio.
    cmd := exec.Command("bash", "scripts/build-mac.sh")
    cmd.Dir = "../.."
    if out, err := cmd.CombinedOutput(); err != nil {
        os.Stderr.Write(out)
        os.Exit(1)
    }
    os.Exit(m.Run())
}

func TestWrapperIsTextScript(t *testing.T) { /* shebang, size < 4KB */ }
func TestRealBinaryIsArm64(t *testing.T)    { /* Mach-O arm64 */ }
func TestWrapperEscapesAS(t *testing.T)      { /* grep sed pattern AppleScript escape */ }
func TestWrapperLaunchesTerminal(t *testing.T) { /* contains "osascript" + "Terminal" + "<app>-bin" */ }
func TestInfoPlistCorrect(t *testing.T)      { /* CFBundleExecutable=<app> */ }
```

## Anti-pattern da evitare

- ❌ `CFBundleExecutable` punta direttamente al binario nudo → silent exit 2 dal Finder
- ❌ Solo escape `'` per shell, senza AppleScript escape → path con `"` rompe `do script`
- ❌ Non testare scenari `@manual` "doppio click reale" prima della consegna (vedi framework promotion `manual-scenarios-deploy-checkpoint`)
- ❌ Lanciare il wrapper da CLI in shell esistente (apre **seconda** finestra Terminal) — documentare nel README

## Limitazioni note

- **Non firmato Apple Developer**: l'utente alla prima apertura deve fare control-clic → "Apri" → conferma (Gatekeeper standard). Documentare nel README della consegna.
- **macOS 11+ richiesto** (LSMinimumSystemVersion=11.0 nel template). Adattare se serve compatibility più ampia.
- **Solo darwin/arm64** nel template. Per `darwin/amd64` (Mac Intel) aggiungere `GOOS=darwin GOARCH=amd64` come secondo build target e usare `lipo` per fat binary (oppure due bundle separati).

## Vedi anche

- `.pipeline/solutions/2026-05-19-bugfix-app-bundle-doppio-click.md` — emersione del pattern
- `.pipeline/bugs/app-bundle-doppio-click-no-output.md` — bug originale
- `scripts/build-mac.sh` (LabNexus) — implementazione concreta
- `internal/bundle/bundle_test.go` (LabNexus) — test integration template
- Framework proposal `.pipeline/proposed-updates/2026-05-19-manual-scenarios-deploy-checkpoint.md`
