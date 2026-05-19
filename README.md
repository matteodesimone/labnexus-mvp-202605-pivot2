# LabNexus

Eseguibile macOS / Linux nativo (Go) che esegue **una capability ispettiva alla volta** del modello AICertus, chiamando un LLM locale (Qwen 3 via Ollama) o, come canale di sviluppo, EUrouter.

Sprint 1 (in corso): motore + 7 profili + prompt meta per generare nuovi profili in autonomia. Vedi `.pipeline/spec.md` e `.pipeline/plan.md`.

## Quick start

### Compilazione

Macchina di sviluppo (Linux/WSL2 con Ollama):
```
bash scripts/build-linux.sh
./dist/labnexus-linux-amd64 --help
```

Consegna a Denis (macOS Apple Silicon):
```
bash scripts/build-zip.sh
# produce dist/labnexus-sprint1-darwin-arm64.zip
```

### Uso

CLI:
```
labnexus list                         # elenca i profili installati
labnexus describe revisione           # mostra il profilo (kb_files, modello, ecc.)
labnexus validate revisione           # controlla lo schema YAML del profilo
labnexus check revisione --input ./test1   # dry-run: parsing + stima token, no LLM
labnexus run --profile revisione --input ./test1 --output ./out
```

Override del provider (FR-12, dual-platform):
```
labnexus run --profile revisione --provider eurouter --input ./test1 --output ./out
# oppure: LABNEXUS_PROVIDER=eurouter labnexus run ...
# Il flag --provider ha precedenza sull'env var.
```

Modalità interattiva (TUI cross-platform):
```
labnexus                              # apre selettore: profilo → input → output
labnexus /Users/denis/Test1           # input pre-selezionato (drag&drop su .app)
```

### Procedura macOS Gatekeeper (prima esecuzione)

Il binario non è firmato con certificato Apple Developer (vedi OOS-5 in `spec.md`). Alla prima apertura:
1. Control-clic sull'icona di `labnexus.app` nel Finder
2. "Apri" → conferma
3. Operazione una tantum

## Provider LLM

| Provider | Quando | Come |
|---|---|---|
| `ollama` | Default Sprint 1, **deliverable Denis (locale mandatory)** | Richiede `ollama serve` in ascolto su `localhost:11434` e il modello `qwen3.6` installato |
| `eurouter` | **Solo sviluppo CTO** (su macOS dove Ollama non gira); facility di debug | Richiede env var `EUROUTER_API_KEY`; **non operativo sui dati reali del SGQ di Denis senza approvazione esplicita** |

## Layout

```
cmd/labnexus/           entry point cobra
internal/
  profile/              schema YAML + Load + Validate + List (FR-3)
  input/                parsing md/txt/csv/pdf/docx/xlsx, walk non ricorsivo (FR-4)
  prompt/               composizione system+user message (FR-5)
  tokens/               stima char/4 + warning 70% / block 100% (FR-6)
  provider/             interface LLMProvider + Ollama (NDJSON) + EUrouter (SSE) (FR-7, FR-12)
  output/               markdown + frontmatter, naming collision EC-15 (FR-9)
  runlog/               log per step + TTY detect (FR-8, NFR-4)
  runner/               orchestratore pipeline stages
  tui/                  flusso huh sequenziale, cross-platform (FR-10)
profili/                file YAML delle capability
KB-ispettore/           (symlink) knowledge base di Denis — system context
features/               scenari BDD godog (acceptance)
scripts/                build-mac.sh, build-linux.sh, build-zip.sh
.pipeline/              spec, plan, fette, state.json, test-data, standards
docs/piano_iniziale/    materiale di bootstrap (brief, piano esperimenti, KB-ispettore originale)
```

## Test

```
go test ./...
```

Unit + integration testati con `testing` stdlib + `httptest` (Ollama fake NDJSON, EUrouter fake SSE).
BDD acceptance con `godog` in `features/`.

## Status Sprint 1

Vedi `.pipeline/state.json` e `.pipeline/fette.md` per lo stato corrente delle fette.

| Fetta | Cosa | Gettoni | Stato |
|---|---|---|---|
| 1 | Scooter — motore + revisione + rilievi | 6 | in corso |
| 2 | review-pack + audit-checklist + equipment-alert | 3 | pianificata |
| 3 | competence-gap + pt-analysis + prompt meta | 3 | pianificata |
