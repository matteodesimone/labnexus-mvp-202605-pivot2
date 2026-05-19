# LabNexus — Guida sviluppatore

Per il CTO + futuro contributor. Build, test, layout, gotchas.

## Prerequisites

- Go 1.22+
- mise (gestione tool versions via `.mise.toml`)
- Per i test BDD: il binario viene compilato da `features/main_test.go` al `TestMain` — niente da installare a parte
- Per `labnexus run` reale: Ollama (`ollama serve`) + modello `qwen3.6` (`ollama pull qwen3.6`); o `EUROUTER_API_KEY` per provider eurouter

## Build

### Linux/amd64 (sviluppo CTO su WSL2)
```
bash scripts/build-linux.sh
./dist/labnexus-linux-amd64 --help
```

### macOS/arm64 (consegna a Denis)
```
bash scripts/build-mac.sh
open labnexus.app
```

### Pacchetto consegna (zip completo)
```
bash scripts/build-zip.sh
# produce dist/labnexus-sprint1-darwin-arm64.zip
# contiene labnexus.app + profili/ + KB-ispettore/ + README.md
```

## Test

### Tutta la suite
```
go test ./...
```

Aspettativa: 10 package `ok`, 0 fail. Su `features/` godog gira ~1.5s, le acceptance hanno tag `~@manual` (escludono scenari hardware-only Finder/drag&drop/Gatekeeper).

### Sotto-categorie

```
# Unit logic puro
go test ./internal/...

# BDD acceptance (godog)
go test ./features/

# Singolo scenario BDD (godog Tags)
go test ./features/ -godog.tags=@motore-cli
```

### Linting
```
golangci-lint run
gofmt -l .
go vet ./...
```

> **Nota**: `golangci-lint` NON è gestito da mise (il backend aqua va in 401 sul rate-limit GitHub API). Installalo a parte:
> - **macOS**: `brew install golangci-lint`
> - **Linux**: vedi <https://golangci-lint.run/welcome/install/#binaries>
> - **Cross-platform**: `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`
>
> Verifica: `golangci-lint --version` (atteso ≥ 2.0).

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
  *.feature                   # 6 file scenari
profili/                      # YAML capability shippable (revisione.yml, rilievi.yml)
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
