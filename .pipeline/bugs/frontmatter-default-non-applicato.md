---
slug: frontmatter-default-non-applicato
severity: MEDIUM
status: open
opened_at: 2026-05-20
opened_by: /v-review loop 1 — unmasked pre-existing failure
fetta_context: Fetta 2 (capability C/D/E), ma il bug è cross-fetta (anche Fetta 1)
discovered_via: M4 assertion BDD (rimossa dagli scenari Fetta 2 per rispettare scope plan)
---

# Bug — `output.frontmatter_default` del profilo non è scritto nell'output

## Sintomo

I profili YAML dichiarano nel campo `output.frontmatter_default` un set di chiavi/valori che dovrebbero comparire nel frontmatter del file di output prodotto da `labnexus run`:

```yaml
output:
  frontmatter_default:
    tipo: equipment_alert
    stato: bozza_da_validare_qm
    profilo_labnexus: equipment-alert
    locale: true
```

In esecuzione reale (verificato via fake Ollama + binario costruito da `TestMain`), l'output **non contiene queste chiavi** nel suo frontmatter. Sono presenti solo le chiavi prodotte direttamente dall'engine (`profile`, `modello`, `provider`, `durata`, `timestamp`, ecc.).

## Riproduzione

```bash
go test ./features/... -godog.paths=features/profili-fetta2.feature  # con M4 assertion attiva
```

Il test fallisce su tutti e 3 gli scenari Fetta 2 al passo "il frontmatter contiene i default del profilo «X»" con errore `frontmatter privo della chiave 'tipo'`.

## Analisi

`grep -rn "frontmatter_default\|FrontmatterDefault" internal/ cmd/` → **0 match**. La struttura `Profile.Output.FrontmatterDefault` non è dichiarata in `internal/profile/profile.go`, e il writer in `internal/output/` non ha logica per leggere/serializzare questo campo.

Il campo è presente in tutti i profili (`revisione.yml`, `rilievi.yml`, `review-pack.yml`, `audit-checklist.yml`, `equipment-alert.yml`) e documentato negli standard di progetto come parte dello schema. È **declarazione senza implementazione**.

## Impatto

- **Tracciabilità per Denis**: Denis riconosce un output L1-validato dal frontmatter (es. `tipo: capa_pack` distingue un rilievo da una revisione). Senza questi tag, il consumo manuale degli output L1 è più rumoroso.
- **Convenzione standard violata**: il file `02-project.md` documenta `output.frontmatter_default` come parte dello schema profilo. Discrepanza fra documentazione e implementazione = bug Layer 0 ("la documentazione è codice").
- **Cross-fetta**: affligge tutti i profili, non solo Fetta 2.

## Fix proposto

Aggiungere in `internal/profile/profile.go`:
```go
type Output struct {
    FrontmatterDefault map[string]string `yaml:"frontmatter_default"`
}
type Profile struct {
    // ...
    Output Output `yaml:"output"`
}
```

In `internal/output/output.go`, prima di scrivere il frontmatter generato dall'engine, **merge** delle chiavi di `profile.Output.FrontmatterDefault` (con precedenza alle chiavi generate dall'engine se conflitti). Aggiungere unit test che verifica la presenza delle chiavi default per ognuno dei 5 profili shippable.

**Stima**: 1 gettone infra (≤ 30 min CTO).

## Note di triage

- **Smascherato in `/v-review` loop 1 di Fetta 2** (2026-05-20) tramite M4 assertion. L'assertion è stata rimossa dai 3 scenari Fetta 2 perché il fix engine è scope-creep rispetto al plan ("Fetta 2 = nessuna nuova logica del motore"). L'helper code (`extractFrontmatter`, `frontmatterContainsKV`, `fmContieneDefaultDelProfilo`, mappa `fetta2ExpectedDefaults`) **rimane in `features/steps_test.go`** + relativi unit test in `assertions_fetta2_test.go` per essere riutilizzati quando il bug viene chiuso.
- Al fix di questo bug, **ri-abilitare** lo step `il frontmatter contiene i default del profilo "<X>"` nei 3 scenari Fetta 2 (e considerarne uno equivalente per Fetta 1).
