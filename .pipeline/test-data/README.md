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
└── edge-pdf-rotto/             # PDF intenzionalmente malformato (EC-1)
```

## Origine dei dati reali (non duplicati qui — referenze)

- **Test 1 input** (`revisione`): `docs/piano_iniziale/materiali-dominio/test-precedenti/test-1-input/`
- **Test 1 golden output** (`revisione`): `docs/piano_iniziale/materiali-dominio/test-precedenti/test-1-output-approvato/`
- **Test 2 input** (`rilievi`): `docs/piano_iniziale/materiali-dominio/test-precedenti/test-2-input/ACIAA A1 rilievi.csv`
- **Test 2 golden output**: ASSENTE (TD-2 — confronto solo col template CAPA Pack)
- **KB-ispettore completa**: `docs/piano_iniziale/materiali-dominio/KB-ispettore/` (usata negli scenari BDD A/B; per test unitari del motore si usa `fixtures/kb-ispettore-minima/`)
