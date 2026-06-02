# Review: nuovi lavori self-service (wizard) + prompt nei profili + de-dup trigger/input

## Summary

La feature è solida e coerente col modello SYSTEM(KB)/TRIGGER(istruzione)/INPUT(dati).
De-dup, wizard `init`, serializzatore `WriteMetadata`, migrazione prompt e cartella
template sono implementati con test e verificati (full suite + `go vet` verdi; de-dup
provato e2e con `check --show-prompt`). I findings esterni in-scope sono stati risolti
inline; la grande maggioranza dei findings esterni riguarda **codice pre-esistente non
toccato** dalla feature (già tracciato come bug) oppure è **falso positivo** (fix di
stamattina che i reviewer non hanno recepito, o deviazioni intenzionali documentate).

## Review Method

- Primary: Claude (full project context)
- External reviewers: **Codex** (quality, security), **Mistral/vibe** (test-coverage, privacy)
- Pre-flight: 1° tentativo abort (mistral timeout transitorio); 2° tentativo **2/3 ready**. **Ollama down** (localhost non in esecuzione, tutti i check falliti).
- Modalità: **partial multi-review** (exit 1) — Ollama non disponibile
- Standards: `.pipeline/standards/` (layered) + `02-project.md` Layer 3

---

## In-scope — FIXED in questo ciclo

### [MEDIUM] De-dup: path di esclusione non normalizzato (Claude)
- **File**: internal/input/input.go (walkFilesDeterministic)
- **Issue**: l'esclusione usava `ToSlash` ma non `filepath.Clean`; un `trigger_prompt_file` non canonico (`./Prompt.txt`) non combaciava col rel-path pulito del walk → de-dup saltata silenziosamente (doppio invio).
- **Fix**: `filepath.ToSlash(filepath.Clean(e))`. Test `TestParseDir_ExcludeNormalizesPath` (red→green).

### [LOW] Messaggio non-TTY fuorviante per `init` (Claude)
- **File**: cmd/labnexus/main.go (initRunE)
- **Fix**: messaggio init-specifico ("richiede un terminale interattivo: Terminal o doppio click su Esegui.command") invece del generico `ErrNonTTY` condiviso.

### [LOW] Nessun test per `RunInit` (Mistral test-coverage)
- **File**: internal/tui/init_test.go
- **Fix**: `TestRunInit_NonTTYReturnsErrNonTTY` (guard d'ingresso in non-TTY, entrambi i rami). La parte interattiva resta manual-only (richiede TTY reale).

---

## Dismissed False Positives

- **Codex [HIGH] main.go:1 descrive Ollama-first** — già corretto stamattina (`Qwen 3 via EUrouter default Sprint 1.5+`). Verificato `cmd/labnexus/main.go:1-3`.
- **Codex [HIGH] guida referenzia bundle `.app`** — già corretto stamattina (righe 9, 74 → binary standalone + `labnexus.command`). Verificato.
- **Mistral [HIGH] manca test XOR `LoadMetadata`** — il test esiste (`trigger_inline_test.go::inline_e_file_insieme_rifiutati_XOR`).
- **Codex [HIGH] root lancia TUI** — design intenzionale Sprint 1.5 (doppio click `.command` → TUI per Denis; `02-project.md:35`). Layer 3 sovrascrive l'overlay cli-tool.
- **Mistral-privacy [CRITICAL] EUrouter senza runtime gate** — deviazione **intenzionale** post-pivot-3, documentata (`privacy-eurouter-gate-mancante.md`, cliente ha approvato). Il reviewer stesso cita il bug file.

## Deferred — pre-esistenti (non introdotti né smascherati dalla feature)

Già tracciati come bug nei cicli precedenti:
- API-key precheck ignora override provider (Codex) → `api-key-precheck-ignora-provider-override.md`
- `Resolve` job-name ambiguo (Codex) → `resolve-job-ambiguo-substring-match.md` — **nota**: il flusso template lo rende più rilevante (copie non rinominate collidono); mitigato dal `LEGGIMI.txt` che insiste sul rinominare.
- Input symlink exfiltration (Codex security, Mistral) → `input-symlink-escape-exfiltration.md`
- Output/log world-readable (Codex security) → `output-log-permessi-world-readable.md`

Backlog GDPR architetturale (Mistral-privacy) — **non** introdotto da questa feature, per triage CTO:
- PII via `--show-prompt` → già `show-prompt-pii-exposure.md`
- Nomi file nei log → già `privacy-paths-personali-in-log.md`
- **Tema nuovo**: data-lifecycle (retention/erasure/consent) — nessuna policy di ritenzione/cancellazione né tracciamento consenso. Item compliance ampio: candidato a spec dedicata, non a fix di feature. Non filato come bug singolo per non duplicare il lavoro di standardizzazione privacy già in corso.

Debito di test pre-esistente (Mistral): parser `.rtf/.xls/.doc`, symlink walk, `drainStream` timeout, `classifyError`, `stripFenceWrapping` — su codice non toccato dalla feature; annotato, non incluso in questo ciclo.

---

## Metrics
- Files reviewed: 10 (5 codice, 2 test, 1 launcher template, 1 profilo, 1 doc)
- Findings totali (esterni + Claude): ~26 (3 critical, 13 high, 8 medium, ~2 low)
- In-scope risolti: 3 (1 medium, 2 low)
- False positive dismessi: 5
- Deferred (pre-esistenti tracciati): 4 + cluster GDPR
- Confermati 2+ reviewer: input symlink (Codex+Mistral)

## Verdict
**PASS** — i findings introdotti da questa feature sono risolti; full suite + `go vet` verdi; de-dup verificato e2e. I findings residui sono pre-esistenti/architetturali, tracciati o annotati per triage separato.
