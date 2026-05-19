# LabNexus — Guida utente

Per Denis (QM) e per chi userà l'eseguibile dopo Sprint 1.

## Cosa fa LabNexus

Esegue **una capability ispettiva alla volta** del modello AICertus. La capability viene scelta tramite un **profilo YAML** (`revisione`, `rilievi`, ecc.). Il modello AI gira **in locale** sul Mac, su Qwen 3 (Ollama). Nessun dato esce dalla macchina.

## Prima volta su un Mac nuovo

L'app non è firmata con certificato Apple Developer (lo sarà in Sprint 2). Alla prima apertura:

1. Apri il Finder, naviga a dove hai estratto `labnexus.app`
2. **Control-clic** sull'icona di `labnexus.app`
3. Scegli "Apri" dal menu
4. macOS chiede conferma — clicca "Apri" di nuovo
5. **Operazione una tantum**: i successivi lanci sono normali

Inoltre serve **Ollama** in esecuzione (`ollama serve`) con il modello `qwen3.6` installato (`ollama pull qwen3.6` la prima volta).

## Tre modi di lanciarlo

### 1. Doppio click (più semplice)

Doppio click su `labnexus.app`: si apre Terminal con una TUI sequenziale. La TUI ti chiede in ordine:

1. **Quale capability** vuoi eseguire? (lista a scorrimento con descrizione)
2. **Cartella di input**: dove sono i file da processare?
3. **Cartella di output**: dove vuoi che venga scritto il risultato?

Quando l'esecuzione termina, il file `.md` risultante è nella cartella di output.

### 2. Drag & drop di una cartella

Trascina una cartella sull'icona di `labnexus.app`: la TUI parte già con quella cartella come input. Ti chiede solo la capability e la cartella di output.

### 3. Riga di comando (per scripting)

```
labnexus run --profile revisione --input /path/Test1 --output /path/Out
```

> **Nota tecnica**: il "comando" `labnexus` nel PATH dovrebbe essere il **binario reale**, NON il wrapper del bundle `.app`. Se hai installato solo `labnexus.app`, usa il binario interno: `./labnexus.app/Contents/MacOS/labnexus-bin`. Lanciare il wrapper (`labnexus` dentro `.app/Contents/MacOS/`) da terminale è valido ma aprirà una **seconda** finestra Terminal — è il comportamento corretto per il doppio click dal Finder, scomodo se sei già in una shell.

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

In Sprint 1 sono validate sull'hardware locale; nuove capability possono essere aggiunte tra sprint via prompt meta (vedi Sprint 2).

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
- **"Ollama non raggiungibile"** → apri Terminal e lancia `ollama serve` in una finestra separata
- **"profile: kb_file ... non esiste"** → il profilo cita un file della KB che non c'è nella cartella `KB-ispettore`. Verifica di avere installato la KB completa
- **"context window superato"** → l'input è troppo grande per il modello. Riduci il numero di file di input
- **PDF non si legge** → il PDF è troppo complesso per il parser interno. Convertilo manualmente con `pandoc input.pdf -o input.md` e riprova
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
