# Plan: Sprint 1 LabNexus — Fetta 1: Scooter (motore + revisione + rilievi)

## Fetta Overview

- **FRs covered**: FR-1, FR-2, FR-3, FR-4, FR-5, FR-6, FR-7, FR-8, FR-9, FR-10, FR-11, FR-12, FR-13, FR-14 (tutti gli FR motore + i 2 profili shakedown).
- **EC covered**: EC-1, EC-2, EC-3, EC-4, EC-5, EC-6, EC-7, EC-8, EC-9, EC-10, EC-13, EC-14, EC-15, EC-16.
- **After this fetta**: il CTO (e Denis, una volta consegnato) può lanciare `labnexus run --profile {revisione, rilievi} --input <dir> --output <dir>` end-to-end, su Mac Apple Silicon e Linux/WSL2, ottenendo un file markdown con frontmatter completo. Il motore è validato contro entrambi i golden file: revisione PG_RISK_LAB rev04→rev05 e rilievi ACIAA A1.
- **Estimated effort**: 6 gettoni = stima 30–45h CTO. Complessità **Complex (> 8h)**.

## Architecture Decision

### Logic/Presentation Separation

Logica in `internal/` (pure Go, no import di cobra, huh, net/http per user-facing): parsing input, validazione schema profilo, composizione prompt, stima token, astrazione provider, scrittura output. Presentation in `cmd/labnexus/` (cobra subcommands) + `internal/tui/` (huh-based, ma resta presentation perché chiama logic). Il test: `import "labnexus/internal/profile"` deve compilare senza alcuna dipendenza CLI/TUI. Vedi sezione "Files to Create".

### Design Patterns Applied

| Pattern | Where | Justification (one line) |
|---|---|---|
| Strategy | `internal/provider` (interface `LLMProvider` + `OllamaProvider`/`EurouterProvider`) | Due provider con protocolli diversi (NDJSON vs SSE), la selezione runtime serve nativamente la dual-platform |
| Logic/Presentation separation | `internal/` vs `cmd/labnexus/` + `internal/tui/` | Mandatory dal framework; consente test della logica con il solo `go test ./internal/...` senza cobra/huh |
| Pipeline (stages) | `cmd/labnexus/run.go` orchestratore | Esecuzione = sequenza deterministica `load profile → parse input → compose prompt → estimate tokens → stream LLM → write output`; ogni stage isolato e loggato (FR-8) |
| Functional options | `internal/provider` opts | Idiomatico Go per `Stream(ctx, system, user, opts...)` |

Skill consultate (`.pipeline/skills/framework/`): `design-patterns.md`, `intuitive-design.md` (TUI cross-platform), `bdd-godog-tag-scoping.md` + `bdd-idempotent-given-steps.md` (per i feature files), `fake-exec-binary-test-pattern.md` (per testare la CLI invocata come subprocess negli step BDD), `optional-logwriter-pattern.md` (log per step), `framework-state-files-frontmatter.md` (frontmatter YAML dell'output coerente con il pattern del framework).

### Embedded Artifact Dependency Check

**Non applicabile in Fetta 1**: nessun artefatto è embedded nel binario via `go:embed`. I profili YAML, la KB-ispettore e i file di input sono **letti dal filesystem a runtime** (decisione di design del brief — consente l'iterazione di Denis sulla KB senza ricompilazione). Verifica registrata per chiusura: nessun ciclo di build da risolvere.

## Files to Create

### Modulo Go (root)

| File | Purpose | Layer | Dependencies |
|---|---|---|---|
| `go.mod` | Module declaration (`module github.com/labnexus/labnexus` o equivalente) | — | — |
| `go.sum` | Dependency lockfile | — | — |
| `.gitignore` (aggiornare) | escludi `dist/`, `*.app/Contents/MacOS/labnexus`, `.DS_Store` | — | — |
| `README.md` (aggiornare) | how-to: install, gatekeeper, comandi | docs | — |

### Logic (internal/)

| File | Purpose | Layer | Dependencies |
|---|---|---|---|
| `internal/profile/schema.go` | Struct YAML del profilo (FR-3); costanti per provider ammessi e soglia trigger_prompt | logic | gopkg.in/yaml.v3 |
| `internal/profile/validate.go` | Validazione schema (campi obbligatori, esistenza kb_files, lunghezza minima trigger_prompt, provider ∈ enum) | logic | profile, filesystem |
| `internal/profile/schema_test.go` | Unit tests | test | profile |
| `internal/input/parser.go` | Walk non ricorsivo + dispatch al parser di formato (FR-4) | logic | input/* |
| `internal/input/md.go`, `txt.go`, `csv.go` | Plain text extractors | logic | stdlib |
| `internal/input/pdf.go` | PDF extraction via `ledongthuc/pdf` (fallback warning, EC-1) | logic | ledongthuc/pdf |
| `internal/input/docx.go` | DOCX extraction via `nguyenthenguyen/docx` (EC-2) | logic | nguyenthenguyen/docx |
| `internal/input/xlsx.go` | XLSX extraction via `xuri/excelize/v2` | logic | xuri/excelize/v2 |
| `internal/input/parser_test.go` | Unit tests + edge cases | test | input |
| `internal/prompt/compose.go` | Composizione `system_message` + `user_message` (FR-5) | logic | — |
| `internal/prompt/compose_test.go` | Unit tests | test | prompt |
| `internal/tokens/estimate.go` | `char_count / 4` + soglie 70%/100% (FR-6) | logic | — |
| `internal/tokens/estimate_test.go` | Unit tests | test | tokens |
| `internal/provider/provider.go` | Interface `LLMProvider`, type `StreamEvent`, `Options`, errori sentinella | logic | context |
| `internal/provider/ollama.go` | `OllamaProvider`, NDJSON streaming su `localhost:11434` (FR-7, EC-7, EC-13) | logic | net/http, encoding/json |
| `internal/provider/eurouter.go` | `EurouterProvider`, SSE streaming OpenAI-compatible (FR-7, EC-6, EC-14) | logic | net/http, bufio |
| `internal/provider/select.go` | Risoluzione provider: profilo + `--provider` flag + `LABNEXUS_PROVIDER` env (FR-12) | logic | provider |
| `internal/provider/*_test.go` | Unit + integration (mock HTTP server, contract tests sui chunk) | test | provider, httptest |
| `internal/output/writer.go` | Naming `<timestamp>_<profile>_<descrittore>.md`, suffisso `_2` su collisione (EC-15) | logic | os, time |
| `internal/output/frontmatter.go` | Costruzione frontmatter YAML coerente con `framework-state-files-frontmatter` skill | logic | yaml.v3 |
| `internal/output/*_test.go` | Unit tests | test | output |
| `internal/runlog/log.go` | Log per step + tempi (FR-8), progress bar TTY-aware (NFR-4) | logic | io, schollz/progressbar/v3 |
| `internal/runlog/log_test.go` | Unit tests | test | runlog |
| `internal/runner/run.go` | Orchestratore pipeline: load → parse → compose → tokens → stream → write. Espone `Run(cfg) error` chiamato sia da `cobra run` sia da TUI | logic | tutte le altre internal/ |
| `internal/runner/run_test.go` | Integration: orchestrazione end-to-end con provider fake | test | runner |

### Presentation (cmd/ + internal/tui/)

| File | Purpose | Layer | Dependencies |
|---|---|---|---|
| `cmd/labnexus/main.go` | Entry point: setup cobra root, dispatch al subcommand o, se no-args, lancia TUI | presentation | cobra, tui |
| `cmd/labnexus/cmd_run.go` | `labnexus run --profile --input --output [--provider]` (FR-1, FR-12) | presentation | runner |
| `cmd/labnexus/cmd_list.go` | `labnexus list` — elenca i profili (FR-1) | presentation | profile |
| `cmd/labnexus/cmd_describe.go` | `labnexus describe <profile>` (FR-1) | presentation | profile |
| `cmd/labnexus/cmd_check.go` | `labnexus check <profile> --input <dir> [--show-prompt]` (FR-1, FR-6) | presentation | runner (dry-run) |
| `cmd/labnexus/cmd_validate.go` | `labnexus validate <profile>` (FR-3) | presentation | profile/validate |
| `cmd/labnexus/cmd_*_test.go` | godog step impl invocazione subprocess CLI (skill `fake-exec-binary-test-pattern`) | test | os/exec |
| `internal/tui/flow.go` | Flusso huh sequenziale profilo → input → output (FR-10); rilevamento TTY (EC-16) | presentation | charmbracelet/huh, term.IsTerminal |
| `internal/tui/flow_test.go` | Test con mock di prompts (huh ha helper) | test | tui |

### Artefatti runtime

| File/Dir | Purpose |
|---|---|
| `profili/revisione.yml` | Profilo `revisione` (kb_files + trigger_prompt + params) — FR-13 |
| `profili/rilievi.yml` | Profilo `rilievi` — FR-14 |
| `KB-ispettore/` | **Copia controllata** della KB di Denis al momento dell'avvio sviluppo, presa da `docs/piano_iniziale/materiali-dominio/KB-ispettore/`. Aggiornata da Denis sulla sua versione locale dopo la consegna (non da noi durante lo sprint). |
| `labnexus.app/Contents/MacOS/labnexus` | Binario darwin/arm64 (prodotto da build) |
| `labnexus.app/Contents/Info.plist` | `CFBundleExecutable=labnexus`, `LSItemContentTypes=public.folder` per drag&drop (FR-11) |
| `dist/labnexus-linux` | Binario linux/amd64 (prodotto da build) |

### Scripts

| File | Purpose |
|---|---|
| `scripts/build-mac.sh` | `GOOS=darwin GOARCH=arm64 go build -o labnexus.app/Contents/MacOS/labnexus ./cmd/labnexus` + bundle Info.plist |
| `scripts/build-linux.sh` | `GOOS=linux GOARCH=amd64 go build -o dist/labnexus ./cmd/labnexus` |
| `scripts/build-zip.sh` | Pipeline completa: build-mac → assembla zip con `labnexus.app` + `profili/` + `KB-ispettore/` + `README.md`. Usata da `SHIP_CMD` in `config.sh`. |

### Test infrastructure

| File | Purpose |
|---|---|
| `features/*.feature` | Copie da `.pipeline/proposed-features/` (curate al gate spec) |
| `features/steps/*_test.go` | Step implementations godog (skill `bdd-godog-tag-scoping`, `bdd-idempotent-given-steps`) |
| `features/main_test.go` | Setup godog (TestFeatures, scenari, init) |
| `.pipeline/test-data/fixtures/profili-validi/*.yml` | Versione di test dei profili (curata da `proposed-test-data/SAMPLE-*.yml`) |
| `.pipeline/test-data/fixtures/profili-invalidi/*.yml` | YAML invalidi per validate negativi |
| `.pipeline/test-data/fixtures/kb-ispettore-minima/` | Snapshot ridotto KB (TD-3) |
| `.pipeline/test-data/scenarios/motore-cli/input-mix-formati/` | Un file per formato (md/txt/csv/pdf/docx/xlsx) |
| `.pipeline/test-data/scenarios/motore-cli/input-vuoto/` | Cartella vuota |
| `.pipeline/test-data/scenarios/motore-cli/input-annidato/` | doc.md + sub/altro.md |
| `.pipeline/test-data/scenarios/revisione-test1/` | Riferimento a `docs/piano_iniziale/materiali-dominio/test-precedenti/test-1-input/` |
| `.pipeline/test-data/scenarios/rilievi-test2/` | Riferimento a `docs/piano_iniziale/materiali-dominio/test-precedenti/test-2-input/` |
| `.pipeline/test-data/golden/revisione/` | Riferimento a `test-1-output-approvato/` |

## Files to Modify

| File | Change | Reason |
|---|---|---|
| `.gitignore` | escludi `dist/`, `*.app/Contents/MacOS/*`, `.DS_Store`, `*.log` | artefatti di build non vanno in repo |
| `README.md` | sostituire placeholder con instructions reali (build, run, Gatekeeper) | finora c'è solo il template; va riempito al gate ship della Fetta 1 |
| `CONTEXT.md` | popolare il template (auto-generated all'init) con i dati di progetto | lo lasciamo a fine Fetta 1 (quando l'architettura è stabile) |
| `.pipeline/proposed-features/*.feature` | spostare in `features/` con eventuali correzioni emerse al gate spec | promozione da proposed a runtime tests |
| `.pipeline/proposed-test-data/SAMPLE-*.yml` | spostare in `.pipeline/test-data/fixtures/profili-validi/` | promozione da proposed a runtime tests |

## Implementation Steps (ordinati)

Pattern: ogni step apre una piccola fetta di BDD (red) + impl (green). Il driver di ogni step sono gli scenari BDD in `.pipeline/proposed-features/` (poi `features/`).

### Sub-fetta 1.A — Motore base (4 gettoni)

1. **Setup modulo Go**
   - `go mod init <module-path>`; dichiara deps principali (cobra, yaml.v3, ledongthuc/pdf, nguyenthenguyen/docx, xuri/excelize, schollz/progressbar/v3, charmbracelet/huh).
   - Crea struttura `cmd/labnexus/`, `internal/`, `features/`, `scripts/`, `profili/`.
   - **Atteso**: `go build ./...` passa; `go test ./...` passa (nessun test).

2. **Schema profilo + comando `validate`** (FR-3)
   - BDD: `motore-cli.feature` scenari `validate accetta…`, `validate rifiuta trigger corto`, `validate rifiuta kb_files inesistente`, `validate rifiuta provider non ammesso`.
   - Impl: `internal/profile/{schema.go, validate.go}` + `cmd/labnexus/cmd_validate.go`.
   - **Atteso**: 4 scenari BDD verdi; unit test schema verde.

3. **Comandi `list` e `describe`** (FR-1, FR-2)
   - BDD: `motore-cli.feature` scenari `list mostra profili`, `describe espone documentazione`, `profilo inesistente fallisce loud`, `flag obbligatori enforced`.
   - Impl: due cobra commands che scansionano `./profili/*.yml`, parsano nome + descrizione.
   - **Atteso**: 4 scenari BDD verdi.

4. **Parsing input multi-formato** (FR-4, EC-1, EC-2, EC-9)
   - BDD: scenari `check estrae testo da MD/TXT/CSV/PDF/DOCX/XLSX`, `formato non supportato skip warning`, `walk non ricorsivo`, EC-1 `PDF rotto`, EC-2 `DOCX corrotto`, EC-9 `cartella vuota` / `solo non supportati`.
   - Impl: `internal/input/{parser.go, md.go, txt.go, csv.go, pdf.go, docx.go, xlsx.go}` + `cmd/labnexus/cmd_check.go` (dry-run base).
   - **Atteso**: ~8 scenari verdi.

5. **Composizione prompt** (FR-5)
   - BDD: scenario `check produce anteprima del prompt composto` (flag `--show-prompt`).
   - Impl: `internal/prompt/compose.go`.
   - **Atteso**: 1 scenario + unit test verde.

6. **Stima token + soglie** (FR-6, EC-3)
   - BDD: `warning sopra 70%`, `errore sopra 100%`, EC-3 `context al 100%`.
   - Impl: `internal/tokens/estimate.go` + integrazione in `cmd_check.go` e `cmd_run.go` (gate pre-call).
   - **Atteso**: 3 scenari verdi.

7. **Provider Ollama (NDJSON streaming)** (FR-7, EC-7, EC-13)
   - BDD: `streaming Ollama produce token incrementali`, `Ollama non disponibile`, EC-13 `chunk malformato isolato` / `≥10 consecutivi`.
   - Impl: `internal/provider/{provider.go, ollama.go}` + test con `net/http/httptest` che simula NDJSON.
   - **Atteso**: 4 scenari verdi. Streaming TTFB < 5s su Qwen reale (NFR-3) — verificato sul box di sviluppo Linux/WSL2 con GPU.

8. **Provider EUrouter (SSE)** (FR-7, FR-12, EC-6, EC-14)
   - BDD: `streaming EUrouter rispetta SSE`, `provider eurouter senza API key`, EC-14 `timeout SSE senza DONE`.
   - Impl: `internal/provider/eurouter.go` + test SSE httptest server. Verifica env `EUROUTER_API_KEY` in `cmd_run.go` pre-call.
   - **Atteso**: 3 scenari verdi.

9. **Override provider** (FR-12)
   - BDD: `override --provider`, `override LABNEXUS_PROVIDER env`, `flag ha precedenza su env`.
   - Impl: `internal/provider/select.go` (precedenza: flag > env > profilo); flag globale `--provider` in cobra.
   - **Atteso**: 3 scenari verdi.

10. **Output markdown + frontmatter** (FR-9, EC-15)
    - BDD: `file output con naming convenzionale`, `frontmatter completo`, EC-15 `collisione → suffisso`, EC-8 `streaming interrotto → output parziale`.
    - Impl: `internal/output/{writer.go, frontmatter.go}`.
    - **Atteso**: 4 scenari verdi.

11. **Log per step + progress bar** (FR-8, NFR-4)
    - BDD: `log per step con tempi`, `progress bar in TTY`, `niente ANSI in non-TTY`.
    - Impl: `internal/runlog/log.go` (TTY-aware progress, log strutturato per step a stderr).
    - **Atteso**: 3 scenari verdi.

12. **Comando `run` orchestratore** (FR-1, FR-7, NFR-5)
    - BDD: `flag obbligatori`, `streaming Ollama produce token`, scenari profili (vedi step 15-16 sotto).
    - Impl: `cmd/labnexus/cmd_run.go` + `internal/runner/run.go` (pipeline stages chiamate in sequenza).
    - Exit codes: 0/1/2 secondo NFR-5.
    - **Atteso**: smoke test end-to-end con provider fake.

13. **Modalità interattiva TUI** (FR-10, EC-16)
    - BDD: scenari `lancio senza argomenti apre TUI`, `argomento posizionale = input pre-selezionato`, `non-TTY rifiuta`, `funziona uguale su macOS e Linux`.
    - Impl: `internal/tui/flow.go` con `charmbracelet/huh` (singolo form sequenziale Select + FilePicker + FilePicker); rilevamento TTY via `golang.org/x/term`.
    - **Atteso**: 5 scenari verdi.

14. **Bundle macOS + binario Linux** (FR-11)
    - BDD: scenari `bundle .app contiene arm64`, `doppio click apre Terminal con TUI`, `drag&drop pre-seleziona input`, `bundle non firmato → Gatekeeper documentato`, `binario Linux standalone`, `CLI puro identico su entrambi`.
    - Impl: `scripts/build-mac.sh` (compila + assembla `.app/Contents/{MacOS,Info.plist}`), `scripts/build-linux.sh`, `scripts/build-zip.sh`. `Info.plist` con `LSItemContentTypes=public.folder` per drag&drop.
    - **Atteso**: 6 scenari verdi.

### Pattern di validazione (applicabile a TUTTE le capability A-G)

Sotto-fette 1.B, 1.C e tutte le sotto-fette delle Fette 2 e 3 seguono lo stesso pattern di validazione a due livelli (vedi NFR-8 della spec e `fette.md`):

- **Livello 1 — Pre-check tecnico CTO (strutturale)**: l'output è ben formato per la capability (frontmatter completo, struttura coerente col template, no allucinazioni macroscopiche). Se fallisce → STOP, debug pipeline. NON sottomettere a Denis.
- **Livello 2 — Validazione formale Denis (qualitativa)**: solo dopo L1 OK, Denis valuta secondo griglia (allucinazioni, qualità ispettiva, profondità, tenuta, utilizzabilità pratica). Esito *validata* / *validata con riserva* / *non validata* registrato nel report di sprint e nel frontmatter (`valutazione_denis`). **Nessuna capability è "validata" senza il giudizio di Denis.**

I gate `/v-deploy` di ciascuna fetta non chiudono finché Denis non ha completato L2 per tutti i profili contenuti. Lo scaffolding della fetta successiva può procedere in parallelo all'attesa del giudizio (OQ-5).

### Sub-fetta 1.B — Profilo `revisione` (1 gettone) — shakedown motore #1

15. **Scrittura del profilo `revisione`** (FR-13)
    - Promuovi `proposed-test-data/SAMPLE-revisione.yml` → `profili/revisione.yml`.
    - Verifica `labnexus validate revisione` → exit 0.
    - **Atteso**: profilo valido, schema OK.

16. **Esecuzione su Test 1 + due livelli di validazione** (FR-13)
    - BDD: `profilo revisione replica Test 1`, `confronto qualitativo col golden`.
    - Esegui `labnexus run --profile revisione --input docs/piano_iniziale/materiali-dominio/test-precedenti/test-1-input --output <out>` su Qwen 3.6 locale (Linux/WSL2 con GPU 4090).

    **Livello 1 — Pre-check tecnico del CTO (strutturale, NON qualitativo).** Verifica oggettiva che il software non abbia rotto qualcosa:
      - file di output esiste, frontmatter YAML completo (FR-9)
      - sezioni di sintesi presenti, callout `[!MODIFICA]` presenti
      - tono identificabile come ispettivo (no "fammi sapere se posso aiutarti", no markdown segnaposto), RT-08 rev03 e rev05 citati esplicitamente
      - **nessuna allucinazione macroscopica** (es. citazioni di norme inesistenti nei file forniti)
    Se il pre-check **fallisce** → **STOP**, debug della pipeline software (prompt assembly, ordine kb_files, parsing degli input, parametri Qwen). Il CTO NON sottomette a Denis un output che non passa il pre-check: è perdere il tempo del cliente su un bug del codice.

    **Livello 2 — Validazione formale di Denis (qualitativa, conformity criterion NFR-8).** Solo dopo che il Livello 1 passa, l'output viene sottomesso a Denis che lo valuta secondo la **griglia condivisa**: allucinazioni, qualità ispettiva, profondità, tenuta, utilizzabilità pratica. Esito ufficiale: *validata* / *validata con riserva* / *non validata*. Il giudizio viene registrato nel report di sprint e nel frontmatter dell'output (campo `valutazione_denis: ...`).

    **Cosa fa scattare il vero gate STOP della Fetta 1**:
      - Se Denis dice *validata* o *validata con riserva* → motore confermato per Capability A, si procede a 1.C (profilo `rilievi`).
      - Se Denis dice *non validata* → il CTO **non si butta a debuggare alla cieca**. Apre un'analisi a tre vie: (a) è problema di motore/codice? (improbabile se il Livello 1 è passato, ma da rivedere), (b) è problema del trigger_prompt? (riformulare), (c) è problema del modello Qwen su questa capability? (improbabile, perché il Test 1 manuale era già stato approvato — ma se accade vuol dire che i parametri di Qwen sono diversi tra chat manuale e motore, es. `qwen3.6` registry vs run manuale). Si rinegozia il piano col cliente prima di procedere a 1.C.

    **Atteso**: 2 scenari `profili.feature` verdi (Livello 1) + giudizio formale di Denis registrato (Livello 2). Lo scaffolding di 1.C (profilo `rilievi`) può procedere in parallelo all'attesa del giudizio di Denis (vedi OQ-5: scaffolding parallelizzabile, esecuzione/valutazione strict).

### Sub-fetta 1.C — Profilo `rilievi` (1 gettone) — shakedown motore #2

17. **Scrittura del profilo `rilievi`** (FR-14)
    - Scrivi `profili/rilievi.yml` (kb_files: CLAUDE.md + come-pensa-un-ispettore + NC-patterns; trigger_prompt allineato al template CAPA Pack del CoWork KB sez. 15.2).
    - Verifica `labnexus validate rilievi` → exit 0.

18. **Esecuzione su Test 2 + due livelli di validazione** (FR-14, TD-2)
    - BDD: `profilo rilievi produce CAPA Pack per ogni riga del CSV`.
    - Esegui `labnexus run --profile rilievi --input docs/piano_iniziale/materiali-dominio/test-precedenti/test-2-input --output <out>` su Qwen 3.6.

    **Livello 1 — Pre-check tecnico del CTO (strutturale).** Verifica che:
      - per ogni riga del CSV ACIAA A1 esiste un blocco CAPA Pack nell'output
      - ogni blocco rispetta la **struttura del template CAPA Pack** (CoWork KB sez. 15.2): meccanismo, estensione, efficacia
      - per ogni rilievo è esplicitato se basta correzione o serve azione correttiva
      - nessuna allucinazione macroscopica (codici/sezioni inventati rispetto al CSV di input)
    Se il pre-check **fallisce** → **STOP**, debug. Non sottomettere a Denis un output strutturalmente rotto.

    **Livello 2 — Validazione formale di Denis (qualitativa, conformity criterion NFR-8).** Sottomettere a Denis per giudizio secondo griglia. Esito *validata* / *validata con riserva* / *non validata* registrato. Niente golden file (TD-2: output Test 2 non conservato); il giudizio si basa sulla **qualità dei CAPA Pack come strumenti operativi**, non sul confronto con un riferimento esistente.

    **Cosa fa scattare il gate STOP della Fetta 1**:
      - Se Denis dice *validata* o *validata con riserva* per **entrambe** le capability (A e B) → motore validato definitivamente → fetta 1 chiusa, si procede a `/v-review` e poi `/v-deploy`.
      - Se Denis dice *non validata* su B (ma A era OK) → analisi: problema del trigger_prompt di `rilievi`? Problema del modello su task di NC? Rinegoziazione del piano (es. ulteriore iterazione su `rilievi.yml` fuori budget, o passaggio a Fetta 2 con Capability B come riserva da rivedere). Non si butta sul debug software cieco.

    **Atteso**: 1 scenario `profili.feature` verde (Livello 1) + giudizio formale di Denis registrato (Livello 2).

## Migration Steps

**Nessuna**. Il progetto non ha database (Database Choice in `02-project.md`: `Engine: none`, `ORM: none`). Tutto è filesystem.

## Test Strategy

Pattern: ogni step di implementazione apre la sua slice di scenari BDD (in `.pipeline/proposed-features/`, poi promossi a `features/`) come "red phase". L'impl segue come "green phase". Refactor preserva il verde di TUTTI gli scenari.

### Test layers per Fetta 1

| Layer | Strumento | Scope |
|---|---|---|
| **Acceptance (BDD)** | `godog` | Tutti i 7 file `.feature` (motore-cli, motore-providers, motore-output, motore-modo-interattivo, profili — solo scenari `revisione` e `rilievi` in Fetta 1, prompt-meta NON in Fetta 1, edge-cases). Step impl in `features/steps/`. CLI invocata come subprocess (skill `fake-exec-binary-test-pattern`). |
| **Integration** | `testing` stdlib | `internal/runner`, `internal/provider/{ollama,eurouter}` contro `httptest.NewServer` con NDJSON/SSE fixture. `internal/input` su veri file di test (PDF, DOCX, XLSX in `.pipeline/test-data/scenarios/motore-cli/input-mix-formati/`). |
| **Unit** | `testing` stdlib + table-driven | Logica pura: schema validate, prompt compose, token estimate, frontmatter render, output naming + collision. |

### Test data fixtures necessari

- **Profili validi**: `revisione.yml`, `rilievi.yml`, `minimal.yml` (motore in isolamento)
- **Profili invalidi**: 4 YAML per i 4 fallimenti di validate (trigger corto, kb mancante, provider errato, campo obbligatorio mancante)
- **Input mix formati**: 6 file dummy uno per formato
- **Input edge**: cartella vuota, cartella con solo `.jpg`, cartella annidata, PDF intenzionalmente rotto, DOCX intenzionalmente corrotto
- **KB-ispettore-minima**: snapshot ridotto (3-4 file) per i test del motore (TD-3)
- **Test 1 e 2 reali**: riferimento (no copia) ai golden file in `docs/piano_iniziale/materiali-dominio/test-precedenti/`
- **Provider fixtures**: registrazioni NDJSON e SSE pronte per il replay nei test di provider

### Coverage atteso

- BDD: tutti gli scenari relativi agli FR/EC coperti dalla Fetta 1 → verdi
- Unit: ≥90% sulle funzioni pure (schema validate, prompt compose, token estimate, frontmatter)
- Integration: provider Ollama e EUrouter verdi su almeno il happy path + 2 errori (no API key, timeout)
- Acceptance shakedown: confronto qualitativo CTO documentato per `revisione` (vs golden) e `rilievi` (vs template)

## Risk Assessment

| Rischio | Probabilità | Impatto | Mitigazione |
|---|---|---|---|
| **Output `revisione` o `rilievi` strutturalmente difettoso** (pre-check L1 fallisce) | Media | Alto | Pre-check tecnico CTO obbligatorio prima di sottomettere a Denis (step 16/18). Se L1 fallisce → debug pipeline (prompt assembly, ordine kb_files, parametri Qwen, parsing input) **prima** di chiamare Denis. Niente "magari Denis lo trova accettabile lo stesso": è perdere tempo del cliente su un bug nostro. |
| **Denis dichiara *non validata* uno dei due profili** (validazione L2 negativa) | Media | Alto | Il giudizio formale di Denis è il conformity criterion (NFR-8). Se *non validata*, **non si debugga alla cieca**: si apre un'analisi a tre vie (motore, trigger_prompt, modello) e si rinegozia il piano col cliente. La fetta 1 può chiudere "validata con riserva" se Denis lo accetta, o restare aperta finché non si trova la causa. **Dipendenza temporale**: il giudizio richiede tempo di Denis — fissare il checkpoint in calendario prima di iniziare l'esecuzione, non dopo. |
| **Streaming Ollama NDJSON: chunk parsing fragile** | Media | Medio | Test con `httptest` su NDJSON fixture registrate. Tollerare singolo chunk malformato (EC-13), abort dopo 10 consecutivi. |
| **`charmbracelet/huh` introduce molte transitive deps** | Bassa | Medio | Valutare al primo `go mod tidy`. Fallback semplice: `manifoldco/promptui` o flusso minimal con `bufio.Scanner`. Decisione **prima** dello step 13. |
| **PDF complessi spaccano `ledongthuc/pdf`** (EC-1) | Alta | Basso | Gestito by design: skip + warning + suggerimento pandoc. Se ≥50% dei file fallisce, exit 1 con messaggio chiaro. |
| **Cross-compile darwin/arm64 da Linux WSL2** | Bassa | Basso | Go ha cross-compile nativo. `CGO_ENABLED=0` (già in `.mise.toml`) elimina problemi di toolchain C; verificare che `pdfcpu`/`ledongthuc/pdf`/`excelize` siano pure-Go (la maggior parte lo è). Se una dep richiede CGO, fare swap precoce con un'alternativa pure-Go. |
| **`qwen3.6` non disponibile in registry Ollama del CTO** | Media | Basso | Verifica al primo `ollama list`. Fallback al primo `qwen3*` presente (OQ-2). Registrare nel frontmatter il modello effettivamente usato. |
| **Drag&drop `.app` macOS non funziona su `Info.plist` minimale** | Media | Basso | Test reale sul Mac Apple Silicon (Denis o CTO se ha accesso). Se non funziona, fallback: lancio dal Finder = doppio click → TUI parte dall'inizio (input chiesto in TUI). Drag&drop è bonus, non blocking. |
| **Temperature 0.9 produce allucinazioni eccessive** | Media | Alto | È un parametro CONFIGURABILE nel profilo (FR-3). Se allo shakedown `revisione` l'output è troppo dispersivo, abbassare a 0.5 e ri-runnare. Il valore 0.9 è ipotesi del CTO da validare empiricamente. |
| **`EUROUTER_API_KEY` esposta per errore in log** | Bassa | Critico | Audit dei log nel runner: nessun `os.Environ()`, nessun dump di header HTTP, nessun echo dei flag globali contenenti chiavi. Scan dei log in `/v-review`. |
| **`go.mod` versioning di librerie BDD/TUI cambia API durante lo sprint** | Bassa | Medio | Lock-in alle versioni testate; nessun `go get -u` aggressivo. |

## Estimated Complexity

**Complex (> 8h).** Realisticamente 30-45h di lavoro CTO distribuite su circa 10 giorni lavorativi del residuo sprint. La parte pesante è l'integrazione cross-platform (TUI + bundle .app) e il debug dello shakedown comparativo col golden file (step 16 — il vero "gate di realtà" del motore).

### Distribuzione gettoni interna proposta (`allocated_per_fetta`)

| Sotto-fetta | Gettoni | Allocazione infra/client |
|---|---|---|
| 1.A motore | 4 | infra 2 / client 2 (50/50 — l'infra è significativa per CLI/TUI/bundle, ma il motore è il prodotto core) |
| 1.B profilo revisione | 1 | infra 0 / client 1 |
| 1.C profilo rilievi | 1 | infra 0 / client 1 |
| **Totale Fetta 1** | **6** | **infra 2 / client 4** |

Allocazione precisa registrata in `state.json` al gate plan.
