# language: it

Funzionalità: Profili F/G (competence-gap, pt-analysis) — Fetta 3 (FR-18/19)
  In quanto CTO che valida le ultime 2 capability nuove dello Sprint 1 su Qwen 3 locale,
  voglio che ciascun profilo produca un output strutturalmente coerente
  col template della capability (Livello 1, BDD),
  così che poi Denis possa validare l'output qualitativamente (Livello 2, fuori BDD).

  Nota fondamentale (NFR-8): questi scenari validano SOLO il LIVELLO 1
  della verifica — pre-check tecnico strutturale via fake Ollama con output
  pre-fabbricato plausibile. Il giudizio QUALITATIVO è il LIVELLO 2 di Denis,
  registrato manualmente nel report di sprint e nel frontmatter (`valutazione_denis`).

  Nota su pt-analysis: l'esito L2 "non validata" è esito ACCETTABILE per
  Sprint 1 (OQ-8 risolta, vedi spec). La parte numerica è la sfida hard del
  Qwen locale; il limite del modello è informazione utile per Sprint 2.

  # --- FR-18: Capability F — competence-gap (gap report + piano + autorizzazioni) ---

  Scenario: profilo competence-gap produce gap report + piano formazione + autorizzazioni da aggiornare
    Dato che il profilo "competence-gap" è installato dal file di progetto
    E la cartella di input contiene la matrice competenze sintetica e la descrizione di procedura nuova
    Quando lancio "labnexus run --profile competence-gap --input <dir> --output <out>"
    Allora exit code è 0
    E il file di output contiene il frontmatter YAML completo
    E il body contiene le 4 sezioni del Competence Gap Report
    E il body cita SOLO nomi tecnici presenti nella matrice
    E il body contiene un piano formazione con tempistiche
    E il body contiene una sezione autorizzazioni da rilasciare o aggiornare
    E il frontmatter contiene i default del profilo "competence-gap"
    E il giudizio formale qualitativo è demandato a Denis (Livello 2, fuori BDD)

  # --- FR-19: Capability G — pt-analysis (Investigation Pack 4 sezioni) ---

  Scenario: profilo pt-analysis produce un Investigation Pack a 4 sezioni con classificazione rischio per metodo
    Dato che il profilo "pt-analysis" è installato dal file di progetto
    E la cartella di input contiene risultati PT 2025, elenco metodi e storico 2023-2024
    Quando lancio "labnexus run --profile pt-analysis --input <dir> --output <out>"
    Allora exit code è 0
    E il file di output contiene il frontmatter YAML completo
    E il body contiene 4 sezioni numerate del PT Investigation Pack
    E il body classifica il rischio per metodo (alto, medio, basso)
    E il body cita SOLO codici PT presenti nei risultati di input
    E il body cita SOLO codici metodo presenti nell'elenco metodi
    E il body include azioni immediate per i casi con z-score fuori soglia
    E il frontmatter contiene i default del profilo "pt-analysis"
    E il giudizio formale qualitativo è demandato a Denis (Livello 2, fuori BDD)
