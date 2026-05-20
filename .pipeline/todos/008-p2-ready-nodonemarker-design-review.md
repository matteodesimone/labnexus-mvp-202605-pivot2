---
id: 008
priority: P2
status: ready
source: bug
created: 2026-05-20
file: internal/provider/ollama.go
bug_ref: .pipeline/bugs/nodonemarker-semantica-completato.md
---

# EOF senza done marker mappa a `stato: completato` — semantica troppo permissiva?

## Issue

In `internal/provider/ollama.go` post-fix `eof-post-content-falsamente-interrotto` (commit 8a41239), un EOF della connessione Ollama dopo aver ricevuto del contenuto MA SENZA done marker viene mappato a `stato: completato` con un warning su stderr. La motivazione del fix era impedire false interruzioni quando Ollama termina normalmente ma senza emettere il done marker (succede su alcuni modelli reasoning).

**Il dubbio (design review)**: questo è troppo permissivo? Un EOF senza marker potrebbe anche significare "il modello ha avuto un crash interno e la connessione si è chiusa a metà output, ma noi NON lo sappiamo". Marcare come "completato" è ottimista e potrebbe nascondere bug del modello / network blip.

## Fix proposto (da discutere)

Tre alternative:
1. **Status quo**: lascia `completato` con warning. Pro: non interrompere lavori legittimi. Con: nasconde bug.
2. **Status `parziale-no-done`** nuovo: distingue `completato` (done marker) da `parziale-no-done` (EOF senza marker). L'utente vede chiaramente l'ambiguità. Pro: trasparenza. Con: nuovo stato da gestire downstream.
3. **Threshold-based**: se la quantità di token ricevuti supera una soglia ragionevole (es. > 100 token o > 30s), mark `completato`; altrimenti `interrotto`. Pro: euristica utile. Con: arbitrario.

Raccomandazione iniziale: opzione 2 (`parziale-no-done` come nuovo stato esplicito). Richiede discussione design col CTO + valutazione UX nella TUI.

## Context

- Bug di design review (non un fault funzionale). Convertito in todo P2 da `/v-triage` 2026-05-20.
- Stima fix opzione 2: 2-3 gettoni infra (~2h CTO + UX TUI + test).
- Buon candidato per design discussion in `/v-brainstorm` prima del fix.
