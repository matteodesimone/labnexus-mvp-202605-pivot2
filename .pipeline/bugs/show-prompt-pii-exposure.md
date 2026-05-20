---
slug: show-prompt-pii-exposure
severity: MEDIUM
status: open
opened_at: 2026-05-20
opened_by: /v-review Fetta 2 (Mistral loop 1 security, Codex loop 2 privacy)
discovered_via: external review (multi-LLM)
---

# Bug — `labnexus check --show-prompt` può esporre PII via stdout/copy-paste

## Sintomo

Il comando `labnexus check <profile> --input <dir> --show-prompt` (FR-5, dry-run) stampa il prompt LLM completamente composto su stdout, includendo:
- Tutti i contenuti dei `kb_files` (es. `come-pensa-un-ispettore.md`, sezioni ISO)
- Tutti i contenuti dei file di input (es. NC dell'anno con nomi tecnici, rapporti SGQ con riferimenti a persone)

Se l'utente esegue:
```bash
labnexus check review-pack --input /private/SGQ-Denis --show-prompt > debug.txt
```
… oppure piasta l'output in un ticket di supporto o lo copia in Slack, **dati personali di colleghi di Denis finiscono in canali non controllati**, violando NFR-6 / NFR-7 e GDPR (data minimization).

## Analisi

- `internal/runner/runner.go:70` (circa) gestisce il dry-run con `--show-prompt`.
- L'output è ESPLICITAMENTE richiesto dall'utente (il flag è un comando di debug), quindi tecnicamente è "informed consent".
- Tuttavia il nome del flag (`--show-prompt`) non avvisa l'utente sull'esposizione potenziale di PII contenuti nei kb_files o input.

## Impatto

- **Sprint 1 reale**: Denis è un singolo utente che lavora con i suoi dati. Il rischio è "incidente di copia/paste" → MEDIO.
- **Sprint 2 / produzione multi-utente**: rischio cresce. Operatori meno esperti potrebbero pipare l'output altrove.

## Fix proposto

Opzione A (soft, ~30 min):
- Documentare warning prominente nel README e nell'help text di `check --show-prompt`: *"⚠ Output contiene PII degli input. Non condividere su canali non controllati."*
- Aggiungere check TTY: se stdout NON è una TTY (pipe / file), stampare prima il warning su stderr.

Opzione B (hard, ~2-3 gettoni infra):
- Flag opzionale `--redact` che maschera nomi propri (regex semplice), numeri di rapporto, email, codici personali.
- Default `check --show-prompt` chiama `--redact` automaticamente; per output integro l'utente deve `--no-redact` con conferma esplicita.

Raccomandazione: Opzione A per Sprint 1 + Opzione B in Sprint 2 quando il prodotto è multi-utente.

## Note di triage

- Convertito in todo P2 da `/v-triage` 2026-05-20.
- Mitigation soft (Opzione A) può essere chiusa rapidamente prima dell'handoff Sprint 1 a Denis.
