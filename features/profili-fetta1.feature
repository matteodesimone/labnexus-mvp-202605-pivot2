# language: it

Funzionalità: Profili A (revisione) e B (rilievi) — shakedown del motore in Fetta 1
  In quanto CTO che valida il motore con due capability già approvate da Denis,
  voglio che ciascuno dei due profili shakedown produca un output strutturalmente coerente
  con il template della capability (Livello 1, BDD),
  così che poi Denis possa validare l'output qualitativamente (Livello 2, fuori BDD).

  Nota fondamentale (NFR-8): questi scenari validano SOLO il LIVELLO 1
  della verifica — il pre-check tecnico strutturale del CTO. Il giudizio QUALITATIVO
  è il LIVELLO 2, validazione formale di Denis, registrato manualmente nel report
  di sprint e nel frontmatter dell'output (`valutazione_denis`).

  # --- FR-13: Capability A — revisione (shakedown motore #1) ---

  Scenario: profilo revisione replica il Test 1 manuale
    Dato che il profilo "revisione" è installato in ./profili/revisione.toml
    E la cartella di input contiene PG_RISK_LAB_Rev_00.docx + DE0779_RT_08rev03.pdf + RT-08-rev.05.pdf
    Quando lancio "labnexus run --profile revisione --input <test1-input> --output <out>"
    Allora exit code è 0
    E il file di output contiene il frontmatter YAML completo
    E il body contiene almeno una sintesi iniziale e un elenco di modifiche
    E il body contiene callout "[!MODIFICA]" per ogni cambiamento normativo applicato
    E il body cita esplicitamente RT-08 rev03 e RT-08 rev05 come fonti di confronto
    E nessuna informazione inventata appare nel body (no requisiti normativi non presenti negli input)

  Scenario: confronto strutturale col golden file di Denis (Livello 1)
    Dato che esiste il golden file "PG_RISK_LAB_Rev_01_BOZZA-qwen3.6-locale.docx" approvato da Denis
    Quando confronto strutturalmente l'output di "labnexus run --profile revisione ..." col golden file
    Allora la struttura (frontmatter + sezioni di sintesi + callout MODIFICA + tono ispettivo) è paragonabile
    E il giudizio formale qualitativo del confronto è demandato a Denis (Livello 2, fuori BDD)

  # --- FR-14: Capability B — rilievi (shakedown motore #2) ---

  Scenario: profilo rilievi produce un CAPA Pack per ogni riga del CSV
    Dato che il profilo "rilievi" è installato
    E la cartella di input contiene "ACIAA A1 rilievi.csv" con N righe (NC, Osservazioni, Commenti)
    Quando lancio "labnexus run --profile rilievi --input <test2-input> --output <out>"
    Allora exit code è 0
    E il body contiene N blocchi CAPA Pack distinti
    E ogni blocco riporta meccanismo, estensione, efficacia
    E per ogni rilievo è esplicitato se basta correzione o se serve azione correttiva
    E il formato segue il template "CAPA Pack" del CoWork KB sezione 15.2
    E il giudizio formale qualitativo è demandato a Denis (Livello 2, fuori BDD)
