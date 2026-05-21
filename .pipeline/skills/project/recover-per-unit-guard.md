# Skill: Recover-Per-Unit Guard per librerie panic-prone

**Created**: 2026-05-21 in sub-fetta 1.5.B (compound)
**Trigger**: integri una lib esterna mature ma instabile (panic su edge cases) e vuoi continuare il processing senza crash.

## Il problema

Sub-fetta 1.5.B doveva parsare file `.xls` reali di Denis con `github.com/extrame/xls`. La lib funziona per la maggior parte delle row, ma panica su indici interni inconsistenti (es. file `M132.00_Rapporto Proficiency Testing.xls` panica su una row su 18). Senza guard, l'intero processo crasha.

Patterns Sprint 1 (EC-1 "warning + skip") richiedono: il parser ritorna error chiaro, il chiamante logga warning e procede col file successivo. Ma con panic la chain è interrotta.

## Il pattern

**Doppio recover stratificato**:

1. **Recover outer (file/risorsa level)**: avvolge l'intera operazione `parseXxx`. Trasforma panic in `error` chiaro per il chiamante. Ritorna `""` + error. Permette al chiamante di applicare EC-1 ("warning + skip").

2. **Recover inner (unit level)**: avvolge la singola unità di processing (row, sheet, entry). Skippa l'unità che panica, ma continua con la successiva. NON ritorna error — lascia l'output parziale.

### Codice (Go example)

```go
// parseXls.go — parser per file Excel 97-2003 OLE binary
func parseXls(path string) (s string, err error) {
    // OUTER recover: trasforma panic in error per il chiamante
    defer func() {
        if r := recover(); r != nil {
            s = ""
            err = fmt.Errorf("input: parse .xls %q panicked: %v", path, r)
        }
    }()
    wb, err := xls.Open(path, "utf-8")
    if err != nil { return "", fmt.Errorf("input: open .xls %q: %w", path, err) }
    var out strings.Builder
    for i := 0; i < wb.NumSheets(); i++ {
        sh := wb.GetSheet(i)
        if sh == nil { continue }
        appendSheetText(&out, sh)
    }
    return out.String(), nil
}

func appendSheetText(out *strings.Builder, sh *xls.WorkSheet) {
    max := int(sh.MaxRow)
    for r := 0; r <= max; r++ {
        appendRowSafely(out, sh, r) // <- INNER recover dentro
    }
}

func appendRowSafely(out *strings.Builder, sh *xls.WorkSheet, r int) {
    // INNER recover: skippa la row che panica, continua col next
    defer func() { _ = recover() }()
    row := sh.Row(r)
    if row == nil { return }
    // ... scrivi celle della row
}
```

### Convenzioni critiche

1. **Naming**: il wrapper inner ha `Safely` nel nome (`appendRowSafely`, `processItemSafely`) per chiamare attenzione al panic-catching.
2. **Outer logga, inner skippa silente**: l'outer mette l'errore nel return chain. L'inner skippa per default. Se serve logging delle row che hanno fatto panic per debugging, aggiungi un parametro `log *runlog.Logger` all'inner.
3. **Outer cattura ANCHE i panic dell'inner**: se l'inner panica, l'inner recover li gestisce. Ma se panica fuori dal recover (es. inside `sh.Row(r)` chiamato prima del recover), l'outer è la rete di sicurezza.
4. **Documenta esplicitamente** il pattern nei commenti: "lib X può panicare su Y. Recover outer + inner per defense-in-depth".

## Anti-pattern

- **Solo outer recover**: tutta la sheet fallisce su una row → output parziale degraded.
- **Solo inner recover**: panic fuori dalla row (es. nell'iterazione `wb.GetSheet`) crasha il processo intero.
- **Recover globale a livello di `main`**: troppo broad. Cattura anche bug del NOSTRO codice, mascherando regressioni.
- **`_ = recover()` senza commento**: futuro maintainer non sa perché. Aggiungi sempre `// inner recover: lib X panic su row indici inconsistenti, vedi commento parseXls`.

## Trade-off

- **Pro**: robusto, no crash su file reali con quirks. Continua processing.
- **Con**: errori silent al livello inner — bug regression nel NOSTRO codice (non lib esterna) potrebbe essere mascherato. Mitigare con test unit dedicati per le row in cui il parsing è critico.
- **Con**: pattern non idiomatico Go (recover è "last resort"). Documenta chiaramente perché è necessario.

## Precedente concreto

Sub-fetta 1.5.B LabNexus — `internal/input/xls.go` con `extrame/xls`. Decision point step 4: testato la lib su `data/CAPABILITY G/M132.00_Rapporto Proficiency Testing.xls`, panic su row interna nil. Risolto con doppio recover in ~15 min, no split fetta necessario.

Vedi `.pipeline/solutions/2026-05-21-sub-fetta-1.5.B-config-toml-parser-nuovi.md` e `internal/input/xls.go`.

## Quando NON applicare

- Lib stable (es. stdlib Go): non panica per design, no recover necessario.
- Codice nostro (no lib esterna): se panica, è un bug — fix it, no mask.
- Performance-critical hot path: il `defer recover()` ha overhead non trascurabile su loop ad alta frequenza.

## Generalizzabilità (verso framework promotion)

Pattern universal: vale per qualsiasi linguaggio con exception/panic semantics (Python try/except per-item, Rust catch_unwind, JVM try/catch). Candidato a `skills/framework/` come pattern defensive di-default per integrazione lib esterne instabili.
