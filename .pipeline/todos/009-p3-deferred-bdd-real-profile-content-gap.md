---
id: 009
priority: P3
status: deferred
source: review
created: 2026-05-20
file: features/steps_test.go
bug_ref:
---

# BDD non esercita end-to-end il contenuto reale dei profili Fetta 2

## Issue

`features/steps_test.go:installFromProject` ora valida lo schema del profilo reale via `labnexus validate` (fix loop 1 review), ma per il successivo `labnexus run` swap-in `minimalProfileYAML(name)` con solo `CLAUDE.md` come kb_file e un trigger_prompt stub. Quindi:

- **Validato automaticamente in BDD**: schema profilo (campi, tipi, trigger_prompt ≥ 50 char, kb_files esistono nella KB reale).
- **NON validato in BDD**: composizione effettiva del prompt con i kb_files multipli reali (sezioni-ISO, NC-patterns, MAPPA) + il `trigger_prompt` REALE del profilo.

Un cambiamento accidentale al trigger_prompt che indebolisce le istruzioni (es. perdita della frase "esattamente 13 sezioni") non sarebbe catturato dalle BDD attuali — solo dal manuale "L1 CTO sul Qwen reale" o dall'L2 di Denis.

## Fix (DEFERRED)

Chiusura piena richiederebbe:
1. Helper test che copia `profili/<name>.yml` reale in `s.profiliDir` (no swap stub).
2. Helper che copia tutti i kb_files referenziati dal profilo reale dentro `s.kbDir`.
3. Adjustment al fake Ollama per non esplodere su prompt molto grandi (e gestire UTF-8 chunking attualmente fragile in `splitInThree`).
4. Possibili nuovi scenari `scenarios real-content` taggati `@slow` (esclusi dal CI normale, runnati on-demand).

Stima: ~50-100 LoC test infra + 30-60min CTO + rischio flakiness da `splitInThree` UTF-8 bug pre-esistente.

## Context — perché DEFERRED

- ACCEPTED RISK documentato in `.pipeline/review.md` Fetta 2 (sezione "Persistent finding").
- **Mitigation in essere**: lo shakedown reale CTO su Qwen 3 (`labnexus run --profile <X> --input <fixture> --output <out>`) esercita la composizione completa del prompt. L1 manuale + L2 Denis chiudono il gap per NFR-8.
- Convertito in todo P3 DEFERRED da `/v-triage` 2026-05-20.
- Promozione a P2 se in Sprint 2 il sistema diventa multi-utente o se l'iterazione sui trigger_prompt richiede regressione automatica.
