---
severity: medium
status: open
created: 2026-05-19
source: review esterna (codex-quality HIGH)
fix: ""
test: ""
---

# Bug (design review): EOF senza done marker mappa a `stato: completato` con warning — semantica troppo permissiva?

## Reported

- **Feature/Area**: `internal/provider/ollama.go:consumeOllamaStream` (emette `Done:true, NoDoneMarker:true`) + `internal/runner/runner.go:drainStream` (mappa `Done → stato=completato`)
- **Trovato da**: review codex-quality HIGH 2026-05-19
- **Sprint impact**: MEDIUM. Design del bugfix-4 (`eof-post-content-falsamente-interrotto`).

## Description

Il fix del bug-4 ha cambiato la semantica di "EOF post-content senza done marker": prima → `stato: interrotto`; adesso → `stato: completato` + warning. Razionale del bug-4: il caso reale era "Ollama emette tutti i token + chiude la connessione TCP senza emettere il chunk finale done:true", quindi l'output era de facto completo.

Codex review obietta:
> EC-8 requires partial/interrupted streams to be written with stato: interrotto and exit 1. Marking uncertain output as completed can send truncated model output into Denis's quality validation as if it were reliable.

Tradeoff:
- **Status quo (post-bug-4)**: NoDoneMarker → completato. Pro: il caso reale produce output usable. Con: un troncamento mid-content viene riportato come "completato" → Denis può ricevere documento incompleto credendo che sia integro.
- **Alternativa**: introdurre `stato: completato_senza_marker` (terzo stato) + exit code dedicato (es. 3) → Denis vede esplicitamente "completezza incerta, da rivedere". Pro: granularità. Con: complicazione UX, schema frontmatter extended.
- **Alternativa stricta**: NoDoneMarker → `stato: interrotto`. Pro: conservative, no falsi positivi. Con: il caso reale comune produce file marcati "interrotto" anche se de facto completi (degrado UX).

## Fix proposto

Decisione di prodotto richiesta (CTO + Denis):
1. Mantenere status quo + warning (semplicità) — accettabile se Denis legge sempre l'output prima di validare
2. Introdurre `completato_senza_marker` (precisione) — più safe per produzione
3. Tornare a `interrotto` (massima cautela) — più strict

## Note

Deferred to backlog: design decision, non bug funzionale. Da rivisitare al checkpoint Sprint 1 con Denis dopo i primi run reali. Probabilmente opzione 2 quando si esce dal prove-out.
