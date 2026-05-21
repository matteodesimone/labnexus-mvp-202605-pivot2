# LabNexus — Guida essenziale (Denis)

Questa è la guida minima per usare LabNexus senza Terminal. Per dettagli avanzati c'è anche `README.md` (più tecnico).

---

## Cosa hai ricevuto

Estraendo lo zip ottieni una cartella con dentro:

| File / cartella | A cosa serve |
|---|---|
| `labnexus` | Il programma vero e proprio (binary Mac) |
| `labnexus.command` | **Quello che lanci tu con doppio click** |
| `labnexus.config.toml` | **Il file della chiave API** (vedi sotto) |
| `profili/` | I 7 profili shippati (revisione, rilievi, review-pack, audit-checklist, equipment-alert, competence-gap, pt-analysis) |
| `lavori/` | 7 cartelle con materiale di test già pronto, una per profilo |
| `KB-ispettore/` | La tua knowledge base (system context, NON modificare) |
| `meta-prompt-genera-profilo.md` | Prompt da incollare in claude.ai per generare nuovi profili (avanzato) |
| `guida-meta-prompt-denis.md` | Guida d'uso del meta-prompt (avanzato) |

---

## Setup iniziale (UNA VOLTA SOLA)

### Passo 1 — Sposta lo zip in una cartella stabile

Estrai lo zip in una cartella che NON sposterai più (es. `~/labnexus/`). Tutti i file devono restare insieme.

### Passo 2 — Metti la tua chiave EUROUTER

**Il file dove va la chiave si chiama `labnexus.config.toml`**, ed è nella cartella principale dello zip estratto (stesso livello di `labnexus.command`).

Per editarlo:

1. **Doppio click** sul file `labnexus.config.toml` nel Finder
2. macOS lo apre con **TextEdit** (o un editor di testo se ne hai un altro di default)
3. Cerca la riga:
   ```
   eurouter_api_key = ""
   ```
4. **Sostituisci la stringa vuota tra le doppie virgolette** con la tua chiave. Esempio:
   ```
   eurouter_api_key = "sk-eurouter-1234567890abcdefghij"
   ```
   ⚠ **Importante**:
   - DEVE essere tra le doppie virgolette (`"..."`)
   - NON aggiungere spazi prima o dopo il segno `=`
   - NON cancellare le altre righe del file
5. **Salva** col menu `File → Salva` (o `Cmd+S`). Chiudi TextEdit.

✓ Fatto. La chiave è memorizzata. Non devi più ri-editare questo file (a meno che la chiave non cambi).

### Passo 3 — Prima esecuzione (Gatekeeper macOS)

Il programma non è firmato Apple Developer, quindi la prima volta macOS chiede conferma:

1. **Control-clic** (o clic-destro) su `labnexus.command` nel Finder
2. Scegli **"Apri"** dal menu contestuale
3. macOS chiede conferma → clicca **"Apri"** ancora
4. Ripeti gli stessi 3 passi su `labnexus` (il binary). Anche questo deve essere autorizzato.

✓ Fatto, **solo la prima volta**. Le volte successive: doppio click normale.

---

## Uso quotidiano

### Lanciare una capability su un materiale di test (i 7 lavori preconfigurati)

1. **Doppio click su `labnexus.command`**. Si apre Terminal e parte una procedura guidata.

2. **Seleziona il profilo**: la procedura ti mostra una lista dei 7 profili shippati. Usa le **frecce su/giù** per scorrere, **invio** per confermare. Esempio: seleziona `revisione`.

3. **Indica la cartella di input**: la procedura ti chiede dove sono i file. Per i materiali di test forniti:
   - **Apri il Finder** in parallelo
   - **Trascina la cartella** `lavori/CAPABILITY A — Profilo revisione` (o quella corrispondente al profilo che hai scelto) dal Finder direttamente sulla finestra Terminal
   - Terminal incolla automaticamente il path. Premi **invio**.

4. **Indica la cartella di output**: la procedura ti chiede dove salvare il risultato. Suggerimento: usa la sotto-cartella `output` dentro la cartella di lavoro:
   - **Trascina la cartella** `lavori/CAPABILITY A — Profilo revisione/output/` su Terminal
   - Se la cartella `output/` non esiste ancora, digita il path manualmente (es. `lavori/CAPABILITY A — Profilo revisione/output`) — verrà creata automaticamente
   - Premi **invio**.

5. **Aspetta**. Vedrai messaggi tipo:
   ```
   [10:32:14] chiamata provider eurouter ...
   [10:32:14] attendo risposta dal modello (warmup può richiedere minuti)...
   [10:32:51] primo token ricevuto, generazione in corso...
   ```
   Poi il **testo del modello scorrerà a video in tempo reale**: vedi il modello che pensa. Tempo tipico: da 1 minuto a 8-10 minuti, dipende dal volume di input.

6. **Risultato finale**: la procedura termina con la riga `output: <path/al/file.md>`. Apri quel file `.md` con il tuo editor (Obsidian, TextEdit, Marked, ecc.) per leggere la bozza prodotta.

### Dove trovi i risultati

Ogni esecuzione produce DUE file nella cartella di output:

| File | Cosa contiene |
|---|---|
| `<timestamp>_<profile>_<input>.md` | **La bozza prodotta dal modello**. È quello che leggi tu. |
| `<timestamp>_<profile>_<input>.log` | Un file tecnico (audit trail). Lo ignori, serve a noi del team se devi segnalare un problema. |

### Provare un'altra capability

Ripeti i 5 passi sopra, scegliendo un altro profilo nella TUI (es. `rilievi`, `review-pack`, ecc.) e trascinando la cartella corrispondente in `lavori/`.

---

## Valutare il risultato (L2)

Quando hai letto la bozza `.md`, aggiungi la tua valutazione modificando il **frontmatter** (la sezione YAML in cima al file, tra i due `---`):

1. Apri il `.md` con TextEdit (doppio click)
2. In cima al file vedrai qualcosa come:
   ```yaml
   ---
   profilo: revisione
   modello: qwen3.6
   provider: eurouter
   data_esecuzione: 2026-05-22T10:32:14+02:00
   ...
   ---
   ```
3. **Aggiungi una riga** prima del secondo `---` con la tua valutazione:
   ```yaml
   valutazione_denis: validata
   ```
   I valori ammessi sono: `validata`, `validata_con_riserva`, `non_validata`.
4. **Salva** (Cmd+S).

Questo file resta nella cartella `lavori/<X>/output/`. Quando il team CTO raccoglie le valutazioni per il checkpoint cliente, legge questo campo da tutti i `.md`.

---

## Quando qualcosa va storto

| Sintomo | Soluzione |
|---|---|
| Doppio click su `labnexus.command` → Terminal si apre e si chiude subito | Probabilmente la chiave EUROUTER non è ancora settata. Riapri `labnexus.config.toml` e controlla la riga `eurouter_api_key = "..."`. La chiave deve essere tra virgolette doppie e non vuota. |
| Doppio click → "macOS non può verificare lo sviluppatore" | È Gatekeeper. Rifai il Passo 3 del setup: control-clic → Apri. |
| Esecuzione che si blocca su "attendo risposta dal modello" per più di 5-10 minuti | Probabilmente la chiave EUROUTER è invalida o eurouter è giù. Chiudi Terminal (Cmd+W), verifica la chiave, riprova. |
| Errore "EUROUTER_API_KEY mancante" | La chiave nel file `labnexus.config.toml` è ancora `""` vuota. Apri il file, inserisci la chiave, salva. |
| Output `.md` con testo bizzarro o vuoto | Probabilmente il modello ha avuto un problema. Riprova lanciando lo stesso lavoro: gli output sono salvati con timestamp diversi, non si sovrascrivono. |
| Voglio aggiungere un nuovo lavoro (non i 7 pre-shippati) | Per ora chiedi al CTO. La procedura "concierge" per aggiungere nuovi lavori in autonomia è documentata in `docs/USER-GUIDE.md` e `guida-meta-prompt-denis.md` (entrambe richiedono un po' di confidenza col Finder + editor di testo). |

---

## Cose che ti chiedi probabilmente

**Q: Devo lasciare il Mac acceso durante l'esecuzione?**
Sì. L'esecuzione gira sul server eurouter ma il programma sul tuo Mac aspetta la risposta. Mac in stand-by interrompe.

**Q: Posso lanciare due capability in parallelo?**
No, una alla volta. Lanciane uno, aspetta la fine, lancia il prossimo.

**Q: I miei dati di input vanno in cloud?**
Sì, vanno al server EUrouter (gateway EU-GDPR su `api.eurouter.ai`). Stefano Fiorina ha autorizzato formalmente questo flusso. Per dati strettamente confidenziali (es. dati personali del laboratorio), valuta caso per caso prima di lanciarli.

**Q: Quanto costa ogni esecuzione?**
Dipende dal volume di input (numero token). Il programma stampa la stima all'inizio (es. "35659 token stimati"). Eurouter fattura al token; verifica il tuo piano.

**Q: Posso usare il modello locale invece del cloud?**
Sì, ma serve hardware adeguato (GPU ≥ 24GB VRAM consigliato) e una procedura di setup separata. Chiedi al CTO se vuoi questa opzione.

---

## Contatti

Per problemi tecnici o domande: **Matteo De Simone (CTO)** — matteo.ds@gmail.com

Per validazione L2 e quality metric: discussione nel prossimo checkpoint Sprint 1.
