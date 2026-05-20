---
id: 006
priority: P2
status: ready
source: bug
created: 2026-05-20
file: internal/profile/profile.go, internal/output/output.go
bug_ref: .pipeline/bugs/frontmatter-default-non-applicato.md
---

# Cablare `output.frontmatter_default` del profilo nell'output writer

## Issue

Tutti i 5 profili shippable (`revisione`, `rilievi`, `review-pack`, `audit-checklist`, `equipment-alert`) dichiarano nello YAML un campo `output.frontmatter_default` con chiavi/valori (es. `tipo: management_review_pack`, `stato: bozza_da_validare_qm`, `profilo_labnexus: review-pack`, `locale: true`) che dovrebbero comparire nel frontmatter del file di output.

In esecuzione reale: `grep -rn frontmatter_default internal/ cmd/` → **0 match**. La struttura non è dichiarata in `Profile`, non viene letta dal writer, nessuno scrive quei default. Discrepanza fra documentazione e implementazione (violazione Layer 0 "documentazione è codice"). Affligge tutte e 5 le capability, non solo Fetta 2.

## Fix

In `internal/profile/profile.go`:
```go
type Output struct {
    FrontmatterDefault map[string]string `yaml:"frontmatter_default"`
}
type Profile struct {
    // ...campi esistenti
    Output Output `yaml:"output"`
}
```

In `internal/output/output.go` (writer del file di output):
1. Prima di scrivere il frontmatter generato dall'engine, merge delle chiavi `profile.Output.FrontmatterDefault`.
2. Precedenza: se collisione, le chiavi generate dall'engine vincono (timestamp, durata, modello effettivo, ecc.).
3. Unit test che verifica la presenza dei default per ognuno dei 5 profili.

Re-abilitare in `features/profili-fetta2.feature` lo step già scritto (e rimosso temporaneamente): `E il frontmatter contiene i default del profilo "<X>"`. Lo step impl + mappa `fetta2ExpectedDefaults` esistono già in `features/steps_test.go` (review loop 1).

## Context

- Smascherato dalla M4 assertion di `/v-review` Fetta 2 loop 1 (vedi `.pipeline/review.md`).
- Aggiornare anche per i profili Fetta 1 (`revisione`, `rilievi`) con default + test.
- Convertito in todo P2 da `/v-triage` 2026-05-20.
- Stima fix: 1-2 gettoni infra (~30min-1h CTO).
- Buon candidato per **`/v-bugfix`** indipendente prima di Fetta 3 o subito dopo L2 Denis di Fetta 2.
