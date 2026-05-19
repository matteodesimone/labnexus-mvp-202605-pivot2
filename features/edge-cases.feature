# language: it

Funzionalità: Edge case del motore LabNexus
  In quanto CTO che esegue capability reali su input reali,
  voglio che gli errori e le situazioni limite siano esplicite, recuperabili o documentate,
  così da non perdere tempo a indovinare cosa è successo.

  # --- EC-1: PDF complesso che spacca il parser ---

  Scenario: PDF non parsabile viene skippato con warning, esecuzione prosegue
    Dato che la cartella di input contiene "documento_buono.pdf" e "documento_rotto.pdf"
    E "documento_rotto.pdf" fa fallire la libreria di parsing
    Quando lancio "labnexus run --profile revisione --input <dir> --output <out>"
    Allora exit code è 0
    E stderr contiene "documento_rotto.pdf"
    E stderr suggerisce "converti manualmente con pandoc"
    E l'output viene generato a partire dal solo "documento_buono.pdf"

  Scenario: oltre il 50% dei file fallisce il parsing → errore esplicito
    Dato che la cartella di input contiene 4 PDF e 3 di essi fanno fallire il parser
    Quando lancio "labnexus run ..."
    Allora exit code è 1
    E stderr contiene "oltre il 50%"
    E nessuna chiamata al provider è stata effettuata

  # --- EC-3: context oltre il 100% ---

  Scenario: context al 100% blocca l'esecuzione
    Dato che il prompt composto eccede il context_window dichiarato nel profilo
    Quando lancio "labnexus run ..."
    Allora exit code è 2
    E stderr contiene "context window"
    E stderr suggerisce di ridurre i file di input o aumentare context_window nel profilo

  # --- EC-7: Ollama down ---

  Scenario: Ollama non raggiungibile produce errore esplicito con istruzione
    Dato che ollama serve non è in ascolto
    Quando lancio "labnexus run --profile revisione ..."
    Allora exit code è 1
    E stderr contiene "Ollama non raggiungibile su localhost:11434"
    E stderr suggerisce "ollama serve"

  # --- EC-8: streaming chiuso senza done marker (rivisto dal bugfix
  # `eof-post-content-falsamente-interrotto` del 2026-05-19) ---

  Scenario: EOF post-content interpretato come "completato (no done marker)"
    Dato che la connessione cade dopo che il provider ha emesso ~50% dei token
    Quando il client gestisce l'interruzione
    Allora il file di output esiste con il contenuto parziale
    E il frontmatter riporta "stato: completato"
    E stderr contiene "stream chiuso senza done marker"

  # --- EC-9: input vuoto o solo formati non supportati ---

  Scenario: cartella di input vuota
    Dato che la cartella di input non contiene alcun file
    Quando lancio "labnexus run ..."
    Allora exit code è 2
    E stderr contiene "cartella di input vuota"

  Scenario: cartella di input con solo formati non supportati
    Dato che la cartella di input contiene solo file .jpg e .mov
    Quando lancio "labnexus run ..."
    Allora exit code è 2
    E stderr elenca i formati supportati

  # --- EC-13: chunk NDJSON Ollama malformato ---

  Scenario: chunk NDJSON malformato isolato viene skippato
    Dato che durante lo streaming Ollama emette 1 chunk NDJSON malformato seguito da chunk validi
    Quando il parser legge lo stream
    Allora il chunk malformato è skippato con warning
    E lo streaming prosegue regolarmente
    E exit code è 0

  Scenario: ≥ 10 chunk NDJSON malformati consecutivi → errore
    Dato che lo streaming Ollama emette 10 chunk malformati consecutivi
    Allora il client interrompe lo streaming
    E exit code è 1
    E stderr contiene "streaming Ollama corrotto"

  # --- EC-14: SSE [DONE] mai ricevuto ---

  Scenario: timeout SSE EUrouter senza terminatore [DONE]
    Dato che EUrouter non emette mai "data: [DONE]" entro il timeout configurato
    Quando il client raggiunge il timeout
    Allora il file di output esiste con il contenuto parziale
    E il frontmatter contiene "stato: timeout"
    E exit code è 1

  # --- EC-15: collisione di timestamp output ---

  @manual
  Scenario: due esecuzioni nello stesso secondo non si sovrascrivono
    Dato che nella cartella di output esiste già "2026-05-18T141022_revisione_doc.md"
    Quando lancio una seconda esecuzione che produrrebbe lo stesso nome
    Allora il nuovo file ha suffisso "_2"
    E il file esistente è intatto

  # --- EC-11: meta-prompt che produce file KB inventati ---

  Scenario: validate intercetta YAML del meta-prompt con kb_files inventati
    Dato che il meta-prompt ha (per errore) prodotto uno YAML che referenzia "fantasma-ISO-99999.md"
    Quando lancio "labnexus validate <nome>"
    Allora exit code è 2
    E stderr contiene "fantasma-ISO-99999.md"
    E stderr contiene "non esiste in ./KB-ispettore/"
