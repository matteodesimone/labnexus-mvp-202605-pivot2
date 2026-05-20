---
id: 001
priority: P1
status: done
source: bug
created: 2026-05-20
file: internal/provider/provider.go
bug_ref: .pipeline/bugs/privacy-eurouter-gate-mancante.md
---

# Eurouter provider va gated dietro approval esplicito

## Issue

`provider.Select` accetta `eurouter` ogni volta che il flag `--provider eurouter`, l'env `LABNEXUS_PROVIDER=eurouter`, o un profilo lo dichiara, senza verificare che l'uso sia esplicitamente autorizzato sui dati reali. Un override accidentale o un env stantio può silenziosamente esfiltrare dati SGQ di Denis al gateway EU OpenAI-compatible (api.eurouter.ai), violando NFR-1 + la promessa contrattuale post-pivot 2 "tutta l'elaborazione AI gira locale".

## Fix

Gate `eurouter` dietro un meccanismo esplicito e difficile da impostare per errore. Opzioni:
1. **Soft gate (Sprint 1)**: env var dedicata `LABNEXUS_ALLOW_CLOUD_PROVIDER=approved-for-synthetic-data` richiesta in aggiunta a `LABNEXUS_PROVIDER=eurouter`. Se manca: rifiuta con messaggio esplicito di blocco.
2. **Hard gate (raccomandato per deploy a Denis)**: in release build (`-tags release` o equivalente) `provider.Select` SOLO `ollama` ammesso, qualsiasi tentativo eurouter restituisce errore. Dev build (default per CTO) accetta eurouter previo soft gate.

Test: BDD scenario "lancio con --provider eurouter senza ALLOW_CLOUD: exit 2 + stderr 'gate non superato'".

## Context

- Confermato da Codex in 2 review consecutive (Fetta 1 review backlog + Fetta 2 loop 1+2).
- **Blocker per la consegna a Denis**: senza questo gate, un errore di configurazione del CTO durante il setup hardware potrebbe inviare dati reali del laboratorio al cloud.
- Convertito in todo P1 da `/v-triage` 2026-05-20.
- Stima fix: 2-3 gettoni infra (~2h CTO + test).
