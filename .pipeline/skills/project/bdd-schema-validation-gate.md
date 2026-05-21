---
slug: bdd-schema-validation-gate
created: 2026-05-21
source_cycle: Fetta 2 review loop 1 (H1+H2 fix)
applies_to: capability declarative con profilo YAML shippable (review-pack, audit-checklist, equipment-alert; e tutte le future Fetta 3 + meta-prompt)
---

# BDD schema-validation gate in install step

## Problema

Le capability declarative di LabNexus consistono in profili YAML (`profili/<name>.yml`) + selezione kb_files + trigger_prompt. Il "deliverable" è il file YAML stesso. Quando i test BDD usano uno stub minimal in tmp dir (per non dover copiare tutta la KB-ispettore), **la correttezza del profilo reale non è esercitata** — uno schema sbagliato (kb_file inesistente, trigger troppo corto, provider sbagliato) passerebbe i test verdi.

## Pattern

Nella step di install BDD per ogni capability (`installReviewPack`, `installAuditChecklist`, ecc.), eseguire **come precondizione** il subprocess `labnexus validate <name>` contro:
- `--profiles-dir <repoRoot>/profili` (il profilo reale shippato)
- `--kb-dir <repoRoot>/docs/piano_iniziale/materiali-dominio/KB-ispettore` (la KB reale)

Se il validate fallisce → la step BDD ritorna errore (RED) PRIMA che il run inizi. Lo scenario fallisce con messaggio diagnostico esplicito ("labnexus validate <name> ha fallito contro la KB reale: ...").

Implementazione canonica: `features/steps_test.go::installFromProject`:

```go
projectProfile := filepath.Join(repoRoot, "profili", name+".yml")
if _, err := os.Stat(projectProfile); err != nil {
    return c, fmt.Errorf("profilo di progetto non ancora implementato: %s mancante", projectProfile)
}
projectProfiliDir := filepath.Join(repoRoot, "profili")
projectKbDir := filepath.Join(repoRoot, "docs", "piano_iniziale", "materiali-dominio", "KB-ispettore")
validateCmd := exec.Command(binaryPath, "validate", name, "--profiles-dir", projectProfiliDir, "--kb-dir", projectKbDir)
if out, err := validateCmd.CombinedOutput(); err != nil {
    return c, fmt.Errorf("labnexus validate %s ha fallito contro la KB reale: %v\n%s", name, err, string(out))
}
// poi: scrivi stub minimal in tmp dir per il run del binary
```

Dopo la validation reale, scrivere uno stub minimal in `s.profiliDir` (perché il run usa `--kb-dir <tmp>` che ha solo `CLAUDE.md`).

## Cosa NON copre

Questo pattern valida lo **schema** del profilo reale (campi, tipi, kb_files exist, trigger ≥ 50 char) ma NON la **content composition** end-to-end (trigger_prompt reale + tutti i kb_files multipli concatenati nel prompt). Quest'ultimo è il "persistent HIGH accepted risk" documentato in `.pipeline/review.md` Fetta 2 e nel todo #009, mitigato dallo shakedown reale CTO + L2 Denis.

## Quando applicarlo

Ogni capability declarative dello sprint:
- ✅ Fetta 1: già coperto via `revisione`, `rilievi` (existing install steps)
- ✅ Fetta 2: implementato in `installFromProject` per review-pack/audit-checklist/equipment-alert
- ▶ Fetta 3: applicare a `competence-gap`, `pt-analysis`, e a `feat-meta`-generated profili (con audit umano del trigger_prompt PRIMA del test)

## Costo

Un subprocess `labnexus validate` per scenario, ~50-150ms. Trascurabile vs il valore di catturare schema regressions automaticamente. Già adottato in Fetta 2 BDD senza impatto perceptible sul tempo della suite.
