---
slug: hallucination-prevention-fixture-set
created: 2026-05-21
source_cycle: Fetta 2 review loop 1 (H3 fix)
applies_to: tutte le capability che citano riferimenti tecnici verificabili (PCM-NN, codici NC, sezioni ISO, ID risk register, nomi fornitori, codici apparecchiature)
---

# Hallucination-prevention via fixture-extracted set

## Problema

La quality metric di Denis (NFR-8) ha "allucinazioni" come killer #1. Le assertion BDD ingenue (es. regex `PCM-\d+` per cercare un riferimento "qualsiasi") passano anche con allucinazioni plausibili (un `PCM-99` inventato che non esiste nella scheda). L'assertion deve verificare:
- (a) almeno UN riferimento citato (no output generico vacante);
- (b) TUTTI i riferimenti citati provengono dal set ammesso (no invenzioni).

## Pattern

Estrarre il set ammesso DALLA FIXTURE che il modello vede in input. Poi confrontare il body dell'output:

```go
func fixtureAllowedPCMCodes() (map[string]bool, error) {
    scheda := filepath.Join(repoRoot, ".pipeline/test-data/scenarios/equipment-alert/scheda-pmt-sintetica/scheda-apparecchiatura.md")
    b, _ := os.ReadFile(scheda)
    allowed := extractPCMCodes(string(b)) // map[string]bool con PCM-01, PCM-04, PCM-09
    return allowed, nil
}

func checkMetodiOnlyFromSet(body string, allowed map[string]bool) error {
    cited := extractPCMCodes(body)
    if len(cited) == 0 { return errors.New("nessun PCM citato") }
    for code := range cited {
        if !allowed[code] { return fmt.Errorf("metodo %s NON presente nella scheda (allucinazione)", code) }
    }
    return nil
}
```

L'estrattore (`extractPCMCodes`) usa una regex condivisa (`PCM-\d+`) sia sul body output sia sulla fixture; il confronto è insieme-vs-insieme.

## Generalizzabile a

Pattern applicabile ogni volta che l'output deve citare elementi da un set FINITO e VERIFICABILE presente nell'input:

| Capability | Set ammesso da | Regex estrattore |
|---|---|---|
| equipment-alert (FR-17) | scheda apparecchiatura.md | `PCM-\d+` |
| audit-checklist (FR-16) | risk-register.csv | colonna `id` → `R-\d+` |
| pt-analysis (FR-G) | risultati-pt.csv | colonna `codice_pt` → `PT-\d{4}-\d+` |
| review-pack (FR-15) | citazioni di NC | colonna `codice_nc` → `NC-\d{4}-\d+` |
| competence-gap (FR-F) | matrice competenze | nomi tecnici (caso speciale: persone) |

Per la matrice competenze, attenzione: i "nomi propri" non sono un set deterministico facilmente regex-estraibile. Soluzione: il fixture YAML/JSON-formatta i nomi come colonna deterministica.

## Cosa NON cattura

- **Allucinazioni di valori semantici** (es. PCM-01 citato correttamente ma con una descrizione errata). Quello è L2 di Denis (NFR-8).
- **Riferimenti mancanti** (output cita 1 PCM su 3 disponibili). Il pattern accetta "almeno uno" — se vuoi forzare "tutti", aggiungi assertion `len(cited) == len(allowed)` ma valuta se è realmente requirement (review-pack potrebbe legittimamente citare solo i PCM critici, non tutti).
- **False positives sull'output reale di Qwen** se il modello produce variant spelling (`PCM 01` con spazio invece di `PCM-01`). Pattern assume normalizzazione lessicale.

## Cost di adozione

~30 LoC per nuova capability: definire `fixture<Cap>AllowedSet()`, riusare `checkOnlyFromSet(body, allowed)`. Riusabile la regex se gli ID seguono lo stesso formato (es. tutti `XXX-\d+`).
