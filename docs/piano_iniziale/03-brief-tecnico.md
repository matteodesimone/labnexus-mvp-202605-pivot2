# Brief per Claude Code — Sviluppo motore LabNexus

> Questo è un brief Prove Out-compliant da passare al framework di vibe coding basato su Claude Code. Contiene Discovery Goal, vincoli, note feature operative. Lascia che il framework spacchetti le fette e proponga il plan: il CTO (Matteo) dà il "go" prima dell'implementazione.

---

## Discovery Goal

```yaml
progetto: LabNexus / AICertus — Sprint 1 / motore + profili + meta-prompt
cliente: LabNexus (Stefano Fiorina, Denis Brazzo)
ruolo_agente: assistente di sviluppo, framework Prove Out
gettoni_budget: 10  # 4 motore + 6 profili (vedi nota gettoni in coda)
durata: 10 giorni lavorativi
stato: avvio
```

### Credo che…

Un singolo eseguibile Go nativo per Mac, configurabile per capability tramite file YAML, possa eseguire i sette task ispettivi LLM-critical del modello AICertus chiamando Qwen 3 in locale via Ollama, producendo output di qualità accettabile per un esperto di certificazioni ACCREDIA. La struttura dei profili YAML, ben formalizzata, è anche generalizzabile in un prompt meta che permette di generare nuovi profili in autonomia tramite Claude.

### Avrò ragione se…

L'eseguibile, lanciato su un input strutturato corrispondente a una delle sette capability, produce un output markdown che:

- è coerente con il template specifico della capability (definito nella KB-ispettore di Denis)
- non contiene allucinazioni rispetto ai dati di input
- mantiene il tono ispettivo richiesto dalla KB-ispettore
- per i primi due profili (revisione, rilievi), è qualitativamente paragonabile agli output del Test 1 e Test 2 manuali già approvati

E inoltre: il prompt meta, quando incollato in una chat Claude insieme alla descrizione testuale di un nuovo task ispettivo, produce un file `<profilo>.yml` valido che l'eseguibile può eseguire senza errori di parsing o di validazione schema.

### C'è un problema se…

- L'output dei profili `revisione` e `rilievi` differisce significativamente dagli output Test 1/Test 2 manuali → problema nella pipeline software (prompt assembly, context injection, streaming)
- Il parsing dei file di input fallisce silenziosamente su formati edge case (PDF complessi, XLSX con formule, DOCX corrotti)
- Il context window viene saturato senza preavviso → manca telemetria sui token consumati
- L'eseguibile non gira correttamente su Mac di consegna a causa di Gatekeeper o problemi di binary distribution
- Il prompt meta produce YAML che pretende di iniettare KB-files inesistenti, o trigger_prompt vuoti, o struttura schema sbagliata

### Non farò…

- **Niente pipeline deterministica** (filesystem walk del SGQ, indicizzazione, scadenze, cross-reference): è materia di Sprint 2
- **Niente orchestrazione automatica** (watcher, trigger, cron): Sprint 2
- **Niente UI web, niente HITL, niente database**: l'eseguibile è solo CLI/macOS dialog, output filesystem markdown
- **Niente firma Apple Developer del binario**: documentiamo la procedura Gatekeeper manuale nel README
- **Niente system prompt scritti da noi per le capability**: usiamo i file della KB-ispettore di Denis (`CLAUDE.md`, `come-pensa-un-ispettore.md`, sezioni ISO, NC-patterns) come system context. L'unico testo scritto da noi per profilo è il `trigger_prompt`, 5-10 righe
- **Niente confronto con Claude/altri provider in produzione**: EUrouter è solo facility tecnica di debug
- **Niente embedding/RAG/vector database**: la KB si carica come context, è piccola abbastanza
- **Niente chunking automatico** se l'input supera il context window: in Sprint 1 i task sono dimensionati per stare dentro, basta loggare un warning a 70% di saturazione
- **Niente PT analysis con calcoli serissimi**: Qwen sui numeri è debole, lo accettiamo come dato di realtà — è anzi parte della validazione dello sprint
- **Niente validazione semantica automatica dell'output del prompt meta**: il file YAML generato dal meta-prompt viene passato al validatore di schema dell'eseguibile (sintassi, campi obbligatori, esistenza file referenziati); la valutazione di qualità del trigger_prompt resta umana

### Starò attento a…

- **Parsing PDF**: librerie Go (`pdfcpu`, `ledongthuc/pdf`) hanno casi limite su PDF complessi. Se un PDF di Denis spacca il parser, primo fallback è chiedergli la conversione in markdown via `pandoc`; secondo fallback è skip con warning esplicito
- **Parsing DOCX**: `unidoc/unioffice` o `nguyenthenguyen/docx` per Go. Idem casi limite
- **Streaming Ollama**: la API `/api/chat` di Ollama è OpenAI-compatible per la parte richiesta/risposta ma con qualche differenza nei chunk SSE. Verificare attentamente
- **Gatekeeper macOS**: testare la procedura manuale di sblocco su un Mac vergine prima della consegna
- **Token counting**: serve una stima dei token di input prima di chiamare Ollama, per loggare il warning a 70% di context. Implementare tokenizer approssimato (è ok una stima per character count diviso 4)
- **Sviluppo nativo Mac**: il binario va compilato sull'ambiente target. La macchina di sviluppo è macOS Apple Silicon, quindi compila e gira lì
- **EUrouter**: API key in env var `EUROUTER_API_KEY`, mai hardcoded mai committato
- **Prompt meta**: deve essere robusto a descrizioni vaghe dell'utente. Se Denis dice solo *"voglio un profilo che gestisca i reclami"*, il meta-prompt deve guidare Claude a fare domande di chiarimento prima di produrre lo YAML

### Lo faremo perché…

Il cliente ha bisogno di sapere se Qwen 3 locale regge il caso d'uso AICertus su capability rappresentative, prima di investire in costruzione del sistema completo (12-18 mesi). Lo strumento è anche il primo deliverable concreto che Denis può usare nel suo lavoro quotidiano. Il prompt meta gli dà un meccanismo per **esplorare ulteriori capability in autonomia tra Sprint 1 e Sprint 2**, generando candidati di profili senza dover aspettare il CTO.

### Lo faremo con…

- **Stack**: Go (la macchina di sviluppo Matteo è Mac Apple Silicon, sviluppo nativo)
- **Provider LLM principale**: Ollama HTTP API su `localhost:11434` (modello Qwen 3, parametri esatti da decidere durante shakedown del profilo `revisione`)
- **Provider LLM secondario**: EUrouter API (eurouter.ai, OpenAI-compatible)
- **Parsing**: librerie Go per docx/xlsx/pdf/csv/md/txt
- **CLI**: `spf13/cobra` o equivalente
- **Progress bar**: `schollz/progressbar/v3` o `vbauerster/mpb`
- **Dialog macOS nativo**: via Cocoa bindings (`progrium/macdriver`) o approccio leggero con `osascript`

### Per rispondere, esplorerò queste domande

Ordinate per incertezza decrescente.

1. **L'eseguibile riproduce gli output Test 1 e Test 2 manuali?** Se sì, il motore è validato. Se no, debug della pipeline software.
2. **Qwen mantiene qualità su task LLM-critical nuovi** (review-pack, audit, equipment, competence, PT)?
3. **Quanto degrada Qwen quando il context si avvicina al limite di 128k**? (Test indipendente fuori motore, in chat manuale Ollama, fuori da questo brief di sviluppo)
4. **Il parsing dei file reali del SGQ di Denis funziona** su tutti i formati che ci passerà?
5. **Il prompt meta produce YAML che funziona davvero al primo colpo**, o richiede iterazione con l'utente?

### Risponderò entro

10 giorni lavorativi dall'avvio.

---

## Note feature operative

Quelle che seguono sono **tracce di note feature** Prove Out — non vincoli sul plan. Il framework può:

- combinare più note feature in una singola fetta (es. setup + primo profilo insieme)
- spaccare una nota in più fette se conviene
- riordinarle se serve

Lo strict ordering è solo: **`feat-001` prima di `feat-002` prima di `feat-003`**, perché il motore deve esistere prima dei profili, e i primi due profili sono lo shakedown del motore. Il `feat-meta` (prompt meta) **viene per ultimo**, dopo aver sviluppato tutti i sette profili, perché astrae il pattern che emerge da quelli.

### feat-001 — Motore LabNexus

```yaml
id: feat-001
tiny_experiment: Un eseguibile Go nativo per Mac legge una cartella di input + un file di configurazione, costruisce un prompt secondo lo schema definito, chiama un provider LLM in streaming, scrive l'output in markdown
vero_se: L'eseguibile, lanciato su un input minimale dummy + un config.yml dummy, produce un file output markdown coerente con il prompt, log completo in console, progress bar visibile durante lo streaming
gettoni: 4
fetta: TBD (decide framework)
stato: da fare
```

**Componenti**:

1. CLI con `cobra`. Comando principale: `labnexus run --profile <name> --input <dir> --output <dir>`. Comando ausiliari: `list`, `describe <profile>`, `check <profile> --input <dir>` (dry-run senza chiamata LLM), `validate <profile>` (controlla schema del file YAML senza eseguire nulla — usato dal meta-prompt per validare l'output)

2. Risoluzione path: i file di profilo stanno in `./profili/<name>.yml` rispetto all'eseguibile. La KB-ispettore sta in `./KB-ispettore/` rispetto all'eseguibile (o path indicato nel config). Tutti i path nei `kb_files:` del config sono **relativi** alla cartella KB-ispettore

3. Lettura cartella input: walk non ricorsivo (input è una cartella piatta con i file da processare). Tipi accettati: `.md`, `.txt`, `.csv`, `.docx`, `.pdf`, `.xlsx`. Estrai testo plain da ognuno

4. Lettura `config.yml`: schema in coda a questo documento. Validazione strict: tutti i `kb_files` esistono fisicamente, `provider` è uno dei valori ammessi, `trigger_prompt` non è vuoto e ha lunghezza minima 50 caratteri. Il comando `validate` espone questa validazione come funzione autonoma

5. Composizione prompt:
   ```
   system_message = concat(read(kb_files), separator="\n\n---\n\n")
   user_message = trigger_prompt + "\n\n## File di input\n\n" + concat(input_files_as_text, separator="\n\n--- FILE: <name> ---\n\n")
   ```

6. Astrazione `LLMProvider`:
   ```go
   type LLMProvider interface {
       Stream(ctx context.Context, system string, user string, opts Options) (<-chan StreamEvent, error)
   }
   ```
   Due implementazioni: `OllamaProvider`, `EurouterProvider`

7. Stima token: prima della chiamata, conta i caratteri / 4 come stima. Se > 70% del `context_window`, logga warning ma procedi. Se > 100%, **errore esplicito**, non chiamare

8. Esecuzione streaming: ricevi token, aggiorna progress bar (basata su rapporto tokens / max_tokens previsti), aggrega in buffer

9. Scrittura output: file markdown in `<output_dir>/<timestamp>_<profile>_<input_descriptor>.md`. Frontmatter YAML con: data esecuzione, profilo, modello usato, provider, durata, token consumati stimati, file di input usati. Body = output del modello

10. Log console: linea per ogni step (caricamento KB, parsing input, stima context, chiamata provider, streaming, scrittura output). Tempo per ciascun step. Errori esplicitati

11. Modalità interattiva: se l'eseguibile viene lanciato senza argomenti (o con un solo argomento posizionale = cartella di input), aprire dialog macOS nativo per chiedere parametri mancanti

12. Distribuzione: target `darwin-arm64`. Build con `go build -o labnexus.app/Contents/MacOS/labnexus`. README di una pagina con procedura Gatekeeper

**Componenti EUrouter (inclusi in `feat-001`)**:
13. Provider `EurouterProvider` implementa `LLMProvider`. API key letta da env var `EUROUTER_API_KEY`. Endpoint `https://api.eurouter.ai/v1/chat/completions`. OpenAI-compatible per request/response shape e streaming SSE
14. Se `provider: eurouter` in un config e la key non è settata, errore esplicito al lancio

**Non farò** (in `feat-001`):
- Più di un input al lancio (ne fa uno alla volta)
- Caching di chiamate, retry intelligente
- Logging strutturato JSON (basta testo)
- Sandbox/permessi avanzati

**Starò attento a**:
- Streaming Ollama: la API restituisce JSON line-delimited (NDJSON), un oggetto per chunk con campo `message.content`. Parsare incrementalmente, non aspettare la response completa
- Streaming EUrouter: SSE standard OpenAI con `data: {...}\n\n`, terminatore `data: [DONE]`
- Dialog macOS: l'approccio `osascript` è il più semplice (`osascript -e 'set folderPath to (choose folder)...'`), accettabile per Sprint 1. Bindings Cocoa veri sono Sprint 2
- Bundle macOS: per il doppio click serve un vero `.app` bundle (con `Info.plist`), non solo un eseguibile. Approfondisci ma minimo

### feat-002 — Profilo `revisione` (shakedown motore #1)

```yaml
id: feat-002
tiny_experiment: Il profilo "revisione" replica il Test 1 manuale (PG_RISK_LAB rev04→rev05 con RT-08) usando l'eseguibile
vero_se: Lanciato su una cartella che contiene PG_RISK_LAB_Rev_00.docx + DE0779_RT_08rev03.pdf + RT-08-rev.05.pdf, l'eseguibile produce una bozza di revisione qualitativamente paragonabile al PG_RISK_LAB_Rev_01_BOZZA-qwen3.6-locale.docx già approvato da Denis
gettoni: 1
fetta: TBD
stato: da fare
```

**Cosa scrivere**:

1. File `profili/revisione.yml`. Schema completo in coda a questo documento
2. `trigger_prompt` di 5-10 righe. Riferimento al template di output che la KB-ispettore di Denis descrive per `/sgq-draft`
3. Selezione `kb_files`: `CLAUDE.md`, `come-pensa-un-ispettore.md`, sezioni 6 e 7 ISO 17025, NC-patterns
4. Esecuzione su input del Test 1 (i 3 file in `KB-ispettore/INPUT_esempio/test01/`)
5. Confronto qualitativo con il file `KB-ispettore/bozze/PG_RISK_LAB_Rev_01_BOZZA-qwen3.6-locale.docx`

**Criteri di confronto** (è un test informale interno, da fare prima della valutazione formale con Denis):
- Frontmatter YAML completo e corretto
- Sezioni di sintesi presenti
- Modifiche citate con callout `[!MODIFICA]`
- Tono ispettivo (no checklist astratte)
- Aderenza ai contenuti delle due versioni RT-08

Se il confronto è OK → motore validato per questa capability. Se è significativamente peggiore → debug della pipeline software prima di andare avanti.

### feat-003 — Profilo `rilievi` (shakedown motore #2)

```yaml
id: feat-003
tiny_experiment: Il profilo "rilievi" replica il Test 2 manuale (rilievi Accredia ACIAA A1) usando l'eseguibile
vero_se: Lanciato sul CSV ACIAA A1, l'eseguibile produce un CAPA Pack per ogni rilievo nel formato definito nel CoWork KB sez. 15.2
gettoni: 1
fetta: TBD
stato: da fare
```

**Cosa scrivere**:
- `profili/rilievi.yml`
- `trigger_prompt`
- `kb_files`: stessi del Test 2 manuale (CLAUDE.md, come-pensa-un-ispettore, NC-patterns)
- Esecuzione sul CSV reale `ACIAA A1 rilievi.csv`

### feat-004 ... feat-008 — Profili nuovi

Per ciascuno dei profili `review-pack`, `audit-checklist`, `equipment-alert`, `competence-gap`, `pt-analysis`:

```yaml
id: feat-XXX
tiny_experiment: Il profilo <nome> esegue la capability <descrizione> come definita nel CoWork KB
vero_se: L'eseguibile produce output coerente con il template della capability nel CoWork KB, senza allucinare rispetto ai dati di input
gettoni: 1
fetta: TBD
stato: da fare
```

**Sviluppo**:
- Selezione `kb_files` rilevanti (CLAUDE.md sempre incluso, più sezioni ISO e/o NC-patterns rilevanti)
- `trigger_prompt` scritto riferendoci alla logica dell'agente nel CoWork KB
- Esecuzione su input reali (forniti da Denis) o sintetici (preparati dal CTO se Denis non ha materiale)

**Ordine**: `review-pack` → `audit-checklist` → `equipment-alert` → `competence-gap` → `pt-analysis` (rischio crescente, ultimo è quello con ragionamento numerico).

I dettagli su input e criteri di successo per ciascun profilo sono nel documento "Piano degli esperimenti".

### feat-meta — Prompt meta per generazione profili (ULTIMO)

```yaml
id: feat-meta
tiny_experiment: Un prompt da incollare in una chat Claude (web o desktop) guida Claude a produrre un file profilo.yml valido per LabNexus, a partire dalla descrizione testuale di un nuovo task ispettivo da parte dell'utente
vero_se: Su 2-3 casi inventati ma realistici (es. "voglio un profilo per gestire reclami clienti", "voglio un profilo per la verifica di conformità di un certificato emesso"), il prompt produce YAML che (a) passa il comando `labnexus validate <profilo>` senza errori, (b) ha un trigger_prompt coerente con la descrizione utente, (c) seleziona file KB-ispettore plausibili
gettoni: 1
fetta: TBD (per ultimo dopo tutti i profili)
stato: da fare
```

**Perché viene per ultimo**: il prompt meta astrae il pattern che emerge dalla scrittura dei sette profili. Tentare di scriverlo prima significa inventare un'astrazione su un problema non ancora risolto. Solo dopo aver scritto sette `trigger_prompt` reali e selezionato sette set di `kb_files` reali, lo schema "come si scrive un buon profilo per LabNexus" è chiaro.

**Cosa il prompt meta deve contenere**:

1. **Identità di Claude in quel contesto**: viene istruito a comportarsi come *"assistente specializzato nella generazione di profili LabNexus a partire da descrizioni di task ispettivi"*
2. **Contesto del sistema**: spiegazione sintetica di cosa è LabNexus (motore + profili YAML), come la KB-ispettore viene utilizzata come system context, come il trigger_prompt funziona
3. **Schema completo del file YAML** che deve produrre (campi obbligatori, campi opzionali, vincoli)
4. **Lista dei file KB-ispettore disponibili** con descrizione di una riga ciascuno (es. *"CLAUDE.md: definisce identità, tono, regole base"*, *"come-pensa-un-ispettore.md: framework di lettura ispettiva"*, ecc.). Questo è cruciale: Claude deve sapere quali file esistono per non inventarne di inesistenti
5. **Esempi few-shot**: 2-3 profili reali già sviluppati (es. `revisione.yml` e `rilievi.yml` come esempi compatti)
6. **Istruzioni di interazione**: il prompt deve istruire Claude a **chiedere chiarimenti all'utente** quando la descrizione iniziale è troppo vaga, prima di produrre lo YAML. Domande tipiche da fare: *"Quale tipo di input l'utente passerà nella cartella di input?"*, *"L'output deve seguire un template specifico del CoWork KB?"*, *"Quali requisiti ISO 17025 sono rilevanti per questo task?"*
7. **Pattern di output finale**: il prompt deve istruire Claude a chiudere la conversazione con **(a) un blocco di codice YAML completo**, **(b) una breve spiegazione delle scelte fatte**, **(c) il comando esatto** che l'utente può eseguire per validare lo YAML con il motore (`labnexus validate <nome>.yml`)

**Cosa il prompt meta NON deve fare**:
- Non scrivere code Go (è un meta-prompt per generare YAML, non codice)
- Non inventare file KB-ispettore che non esistono
- Non bypassare il vincolo "trigger_prompt scritto da umano": il trigger deve essere coerente con il task descritto, ma resta valore umano dietro le quinte (Claude lo scrive, l'utente lo rivede e adatta)
- Non promettere capability di Sprint 2 (es. orchestrazione automatica)

**Test del prompt meta**: provarlo su almeno tre casi inventati di crescente complessità.

- Caso semplice: *"voglio un profilo che, dato un certificato di taratura, genera una scheda di verifica della conformità"*. Output atteso: YAML valido, trigger_prompt chiaro, kb_files = CLAUDE.md + sezione 6 ISO + come-pensa-un-ispettore
- Caso medio: *"voglio un profilo che, dato un set di rapporti di prova, identifica eventuali deviazioni dai limiti di accettabilità del metodo"*. Output atteso: il meta-prompt fa domande di chiarimento prima di produrre lo YAML
- Caso complesso/ambiguo: *"voglio un profilo che mi aiuti con la qualifica dei fornitori"*. Output atteso: il meta-prompt riconosce l'ambiguità, fa più domande, produce uno YAML solo dopo aver capito di quale sub-task si tratta

**Deliverable**: un file `meta-prompt-genera-profilo.md` di 1-2 pagine, con il prompt pronto da copiare-incollare in una chat Claude, più una pagina di istruzioni d'uso per Denis (*"come usarlo: apri claude.ai, incolla questo prompt, descrivi il tuo task, segui le domande di Claude, alla fine ti darà uno YAML che metterai nella cartella profili/ del tuo eseguibile"*).

---

## Schema completo `config.yml` di un profilo

Questo schema è quello che il comando `labnexus validate` controlla. È anche la "specifica" che il prompt meta deve seguire per produrre YAML validi.

```yaml
# Profilo: nome (deve coincidere con il filename: profili/<nome>.yml)
# OBBLIGATORIO
profilo: revisione

# Descrizione di una riga, visualizzata in `labnexus list`
# OBBLIGATORIO
descrizione: Revisione documentale in seguito a cambio di norma di riferimento

# Provider LLM. Scelte ammesse:
#   - ollama  : modello locale via Ollama (porta 11434). Default per Sprint 1.
#   - eurouter: gateway europeo OpenAI-compatible. Solo per debug.
# OBBLIGATORIO
provider: ollama

# Modello specifico:
#   - per provider=ollama:   nome locale (es. "qwen3:14b-q5_K_M")
#   - per provider=eurouter: nome OpenRouter-style (es. "mistralai/mistral-large")
# OBBLIGATORIO
modello: qwen3:14b-q5_K_M

# Parametri opzionali (default ragionevoli se omessi)
temperature: 0.2
max_tokens: 8192
context_window: 128000

# File della KB-ispettore di Denis da concatenare come system context.
# Path relativi alla cartella KB-ispettore configurata.
# I file vengono trattati come prompt normali, nell'ordine indicato.
# OBBLIGATORIO (almeno un file)
# Tutti i file referenziati devono esistere fisicamente
kb_files:
  - CLAUDE.md
  - come-pensa-un-ispettore.md
  - sezioni-ISO/sezione-6-risorse-personale-dotazioni-domande.md
  - sezioni-ISO/sezione-7-processo-metodi-validazione-domande.md
  - NC-patterns/come-rispondere-NC-ACCREDIA.md
  - NC-patterns/errori-fatali-da-evitare.md

# Trigger prompt: l'istruzione di innesco del task.
# OBBLIGATORIO. Minimo 50 caratteri. Scritto da umano (o generato da prompt meta + revisione umana)
trigger_prompt: |
  Sulla base della tua identità e delle regole definite nei file di system context
  che hai appena letto, esegui ora questo compito specifico.
  
  Ti vengono forniti tre file nella cartella di input:
  - un documento del SGQ da revisionare (probabilmente .docx o .md)
  - la versione attualmente in vigore della norma di riferimento (.pdf o .md)
  - la versione nuova della stessa norma (.pdf o .md)
  
  Produci una bozza di revisione del documento allineata alla nuova versione
  della norma. Applica la logica di lettura ispettiva, usa i callout Obsidian
  [!MODIFICA], [!ATTENZIONE], [!ISPETTORE] dove appropriato. Frontmatter YAML
  completo. Tono di senior technical advisor. Non inventare requisiti normativi
  che non sono nei file forniti.

# Documentazione per l'utente: cosa aspettarsi nella cartella di input
# OPZIONALE ma raccomandato
input:
  pattern_attesi:
    - "documento da revisionare (docx|md)"
    - "norma versione vecchia (pdf|md)"
    - "norma versione nuova (pdf|md)"
  istruzioni_se_mancante: |
    Per la revisione servono 3 file: il documento da revisionare, la versione 
    vecchia della norma, la versione nuova.

# Configurazione output
# OPZIONALE
output:
  frontmatter_default:
    tipo: bozza_revisione
    stato: bozza_da_validare_qm
    profilo_labnexus: revisione
    locale: true
```

---

## Quello che il framework Prove Out deve fare adesso

1. Leggere questo brief in modo critico
2. Leggere i file della KB-ispettore di Denis (passati separatamente, in un secondo turno quando il piano è stato approvato)
3. Costruire il piano di sviluppo: spacchettare in fette verticali, allocare gettoni in modo coerente, definire l'ordine
4. Identificare la **prima fetta** — quale parte di `feat-001` attraversa per intero la pipeline minima (anche con tutto mockato dove serve)
5. Presentare il piano al CTO in formato markdown leggibile, **prima di scrivere codice**

Quando il CTO dice "OK procedi", il framework comincia con la prima fetta.

---

## Note gettoni (per il framework)

Il budget di 10 gettoni si scompone così:

- Motore (`feat-001`, include EUrouter): 4 gettoni
- Profilo revisione (`feat-002`): 1
- Profilo rilievi (`feat-003`): 1
- Profilo review-pack (`feat-004`): 1
- Profilo audit-checklist (`feat-005`): 1
- Profilo equipment-alert (`feat-006`): 1
- Profilo competence-gap (`feat-007`): 1
- Profilo pt-analysis (`feat-008`): 1
- Prompt meta per generazione profili (`feat-meta`): 1

**Totale: 10**. Il gettone è un'unità di complessità relativa, non un'ora; il minimo è 1.

Fuori da questo budget e gestiti separatamente dal CTO:
- Test indipendente di tenuta su context lungo (chat manuale Ollama, no codice)
- Report finale e walkthrough

Se durante lo sviluppo emerge che una feature costa più di quanto stimato, il framework deve **fermarsi e proporre il trade-off** al CTO (es. taglio di un profilo successivo, riduzione di scope della feature corrente, escalation per decisione).

Se invece il lavoro va più veloce, i gettoni risparmiati restano accumulati per imprevisti.

---

## Compound a fine sessione

Ogni sessione significativa di sviluppo termina con un compound: cosa ho imparato, cosa cambia nei documenti, cosa va aggiornato nel contesto del progetto per la prossima sessione. È prerogativa del framework — gestione automatica nel suo flusso, non richiesta esplicita del CTO.
