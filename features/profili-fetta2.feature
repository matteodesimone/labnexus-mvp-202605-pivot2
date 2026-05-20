# language: it

Funzionalità: Profili C/D/E (review-pack, audit-checklist, equipment-alert) — Fetta 2 (FR-15/16/17)
  In quanto CTO che valida 3 nuove capability dichiarative su Qwen 3 locale,
  voglio che ciascun profilo produca un output strutturalmente coerente
  col template della capability (Livello 1, BDD),
  così che poi Denis possa validare l'output qualitativamente (Livello 2, fuori BDD).

  Nota fondamentale (NFR-8, identica a Fetta 1): questi scenari validano SOLO il LIVELLO 1
  della verifica — il pre-check tecnico strutturale del CTO via fake Ollama con output
  pre-fabbricato plausibile. Il giudizio QUALITATIVO è il LIVELLO 2 di Denis,
  registrato manualmente nel report di sprint e nel frontmatter (`valutazione_denis`).

  # --- FR-15: Capability C — review-pack (Management Review Pack, 13 sezioni) ---

  Scenario: profilo review-pack produce il Management Review Pack a 13 sezioni
    Dato che il profilo "review-pack" è installato dal file di progetto
    E la cartella di input contiene il pacchetto sintetico anno-2025 per Management Review Pack
    Quando lancio "labnexus run --profile review-pack --input <dir> --output <out>"
    Allora exit code è 0
    E il file di output contiene il frontmatter YAML completo
    E il body contiene 13 sezioni numerate del Management Review Pack
    E il body contiene una sintesi iniziale coerente con le decisioni finali del Pack
    E il frontmatter contiene i default del profilo "review-pack"
    E il giudizio formale qualitativo è demandato a Denis (Livello 2, fuori BDD)

  # --- FR-16: Capability D — audit-checklist (15–30 domande raggruppate per area) ---

  Scenario: profilo audit-checklist produce una checklist tra 15 e 30 domande raggruppate per area
    Dato che il profilo "audit-checklist" è installato dal file di progetto
    E la cartella di input contiene il risk-register sintetico (18 voci) e lo scope d'audit
    Quando lancio "labnexus run --profile audit-checklist --input <dir> --output <out>"
    Allora exit code è 0
    E il file di output contiene il frontmatter YAML completo
    E il body contiene tra 15 e 30 domande d'audit raggruppate per area
    E il body cita campioni documentali da richiedere per ciascuna area
    E il body riporta un livello di rischio per ciascuna area
    E il frontmatter contiene i default del profilo "audit-checklist"
    E il giudizio formale qualitativo è demandato a Denis (Livello 2, fuori BDD)

  # --- FR-17: Capability E — equipment-alert (5 sezioni del template + metodi specifici) ---

  Scenario: profilo equipment-alert produce un Equipment Alert con 5 sezioni e cita i metodi associati
    Dato che il profilo "equipment-alert" è installato dal file di progetto
    E la cartella di input contiene la scheda PMT sintetica e la descrizione evento taratura
    Quando lancio "labnexus run --profile equipment-alert --input <dir> --output <out>"
    Allora exit code è 0
    E il file di output contiene il frontmatter YAML completo
    E il body contiene le 5 sezioni del template Equipment Alert
    E il body cita i metodi di prova specifici associati all'apparecchiatura
    E il body include un ragionamento causale tra stato apparecchiatura e impatto metodi
    E il frontmatter contiene i default del profilo "equipment-alert"
    E il giudizio formale qualitativo è demandato a Denis (Livello 2, fuori BDD)
