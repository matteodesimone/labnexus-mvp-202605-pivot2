# language: it

Funzionalità: Sprint 1.5.C — Concierge mode (auto-discovery lavori/ + audit log + streaming live)
  Come Denis (QM laboratorio), voglio usare labnexus come concierge preliminare per il mio
  lavoro quotidiano di Quality Management: aggiungere/modificare/rimuovere cartelle di lavoro
  in lavori/ senza toccare config centrali, vedere il modello pensare in tempo reale,
  avere un audit trail completo per ogni esecuzione.

  # --- FR-25, FR-26: cartella lavori/ + _labnexus.toml per cartella ---

  Scenario: una cartella in lavori/ con _labnexus.toml è riconosciuta come job
    Dato che esiste "lavori/CAPABILITY A — Profilo revisione/_labnexus.toml" con:
      """
      profile = "revisione"
      trigger_prompt_file = "Prompt_INPUT_Rev.00.rtf"
      """
    E nella stessa cartella esistono i file di input (es. Documento_da_revisionare/, Nuovi_Requisiti/)
    Quando lancio "labnexus jobs"
    Allora stdout contiene una riga per il job "revisione" con path "lavori/CAPABILITY A — Profilo revisione"

  # --- FR-27: fallback convention naming ---

  Scenario: cartella senza _labnexus.toml ma con convention naming è riconosciuta
    Dato che esiste "lavori/Nuovo Lavoro — Profilo audit-checklist/" SENZA _labnexus.toml
    E la cartella contiene file di input
    Quando lancio "labnexus jobs"
    Allora stdout contiene il job risolto con profile "audit-checklist" (estratto da regex "Profilo (.+)$")
    E stderr contiene un warning "convention naming usato (manca _labnexus.toml)"

  Scenario: cartella senza _labnexus.toml e senza convention naming è skippata
    Dato che esiste "lavori/cartella-random/" senza _labnexus.toml e senza pattern "Profilo ..."
    Quando lancio "labnexus jobs"
    Allora stdout NON contiene quella cartella nella lista jobs
    E stderr contiene un warning "lavori/cartella-random skipped (no _labnexus.toml, no convention match)"

  # --- FR-28: labnexus jobs subcommand ---

  Scenario: labnexus jobs lista i jobs auto-discovery
    Dato che lavori/ contiene 7 cartelle con _labnexus.toml valido
    Quando lancio "labnexus jobs"
    Allora stdout elenca i 7 jobs con nome + profile + input_dir
    E exit code è 0

  Scenario: labnexus jobs --json produce output strutturato
    Dato che lavori/ contiene 7 cartelle
    Quando lancio "labnexus jobs --json"
    Allora stdout è JSON parsabile
    E contiene un array di 7 oggetti con campi {"name", "profile", "input_dir"}

  Scenario: lavori/ vuota o assente → labnexus jobs ritorna lista vuota senza crash
    Dato che lavori/ NON esiste o è vuota
    Quando lancio "labnexus jobs"
    Allora stdout contiene "Nessun lavoro trovato in lavori/"
    E exit code è 0 (non è un errore)

  # --- FR-29: labnexus run --job <name> ---

  Scenario: labnexus run --job esegue headless via auto-discovery
    Dato che lavori/ contiene un job auto-discovery "revisione" con _labnexus.toml e profile "revisione"
    Quando lancio "labnexus run --job revisione"
    Allora exit code è 0
    E nessuna TUI è aperta (headless)
    E output va a "lavori/CAPABILITY A — Profilo revisione/output/<timestamp>_revisione_*.md"

  Scenario: labnexus run --profile/--input/--output bypass auto-discovery (Sprint 1 invariato)
    Dato che lavori/ ha jobs auto-discovery
    Quando lancio "labnexus run --profile revisione --input /custom/path --output /custom/out"
    Allora il sistema NON guarda lavori/ (path espliciti hanno precedenza)
    E exit code è 0

  # --- FR-30: output dentro lavori/<X>/output/ ---

  Scenario: output di default va dentro la cartella del job
    Dato che lavori/CAPABILITY-X contiene un job
    Quando lancio "labnexus run --job <X>"
    Allora viene creata la sottocartella "lavori/CAPABILITY-X/output/" se non esiste
    E i file .md e .log sono scritti lì
    E NON sono toccate cartelle al di fuori di lavori/CAPABILITY-X/

  # --- FR-31: TUI estesa con menu auto-discovery ---

  Scenario: TUI mostra menu pre-mappato dei jobs auto-discovery
    Dato che lavori/ contiene 5 jobs e il binary è invocato senza args in TTY
    Quando lancio "labnexus"
    Allora la TUI si apre con un menu che mostra i 5 jobs come opzioni selezionabili
    E ciascuna voce mostra nome + profile + path
    E Denis può scegliere uno e il sistema esegue come "run --job <selezione>"

  # --- FR-32, FR-33: logger esteso multi-writer + step granulari ---

  Scenario: log file viene generato accanto al .md per ogni run
    Dato che lavori/<X>/ contiene un job
    Quando lancio "labnexus run --job <X>"
    Allora viene generato sia "<timestamp>_<profile>_*.md" SIA "<timestamp>_<profile>_*.log"
    E i due file hanno lo stesso prefisso di naming
    E sono nella stessa cartella "lavori/<X>/output/"

  Scenario: log file contiene i dettagli per step
    Dato che ho appena eseguito un run con successo
    Quando apro il .log file generato
    Allora contiene linee con [HH:MM:SS] per ogni step (caricamento profilo, parsing input, ecc.)
    E contiene linee con "file: <name> (<dim> bytes, ~<token> token)"
    E contiene linee con "kb_files caricati: <N>"
    E contiene linee con "provider chiamato: eurouter (https://api.eurouter.ai/...)"
    E contiene linee con "modello: qwen3.6, temperature: 0.9"

  # --- FR-34: frontmatter MD include log_file ---

  Scenario: frontmatter MD output include riferimento al log_file
    Dato che ho appena eseguito un run
    Quando apro il .md generato
    Allora il frontmatter YAML contiene una chiave "log_file" con path relativo al .log accoppiato

  # --- FR-35: body streaming live in stderr + log mirror ---

  Scenario: durante run in TTY, body streaming è visibile incrementalmente in stderr
    Dato che lancio "labnexus run --job revisione" da Terminal interattivo (TTY)
    Allora durante la chiamata provider, stderr riceve incrementalmente chunk di testo del modello
    E il log file riceve gli stessi chunk in append
    E al termine, il .md output contiene il body completo (uguale a quello accumulato dallo streaming)

  Scenario: durante run non-TTY (pipe/redirect), body streaming va solo al log file
    Dato che lancio "labnexus run --job revisione 2>&1 | cat"
    Allora stderr (catturato via pipe) NON contiene chunk live del body (per non sporcare lo stdout deterministico downstream)
    E il log file SÌ contiene il body completo
    E il .md output è generato regolarmente

  # --- FR-36: pre-check API key in labnexus.command ---

  Scenario: API key mancante → errore esplicito al doppio click .command
    Dato che "labnexus.config.toml" contiene "eurouter_api_key = \"\""
    E il provider effettivo è eurouter
    Quando faccio doppio click su "labnexus.command"
    Allora Terminal mostra un errore italiano "Apri labnexus.config.toml e inserisci la tua chiave EUROUTER_API_KEY"
    E exit code è 2
    E NON viene tentata nessuna chiamata di rete

  # --- FR-37: README zip-level italiano concierge ---

  Scenario: README.md guida Denis al setup in 3 passi
    Dato che leggo "README.md" al root del zip estratto
    Allora contiene una sezione "Inizio rapido" con 3 step numerati
    E lo step 1 menziona "estrai zip"
    E lo step 2 menziona "edit labnexus.config.toml con la API key"
    E lo step 3 menziona "doppio click su labnexus.command"

  # --- FR-38: USER-GUIDE concierge mode ---

  Scenario: USER-GUIDE descrive workflow concierge Denis
    Dato che leggo "docs/USER-GUIDE.md"
    Allora contiene una sezione su come Denis aggiunge un nuovo lavoro a lavori/
    E descrive il pattern _labnexus.toml
    E spiega come leggere il file output + valutare L2 + annotare valutazione_denis

  # --- FR-39: guida-meta-prompt tono concierge ---

  Scenario: guida-meta-prompt usa terminologia "lavoro" non "test"
    Dato che leggo "docs/guida-meta-prompt-denis.md"
    Allora contiene "aggiungi nuovi tipi di lavoro" (non "aggiungi nuovi test")
    E lo schema profilo mostrato è in TOML (non YAML)
    E menziona "trigger_prompt_file" come pattern opzionale

  # --- EC-1, EC-2, EC-3, EC-4: edge cases auto-discovery ---

  Scenario: profile referenziato in _labnexus.toml non esiste
    Dato che "lavori/X/_labnexus.toml" dichiara profile = "inesistente"
    Quando lancio "labnexus jobs"
    Allora stdout lista il job con warning "profile 'inesistente' non trovato"
    E stderr contiene il warning

  Scenario: trigger_prompt_file referenziato non esiste nella cartella input
    Dato che "lavori/X/_labnexus.toml" dichiara trigger_prompt_file = "missing.rtf"
    E "lavori/X/missing.rtf" NON esiste
    Quando lancio "labnexus run --job X"
    Allora exit code è 2
    E stderr contiene "trigger_prompt_file 'missing.rtf' non trovato in lavori/X/"
