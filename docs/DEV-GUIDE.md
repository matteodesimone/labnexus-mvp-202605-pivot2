# LabNexus — Guida sviluppatore

Per il CTO + futuro contributor. Build, test, layout, gotchas.

## Prerequisites

- Go 1.22+
- mise (gestione tool versions via `.mise.toml`)
- `make` (preinstallato su macOS e Linux)
- Per i test BDD: il binario viene compilato da `features/main_test.go` al `TestMain` — niente da installare a parte
- Per `labnexus run` reale: Ollama (`ollama serve`) + modello `qwen3.6` (`ollama pull qwen3.6`); o `EUROUTER_API_KEY` per provider eurouter

## Makefile — entry point unico

**Tutto si fa via `make`.** Lancia `make` (senza argomenti) per l'help colorato con tutti i target. Sotto il cofano i target chiamano `scripts/build-*.sh` e `go ...`; gli script restano comunque utilizzabili a mano se preferisci.

### Workflow tipici

| Quando | Comando | Cosa fa |
|---|---|---|
| Pre-commit / CI | `make ci` | `tidy + vet + test` — ~1.5s su cache calda |
| Iterazione locale | `make build && make smoke` | Compila + smoke test (list/describe/validate/check su Test 1 reale, no LLM) |
| Debug test rotti | `make test-verbose` | Output `-v` per individuare lo scenario fallito |
| Consegna Denis | `make ship` | ⭐ One-shot: `clean + test + vet + build-mac + zip` → `dist/labnexus-sprint1-darwin-arm64.zip` |
| Pulizia | `make clean` | Rimuove `bin/`, `dist/`, `labnexus.app/`, `labnexus-smoke`, `labnexus-final` |

### Lista completa dei target

```
make                Help colorato (default)

Build:
  build             Binario per la macchina corrente → bin/labnexus
  build-mac         Cross-compile darwin/arm64 → labnexus.app
  build-linux       Cross-compile linux/amd64 → dist/labnexus-linux-amd64
  build-all         build-mac + build-linux (entrambi i target Sprint 1)

Test:
  test              go test ./...
  test-unit         Solo internal/... (unit + integration)
  test-bdd          Solo features/ (godog BDD, tag ~@manual)
  test-verbose      Output -v per debug fallimenti

Sanity / Lint:
  vet               go vet ./...
  fmt               go fmt ./...
  lint              golangci-lint run (richiede install separato, vedi sotto)
  tidy              go mod tidy
  ci                tidy + vet + test (shortcut CI)

Smoke:
  smoke             Build + list + describe + validate + check su Test 1 reale (dry-run, no LLM)

Packaging:
  package           build-mac + zip → dist/labnexus-sprint1-darwin-arm64.zip
  ship              clean + test + vet + package (pipeline completa di consegna)

Cleanup:
  clean             Rimuove bin/, dist/, labnexus.app/, smoke binaries
  distclean         clean + go clean -testcache
```

### Versioning automatico

`make build` injetta `-X main.version=<git-describe>` nel binario. Quando taggherai (es. `git tag v0.1.0-fetta1`), `labnexus --version` riporterà il tag. Senza tag mostra l'hash del commit corrente + suffisso `-dirty` se ci sono modifiche non committate.

### Singolo scenario BDD (oltre Makefile)

Se vuoi runnare un singolo scenario filtrando per tag godog:

```
go test ./features/ -godog.tags=@motore-cli
go test ./features/ -godog.tags="not @manual and @ollama"
```

### Linting — install di `golangci-lint`

`golangci-lint` NON è gestito da mise (il backend aqua va in 401 sul rate-limit GitHub API). Installalo separatamente:

- **macOS**: `brew install golangci-lint`
- **Linux/WSL2**: `curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.5.0`
- **Cross-platform**: `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`

Verifica: `golangci-lint --version` (atteso ≥ 2.0). Poi `make lint` funziona.

### Quando NON usare Makefile

I target make sono shortcut idiomatici; sotto il cofano sono `go build`/`go test`/`bash scripts/*.sh`. Se hai bisogno di flag custom (es. `go test -count=10 -race`), invoca direttamente il comando — il Makefile non è una camicia di forza.

## Layout codice

```
cmd/labnexus/main.go          # presentation: cobra root + classifyError NFR-5
internal/                     # tutta la logica, no import I/O framework
  profile/                    # schema YAML + validate (FR-3) + path-traversal guard
  input/                      # parser multi-formato (md/txt/csv/pdf/docx/xlsx)
                              #   - magic-bytes check su PDF
                              #   - soglia strict-majority `> 50%` (EC-1)
                              #   - walk non ricorsivo
  prompt/                     # Compose: system = join(kb, "---") + user = trigger + "## File di input" + blocks
  tokens/                     # Estimate (char/4) + Check (warn ≥0.7, block ≥1.0)
  provider/                   # LLMProvider interface + 2 impl + Select(flag>env>profile)
    provider.go               #   Select + types
    ollama.go                 #   NDJSON streaming, tolerance EC-13
    eurouter.go               #   SSE OpenAI-compatible, EC-14 timeout
  output/                     # Frontmatter + Write con collision suffix _2/_3 (EC-15)
  runlog/                     # TTY-aware logger + progress bar
  runner/                     # orchestratore Pipeline stages (refactored < 25 righe ciascuna)
  tui/                        # charmbracelet/huh sequential, ErrNonTTY se non-interattivo (EC-16)
features/                     # godog BDD
  main_test.go                # TestMain: build binary in tmpdir + scenario hook
  helpers_test.go             # scenarioState, fake server Ollama/EUrouter, subprocess invoke
  *_steps_test.go             # step impl in italiano con (?:che )? prefix per Gherkin
  *.feature                   # 7 file scenari
profili/                      # YAML capability shippable (revisione, rilievi, review-pack, audit-checklist, equipment-alert, competence-gap, pt-analysis)
docs/                         # meta-prompt-genera-profilo.md + guida-meta-prompt-denis.md (FR-20/21, deliverable Sprint 1)
KB-ispettore                  # symlink → docs/piano_iniziale/materiali-dominio/KB-ispettore
scripts/                      # build-mac.sh, build-linux.sh, build-zip.sh
```

## Workflow

Il progetto segue la pipeline `vibbly`:

```
spec → plan → test-scaffold → implement → review → deploy → compound
```

Ogni step ha un comando dedicato (`/v-spec`, `/v-plan`, ecc.) eseguito in Claude Code. Stato in `.pipeline/state.json`. Branch: `development` (mai commit diretti su `main`).

## Convenzioni Go

- **Logic/presentation separation**: `internal/` non importa cobra, huh, fmt.Println per user output. La presentation chiama logic, mai il contrario.
- **Errori come valori**: niente panic. Funzioni ritornano `(*T, error)`. Errori sentinella (`ErrOllamaUnreachable`, `ErrMissingAPIKey`, `ErrNonTTY`) per casi tipici.
- **Funzioni ≤ 20 righe** (`MAX_FUNCTION_LINES`). `runner.Run` è stato refactorato in step-functions per rispettarlo.
- **Accept interfaces, return structs**.
- **Context.Context come primo arg** per cose cancellabili (Stream).
- **No `init()`** salvo necessità (test setup è in `TestMain`).

## Strategy pattern: LLMProvider

`provider.LLMProvider` è l'interface unica. Per aggiungere un nuovo provider:
1. Crea `internal/provider/<nome>.go` che implementa `Name()` e `Stream(ctx, system, user, opts) (<-chan StreamEvent, error)`
2. Estendi `Select(...)` aggiungendo il nuovo case
3. Aggiorna `profile.AllowedProviders`
4. Aggiungi integration test con `httptest.NewServer`

## Mock zones

Seguiamo le zone di `01-methodology.md`:
- **Zone 1 (mai mockare)**: filesystem, profili YAML reali, KB letta da disco
- **Zone 2 (mock al boundary)**: HTTP server per provider (`httptest.NewServer`)
- **Zone 3 (real app)**: BDD acceptance invocano il binario come subprocess

## Debugging

### Ollama (Linux/WSL2 con GPU)
```
ollama serve  # avvia server
ollama list   # verifica modelli installati
ollama run qwen3.6 "ciao"  # smoke test del modello
```

### EUrouter (CTO su macOS — Ollama non disponibile)
```
export EUROUTER_API_KEY="..."
export LABNEXUS_PROVIDER=eurouter
./labnexus run --profile revisione --input ... --output ...
```

oppure `--provider eurouter` come flag (precedenza maggiore di env).

### Test dry-run senza Ollama
```
./labnexus check revisione --input docs/piano_iniziale/materiali-dominio/test-precedenti/test-1-input --show-prompt
```
Stampa il prompt composto su stdout senza chiamare alcun modello. Utile per verificare la KB iniettata e la stima token.

## Aggiungere una nuova capability (profilo)

Sprint 1: il prompt meta è una feature di Fetta 3. In manuale:
1. Crea `profili/<nome>.yml` con schema FR-3 (vedi `internal/profile/profile.go`)
2. Verifica con `./labnexus validate <nome>`
3. Esegui `./labnexus describe <nome>` per controllare i campi
4. Dry-run con `./labnexus check <nome> --input <dir>`
5. Run reale con `./labnexus run --profile <nome> --input <dir> --output <dir>`
6. Valutazione di Denis (L2)

## Pattern di validazione (NFR-8)

Ogni capability passa per due livelli:
- **L1 (CTO, automatico)**: pre-check tecnico strutturale via BDD acceptance — frontmatter completo, struttura coerente col template, no allucinazioni macroscopiche
- **L2 (Denis, manuale)**: validazione formale qualitativa secondo griglia (allucinazioni, qualità ispettiva, profondità, tenuta, utilizzabilità pratica)

L'esito di L2 va registrato nel report di sprint e nel frontmatter (`valutazione_denis`).

## Gotchas

- **`pdf.Open` su contenuto non-PDF**: succeede silenziosamente con 0 pagine. Per questo abbiamo un magic-bytes check in `extractPDF` che fail-fast se i primi 4 byte non sono `%PDF`.
- **Soglia parsing failures**: `> 0.5` (strict majority), NON `≥ 0.5`. Quindi 1 PDF rotto su 2 → procede.
- **Gherkin italiano**: `Dato che X` → testo step diventa `che X`. Tutti i pattern step in `features/*_steps_test.go` usano prefisso `(?:che )?` opzionale.
- **TTY detection**: `golang.org/x/term.IsTerminal(int(*os.File.Fd()))`. `bytes.Buffer` non è un `*os.File` → non-TTY.
- **macOS port 1**: su alcuni sistemi macOS port 1 (TCP) può accettare connessioni e ritornare HTTP 404. Per testare "Ollama down" usiamo un endpoint qualunque + HTTP 404/5xx → trattato come `ErrOllamaUnreachable` (vedi `OllamaProvider.doRequest`).
- **Symlink KB-ispettore**: relativo `docs/piano_iniziale/materiali-dominio/KB-ispettore` (no `../`). `build-zip.sh` lo risolve via `zip -ry`.

## Documenti correlati

- `CLAUDE.md` — istruzioni per AI builder agent (Golden Rules)
- `CONTEXT.md` — context per qualunque LLM agent
- `README.md` — quick-start per utente generico
- `.pipeline/spec.md` — requisiti formali Sprint 1
- `.pipeline/plan.md` — piano implementativo Fetta 1
- `.pipeline/standards/` — quality bar layered (Layer 0-3)
