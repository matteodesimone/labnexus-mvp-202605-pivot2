# Spec: LabNexus Sprint 1 — motore + 7 profili + prompt meta

> Project type: **prove-out** (S, 20 gettoni / 3 settimane; 8 già spesi, tetto interno residuo 15).
> Domanda contrattuale di sprint: *un modello locale (Qwen 3 via Ollama, context 128k) produce output di qualità accettabile per Denis come QM sulle capability ispettive del modello AICertus?*
> Discovery Goal di riferimento: `docs/piano_iniziale/03-brief-tecnico.md` (sezione Discovery Goal).

## Summary

Costruire `labnexus`, un eseguibile macOS nativo (Go, Apple Silicon) che esegue **una capability ispettiva alla volta** del modello AICertus, chiamando **Qwen 3 in locale via Ollama** in streaming, e produce output markdown su filesystem. La capability è descritta da un **profilo YAML** che dichiara quali file della KB-ispettore di Denis caricare come system context, quale `trigger_prompt` usare come user prompt, e i parametri tecnici. Sprint 1 produce **7 profili** (le capability A-G del modello AICertus, in ordine di rischio crescente) e un **prompt meta** che, incollato in una chat Claude esterna, guida la generazione di nuovi profili YAML. Provider secondario **EUrouter** disponibile come facility di debug interno (non operativo sui dati reali). Modalità d'uso: doppio click (dialog macOS), drag&drop di cartella sull'icona, riga di comando.

Il deliverable è lo **strumento per condurre i 7 esperimenti di validazione** descritti in `01-piano-esperimenti.md`. La valutazione finale di ciascuna capability spetta a Denis (griglia condivisa: allucinazioni, qualità ispettiva, profondità, tenuta, utilizzabilità pratica).

## Actors

- **Denis Brazzo** — QM esperto SGQ + ispettore ACCREDIA. Utente operativo dell'eseguibile, autore della KB-ispettore, giudice formale sulla qualità degli output (conformity criterion del cliente).
- **Stefano Fiorina** — referente cliente LabNexus s.r.l. Stakeholder strategico, partecipa ai checkpoint.
- **Matteo De Simone (CTO)** — sviluppa l'eseguibile, prepara dati sintetici dove Denis non fornisce reali, esegue gli esperimenti e raccoglie i giudizi formali, redige il report finale.
- **Claude (chat web)** — esecutore del `prompt meta`: NON parte dell'eseguibile, ma utente esterno del prompt prodotto. Genera bozze di profilo YAML che poi `labnexus validate` controlla a livello di schema.

## Functional Requirements

Le FR si raggruppano in tre famiglie: **motore** (`feat-001`), **profili** (`feat-002`…`feat-008`), **prompt meta** (`feat-meta`). Lo strict ordering tra feat resta vincolo del `/v-plan`.

### Motore — `feat-001` (4 gettoni)

- **FR-1** L'eseguibile espone una CLI con i sottocomandi `run`, `list`, `describe <profile>`, `check <profile> --input <dir>` (dry-run senza chiamata LLM), `validate <profile>` (controllo schema YAML senza esecuzione). Il comando principale `run` accetta `--profile <name>`, `--input <dir>`, `--output <dir>` come flag obbligatori e fallisce loud se assenti o invalidi.
- **FR-2** L'eseguibile risolve i path relativi alla propria posizione: profili in `./profili/<nome>.yml`, KB-ispettore in `./KB-ispettore/`. I `kb_files` dichiarati nel profilo sono path relativi alla cartella KB-ispettore.
- **FR-3** L'eseguibile carica il file di profilo YAML e ne **valida lo schema** (campi obbligatori: `profilo`, `descrizione`, `provider`, `modello`, `kb_files`, `trigger_prompt` ≥ 50 caratteri; vincolo: tutti i file in `kb_files` devono esistere fisicamente; `provider` ∈ {`ollama`, `eurouter`}). Lo stesso validatore è invocabile come comando autonomo (`labnexus validate <profile>`), usato sia dall'utente sia dal flusso del `prompt meta`.
- **FR-4** L'eseguibile legge la cartella di input con **walk non ricorsivo**. Accetta `.md`, `.txt`, `.csv`, `.docx`, `.pdf`, `.xlsx` ed estrae il testo plain di ogni file. Su formato non supportato, log esplicito e skip del file (warning, non errore fatale).
- **FR-5** L'eseguibile compone il prompt secondo lo schema:
  - `system_message` = concatenazione dei `kb_files` letti in ordine, separati da `\n\n---\n\n`
  - `user_message` = `trigger_prompt` + `\n\n## File di input\n\n` + concatenazione del testo estratto, separato da `\n\n--- FILE: <name> ---\n\n`
- **FR-6** L'eseguibile stima il numero di token in input come `char_count / 4`. Se la stima supera il **70%** del `context_window` del profilo, logga un warning ma procede. Se supera il **100%**, **errore esplicito** e nessuna chiamata al provider.
- **FR-7** L'eseguibile espone l'astrazione `LLMProvider` con metodo `Stream(ctx, system, user, opts) (<-chan StreamEvent, error)` e due implementazioni:
  - `OllamaProvider` — chiama Ollama su `localhost:11434`, API `/api/chat`, **NDJSON line-delimited streaming** (un oggetto JSON per chunk con campo `message.content`).
  - `EurouterProvider` — chiama `https://api.eurouter.ai/v1/chat/completions`, **SSE OpenAI-compatible** (`data: {...}\n\n`, terminatore `data: [DONE]`), API key letta dall'env var `EUROUTER_API_KEY`.
- **FR-8** Durante lo streaming, l'eseguibile aggiorna una **progress bar** (basata sul rapporto token ricevuti / `max_tokens` previsti) e mostra log per ogni step: caricamento KB, parsing input, stima context, chiamata provider, streaming, scrittura output, con tempo di ogni step.
- **FR-9** L'eseguibile scrive l'output come singolo file markdown in `<output_dir>/<timestamp>_<profile>_<input_descriptor>.md` con frontmatter YAML contenente: data esecuzione, profilo, modello, provider, durata, token consumati stimati, lista file di input.
- **FR-10** **Modalità interattiva (TUI cross-platform)**: se lanciato senza argomenti, l'eseguibile presenta un **flusso TUI sequenziale** (libreria Go cross-platform — default proposto `charmbracelet/huh`, scelta finale in `/v-plan`) che chiede: profilo (lista da `list` con descrizioni), cartella di input (path con autocomplete o input testuale), cartella di output. Se lanciato con un singolo argomento posizionale interpretabile come cartella esistente, lo usa come input e chiede solo profilo e output. Se la TUI rileva di non essere in TTY (stdin/stderr non-interattivi), termina con exit 2 e messaggio "modalità interattiva non disponibile in non-TTY — usa `labnexus run` con i flag".
- **FR-11** L'eseguibile è distribuito in due forme: (a) bundle **`labnexus.app`** per macOS Apple Silicon (`darwin/arm64`) — il doppio click apre Terminal con la TUI di FR-10 in modalità live; il drag&drop di una cartella sull'icona `.app` passa il path come argomento posizionale (input pre-selezionato). Il bundle minimo include `Info.plist` valido con `CFBundleExecutable`, `LSHandlerRank`, e dichiarazione `Folder` per drag&drop. (b) binario **`labnexus`** standalone per `linux/amd64` — eseguibile da terminale (CLI o TUI), nessun bundle, nessun drag&drop GUI (su Linux/WSL2 non c'è equivalente Finder universalmente disponibile).
- **FR-12** Entrambi i provider sono first-class. **Override del provider del profilo** disponibile via flag `--provider <ollama|eurouter>` (o env var `LABNEXUS_PROVIDER`) per permettere allo stesso profilo di girare su Ollama (deliverable Denis) o su EUrouter (sviluppo CTO su macOS) senza mantenere due copie del profilo. Se `provider: eurouter` (per dichiarazione del profilo o per override) e `EUROUTER_API_KEY` non è settata, errore esplicito al lancio prima di qualsiasi chiamata di rete. Il provider effettivamente usato in ogni esecuzione viene registrato nel frontmatter dell'output (FR-9).

### Profili — `feat-002` … `feat-008` (1 gettone ciascuno)

Lo strict ordering è: A → B → C → D → E → F → G. Le prime due (A, B) sono **shakedown del motore** contro i golden file dei Test 1 e 2 manuali già approvati da Denis.

- **FR-13** Profilo `revisione` (Capability A, shakedown #1): file `profili/revisione.yml`, `kb_files` = `CLAUDE.md` + `come-pensa-un-ispettore.md` + sezioni 6/7 ISO 17025 + `NC-patterns/come-rispondere-NC-ACCREDIA.md` + `NC-patterns/errori-fatali-da-evitare.md`. `trigger_prompt` di 5-10 righe coerente con il template `/sgq-draft` descritto nella KB-ispettore. **Vero se**: lanciato sui 3 file del Test 1 (`PG_RISK_LAB_Rev_00` + `DE0779_RT_08rev03.pdf` + `RT-08-rev.05.pdf`), l'output è **qualitativamente paragonabile** al golden `PG_RISK_LAB_Rev_01_BOZZA-qwen3.6-locale.docx` (frontmatter YAML completo, sezioni di sintesi presenti, modifiche citate con callout `[!MODIFICA]`, tono ispettivo, aderenza ai contenuti delle due versioni RT-08, no allucinazioni). Se NO → debug del motore prima di passare a `feat-003`.
- **FR-14** Profilo `rilievi` (Capability B, shakedown #2): file `profili/rilievi.yml`, `kb_files` = `CLAUDE.md` + `come-pensa-un-ispettore.md` + NC-patterns. `trigger_prompt` allineato al template `CAPA Pack` del CoWork KB sez. 15.2. **Vero se**: lanciato sul CSV reale `ACIAA A1 rilievi.csv` del Test 2, l'output produce **un CAPA Pack per ciascun rilievo** secondo la triade meccanismo / estensione / efficacia, distinguendo correttamente quando basta correzione e quando serve azione correttiva. (Golden file di output non conservato: confronto qualitativo col template.)
- **FR-15** Profilo `review-pack` (Capability C, prima capability nuova): `kb_files` rilevanti incluso `CLAUDE.md`. `trigger_prompt` riferito al template `Management_Review_Pack` del CoWork KB sez. 15.4. **Vero se**: l'eseguibile, lanciato su dati strutturati di un anno operativo di un laboratorio (NC, audit interni, reclami, PT con z-score, stato apparecchiature, KPI), produce il Pack a **13 sezioni** senza degrado tra sezione 1 e 13 e con coerenza tra sintesi e decisioni finali.
- **FR-16** Profilo `audit-checklist` (Capability D): `trigger_prompt` riferito all'agente AuditAgent. **Vero se**: lanciato su un risk register (15-25 voci) + descrizione di scope d'audit, produce 15-30 domande ispettive raggruppate per area, con campioni documentali da verificare e rilievi anticipati con livello di rischio.
- **FR-17** Profilo `equipment-alert` (Capability E): `trigger_prompt` riferito al template `Equipment_Alert` del CoWork KB sez. 15.3. **Vero se**: lanciato su scheda apparecchiatura + descrizione di evento (taratura scaduta o NC su certificato), produce stato attuale, rischio tecnico, azioni proposte, bozza email fornitore, checklist al rientro, con ragionamento causale documentabile su metodi di prova associati.
- **FR-18** Profilo `competence-gap` (Capability F): `trigger_prompt` riferito all'agente CompetenceAgent. **Vero se**: lanciato su matrice competenze (10-15 righe × 8-12 colonne) + descrizione di procedura/metodo nuovo, produce gap report (autorizzati / non autorizzati), piano formazione, autorizzazioni da rilasciare/aggiornare.
- **FR-19** Profilo `pt-analysis` (Capability G, ultima): `trigger_prompt` riferito all'agente PTQCAgent. **Vero se**: lanciato su risultati PT con z-score + elenco metodi + idealmente storico 2-3 anni, produce Investigation Pack con riassunto, classificazione rischio per metodo, trend storico, azioni immediate per z-score fuori soglia.

### Prompt meta — `feat-meta` (1 gettone, ultimo)

- **FR-20** Esistenza del deliverable `meta-prompt-genera-profilo.md` (1-2 pagine) da incollare in una chat Claude esterna. Contiene: (a) identità di Claude come *assistente specializzato in profili LabNexus*, (b) descrizione sintetica del sistema motore + profili YAML + KB-ispettore come system context, (c) **schema completo del file YAML** atteso, (d) **lista dei file della KB-ispettore disponibili** con descrizione di una riga ciascuno (per evitare invenzioni), (e) 2-3 profili reali come **esempi few-shot** (es. `revisione.yml`, `rilievi.yml`), (f) istruzioni a Claude di **fare domande di chiarimento** quando la descrizione utente è vaga, (g) **pattern di output finale**: blocco di codice YAML + breve spiegazione delle scelte + comando esatto `labnexus validate <nome>.yml`.
- **FR-21** Insieme al meta-prompt, breve guida d'uso per Denis (1 pagina) con istruzioni: apri claude.ai, incolla, descrivi il task, segui le domande, salva lo YAML in `profili/`, valida con `labnexus validate`.
- **FR-22** Il prompt meta è validato funzionalmente su **3 casi inventati di crescente complessità**:
  - *Caso semplice*: «profilo per verifica di conformità di un certificato di taratura». Output atteso: YAML valido, trigger_prompt chiaro, kb_files plausibili al primo colpo.
  - *Caso medio*: «profilo per identificare deviazioni dai limiti di accettabilità in un set di rapporti di prova». Output atteso: il meta-prompt **fa almeno una domanda di chiarimento** prima di produrre lo YAML.
  - *Caso ambiguo*: «profilo per la qualifica dei fornitori». Output atteso: il meta-prompt riconosce l'ambiguità, fa più domande, produce YAML solo dopo aver capito il sub-task.

Per tutti e tre, lo YAML prodotto **passa `labnexus validate`** senza errori di schema.

## Non-Functional Requirements

- **NFR-1 — Locale mandatory _per il deliverable consegnato_.** Il binario consegnato a Denis è configurato per girare con `provider: ollama` su Qwen 3 in locale, sul suo hardware macOS Apple Silicon. Questa è la domanda contrattuale dello sprint post-pivot 2 e non si tocca. **Durante lo sviluppo del CTO, entrambi i provider sono first-class citizens, per ragioni operative**: su macOS del CTO Ollama non è disponibile → si usa `provider: eurouter` su dati sintetici o sui golden file già approvati; su Linux/WSL2 del CTO con NVIDIA RTX 4090 si usa `provider: ollama` con accelerazione GPU. La restrizione resta una: **nessun dato reale del SGQ di Denis viene mandato a EUrouter senza approvazione esplicita** (i materiali reali in `materiali-dominio/test-precedenti/` sono già stati pubblicamente discussi col cliente, quindi utilizzabili come golden; i materiali futuri di Denis no).
- **NFR-2 — Portability.** Due target Sprint 1: **`darwin/arm64`** (consegna Denis, macOS Apple Silicon) e **`linux/amd64`** (sviluppo CTO su Windows WSL2 con NVIDIA RTX 4090 24GB per Ollama accelerato via GPU). Cross-compile Go nativo (`GOOS`/`GOARCH`). Stessa codebase, due artefatti: `labnexus.app` bundle per macOS, binario `labnexus` standalone per Linux. `darwin/amd64` (Mac Intel) non in scope salvo trigger esplicito al checkpoint.
- **NFR-3 — Performance/UX.** Lo streaming è incrementale: niente attesa della response completa prima di mostrare la progress bar o di iniziare a scrivere il buffer di output. Tempo da invocazione a primo token visibile ≤ 5 s su input nominale.
- **NFR-4 — Determinismo CLI per piping.** stdout = output strutturato (log line-oriented, NDJSON-compatibile dove utile per scripting). stderr = errori e diagnostica. Niente ANSI/spinner quando non in TTY (vedi standard cli-tool del framework).
- **NFR-5 — Exit codes semantici.** `0` = successo (output generato), `1` = fallimento generico (rete/provider/parsing), `2` = misuse/invalid invocation (flag mancante, profilo inesistente, schema invalido, file di input vuoto).
- **NFR-6 — Sicurezza.** `EUROUTER_API_KEY` letta esclusivamente da env var, mai loggata, mai committata. Nessun secret in codice. Path traversal protetto: i `kb_files` non possono uscire da `./KB-ispettore/`.
- **NFR-7 — Privacy.** Vedi sezione "Data & Privacy". I file di input sono letti **in sola lettura**; nessun file della cartella di input viene modificato o copiato altrove se non riportato in chiaro all'utente.
- **NFR-8 — Quality metric / conformity criterion.** Per **ciascuno dei 7 profili** la validazione si svolge in **due livelli**:
  - **Livello 1 — Pre-check tecnico del CTO (strutturale)**: verifica oggettiva che il software abbia prodotto un output ben formato (frontmatter completo, struttura attesa per la capability, no allucinazioni macroscopiche rispetto agli input forniti). Se fallisce → debug della pipeline software, l'output **non viene sottomesso a Denis**.
  - **Livello 2 — Validazione formale di Denis (qualitativa, conformity criterion)**: solo dopo il pass di L1, Denis valuta secondo la griglia condivisa: allucinazioni, qualità ispettiva, profondità, tenuta, utilizzabilità pratica. Esito ufficiale: *validata* / *validata con riserva* / *non validata*. Registrato nel report di sprint e nel frontmatter dell'output (`valutazione_denis: ...`).

  Per `revisione` (Capability A) il Livello 1 include il confronto strutturale col golden file del Test 1; per `rilievi` (B) solo col template CAPA Pack (TD-2). Per le capability C-G nuove il Livello 1 è strutturale rispetto al template CoWork KB della capability; il giudizio qualitativo resta interamente a Denis al Livello 2. **Nessuna capability è "validata" senza il Livello 2 di Denis**.
- **NFR-9 — Stabilità di scope.** Niente nuovi rami di sviluppo durante Sprint 1 senza motivazione esplicita su dati raccolti (vedi `02-project.md`, sezione Key Decisions).

## Input/Output

### Input

- **Cartella di input** (`--input <dir>`): cartella piatta (no ricorsione) contenente i file da processare per la capability. Formati: `.md`, `.txt`, `.csv`, `.docx`, `.pdf`, `.xlsx`. Il numero e la natura dei file dipendono dalla capability:
  - `revisione`: 3 file (documento da revisionare, norma vecchia, norma nuova)
  - `rilievi`: 1 CSV (lista rilievi)
  - `review-pack`: N file strutturati (NC dell'anno, audit interni, reclami, PT, apparecchiature, KPI)
  - `audit-checklist`: risk register + descrizione scope
  - `equipment-alert`: scheda apparecchiatura + descrizione evento
  - `competence-gap`: matrice competenze + descrizione procedura/metodo
  - `pt-analysis`: risultati PT + elenco metodi + idealmente storico
- **Profilo** (`--profile <name>`): nome che risolve in `./profili/<name>.yml`.
- **KB-ispettore**: cartella `./KB-ispettore/` letta **al runtime di ogni esecuzione** (no caching, no ricompilazione).
- **Env var**: `EUROUTER_API_KEY` se e solo se un profilo dichiara `provider: eurouter`.

### Output

- **Cartella di output** (`--output <dir>`): un singolo file markdown per esecuzione, denominato `<timestamp>_<profile>_<input_descriptor>.md`, con frontmatter YAML (data, profilo, modello, provider, durata, token stimati, file di input) e body con l'output completo del modello.
- **stdout / stderr**: log strutturato per step (vedi FR-8) con tempi.
- **Exit code**: 0/1/2 secondo NFR-5.

## Interface & Intuitiveness

### Main Instructions

| Screen/Command | Main Instruction |
|---|---|
| `labnexus` (no args, doppio click) | L'utente sceglie nel dialog quale capability eseguire, poi seleziona cartella input e cartella output. |
| Drag&drop di cartella su `labnexus.app` | L'utente trascina una cartella, il dialog chiede solo profilo e cartella output. |
| `labnexus run --profile <name> --input <dir> --output <dir>` | L'utente lancia direttamente l'esecuzione di una capability su input noto. |
| `labnexus list` | L'utente vede l'elenco dei profili disponibili con descrizione di una riga ciascuno. |
| `labnexus describe <profile>` | L'utente legge cosa fa quel profilo, quali input si aspetta, quale KB carica. |
| `labnexus check <profile> --input <dir>` | L'utente verifica che l'input sia processabile e che la stima token rientri nel context, senza chiamare il modello. |
| `labnexus validate <profile>` | L'utente (o il flusso del prompt meta) verifica che lo YAML del profilo rispetti lo schema, senza eseguire nulla. |

### Intuitiveness Checklist

Per la primaria interazione (`labnexus run` da CLI; modalità interattiva via doppio click):
- [x] **Discoverable**: `--help` lista tutti i sottocomandi e tutti i flag obbligatori (cobra `MarkFlagRequired`); `list` mostra i profili senza dover aprire file YAML.
- [x] **Comprehensible**: ogni sottocomando ha una one-liner; `describe` documenta cosa fa il profilo prima di eseguirlo; `check` è il dry-run esplicito.
- [x] **Forgiving**: `check` permette di sondare l'input senza spendere tempo/token; `validate` permette di sondare lo schema; saturazione del context oltre 100% blocca con errore esplicito invece di troncare silenziosamente; flag mancante fallisce loud con elenco delle opzioni valide.
- [x] **Feedback**: log per step con tempi, progress bar durante streaming, frontmatter YAML che traccia parametri usati, Finder che si apre sull'output a fine esecuzione (modalità interattiva).

La review completa contro gli 8 attributi McKay avviene in `/v-review`.

## Edge Cases

- **EC-1 — PDF complesso che spacca il parser.** Librerie Go (`pdfcpu`, `ledongthuc/pdf`) hanno casi limite. Comportamento: warning esplicito che identifica il file, suggerimento di conversione manuale via `pandoc` come fallback, **skip del singolo file** non interruzione dell'intera esecuzione. Se ≥ 50% dei file di input falliscono il parsing, errore esplicito (exit 1) — l'esecuzione non avrebbe senso.
- **EC-2 — DOCX corrotto o non leggibile.** Analogo a EC-1.
- **EC-3 — Context window saturato a > 70%.** Warning, esecuzione procede. A > 100% errore esplicito (exit 2, non si chiama il provider).
- **EC-4 — `kb_files` referenziato non esiste.** `labnexus validate` lo rileva e fallisce; `labnexus run` lo rileva al caricamento e fallisce con exit 2 prima di chiamare il provider.
- **EC-5 — `trigger_prompt` vuoto o < 50 caratteri.** `labnexus validate` lo rileva e fallisce.
- **EC-6 — `provider: eurouter` con `EUROUTER_API_KEY` mancante.** Errore esplicito al lancio, prima di qualsiasi tentativo di rete.
- **EC-7 — Ollama non in ascolto su `localhost:11434`.** Errore di connessione esplicito con istruzione "verifica che `ollama serve` sia attivo".
- **EC-8 — Streaming interrotto a metà.** Il buffer parziale viene comunque scritto come file output con frontmatter che segnala `stato: interrotto`, exit 1.
- **EC-9 — Cartella di input vuota o tutti file non supportati.** Exit 2 con messaggio esplicito.
- **EC-10 — Gatekeeper macOS blocca il binario al primo lancio.** Solo target `darwin/arm64`. Comportamento atteso: documentato nel README, procedura standard control-clic → "Apri". Non è un bug. Su `linux/amd64` Gatekeeper non esiste; il binario è solo `chmod +x`.
- **EC-16 — Modalità interattiva invocata in non-TTY** (es. `labnexus | cat` senza argomenti). La TUI non può funzionare. Exit 2 con messaggio chiaro "usa `labnexus run` con i flag in modalità non-interattiva". Vedi FR-10.
- **EC-11 — Profilo YAML prodotto dal prompt meta con campi inventati.** `labnexus validate` lo intercetta prima che entri in produzione.
- **EC-12 — Output del modello che non rispetta il template di capability.** Non è un errore di codice — è un dato di validazione del modello. Lo registra Denis nella griglia.
- **EC-13 — Streaming Ollama: chunk NDJSON malformato.** Parser tollerante: skip del chunk con warning, prosegue. Se ≥ 10 chunk consecutivi malformati, errore.
- **EC-14 — Streaming EUrouter: terminatore `data: [DONE]` mai ricevuto.** Timeout configurabile; al timeout, scrive output parziale + frontmatter `stato: timeout`, exit 1.
- **EC-15 — Output file con timestamp collidente** (esecuzioni concorrenti improbabili ma possibili). Suffisso incrementale `_2`, `_3` per disambiguare. Mai overwrite.

## Data & Privacy

**Personal data involved**: i file di input del SGQ di Denis **possono contenere dati personali limitati** del personale del laboratorio (nomi tecnici, qualifiche, autorizzazioni), dei fornitori (ragioni sociali, referenti) e dei clienti del laboratorio (nominativi nei rapporti di prova). Sono dati professionali, non sensibili ai sensi del GDPR, ma rientrano comunque nel regime di trattamento.

- **Purpose**: produrre output ispettivi (revisioni, CAPA, review pack, ecc.) per il QM. Trattamento legittimo per finalità professionali del laboratorio.
- **Retention**: l'eseguibile **non conserva alcun dato**. Legge i file di input in sola lettura, scrive l'output nella cartella indicata dall'utente. Nessun database, nessuna cache, nessun upload. La retention è gestita dall'utente sul filesystem.
- **Deletion**: l'utente cancella manualmente i file di output. Il filesystem locale è l'unica copia.
- **Third parties**: con `provider: ollama` (Sprint 1 default) **nessun dato lascia la macchina**. Con `provider: eurouter` (facility di debug) i dati transitano per gateway EU OpenAI-compatible (eurouter.ai, dati residenti in UE, GDPR-compliant). **Vincolo NFR-1**: EUrouter non viene usato sui dati reali del SGQ senza approvazione esplicita; in Sprint 1 resta strumento di debug interno su dati sintetici o anonimizzati.
- **Secrets**: `EUROUTER_API_KEY` solo in env var, mai loggata, mai committata.
- **Logging**: i log per step (FR-8) non riportano contenuto dei file di input; solo nomi file, dimensioni, durate, token stimati. Il frontmatter YAML dell'output riporta i **nomi** dei file di input, non il contenuto.

## Out of Scope

Lo Sprint 1 **non** consegna:

- **OOS-1** Pipeline deterministica (filesystem walk del SGQ, indicizzazione documentale, scadenze, cross-reference, riferimenti normativi obsoleti, audit trail). → Sprint 2.
- **OOS-2** Orchestrazione automatica (watcher, trigger su filesystem, cron, esecuzione multi-capability in sequenza). Una esecuzione = una capability. → Sprint 2.
- **OOS-3** UI web, dashboard HITL, database centralizzato, multi-tenancy. → Sprint 2.
- **OOS-4** ISO 17034. Tutti gli esperimenti su perimetro **ISO 17025**.
- **OOS-5** Firma Apple Developer del binario. Procedura Gatekeeper manuale documentata.
- **OOS-6** Embedding / RAG / vector database. KB caricata come prompt. Nessun chunking automatico.
- **OOS-7** Caching delle chiamate LLM, retry intelligente, sandbox/permessi avanzati.
- **OOS-8** Provider LLM diversi da Ollama ed EUrouter. Niente OpenAI diretto, Anthropic diretto, gateway terzi. I due provider supportati restano i due dichiarati nello schema.
- **OOS-9** Più di un input per lancio. Logging strutturato JSON. Modifica dei file di input.
- **OOS-10** Validazione semantica automatica dell'output del prompt meta. Solo validazione di schema; la qualità del `trigger_prompt` resta giudizio umano (revisione manuale prima dell'uso).
- **OOS-11** Test di tenuta su context lungo. **Gestito separatamente dal CTO in chat manuale Ollama**, fuori dall'eseguibile e fuori dal budget pipeline. 1 gettone interno, non passa dal `/v-plan`.
- **OOS-12** Report finale di sprint + walkthrough conclusivo. **Gestito separatamente dal CTO**, fuori dal budget pipeline. 2 gettoni interni.
- **OOS-13** Aggiunta di nuove capability oltre le 7. Il prompt meta abilita Denis a esplorare candidati in autonomia tra sprint; la messa in produzione resta lavoro tecnico del CTO in Sprint 2.

## Open Questions

Domande aperte da chiarire **prima** del `/v-plan` o nella sua fase iniziale. Annotate per non perderle.

- ~~**OQ-1** parametri esatti Qwen~~ → **RISOLTA (2026-05-18, decisione CTO)**: **i parametri sono configurabili nel profilo** (campi `temperature`, `max_tokens`, `context_window`, `top_p`, `top_k`, `num_ctx`). Default Sprint 1: `temperature: 0.9`, `max_tokens: 8192`, `context_window: 128000`. `top_p`/`top_k`/`num_ctx` non vengono toccati esplicitamente in Sprint 1 (si lascia il default di Ollama/Qwen). I valori si **affinano durante lo shakedown del profilo `revisione`**: se l'output del Test 1 è qualitativamente OK, i parametri si propagano come default a tutti i profili. Override per singolo profilo solo se motivato da risultati osservati. (Nota tecnica: temperature 0.9 è alta per task ispettivi — scelta deliberata del CTO per dare a Qwen più libertà generativa; si valuterà al primo lancio se serve abbassare.)
- ~~**OQ-2** modello esatto Qwen~~ → **RISOLTA (2026-05-18, decisione CTO)**: modello = **`qwen3.6`** (stringa Ollama del registry locale; la quantizzazione e il nome esatto del tag verranno confermati al primo `ollama list` su ognuno degli hardware target — CTO Windows WSL2 con NVIDIA RTX 4090 24GB, Denis macOS Apple Silicon). Il nome del modello è un **campo del profilo**, configurabile per capability. Decisione di runtime: se `qwen3.6` non è disponibile, fallback al primo `qwen3*` presente nel registry locale, registrato nel frontmatter dell'output. Non blocco strutturale del piano.
- ~~**OQ-3** materiali nuovi C-G~~ → **RISOLTA (2026-05-18)**: dati **sintetici plausibili come default** preparati dal CTO sul modello dei pattern tipici di un laboratorio ISO 17025; se Denis consegna materiali reali anonimizzati prima dello shakedown di una capability, **li sostituiamo capability-per-capability**. Decisione coerente con il piano esperimenti.
- ~~**OQ-4** golden Test 2 non conservato~~ → **RISOLTA (2026-05-18)**: nessuna rigenerazione del Test 2. La validazione del profilo `rilievi` (FR-14) si basa esclusivamente sul **confronto col template `CAPA Pack` del CoWork KB sez. 15.2** (struttura, triade meccanismo/estensione/efficacia, distinzione correzione vs azione correttiva). Lo shakedown del motore poggia operativamente sul confronto con golden Test 1 della Capability A; la Capability B confronta col template di output. Trade-off accettato dal CTO.
- ~~**OQ-5** strict ordering vs scaffolding parallelo~~ → **RISOLTA (2026-05-18, decisione CTO)**: lo strict ordering A→B→C→D→E→F→G vale **solo per "esecuzione + valutazione di Denis"**, ovvero per la fase in cui si interpreta l'output del modello. Lo **scaffolding non-LLM** (struttura YAML del profilo, selezione dei `kb_files`, prima bozza del `trigger_prompt`, preparazione dei dati di input sintetici) **può procedere in parallelo**: è lavoro deterministico del CTO, non valida nulla del modello. Disciplina di validazione strict resta intatta sull'output. Razionale: tenere lo scaffolding serializzato sarebbe pura attesa e brucerebbe ore senza ridurre rischio.
- ~~**OQ-6** target platforms~~ → **RISOLTA (2026-05-18, decisione CTO)**: **due target di Sprint 1**: `darwin/arm64` (consegna Denis, macOS Apple Silicon) **e `linux/amd64`** (sviluppo CTO su Windows WSL2 con NVIDIA RTX 4090 24GB per Ollama accelerato). Cross-compile Go: trivial via `GOOS`/`GOARCH`. Conseguenza importante: la modalità interattiva NON può dipendere da `osascript` (macOS-only) — vedi OQ-7. `darwin/amd64` (Mac Intel) **non in scope** se non emerge bisogno specifico al checkpoint. Vedi NFR-2 aggiornato.
- ~~**OQ-7** modalità interattiva cross-platform~~ → **RISOLTA (2026-05-18, decisione CTO)**: data la necessità di girare anche su `linux/amd64` (OQ-6), `osascript` non è praticabile come unica modalità. **Si adotta una libreria TUI Go cross-platform** (default proposto: `charmbracelet/huh` per il flusso prompt sequenziale di selezione profilo / input / output; in alternativa `manifoldco/promptui` se `huh` si rivela troppo grosso). Selezione precisa della libreria delegata al `/v-plan`. Il flusso interattivo resta sequenziale (profilo → input → output). Su macOS il bundle `.app` apre Terminal e lancia la TUI; su Linux/WSL2 il binario gira direttamente in console. **Drag&drop di cartella su `.app` macOS** preserva l'input passandolo come argomento posizionale (feature macOS-native opportunistica, NON disponibile su Linux). Bindings Cocoa nativi e dialog `osascript` restano fuori scope Sprint 1. Vedi FR-10, FR-11 aggiornati.
- ~~**OQ-8** pt-analysis pre-processing deterministico~~ → **RISOLTA (2026-05-18, decisione CTO)**: **confermato — niente pre-processing deterministico in Sprint 1**, anche se Qwen sbaglia i calcoli. Razionale: la domanda contrattuale dello sprint è *"Qwen regge sulla capability"*; aggiungere pre-processing risponderebbe a una domanda diversa ("come rendere robusto un pipeline ibrido"). Il giudizio di Denis è il verifier vero della capability G; un esito *"non validata"* è informazione utile per il dimensionamento di Sprint 2 (dove la parte numerica potrà essere appoggiata a pre-processing deterministico), non un bug del codice di Sprint 1.

Le risposte a queste domande vanno annotate prima di chiudere il gate `/v-spec` (o subito dopo, prima di `/v-plan`).

## Decisioni sui test data (gate spec)

- **TD-1** Per le 5 capability nuove (C-G) si parte con dataset **sintetici plausibili** preparati dal CTO; sostituiti da reali anonimizzati di Denis se consegnati prima dello shakedown di quella capability. Vedi `.pipeline/proposed-test-data/README.md` per i pattern previsti.
- **TD-2** Per il profilo `rilievi` (B) **nessun golden file**, solo confronto col template `CAPA Pack` del CoWork KB 15.2.
- **TD-3** Per i test unitari/integrazione del motore (parsing, composizione prompt, validate) si usa uno **snapshot ridotto** della KB in `.pipeline/test-data/fixtures/kb-ispettore-minima/` (3-4 file). Gli scenari BDD dei profili `revisione` e `rilievi` usano invece la KB **reale completa** in `docs/piano_iniziale/materiali-dominio/KB-ispettore/` per fedeltà al prompt effettivo.
