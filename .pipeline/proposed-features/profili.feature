# language: it

Funzionalità: I sette profili LabNexus — capability ispettive A-G
  In quanto CTO che valida una capability alla volta,
  voglio che ciascun profilo, lanciato sull'input previsto,
  produca un output coerente col template descritto nel CoWork KB di Denis,
  così che Denis possa giudicarlo secondo la griglia condivisa
  (allucinazioni, qualità ispettiva, profondità, tenuta, utilizzabilità pratica).

  Nota fondamentale (NFR-8): questi scenari BDD validano SOLO il LIVELLO 1
  della verifica — il pre-check tecnico strutturale del CTO (output esiste,
  frontmatter completo, struttura coerente col template della capability,
  no allucinazioni macroscopiche). Il giudizio QUALITATIVO di ciascun output
  (è davvero utilizzabile? è ispettivo? regge in profondità?) è il LIVELLO 2,
  validazione formale di Denis, fuori dal perimetro automatizzabile da godog.
  L'esito di Livello 2 viene registrato manualmente nel report di sprint e
  nel frontmatter dell'output (campo `valutazione_denis`).

  Una capability si dice "validata" SOLO dopo che Denis ha emesso il giudizio
  formale al Livello 2. Nessuno scenario BDD può sostituirsi a quel giudizio.

  # --- FR-13: Capability A — revisione (shakedown motore #1) ---

  Scenario: profilo revisione replica il Test 1 manuale
    Dato che il profilo "revisione" è installato in ./profili/revisione.yml
    E la cartella di input contiene PG_RISK_LAB_Rev_00.docx + DE0779_RT_08rev03.pdf + RT-08-rev.05.pdf
    Quando lancio "labnexus run --profile revisione --input <test1-input> --output <out>"
    Allora exit code è 0
    E il file di output contiene il frontmatter YAML completo
    E il body contiene almeno una sintesi iniziale e un elenco di modifiche
    E il body contiene callout "[!MODIFICA]" per ogni cambiamento normativo applicato
    E il body cita esplicitamente RT-08 rev03 e RT-08 rev05 come fonti di confronto
    E nessuna informazione inventata appare nel body (no requisiti normativi non presenti negli input)

  Scenario: confronto qualitativo col golden file di Denis
    Dato che esiste il golden file "PG_RISK_LAB_Rev_01_BOZZA-qwen3.6-locale.docx" approvato da Denis
    Quando confronto l'output di "labnexus run --profile revisione ..." col golden file
    Allora la struttura (frontmatter + sezioni di sintesi + callout MODIFICA + tono ispettivo) è paragonabile
    E il giudizio formale del confronto è registrato nel report di sprint

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

  # --- FR-15: Capability C — review-pack (Management Review Pack annuale) ---

  Scenario: review-pack produce le 13 sezioni del Pack
    Dato che il profilo "review-pack" è installato
    E la cartella di input contiene dati strutturati di un anno operativo (NC, audit interni, reclami, PT con z-score, apparecchiature, KPI)
    Quando lancio "labnexus run --profile review-pack --input <dir> --output <out>"
    Allora exit code è 0
    E il body contiene esattamente 13 sezioni numerate del template Management Review Pack (CoWork KB 15.4)
    E la sintesi iniziale è coerente con le decisioni finali della direzione
    E nessun dato è inventato rispetto ai file di input

  # --- FR-16: Capability D — audit-checklist da risk register ---

  Scenario: audit-checklist genera domande ispettive raggruppate per area
    Dato che il profilo "audit-checklist" è installato
    E la cartella di input contiene un risk register (15-25 voci) e una descrizione di scope d'audit
    Quando lancio "labnexus run --profile audit-checklist --input <dir> --output <out>"
    Allora exit code è 0
    E il body contiene tra 15 e 30 domande ispettive
    E le domande sono raggruppate per area
    E per ogni gruppo è elencato almeno un "campione documentale da verificare"
    E per ogni area è indicato il livello di rischio dei rilievi anticipati

  # --- FR-17: Capability E — equipment-alert ---

  Scenario: equipment-alert ragiona su impatto della NC sui metodi associati
    Dato che il profilo "equipment-alert" è installato
    E la cartella di input contiene una scheda apparecchiatura con metodi associati e una descrizione di evento (certificato taratura rientrato con NC)
    Quando lancio "labnexus run --profile equipment-alert --input <dir> --output <out>"
    Allora exit code è 0
    E il body contiene le sezioni del template Equipment Alert (CoWork KB 15.3)
    E è presente la sezione "stato attuale"
    E è presente la sezione "rischio tecnico" che cita metodi specifici
    E è presente la sezione "azioni proposte"
    E è presente una "bozza email fornitore"
    E è presente una "checklist al rientro"

  # --- FR-18: Capability F — competence-gap ---

  Scenario: competence-gap analizza la matrice e produce il gap report
    Dato che il profilo "competence-gap" è installato
    E la cartella di input contiene una matrice competenze (10-15 righe × 8-12 colonne) e una descrizione di procedura/metodo nuovo
    Quando lancio "labnexus run --profile competence-gap --input <dir> --output <out>"
    Allora exit code è 0
    E il body riporta chi è autorizzato e chi no per la procedura/metodo nuovo
    E il body propone un piano di formazione concreto
    E il body elenca le autorizzazioni da rilasciare o aggiornare

  # --- FR-19: Capability G — pt-analysis (rischio Qwen sui numeri) ---

  Scenario: pt-analysis produce l'Investigation Pack
    Dato che il profilo "pt-analysis" è installato
    E la cartella di input contiene risultati di una campagna PT con z-score, elenco metodi, e idealmente storico 2-3 anni
    Quando lancio "labnexus run --profile pt-analysis --input <dir> --output <out>"
    Allora exit code è 0
    E il body contiene un riassunto risultati
    E il body classifica il rischio per metodo
    E il body riporta il trend storico se i dati storici sono presenti
    E il body propone azioni immediate per ogni z-score fuori soglia

  Scenario: pt-analysis con errori numerici di Qwen è ATTESO
    Dato che Qwen potrebbe sbagliare calcoli numerici (debolezza nota degli LLM)
    Quando lancio pt-analysis e il modello produce un calcolo z-score inaccurato
    Allora il software NON aggiunge pre-processing deterministico (vedi OOS-1 e OQ-8)
    E il giudizio sull'usabilità della capability spetta a Denis
    E un esito "non validata" su pt-analysis è informazione utile per Sprint 2, non un bug del codice
