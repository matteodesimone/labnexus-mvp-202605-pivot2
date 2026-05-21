# LabNexus — Guida utente

Per Denis (QM) e per chi userà l'eseguibile dopo Sprint 1.

## Cosa fa LabNexus

Esegue **una capability ispettiva alla volta** del modello AICertus. La capability viene scelta tramite un **profilo YAML** (`revisione`, `rilievi`, ecc.). Il modello AI è **Qwen 3** via **EUrouter** (gateway cloud EU-GDPR su `api.eurouter.ai`) per default — Sprint 1.5 post-pivot-3. I dati di input transitano a EUrouter; il cliente (Stefano Fiorina) ha formalmente approvato questo flusso. DPA contrattuale è gestito lato cliente.

**Opzione locale** (per chi dispone di hardware adeguato, es. workstation con GPU ≥ 24 GB VRAM): si può forzare il provider Ollama locale con `--provider ollama` (CLI) o `LABNEXUS_PROVIDER=ollama` (env). In quel caso il dato resta sulla macchina, ma è responsabilità dell'utente verificare che `ollama serve` sia attivo e il modello `qwen3.6` installato.

## Prima volta su un Mac nuovo

Sprint 1.5 (post-pivot-3): nessun bundle `.app`. Il deliverable è un binary `labnexus` standalone + un piccolo launcher `labnexus.command` per il doppio click dal Finder.

Setup iniziale:

1. Estrai lo zip in una cartella di tua scelta.
2. Apri Terminal nella cartella dello zip (oppure naviga col `cd`).
3. (Solo se hai scompattato lo zip via Finder) Rimuovi il quarantine flag di macOS:
   ```bash
   xattr -d com.apple.quarantine labnexus labnexus.command
   ```
   In alternativa: control-clic sull'icona di `labnexus` (e `labnexus.command`) nel Finder → "Apri" → conferma. Operazione una tantum.
4. **Sprint 1.5.B** (in arrivo): edita `labnexus.config.toml` per inserire la tua `eurouter_api_key`. **Sprint 1.5.A** (corrente): l'API key resta come env var temporanea (`export EUROUTER_API_KEY=...`) finché 1.5.B non introduce il config master.

Nessuna installazione di Ollama richiesta (post-pivot-3: il default è eurouter cloud EU-GDPR).

## Due modi di lanciarlo

### 1. Doppio click su `labnexus.command` (più semplice)

Doppio click su `labnexus.command` dal Finder: si apre Terminal con una TUI sequenziale. La TUI ti chiede in ordine:

1. **Quale capability** vuoi eseguire? (lista a scorrimento con descrizione)
2. **Cartella di input**: dove sono i file da processare?
3. **Cartella di output**: dove vuoi che venga scritto il risultato?

Quando l'esecuzione termina, il file `.md` risultante è nella cartella di output.

### 2. Riga di comando (per scripting o utenti CLI)

```
./labnexus run --profile revisione --input /path/Test1 --output /path/Out
```

Oppure invoca la TUI con un input pre-selezionato:

```
./labnexus /Users/denis/Test1
```

Il path posizionale viene preso come input; la TUI chiede solo capability + output.

Comandi accessori:
- `labnexus list` — elenca i profili installati con descrizione
- `labnexus describe revisione` — dettagli del profilo (provider, modello, file della KB usati)
- `labnexus check revisione --input /path/Test1` — **dry-run**: verifica parsing + stima token, NON chiama il modello
- `labnexus validate revisione` — controlla che lo YAML del profilo sia valido

## Le 7 capability dello Sprint 1

| Sigla | Profilo | Cosa fa |
|---|---|---|
| A | `revisione` | Revisione documentale dopo cambio di norma |
| B | `rilievi` | Gestione rilievi Accredia (CAPA Pack per ogni rilievo) |
| C | `review-pack` | Management Review Pack annuale (13 sezioni) |
| D | `audit-checklist` | Checklist d'audit da risk register |
| E | `equipment-alert` | Equipment alert per scadenze tarature |
| F | `competence-gap` | Gap competenze da matrice + procedura nuova |
| G | `pt-analysis` | Analisi z-score PT |

In Sprint 1 sono validate sull'hardware locale. **Nuove capability** puoi aggiungerle tra sprint in autonomia tramite il **prompt meta** (Sprint 1 Fetta 3, FR-20): apri `claude.ai`, incolla `docs/meta-prompt-genera-profilo.md`, descrivi il task, segui le domande di chiarimento, salva lo YAML in `profili/` e lancia `labnexus validate <nome>`. Guida operativa step-by-step: `docs/guida-meta-prompt-denis.md`.

## Cosa vedo durante l'esecuzione

Dopo che hai scelto capability, input e output, sul terminale appaiono linee di log come:

```
[10:32:14] chiamata provider eurouter (https://api.eurouter.ai/v1/chat/completions) ...
[10:32:14] attendo risposta dal modello (warmup può richiedere minuti su modelli grandi)...
[10:32:24] ...ancora in attesa del primo token (10s elapsed)
[10:32:34] ...ancora in attesa del primo token (20s elapsed)
[10:32:51] primo token ricevuto (TTFT 37s), generazione in corso...
  > 142 token, 23s elapsed (6.1 tok/s)        ← aggiornata live
[10:38:47] streaming completato: 1842 token in 6m33s (4.7 tok/s)
[10:38:47] output: /Users/denis/Out/2026-05-19_103214_revisione_Test1.md
```

Le righe `...ancora in attesa` significano che il modello deve essere allocato lato server (cold start su EUrouter, o warmup di RAM se hai forzato Ollama locale). È normale soprattutto la prima richiesta dopo qualche minuto di inattività, o su modelli grandi (qwen3.6 36B = 23 GB). Una volta arrivato il primo token, la riga `> X token` si aggiorna in-place ogni mezzo secondo mostrando token cumulati, tempo trascorso e velocità (tok/s).

Se non vedi nessun aggiornamento per più di 2-3 minuti dopo il primo token, c'è davvero un problema (vedi sezione "Quando qualcosa va storto").

## Come è fatto l'output

Ogni esecuzione produce **un file markdown** in `<output>/<timestamp>_<profilo>_<input>.md`. Inizia con un blocco YAML di metadati:

```yaml
---
profilo: revisione
modello: qwen3.6
provider: ollama
data_esecuzione: 2026-05-19T08:45:00+02:00
durata_secondi: 42.3
token_stimati: 1234
file_input: [PG_RISK_LAB_Rev_00.docx, DE0779_RT_08rev03.pdf, RT-08-rev.05.pdf]
stato: completato
---

# corpo dell'output del modello…
```

Dopo che hai valutato l'output, puoi aggiungere al frontmatter `valutazione_denis: validata` / `validata con riserva` / `non validata` per il report di sprint.

## Quando qualcosa va storto

- **macOS non apre l'app** → Gatekeeper, vedi sopra ("Prima volta su un Mac nuovo")
- **"Ollama non raggiungibile" + `context deadline exceeded`** → due possibilità:
  - Ollama davvero down: apri Terminal e lancia `ollama serve` in una finestra separata
  - Ollama non risponde con gli HTTP headers entro il TTFB timeout (30 min di default). Possibile causa: modello molto grande in warmup, o server bloccato. **Lo streaming del body è illimitato** (può durare ore senza essere cancellato), il timeout protegge solo dal TTFB. Se serve aumentare il TTFB: `LABNEXUS_HTTP_TIMEOUT=3600 ./labnexus run ...` (3600s = 1 ora)
  - Verifica veloce: `curl http://localhost:11434/api/tags` deve rispondere con la lista dei modelli
- **Output troncato a metà frase** → non è un timeout (il body streaming è illimitato post-headers). Possibili cause: il modello ha emesso `done:true` prematuramente (raro), Ollama ha chiuso la connessione, oppure `max_tokens` raggiunto. Aumenta `max_tokens` nel profilo YAML se serve output più lungo.
- **"profile: kb_file ... non esiste"** → il profilo cita un file della KB che non c'è nella cartella `KB-ispettore`. Verifica di avere installato la KB completa
- **"context window superato"** → l'input è troppo grande per il modello. Riduci il numero di file di input
- **PDF non si legge** → il PDF è troppo complesso per il parser interno. Convertilo manualmente con `pandoc input.pdf -o input.md` e riprova
- **Il PDF passa la lettura ma l'output del modello è confuso** → il testo estratto dal PDF è incompleto (probabile su PDF scansionati, multi-colonna fitto, formule). Esegui `labnexus check <profilo> --input <dir> --show-prompt` per vedere il testo estratto; se è sbagliato, converti il PDF in `.md` manualmente prima
- **Output strano** → segnala al CTO; potrebbe essere un problema del modello su quella capability specifica

## Dati e privacy

- I file di input vengono letti **in sola lettura**; non vengono modificati né copiati altrove
- L'output è scritto **solo** nella cartella che hai indicato
- Niente dati lasciano il Mac (provider default `ollama` = locale)
- L'app non conserva alcuna copia: nessun database, nessuna cache, nessun upload

## Limitazioni note

- Una esecuzione = una capability (no orchestrazione automatica, no watcher: sono Sprint 2)
- Una sola finestra alla volta (no esecuzioni parallele)
- ISO 17034 non coperta in Sprint 1 (solo ISO 17025)
- Il modello `pt-analysis` può sbagliare i calcoli numerici: l'output va sempre verificato

## Contatti

Per problemi: CTO Matteo De Simone.
