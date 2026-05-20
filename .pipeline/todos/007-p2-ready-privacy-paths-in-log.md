---
id: 007
priority: P2
status: ready
source: bug
created: 2026-05-20
file: internal/runlog/runlog.go, internal/runner/runner.go
bug_ref: .pipeline/bugs/privacy-paths-personali-in-log.md
---

# Paths personali stampati verbatim nei log step (PII leak)

## Issue

`internal/runner` log lines stampano verbatim path assoluti come `/Users/denis/Documents/SGQ-2025-Q4/laboratorio-Mario-Rossi/...`. Questi path possono contenere:
- Nome utente (Denis) — già noto al cliente, ma comunque PII
- Nome del laboratorio / cliente
- Eventuali codici personali nei nomi cartella

Quando Denis copia/pasta l'output del terminale in un ticket di supporto o in Slack, queste PII finiscono in canali non controllati.

## Fix

In `internal/runlog/runlog.go` o nel point-of-log in `internal/runner`:
1. Funzione `redactPath(p string) string` che ritorna solo il `filepath.Base(p)` (basename) OPPURE un path relativo alla cartella output scelta dall'utente.
2. Applicare a tutti i log step che menzionano paths.
3. Aggiungere flag `--verbose` che mantiene path completi per debug interno CTO.

Test:
- Unit: `redactPath("/Users/denis/SGQ/file.csv")` → `"file.csv"` (default).
- BDD: scenario lancia con `--verbose`: stderr contiene full path. Senza: solo basename.

## Context

- Confermato da Codex L2 review Fetta 2 + già filato in backlog Fetta 1.
- Mitigation soft (basename only) è low-risk per il debug CTO (il CTO ha già accesso al filesystem); è solo cosmetica per Denis.
- Convertito in todo P2 da `/v-triage` 2026-05-20.
- Stima fix: 1-2 gettoni infra (~30min-1h CTO).
