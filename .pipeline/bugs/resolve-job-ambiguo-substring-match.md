---
severity: high
status: open
created: 2026-06-02
source: external reviewers (Codex + Mistral) durante /v-review bugfix trigger-override-job
area: internal/jobs/jobs.go (Resolve)
---

# Bug: `Resolve` può eseguire il job sbagliato (substring match, primo che vince)

## Description

`jobs.Resolve` (internal/jobs/jobs.go:~179) fa case-insensitive **substring**
match sul nome cartella e ritorna il **primo** match in ordine alfabetico. Se due
job contengono lo stesso termine di ricerca (es. `--job revisione` con cartelle
"CAPABILITY A — Profilo revisione" e "Revisione straordinaria — Profilo revisione"),
viene scelto silenziosamente il primo, senza segnalare l'ambiguità.

Per il workflow di valutazione formale di Denis significa: input SGQ sbagliati al
modello e output scritto nella cartella sbagliata, senza alcun avviso.

Confermato da 2 reviewer (Codex HIGH, Mistral LOW).

## Expected

Precedenza: match esatto case-insensitive → altrimenti substring solo se UN solo
candidato → altrimenti errore esplicito di ambiguità che elenca i match validi.

## Note

Pre-esistente, non introdotto né smascherato dal bugfix trigger-override. È un
cambio di comportamento (semantica di matching) → richiede pipeline + test
dedicati, non fixabile inline in un bugfix scoped. Deferred.
