# DRAFT test data per Sprint 1.5 — fixture aggiuntive necessarie

Questo file documenta le **fixture aggiuntive** che Sprint 1.5 introduce, oltre a quelle
già esistenti per Sprint 1 (Fetta 1-3).

## Sample artefatti DRAFT (in questa cartella)

| File | Per FR | Note |
|---|---|---|
| `SAMPLE-labnexus-config.toml` | 1.5.B / FR-11 | Esempio config master. Va shipped nello zip al root con `eurouter_api_key = ""` placeholder. Denis lo edita. |
| `SAMPLE-labnexus-job-metadata.toml` | 1.5.C / FR-26 | Esempio `_labnexus.toml` per cartella `lavori/<X>/`. Denis aggiunge uno per ogni cartella nuova (o lo lascia generare da un futuro `labnexus init-job`). |
| `SAMPLE-revisione.toml` | 1.5.B / FR-12, FR-13, FR-16 | Esempio profilo TOML post-refactor. Mostra campi opzionali ereditati + `trigger_prompt_file` indirezione. Pattern per i 7 profili shippati. |

## Fixture nuove necessarie per testare 1.5.B (parser nuovi formati)

Per testare i parser `.xls`, `.doc`, `.rtf` servono file campione **sintetici plausibili**:

### `.xls` (Excel 97-2003 OLE)

Path proposto: `.pipeline/test-data/fixtures/parser-formati/sample.xls`

Contenuto: 2 sheet, ciascuna con ~5-10 righe × 3-5 colonne di dati realistici laboratorio
(es. sheet "PT-Risultati" con codice + parametro + z-score + valutazione; sheet "Metodi"
con codice + descrizione).

**Generazione**: il CTO crea un file `.xls` reale con LibreOffice Calc (salvataggio in
formato Excel 97-2003) o usando uno script Go con la libreria scelta (`shakinm/xlsReader`
ha funzioni di scrittura? Verifica). In alternativa: prendere uno dei file `.xls` reali
in `data/` e ridurlo a 5-10 righe sintetiche.

### `.doc` (Word 97-2003 OLE)

Path proposto: `.pipeline/test-data/fixtures/parser-formati/sample.doc`

Contenuto: ~1-2 pagine di testo plain con paragrafi, eventuale tabella semplice.
Realtà del laboratorio: scheda personale, lettera fornitore, procedura.

**Generazione**: LibreOffice Writer → "Salva con nome" → formato "Word 97-2003 (.doc)".
Oppure: prendere `M115.00_Scheda personale_Alessandro Bidossi.doc` di Denis e
anonimizzarlo (sostituire nome reale con "Mario Rossi" o equivalente sintetico).

### `.rtf` (Rich Text Format)

Path proposto: `.pipeline/test-data/fixtures/parser-formati/sample.rtf`

Contenuto: stesso pattern dei `Prompt_INPUT_Rev.00.rtf` di Denis — TASK + INPUT FORNITI
+ ATTIVITÀ + OUTPUT ATTESO. Ma sintetico (caso generico).

**Generazione**: TextEdit Mac → "Salva con nome" → formato RTF. O semplicemente un
RTF file scritto a mano (RTF è formato testuale ASCII-friendly).

## Fixture nuove necessarie per testare 1.5.C (auto-discovery)

Per testare `labnexus jobs` auto-discovery servono cartelle stub:

### `lavori/<X>/_labnexus.toml` test cases

Path proposto: `.pipeline/test-data/scenarios/concierge-autodiscovery/`

Strutture da scaffold (e da committare):

```
.pipeline/test-data/scenarios/concierge-autodiscovery/
├── lavori-happy-path/
│   ├── CAPABILITY A — Profilo revisione/
│   │   ├── _labnexus.toml         # profile=revisione, trigger_prompt_file=Prompt.rtf
│   │   ├── Prompt.rtf
│   │   └── doc-stub.txt
│   ├── CAPABILITY B — Profilo rilievi/
│   │   ├── _labnexus.toml         # profile=rilievi
│   │   └── rilievi-stub.csv
│   └── ...
├── lavori-fallback-convention/
│   └── Cartella ad hoc — Profilo audit-checklist/
│       └── input-stub.md           # NO _labnexus.toml, convention naming fallback
├── lavori-skipped-no-metadata/
│   └── cartella-random/
│       └── stuff.txt               # né _labnexus.toml né convention naming → skip
├── lavori-malformed-toml/
│   └── CAPABILITY F — Profilo competence-gap/
│       └── _labnexus.toml           # TOML syntax invalido
└── lavori-trigger-prompt-file-missing/
    └── CAPABILITY G — Profilo pt-analysis/
        └── _labnexus.toml           # trigger_prompt_file=missing.rtf, ma missing.rtf non esiste
```

## Note di triage per /v-test-scaffold

- Le fixture **parser formati** sono prerequisito per i BDD scenari `.xls`, `.doc`, `.rtf` (1.5.B
  proposed-features). Vanno create prima del red phase.
- Le fixture **concierge auto-discovery** sono prerequisito per i BDD scenari di auto-discovery
  e edge cases (1.5.C proposed-features).
- I `SAMPLE-*.toml` in questa cartella sono **template per scrittura finale** (1.5.B implementation):
  diventano la versione shipped nello zip (master config) + il pattern per `_labnexus.toml`
  generation manuale (Denis lo scrive a mano, brevissimo) o futuro `labnexus init-job` (Sprint 2).

## Source dei dati reali di Denis (non duplicati qui — referenze)

- **Cartella `data/`** del CTO contiene gli scenari reali di Denis (7 cartelle CAPABILITY A-G).
  Sprint 1.5.C **rinomina** questa cartella in `lavori/` come parte del refactor.
- **File reali dentro data/**: `.xls` (M132.00, M134.01, M114.00, M147.00, M143.01, ecc.),
  `.doc` (M115.00 Scheda personale), `.rtf` (Prompt_INPUT_*.rtf per ogni capability),
  `.xlsx`, `.docx`, `.pdf`, `.txt`. Sono i banchi di prova reali per i parser.

**Per la `/v-test-scaffold`**: prima di scrivere i BDD scenari, scaffold le fixture stub
(con dati sintetici plausibili) nelle path proposte sopra. Le fixture reali di `data/`
restano la "ground truth" finale ma sono coperte dalle esecuzioni reali (shakedown CTO +
L2 Denis), non dai BDD test automatici.
