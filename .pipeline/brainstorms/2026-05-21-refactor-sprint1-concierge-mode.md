# Brainstorm: Refactor Sprint 1 — pivot 3 + concierge mode pre-handoff Denis

Date: 2026-05-21

## The Problem

Lo zip Sprint 1 deliverable (Fetta 1+2+3 chiuso + 5 bugfix security) è tecnicamente completo ma **non utilizzabile sul Mac di Denis** così com'è. Tre fattori critici emersi dopo la chiusura formale Sprint 1:

1. **Hardware Denis insufficiente per Ollama+Qwen3** (~22-28 GB RAM). Il deliverable lo presuppone come default (NFR-1 "locale mandatory"), ma non è eseguibile in pratica. Il cliente è informato e d'accordo a pivottare su eurouter cloud EU-GDPR (paga lui l'API).
2. **`labnexus.app` bundle macOS** introduce verbosità inutile (wrapper bash + osascript + Info.plist + LSHandlerRank) — Matteo CTO vuole un CLI puro multipiattaforma macOS+Linux per coerenza e semplicità.
3. **La cartella `data/`** appena consegnata da Denis contiene **scenari di utilizzo reale**, non test fixture. Denis userà LabNexus come **concierge preliminare** per il suo lavoro quotidiano di Quality Manager + ispettore ACCREDIA: bozze di revisione, gap report, audit checklist che lui consumerà sul serio. In ottica Sprint 2: MVP più completo e automatizzato.

Il refactor mira a chiudere questi 3 gap come **Sprint 1.5 pre-handoff**, a 0 gettoni cliente (offerto dal CTO).

## Key Insights from Discussion

- **Pivot 3 è confermato dal cliente, non un'idea CTO unilaterale**: Stefano Fiorina + Denis informati e d'accordo. Pagano l'API eurouter. Sprint 1 NFR-1 va archiviato come "spec amendment post-pivot-3".
- **TOML per config user-edited è oggettivamente più safe di YAML**: indent-based fail silenzioso è il principale generatore di errori in config DevOps. Denis non è programmatore. TOML è "pensato per human-edited config" (claim del creator T. Preston-Werner). Frontmatter MD output resta YAML (convenzione Obsidian/Jekyll/Hugo universal, non editato da Denis).
- **Il `Prompt_INPUT_Rev.00.rtf` di Denis È il vero trigger_prompt**, non un file di input qualsiasi. Lui ha fatto il prompt engineering al posto del CTO. Architettura naturale: `trigger_prompt_file` nel profilo come indirezione (XOR con `trigger_prompt` inline).
- **Auto-discovery + flag CLI bypass coesistono pulitamente**: la TUI fa auto-discovery del `lavori/` per il menu Denis; il backend server passa `--job <name>` o `--profile X --input Y --output Z` per headless. Niente lock-in al test-runner mode.
- **Cartella `lavori/` (italiano per Denis) + subcommand+flag inglesi `jobs/--job` (CTO + backend)**: compromesso pulito. Denis vede italiano nei suoi materiali; il CTO vede inglese tecnico nei comandi.
- **NON HARDCODING dei nomi cartelle**: Denis aggiungerà/rinominerà/rimuoverà cartelle. Soluzione: `_labnexus.toml` 2-righe dentro ogni cartella (profile + opzionale trigger_prompt_file). Lo script `.command` lavora relativamente al proprio path.
- **Output dentro ogni cartella `lavori/<X>/output/`**: niente parsing "lettera + profilo" perché il path è risolto dalla cartella stessa.
- **Logger come audit trail (non solo debug)**: ogni run produce `.md` (bozza) + `.log` (audit trail accoppiato). Frontmatter `log_file` field nel MD punta al log. Coerente con ISO 17025.
- **Streaming body live in stderr + log file**: i provider EUrouter SSE + Ollama NDJSON già emettono chunk via channel. Estensione naturale: visualizzare i chunk in console mentre il modello pensa. UX significativa (ansia ridotta + debug visibile).
- **Mental model "concierge" vs "test runner"**: cambia naming + tono documenti, non architettura.

## Possible Approaches

### Approach A: Big-bang refactor (un unico cycle Sprint 1.5)
- Tutti i punti (1a+b+c+d + 2a-k) in un singolo cycle spec→plan→test-scaffold→implement→review→deploy.
- **Pro**: zip finale coerente, una sola validazione, una sola consegna.
- **Con**: scope grande (~23-36h CTO), un solo punto di fallimento, review massiccia.
- **Complessità**: complex (>8h).

### Approach B: Refactor in 3 sub-fette ordinate
- Sub-fetta 1.5.A: CLI puro + eurouter default + rimozione gate (Pivot 3 core) — ~5-7h.
- Sub-fetta 1.5.B: Config master TOML + profile-override + parser nuovi formati — ~9-13h.
- Sub-fetta 1.5.C: Concierge mode (auto-discovery `lavori/` + `--job` + log audit + streaming live + UI/docs) — ~9-16h.
- **Pro**: rischio diluito, ciascuna sub-fetta validabile + shippable indipendentemente.
- **Con**: 3 cycle di review/deploy = overhead pipeline; potenziale conflict di state intermedi.
- **Complessità**: medium per sub-fetta.

### Approach C: Refactor + L2 Denis paralleli
- Sub-fetta 1.5.A (Pivot 3) ship subito → Denis può cominciare i test reali con quella versione.
- Sub-fetta 1.5.B+C in parallelo a L2 Denis: il CTO continua a sviluppare mentre Denis valida le 5 capability nuove + FR-22 sui 3 casi Claude.
- **Pro**: massima velocità per Denis (riceve subito quello che gli serve per cominciare).
- **Con**: Denis lavora con un binary intermedio (no log audit trail, no streaming live, no auto-discovery) — esperienza degradata nelle prime esecuzioni.
- **Complessità**: medium ma coordinazione complicata.

## Recommended Approach

**Approach B — 3 sub-fette ordinate**.

Razionale (sottrarre è moltiplicare):
- A è solido ma rischioso: ~30h in un singolo cycle senza checkpoint intermedi = se qualcosa va male a metà, tutto bloccato.
- C è seducente ma frammenta troppo: Denis riceve un binary intermedio scomodo, perde la UX concierge (streaming + log audit) proprio quando ne avrebbe più bisogno.
- B mantiene la disciplina pipeline (review per ciascuna sub-fetta) + permette di shippare incrementalmente se serve. Le 3 sub-fette sono naturalmente disaccoppiate:
  - 1.5.A è prerequisito per qualunque uso reale (eurouter default).
  - 1.5.B abilita 1.5.C (config master + nuovi formati senza cui Denis non può manco aprire i suoi .xls).
  - 1.5.C è la finalizzazione concierge (UX + auto-discovery + audit trail).

**Sequenza implementativa**:

1. **Sub-fetta 1.5.A — Pivot 3 + CLI puro** (~5-7h):
   - Eliminare `.app` bundle: `scripts/build-mac.sh` semplificato a `go build cross-compile darwin/arm64` standalone binary.
   - `labnexus.command` minimal (2 righe bash: `cd "$(dirname "$0")"; ./labnexus`).
   - Eliminare 7+ BDD scenari `@manual` (Info.plist, drag&drop, bundle structure).
   - Rimuovere `ErrEurouterGateMissing` + `requireLoopbackOrCloudGate` da `internal/provider/provider.go`.
   - Rimuovere 4+6 test in `select_test.go`.
   - Aggiornare `helpers_test.go` (no più auto-set `LABNEXUS_ALLOW_CLOUD_PROVIDER`).
   - Bug `.pipeline/bugs/privacy-eurouter-gate-mancante.md`: status → "intentional deviation post-pivot-3".
   - I 7 profili: cambiare `provider: ollama` → `provider: eurouter` (transitorio, sparirà in 1.5.B con il config master).
   - Documenti: README zip-level (italiano Denis), eliminare riferimenti a `.app` / Gatekeeper bundle.

2. **Sub-fetta 1.5.B — Config master TOML + profile-override + parser nuovi formati** (~9-13h):
   - Nuovo pacchetto `internal/config/`: load `labnexus.config.toml` master.
   - Migration YAML → TOML per `profili/*.toml` (7 file) + nuovo `labnexus.config.toml` shipped.
   - `profile.Profile`: campi `Provider`, `Modello`, `Temperature`, `MaxTokens`, `ContextWindow` diventano opzionali (zero-value detection).
   - `runner.Run`: load master prima, merge profile-over-master, esegue.
   - `profile.Validate`: validazione su merged profile.
   - `trigger_prompt_file` field nel profile schema (mutually exclusive con `trigger_prompt`, path-safe relativo a `--input`).
   - Parser nuovi formati in `internal/input/`:
     - `.xls` via `shakinm/xlsReader` o `extrame/xls`.
     - `.rtf` via parser custom semplice o `scelto/go-rtf-parser`.
     - `.doc` via `richardlehane/mscfb` + parsing manuale (Old Word binary, lib Go meno matura ma fattibile).
   - Update meta-prompt-genera-profilo.md: schema in TOML + 3 esempi few-shot ricalibrati.
   - Test suite update (BDD helpers usano stub TOML, unit tests profile/config nuovi).

3. **Sub-fetta 1.5.C — Concierge mode + auto-discovery + log audit + streaming live + docs** (~9-16h):
   - Cartella `lavori/` (rinominata da `data/`) shipped nello zip.
   - `_labnexus.toml` dentro ogni cartella `lavori/<X>/` (2 righe: `profile`, opzionale `trigger_prompt_file`).
   - Nuovo subcommand `labnexus jobs`: auto-discovery `lavori/*/`, lista (eventuale `--json` per backend).
   - Nuovo flag `labnexus run --job <name>`: risolve via auto-discovery, esegue headless.
   - TUI esteso (`charmbracelet/huh`): menu pre-mappato dai job auto-discovery.
   - Output dentro `lavori/<X>/output/` (niente parsing nome cartella).
   - Logger esteso multi-writer: stderr + file `.log` per ogni run.
   - Frontmatter `log_file` field nel `.md` di output (audit trail).
   - Streaming body live in stderr (oltre a progress bar esistente); log mirror.
   - Log dettagliato per step: file caricati con dim/token, kb_files con path/lunghezza, prompt composto con N token, provider+endpoint+modello, streaming throughput tok/s, output path + dimensione.
   - Documenti aggiornati con tono concierge: README zip-level italiano, USER-GUIDE workflow reale, guida-meta-prompt-denis.md.

**Spec amendments necessari** (la spec Sprint 1 è archiviata; vanno tracciati come "Sprint 1.5 amendments" in un nuovo spec):
- FR-2: aggiunti `.xls`, `.doc`, `.rtf` ai formati supportati.
- FR-3: `trigger_prompt_file` opzionale (XOR con `trigger_prompt`).
- FR-10/FR-11: macOS `.app` superato → `.command` minimal launcher.
- NFR-1: locale opt-in, cloud EU-GDPR default (post-pivot-3, cliente informato).
- NFR-6: path traversal protection estesa a `trigger_prompt_file`.
- NFR-7: i dati reali transitano a eurouter (cliente d'accordo, DPA da formalizzare).
- NFR-8: "L2 Denis" diventa "uso reale + correzione manuale" (concierge mode).
- NUOVO FR: `labnexus jobs` + `labnexus run --job <name>` (auto-discovery `lavori/`).
- NUOVO FR: log file per ogni run (audit trail ISO 17025 compatible).
- NUOVO FR: body streaming live in stderr + log mirror.
- NUOVO FR: config master TOML con profile-override.

## Stress Test Results

5 obiezioni grillate, tutte sopravvissute (vedi flusso completo conversazione):

- **Assumption tested**: "Hardware Denis insufficiente per Ollama" → **HELD** (confermato sempre da Matteo, non assunzione).
- **Assumption tested**: "Cliente informato e d'accordo al pivot 3" → **HELD** (Stefano Fiorina d'accordo, cliente paga API eurouter).
- **Assumption tested**: "Su eurouter c'è qwen3.6 quindi golden file Test 1 transferisce" → **HELD** (Matteo conferma, e in ogni caso la libertà di config gli permette di calibrare).
- **Assumption tested**: "Gate `LABNEXUS_ALLOW_CLOUD_PROVIDER` può essere rimosso senza creare regressioni" → **HELD** (post-pivot-3 il pattern di protezione cambia natura, la decisione è solo via configurazione).
- **Assumption tested**: "Drag&drop iconcina perso è regressione UX accettabile" → **HELD** (FR-11 era UX bonus, non requirement contrattuale; workflow concierge via TUI menu compensa).
- **Assumption tested**: "Hardcoding nomi cartelle vincola Denis" → **REVISED** (proposta `_labnexus.toml` per-cartella + auto-discovery).
- **Assumption tested**: "Auto-discovery TUI-only limita backend server use" → **REVISED** (chiarito che `--profile X --input Y --output Z` resta bypass headless, e `--job <name>` è alternativa pulita).
- **Assumption tested**: "Cartella `data/` è test fixture" → **REVISED** (è scenario reale di lavoro Denis → rinominata `lavori/` con tono concierge).

## Ready for Spec

**Yes**.

Key inputs per `/v-spec`:

- **Core requirement**: refactor Sprint 1 pre-handoff Denis in 3 sub-fette ordinate (Pivot 3 + CLI puro → Config master TOML + nuovi parser → Concierge mode con auto-discovery + audit log + streaming live), per chiudere il gap tra il deliverable Sprint 1 chiuso lato AI e l'uso reale di Denis come QM laboratorio italiano.
- **Primary actor**: Denis Brazzo (QM laboratorio ISO 17025 + ispettore ACCREDIA), uso quotidiano come concierge preliminare per bozze di revisione, gap report, audit checklist da consumare sul serio. Secondari: Stefano Fiorina (cliente, paga API), Matteo De Simone (CTO sviluppo + futuro server backend Sprint 2).
- **Key constraints**:
  - 0 gettoni cliente, ~23-36h CTO (offerto da Matteo).
  - 3 sub-fette sequenziali (1.5.A → 1.5.B → 1.5.C), ciascuna con cycle pipeline completo.
  - Coerenza con Sprint 2 vision (MVP automatizzato, server backend, watcher).
  - Niente hardcoding di nomi cartelle / profili nel codice o negli script.
  - Italiano per Denis-facing (cartella `lavori/`, doc README/USER-GUIDE/guida-meta-prompt); inglese per CLI/CTO/backend (`labnexus jobs`, `--job`).
- **Out of scope** (rimandati a Sprint 2 / future fette):
  - Watcher cartella + trigger automatici (Sprint 2 MVP).
  - DPA formale eurouter (compito Stefano Fiorina lato contrattuale, non tecnico).
  - Pre-processing deterministico per `pt-analysis` (Sprint 2 hardening numerical).
  - Backup / sync output cloud (out-of-scope, Denis gestisce filesystem locale).
  - Multi-tenant / multi-user (single-user Denis).
  - Auto-update binary (Sprint 2 distribution).
  - Re-introduzione `.app` bundle drag&drop (se Denis lo richiede dopo concierge mode, valutiamo Sprint 2).
  - Refactor split `features/steps_test.go` ~2400 LoC (tech debt accepted in Fetta 3 compound).

## Related Items in Backlog

- `.pipeline/todos/007-p2-ready-privacy-paths-in-log.md`: redactPath helper per log paths personali. **Naturalmente integrato** in 1.5.C (logger esteso) — il pattern di redaction sul log audit trail fa parte del refactor.
- `.pipeline/todos/008-p2-ready-nodonemarker-design-review.md`: design review semantica EOF post-content. Indipendente dal refactor; può essere chiuso prima o dopo Sprint 1.5.
- `.pipeline/todos/009-p3-deferred-bdd-real-profile-content-gap.md`: BDD non e2e real content. Sprint 1.5.B (parser nuovi formati + config master) **non risolve** ma riduce parte del gap. Resta P3 deferred per Sprint 2.
- `.pipeline/ideas/2026-05-21-frontmatter-key-ownership-table.md`: tabella ownership chiavi frontmatter. Naturalmente integrabile in 1.5.B (config master schema).
- `.pipeline/proposed-updates/`: 3 framework promotions Fetta 2+3 (testmain-mrun-trap, unicode-word-boundary, meta-prompt-deliverable-pattern) — invariati, da triagiare separatamente.

---

**Brainstorm captured. When ready, run `/v-spec refactor sprint 1 pre-handoff Denis (pivot 3 + concierge mode in 3 sub-fette)` per formalizzare lo spec dello Sprint 1.5.**
