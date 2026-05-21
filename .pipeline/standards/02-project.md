<!-- Auto-filled by `vibbly init`. Edit freely; vibbly doctor reads the key:value lines below. -->
domain: cli-tool
Stack: go-cli

# Project Standards

## Project Context

- **What**: Eseguibile macOS nativo (Go) che esegue una capability ispettiva alla volta del modello AICertus, chiamando Qwen 3 in locale via Ollama. A corredo, un prompt meta per generare nuovi profili in autonomia e un report di validazione di Qwen sulle 7 capability.
- **Why**: rispondere alla domanda contrattuale dello Sprint 1 — *un modello locale (Qwen 3 via Ollama, context 128k) produce output di qualità accettabile per un QM esperto sulle capability LLM-critical del sistema AICertus?* — prima di investire 12-18 mesi nella costruzione del sistema completo.
- **Who**: utente operativo è **Denis Brazzo** (QM esperto SGQ + ispettore ACCREDIA, laboratorio del cliente). Stakeholder cliente: **Stefano Fiorina**. CTO interno: **Matteo De Simone**.
- **Stage**: Sprint 1 chiuso AI-side 2026-05-21; Sprint 1.5 (refactor pre-handoff Denis) **in corso**, sub-fetta 1.5.A "pivot 3 + CLI puro" implementata e in review.

### Sprint 1.5.C amendments (concierge mode + audit trail + lavori/ rename)

> **Aggiunte 1.5.C (2026-05-22)** rispetto alle voci storiche più sotto:

- **Cartella `lavori/` con auto-discovery** (FR-25, FR-28, NFR-10): la cartella `data/` è stata rinominata `lavori/`. Ogni sotto-cartella ha un `_labnexus.toml` (2 righe: `profile` + opzionale `trigger_prompt_file`) che la identifica come job auto-discovered. Fallback convention naming: regex `Profilo (.+)$` sul nome cartella se metadata manca. **NIENTE hardcoding di nomi cartelle nello script `.command` o nel codice** — Denis può aggiungere/rimuovere/rinominare lavori senza modifiche a config centrale.
- **Subcommand `labnexus jobs` + flag `--job`** (FR-28/29): `labnexus jobs` lista i lavori auto-discovered (default human-readable; `--json` per backend). `labnexus run --job <nome>` esegue headless via auto-discovery: profile + input_dir + output_dir tutti risolti dal job. Coesiste con `--profile/--input/--output` (Sprint 1 invariato).
- **Audit trail multi-writer** (FR-32/33/34, NFR-11): per ogni run, il logger scrive contemporaneamente su `os.Stderr` (live) e su `<output_dir>/<basename>.log` (audit trail file). Pattern `io.MultiWriter`. Frontmatter MD output include nuovo campo `log_file: <path-relativo>` per cross-reference bozza ↔ audit (FR-34). Engine-owned per ownership table (FR-24).
- **Body streaming live in stderr** (FR-35): durante chiamata provider, ogni token chunk viene scritto raw su stderr (TTY) e sempre sul log file. Denis vede il modello "pensare" in tempo reale.
- **Pre-check API key all'avvio** (FR-36): `labnexus.command` shell-side controlla che `labnexus.config.toml` abbia `eurouter_api_key` non vuoto (o env `EUROUTER_API_KEY` setted) PRIMA del binary fork. Messaggio operativo italiano + exit 2 se manca.

### Sprint 1.5 amendments (post-pivot-3 + decommissione bundle macOS)

> **Queste decisioni invalidano alcune voci di "Architecture Overview" e "Key Decisions Already Made" più sotto.** Le voci originali sono mantenute per tracciabilità storica; in caso di conflitto, gli amendments di Sprint 1.5 prevalgono.

- **Provider default invertito**: post-pivot-3, il deliverable per Denis usa **`provider: eurouter` (cloud EU-GDPR via `api.eurouter.ai`) di default**. Motivazione: Denis non dispone di hardware sufficiente per Ollama locale (Qwen 3 36B MoE). Cliente (Stefano Fiorina) formalmente informato e d'accordo; transito dati reali del SGQ a EUrouter approvato esplicitamente. DPA contrattuale è onere lato cliente, non blocker tecnico per la consegna. Ollama resta provider opzionale via flag/env/profile per utenti con hardware adeguato.
- **Gate cloud rimosso**: la variabile `LABNEXUS_ALLOW_CLOUD_PROVIDER` + la funzione `requireLoopbackOrCloudGate` sono stati decommissionati da `internal/provider/provider.go` in sub-fetta 1.5.A. La selezione provider (`Select`) ora si limita alla precedenza flag > env > profilo (FR-12 invariato), senza gate runtime di approvazione. Vedi `.pipeline/bugs/privacy-eurouter-gate-mancante.md` (status: `intentional_deviation_post_pivot_3`).
- **Bundle macOS `.app` decommissionato**: il deliverable Sprint 1.5 è un **binary CLI puro** standalone (`labnexus-darwin-arm64`) + un launcher minimale `scripts/labnexus.command` (2 righe bash) per il doppio click dal Finder. Eliminati: `osascript`, `Info.plist`, `CFBundleExecutable`, drag&drop di cartella sull'icona `.app`, wrapper bash. Lo zip di consegna include i due file al root insieme ai profili e KB. Workaround Gatekeeper: `xattr -d com.apple.quarantine` su entrambi i file (operazione una tantum, documentata in README + USER-GUIDE).
- **NFR-7 ribaltato**: l'affermazione "nessun dato esce dalla macchina locale" NON è più vera per il deliverable Sprint 1.5. Il marcatore audit `output.frontmatter_default.locale` nei 7 profili **deve riflettere il provider effettivo a runtime**, non essere hardcoded a `true`.

### Architecture Overview

Singolo binario Go (cross-platform `darwin/arm64` + `linux/amd64`) lanciabile da terminale o, su macOS, anche con doppio click su `scripts/labnexus.command` (launcher bash 2 righe che apre Terminal nella cartella dello zip). Per ciascuna delle 7 capability esiste un file di profilo TOML (`profili/<nome>.toml`, Sprint 1.5.B migration da YAML) che dichiara: file della KB-ispettore da iniettare come system context, `trigger_prompt` come user prompt, e opzionalmente provider/modello/parametri (ereditati dal master `labnexus.config.toml` al root se non specificati). A ogni esecuzione la KB-ispettore viene **riletta dal filesystem** (no embedding, no RAG, no DB). Provider default Sprint 1.5: **EUrouter cloud EU-GDPR** su `api.eurouter.ai` (SSE OpenAI-compatible); provider opzionale Ollama su `localhost:11434` (NDJSON streaming) per utenti con hardware adeguato. Output markdown in cartella dedicata con frontmatter YAML (convenzione Obsidian/Hugo invariata). Streaming con progress bar.

> **Storico (pre-Sprint-1.5):** il deliverable Sprint 1 era un bundle `.app` macOS con tre modalità (doppio click + `osascript` dialog, drag&drop di cartella sull'icona, riga di comando) e provider locale Ollama. Decommissionato in sub-fetta 1.5.A.

### Key Decisions Already Made

- **Locale mandatory** ~~tutta l'elaborazione AI gira sul Mac del cliente, no cloud~~ — **invalidato dal pivot 3 (Sprint 1.5, 2026-05-21)**: vedi "Sprint 1.5 amendments" sopra. La voce è preservata per tracciabilità storica. Domanda contrattuale ricalibrata in riunione Fiorina + Brazzo + CTO post-Sprint-1 (Denis non ha hardware Ollama).
- **Cross-platform Go: due target `darwin/arm64` + `linux/amd64`** — `darwin/arm64` è il deliverable per Denis (macOS Apple Silicon), `linux/amd64` è l'ambiente di sviluppo del CTO (Windows WSL2 con NVIDIA RTX 4090 24GB per Ollama accelerato). Stesso codice, due artefatti via `GOOS`/`GOARCH`. Niente UI web, niente database, niente runtime esterni da installare.
- **KB-ispettore scritta da Denis, NON da noi** — i file della KB sono trattati come fonte di verità sul tono ispettivo e i requisiti ISO 17025. L'unico testo scritto da noi per profilo è il `trigger_prompt` (5-10 righe, min 50 caratteri). Selezioniamo quali file iniettare, non li riscriviamo.
- **Ordine di sviluppo dei 7 profili è vincolato e non negoziabile** (rischio crescente, per separare problemi del motore da problemi del modello):
  1. `revisione` (Capability A) — shakedown motore su Test 1 già approvato
  2. `rilievi` (Capability B) — shakedown motore su Test 2 già approvato
  3. `review-pack` (C) — prima capability nuova
  4. `audit-checklist` (D)
  5. `equipment-alert` (E)
  6. `competence-gap` (F)
  7. `pt-analysis` (G) — ultima apposta (ragionamento numerico, lato debole degli LLM)
- **`feat-meta` viene per ultimo** — il prompt meta astrae il pattern emerso dai 7 profili reali, non lo precede.
- **Test indipendente di tenuta su context lungo è fuori eseguibile** — si esegue in chat manuale Ollama per isolare la variabile "context lungo" dal codice. Eseguito dopo i primi 2-3 profili.
- **No firma Apple Developer** in Sprint 1 — procedura Gatekeeper manuale documentata nel README. Investimento 100€/anno valutabile in Sprint 2.
- **No embedding/RAG/vector DB** — la KB-ispettore è abbastanza piccola da entrare nel context come prompt. No chunking automatico: gli input dello sprint sono dimensionati per stare dentro 128k.
- **Dual-provider** ~~Ollama-only nel deliverable~~ — **invalidato dal pivot 3 (Sprint 1.5)**: il deliverable Denis ora usa `provider: eurouter` (cloud EU-GDPR) di default; Ollama resta provider opzionale. Restrizione "nessun dato del SGQ a EUrouter senza approvazione esplicita" **ribaltata** post-pivot-3: Stefano Fiorina ha formalmente approvato il transito dei dati reali a EUrouter. Vedi "Sprint 1.5 amendments" sopra. Override del provider via flag `--provider` o env var `LABNEXUS_PROVIDER` invariato.
- **ISO 17034 fuori scope** dello sprint per decisione del cliente: KB-ispettore copre 17025 in modo solido, 17034 non è formalizzato. Tutti gli esperimenti restano su perimetro **ISO 17025**.
- **Stabilità di scope sullo Sprint 1** — il cliente ha accettato il principio: nessun nuovo ramo di sviluppo senza motivazione esplicita basata su dati raccolti.

### Pivot history (contesto stabile)

- **Assestamento iniziale** (fine settimana 1): da agente conversazionale cloud generico a strumento locale specializzato per i task ispettivi del QM. Domanda contrattuale invariata, esecuzione ricalibrata. Costo: 5 gettoni della prima fetta archiviata.
- **Pivot vero** (inizio settimana 2): emerso esplicitamente in riunione congiunta (Fiorina + Brazzo + CTO) che **tutta l'elaborazione AI deve girare in locale, non solo l'interfaccia**. La domanda contrattuale è cambiata da "LLM in cloud" a "modello locale (Qwen 3 via Ollama, context 128k)". È da qui che parte il lavoro tecnico residuo.

### Domain Vocabulary

- **gettone** — unità di complessità relativa, NON un'ora. Minimo 1, nessuna attività vale meno.
- **profilo** — file TOML in `profili/<nome>.toml` (Sprint 1.5.B; era `.yml` in Sprint 1) che configura una capability dell'eseguibile. Campi opzionali ereditati dal config master.
- **KB-ispettore** — knowledge base scritta da Denis Brazzo, iniettata come system context. Vive in `KB-ispettore/`. **Fonte canonica** durante lo sprint; iterazioni successive le fa Denis localmente.
- **trigger_prompt** — user prompt scritto dal team CTO, l'unico testo nostro per profilo. 5-10 righe, minimo 50 caratteri.
- **capability** — task LLM-critical del sistema AICertus completo. Le 7 dello sprint hanno sigle A-G:
  - A = `revisione` (DocAgent + NormAgent)
  - B = `rilievi` (CAPAAgent)
  - C = `review-pack` (ReviewAgent)
  - D = `audit-checklist` (AuditAgent)
  - E = `equipment-alert` (EquipmentAgent)
  - F = `competence-gap` (CompetenceAgent)
  - G = `pt-analysis` (PTQCAgent)
- **fetta** — unità di sviluppo verticale del framework Prove Out.
- **QM** — Quality Manager (Denis).
- **SGQ** — Sistema di Gestione della Qualità del laboratorio.
- **ACCREDIA** — ente italiano di accreditamento.
- **rilievi** — esiti di un audit ACCREDIA: NC (Non Conformità), Osservazioni, Commenti.
- **LabNexus** = nome dell'eseguibile. **AICertus** = nome del sistema/modello completo (roadmap pluriennale).
- **CoWork KB** = `LabNexus_CoWork_KB.md`, visione del prodotto finale; **WikiLLM Orchestrator** = `LabNexus_WikiLLM_Orchestrator.md`. Entrambi sono **visione di lungo periodo, NON spec di Sprint 1**: per Sprint 1 sono fonte autorevole per tono e template di output dei trigger_prompt, e per la tassonomia degli agenti del prompt meta.
- **Famiglia deterministica / LLM-critical / mista** — classificazione dei task del sistema AICertus. Sprint 1 valida solo la **LLM-critical** (più la parte LLM dei task misti). I task deterministici sono ingegneria di Sprint 2.

## Project Identity

- **Name**: LabNexus (eseguibile) — Sprint 1 del progetto AICertus
- **Client**: LabNexus s.r.l. (Stefano Fiorina, Denis Brazzo)
- **Quality metric**: per ciascuno dei 7 esperimenti, **valutazione formale di Denis** secondo griglia condivisa (allucinazioni, qualità ispettiva, profondità, tenuta, utilizzabilità pratica). Esito per capability: *validata* / *validata con riserva* / *non validata*. Per i profili `revisione` e `rilievi`, criterio aggiuntivo: output qualitativamente paragonabile ai golden file dei Test 1 e 2 manuali approvati da Denis.
- **Scale**: utente singolo (Denis) su laptop macOS Apple Silicon del laboratorio. Post-pivot-3 (Sprint 1.5): rete attiva per chiamate a `api.eurouter.ai` (provider default). Dati on-device solo se l'utente abilita esplicitamente il provider opzionale Ollama.

## Project Type

type: prove-out

<!-- Sprint 1 con budget gettoni dichiarato al cliente, gate human, deliverable di sprint. -->

## Domain

domain: cli-tool

## Stack Selection

Stack: go-cli

## UI Component Library

UI_COMPONENTS: none

## Database Choice

- **Engine**: none — tutto filesystem markdown
- **ORM**: none

## Privacy & Data Handling

- **Personal data collected**: i dati operativi del SGQ di Denis (procedure, NC, rapporti, matrici competenze, risultati PT) **possono contenere riferimenti a persone** del laboratorio. **Vincolo storico Sprint 1**: i dati non escono dalla macchina locale (Ollama). **Vincolo aggiornato Sprint 1.5** post-pivot-3 (2026-05-21): cliente (Stefano Fiorina) ha formalmente approvato il transito dei dati reali a `api.eurouter.ai` (cloud EU-GDPR). EUrouter è ora il provider default; DPA contrattuale è onere lato cliente. Ollama resta opzionale per utenti con hardware adeguato.
- **Secrets**: `EUROUTER_API_KEY` in env var, **mai hardcoded, mai committato**.

## Project-Specific Conventions

- **Layout filesystem dell'eseguibile** (path relativi alla cartella di `labnexus`):
  - `labnexus.config.toml` — config master con defaults globali (Sprint 1.5.B+, ereditati dai profili)
  - `profili/<nome>.toml` — file di configurazione per capability (Sprint 1.5.B+; era `.yml` in Sprint 1)
  - `KB-ispettore/` — knowledge base di Denis (path nei `kb_files:` dei profili sono relativi qui)
- **Nome file di output**: `<output_dir>/<timestamp>_<profile>_<input_descriptor>.md` con frontmatter YAML (data esecuzione, profilo, modello, provider, durata, token stimati, file di input).
- **Provider names ammessi**: `ollama` | `eurouter`. Default Sprint 1: `ollama`. **Default Sprint 1.5+ (post-pivot-3): `eurouter`**. Vedi "Sprint 1.5 amendments".
- **Modello default**: stringa `qwen3.6` (tag e quantizzazione esatti confermati al primo `ollama list` su ogni hardware target). Default temperature `0.9`, `max_tokens 8192`, `context_window 128000`. Parametri affinati durante shakedown del profilo `revisione`.
- **Comandi CLI**: `labnexus run --profile <name> --input <dir> --output <dir>` (principale); `list`, `describe <profile>`, `check <profile> --input <dir>` (dry-run), `validate <profile>` (controllo schema YAML, usato anche dal prompt meta).
- **Stima token**: approssimazione `char_count / 4`. Warning a 70% del `context_window`, **errore esplicito** a 100% (non chiamare il provider).
- **Lettura cartella input**: walk **non ricorsivo** (cartella piatta). Formati accettati: `.md`, `.txt`, `.csv`, `.docx`, `.pdf`, `.xlsx` — estrarre testo plain.
- **Composizione prompt**:
  ```
  system_message = concat(read(kb_files), separator="\n\n---\n\n")
  user_message   = trigger_prompt + "\n\n## File di input\n\n" + concat(input_files_as_text, sep="\n\n--- FILE: <name> ---\n\n")
  ```
- **Path dei `kb_files:`** nei profili: **relativi alla cartella KB-ispettore** (non al CWD).
- **Linguaggio**: italiano per documentazione, commit message convenzionali in italiano accettabili. Identificatori di codice in inglese (convenzione Go).
- **Trigger_prompt scritto da umano**: anche se prodotto dal prompt meta, **richiede revisione umana** prima di essere usato in produzione.

## Intentional Deviations

- **Binario non firmato Apple Developer** → procedura Gatekeeper manuale documentata. Volutamente: investimento 100€/anno rimandato a Sprint 2.
- **Nessun retry / nessun caching delle chiamate LLM** → semplicità di Sprint 1. Errore esplicito su fallimento, l'utente rilancia.
- **Logging testo, non JSON strutturato** → testo basta per Sprint 1.
- **Stima token approssimata (`char_count / 4`)** → sufficiente per il warning a 70%, non serve precisione di un vero tokenizer in Sprint 1.
- **Output Test 2 non conservato** → la validazione del profilo `rilievi` userà criteri qualitativi descritti in `feat-003` (confronto col template `CAPA_Pack` del CoWork KB sez. 15.2), non un golden file.
- **CTO assorbe internamente 3 gettoni di sforo** (budget interno 15 vs 12 dichiarati al cliente con sforo dichiarato di 1 = 13 comunicati). Per il framework il **tetto operativo reale è 15 gettoni**, non 12. Trade-off / escalation / trasferimenti a sprint successivi vanno pensati su 15.
- **Sviluppo del prompt meta DOPO i 7 profili**, non prima — astrazione su pattern emerso, non a priori. Volutamente in apparente violazione di "DRY first".
- **Niente validazione semantica automatica dell'output del prompt meta** → solo validazione di schema (`labnexus validate`). La qualità del `trigger_prompt` resta giudizio umano.

## Out of Scope

Tutto ciò che segue è **dichiaratamente fuori Sprint 1**. Materiale prezioso, parte della roadmap Sprint 2 e successivi (vedi `LabNexus_CoWork_KB.md` e `LabNexus_WikiLLM_Orchestrator.md`):

- **Pipeline deterministica** — filesystem walk del SGQ, indicizzazione documentale, mappa documentale, scadenze documenti/tarature, cross-reference moduli mancanti, riferimenti normativi obsoleti, audit trail, incrocio nomi tecnici-fornitori per imparzialità, detection celle Excel non protette, calcoli aritmetici su z-score PT.
- **Orchestrazione automatica** — watcher di cartella, trigger su filesystem, cron, pipeline multi-capability in sequenza.
- **Una esecuzione = una capability**. Nessuna composizione.
- **Interfaccia web / HITL** — niente UI web, niente human-in-the-loop conversazionale.
- **Database centralizzato** — niente DB.
- **Multi-tenancy** — utente singolo.
- **Modifica file SGQ originale** — lettura input in sola lettura, scrittura solo nella cartella output.
- **Sostituzione del giudizio del QM** — ogni output è bozza, etichettata come da approvare manualmente.
- **ISO 17034** — fuori scope dello Sprint 1 per decisione di Denis (KB non formalizzata).
- **Confronto in produzione con Claude o altri provider** — ~~EUrouter resta facility di debug interno~~ **(invalidato Sprint 1.5: EUrouter è ora provider default deliverable, vedi amendments).** Resta out-of-scope il confronto con Claude direttamente o con provider che non siano `ollama` o `eurouter`.
- **Aggiunta di nuove capability oltre le 7** — è attività di Sprint 2, anche se il prompt meta abilita Denis a esplorare candidati in autonomia tra sprint (CTO resta nel loop per validazione tecnica e messa in produzione).
- **Più di un input al lancio**, retry intelligente, sandbox/permessi avanzati, chunking automatico se input supera context.
- **Test indipendente di tenuta su context lungo** — gestito separatamente dal CTO in chat manuale Ollama, fuori dal codice dell'eseguibile e fuori dal budget gettoni di sviluppo.
- **Report finale + walkthrough conclusivo** — gestito separatamente dal CTO, fuori dal budget gettoni di sviluppo.

## Deploy

- **Target Sprint 1.5+**: singolo file zip contenente il binary CLI puro `labnexus` (Mach-O `darwin/arm64` non firmato), il launcher `labnexus.command` (2 righe bash per doppio click Finder), 7 file profilo YAML, copia controllata della KB-ispettore al momento della consegna, README di una pagina, deliverable meta-prompt (`docs/meta-prompt-genera-profilo.md` + `docs/guida-meta-prompt-denis.md`), esempi di input per ciascuna capability. **Storico Sprint 1**: target era `labnexus.app` bundle macOS (decommissionato in sub-fetta 1.5.A).
- **Environments**: dev locale (CTO, macOS Apple Silicon + Linux/WSL2) → consegna su laptop di Denis (macOS Apple Silicon). Niente staging. Post-pivot-3 (Sprint 1.5): rete attiva sul laptop di Denis per chiamate al provider default `api.eurouter.ai`.
- **CI/CD**: nessuna pipeline CI/CD. Build locale via `make build-mac` / `bash scripts/build-mac.sh` (singolo `go build` cross-compile, niente bundle). Distribuzione zip manuale.
- **Procedura Gatekeeper**: documentata nel README (control-clic → "Apri" → conferma, una sola volta).

## Tooling

- **Runtime management**: mise (`.mise.toml` presente — Go 1.22, golangci-lint, goose; `CGO_ENABLED=0`).
- **Platforms**: due target Sprint 1 — `darwin/arm64` (consegna Denis) + `linux/amd64` (sviluppo CTO su Windows WSL2). Build cross-compile Go via `GOOS`/`GOARCH`.
- **Librerie chiave previste**:
  - CLI: `spf13/cobra`
  - Progress bar: `schollz/progressbar/v3` o `vbauerster/mpb`
  - Modalità interattiva: **TUI cross-platform** (default proposto `charmbracelet/huh`, scelta finale al `/v-plan`) — gira identica su macOS e Linux/WSL2, no dipendenze osascript/Cocoa
  - Parsing: `pdfcpu` / `ledongthuc/pdf` (PDF), `unidoc/unioffice` o `nguyenthenguyen/docx` (DOCX), libreria standard per CSV/XLSX/MD/TXT
- **Provider LLM**: Ollama HTTP API su `localhost:11434` (NDJSON streaming), EUrouter SSE OpenAI-compatible su `https://api.eurouter.ai/v1/chat/completions`.
- **Health check**: `labnexus check <profile> --input <dir>` (dry-run senza chiamata LLM); `labnexus validate <profile>` (schema-only).

## Project Extensions

### Custom Commands
(none — comandi standard pipeline vibbly)

### Custom Review Checks
(none in aggiunta ai default `quality`, `test-coverage`, `security`, `privacy` già attivi in `config.sh`)

### Custom Standards
(none — si seguono gli standard `stacks/go-cli.md` e `coding-base` del framework)

---

## Note di onboarding (ambiguità e conflitti rilevati)

Annotazioni del primo onboarding (2026-05-18) da risolvere/confermare con il CTO prima del `/v-plan`:

1. **`config.sh` ha `PROJECT_TYPE="internal-tool"` ma il progetto è chiaramente `prove-out`** (budget gettoni dichiarati al cliente, sprint con deliverable concordati, fette, gate human). Va corretto a `PROJECT_TYPE="prove-out"` per allineare il framework allo stato reale. `BUDGET_TOTAL=20` e `SPRINT_DURATION_WEEKS=3` in `config.sh` sono già coerenti con il pacchetto S; vanno però rifletti **8 gettoni già spesi** e **15 gettoni come tetto operativo interno** (vs 12 dichiarati).
2. **Discrepanza numerica budget brief tecnico**: il brief (`03-brief-tecnico.md`) dichiara `gettoni_budget: 10` ma la somma in coda è `4 + 7 + 1 = 12` (motore + 7 profili + meta), e il file 00 lo conferma come `12 + 1 test tenuta + 2 report = 15`. La cifra "10" nel brief è probabilmente un refuso. **Fonte canonica: file 00-stato-e-pivot.md → 15 gettoni interno, di cui 12 sviluppo (4+7+1) coperti da questo flusso pipeline, e 3 fuori (1 test tenuta + 2 report) gestiti separatamente dal CTO.**
3. **`SHIP_CMD="goreleaser release --clean"`** è default cli-tool del framework, ma la consegna reale è uno **zip manuale** con `.app` bundle non firmato, non una GitHub release pubblica. Va sostituito (o disattivato) per Sprint 1; goreleaser può essere reintrodotto in Sprint 2 se serve.
4. **Strict ordering delle fette**: il brief impone `feat-001 → feat-002 → feat-003 → feat-004…008 → feat-meta`. Il framework dovrà rispettare questo vincolo nella fase `/v-plan`, eventualmente riconducendo singole fette a granularità più fine ma **senza riordinare**.
5. **Cartella `docs/piano_iniziale/` non è in `.pipeline/onboarding/`** (il progetto non è stato inizializzato con `vibbly init` su una base preesistente). Il materiale resta dov'è come riferimento canonico; i contenuti operativi sono stati assorbiti qui. Il **brief tecnico `03-brief-tecnico.md` è il candidato naturale per essere passato a `/v-brainstorm` o direttamente a `/v-spec`** come ingresso del flusso Prove Out.
