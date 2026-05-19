# Plan: Sprint 1 LabNexus — Fetta 2: 3 capability nuove a rischio basso/medio

## Fetta Overview

- **FRs covered**: FR-15 (`review-pack`, Capability C), FR-16 (`audit-checklist`, Capability D), FR-17 (`equipment-alert`, Capability E).
- **After this fetta**: l'utente può eseguire `labnexus run --profile {review-pack, audit-checklist, equipment-alert} --input <dir> --output <dir>` end-to-end su Qwen 3 locale e ricevere un file markdown strutturato secondo il template CoWork KB della rispettiva capability. **Nessuna nuova logica del motore**: questa fetta è interamente "data + prompt engineering" sui profili YAML e i loro trigger.
- **Estimated effort**: 3 gettoni = ~6-12h CTO (1 gettone/profilo, ognuno include scrittura profilo + sintetici plausibili + BDD scenario L1 + esecuzione real su Ollama + giudizio L2 Denis).
- **Complessità**: Medium (2-8h).

## Architecture Decision

Nessun cambio architetturale al motore. Tutto il lavoro di Fetta 2 è **declarative** (file YAML profilo + Markdown KB-files + test fixtures + Gherkin scenari). La pipeline `runner.Run` esistente carica i nuovi profili senza modifiche di codice.

### Logic/Presentation Separation

Invariato. Tutta la **logica** vive in `internal/` (già implementata in Fetta 1) — i nuovi profili YAML sono **configurazione**, non codice. La **presentation** (`cmd/labnexus/` + `internal/tui/`) non cambia.

**Inserimento dei nuovi profili**:
- `profili/<name>.yml` — 3 nuovi file dichiarativi (provider: ollama, modello: qwen3.6, kb_files curati, trigger_prompt 5-10 righe)
- `internal/profile/profile.go` esistente li carica e valida con la stessa logica di Fetta 1

### Design Patterns Applied

| Pattern | Where | Justification |
|---|---|---|
| Validazione L1+L2 (skill di progetto) | BDD + gate L2 Denis | Pattern documentato in `.pipeline/skills/project/validazione-due-livelli-cto-cliente.md`. Stesso approccio Fetta 1 per le 3 nuove capability. |
| Dual-provider env override (skill) | BDD fake server + subprocess test | `.pipeline/skills/project/dual-provider-env-override.md`. Stesso approccio Fetta 1 per i nuovi scenari BDD. |
| (nessun nuovo pattern di codice) | — | Riuso architettura Fetta 1. Niente Strategy/Repository/Factory aggiunto. |

### Embedded Artifact Dependency Check

**Non applicabile**: nessun `go:embed` o artefatto compilato nella Fetta 2. I profili YAML e i fixture sono letti runtime dal filesystem (decisione di Sprint 1 — abilita iterazione del cliente sulla KB senza ricompilazione).

## Files to Create

### Profili (artefatti runtime — invariati dopo build)

| File | Purpose | Layer | Dependencies |
|---|---|---|---|
| `profili/review-pack.yml` | Capability C — Management Review Pack annuale (CoWork KB 15.4, 13 sezioni) | declarativo | `internal/profile/validate` (esistente) |
| `profili/audit-checklist.yml` | Capability D — checklist d'audit da risk register | declarativo | idem |
| `profili/equipment-alert.yml` | Capability E — equipment alert (CoWork KB 15.3) | declarativo | idem |

Selezione `kb_files` per ognuno (preliminare, da affinare al lancio reale):
- **review-pack**: `CLAUDE.md` + `come-pensa-un-ispettore.md` + `MAPPA_ESPLOSA_REQUISITI_ISO17025.md` + `sezioni-ISO/sezione-8-sgq-domande.md` (riesame direzione) + `NC-patterns/*` (per la gestione NC del Pack)
- **audit-checklist**: `CLAUDE.md` + `come-pensa-un-ispettore.md` + `MAPPA_ESPLOSA_REQUISITI_ISO17025.md` + tutte le `sezioni-ISO/*.md` (audit copre tutto)
- **equipment-alert**: `CLAUDE.md` + `come-pensa-un-ispettore.md` + `sezioni-ISO/sezione-6-risorse-personale-dotazioni-domande.md` + `NC-patterns/come-rispondere-NC-ACCREDIA.md`

### BDD scenari (test L1 strutturale)

| File | Purpose |
|---|---|
| `features/profili-fetta2.feature` | 3 scenari (uno per capability) — pre-check tecnico strutturale L1, demanda L2 a Denis |
| `features/steps_test.go` (extend) | 3 nuovi step impl: `installReviewPack`, `installAuditChecklist`, `installEquipmentAlert` + step domain-specifici per le sezioni del Pack 13 (review-pack), il count delle domande (audit), le sezioni del template Equipment Alert |

### Test data fixtures (TD-1: sintetici plausibili default)

| Path | Purpose |
|---|---|
| `.pipeline/test-data/scenarios/review-pack/anno-2025-sintetico/` | NC dell'anno (CSV), esiti audit interni (md), reclami (csv), risultati PT con z-score (csv), stato apparecchiature (csv), KPI (md) — tutti sintetici plausibili |
| `.pipeline/test-data/scenarios/audit-checklist/risk-register-sintetico/` | Risk register CSV con 18 voci + descrizione scope audit (md) |
| `.pipeline/test-data/scenarios/equipment-alert/scheda-pmt-sintetica/` | Scheda apparecchiatura (md) + descrizione evento "taratura rientrata con NC su intervallo X-Y" (md) |

Se Denis consegna materiali reali anonimizzati prima dello shakedown di una specifica capability, il CTO li sostituisce capability-per-capability (decisione TD-1, vedi spec).

## Files to Modify

| File | Change | Reason |
|---|---|---|
| `.pipeline/test-data/README.md` | Aggiungere descrizione dei 3 nuovi scenari di Fetta 2 | Documentazione = codice (Layer 0) |
| `features/steps_test.go` | Aggiungere 3 step di scaffolding profilo + step di assertion strutturale (`bodyContiene13Sezioni`, `bodyContieneNDomande`, `bodyContieneTemplateEquipmentAlert`) | Estensione BDD coverage |
| `README.md` (opzionale) | Aggiornare la tabella delle 7 capability con stato "implementata" per C/D/E | Documentazione utente allineata |

## Implementation Steps (ordinati per strict ordering rischio crescente)

### Sub-fetta 2.A — `review-pack` (1 gettone, Capability C)

1. **Scrittura `profili/review-pack.yml`**
   - Schema YAML completo (vedi FR-3) con `kb_files` come sopra
   - `trigger_prompt` 5-10 righe riferito al `TEMPLATE_Management_Review_Pack` (CoWork KB sez. 15.4): istruisce il modello a produrre 13 sezioni numerate, con la sintesi iniziale coerente con le decisioni finali
   - Verifica `labnexus validate review-pack` → exit 0

2. **Costruzione test data sintetico** in `.pipeline/test-data/scenarios/review-pack/anno-2025-sintetico/`
   - Tabelle CSV plausibili: NC dell'anno (~15 righe), reclami (~5), z-score PT (~10), apparecchiature (~12), KPI (~6 indicatori)
   - File markdown narrativo: esiti audit interni (~3 audit)
   - Tutti i dati FITTIZI ma realistici (Layer 1 methodology: "Maria Rossi", non "test_user_1")

3. **Scenario BDD L1** in `features/profili-fetta2.feature`:
   - `Scenario: profilo review-pack produce il Management Review Pack a 13 sezioni`
   - Step: cartella input + lancio + frontmatter completo + body contiene 13 sezioni numerate + sintesi iniziale coerente con decisioni finali

4. **Esecuzione reale su Qwen 3 locale** (sull'hardware del CTO, fuori CI)
   - `labnexus run --profile review-pack --input <dir-sintetica> --output <out>`
   - Pre-check L1 (CTO): frontmatter OK + 13 sezioni + coerenza sintesi/decisioni — se fail → STOP, debug (kb_files? trigger_prompt? parametri Qwen?)
   - **Solo se L1 passa**: sottomettere output a Denis per L2

5. **L2 di Denis (fuori pipeline)**: registrare valutazione formale nel frontmatter (`valutazione_denis: validata | validata con riserva | non validata`)

### Sub-fetta 2.B — `audit-checklist` (1 gettone, Capability D)

6. **Scrittura `profili/audit-checklist.yml`** con `trigger_prompt` riferito all'agente AuditAgent (CoWork KB). Istruzione: 15-30 domande ispettive raggruppate per area + campioni documentali + rilievi anticipati con livello di rischio.

7. **Test data sintetico** in `.pipeline/test-data/scenarios/audit-checklist/risk-register-sintetico/`:
   - Risk register CSV (18 voci con livello rischio, area, descrizione)
   - Descrizione scope audit (md): es. "audit su §6.4 dotazioni con focus su tarature esterne"

8. **Scenario BDD L1**: lancio + body contiene tra 15 e 30 domande + raggruppamento per area + campioni documentali + livello rischio per area.

9. **Esecuzione reale + L1 + L2** (analogo a 2.A step 4-5).

### Sub-fetta 2.C — `equipment-alert` (1 gettone, Capability E)

10. **Scrittura `profili/equipment-alert.yml`** con `trigger_prompt` riferito al `TEMPLATE_Equipment_Alert` (CoWork KB sez. 15.3): stato attuale + rischio tecnico (con riferimento ai metodi specifici) + azioni proposte + bozza email fornitore + checklist al rientro.

11. **Test data sintetico** in `.pipeline/test-data/scenarios/equipment-alert/scheda-pmt-sintetica/`:
    - Scheda apparecchiatura (md): codice, descrizione, ultima taratura, prossima scadenza, fornitore, **metodi di prova associati** (almeno 3)
    - Descrizione evento (md): es. "Certificato di taratura rientrato con NC su intervallo 0-50 °C"

12. **Scenario BDD L1**: body contiene le 5 sezioni del template Equipment Alert + cita metodi specifici + rischio tecnico ragionato causalmente.

13. **Esecuzione reale + L1 + L2** (analogo).

### Gate di chiusura Fetta 2

14. Tutte e 3 le capability passano L1 (BDD) + L2 (Denis: *validata* o *validata con riserva*).
15. Aggiornare `.pipeline/test-data/README.md` con la sezione dei nuovi scenari.
16. Aggiornare README.md root con stato "implementata" per C/D/E (opzionale, gestibile in compound).
17. Procedere a `/v-review`.

## Migration Steps

**Nessuna**. Niente database; tutti gli artefatti sono filesystem.

## Test Strategy

Pattern identico a Fetta 1 (riuso skill `dual-provider-env-override` per BDD):

| Layer | Strumento | Scope Fetta 2 |
|---|---|---|
| **Acceptance (BDD)** | godog | 3 nuovi scenari in `features/profili-fetta2.feature` — uno per capability — testano L1 strutturale via fake Ollama che produce output pre-fabbricato plausibile. La validazione qualitativa L2 è dichiarata "fuori BDD" (skill `validazione-due-livelli-cto-cliente`). |
| **Integration** | n/a | Nessuna nuova integrazione di codice. Il flusso `runner.Run` esistente è già coperto. |
| **Unit** | n/a | Nessuna logica nuova di codice. |

### Test data fixtures necessarie

Vedi sezione "Files to Create" → 3 cartelle sintetiche in `.pipeline/test-data/scenarios/`. Costruite dal CTO secondo i pattern tipici di un laboratorio ISO 17025 (decisione TD-1).

### Coverage atteso

- BDD: 3 nuovi scenari L1 verdi (lancio + struttura body verificata)
- Esecuzione reale: 3 output markdown prodotti su Qwen 3 locale, sottomessi a Denis
- L2: 3 esiti formali registrati (in attesa di Denis al checkpoint)

## Risk Assessment

| Rischio | Probabilità | Impatto | Mitigazione |
|---|---|---|---|
| **`review-pack`: Qwen perde coerenza tra sezione 1 e 13** (output narrativo lungo) | Media | Alto | `trigger_prompt` esplicito sulla coerenza sintesi/decisioni; pre-check L1 verifica entrambe le sezioni; monitorare warning context (Pack lungo + KB → potrebbe sfiorare il 70%). Se Denis dice "non validata" → analisi causa (modello vs trigger vs context). |
| **`audit-checklist`: domande generiche invece di mirate al risk register** | Media | Medio | Risk register sintetico ben strutturato + trigger esplicito su "tradurre rischi specifici in domande puntuali". Denis valuta utilizzabilità operativa. |
| **`equipment-alert`: ragionamento causale debole su impatto metodi** | Bassa | Medio | Scheda apparecchiatura con "metodi associati" esplicitati + trigger su "ragionamento causale": *se questa apparecchiatura è fuori stato, allora questi metodi sono coinvolti…*. Denis valuta qualità del nesso. |
| **Dati sintetici troppo poco realistici** | Media | Medio | Pattern tipici ISO 17025 (NC, audit, PT, taratura). Se Denis al checkpoint dice "non riconosco il laboratorio", sostituire con reali anonimizzati (TD-1). |
| **Context window al 70%** in `review-pack` per via della KB pesante (sezioni-ISO + NC-patterns + mappa requisiti) | Media | Basso | `runlog` già emette warning a 70%; se passa il 100%, abort esplicito. Mitigation: selezionare solo i kb_files strettamente necessari per `review-pack` (no tutte le sezioni ISO; le più rilevanti per riesame). |
| **L2 Denis ritardato** | Media | Basso | Lo scaffolding delle 3 capability è parallelizzabile (OQ-5); l'esecuzione reale procede capability-per-capability ma in attesa di Denis si sviluppa la successiva in parallelo. |
| **Tag esatto `qwen3.6` su Ollama** | Bassa | Basso | Vedi `.pipeline/ideas/2026-05-19-auto-derive-modello-da-ollama-list.md`. Pre-volo manuale al primo `ollama list`. |

## Estimated Complexity

**Medium (2-8h)** complessivamente, distribuiti su 3 sub-fette di ~1-3h ciascuna. La parte "scrittura profilo + sintetici" è ~30min-1h per capability; la "esecuzione reale + iterazione del trigger_prompt fino al L1 passing" può richiedere 1-2h per capability se Qwen all'output strutturato richiede tuning. **L2 di Denis (fuori budget)** dipende dalla sua disponibilità.

### Distribuzione gettoni interna (`allocated_per_fetta` da state.json)

| Sotto-fetta | Gettoni | Infra/Client |
|---|---|---|
| 2.A review-pack | 1 | 0/1 (puro client) |
| 2.B audit-checklist | 1 | 0/1 |
| 2.C equipment-alert | 1 | 0/1 |
| **Totale Fetta 2** | **3** | **0/3** |

Coerente con `allocated_per_fetta["2"] = {total: 3, infra: 0, client: 3}` in `state.json`.

## Pattern di validazione (richiamo)

Tutte e 3 le capability di Fetta 2 seguono il pattern L1+L2 codificato nella skill `validazione-due-livelli-cto-cliente.md` (vedi `.pipeline/skills/project/`):

- **L1 (CTO, BDD)**: pre-check strutturale automatizzato — frontmatter completo, struttura coerente col template, no allucinazioni macroscopiche
- **L2 (Denis, manuale)**: griglia condivisa (allucinazioni, qualità ispettiva, profondità, tenuta, utilizzabilità) → esito *validata* / *validata con riserva* / *non validata*

Nessuna capability è "validata" senza L2 di Denis. Il gate `/v-deploy` di Fetta 2 si chiude quando Denis ha emesso L2 sulle 3 capability (o, in caso di non validata, dopo rinegoziazione del piano).
