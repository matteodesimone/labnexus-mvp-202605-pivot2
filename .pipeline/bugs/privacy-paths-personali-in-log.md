---
severity: medium
status: open
created: 2026-05-19
source: review esterna (codex-privacy + mistral-privacy)
fix: ""
test: ""
---

# Bug: file paths personali stampati verbatim su stderr (PII in log)

## Reported

- **Feature/Area**: `internal/runner/runner.go:parseInputDir` (warning skipped file) + `runner.go:230` (output path) + `runner.go:printPrompt` (--show-prompt)
- **Trovato da**: review multi-LLM 2026-05-19 (codex-privacy CRITICAL + mistral-privacy MEDIUM)
- **Sprint impact**: MEDIUM. PII (`/Users/<nome>`, nomi file con persone reali) in log stderr.

## Description

Log emessi:
- `warning: .DS_Store: formato non supportato` — innocuo
- `warning: DE0779_RT_08rev03.pdf: ...` — può contenere riferimenti a persone
- `output: /Users/matteodesimone/.../2026-05-19_..._Denis_Brazzo_NC.md` — path completo con identificativi

Logs accidentalmente condivisi (ticket support, screenshot) espongono PII.

## Fix proposto

1. Warning skipped file: usare `filepath.Base(name)` invece di full path
2. Output path: rendere relativo (es. `output/<filename>.md`) — già relative se OutputDir fornito relativo, ma il default è absoluto
3. `--show-prompt`: aggiungere warning stderr "this prints input content verbatim, may contain PII"

## Risk of fix

**Basso**. Modifiche localizzate, fix in ~30 min + test.

## Note

Deferred to backlog: pre-esistente.
