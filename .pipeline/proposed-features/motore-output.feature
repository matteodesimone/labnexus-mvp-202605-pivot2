# language: it

Funzionalità: Motore LabNexus — output markdown con frontmatter tracciabile
  In quanto CTO che raccoglie gli output per la valutazione di Denis,
  voglio che ogni esecuzione produca un singolo file markdown con metadati completi
  così da poter riconciliare output e parametri usati senza ambiguità.

  # --- FR-9: output markdown + frontmatter ---

  Scenario: file di output viene creato con naming convenzionale
    Quando lancio "labnexus run --profile revisione --input /tmp/test1-input --output /tmp/out" con successo
    Allora nella cartella "/tmp/out" esiste un file con pattern "<timestamp>_revisione_<descrittore>.md"
    E il timestamp ha formato ISO compatto "YYYY-MM-DDTHHMMSS"
    E exit code è 0

  Scenario: frontmatter del file di output è completo
    Quando lancio una capability "revisione" con successo
    Allora il file di output inizia con un blocco YAML "---"
    E il frontmatter contiene la chiave "profilo" con valore "revisione"
    E il frontmatter contiene la chiave "modello" con il modello effettivamente chiamato
    E il frontmatter contiene la chiave "provider" con valore "ollama" o "eurouter"
    E il frontmatter contiene la chiave "data_esecuzione" in ISO 8601
    E il frontmatter contiene la chiave "durata_secondi"
    E il frontmatter contiene la chiave "token_stimati"
    E il frontmatter contiene la chiave "file_input" con la lista dei nomi (non i contenuti)
    E dopo il frontmatter c'è il body con l'output completo del modello

  Scenario: collisione di timestamp produce suffisso incrementale
    Dato che nella cartella "/tmp/out" esiste già il file "2026-05-18T141022_revisione_doc.md"
    Quando lancio una nuova esecuzione "revisione" nello stesso secondo con lo stesso input descriptor
    Allora il nuovo file ha nome "2026-05-18T141022_revisione_doc_2.md"
    E nessun file esistente viene sovrascritto

  Scenario: streaming interrotto produce comunque un file parziale
    Dato che la connessione al provider viene interrotta dopo metà output
    Quando l'eseguibile chiude la chiamata
    Allora il file di output esiste con il contenuto parziale
    E il frontmatter contiene "stato: interrotto"
    E exit code è 1
