# CONTEXT.md — LLM Project Context

> Single source of truth per qualsiasi LLM agent che lavora con questa codebase.
> Allineato al codice; non un template generico.

## Project

- **Name**: LabNexus (eseguibile macOS/Linux) — Sprint 1 del progetto AICertus
- **Description**: CLI nativo Go che esegue una capability ispettiva alla volta del modello AICertus chiamando un LLM locale (Qwen 3 via Ollama) o, come canale di sviluppo, EUrouter
- **Type**: prove-out (Sprint 1 di pacchetto S — 20 gettoni, 3 settimane)
- **Stack**: Go 1.22+ (cli-tool domain)
- **Domain**: cli-tool

## Architecture

```
cmd/labnexus/        # presentation: cobra root + 5 sottocomandi (run/list/describe/check/validate)
internal/
  profile/           # logic: schema YAML + Load + Validate + List (FR-3)
  input/             # logic: walk non ricorsivo + estrazione testo da md/txt/csv/pdf/docx/xlsx (FR-4)
  prompt/            # logic: composizione system+user message (FR-5)
  tokens/            # logic: stima char/4 + soglie warn/block (FR-6)
  provider/          # logic: Strategy pattern, LLMProvider interface + Ollama (NDJSON) + EUrouter (SSE) (FR-7, FR-12)
  output/            # logic: file markdown + frontmatter + naming collision (FR-9, EC-15)
  runlog/            # logic: log per step + TTY-aware progress (FR-8, NFR-4)
  runner/            # logic: orchestratore Pipeline stages
  tui/               # presentation: huh-based sequential flow cross-platform (FR-10)
features/            # godog BDD acceptance (7 feature files + step impl)
profili/             # YAML delle capability (revisione, rilievi, review-pack, audit-checklist, equipment-alert, competence-gap, pt-analysis)
KB-ispettore/        # symlink alla KB scritta da Denis (system context, NON modificato da noi)
scripts/             # build-mac.sh, build-linux.sh, build-zip.sh
.pipeline/           # spec, plan, fette, state.json, test-data, standards
docs/piano_iniziale/ # materiale di bootstrap (brief, KB, golden file Test 1)
```

### Key Components

- **`runner.Run`**: orchestratore. Pipeline stages: load profile → load KB → parse input → compose prompt → estimate tokens → stream LLM → write output. Refactorato in step-functions per rispettare `MAX_FUNCTION_LINES=20`.
- **`provider.LLMProvider`**: interface Strategy. Due impl: `OllamaProvider` (NDJSON streaming) ed `EurouterProvider` (SSE OpenAI-compatible).
- **`provider.Select`**: precedenza override flag `--provider` > env `LABNEXUS_PROVIDER` > campo `provider` del profilo (FR-12).
- **`input.ParseDir`**: walk non ricorsivo; soglia strict-majority `> 50%` per failure ratio (EC-1).
- **`tui.Run`**: rilevamento TTY via `golang.org/x/term`; ritorna `ErrNonTTY` se non-interattivo (EC-16).

### Data Flow

`cobra subcommand` → `runner.Config` → `runner.Run` → step pipeline → output `.md` file con frontmatter YAML.

## CLI

| Comando | Funzione |
|---|---|
| `labnexus list` | elenca profili installati |
| `labnexus describe <profile>` | mostra dettagli del profilo (provider, modello, kb_files) |
| `labnexus validate <profile>` | controlla schema YAML senza eseguire |
| `labnexus check <profile> --input <dir>` | dry-run: parsing + stima token, no LLM |
| `labnexus run --profile <name> --input <dir> --output <dir>` | esegue la capability end-to-end |
| `labnexus` (no args) | apre TUI sequenziale (cross-platform, charmbracelet/huh) |

**Flag globali**: `--profiles-dir`, `--kb-dir`, `--provider` (override). **Env var**: `LABNEXUS_PROVIDER`, `LABNEXUS_OLLAMA_ENDPOINT`, `LABNEXUS_EUROUTER_ENDPOINT`, `EUROUTER_API_KEY`.

**Exit codes** (NFR-5): `0`=successo, `1`=runtime (provider, streaming, parsing data quality), `2`=misuse (flag mancanti, profilo invalido, context window superato).

## Dependencies

| Dep | Perché |
|---|---|
| `spf13/cobra` | CLI subcommand framework idiomatic Go |
| `gopkg.in/yaml.v3` | parsing profili YAML |
| `ledongthuc/pdf` | estrazione testo PDF (pure-Go) |
| `charmbracelet/huh` | TUI sequenziale cross-platform |
| `golang.org/x/term` | rilevamento TTY |
| `cucumber/godog` | BDD acceptance test |

## Configuration

- `.pipeline/config.sh`: `PROJECT_TYPE=prove-out`, `REVIEW=all` con loop max 3, `SHIP_CMD=./scripts/build-zip.sh`, soglie `MAX_FUNCTION_LINES=20`, `MAX_CYCLOMATIC_COMPLEXITY=5`, `MAX_NESTING_DEPTH=2`.
- Profili YAML: schema definito in `internal/profile/profile.go` (FR-3). Default Sprint 1: `temperature: 0.9`, `max_tokens: 8192`, `context_window: 128000`, `modello: qwen3.6`.

## Quality Standards

Layered (`.pipeline/standards/`):
- Layer 0 `00-principles.md` — 8 principi (sottrarre, c'è un modo, decidi/fai/osserva, obsess quality, tempo fisso, sempre funzionante, sicurezza by design, doc è codice)
- Layer 1 `01-methodology.md` — TDD+BDD, mock zones, modularity balanced, intuitive interfaces
- Layer 2 `stacks/backend-cli.md` — Go conventions (no global state, internal/, errors as values, cobra patterns)
- Layer 3 `02-project.md` — LabNexus specifics (dual-provider, locale mandatory, ordering A-G, validazione L1/L2, ecc.)

**Quality metric**: per ciascuna capability, valutazione formale di Denis (L2) secondo griglia: allucinazioni, qualità ispettiva, profondità, tenuta, utilizzabilità. Pre-check tecnico CTO (L1) verifica struttura ma NON sostituisce L2.

## Current State

### Completed (Fetta 1 — Scooter)

- Motore CLI completo (run/list/describe/check/validate) con cobra
- 9 package `internal/` (profile/input/prompt/tokens/provider×3/output/runlog/runner/tui)
- 2 profili shippable: `revisione` (Capability A) + `rilievi` (Capability B)
- Bundle macOS `labnexus.app` (build-mac.sh) + binario Linux/amd64 (build-linux.sh)
- KB-ispettore symlink alla versione completa in `docs/piano_iniziale/materiali-dominio/`
- Test: tutti i package internal coperti + BDD acceptance 40/43 passing (3 pending, 10 @manual hardware-only)

### In Progress

Fetta 1 in attesa di shakedown reale su Ollama+Qwen3.6 sull'hardware del CTO (Linux/WSL2 con NVIDIA RTX 4090) e validazione formale di Denis (L2 per Capability A vs golden Test 1 e B vs template CAPA Pack).

### Planned

- Fetta 2 (3 gettoni): `review-pack` + `audit-checklist` + `equipment-alert`
- Fetta 3 (3 gettoni): `competence-gap` + `pt-analysis` + `prompt meta`

## Known Limitations

- Binario macOS NON firmato Apple Developer → procedura Gatekeeper manuale (`OOS-5`)
- TUI in non-TTY ritorna `ErrNonTTY` → bisogna usare CLI con flag
- `pt-analysis` (Capability G) accetta calcoli Qwen errati come dato di realtà; no pre-processing deterministico Sprint 1
- `ISO 17034` fuori scope dello sprint
- Una esecuzione = una capability (no orchestrazione, no watcher)
- Solo target `darwin/arm64` e `linux/amd64` (Mac Intel non in scope, trigger riunione se serve)

## How to Work With This Project

1. Read `CLAUDE.md` — istruzioni per agent builder + Golden Rules
2. Read `.pipeline/standards/` Layer 0-3 — qualità
3. Read `.pipeline/spec.md` — requisiti correnti (Sprint 1)
4. Read `.pipeline/plan.md` — piano della fetta corrente
5. Branch: lavora su `development`, NON su `main`
6. TDD: test rossi prima, implementazione poi (`/v-test-scaffold` → `/v-implement`)
7. Quality bar: NFR-8 L2 = giudizio formale di Denis sui 7 output capability
