---
id: 005
priority: P2
status: ready
source: review
created: 2026-05-20
file: cmd/labnexus/main.go (subcommand check)
bug_ref: .pipeline/bugs/show-prompt-pii-exposure.md
---

# `labnexus check --show-prompt` warning PII su non-TTY + nel README

## Issue

Il flag `--show-prompt` (FR-5, dry-run) stampa il prompt LLM completamente composto su stdout, inclusi tutti i kb_files (con tono ispettore, contenuti procedure) e tutti gli input (NC, audit, rapporti SGQ con riferimenti a persone). Se l'utente pipe l'output (`> debug.txt`) o lo incolla in un ticket di supporto, PII di colleghi finisce in canali non controllati. Violazione potenziale NFR-6/7 + GDPR data minimization.

## Fix (Sprint 1 — soft mitigation)

1. In `cmd/labnexus/main.go` (subcommand `check`), quando `--show-prompt` è settato AND `stdout` NON è una TTY:
   - Stampa su stderr (PRIMA del prompt) il warning:
     ```
     ⚠ --show-prompt output contiene contenuti dei kb_files e degli input.
       Se questi includono dati personali del SGQ, NON ridirezionare a file
       né condividere su canali non controllati. Continuare? [y/N]
     ```
   - In TTY: prompt interattivo; in non-TTY: rifiuta a meno che non sia settato `LABNEXUS_SHOW_PROMPT_CONFIRMED=1`.
2. Aggiornare `README.md` con avviso nella sezione comando `check`: *"⚠ `--show-prompt` espone il prompt completo (KB + input). Usa solo per debug, non condividere l'output."*
3. (Sprint 2) implementare flag `--redact` che maschera nomi propri, codici personali, email con regex semplice.

## Context

- Confermato da Mistral L1 + Codex L2.
- Mitigation soft è economica (~30 min) e protegge Denis dall'incidente comune "copio l'output di check in un ticket di supporto".
- Convertito in todo P2 da `/v-triage` 2026-05-20.
- Stima fix Sprint 1 soft: 1 gettone infra (~30 min CTO).
