# language: it

Funzionalità: Prompt meta per generazione autonoma di nuovi profili
  In quanto Denis che, tra Sprint 1 e Sprint 2,
  vuole esplorare candidati di nuove capability senza dover aspettare il CTO,
  voglio un prompt da incollare in una chat Claude
  che mi guidi a produrre un file profilo.yml valido,
  con domande di chiarimento quando la mia descrizione è vaga.

  # --- FR-20: esistenza e contenuto del meta-prompt ---

  Scenario: il file meta-prompt esiste con il contenuto richiesto
    Quando consulto il deliverable "meta-prompt-genera-profilo.md"
    Allora il file ha lunghezza tra 1 e 2 pagine
    E contiene una sezione "Identità di Claude in questo contesto"
    E contiene una descrizione sintetica del sistema motore + profili YAML + KB-ispettore
    E contiene lo schema completo del file YAML di profilo (campi obbligatori e opzionali)
    E contiene la lista dei file della KB-ispettore disponibili, ognuno con descrizione di una riga
    E contiene 2-3 esempi few-shot di profili reali (es. revisione.yml, rilievi.yml)
    E istruisce Claude a fare domande di chiarimento quando la descrizione utente è vaga
    E specifica il pattern di output finale: blocco YAML + spiegazione delle scelte + comando "labnexus validate <nome>.yml"

  # --- FR-21: guida d'uso per Denis ---

  Scenario: la guida d'uso istruisce Denis passo-passo
    Quando consulto la guida d'uso allegata al meta-prompt
    Allora la guida ha lunghezza di circa 1 pagina
    E indica di aprire claude.ai
    E indica di incollare il meta-prompt come primo messaggio
    E indica di descrivere il task ispettivo a Claude
    E indica di rispondere alle domande di chiarimento
    E indica di salvare lo YAML risultante in "profili/<nome>.yml" del proprio eseguibile
    E indica di lanciare "labnexus validate <nome>" come ultimo passo

  # --- FR-22: i 3 casi di validazione ---

  Scenario: caso semplice — verifica conformità di un certificato di taratura
    Dato che incollo il meta-prompt in una chat Claude
    Quando descrivo "voglio un profilo che, dato un certificato di taratura, generi una scheda di verifica della conformità"
    Allora Claude produce uno YAML senza fare domande superflue
    E lo YAML ha un trigger_prompt chiaro coerente con il task
    E i kb_files selezionati includono CLAUDE.md e file plausibili (es. sezione 6 ISO 17025, come-pensa-un-ispettore.md)
    E lanciando "labnexus validate <nome>" sullo YAML salvato, exit code è 0

  Scenario: caso medio — deviazioni dai limiti in rapporti di prova
    Dato che incollo il meta-prompt in una chat Claude
    Quando descrivo "voglio un profilo che, dato un set di rapporti di prova, identifichi deviazioni dai limiti di accettabilità del metodo"
    Allora Claude pone almeno una domanda di chiarimento prima di produrre lo YAML
    E le domande riguardano formato dell'input, template di output, requisiti ISO rilevanti
    E dopo le risposte, Claude produce uno YAML che passa "labnexus validate"

  Scenario: caso ambiguo — qualifica fornitori
    Dato che incollo il meta-prompt in una chat Claude
    Quando descrivo "voglio un profilo che mi aiuti con la qualifica dei fornitori"
    Allora Claude riconosce l'ambiguità e pone più domande di chiarimento
    E le domande discriminano tra sub-task possibili (qualifica iniziale vs riqualifica periodica vs valutazione di non conformità del fornitore)
    E lo YAML è prodotto solo dopo che il sub-task è chiaro
    E lo YAML passa "labnexus validate"

  Scenario: il meta-prompt NON inventa file KB inesistenti
    Quando Claude propone uno YAML in uno qualsiasi dei 3 casi
    Allora ogni file in kb_files corrisponde a un file realmente presente nella lista del meta-prompt
    E "labnexus validate" non fallisce su "file kb_files non esiste"
