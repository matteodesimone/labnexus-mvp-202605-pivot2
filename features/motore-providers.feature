# language: it

Funzionalità: Motore LabNexus — provider Ollama e EUrouter, streaming, progress
  In quanto CTO che lancia capability su Qwen 3 locale (o, in debug, su EUrouter),
  voglio che lo streaming sia incrementale, la progress bar visibile,
  e il log per step diagnosticabile,
  così da non aspettare al buio e individuare problemi senza ricominciare da capo.

  # --- FR-7: provider Ollama (NDJSON) ---

  Scenario: streaming Ollama produce token incrementali
    Dato che ollama serve è in ascolto su localhost:11434
    E "./profili/revisione.toml" dichiara provider "ollama" e modello "qwen3:14b-q5_K_M"
    E la cartella di input contiene un input minimale di test
    Quando lancio "labnexus run --profile revisione --input /tmp/in --output /tmp/out"
    Allora il primo token visibile in stdout arriva entro 5 secondi
    E il flusso di token continua incrementalmente fino al termine
    E exit code è 0

  Scenario: Ollama non disponibile fallisce con istruzione operativa
    Dato che ollama serve NON è in ascolto su localhost:11434
    Quando lancio "labnexus run --profile revisione --input /tmp/in --output /tmp/out"
    Allora exit code è 1
    E stderr contiene "Ollama"
    E stderr suggerisce di "verifica che 'ollama serve' sia attivo"

  # --- FR-7: provider EUrouter (SSE) ---

  Scenario: streaming EUrouter rispetta SSE OpenAI-compatible
    Dato che "./profili/debug-eurouter.toml" dichiara provider "eurouter" e modello "mistralai/mistral-large"
    E l'env var "EUROUTER_API_KEY" è impostata a un valore valido (test)
    E la cartella di input contiene un input sintetico
    Quando lancio "labnexus run --profile debug-eurouter --input /tmp/in --output /tmp/out"
    Allora il client invia una request a "https://api.eurouter.ai/v1/chat/completions"
    E i chunk SSE ricevuti hanno la forma "data: {...}"
    E il chunk finale è "data: [DONE]"
    E exit code è 0

  # --- FR-8: progress bar e log per step ---

  Scenario: log per step con tempi
    Quando lancio "labnexus run --profile revisione --input /tmp/in --output /tmp/out"
    Allora stderr contiene una riga per ognuno degli step "caricamento KB", "parsing input", "stima context", "chiamata provider", "streaming", "scrittura output"
    E ciascuna riga riporta la durata dello step in millisecondi

  Scenario: progress bar visibile durante streaming
    Quando lancio "labnexus run --profile revisione --input /tmp/in --output /tmp/out" in modalità TTY
    Allora durante lo streaming stderr mostra una progress bar che avanza
    E la progress bar scompare al termine dello streaming

  Scenario: nessuna progress bar in modalità non-TTY (piping)
    Quando lancio "labnexus run --profile revisione --input /tmp/in --output /tmp/out 2>&1 | cat"
    Allora stderr non contiene caratteri ANSI di progress bar
    E stderr contiene comunque i log per step
