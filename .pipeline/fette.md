# Fette di Sprint 1 — LabNexus motore + 7 profili + prompt meta

**Slicing approvato dal CTO**: 3 fette equilibrate, 12 gettoni residui di sviluppo pipeline.
**Strict ordering vincolante**: A → B → C → D → E → F → G → meta (vedi OQ-5 risolto: solo "esecuzione + valutazione di Denis" è strict, lo scaffolding può procedere parallelo).

| # | Nome | Gettoni | Feat coperti | Cosa funziona alla fine | Stato |
|---|---|---|---|---|---|
| 1 | Scooter — motore + shakedown | 6 | `feat-001` + `feat-002` + `feat-003` | L'utente esegue `labnexus run --profile {revisione, rilievi} ...` end-to-end. Il motore è validato contro entrambi i golden file (Test 1 PG_RISK_LAB, Test 2 ACIAA A1). | in corso |
| 2 | 3 capability nuove a rischio basso/medio | 3 | `feat-004` + `feat-005` + `feat-006` | Disponibili: `review-pack`, `audit-checklist`, `equipment-alert`. Denis può valutarle. | pianificata |
| 3 | 2 capability rischiose + prompt meta | 3 | `feat-007` + `feat-008` + `feat-meta` | Disponibili: `competence-gap`, `pt-analysis`. Inoltre `meta-prompt-genera-profilo.md` consente a Denis di esplorare nuove capability in autonomia. | pianificata |

**Totale**: 6 + 3 + 3 = 12 gettoni di sviluppo pipeline.

## Razionale dello slicing

- **Fetta 1 (scooter ricco, 6g)**: la "user journey" minimale dello sprint richiede di **validare il motore prima di costruire capability nuove**. Lo shakedown #1 (revisione, Test 1) e #2 (rilievi, Test 2) sono entrambi shakedown del motore, non valutazione delle capability A/B (già validate sul modello in chat manuale). Una fetta che includa solo Capability A lascerebbe il motore non completamente verificato. Capability B aggiunge 1 gettone e chiude il gate "motore validato definitivamente" prima di passare al rischio modello.
- **Fetta 2 (3 capability nuove, 3g)**: tre capability nuove con rischio basso/medio (testo, struttura narrativa, ragionamento causale qualitativo). La granularità "1 fetta per 1 gettone" sarebbe overhead pipeline ingiustificato; tre profili nuovi in un solo ciclo riducono il costo di gate senza perdere la possibilità di fermarsi (se Denis valuta "non validata" la prima, ridiscutiamo prima delle altre).
- **Fetta 3 (rischio alto + meta, 3g)**: `competence-gap` richiede lettura di griglie strutturate (rischio medio), `pt-analysis` è la capability più rischiosa per Qwen (ragionamento numerico — vedi OQ-8: niente pre-processing deterministico, l'esito negativo è dato utile per Sprint 2). Il `prompt meta` chiude lo sprint perché astrae il pattern emerso dai 7 profili reali.

## Pattern di validazione (vale per tutte le capability A-G)

Ogni esecuzione di profilo passa attraverso **due livelli** di validazione (vedi NFR-8 nella spec):

- **L1 — Pre-check tecnico CTO (strutturale)**: il CTO verifica che il software abbia prodotto un output ben formato (frontmatter completo, struttura coerente col template della capability, no allucinazioni macroscopiche). Se L1 fallisce → STOP, debug pipeline. **L'output non va a Denis se L1 non passa**: è perdere tempo del cliente su un bug del codice.
- **L2 — Validazione formale Denis (qualitativa)**: solo dopo L1 OK, l'output è sottomesso a Denis che giudica secondo griglia condivisa. Esito *validata* / *validata con riserva* / *non validata* registrato nel report di sprint.

Una capability è **validata** solo dopo che Denis l'ha approvata al Livello 2. Le Fette 1, 2, 3 condividono questo pattern: ogni profilo dentro una fetta passa per L1+L2.

## Checkpoint cliente

- **Fine Fetta 1**: L2 di Denis su `revisione` e `rilievi` — motore validato definitivamente? Procediamo con capability nuove?
- **Fine Fetta 2**: L2 di Denis su `review-pack`, `audit-checklist`, `equipment-alert`. Se qualcuna è *non validata*, rivalutazione del piano della Fetta 3.
- **Fine Fetta 3**: L2 di Denis su `competence-gap` e `pt-analysis` (la più rischiosa, vedi OQ-8) + walkthrough finale + report di sprint.

## Possibilità di ricalibrazione

Se a fine Fetta 1 emerge che il motore richiede più lavoro del previsto (es. parsing PDF complessi, streaming, TUI cross-platform), il CTO presenta il trade-off in gettoni e il cliente decide priorità (es. taglio di una capability di Fetta 2 o 3, o spostamento del meta a Sprint 2). Stabilità di scope (NFR-9) resta valida: niente nuovi rami senza dati.
