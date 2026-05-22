---
name: external-binaries-bundled-in-dist
date: 2026-05-22
source: Sprint 1.5.D — PDF/DOCX rendering pivot
tags: [build, deploy, cli-tool, go, bundle, pandoc, typst]
---

# External binaries bundled in dist/

## When to use

Go cli-tool che richiede tool nativi non scrivibili in Go (pandoc, typst, ffmpeg, imagemagick, latex engines, …) per features critiche (rendering documenti, conversioni media, ecc.). Vincolo architetturale: **niente runtime esterno richiesto all'utente finale**.

## Pattern

### 1. Script download idempotente

`scripts/download-<tools>.sh` con:
- URL pinati per versione specifica (no `latest` flotting)
- Cache locale `.cache/<tools>/` per i `.zip`/`.tar` scaricati
- Estrazione e copia binary in `dist/<pkg>/bin/`
- Skip se binary già presente (`-x` test) — idempotenza
- `chmod +x` post-copy
- `xattr -d com.apple.quarantine` per evitare Gatekeeper su macOS

```bash
if [[ -x "${DIST_BIN}/pandoc" ]]; then
  echo "[pdf-tools] pandoc già presente — skip"
else
  curl -L --fail "${PANDOC_URL}" -o "${CACHE}/pandoc.zip"
  unzip -q "${CACHE}/pandoc.zip" -d "${EXTRACT}"
  cp "$(find ${EXTRACT} -name pandoc -type f -perm +111 | head -1)" "${DIST_BIN}/pandoc"
  chmod +x "${DIST_BIN}/pandoc"
  xattr -d com.apple.quarantine "${DIST_BIN}/pandoc" 2>/dev/null || true
fi
```

### 2. Embed assets via go:embed

Template/config/lua-filter che accompagnano il tool sono embedded nel binary Go:
```go
//go:embed template.typ
var templateTypst string

//go:embed callout-filter.lua
var calloutFilterLua string

//go:embed reference.docx
var referenceDocx []byte
```
A runtime li scrivi in `os.MkdirTemp` e li passi al tool via `--filter` / `--reference-doc`.

### 3. Risoluzione runtime con fallback dev

```go
func resolveTool(name string) (string, error) {
    // Env override per test
    if p := os.Getenv("LABNEXUS_"+strings.ToUpper(name)+"_PATH"); p != "" {
        return p, nil
    }
    // Path canonico: <dir-binary>/bin/<tool>
    exe, _ := os.Executable()
    exe, _ = filepath.EvalSymlinks(exe)
    p := filepath.Join(filepath.Dir(exe), "bin", name)
    if isExecutable(p) {
        return p, nil
    }
    // Fallback dev: risali dalla CWD cercando dist/<pkg>/bin/
    if dev, ok := findRepoDistBin(); ok {
        if p := filepath.Join(dev, name); isExecutable(p) {
            return p, nil
        }
    }
    return "", ErrToolMissing
}
```

### 4. Makefile target + ship order

```makefile
.PHONY: pdf-tools
pdf-tools: ## Scarica tool external nativi in dist/.../bin/
	bash scripts/download-pdf-tools.sh

ship: clean pdf-tools test vet package
```

**Critico**: `pdf-tools` PRIMA di `test` nella catena `ship`. Se `clean` rimuove dist e i test richiedono i tool, devi popolare bin/ prima.

### 5. Test resilient

I test che usano i tool devono **skippare graciously** quando mancano:
```go
err := tool.Render(...)
if err == ErrToolMissing || strings.Contains(err.Error(), "pandoc") {
    t.Skipf("pandoc non disponibile (esegui make pdf-tools)")
}
```

Così `go test ./...` funziona anche su working tree senza tool scaricati.

### 6. build-zip include bin/

`scripts/build-zip.sh` deve:
- Verificare `bin/` esiste (errore esplicito se mancano i tool)
- Popolare `dist/<pkg>/` con tutto il contenuto **preservando** `bin/` esistente (`find ... ! -name bin -exec rm -rf`)
- Zippare la cartella popolata (non un STAGE_DIR temporaneo)

Output: sia il `.zip` per consegna sia la cartella estratta pronta per smoke locale (doppio click senza extract).

## Costo dist

Pandoc 180MB + Typst 39MB = ~210MB aggiunti al deliverable. Single-binary value preservato logicamente (tutto in un solo archivio).

## Quando NON usare

- Se la feature è **opzionale** (es. PDF in Sprint 2 ma non Sprint 1) e l'aggiunta di 200MB+ ritarda il primo shipping
- Se il tool ha versioni nettamente diverse per OS (es. wkhtmltopdf macOS x64 vs arm64 nativo non esiste, va via Rosetta → no, scartato in Sprint 1.5.D)
- Se esiste alternativa Go pura "good enough" per il use case (es. goldmark+fpdf per PDF semplici — ma attenzione: bug visuali cumulativi vanno valutati prima del pivot, vedi solution Sprint 1.5.D).

## Anti-patterns visti

- **Download a runtime al primo uso**: l'utente la prima volta non ha rete → fallisce. Pre-scaricare in dist è la via.
- **PATH system pollution**: `brew install pandoc` o `apt install`: fragile, version skew, niente garanzia presenza. Bundling vince.
- **go:embed del binary tool**: tecnicamente possibile, embed in binary Go = 200MB di Go binary che estrae a tempdir al primo run. Hai due problemi (binary enorme + extract logic) invece di uno (download script). Sconsigliato.
