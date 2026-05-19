---
date: 2026-05-19
status: planned
source: compound-learning-fetta1
priority: medium
tags: ollama, robustness, sprint-2, oq-2
---

# Auto-derive del modello esatto da `ollama list`

## Idea

All'avvio di `labnexus run`, derivare il tag esatto del modello Qwen dal registry Ollama locale (`ollama list`) invece di sperare nel valore hardcoded del profilo. Risolverebbe strutturalmente OQ-2 (LabNexus Sprint 1) e il finding L4 della review.

## Razionale

I profili shippable (`profili/revisione.yml`, `profili/rilievi.yml`) hanno `modello: qwen3.6`. Ma il tag esatto nel registry Ollama dell'utente potrebbe essere:
- `qwen3.6`
- `qwen3.6:latest`
- `qwen3.6:14b-q5_K_M`
- `qwen3:latest` (se Denis ha qwen3 ma non qwen3.6)

Attualmente:
1. L'utente lancia `labnexus run`
2. Il binario chiama Ollama con `model: "qwen3.6"`
3. Se Ollama non ha quel tag esatto → errore HTTP 404 da Ollama
4. L'utente deve manualmente aggiustare il campo `modello:` nei file YAML

L'esperienza è scadente. Il CTO ha documentato in DEV-GUIDE.md che l'utente deve fare `ollama list | grep qwen3` al primo lancio. Questa idea propone di automatizzare la cosa.

## Approccio

1. Aggiungere `internal/provider/discover.go` con `DiscoverOllamaModel(prefix string) (string, error)`:
   - GET `http://localhost:11434/api/tags` → JSON con la lista modelli installati
   - Cerca match esatto del campo `modello:` del profilo
   - Se no match esatto: cerca match per prefisso (`qwen3` matcha `qwen3.6:14b-q5_K_M`)
   - Se ambiguo: prefer la variante più recente / più grande
   - Se nessun match: errore esplicito con suggerimento `ollama pull <modello>`
2. Aggiungere flag `--auto-model` (default true per ollama provider, false per eurouter)
3. Log esplicito del modello effettivamente usato: `[runlog] modello: profilo=qwen3.6 → registry=qwen3.6:14b-q5_K_M`
4. Frontmatter dell'output: `modello: qwen3.6:14b-q5_K_M` (tag risolto, non quello del profilo)

## Costo stimato

- ~50 righe di codice in `provider/discover.go` + unit test
- 1 nuovo scenario BDD ("Auto-derive risolve qwen3.6 → tag specifico")
- Aggiornamento DEV-GUIDE.md e CONTEXT.md
- Stima: ~1-2 ore CTO

## Quando affrontarla

- **Pre-Sprint 2 oppure all'inizio di Sprint 2**: il valore aggiunto è alto (riduce frizione operativa per Denis), il costo è basso
- **Pre-requisito**: nessuno
- **Trigger esplicito**: al primo `labnexus run` reale, se Denis riporta errore "modello non trovato"

## Beneficio collaterale

- Diagnostica migliore: il log riporta sempre il tag esatto, utile per riprodurre run passati anche dopo un `ollama pull` che aggiorna il modello
- Robustezza al re-pull: se Denis fa `ollama pull qwen3.6` e Ollama scarica una nuova versione con tag diverso, il binario si adatta da solo

## Vedi anche

- OQ-2 in `.pipeline/spec.md`
- Finding L4 in `.pipeline/archive/2026-05-19-fetta1-review.md` (sarà archiviato in questo compound)
- `internal/provider/ollama.go` — `OllamaProvider.doRequest`
- `.pipeline/solutions/2026-05-19-fetta1-scooter.md` — emersione

## Notes

(vuoto)
