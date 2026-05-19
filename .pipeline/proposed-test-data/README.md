# Proposed Test Data — Sprint 1 LabNexus

**DRAFT** — bozza di test data per la spec di Sprint 1. Va curata e approvata insieme al gate spec, poi formalizzata in `.pipeline/test-data/` durante `/v-test-scaffold`.

## Principio

Ogni scenario BDD in `.pipeline/proposed-features/` deve poter girare con dati deterministici. Quando i dati reali esistono (golden file dei Test 1 e 2) li usiamo come fonte. Quando non esistono (capability C-G), prepariamo **dati sintetici plausibili** sul modello descritto nel piano esperimenti, su input/schema indicati da Denis (o, in mancanza, ricostruiti dal CTO secondo i pattern tipici di un laboratorio ISO 17025).

Tutti i dati sono **realistici ma fittizi**. Nessun dato reale di laboratori esistenti, di persone identificabili o di fornitori reali entra nel repo.

## Struttura proposta

```
.pipeline/test-data/
├── README.md
├── fixtures/
│   ├── profili-validi/             # YAML profilo validi minimi
│   │   ├── revisione.yml           # versione di test del profilo
│   │   ├── rilievi.yml
│   │   └── minimal.yml             # profilo minimo per testare il motore in isolamento
│   ├── profili-invalidi/           # per gli scenari validate negativi
│   │   ├── trigger-corto.yml       # trigger_prompt < 50 caratteri
│   │   ├── kb-mancante.yml         # referenzia un kb_file inesistente
│   │   ├── provider-errato.yml     # provider="openai"
│   │   └── senza-trigger.yml       # campo obbligatorio mancante
│   └── kb-ispettore-minima/        # snapshot RIDOTTO della KB per i test di motore
│       ├── CLAUDE.md               # versione di test, breve
│       ├── come-pensa-un-ispettore.md
│       └── ...
├── scenarios/
│   ├── motore-cli/
│   │   ├── input-mix-formati/      # un file per ogni formato supportato
│   │   ├── input-annidato/         # per testare walk non ricorsivo
│   │   ├── input-vuoto/            # per EC-9
│   │   └── input-solo-non-supportati/
│   ├── revisione-test1/            # PUNTA al materiale reale (vedi sotto)
│   ├── rilievi-test2/              # PUNTA al materiale reale
│   ├── review-pack/                # SINTETICO da costruire
│   ├── audit-checklist/            # SINTETICO da costruire
│   ├── equipment-alert/            # SINTETICO da costruire
│   ├── competence-gap/             # SINTETICO da costruire
│   ├── pt-analysis/                # SINTETICO da costruire
│   ├── edge-pdf-rotto/             # un PDF intenzionalmente corrotto
│   └── edge-streaming-interrotto/  # mock provider script
└── golden/
    ├── revisione/                  # PUNTA a materiali-dominio/test-precedenti/test-1-output-approvato/
    └── (rilievi: assente, vedi OQ-4)
```

## Fonti per ciascun gruppo

| Scenario / Profilo | Origine dei dati | Note |
|---|---|---|
| `motore-cli/input-mix-formati` | CTO costruisce 6 file dummy (uno per formato) | Contenuto irrilevante, basta che sia parsabile |
| `motore-cli/input-vuoto` | cartella vuota | Trivial |
| `motore-cli/input-solo-non-supportati` | CTO crea 2-3 file `.jpg`/`.mov` placeholder | Trivial |
| `motore-cli/input-annidato` | CTO costruisce struttura "doc.md + sub/altro.md" | Per walk non ricorsivo |
| `revisione-test1` | **materiale di dominio reale** in `docs/piano_iniziale/materiali-dominio/test-precedenti/test-1-input/` | Già presente. NON duplicare: i test referenziano questo path. |
| `golden/revisione` | **materiale di dominio reale** in `docs/piano_iniziale/materiali-dominio/test-precedenti/test-1-output-approvato/` | Già presente. Approvato da Denis. |
| `rilievi-test2` | `docs/piano_iniziale/materiali-dominio/test-precedenti/test-2-input/ACIAA A1 rilievi.csv` | Già presente (1 file CSV) |
| `golden/rilievi` | **ASSENTE** (output Test 2 non conservato) | Vedi **OQ-4** in spec |
| `review-pack` | **SINTETICO** da CTO, su schema fornito o ricostruito | NC/audit/reclami/PT/equipment/KPI |
| `audit-checklist` | **SINTETICO** da CTO | risk register 15-25 voci + descrizione scope audit |
| `equipment-alert` | **SINTETICO** da CTO | scheda apparecchiatura + evento taratura |
| `competence-gap` | **SINTETICO** da CTO | matrice 10-15 × 8-12 + procedura/metodo nuovo |
| `pt-analysis` | **SINTETICO** da CTO | risultati PT con z-score + metodi + storico 2-3 anni |
| `edge-pdf-rotto` | PDF generato malformato | Trivial |
| `profili-invalidi/*` | CTO scrive a mano 4 YAML | Trivial |
| `kb-ispettore-minima` | CTO crea snapshot ridotto della KB per test unitari del motore | Per evitare di caricare l'intera KB ad ogni test unit |

## Domande aperte sui dati

Vedi domande in calce a questo doc — saranno consolidate insieme alle Open Questions della spec.
