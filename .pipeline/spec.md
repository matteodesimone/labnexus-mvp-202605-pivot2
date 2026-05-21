# Spec: LabNexus Sprint 1.5 — refactor pre-handoff Denis (pivot 3 + concierge mode)

> Project type: **prove-out** (Sprint 1.5 = refactor a 0 gettoni cliente, ~23-36h CTO offerti dal CTO; budget tracking infra-side).
> Discovery Goal: Sprint 1 chiuso ma deliverable non eseguibile sul Mac di Denis (no hardware Ollama). Cliente informato e d'accordo a pivottare su eurouter cloud EU-GDPR di default. Sprint 1.5 chiude il gap pre-handoff.
> Brainstorm di riferimento: `.pipeline/brainstorms/2026-05-21-refactor-sprint1-concierge-mode.md`.
> Idea assorbita: `.pipeline/ideas/2026-05-21-frontmatter-key-ownership-table.md` (status: in-progress → integrata in 1.5.B config master).

## Summary

Refactor strutturale del deliverable Sprint 1 LabNexus per renderlo **eseguibile sul Mac di Denis** (no hardware Ollama → eurouter cloud EU-GDPR default) e **utilizzabile come concierge preliminare** di lavoro quotidiano per Quality Manager + ispettore ACCREDIA. Tre sub-fette ordinate:
- **1.5.A** — Pivot 3: eurouter default + CLI puro multipiattaforma (no `labnexus.app` bundle, sì `labnexus.command` minimal launcher).
- **1.5.B** — Config master TOML + override locale dei profili + parser nuovi formati Office (`.xls`, `.doc`, `.rtf`).
- **1.5.C** — Concierge mode: cartella `lavori/` con auto-discovery via `_labnexus.toml`, subcommand `labnexus jobs`, flag `--job <name>`, logger come audit trail ISO 17025, body streaming live in console e log.

## Actors

- **Denis Brazzo** (primary). QM laboratorio ISO 17025 + ispettore ACCREDIA. Uso quotidiano come concierge preliminare per bozze di revisione, gap report, audit checklist da consumare sul serio. Non programmatore, ma metodico e abituato a documentazione tecnica. Sistema d'esercizio: Mac Apple Silicon standard (8-16 GB RAM, no GPU dedicata).
- **Stefano Fiorina** (cliente). Informato e d'accordo al pivot 3. Paga API key eurouter. DPA da formalizzare lato contrattuale (fuori scope tecnico).
- **Matteo De Simone** (CTO). Sviluppo del refactor + futuro server backend Sprint 2 (uso headless di `labnexus`).
- **Backend Sprint 2** (futuro, non-utente attuale). Invocherà `labnexus` in modalità headless via `run --profile X --input Y --output Z` o `run --job <name>`. Niente TUI, niente interattività.
- **Claude esterno** (FR-22 Fetta 3, invariato). Consumer del meta-prompt per generare nuovi profili.

## Functional Requirements

### Sub-fetta 1.5.A — Pivot 3 + CLI puro (~5-7h)

- **FR-1** — Eliminare `labnexus.app` bundle macOS. Audit: bundle è feature genuina (wrapper bash + osascript + Info.plist + LSHandlerRank=Folder per drag&drop). **Decisione strategica di rimozione**, non dead code. Migration verso CLI puro standalone. Files toccati: `scripts/build-mac.sh` (18 referenze a `.app/osascript/Info.plist/CFBundleExecutable/LSHandlerRank` da rimuovere), `scripts/build-zip.sh` (layout zip senza `.app`).
- **FR-2** — `scripts/build-mac.sh` cross-compila un singolo binary Mach-O arm64 `labnexus` (no bundle, no wrapper). `scripts/build-linux.sh` resta sostanzialmente invariato (già produceva binary CLI puro).
- **FR-3** — `labnexus.command` minimal launcher (2 righe bash): `cd "$(dirname "$0")" && ./labnexus`. Doppio click dal Finder macOS apre Terminal e lancia il binary; la TUI esistente parte automaticamente (FR-10 Sprint 1 invariato).
- **FR-4** — Audit confermato: rimuovere `ErrEurouterGateMissing` (feature genuina ex-CRITICAL Fetta 2, decommissione strategica post-pivot-3) + `requireLoopbackOrCloudGate` (idem ex-HIGH Fetta 2). Coerente con la decisione "cliente decide via configurazione, niente runtime gate".
- **FR-5** — Rimuovere i test associati ai gate in `internal/provider/select_test.go`: 4 test eurouter gate (TestSelect_EurouterRejectedWithoutApprovalGate + 3 varianti) + 6 test ollama loopback (TestSelect_OllamaRejectsRemoteEndpoint + 5 varianti).
- **FR-6** — `features/helpers_test.go::buildEnv` non setta più `LABNEXUS_ALLOW_CLOUD_PROVIDER` (no più gate da bypassare).
- **FR-7** — I 7 profili shippati (`profili/{revisione,rilievi,review-pack,audit-checklist,equipment-alert,competence-gap,pt-analysis}.yml`) cambiano `provider: ollama` → `provider: eurouter` come default transitorio (verrà spostato in `labnexus.config.toml` master in 1.5.B).
- **FR-8** — Bug `.pipeline/bugs/privacy-eurouter-gate-mancante.md` aggiornato a `status: intentional_deviation_post_pivot_3` con nota di rationale. Stesso pattern per `.pipeline/bugs/ollama-endpoint-env-override-leak.md`.
- **FR-9** — BDD scenari `@manual` su `.app` bundle / Info.plist / drag&drop (7+ scenari in `features/motore-modo-interattivo.feature` + altri sparsi) **eliminati** (audit: erano `godogPending`, dead-ish in CI; rimozione coerente con rimozione del bundle).
- **FR-10** — README zip-level (rilevante per Denis): rimuovere sezione Gatekeeper-per-`.app`, aggiungere Gatekeeper-per-binary nudo (workaround `xattr -d com.apple.quarantine labnexus` se Finder ha aggiunto quarantine flag). Tono italiano, concierge.

### Sub-fetta 1.5.B — Config master TOML + override + parser nuovi formati (~9-13h)

- **FR-11** — Nuovo file `labnexus.config.toml` (config master) al root dello zip. Contiene defaults globali: `provider`, `modello`, `temperature`, `max_tokens`, `context_window`, `eurouter_api_key`, `ollama_endpoint`. Editato da Denis a mano dopo estrazione zip (almeno per inserire `eurouter_api_key`).
- **FR-12** — Migration formato user-edited da YAML a TOML: `labnexus.config.toml` + `profili/*.toml` (7 file rinominati e riscritti). Lib Go: `github.com/BurntSushi/toml`. Frontmatter Markdown output **resta YAML** (convenzione Obsidian/Jekyll/Hugo, Denis lo legge non edita).
- **FR-13** — `internal/profile/profile.go`: i campi `Provider`, `Modello`, `Temperature`, `MaxTokens`, `ContextWindow` diventano **opzionali** nel profile YAML/TOML (zero-value detection). Se vuoti, vengono ereditati dal config master.
- **FR-14** — Nuovo pacchetto `internal/config/`: load + parse `labnexus.config.toml` master. Funzioni `config.Load(path)` + `config.Merge(master, profile)` per merge profile-over-master con precedenza profile (se field setted).
- **FR-15** — `internal/runner.Run` orchestra: (1) load master → (2) load profile → (3) merge → (4) validate merged → (5) esegue. `profile.Validate` valida il profilo MERGED, non isolato (un profile può avere `provider` vuoto se master lo fornisce).
- **FR-16** — Aggiunto campo `trigger_prompt_file` al profile schema. Mutually exclusive con `trigger_prompt` inline (XOR validation: errore se entrambi setted, errore se entrambi assenti). Path relativo a `--input` dir (cartella del lavoro). Path-safety: validazione `filepath.EvalSymlinks` per impedire escape (NFR-6 esteso).
- **FR-17** — Se `trigger_prompt_file` è setted, `internal/input` parsa il file referenziato e usa il testo come trigger per quella esecuzione. Supporta formati `.rtf`, `.txt`, `.md`. Lo stesso parser è riusato anche per i file di input domain-specific.
- **FR-18** — Parser `.xls` (Excel 97-2003 OLE binary): integrazione in `internal/input/parse.go`. Lib Go: `shakinm/xlsReader` o `extrame/xls`. Estrae testo da tutte le sheet con separator `\n\n--- SHEET: <name> ---\n\n`.
- **FR-19** — Parser `.doc` (Word 97-2003 OLE binary): integrazione in `internal/input/parse.go`. Lib Go: `richardlehane/mscfb` + parser custom per il content stream (lib Go meno matura di `.xlsx`/`.docx` ma fattibile). Fallback in caso di file corrotto: warning + skip (pattern esistente EC-1 Sprint 1).
- **FR-20** — Parser `.rtf`: integrazione in `internal/input/parse.go`. Estrae plain text da control words RTF (strip formattazione). Implementazione custom semplice (RTF è formattato testuale) o lib esterna (`scelto/go-rtf-parser` se disponibile).
- **FR-21** — Spec amendment formati supportati (sostituisce FR-2 Sprint 1): `.md`, `.txt`, `.csv`, `.docx`, `.pdf`, `.xlsx`, **`.xls`**, **`.doc`**, **`.rtf`**.
- **FR-22** — Aggiornare `docs/meta-prompt-genera-profilo.md`: schema in TOML invece di YAML, 3 esempi few-shot ricalibrati in TOML, riferimenti a `trigger_prompt_file` come pattern opzionale, lista kb_files invariata.
- **FR-23** — Test suite update: `features/steps_test.go::minimalProfileYAML` rinominato `minimalProfileTOML`, output stub TOML. BDD helpers che generano profili dinamici aggiornati. Unit test `internal/profile/profile_test.go` + nuovo `internal/config/config_test.go`.
- **FR-24** — Tabella ownership chiavi frontmatter (idea assorbita da `.pipeline/ideas/...frontmatter-key-ownership...`): documentata nel commento di `internal/output/output.go::Frontmatter`. Engine-owned vs profile-default-owned vs external-owned (Denis L2). Previene future collisioni come `stato` (smascherata Fetta 2 review).

### Sub-fetta 1.5.C — Concierge mode + auto-discovery + audit log + streaming (~9-16h)

- **FR-25** — Cartella `lavori/` shippata nello zip. Rinomina della cartella `data/` esistente (Denis-facing). Ogni `lavori/<X>/` contiene gli input domain-specific + un file metadata + (dopo esecuzione) una sottocartella `output/`.
- **FR-26** — File metadata `_labnexus.toml` (2 righe minime) dentro ogni cartella `lavori/<X>/`:
  ```toml
  profile = "<nome-profilo>"
  trigger_prompt_file = "<path-relativo>"   # opzionale
  ```
  Identifica univocamente il job (quale profile usare, eventualmente quale file di prompt).
- **FR-27** — Fallback convention naming: se `_labnexus.toml` manca, `labnexus jobs` tenta regex `Profilo (.+)$` sul nome cartella per estrarre il profilo. Se nessun match, cartella skippata con warning.
- **FR-28** — Nuovo subcommand `labnexus jobs`: scansiona `lavori/*/`, legge `_labnexus.toml` (o applica fallback), produce lista jobs disponibili. Output human-readable di default; flag `--json` per output strutturato (backend-friendly).
- **FR-29** — Nuovo flag `labnexus run --job <name>`: cerca via auto-discovery un job con nome `<name>`, risolve `profile` + `input_dir` + `output_dir` (default `<input_dir>/output`), esegue headless. Coesiste con `--profile X --input Y --output Z` (Sprint 1 invariato).
- **FR-30** — Output dentro `lavori/<X>/output/` (default se `--output` non specificato). Niente parsing "lettera + profilo" — il path è risolto dalla cartella stessa. Naming file esistente preservato: `<timestamp>_<profile>_<input_descriptor>.md`.
- **FR-31** — TUI esistente (`internal/tui/`, FR-10 Sprint 1) estesa con menu pre-mappato: `labnexus` senza args mostra la lista jobs auto-discovery con scelta interattiva (oltre al flusso esistente "profilo → input → output").
- **FR-32** — Logger esteso multi-writer: per ogni run, `runlog.Logger` scrive contemporaneamente su `os.Stderr` (visualizzazione live) e su `<output_dir>/<timestamp>_<profile>_<input_descriptor>.log` (audit trail file). Pattern `io.MultiWriter`.
- **FR-33** — Log dettagliato per step (espansione di Sprint 1 FR-8): file caricati con dim/byte/token estimate per file, kb_files caricati con path + lunghezza testo, prompt composto con char count + token estimate, provider chiamato (URL endpoint + modello + parametri), streaming throughput (tok/s), output finale (path + dimensione bytes).
- **FR-34** — Frontmatter MD output include nuovo campo `log_file: <path-relativo>` che punta al `.log` accoppiato. Tracciabilità tra bozza e audit trail.
- **FR-35** — Body streaming live in stderr durante chiamata provider. Oltre alla progress bar esistente (`runlog.StreamProgress` Fetta 1 UX patch), ogni chunk di testo ricevuto via `<-chan StreamEvent` viene scritto incrementalmente su stderr (in TTY) e sul log file (sempre). Denis vede il modello "pensare" in tempo reale. In non-TTY (pipe/redirect): solo log file riceve i chunk, stderr resta determinista per scripting.
- **FR-36** — Pre-check API key in `labnexus.command` (e in altri entry point): all'avvio, se `labnexus.config.toml` ha `eurouter_api_key = ""` o assente E provider effettivo è `eurouter`, **errore esplicito** con messaggio operativo italiano: *"Apri `labnexus.config.toml` e inserisci la tua chiave EUROUTER_API_KEY prima di lanciare."* Exit code 2 (NFR-5 invariato).
- **FR-37** — README zip-level riscritto in italiano, tono concierge: setup 3 passi (estrai zip → edit `labnexus.config.toml` con la API key → doppio click su `labnexus.command`), workflow Denis (TUI menu, esecuzione, lettura output in `lavori/<X>/output/`), guida per aggiungere nuovi lavori (duplica cartella + edit `_labnexus.toml`).
- **FR-38** — `docs/USER-GUIDE.md` aggiornato per concierge mode: workflow di lavoro reale Denis (no scenari sintetici), esempi su materiali tipo `data/CAPABILITY *` (rinominati `lavori/*`), sezione su come Denis legge gli output + valuta L2 + annota `valutazione_denis` nel frontmatter.
- **FR-39** — `docs/guida-meta-prompt-denis.md` aggiornato: "aggiungi nuovi tipi di lavoro" (non "aggiungi nuovi test"), schema TOML invece di YAML, riferimento a `trigger_prompt_file` come pattern per controllo runtime sul prompt.

## Non-Functional Requirements

Spec amendments rispetto a Sprint 1 NFR (la spec Sprint 1 è archiviata; questi amendments sono il delta):

- **NFR-1 amendment** — **Locale opzionale, cloud EU-GDPR default**. Post-pivot-3: il binario consegnato a Denis è configurato per girare con `provider: eurouter` di default. Locale `ollama` resta supportato come opzione via config master o profilo override. Cliente (Stefano Fiorina + Denis) informati e d'accordo; DPA da formalizzare lato contrattuale. **I dati reali del SGQ di Denis transiteranno a `api.eurouter.ai`** — questo è il nuovo statement contrattuale.
- **NFR-2 — Portability**: invariato. Due target Sprint 1.5: `darwin/arm64` (Denis) + `linux/amd64` (CTO + futuro server backend Sprint 2). Cross-compile Go.
- **NFR-3 — Performance/UX**: invariato. Streaming incrementale, primo token visibile ≤ 5s. Con eurouter (rete) la latenza dipende da provider; comunque < locale Ollama warmup (su Mac Denis era impraticabile).
- **NFR-4 — Determinismo CLI per piping**: invariato. stdout deterministico. **AMPLIATO**: body streaming live va in stderr (no stdout) per non sporcare lo stdout deterministico.
- **NFR-5 — Exit codes semantici**: invariato (0/1/2).
- **NFR-6 amendment — Sicurezza path traversal**: **esteso a `trigger_prompt_file`** (FR-16). Il file referenziato non può uscire dalla cartella `--input` (pattern analogo a `kb_files` Sprint 1 + bugfix #002 KB symlink escape Fetta 2).
- **NFR-7 amendment — Privacy**: i log per step **NON** stampano contenuto dei file di input (invariato). I log per step **STAMPANO** il body streaming live in stderr quando `--show-prompt` o run normale è TTY (FR-35). Questo è informed-consent del CTO (lui vede a video, sa cosa sta facendo). Pre-fix #005 PII warning su `--show-prompt` resta valido per dry-run pipato. **I file di input transitano a eurouter** come parte del prompt (NFR-1 amendment).
- **NFR-8 amendment — Quality metric / conformity criterion**: invariato strutturalmente (L1 strutturale CTO + L2 qualitativa Denis). **Nuovo significato di L2**: "uso reale + correzione manuale di Denis" (concierge mode), non sessione formale di validazione one-shot. Frontmatter `valutazione_denis` field invariato.
- **NFR-9 — Stabilità di scope**: invariato.
- **NFR-10 (NUOVO) — Niente hardcoding di nomi cartelle/profili nel codice o negli script `.command`**. Il sistema deve essere robusto a riorganizzazioni del filesystem `lavori/` da parte di Denis (rinomine, aggiunte, rimozioni) senza necessità di modifiche a config centrale, codice, o re-build. Solo `_labnexus.toml` per-cartella + convention naming fallback come fonti.
- **NFR-11 (NUOVO) — Audit trail**: per ogni esecuzione, il sistema genera SIA il file output MD (bozza per Denis) SIA un file `.log` accoppiato (audit trail con dettagli per step). Tracciabilità ISO 17025: chi (modello+provider), quando (timestamp), come (kb_files iniettati, trigger_prompt, file input), cosa (output prodotto, dimensione, token consumati). Frontmatter MD include `log_file` field per cross-reference.

## Input/Output

### Input

- **Cartella `lavori/`**: contiene N sottocartelle, una per ogni lavoro/job che Denis vuole eseguire. Ogni `lavori/<X>/`:
  - Contiene file domain-specific (input reali) in formati `.md, .txt, .csv, .docx, .pdf, .xlsx, .xls, .doc, .rtf`.
  - Contiene un file `_labnexus.toml` (2-3 righe) che dichiara quale profilo usare e opzionalmente il `trigger_prompt_file`.
  - Eventualmente sottocartelle di organizzazione (es. `Documento_da_revisionare/`, `Nuovi_Requisiti/`).
- **Config master `labnexus.config.toml`** al root: defaults globali (provider, modello, eurouter_api_key, ecc.). Editato da Denis dopo estrazione zip.
- **Profili `profili/*.toml`** (7 file): config per-capability con campi opzionali (ereditano da master).
- **KB-ispettore `KB-ispettore/`**: invariata da Sprint 1, knowledge base di Denis.

### Output

- **Per ogni esecuzione di `labnexus run`** (sia via `--job` sia via `--profile`):
  - File markdown `<output_dir>/<timestamp>_<profile>_<input_descriptor>.md` con frontmatter YAML completo (incluso `log_file`).
  - File log `<output_dir>/<timestamp>_<profile>_<input_descriptor>.log` (audit trail dettagliato per step + body streaming).
  - Default `output_dir` = `<input_dir>/output/` quando run via `--job` o `--input`. Override esplicito via `--output`.
- **stdout deterministico**: path del file MD generato (per scripting backend).
- **stderr**: log per step + body streaming live (TTY) + progress bar (TTY).

## Interface & Intuitiveness

### Main Instructions

| Screen/Command | Main Instruction |
|---|---|
| `labnexus.command` (doppio click Finder) | Denis fa doppio click; Terminal si apre con la TUI che lista i lavori disponibili. |
| `labnexus` (no args) | TUI interattiva: scegli un lavoro disponibile dal menu, esegui. |
| `labnexus jobs` | Lista jobs auto-discovery dalla cartella `lavori/`. |
| `labnexus jobs --json` | Output strutturato JSON per backend / scripting. |
| `labnexus run --job <name>` | Esegui in headless un job by name (per backend o CLI esperto). |
| `labnexus run --profile X --input Y --output Z` | Esegui in headless con path espliciti (Sprint 1 invariato, server-friendly). |
| `labnexus list` | Elenca i profili installati (Sprint 1 invariato). |
| `labnexus describe <profile>` | Mostra dettagli del profilo (Sprint 1 invariato). |
| `labnexus check <profile> --input <dir>` | Dry-run senza chiamata LLM (Sprint 1 invariato). |
| `labnexus validate <profile>` | Schema check del profilo (Sprint 1 invariato, ma valida MERGED profile post-master). |

### Intuitiveness Checklist (primary interaction = doppio click su `labnexus.command`)

- [x] **Discoverable**: `labnexus.command` al root dello zip, accanto al README. Denis apre zip → vede il file.
- [x] **Comprehensible**: README zip-level (italiano) spiega in 3 passi: estrai → API key → doppio click. La TUI mostra il menu con nomi capability auto-discovery (es. "Revisione documenti SGQ", "Gestione rilievi ACCREDIA"). Denis sa cosa sta per accadere.
- [x] **Forgiving**: API key mancante → errore esplicito italiano che dice cosa fare. EC-15 (Sprint 1) gestisce collisione output (suffisso `_2`, `_3`). Niente overwrite. Denis può ri-lanciare lo stesso lavoro senza paura di perdere output precedenti.
- [x] **Feedback**: streaming live in console (FR-35) — Denis vede il modello produrre output in tempo reale. Progress bar (FR-8 Sprint 1) mostra throughput. Log file accoppiato per audit post-fatto.

## Edge Cases

- **EC-1** — `labnexus.config.toml` ha `eurouter_api_key = ""` o assente, provider effettivo è eurouter → errore esplicito italiano in stderr, exit 2. Niente chiamata di rete.
- **EC-2** — Cartella `lavori/` vuota o assente → `labnexus jobs` mostra "Nessun lavoro trovato in `lavori/`"; TUI menu vuoto. Niente crash.
- **EC-3** — `_labnexus.toml` in una cartella `lavori/<X>/` referenzia profilo inesistente → `labnexus jobs` lista il job con warning "profilo X non trovato"; esecuzione (se tentata) fallisce con errore profile validation.
- **EC-4** — Cartella `lavori/<X>/` senza `_labnexus.toml` e nome non matcha convention regex → cartella skippata con warning in stderr; non appare nel menu jobs.
- **EC-5** — `trigger_prompt_file` referenzia file inesistente nella cartella `--input` → errore al run-time (path resolution failure). Profile validation static (`labnexus validate`) non può catturare (il file vive in `--input` dinamica).
- **EC-6** — `trigger_prompt_file` punta a path che esce da `--input` (es. `../../etc/passwd`) → errore path-traversal (NFR-6 amendment), exit 2.
- **EC-7** — File `.doc` corrotto / non parsabile (libreria Go meno matura) → warning specifico "file .doc non parsabile, conversione manuale a .docx consigliata", skip del singolo file. Pattern come EC-1 Sprint 1 per PDF rotti.
- **EC-8** — File `.xls` con macro / encrypted → warning + skip. Pattern come EC-1.
- **EC-9** — File `.rtf` con encoding non-UTF8 → tentare conversione, se fallisce: warning + skip.
- **EC-10** — Config master `labnexus.config.toml` malformato (TOML syntax error) → errore esplicito al primo `labnexus` invocation, con line:column del problema (lib `BurntSushi/toml` fornisce posizione). Exit 2.
- **EC-11** — Profile `profili/X.toml` malformato → analogo errore TOML con file:line.
- **EC-12** — Profile dichiara sia `trigger_prompt` sia `trigger_prompt_file` (XOR violato) → errore in `profile.Validate`, exit 2 al `labnexus validate` o al run.
- **EC-13** — Profile dichiara né `trigger_prompt` né `trigger_prompt_file` → errore analogo.
- **EC-14** — Eurouter API rate limit o errore network → errore esplicito stderr (Sprint 1 pattern invariato). Body parziale eventualmente già scritto al log file (audit trail).
- **EC-15** — Output collision (rilancio stesso lavoro nello stesso secondo) → suffisso `_2`, `_3` (Sprint 1 EC-15 invariato).
- **EC-16** — Denis estrae lo zip via Finder (browser download) → quarantine flag macOS. Al primo `./labnexus` Gatekeeper blocca. Workaround: `xattr -d com.apple.quarantine labnexus labnexus.command` (documentato in README). Oppure control-clic → "Apri" la prima volta. **Già nota da Sprint 1** per `.app`; ora si applica a binary nudo + `.command`.

## Data & Privacy

- **Personal data involved**: SÌ. I dati reali del SGQ di Denis in `lavori/<X>/` (NC, audit interni, rapporti di prova, matrici competenze con nomi tecnici, schede personale, ecc.) **contengono riferimenti a persone del laboratorio** (operatori, tecnici, fornitori, clienti).
- **Purpose**: generare bozze di lavoro QM (revisioni, gap report, audit checklist) come concierge preliminare. Denis le rivede + le usa nel suo lavoro di Quality Manager.
- **Retention**: gli input restano sul filesystem locale di Denis (cartella `lavori/`) fino a quando lui decide di rimuoverli. Gli output generati restano in `lavori/<X>/output/` fino a quando lui decide. Niente storage centralizzato da parte di LabNexus.
- **Deletion**: Denis cancella manualmente le cartelle / file. LabNexus non gestisce ciclo di vita dei dati.
- **Third parties**: **SÌ. Pivot 3**: il contenuto dei file di input + i contenuti dei kb_files transitano a `api.eurouter.ai` (gateway cloud EU-GDPR, OpenAI-compatible) come parte del prompt LLM. Modello processante: Qwen 3.6 36B MoE (o equivalente disponibile su eurouter). Eurouter è third-party processor.
  - **Cliente (Stefano Fiorina + Denis) informati e d'accordo**.
  - **DPA da formalizzare** (azione lato cliente, fuori scope tecnico Sprint 1.5).
  - Gli output prodotti dal modello NON vengono restituiti al cliente come "dati persistenti" su eurouter — eurouter non fa caching/storage by design (statement da verificare nel DPA).
- **Sicurezza credentials**: `eurouter_api_key` vive in `labnexus.config.toml` in chiaro sul filesystem Denis. Single-user, no rete locale, no condivisione. Accettabile per Sprint 1.5. **Sprint 2 considerazione**: keychain Mac / env var / file `.env` con permessi 0600. Out-of-scope per ora.
- **Log audit trail**: i `.log` file in `lavori/<X>/output/` contengono dettagli per step + body streaming. **NON contengono il prompt composto completo** (NFR-7 invariato — il prompt non è stampato nel log normale). **Contengono il body output del modello** (è ciò che Denis vede comunque nel `.md`). Plus log accoppia path file input, kb_files caricati, durate, throughput.

## Out of Scope

- **Watcher cartella + trigger automatici**: Sprint 2 MVP. Sprint 1.5 si ferma a invocazione manuale o backend headless una-tantum.
- **DPA formale eurouter**: azione contrattuale lato cliente (Stefano Fiorina). Fuori scope tecnico.
- **Pre-processing deterministico per `pt-analysis`** (OQ-8 Sprint 1): rimane per Sprint 2 hardening numerical reasoning.
- **Backup / sync output cloud**: Denis gestisce filesystem locale. Niente sync cloud automatico.
- **Multi-tenant / multi-user**: single-user Denis. Sprint 2 può considerare multi-tenancy se serve.
- **Auto-update binary**: Sprint 2 distribution / update mechanism.
- **Re-introduzione drag&drop iconcina `.app`**: se Denis lo richiede dopo prova del concierge mode, valutiamo per Sprint 2. Sprint 1.5 si impegna a "doppio click + TUI menu".
- **Refactor split `features/steps_test.go` ~2400 LoC**: tech debt accettato in Fetta 3 compound, rimane.
- **Re-implementazione `splitInThree` UTF-8 boundary-safe**: pre-existing bug noto, todo #009 P3 deferred. Resta fuori scope Sprint 1.5.
- **Keychain integration per eurouter_api_key**: Sprint 2.
- **Conversione automatica `.doc` → `.docx`** in caso di failure parser: out-of-scope. Denis può convertirlo manualmente se serve.
- **TUI Streaming display body live** in modo "fancy" (es. terminal full-screen rendering con header/footer fissi): Sprint 1.5 implementa append semplice in stderr, no fancy UI.

## Open Questions

(Tutte risolte nel brainstorm. Spec non ha open questions blockanti.)

- Eventuale ri-introduzione drag&drop iconcina: **decisione rimandata a feedback Denis post-handoff** (se richiesto, valutiamo Sprint 2).
- Modello specifico su eurouter (qwen3.6 confermato disponibile da CTO; eventuali fallback se eurouter cambia disponibilità modelli): **lasciato a config master + responsabilità CTO**.

---

**Brainstorm sorgente**: `.pipeline/brainstorms/2026-05-21-refactor-sprint1-concierge-mode.md`
**Idea assorbita**: `.pipeline/ideas/2026-05-21-frontmatter-key-ownership-table.md` (in-progress)
**Sprint 1 spec di riferimento (archiviata)**: `.pipeline/archive/2026-05-21-sprint1-labnexus-spec.md`
