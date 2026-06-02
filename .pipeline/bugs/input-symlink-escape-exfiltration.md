---
severity: high
status: open
created: 2026-06-02
source: external reviewer (Codex security) durante /v-review bugfix trigger-override-job
area: internal/input/input.go (ParseDir / extract)
---

# Bug: i file di input via symlink possono esfiltrare file fuori dalla cartella job

## Description

`ParseDir` (internal/input/input.go:~148) cammina l'albero della cartella di
input, registra i file (anche symlink) e li legge via `extract(...)` **senza** il
controllo di contenimento `Abs` + `EvalSymlinks` già usato in
`ParseTriggerPromptFile`. Un symlink dentro la cartella di lavoro, es.
`lavori/<job>/segreto.txt -> ~/.ssh/config` o un altro file SGQ leggibile,
verrebbe parsato e **concatenato nel prompt inviato al provider** (default
EUrouter cloud).

Distinto da `kb-symlink-escape.md` (quello riguarda la cartella KB; questo la
cartella di input/lavoro).

## Expected

Applicare lo stesso check di contenimento `Abs` + `EvalSymlinks` di
`ParseTriggerPromptFile` a ogni file di input prima dell'estrazione, oppure
skippare i symlink di default con warning esplicito.

## Note

Pre-esistente, non correlato al trigger-override (il path del trigger file È già
containment-checked). Finding di sicurezza reale → pipeline dedicata. Deferred.
