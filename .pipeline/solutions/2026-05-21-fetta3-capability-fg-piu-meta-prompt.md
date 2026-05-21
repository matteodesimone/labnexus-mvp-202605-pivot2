---
date: 2026-05-21
pipeline: new-feature
feature: Fetta 3 — competence-gap + pt-analysis + feat-meta (chiusura AI-side Sprint 1)
stack: go-cli (godog BDD)
tags: declarative-profiles, llm-meta-prompt, regression-test-pattern, hallucination-prevention, markdown-emphasis
---

# Fetta 3: 2 capability + prompt meta deliverable (chiusura AI-side Sprint 1)

## Problem

Chiudere Sprint 1 con le 2 capability nuove più rischiose (competence-gap = F, pt-analysis = G — quest'ultima rischio max per Qwen numerical) + un deliverable testuale che permetta a Denis di generare nuovi profili in autonomia in chat Claude esterna (feat-meta = FR-20/21/22).

## Solution

Cycle Fetta 3: plan → test-scaffold → implement → review (1 loop) → deploy → compound. Pattern Fetta 2 ri-applicato senza friction (3 skills usate as-is + 1 estesa). NUOVO pattern emerso per feat-meta: "deliverable testuale + N functional-test executions" — il deliverable è MD non Go code; la validazione FR-22 è ibrida (3 esecuzioni umane Claude esterno → YAML committed → test Go regression-safe). 1 fix inline review (#4 PT trend), 1 fix inline review (#8 meta-prompt content regression). `TestMetaPromptCases` lasciato RED per design (gate human, esplicita scelta utente).

## What Worked

- **Pattern Fetta 2 ri-applicato a zero friction**: schema validation gate, hallucination-prevention via fixture-set, cloud opt-in gate. Le 3 skills create in Fetta 2 compound hanno coperto direttamente le esigenze F+G. Skill hallucination-prevention ESTESA a 3 classi (nomi tecnici, codici PT, codici metodo) senza riscrittura.
- **Plan accuracy alta**: 3 sub-fette × 1 gettone si è confermata corretta. Solo 1 fix engine inevitabile (#4 trend assertion) + 1 test mancante (#8 meta-prompt regression).
- **Review 1 loop sufficiente** (vs 2 di Fetta 2): pattern ormai stabili = meno finding nuovi. False positives mistral ben classificabili (PRE-EXISTING / TD-1 design).
- **Markdown-emphasis-tolerant regex caught early**: la fail iniziale `rischio **alto**` ha fatto emergere la lezione "regex BDD su output LLM devono essere markdown-aware da subito". Fix preventivo per real Qwen output.

## What Didn't Work

- **Naming temporale "fetta2ExpectedDefaults"**: ho dovuto rinominare a `profileExpectedDefaults` mentre estendevo con F+G. Lesson: nomi semantici da subito (`profileExpectedDefaults`), non temporali (`fettaNExpected*`). Nomi temporali invecchiano male appena la mappa cresce oltre la fetta originale.
- **Documentazione del red-attended-by-human**: `TestMetaPromptCases` red era ambiguo finché non documentato in 4 punti (test message + state.json + README casi-inventati + review.md). Lesson per pattern futuri "deliverable testuale + functional test": **documentare red-by-design ESPLICITAMENTE quando si scrive il test**, non dopo.
- **Markdown emphasis regex tight**: `rischio\s+alto` mi è costato ~5 min di debug iniziale. Va in coda alla lesson sopra (pattern: regex su LLM output va sempre testata con markdown formatting da subito).

## Reusable Pattern

**Pattern NUOVO**: "Deliverable testuale + N functional-test executions" (vedi `.pipeline/skills/project/meta-prompt-deliverable-pattern.md` per dettagli):

- Il deliverable è MD (es. prompt da incollare in LLM esterno), non Go code.
- N casi predefiniti (FR-22 = 3 casi: semplice/medio/ambiguo) sono il "test set".
- Validazione = test Go table-driven che cicla i N YAML committati e li verifica contro lo schema vero.
- Gate-human esplicito (RED → GREEN solo quando l'umano committa i YAML).
- Regression-safe per cambi schema futuri (se lo schema profile cambia in Sprint 2, i casi committed esplodono come red).

## System Updates Applied

- **1 skill di progetto creata**: `.pipeline/skills/project/meta-prompt-deliverable-pattern.md` (vedi sezione sopra).
- **1 skill estesa**: `.pipeline/skills/project/hallucination-prevention-fixture-set.md` Fetta 2 — pattern applicato a 3 nuove classi (nomi propri italiani, codici PT-YYYY-NN, codici PCM-NN). Skill rimane valida invariata; le tabelle in essa diventano referenza per future fette.
- **Naming convention**: rinominato `fetta2ExpectedDefaults` → `profileExpectedDefaults`. Implicita lesson per future mappe condivise cross-fetta: nomi semantici, non temporali.

## Tech Debt

1. **`features/steps_test.go` ~2400 lines** — già notato in Fetta 2 compound come "borderline". Ora con Fetta 3 cresciuto a ~2400 LoC. Tutto verde + funzionante; refactor sarebbe split in `profili_steps_test.go` + `provider_steps_test.go` + `assertions_test.go`. **Acceptable per Sprint 1 (chiusura)**; se Sprint 2 aggiunge BDD nuovi, refactor diventa raccomandato.
2. **Regex parsing CSV in `fixtureAllowedNomiTecnici`**: parser ingenuo (split su prima virgola). Funziona per le fixture attuali (nomi senza virgole). Se Sprint 2 ha nomi con virgole interne (es. "Rossi, Maria"), va sostituito con `encoding/csv`. **Acceptable** — caso non realistico.
3. **`splitInThree` pre-existing UTF-8** — già noto da Fetta 2 compound, sempre pending. Non smascherato in Fetta 3 (canned content abbastanza piccolo). **Tracciato in todo #009 (P3 deferred)**.

**Raccomandazione**: tech debt non significativa, no `/v-refactor` necessario. Sprint 1 è chiuso AI-side, eventual refactor è discussione Sprint 2.

## Sprint 1 closure status (AI-side)

- ✅ Motore (Fetta 1): shipped + stabilizzato (6 bugfix + 1 UX patch).
- ✅ 7 profili (Fetta 1 A/B + Fetta 2 C/D/E + Fetta 3 F/G): tutti shippable, validate schema OK x7, BDD L1 strutturale verde x7.
- ✅ Prompt meta (FR-20/21): docs deliverable in zip al root.
- ✅ Security/hardening: 5 bugfix P1+P2 chiusi (eurouter gate, KB symlink, ollama loopback, prompt injection, frontmatter_default).
- 🟡 Gate fuori AI: L2 Denis 5 capability nuove + FR-22 3 casi Claude + L2 usabilità meta-prompt. **3 gettoni client Fetta 2 + 3 gettoni Fetta 3 = 6 gettoni pending consumption nello shakedown reale**.

Budget Sprint 1: 14/20 spent (8 pre-sprint + 6 Fetta 1). 6 gettoni residui = quelli del lavoro umano post-deploy. Tetto interno 15 mai sforato.

## Pipeline cycle complete: spec → plan → test-scaffold → implement → review (PASS WITH NOTES, 1 loop) → deploy → compound ✓
