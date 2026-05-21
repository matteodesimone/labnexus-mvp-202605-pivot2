# LabNexus

Eseguibile macOS / Linux nativo (Go) che esegue **una capability ispettiva alla volta** del modello AICertus. **Sprint 1.5 (default deliverable)**: chiama Qwen 3 via EUrouter (cloud EU-GDPR, gateway `api.eurouter.ai`). **Opzionale**: Qwen 3 locale via Ollama per chi dispone di hardware adeguato (override via flag `--provider ollama` o env `LABNEXUS_PROVIDER=ollama`).

Sprint 1 (in corso): motore + 7 profili + prompt meta per generare nuovi profili in autonomia. Vedi `.pipeline/spec.md` e `.pipeline/plan.md`.

## Quick start

### Compilazione

Tutto via `make` (lancia `make` senza argomenti per la lista completa dei target):

```
make build        # binario per la macchina corrente → bin/labnexus
make build-linux  # cross-compile Linux/amd64 → dist/labnexus-linux-amd64
make build-mac    # cross-compile macOS arm64 → bin/labnexus-darwin-arm64
make ship         # pacchetto completo di consegna (clean + test + zip)
                  # → dist/labnexus-sprint1-darwin-arm64.zip
make smoke        # smoke test del binario su Test 1 reale (dry-run, no LLM)
```

### Uso

CLI:
```
labnexus list                         # elenca i profili installati
labnexus describe revisione           # mostra il profilo (kb_files, modello, ecc.)
labnexus validate revisione           # controlla lo schema YAML del profilo
labnexus check revisione --input ./test1   # dry-run: parsing + stima token, no LLM
labnexus check revisione --input ./test1 --show-prompt   # dry-run + stampa prompt composto
                                                          # ⚠ output contiene contenuti dei kb_files
                                                          # e degli input (potenziali PII del SGQ).
                                                          # In non-TTY emette warning su stderr.
                                                          # Solo per debug interno — NON condividere
                                                          # l'output su canali non controllati.
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
labnexus /Users/denis/Test1           # input pre-selezionato (path posizionale CLI)
```

Su macOS, doppio click su `labnexus.command` (launcher 2 righe accanto al binary) apre Terminal e lancia la TUI.

### Procedura macOS Gatekeeper (prima esecuzione)

Il binario non è firmato con certificato Apple Developer (vedi OOS-5 in `spec.md`).
Sprint 1.5.A: niente bundle `.app`; il binary è standalone. Il workflow Gatekeeper post-pivot-3 è:

**Workaround consigliato (Terminal)**:
```bash
xattr -d com.apple.quarantine labnexus labnexus.command
```
Rimuove il quarantine flag aggiunto da macOS quando il file è stato scaricato/scompattato via Finder. Operazione una tantum.

**In alternativa (Finder)**:
1. Control-clic sull'icona di `labnexus` (o `labnexus.command`) nel Finder
2. "Apri" → conferma
3. Operazione una tantum per ciascuno dei 2 file

## Provider LLM

Sprint 1.5.A (post-pivot-3): **eurouter è il provider di default** del deliverable. Cliente (Stefano Fiorina + Denis) formalmente d'accordo: i dati reali transitano a `api.eurouter.ai` (gateway cloud EU-GDPR). DPA contrattuale lato cliente.

| Provider | Quando | Come |
|---|---|---|
| `eurouter` | **Default deliverable Sprint 1.5** | Richiede chiave `EUROUTER_API_KEY` (Sprint 1.5.A: env var; Sprint 1.5.B: configurabile in `labnexus.config.toml`). Modello: `qwen3.6` su eurouter EU-GDPR gateway |
| `ollama` | Opzionale via flag/env/profile, scenari sviluppo CTO o utenti con hardware locale sufficiente | Richiede `ollama serve` in ascolto su `localhost:11434` e il modello `qwen3.6` installato |

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
| 1 | Scooter — motore + revisione + rilievi | 6 | completata e stabilizzata |
| 2 | review-pack + audit-checklist + equipment-alert | 3 | implementata (in attesa di L2 Denis) |
| 3 | competence-gap + pt-analysis + prompt meta | 3 | implementata (in attesa di esecuzione reale Qwen + L2 Denis + 3 casi inventati Claude) |
