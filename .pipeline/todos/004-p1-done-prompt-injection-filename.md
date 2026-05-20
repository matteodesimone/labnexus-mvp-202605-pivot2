---
id: 004
priority: P1
status: done
source: review
created: 2026-05-20
file: internal/prompt/prompt.go
bug_ref: .pipeline/bugs/prompt-injection-via-filename.md
---

# Sanitize filenames nel prompt composer (prompt injection prevention)

## Issue

`internal/prompt/composeUserMessage` interpola direttamente il nome di ogni file di input come marker `--- FILE: <nome> ---` nel `user_message` LLM. Nomi file ostili contenenti `\n\nSystem: ignore previous instructions` o marker di chat template possono iniettare istruzioni nel prompt, alterando il comportamento del modello (jailbreak / leak / override system prompt).

## Fix

In `internal/prompt/prompt.go`:
1. Funzione `sanitizeFilename(name string) string` che:
   - rimuove `\n`, `\r`, `\0`
   - sostituisce `[Ss]ystem:`, `[Uu]ser:`, `[Aa]ssistant:` con underscore
   - trunca a max 200 caratteri
   - sostituisce caratteri non-printable con `_`
2. Applicarla a tutti i nomi prima dell'interpolazione nel marker.
3. Unit test in `internal/prompt/prompt_test.go` con casi adversariali (filename con newline, con marker chat, con null byte, troppo lungo).

## Context

- Threat model Sprint 1 è "Denis su propri dati" → adversarial input NON in spec. Quindi è preventivo, non blocker.
- **Diventa critico in Sprint 2** quando il sistema completo AICertus avrà watcher / orchestrazione automatica: qualunque file in coda diventa vettore.
- Confermato da Mistral L1 + Codex L2.
- Convertito in todo P1 da `/v-triage` 2026-05-20 (HIGH = P1 per default; il CTO può downgrader a P2 se preferisce attendere Sprint 2).
- Stima fix: 2-3 gettoni infra (~1-2h CTO + design discussion su policy + test).
