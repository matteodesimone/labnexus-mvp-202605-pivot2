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
├── equipment-alert/
│   └── scheda-pmt-sintetica/   # Fetta 2 — Capability E (FR-17): scheda apparecchiatura + evento
│       ├── scheda-apparecchiatura.md   # T-007 con metodi associati (PCM-01/04/09)
│       └── evento.md                   # certificato di taratura rientrato con NC su intervallo
├── competence-gap/
│   └── matrice-sintetica/      # Fetta 3 — Capability F (FR-18): matrice competenze + procedura nuova
│       ├── matrice-competenze.csv      # 12 tecnici × 10 metodi (autorizzato/non-autorizzato/in-qualifica)
│       └── procedura-nuova.md          # PCM-12 IPA sedimento da introdurre Q1-2026
├── pt-analysis/
│   └── risultati-anno-2025/    # Fetta 3 — Capability G (FR-19): PT 2025 + metodi + storico per trend
│       ├── risultati-pt-2025.csv       # 12 circuiti PT con z-score (1 questionabile, 1 non soddisfacente)
│       ├── elenco-metodi.md            # 9 metodi accreditati + 1 in procedimento (PCM-12)
│       └── storico-2023-2024.csv       # 9 z-score storici sui metodi peggiori per trend statistico
└── feat-meta/
    └── casi-inventati/         # Fetta 3 — feat-meta (FR-22): YAML prodotti da Claude esterno via meta-prompt
        ├── README.md                   # istruzioni operative per Matteo (3 casi semplice/medio/ambiguo)
        └── caso-{1-semplice,2-medio,3-ambiguo}.yml  # da committare DOPO esecuzione manuale Claude
```

### Note sulle fixture di Fetta 2

Le 3 cartelle di Fetta 2 sono dati **sintetici plausibili** costruiti dal CTO secondo i pattern tipici
di un laboratorio ISO 17025 (decisione TD-1 della spec): NC dell'anno con causa radice e azione, audit
interni con rilievi proporzionati, PT con z-score realistici (≤|3|), apparecchiature con metodi
associati e tracciabilità. I nomi sono fittizi (Maria Rossi, Giulio Bianchi…) ma riconoscibili al
laboratorio. Se Denis consegna materiali reali anonimizzati prima dello shakedown di una capability,
queste fixture vanno sostituite capability-per-capability (decisione TD-1).

### Note sulle fixture di Fetta 3

- **competence-gap**: matrice 12 tecnici × 10 metodi PCM con valori `{autorizzato, non-autorizzato, in-qualifica}`. Riusa i nomi sintetici di Fetta 2 (Maria Rossi/Giulio Bianchi/Anna Conti/Luca Ferri + altri 8 plausibili) per continuità narrativa. La procedura nuova introduce PCM-12 IPA sedimento (metodo realmente in roadmap accreditamento per molti laboratori di chimica acque).
- **pt-analysis**: estende e differenzia i risultati PT di review-pack (Fetta 2) per coprire un trend 2023-2025 sui metodi peggiori (PCM-08 tensioattivi, PCM-09 metalli Cd) — necessario per FR-19 "trend storico". 1 PT non soddisfacente (z=3.2) per esercitare le azioni immediate fuori soglia.
- **feat-meta/casi-inventati/**: la cartella ospita i 3 YAML prodotti DA CLAUDE ESTERNO via meta-prompt FR-22. Al test-scaffold contiene solo `README.md` con istruzioni. I 3 `caso-*.yml` vengono committati DOPO esecuzione manuale di Matteo (FUORI pipeline AI). Il test Go `internal/profile/profile_meta_cases_test.go` cicla quei file e li valida contro KB reale.

## Origine dei dati reali (non duplicati qui — referenze)

- **Test 1 input** (`revisione`): `docs/piano_iniziale/materiali-dominio/test-precedenti/test-1-input/`
- **Test 1 golden output** (`revisione`): `docs/piano_iniziale/materiali-dominio/test-precedenti/test-1-output-approvato/`
- **Test 2 input** (`rilievi`): `docs/piano_iniziale/materiali-dominio/test-precedenti/test-2-input/ACIAA A1 rilievi.csv`
- **Test 2 golden output**: ASSENTE (TD-2 — confronto solo col template CAPA Pack)
- **KB-ispettore completa**: `docs/piano_iniziale/materiali-dominio/KB-ispettore/` (usata negli scenari BDD A/B; per test unitari del motore si usa `fixtures/kb-ispettore-minima/`)
