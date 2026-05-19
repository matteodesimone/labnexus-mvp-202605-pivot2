# Review: bugfix `app-bundle-doppio-click-no-output`

## Summary

Bug HIGH risolto via fix idiomatic macOS (CFBundleExecutable wrapper bash + osascript per aprire Terminal + `labnexus-bin` come binario reale). Red-phase TDD rigorosa (5 test integration in `internal/bundle/`, build tag `darwin`). Tre finding del loop 1 risolti in loop 2 (escape AppleScript per path con `"`/`\`; scenario BDD stale aggiornato; comportamento CLI documentato in USER-GUIDE). Bonus emerso: 1 **unmasked failure pre-esistente** (scenario BDD "Ollama non disponibile" flaky a 87s) fixato inline come da procedura framework (≤30 min). Tutti i 11 package Go GREEN, 6/6 nuovi test del bundle, suite acceptance BDD 40 scenari attivi + 3 pending + 13 manual rispettati.

## Review Method

- **Primary**: Claude (full project context — bugfix cycle)
- **External reviewers**: NESSUNO disponibile (`.pipeline/bin/vobbly` assente). **Single-reviewer mode, external reviewers unavailable.** Stessa situazione di Fetta 1 — installazione di vobbly resta TODO operativo del CTO.
- **Pre-flight**: non eseguito (binario assente).
- **Mode**: `REVIEW="all"` con `REVIEW_MAX_LOOPS=3`. Loop corrente: 2/3 (PASS al loop 2, niente loop 3).
- **Standards**: Layer 0 (8 principi) + Layer 1 (TDD/BDD, mock zones) + Layer 2 (`backend-cli.md`) + Layer 3 (`02-project.md` + 2 skill di progetto emerse da compound Fetta 1) + `checks.yaml` (5 check default).
- **Scope**: solo i file del bugfix cycle (commit 772cfc5 + fix loop 2 — vedi sotto).

---

## Findings — stato risoluzione (loop 1 → loop 2)

### MEDIUM (1/1 risolto)

#### ✅ M1 — Wrapper non escapa `"` e `\` per AppleScript string literal
- **File**: `scripts/build-mac.sh` (wrapper bash heredoc)
- **Issue**: il loop 1 escapava solo `'` (shell single-quoting). Un path tipo `Test "1"` (con virgolette doppie, raro ma legale su macOS) avrebbe rotto la stringa `do script "..."` di AppleScript.
- **Fix applicato**: nuova funzione `escape_for_osascript()` che applica due livelli di escape: shell single-quote per `'` → `'\''` E AppleScript string literal per `\` → `\\` e `"` → `\"`. Helper `build_command()` riusa l'escape sia per il path del binario sia per l'argomento posizionale.
- **Test red phase**: nuovo `TestBundle_WrapperEscapesApplescriptStringChars` verifica via grep dei pattern sed `s/"/\\"/g` e `s/\\/\\\\/g` nel wrapper. Falliva prima del fix, passa adesso.

### LOW (2/2 risolti)

#### ✅ L1 — Scenario BDD `bundle .app contiene l'eseguibile arm64` stale
- **File**: `features/motore-modo-interattivo.feature`
- **Issue**: lo scenario `@manual` affermava "labnexus eseguibile arm64" — ma dopo il fix `labnexus` è il **wrapper bash**, non il binario. Documentazione disallineata col codice (Layer 0 "doc è codice").
- **Fix applicato**: scenario rinominato in `bundle .app contiene wrapper bash + binario reale e Info.plist valido`. Step aggiornati per riflettere `labnexus` come wrapper di testo + `labnexus-bin` come Mach-O arm64 + presenza di `osascript`. Aggiunta nota cross-reference a `internal/bundle/bundle_test.go` come verifica automatizzata (lo scenario resta `@manual` per validazione end-to-end Finder+Terminal reale, non automatizzabile in CI).

#### ✅ L2 — Wrapper invocato da CLI apre Terminal in finestra separata (non documentato)
- **File**: `docs/USER-GUIDE.md`
- **Issue**: se l'utente lancia `./labnexus.app/Contents/MacOS/labnexus` da una shell, il wrapper apre una **seconda** finestra Terminal. Comportamento corretto per il doppio click dal Finder, ma sorprendente da CLI.
- **Fix applicato**: aggiunta nota tecnica in `docs/USER-GUIDE.md` sezione "Riga di comando" che dice: usa `labnexus-bin` (binario reale) da terminale; il `labnexus` wrapper è destinato al Finder.

---

## Check-by-check (checks.yaml domain)

### quality
- ✓ Wrapper bash conciso (~40 righe), funzioni `escape_for_osascript` + `build_command` ben separate (Layer 1 modularity balanced).
- ✓ Naming chiaro (`bin_safe`, `input_safe`, `CMD`).
- ✓ Niente over-engineering: il wrapper resta < 50 righe.

### test-coverage
- ✓ 6 test integration in `internal/bundle/` (era 5 al loop 1, +1 per l'escape AppleScript).
- ✓ `TestMain` esegue `scripts/build-mac.sh` una sola volta, tutti i 6 test fanno assertion sul filesystem prodotto.
- ✓ Build tag `darwin` esclude su Linux CI (evitiamo race con cross-compile).
- ⚠ Limit: il test "doppio click reale dal Finder" resta `@manual` — non automatizzabile senza simulazione Launch Services. Test integration ne coprono la struttura ma non l'end-to-end vero.

### security
- ✓ `escape_for_osascript` previene injection via path con caratteri speciali (`'`, `"`, `\`).
- ✓ Path traversal: il wrapper risolve `BIN="${DIR}/labnexus-bin"` dove `DIR` è il path canonical del bundle stesso — niente arbitrary code execution.
- ✓ Nessun secret nel wrapper, niente env dump.
- ✓ HTTP timeout esplicito in `OllamaProvider.doRequest` (90s) — previene DoS via slow-loris se un attaccante punta l'endpoint a un server lento.

### intuitiveness (CLI/GUI)
- ✓ **Discoverable**: doppio click ora produce visibilmente una finestra Terminal con la TUI (era invisibile prima).
- ✓ **Feedback**: l'utente vede subito la TUI in azione.
- ✓ **Forgiving**: drag&drop di una cartella pre-seleziona l'input (path con caratteri speciali ora gestito correttamente).
- ✓ **Predictable**: il binario reale `labnexus-bin` è invocabile da CLI con gli stessi flag CLI di prima.

### privacy
- ✓ Nessun PII nel wrapper o nei log. Path utente passato come argomento ma non loggato esplicitamente nel processo wrapper.
- ✓ `quoted form of` non applica encoding ma preserve l'integrità del path.

---

## Coherence con il piano Sprint 1

| Aspetto | Stato |
|---|---|
| FR-10 (TUI cross-platform) | ✓ Ora effettivamente raggiungibile via doppio click su macOS (era teorica) |
| FR-11 (bundle .app + binario Linux) | ✓ Bundle ora ha layout idiomatic; binario Linux invariato |
| NFR-2 (portability) | ✓ macOS arm64 + linux/amd64 invariati |
| NFR-1 (locale mandatory) | ✓ Wrapper non altera il flusso provider (ollama/eurouter) |
| NFR-5 (exit codes) | ✓ Wrapper propaga exit del binario reale via osascript (default behavior) |
| L1+L2 pattern (NFR-8) | ✓ L1 strutturale coperto dal bundle test; L2 (validazione Denis al doppio click reale) resta TODO operativo CTO |

## Unmasked failures — gestione

1 unmasked failure fixato inline come da regola framework (≤30 min):
- **BDD `Ollama non disponibile`**: era diventato flaky (87s + exit 0 invece di 1). Causa: macOS port 1 (tcpmux) accetta TCP, `OllamaProvider` senza HTTP timeout. Fix: `http.Client.Timeout=90s` esplicito + endpoint RFC 5737 (`192.0.2.1`) per fail-fast deterministico. Tempo: ~15 min, dentro budget.

Nessun altro unmasked failure rilevato.

## Metrics

- Files reviewed: 6 (`scripts/build-mac.sh`, `internal/bundle/bundle_test.go`, `internal/provider/ollama.go`, `features/steps_test.go`, `.pipeline/bugs/app-bundle-doppio-click-no-output.md`, +`features/motore-modo-interattivo.feature` e `docs/USER-GUIDE.md` aggiunti in loop 2)
- Findings loop 1: **3** (0 critical, 0 high, 1 medium, 2 low)
- Findings loop 2: **0** (tutti risolti)
- Confirmed by both: n/a (single-reviewer)
- Claude-only: 3 → 3 risolti
- External-only validated: 0
- False positives dismissed: 0

## Verdict

**PASS** (loop 2/3, `REVIEW="all"`).

Tutte le findings risolte, suite GREEN, unmasked failure pre-esistente fixato inline, documentazione allineata, BDD scenario aggiornato. Pronto per `/v-deploy` (rebuild dello zip di consegna con il bundle corretto).

---

## Action items operativi residui (non bloccanti)

1. Installare `.pipeline/bin/vobbly` per riabilitare review multi-reviewer (carry-over Fetta 1)
2. **Test `@manual` "doppio click reale dal Finder"** prima della consegna a Denis — vincolo nuovo emerso da questo bug: aggiungere checkpoint CTO esplicito pre-`/v-deploy` per i scenari `@manual` (vedi compound Fetta 1 → idee `.pipeline/ideas/2026-05-19-pseudo-tty-bdd-testing.md` per renderli automatizzabili in futuro)
3. Re-fare lo zip di consegna (`make ship`) per produrre il nuovo `dist/labnexus-sprint1-darwin-arm64.zip` con il bundle fixato — è il deliverable corretto per Denis
