# LabNexus / AICertus — Materiale di bootstrap progetto

**Per**: assistente AI di sviluppo (framework vobbly / Prove Out)
**Da**: Matteo De Simone (CTO)
**Data**: {{data}}

## Come leggere questa cartella

Questo è materiale di bootstrap per un progetto che **non parte da zero**. Lo Sprint 1 è già avviato, ha subito un assestamento iniziale e un pivot di scope. Una parte significativa del budget gettoni è già stata spesa. Il piano corrente è frutto di decisioni consensuali con il cliente.

**Ordine di lettura suggerito**:

1. `00-stato-e-pivot.md` — stato del progetto, pivot avvenuti, gettoni spesi e residui. Indispensabile per capire dove siamo.
2. `01-piano-esperimenti.md` — i 7 esperimenti che lo Sprint 1 deve condurre e materiali necessari.
3. `02-deliverable-software.md` — cosa l'eseguibile fa, come si usa, cosa consegniamo a fine sprint.
4. `03-brief-tecnico.md` — il brief operativo Prove Out completo, con Discovery Goal espanso, note feature, schema YAML.
5. `materiali-dominio/` — KB del cliente, materiali di visione, golden file dei test precedenti. **Vedi il README dentro per capire come trattare ciascun gruppo.**

## Cosa devi fare con questo materiale

Decidi tu dove smistare i pezzi nella struttura del framework (`.pipeline/standards/`, `.pipeline/ideas/`, `.pipeline/skills/`, ecc.). Le possibilità a grandi linee:

- Lo **stato del progetto e i pivot** (file 00) vanno nella sezione Project Context di `02-project.md` per essere letti a ogni sessione. È il contesto stabile.
- Il **brief tecnico** (file 03) è la fonte primaria per `/v-brainstorm` o `/v-spec`. Tipicamente va in `.pipeline/ideas/`.
- I documenti **piano esperimenti** e **deliverable software** (file 01 e 02) sono input strategici per fette successive. Decidi tu se metterli in ideas o riassumerli in project context.
- I **materiali di dominio** (cartella `materiali-dominio/`) sono contesto di conoscenza del cliente. Hanno un proprio README che spiega come trattarli.

## Vincoli operativi non negoziabili

Prima di proporre il piano di sviluppo:

- **Tempo e budget sono fissi.** Restano 12 gettoni dichiarati al cliente; il CTO assorbe internamente 3 gettoni di sforo aggiuntivi (il piano interno è 15 ma comunicabile come 12 con sforo dichiarato di 1). Vedi file 00 per il dettaglio.
- **Strategia di rischio crescente sui profili.** L'ordine di sviluppo dei 7 profili è vincolato: prima i 2 profili che replicano test già approvati (per validare il motore), poi le capability nuove in ordine crescente di rischio. Vedi file 03.
- **KB di Denis vs trigger nostro.** I file della KB-ispettore sono scritti da Denis e usati come system context. L'unico testo scritto da noi per profilo è il `trigger_prompt`. Vedi file 03.
- **Locale mandatory.** Tutta l'elaborazione AI su Qwen via Ollama in locale. EUrouter come provider alternativo è facility tecnica di debug, non canale operativo sui dati reali.

## Procedere

Leggi tutto il materiale. Proponi il piano di sviluppo secondo il flusso Prove Out del framework. **Non scrivere codice finché il CTO non ha approvato il piano.**

Se durante la lettura emergono ambiguità o cose non chiare, fai domande puntuali al CTO **prima** di proporre il piano. Meglio una domanda in più adesso che un piano da rifare dopo.
