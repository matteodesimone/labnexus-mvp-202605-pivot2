# Skill — Validazione a due livelli (CTO + cliente/esperto)

Pattern: separare il pre-check tecnico automatico (CTO, in CI) dal giudizio qualitativo formale del committente/esperto di dominio (fuori CI). Emerso da Sprint 1 LabNexus NFR-8 ma applicabile a qualunque progetto dove l'output finale richiede un giudizio umano qualificato.

## Quando applicarlo

- Il deliverable produce output che un esperto di dominio (cliente, QM, revisore) deve giudicare qualitativamente
- Il giudizio "il prodotto funziona" richiede conoscenza che il CTO non ha o non può automatizzare
- Esempi tipici: report ispettivi, contenuti generativi LLM, traduzioni, sintesi mediche/legali, output di simulazioni che vanno interpretati

## Pattern

**Livello 1 — Pre-check tecnico CTO (in CI, BDD acceptance)**:
- Verifica oggettiva strutturale: frontmatter completo, sezioni attese presenti, marker richiesti, no allucinazioni macroscopiche rispetto a input verificabili
- Se L1 fallisce: STOP, debug pipeline. **L'output NON va sottomesso all'esperto** (è perdere il suo tempo su un bug del software)
- Implementabile come scenari BDD: `Allora il body contiene callout "[!MODIFICA]"`, `il body cita esplicitamente <fonte>`, `nessuna informazione inventata appare`

**Livello 2 — Validazione formale esperto (manuale, fuori CI)**:
- Solo dopo che L1 passa
- Giudizio qualitativo secondo griglia condivisa (allucinazioni, qualità, profondità, tenuta, utilizzabilità pratica)
- Esito formale tracciato: *validata* / *validata con riserva* / *non validata*
- Registrato nel report finale di sprint **e** nel frontmatter dell'output (es. `valutazione_<esperto>: validata`)

## Implementazione

### Spec / requirements

In `02-project.md` o equivalente, codificare la quality metric come:
```
NFR-X — Quality metric (validazione due livelli):
  L1 (CTO, automatico): <criteri strutturali enumerati>
  L2 (esperto, manuale): <griglia condivisa con esperto, registrata nel report>
```

### BDD acceptance

Nel feature file della capability:
```gherkin
Nota fondamentale: questi scenari validano SOLO il LIVELLO 1 strutturale.
Il giudizio QUALITATIVO è il LIVELLO 2, validazione formale dell'esperto,
fuori dal perimetro automatizzabile da godog/pytest-bdd/cucumber.
```

Gli step che corrispondono a "il giudizio formale qualitativo è demandato a <esperto>" sono **no-op intenzionali** (ritornano `nil`) — segnalano la separazione di responsabilità.

### Output artefatto

Includi nel frontmatter un campo per la valutazione esperto (default vuoto, riempito manualmente):
```yaml
---
valutazione_<esperto>: ""  # validata | validata con riserva | non validata
---
```

### Gate di chiusura

Una capability è **validata** SOLO dopo che l'esperto ha emesso L2. Il `/v-deploy` chiude il gate tecnico; la "validazione vera" del prodotto resta nel report di sprint, registrata dopo il checkpoint con l'esperto.

## Anti-pattern da evitare

- ❌ Pretendere di "validare il modello LLM" via BDD: impossibile senza giudizio umano
- ❌ Skippare L1 e mandare subito a L2: l'esperto rifiuta un output strutturalmente sbagliato perdendo tempo
- ❌ Mescolare L1 e L2 nello stesso assertion: chi legge il test non capisce dove finisce l'oggettività
- ❌ Fare L2 senza griglia esplicita: il giudizio diventa soggettivo, irriproducibile, conflittuale al checkpoint

## Esempio in LabNexus Sprint 1

- L1 per `revisione`: BDD verifica frontmatter YAML completo + callout `[!MODIFICA]` + citazione esplicita RT-08 rev03/rev05 + no allucinazioni
- L2 per `revisione`: Denis confronta qualitativamente con golden file `PG_RISK_LAB_Rev_01_BOZZA-qwen3.6-locale.docx` e emette esito secondo griglia condivisa
- L1 per `rilievi`: BDD verifica struttura CAPA Pack (meccanismo/estensione/efficacia, correzione vs azione correttiva)
- L2 per `rilievi`: Denis valuta utilizzabilità operativa (no golden file disponibile, vedi TD-2)

## Vedi anche

- `.pipeline/solutions/2026-05-19-fetta1-scooter.md` — emersione del pattern
- Spec NFR-8 originale
- `features/profili-fetta1.feature` — esempio di L1 in Gherkin
