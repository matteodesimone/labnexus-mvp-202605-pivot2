# Guida d'uso — Generare nuovi profili LabNexus con Claude

> **Obiettivo**: aggiungere nuove capability a `labnexus` senza dipendere dal CTO per ogni nuovo profilo. Tempo medio per profilo: 10-20 minuti.

## Cosa ti serve

- Un browser e l'account `claude.ai`.
- `labnexus` installato sul tuo Mac (lo zip già consegnato).
- L'eseguibile `labnexus` accessibile dal terminale (o dentro `labnexus.app/Contents/MacOS/labnexus-bin`).
- Il file `docs/meta-prompt-genera-profilo.md` (questo è incluso nello zip).

## Passi

### 1. Apri Claude

Vai su `https://claude.ai` e accedi.

### 2. Inizia una nuova conversazione e incolla il meta-prompt

Apri `docs/meta-prompt-genera-profilo.md`, seleziona TUTTO il contenuto (`Cmd+A`, `Cmd+C`) e incollalo come **primo messaggio** nella chat. Invia.

Claude risponderà con un messaggio breve ("Sei pronto. Aspetta che l'utente descriva la capability…") confermando di aver letto le istruzioni.

### 3. Descrivi la capability che vuoi

Nel **secondo messaggio**, scrivi cosa vuoi che il profilo faccia. Esempi:

- *"Voglio un profilo per verificare la conformità di un certificato di taratura."*
- *"Voglio un profilo che identifichi deviazioni dai limiti di accettabilità in un set di rapporti di prova."*
- *"Voglio un profilo per la qualifica iniziale di un fornitore critico."*

Più sei specifico, meno domande Claude ti farà.

### 4. Rispondi alle domande di chiarimento

Se la descrizione è generica, Claude ti farà 1-3 domande prima di produrre lo YAML. Tipicamente:

- Quali file ricevi tipicamente in input?
- Cosa deve produrre l'output? Una struttura specifica?
- Esistono template di riferimento?
- Ci sono vincoli di non-allucinazione (es. *cita solo metodi presenti nella scheda*)?

Rispondi con il minimo necessario. Non serve essere lunghi.

### 5. Salva lo YAML che Claude produce

Quando hai informazioni sufficienti, Claude ti produrrà un blocco markdown del tipo:

```
Ecco lo YAML per il profilo `<nome>`:

​```yaml
profilo: <nome>
...
​```

**Note sulle scelte**: ...

**Verifica dello schema**:
​```bash
labnexus validate <nome> --profiles-dir ./profili --kb-dir ./KB-ispettore
​```
```

Copia il **contenuto del blocco yaml** (senza i ` ```yaml `) e salvalo in `profili/<nome>.yml`.

### 6. Valida lo schema

Apri il terminale, vai nella cartella di `labnexus` (quella con `profili/` e `KB-ispettore/`) e lancia:

```bash
./labnexus validate <nome>
```

(oppure dentro il bundle `.app`: `./labnexus.app/Contents/MacOS/labnexus-bin validate <nome>`)

Output atteso: `schema OK` con exit code 0.

**Se vedi un errore**:

- *"trigger_prompt troppo corto"*: chiedi a Claude di espandere il `trigger_prompt`.
- *"kb_file X non esiste in ./KB-ispettore/"*: Claude ha inventato un path. Chiedigli di ri-selezionare dalla lista esatta nella sezione 4 del meta-prompt.
- *"provider X non ammesso"*: deve essere `ollama` o `eurouter`. Normalmente `ollama`.

Riportato l'errore a Claude, lui correggerà e ti darà il nuovo YAML. Sostituisci e ri-valida.

### 7. (Opzionale) Esegui un test rapido

Se vuoi testare il profilo con un input sintetico:

```bash
./labnexus check <nome> --input <cartella-test>
```

Questo fa un dry-run (parsing + stima token, senza chiamare il modello). Se vuoi vedere il prompt composto:

```bash
./labnexus check <nome> --input <cartella-test> --show-prompt
```

⚠ **Attenzione PII**: `--show-prompt` espone il contenuto dei file di input + KB. Non condividere l'output su canali non controllati (ticket, Slack, log condivisi).

### 8. Quando lanciare il profilo reale

```bash
./labnexus run --profile <nome> --input <cartella-input> --output <cartella-output>
```

L'output è un file markdown nella cartella scelta. Apri il file, leggi, valuta L1 strutturalmente (le sezioni ci sono? L'output è ben formato?), poi se OK passa a L2 (giudizio qualitativo tuo).

## Quando NON usare il meta-prompt

- Profili che richiedono **calcoli numerici esatti** (es. statistica di z-score, audit trail di formule): Sprint 1 non li copre. Sprint 2 introdurrà pre-processing deterministico per la parte numerica.
- Profili che richiedono **integrazione con LIMS o database**: out-of-scope Sprint 1.
- Profili che usano dati **non discutibili con un fornitore cloud esterno** come parte di test: usa il profilo solo on-device con Ollama; non condividere materiali reali in Claude per fini di scrittura prompt.

## Quando contattare il CTO

- Se Claude ti chiede di modificare la KB-ispettore (NON dovrebbe).
- Se ti sembra che ti servirebbe un nuovo file della KB (es. *"servirebbe una sezione su X"*): aprire un ticket col CTO; le modifiche alla KB sono iterazioni tue offline tra sprint.
- Se il profilo produce output di qualità ridotta su Qwen 3 locale ma sembra OK su Claude esterno: probabile mismatch trigger_prompt vs modello locale; ricalibrazione con CTO.
- Se vedi `provider eurouter non ammesso senza approval esplicito` durante un test: NON aggirare. Significa che il profilo prova a chiamare il cloud — discutere col CTO PRIMA di forzare.

---

**Riferimento tecnico**: lo schema completo del profilo + lista KB-ispettore è in `docs/meta-prompt-genera-profilo.md` sezione 3 e 4. Il README operativo è in `README.md`.
