---
date: 2026-05-21
pipeline: new-feature + 5 bugfix consecutivi
feature: Fetta 2 — review-pack + audit-checklist + equipment-alert + hardening security pre-handoff Denis
stack: go-cli (godog BDD)
tags: declarative-profiles, multi-llm-review, security-hardening, bdd-install-pattern, unicode-italian, frontmatter-default, prompt-injection, kb-symlink-escape, loopback-gate
---

# Fetta 2: 3 capability declarative + hardening security pre-handoff

## Problem

Implementare 3 capability nuove (review-pack, audit-checklist, equipment-alert) come **declarative-only** (profili YAML + KB selection + trigger_prompt) senza toccare l'engine. Validare il L1 strutturale via BDD canned, lasciando L2 a Denis. Subito dopo, chiudere il backlog di security findings emerso dalla review (eurouter exfiltration, KB symlink escape, Ollama remote, prompt injection via filename, frontmatter_default non cablato).

## Solution

Cycle Fetta 2: spec (Fetta 1 shared) → plan (3 sub-fette dichiarative) → test-scaffold (3 BDD scenari + 10 fixture sintetiche) → implement (3 profili YAML) → review (REVIEW="all", 2 loop, PASS WITH NOTES) → deploy (zip rebuild con 5 profili) → 5 bugfix consecutivi (4 P1 security + 1 P2 engine) → compound.

L'approccio chiave è stato il pattern "schema-validation gate in BDD install step": `installFromProject` esegue `labnexus validate <name>` contro profili/ + KB reale come precondizione del run. Questo cattura regressioni dello schema (kb_files inesistente, trigger troppo corto, provider sbagliato) AUTOMATICAMENTE, senza esercitare end-to-end il content composition (lasciato a shakedown manuale + L2 Denis).

Per i bugfix security, pattern unificato "Cloud opt-in gate" via `LABNEXUS_ALLOW_CLOUD_PROVIDER`: eurouter è rifiutato senza gate, e anche `LABNEXUS_OLLAMA_ENDPOINT` non-loopback richiede lo stesso gate. Singolo punto di decisione, due difese.

## What Worked

- **Plan accurate**: "Fetta 2 = no engine logic" si è confermata corretta (3 profili declarative bastavano per L1). 3 gettoni client allocati = 0 gettoni AI-side spesi (budget-free per la pipeline; i gettoni reali si spendono nello shakedown Qwen + L2 Denis).
- **Multi-LLM review economico**: codex + mistral hanno trovato 7 finding in-scope reali (3 HIGH + 3 MEDIUM + 1 unmasked engine bug), oltre a confermare 5 backlog Fetta 1. 9 false positives dismissibili con chiara motivazione.
- **Triage → /v-bugfix all P1 in batch**: 4 cicli reproduce-test → minimal-fix → suite-verify chiusi senza checkpoint per-bug, ~45 min totali. Pattern ripetibile per future "hardening pass" prima di deploy.
- **Pattern `bodyCitaMetodiAssociati(only-from-fixture-set)`**: il fix loop 1 H3 trasforma una regex generic in un check che parsa il fixture, estrae il set ammesso, e rifiuta allucinazioni. Direttamente applicabile a tutte le capability future con riferimenti tecnici verificabili (PCM-NN, codici NC, sezioni ISO).

## What Didn't Work

- **TestMain trap**: `features/main_test.go` non chiamava `m.Run()` dopo `godog.TestSuite.Run()`. Risultato: TUTTI gli unit test in `features/*_test.go` venivano silenziosamente skippati per 2 fette. Scoperto solo in review loop 1 quando ho cercato di vedere i risultati dei miei nuovi unit test. **CRITICAL trap** per progetti Go che mescolano godog + unit test nello stesso package.
- **Go `\b` ASCII-only**: il mio primo fix di `hasCausalConnective` usava `(?i)\b(poiché|quindi|perché|se|implica)\b`. Il regex Go non riconosce `é/à/ò/ù/ì` come word char → "poiché succede X" non matchava. Soluzione: lookbehind/lookahead via classi Unicode `(^|[^\p{L}])` + match + `([^\p{L}]|$)`.
- **`stato` semantic collision**: i 5 profili dichiarano `output.frontmatter_default.stato: bozza_da_validare_qm` (workflow QM); l'engine assegna `stato: completato/timeout` (runtime status). Stessa chiave, semantica diversa. Smascherato solo dopo aver cablato `frontmatter_default` nell'engine (review loop 1 M4 → bugfix #006). Rinominato profili in `stato_qm`. Lesson: l'introduzione di defaults richiede un audit del namespace condiviso engine/profile.
- **2 Mistral CRITICAL false positives** consecutivi (loop 1 + loop 2): primo era una misread (skip list letta come include list); secondo cross-referenziava step Fetta 1 non in scope Fetta 2. Costo: ~5 min di analisi per dismiss. **Pattern emerso**: Mistral tende a salti di contesto quando lavora su file molti diversi tipi (codice + YAML + Markdown spec) — vale la pena fare cross-check rapido su CRITICAL findings prima di trattarli come bloccanti.
- **Persistent HIGH accepted-risk**: BDD non esercita end-to-end il content composition reale (trigger_prompt + kb_files multipli). Schema validation gate è chiuso; full e2e richiede ~50 LoC test infra e potrebbe sbattere su splitInThree UTF-8 pre-existing bug. Accettato + mitigato dallo shakedown reale CTO + L2 Denis. Promote a P1 quando il sistema diventerà multi-utente in Sprint 2.

## Reusable Patterns

1. **Schema validation gate in BDD install step**: per ogni capability declarative, la step di install verifica l'esistenza E la schema-validity del profilo di progetto via subprocess `labnexus validate`. RED gate per regressioni schema senza richiedere e2e content. Vedi `features/steps_test.go::installFromProject`.

2. **Cloud opt-in gate unico**: `LABNEXUS_ALLOW_CLOUD_PROVIDER` come singolo env var per autorizzare ANY non-local provider (eurouter, ollama remoto). Default = strict local. Soft + economico, copre eurouter e ollama-endpoint-leak in un solo concept. Vedi `internal/provider/provider.go::ErrEurouterGateMissing` + `requireLoopbackOrCloudGate`.

3. **Hallucination-prevention via fixture-extracted set**: per validare riferimenti tecnici verificabili in output LLM, parsare il fixture, estrarre il set ammesso, e asserire che TUTTI i riferimenti nell'output appartengano a quel set. L'esempio è `bodyCitaMetodiAssociati` per PCM codes dalla scheda apparecchiatura. Pattern generalizzabile a: codici NC, sezioni ISO, ID risk register, nomi fornitori.

4. **Filename sanitization in LLM prompt composer**: filename ostili possono iniettare istruzioni nel `user_message`. `sanitizeFilename` rimuove control chars, neutralizza chat role markers (`System:/User:/Assistant:/Developer:`), trunca a 200 char. Default-deny anche per sistemi non-adversarial — buona igiene.

5. **EvalSymlinks + post-resolve containment**: per validare che file referenziati restino dentro una cartella designata, usare `filepath.EvalSymlinks` su entrambi (root + candidato) PRIMA del prefix check. Caso legittimo "root stesso è symlink" automaticamente coperto.

## System Updates Applied

Skill di progetto da creare in `.pipeline/skills/project/` (3 pattern di alto valore):

1. `bdd-schema-validation-gate.md` — pattern Schema validation gate in BDD install step
2. `cloud-opt-in-gate.md` — pattern Cloud opt-in gate unico
3. `hallucination-prevention-fixture-set.md` — pattern Hallucination-prevention via fixture-extracted set

Vedi sezione "Framework promotion candidates" sotto per gli altri.

## Framework Promotion Candidates (proposed)

Tutti i 5 pattern sopra hanno valenza generale oltre LabNexus. Candidati:

- **`testmain-mrun-trap.md`** (skill o standard layer 2 go-cli stack): warning su `TestMain` Go che chiama `godog.TestSuite.Run()` senza `m.Run()` — gli unit test del package vengono silenziosamente skippati. Da aggiungere come check in `.pipeline/checks.yaml` test-coverage lens, o nello standards/stacks/go-cli.md.
- **`unicode-word-boundary-for-non-english.md`** (skill general coding): Go `\b` è ASCII-only. Pattern `(?i)(^|[^\p{L}])...([^\p{L}]|$)` per accenti italiani/francesi/spagnoli/portoghesi.
- **`filename-sanitization-llm-prompt.md`** (skill stack agent o domain LLM): sanitize filename PRIMA di interpolare in prompt LLM. Standard practice anche fuori adversarial threat model.
- **`schema-validation-gate-in-bdd.md`** (skill stack cli-tool con declarative configs): pattern install step che esegue il binary's validate sui file di progetto reali come precondizione di BDD scenari.
- **`cloud-opt-in-gate-pattern.md`** (skill security cross-stack): gate cloud unico per multiple difese (provider remoto, endpoint override). Singolo punto di decisione user-visible.

## Future Ideas

- **BDD real-profile-content e2e** (todo #009 deferred): closure del persistent HIGH quando il sistema sarà multi-utente. Stima: ~50 LoC + investment per fixare splitInThree UTF-8 pre-existing.
- **Audit del namespace key engine-vs-profile** (nuovo): tabella documentata in 02-project.md di chi possiede ogni chiave del frontmatter (engine vs profile default). Previene future collisioni come `stato`.
- **Pattern `--redact` per `--show-prompt`** (todo #005 collegato): Sprint 2 quando multi-utente. Maschera nomi propri / codici personali / email con regex semplice.

## Tech Debt

Audit del codice toccato in questo cycle:

1. **`features/steps_test.go` ~2800 lines** — il file è cresciuto significativamente con Fetta 2 (BDD scenari + helper) + bugfix pass (test env helpers). NON critico ora, ma da considerare uno split per dominio (`profili_steps_test.go`, `provider_steps_test.go`, `assertions_test.go`) prima di Fetta 3 se la crescita continua. **Acceptable, proceed.**
2. **`splitInThree` pre-existing UTF-8 fragility** — il fake Ollama splitta byte-wise, può spezzare caratteri accentati italiani nel canned content. Pre-existing (Fetta 1), non smascherato in Fetta 2 grazie alle canned content corte. Da fixare prima di test e2e con real-profile-content (todo #009). **Acceptable, proceed; tracciato implicitamente in #009.**
3. **`internal/output/output.go::mergeFrontmatterDefaults` regex-based** — non gestisce frontmatter con chiavi nested (es. `output.frontmatter_default.composite_key.subkey`). Sprint 1 non ha defaults nested, ma se Sprint 2 dovesse introdurli serve refactor a yaml.Node manipulation. **Acceptable for current scope.**

Raccomandazione: **proceed**. Tech debt non significativa, no `/v-refactor` necessario prima di Fetta 3.

## Pipeline Cycle Closure

- Fetta 2/3 chiusa: spec → plan → test-scaffold → implement → review (PASS WITH NOTES, 2 loop) → deploy → compound ✓
- Bugfix pass: 5 bug chiusi (privacy-eurouter-gate-mancante, kb-symlink-escape, ollama-endpoint-env-override-leak, prompt-injection-via-filename, frontmatter-default-non-applicato).
- Backlog rimanente: 3× P2 ready (#005 show-prompt PII, #007 paths in log, #008 nodonemarker design) + 1× P3 deferred (#009).
- Commit del cycle: cc22aa3, 55b3fa2, d3ca09e, e15a54a (4 commit, ~250 file modificati/aggiunti).
- Next: Fetta 3 (competence-gap + pt-analysis + prompt meta) con `/v-plan`. L2 di Denis su Fetta 2 capability è gate fuori pipeline AI.
