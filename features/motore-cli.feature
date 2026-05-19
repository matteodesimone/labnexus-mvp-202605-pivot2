# language: it

Funzionalità: Motore LabNexus — CLI, schema profilo, parsing input, prompt, token
  In quanto CTO che usa labnexus per condurre gli esperimenti dello sprint,
  voglio una CLI prevedibile con sottocomandi mirati,
  così da poter validare i profili, esplorare gli input e lanciare le capability
  senza ambiguità e senza spendere token per errore.

  Sfondo:
    Dato che l'eseguibile "labnexus" è installato nella cartella corrente
    E la cartella "./profili/" contiene almeno "revisione.yml" e "rilievi.yml"
    E la cartella "./KB-ispettore/" esiste con i file della KB di Denis

  # --- FR-1 / FR-2: CLI subcommands & path ---

  Scenario: list mostra i profili disponibili
    Quando lancio "labnexus list"
    Allora exit code è 0
    E stdout contiene la riga "revisione"
    E stdout contiene la riga "rilievi"
    E ciascuna riga ha una descrizione di una riga di accanto

  Scenario: describe espone la documentazione di un profilo
    Quando lancio "labnexus describe revisione"
    Allora exit code è 0
    E stdout riporta provider, modello, lista dei kb_files e i pattern di input attesi

  Scenario: flag obbligatori in run sono enforced
    Quando lancio "labnexus run" senza altri argomenti
    Allora exit code è 2
    E stderr contiene "profile"
    E stderr contiene "input"
    E stderr contiene "output"

  Scenario: profilo inesistente fallisce loud
    Quando lancio "labnexus run --profile inesistente --input /tmp/in --output /tmp/out"
    Allora exit code è 2
    E stderr contiene "profilo non trovato"
    E stderr elenca i profili disponibili

  # --- FR-3: schema YAML e comando validate ---

  Scenario: validate accetta un profilo valido
    Dato che esiste "./profili/revisione.yml" con tutti i campi obbligatori e trigger_prompt di 200 caratteri
    Quando lancio "labnexus validate revisione"
    Allora exit code è 0
    E stdout contiene "schema OK"

  Scenario: validate rifiuta un profilo con trigger_prompt troppo corto
    Dato che esiste "./profili/test-corto.yml" con trigger_prompt di 20 caratteri
    Quando lancio "labnexus validate test-corto"
    Allora exit code è 2
    E stderr contiene "trigger_prompt"
    E stderr contiene "almeno 50 caratteri"

  Scenario: validate rifiuta un profilo con kb_files inesistente
    Dato che esiste "./profili/test-kb-mancante.yml" con kb_files che include "non-esiste.md"
    Quando lancio "labnexus validate test-kb-mancante"
    Allora exit code è 2
    E stderr contiene "non-esiste.md"
    E stderr contiene "non esiste"

  Scenario: validate rifiuta provider non ammesso
    Dato che esiste "./profili/test-provider.yml" con provider "openai"
    Quando lancio "labnexus validate test-provider"
    Allora exit code è 2
    E stderr contiene "provider"
    E stderr contiene "ollama"
    E stderr contiene "eurouter"

  # --- FR-4: parsing input multi-formato ---

  Scenario: check estrae testo da MD, TXT, CSV, PDF, DOCX, XLSX
    Dato che la cartella "/tmp/input-mix" contiene un file per ognuno dei formati .md, .txt, .csv, .pdf, .docx, .xlsx
    Quando lancio "labnexus check revisione --input /tmp/input-mix"
    Allora exit code è 0
    E stdout elenca 6 file processati con relativa dimensione in caratteri estratti

  Scenario: formato non supportato viene skippato con warning
    Dato che la cartella "/tmp/input-misto" contiene "documento.md" e "foto.jpg"
    Quando lancio "labnexus check revisione --input /tmp/input-misto"
    Allora exit code è 0
    E stderr contiene "foto.jpg"
    E stderr contiene "formato non supportato"
    E stdout indica 1 file processato

  Scenario: walk non è ricorsivo
    Dato che la cartella "/tmp/input-annidato" contiene "doc.md" e la sottocartella "sub/" con "altro.md"
    Quando lancio "labnexus check revisione --input /tmp/input-annidato"
    Allora stdout indica 1 file processato
    E stdout non menziona "altro.md"

  # --- FR-5: composizione prompt ---

  Scenario: check produce un'anteprima del prompt composto
    Dato che "./profili/revisione.yml" dichiara kb_files = ["CLAUDE.md", "come-pensa-un-ispettore.md"]
    E la cartella "/tmp/input-test1" contiene 3 file di test
    Quando lancio "labnexus check revisione --input /tmp/input-test1 --show-prompt"
    Allora stdout mostra un system_message che è la concatenazione di CLAUDE.md e come-pensa-un-ispettore.md separati da "---"
    E stdout mostra un user_message che inizia con il trigger_prompt
    E user_message contiene la sezione "## File di input"
    E user_message contiene per ogni file di input un blocco preceduto da "--- FILE: <nome> ---"

  # --- FR-6: stima token ---

  Scenario: warning sopra il 70% del context window
    Dato che "./profili/revisione.yml" ha context_window = 128000
    E la composizione del prompt produce ~96000 token stimati (char_count/4)
    Quando lancio "labnexus check revisione --input /tmp/input-pesante"
    Allora exit code è 0
    E stderr contiene "warning"
    E stderr contiene "75%" oppure una percentuale > 70%

  Scenario: errore sopra il 100% del context window
    Dato che "./profili/revisione.yml" ha context_window = 128000
    E la composizione del prompt produce ~140000 token stimati
    Quando lancio "labnexus run --profile revisione --input /tmp/input-troppo-grande --output /tmp/out"
    Allora exit code è 2
    E stderr contiene "context window"
    E stderr contiene "supera"
    E nessuna chiamata di rete al provider è stata effettuata

  # --- FR-12: env var EUrouter ---

  Scenario: provider eurouter senza API key
    Dato che "./profili/debug-eurouter.yml" dichiara provider "eurouter"
    E l'env var "EUROUTER_API_KEY" non è impostata
    Quando lancio "labnexus run --profile debug-eurouter --input /tmp/in --output /tmp/out"
    Allora exit code è 2
    E stderr contiene "EUROUTER_API_KEY"
    E nessuna chiamata di rete è stata effettuata

  Scenario: override del provider via flag --provider
    Dato che "./profili/revisione.yml" dichiara provider "ollama"
    E l'env var "EUROUTER_API_KEY" è impostata
    Quando lancio "labnexus run --profile revisione --provider eurouter --input /tmp/in --output /tmp/out"
    Allora la chiamata di rete va a "api.eurouter.ai"
    E nessuna chiamata viene fatta a "localhost:11434"
    E il frontmatter dell'output riporta provider: "eurouter"

  Scenario: override del provider via env var LABNEXUS_PROVIDER
    Dato che "./profili/revisione.yml" dichiara provider "ollama"
    E l'env var "LABNEXUS_PROVIDER" è impostata a "eurouter"
    E l'env var "EUROUTER_API_KEY" è impostata
    Quando lancio "labnexus run --profile revisione --input /tmp/in --output /tmp/out"
    Allora il provider effettivo è "eurouter"
    E il frontmatter dell'output lo registra

  Scenario: flag --provider ha precedenza sull'env var
    Dato che l'env var "LABNEXUS_PROVIDER" è impostata a "eurouter"
    Quando lancio "labnexus run --profile revisione --provider ollama --input /tmp/in --output /tmp/out"
    Allora il provider effettivo è "ollama"
