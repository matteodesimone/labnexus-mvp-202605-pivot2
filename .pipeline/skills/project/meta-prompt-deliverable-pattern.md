---
slug: meta-prompt-deliverable-pattern
created: 2026-05-21
source_cycle: Fetta 3 (FR-20/21/22 feat-meta)
applies_to: deliverable testuali che guidano un LLM esterno a generare codice/config destinato a essere consumato dal sistema (es. profili YAML, migration SQL, configurazioni di routing)
---

# Pattern — Deliverable testuale + N functional-test executions

## Problema

Sprint 1 di LabNexus ha un deliverable diverso da tutti gli altri profili: **`feat-meta`** è un file MD da incollare in un LLM ESTERNO (Claude chat) per generare nuovi profili YAML. Non è codice eseguibile dal binario `labnexus`. La sua "correttezza" non si verifica con un go test su una pure function — si verifica funzionalmente eseguendo N casi predefiniti nel LLM esterno e controllando che gli output passino lo schema validation del sistema.

Questo è un **pattern strutturale generalizzabile** a qualsiasi situazione in cui:

- Il deliverable è documentazione (prompt, template, istruzioni) destinato a un AI esterno o un utente che produce input strutturato.
- La "validazione" passa per N esecuzioni esterne i cui output sono parsabili dal sistema.
- Servono regression-safe contro cambi schema futuri (se lo schema cambia, gli output committed esplodono e segnalano che il deliverable testuale va aggiornato).

## Pattern

Composto da 5 elementi:

1. **Deliverable testuale primario** (es. `docs/meta-prompt-genera-profilo.md`)
   - MD strutturato, sezioni numerate per facile reference.
   - Istruzioni complete e self-contained per l'LLM esterno (identità, schema atteso, lista risorse disponibili, esempi few-shot, regole di clarification, pattern di output).
   - **Schema e lista risorse esplicite** (no "inventa path KB" — Claude tende ad allucinare path).
   - Pattern di output con comando di verifica esplicito (es. *"l'output deve includere il comando `labnexus validate <nome>` come ultima riga"*).

2. **Guida d'uso** (es. `docs/guida-meta-prompt-denis.md`)
   - 1 pagina operativa per l'utente finale.
   - Step esatti: apri LLM → incolla deliverable → descrivi task → segui clarification → salva output → valida.

3. **N casi inventati di complessità crescente** (FR-22 dello spec)
   - **Caso semplice**: descrizione chiara → atteso output al primo colpo, 0 clarification.
   - **Caso medio**: descrizione plausibile ma con ambiguità minore → atteso ≥1 clarification.
   - **Caso ambiguo**: descrizione vaga → atteso multiple clarification, riconoscimento dell'ambiguità prima dello YAML.

4. **Cartella `casi-inventati/`** con README operativo + placeholders dei N output
   - README documenta i casi, le aspettative, e i passi operativi per il CTO/Power User.
   - Gli N output (es. `caso-1-semplice.yml`) vengono **committati DOPO esecuzione manuale**.

5. **Test Go table-driven regression-safe**
   - Cicla i file `caso-*.<ext>` nella cartella, parsa, valida contro lo schema reale.
   - In RED phase (cartella vuota o < N file): FAIL con messaggio diagnostico esplicito che dice all'umano cosa fare.
   - In GREEN phase (N file committati): pass.
   - Se lo schema cambia in futuro e gli output diventano invalidi, il test esplode → segnale che il deliverable testuale va aggiornato.

## Implementazione canonica (LabNexus Fetta 3)

```
docs/
  meta-prompt-genera-profilo.md       # deliverable primario, 8 sezioni
  guida-meta-prompt-denis.md          # guida utente 1 pagina

.pipeline/test-data/scenarios/feat-meta/casi-inventati/
  README.md                            # istruzioni operative + tabella dei 3 casi
  caso-1-semplice.yml                  # committed post-esecuzione manuale
  caso-2-medio.yml                     # committed post-esecuzione manuale
  caso-3-ambiguo.yml                   # committed post-esecuzione manuale

internal/profile/profile_meta_cases_test.go
  TestMetaPromptCases       # table-driven sui caso-*.yml, valida contro KB reale
  TestMetaPromptContent     # FR-20 regression: il deliverable MD ha N sezioni richieste, K risorse citate, ≥M esempi few-shot
  TestGuidaMetaPromptContent # FR-21 regression: la guida menziona keyword operative chiave
```

## Documentare il "red-by-design"

Il `TestMetaPromptCases` resta **deliberatamente RED** finché l'umano non committa i N output. Questo non è una regressione ma un **gate human esplicito**. Per evitare confusione, va documentato in **4 punti**:

1. **Messaggio diagnostico del test**: deve essere educativo, citare i passi operativi esatti, e il path dove committare gli output.
2. **`state.json` field `pending_human_work`**: descrive il blocco.
3. **`README.md` della cartella `casi-inventati/`**: istruzioni operative dettagliate.
4. **`review.md`**: marker esplicito "[HIGH/EXPECTED]" o "[STATUS: gate human]" per evitare che la review futura interpreti il red come bug.

## Trade-off accettato

- **Costo**: il CTO/Power User deve eseguire manualmente N casi nell'LLM esterno. Stima Sprint 1: ~30-60 min per i 3 casi FR-22.
- **Beneficio**: regression-safe contro schema changes futuri (1 minuto per scoprire che gli output devono essere ri-generati) + il deliverable testuale non è un "blob non testato".

## Generalizzazioni future

Il pattern si applica direttamente a:

- **Template SQL/migration** da generare via LLM esterno (test = run del schema validation sui file committati).
- **Config routing/ingress** generati via LLM (test = parse + apply dry-run).
- **Documentazione strutturata** (es. ADR) generata via LLM (test = check di sezioni canoniche presenti).

Anti-pattern da evitare:

- **NON** lasciare il deliverable testuale senza test funzionale ("è solo un MD"). Test Go regression-safe è economico (~20 LoC) e protegge dalle modifiche silenziose.
- **NON** committare i N output di esempio insieme al deliverable, perché perdi la separation "ho creato il deliverable" vs "ho dimostrato che il deliverable funziona". Tieni i due commit separati.
