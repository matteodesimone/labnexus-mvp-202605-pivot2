# language: it

Funzionalità: Sprint 1.5.B — Config master TOML + profile-override + parser nuovi formati
  Come CTO che deve permettere a Denis di configurare labnexus senza editare 7 profili separati,
  voglio un file di configurazione master TOML che fornisce defaults globali (provider, modello,
  API key, ecc.) ereditati dai profili,
  così che Denis edita un solo file per cambiare comportamento globale, e i profili restano snelli.
  In più, devo supportare i formati di file Office reali che Denis usa nel suo lavoro
  (.xls, .doc, .rtf oltre a quelli già supportati).

  # --- FR-11, FR-12: config master TOML ---

  Scenario: labnexus.config.toml master fornisce defaults globali
    Dato che esiste "labnexus.config.toml" al root con:
      """
      provider = "eurouter"
      modello = "qwen3.6"
      temperature = 0.9
      max_tokens = 8192
      context_window = 128000
      eurouter_api_key = "sk-test-key"
      """
    E un profilo "test-min.toml" che dichiara SOLO "profilo", "descrizione", "kb_files", "trigger_prompt"
    Quando lancio "labnexus run --profile test-min --input /tmp/in --output /tmp/out"
    Allora exit code è 0
    E il provider effettivo usato è "eurouter" (ereditato dal master)
    E il modello effettivo è "qwen3.6" (ereditato dal master)

  Scenario: profile override del master prevale
    Dato che il master dichiara "provider = \"eurouter\""
    E un profilo "test-ollama.toml" dichiara "provider = \"ollama\""
    Quando lancio "labnexus run --profile test-ollama --input /tmp/in --output /tmp/out"
    Allora il provider effettivo è "ollama" (override del profilo vince sul master)

  Scenario: formato TOML errato produce errore chiaro con file:line
    Dato che "labnexus.config.toml" contiene un errore di sintassi TOML alla linea 5
    Quando lancio "labnexus list"
    Allora exit code è 2
    E stderr contiene il messaggio di errore con riferimento a "labnexus.config.toml:5"

  # --- FR-13: campi profile opzionali ---

  Scenario: un profilo può omettere campi opzionali e li eredita dal master
    Dato che il master fornisce provider, modello, temperature, max_tokens, context_window
    E un profilo dichiara SOLO profilo + descrizione + kb_files + trigger_prompt
    Quando eseguo "labnexus describe test-profile"
    Allora stdout mostra provider, modello, temperature, ecc. ereditati dal master
    E i kb_files dichiarati nel profilo

  # --- FR-14, FR-15: pacchetto internal/config + merge logic ---

  Scenario: profile.Validate valida il profilo merged, non isolato
    Dato che il master fornisce provider = "eurouter"
    E un profilo NON dichiara provider esplicitamente
    Quando lancio "labnexus validate test-profile"
    Allora exit code è 0 (provider valido post-merge da master)

  Scenario: profile.Validate fallisce se profilo merged ha campo invalido
    Dato che il master fornisce provider = "eurouter"
    E un profilo override con provider = "openai"
    Quando lancio "labnexus validate test-profile"
    Allora exit code è 2
    E stderr contiene "provider 'openai' non ammesso"

  # --- FR-16, FR-17: trigger_prompt_file indirezione ---

  Scenario: trigger_prompt_file inline-XOR con trigger_prompt — entrambi setted → errore
    Dato che un profilo dichiara SIA "trigger_prompt = ..." SIA "trigger_prompt_file = ..."
    Quando lancio "labnexus validate test-profile"
    Allora exit code è 2
    E stderr contiene "mutually exclusive" o messaggio equivalente

  Scenario: trigger_prompt_file inline-XOR con trigger_prompt — entrambi assenti → errore
    Dato che un profilo NON dichiara né "trigger_prompt" né "trigger_prompt_file"
    Quando lancio "labnexus validate test-profile"
    Allora exit code è 2
    E stderr contiene "trigger_prompt obbligatorio (inline o file)"

  Scenario: trigger_prompt_file viene parsato e usato come trigger
    Dato che un profilo dichiara "trigger_prompt_file = \"Prompt_INPUT.rtf\""
    E la cartella di input "/tmp/in/" contiene "Prompt_INPUT.rtf" con contenuto plain text
    Quando lancio "labnexus run --profile test-trigger-file --input /tmp/in --output /tmp/out"
    Allora il trigger usato dal modello è il contenuto del file RTF (plain text estratto)
    E exit code è 0

  Scenario: trigger_prompt_file path traversal è rifiutato
    Dato che un profilo dichiara "trigger_prompt_file = \"../../etc/passwd\""
    Quando lancio "labnexus run --profile test-trigger-malicious --input /tmp/in --output /tmp/out"
    Allora exit code è 2
    E stderr contiene "path non consentito" o "esce dalla cartella di input"

  # --- FR-18, FR-19, FR-20, FR-21: parser nuovi formati ---

  Scenario: .xls file viene parsato come Excel legacy
    Dato che la cartella di input contiene "Rapporto.xls" (Excel 97-2003 OLE)
    E il file ha 2 sheet con celle compilate
    Quando lancio "labnexus check test-profile --input /tmp/in --show-prompt"
    Allora stdout contiene il testo estratto da entrambe le sheet
    E ogni sheet è separata da "--- SHEET: <name> ---"

  Scenario: .doc file viene parsato come Word legacy
    Dato che la cartella di input contiene "Scheda.doc" (Word 97-2003 OLE)
    Quando lancio "labnexus check test-profile --input /tmp/in --show-prompt"
    Allora stdout contiene il testo estratto dal documento
    E NON sono inclusi marcatori di formattazione (font, colori)

  Scenario: .rtf file viene parsato come Rich Text
    Dato che la cartella di input contiene "Prompt_INPUT.rtf" (RTF)
    Quando lancio "labnexus check test-profile --input /tmp/in --show-prompt"
    Allora stdout contiene il plain text estratto dall'RTF
    E NON contiene control words RTF tipo "\\rtf1\\ansi"

  Scenario: file .doc corrotto → warning + skip, no crash
    Dato che la cartella di input contiene "doc-rotto.doc" non parsabile + "doc-buono.doc" valido
    Quando lancio "labnexus run --profile test-profile --input /tmp/in --output /tmp/out"
    Allora exit code è 0 (continua con i file validi)
    E stderr contiene un warning specifico su "doc-rotto.doc"
    E l'output viene generato dal solo "doc-buono.doc"

  # --- FR-22: meta-prompt aggiornato a TOML ---

  Scenario: docs/meta-prompt-genera-profilo.md usa TOML come formato profilo
    Dato che leggo "docs/meta-prompt-genera-profilo.md"
    Allora la sezione "Schema YAML" è sostituita da "Schema TOML"
    E i 3 esempi few-shot sono in formato TOML
    E è menzionato esplicitamente "trigger_prompt_file" come pattern opzionale

  # --- FR-23: test suite update ---

  Scenario: stub profile per BDD usa TOML invece di YAML
    Dato che leggo "features/steps_test.go::minimalProfileTOML"
    Allora la funzione esiste (rinominata da minimalProfileYAML)
    E produce output TOML valido
    E i test BDD esistenti continuano a passare con i nuovi stub

  # --- FR-24: tabella ownership key frontmatter ---

  Scenario: tabella ownership chiavi frontmatter documentata in output.go
    Dato che leggo "internal/output/output.go::Frontmatter"
    Allora un commento documenta quali chiavi sono engine-owned (profilo, modello, provider, data_esecuzione, durata_secondi, token_stimati, file_input, stato, log_file)
    E quali sono profile-default-owned (tipo, stato_qm, profilo_labnexus, locale)
    E quali sono external-owned (valutazione_denis, manuale Denis L2)
