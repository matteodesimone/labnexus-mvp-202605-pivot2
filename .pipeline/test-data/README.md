# .pipeline/test-data/ — curated fixtures and scenarios for Sprint 1 tests

Questa cartella contiene i **fixture curati** e gli **scenario di test** usati da:

- Unit/integration tests (Go): caricati da test helper nei pacchetti `internal/...`
- BDD acceptance (godog): caricati dagli step in `features/`

Vedi `.pipeline/proposed-test-data/README.md` per il razionale completo (sintesi delle scelte TD-1, TD-2, TD-3 fatte al gate spec).

## Struttura

```
fixtures/
├── profili-validi/             # profili YAML validi (revisione, ecc.)
├── profili-invalidi/           # profili con errori intenzionali per testare validate
└── kb-ispettore-minima/        # snapshot ridotto KB per test unitari motore (TD-3)
scenarios/
├── motore-cli/
│   ├── input-mix-formati/      # 6 file dummy uno per formato supportato
│   ├── input-vuoto/            # cartella vuota
│   ├── input-annidato/         # walk non ricorsivo
│   └── input-solo-non-supportati/
├── edge-pdf-rotto/             # PDF intenzionalmente malformato (EC-1)
├── review-pack/
│   └── anno-2025-sintetico/    # Fetta 2 — Capability C (FR-15): pacchetto annuale per Management Review Pack
│       ├── nc-anno-2025.csv          # 15 NC dell'anno (12 chiuse, 3 in corso)
│       ├── audit-interni-2025.md     # esiti 3 audit interni (Q1/Q2/Q3)
│       ├── reclami-2025.csv          # 5 reclami con esito e tempi
│       ├── risultati-pt-2025.csv     # 10 circuiti PT con z-score
│       ├── stato-apparecchiature-2025.csv  # 12 apparecchiature con stato + metodi
│       └── kpi-2025.md               # 6 indicatori di prestazione vs target
├── audit-checklist/
│   └── risk-register-sintetico/  # Fetta 2 — Capability D (FR-16): risk register + scope audit
│       ├── risk-register.csv         # 18 voci di rischio per area ISO 17025
│       └── scope-audit.md            # scope dell'audit interno (oggetto, riferimenti, modalità)
└── equipment-alert/
    └── scheda-pmt-sintetica/   # Fetta 2 — Capability E (FR-17): scheda apparecchiatura + evento
        ├── scheda-apparecchiatura.md   # T-007 con metodi associati (PCM-01/04/09)
        └── evento.md                   # certificato di taratura rientrato con NC su intervallo
```

### Note sulle fixture di Fetta 2

Le 3 cartelle di Fetta 2 sono dati **sintetici plausibili** costruiti dal CTO secondo i pattern tipici
di un laboratorio ISO 17025 (decisione TD-1 della spec): NC dell'anno con causa radice e azione, audit
interni con rilievi proporzionati, PT con z-score realistici (≤|3|), apparecchiature con metodi
associati e tracciabilità. I nomi sono fittizi (Maria Rossi, Giulio Bianchi…) ma riconoscibili al
laboratorio. Se Denis consegna materiali reali anonimizzati prima dello shakedown di una capability,
queste fixture vanno sostituite capability-per-capability (decisione TD-1).

## Origine dei dati reali (non duplicati qui — referenze)

- **Test 1 input** (`revisione`): `docs/piano_iniziale/materiali-dominio/test-precedenti/test-1-input/`
- **Test 1 golden output** (`revisione`): `docs/piano_iniziale/materiali-dominio/test-precedenti/test-1-output-approvato/`
- **Test 2 input** (`rilievi`): `docs/piano_iniziale/materiali-dominio/test-precedenti/test-2-input/ACIAA A1 rilievi.csv`
- **Test 2 golden output**: ASSENTE (TD-2 — confronto solo col template CAPA Pack)
- **KB-ispettore completa**: `docs/piano_iniziale/materiali-dominio/KB-ispettore/` (usata negli scenari BDD A/B; per test unitari del motore si usa `fixtures/kb-ispettore-minima/`)
