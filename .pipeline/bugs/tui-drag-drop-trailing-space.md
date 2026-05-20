---
severity: medium
status: fixed
created: 2026-05-19
fixed: 2026-05-19
source: developer (CTO smoke test post-bugfix-5, UX TUI con drag&drop Finder)
fix: "Aggiunta cleanPath(s) (TrimSpace + Trim quote singole/doppie + TrimSpace) chiamata all'ingresso di validateExistingDir e validateNonEmpty in internal/tui/tui.go. Il TrimSpace post-form.Run() esisteva già ma veniva applicato DOPO la validate, troppo tardi. cleanPath gestisce 3 pattern reali di paste/drag: trailing space puro, single quote + space (drag&drop macOS Finder→Terminal), double quote."
test: "internal/tui/tui_test.go: TestCleanPath_DragDropFinderPatterns (8 sub-test su input/output normalizzati), TestValidateExistingDir_HandlesTrailingSpace, TestValidateExistingDir_HandlesSingleQuotes. Tutti GREEN."
---

# Bug: TUI non fa trim dell'input → drag&drop Finder produce path con trailing space → "cartella non esiste"

## Reported

- **Feature/Area**: FR-10 (flusso TUI sequenziale) — `internal/tui/tui.go`
- **Trovato da**: CTO smoke test post-bugfix-5, mentre testava il riavvio binario
- **Sprint impact**: MEDIUM. UX-only: il binario funziona se l'utente digita il path manualmente. Ma il pattern drag&drop dal Finder è quello che Denis userà di default, quindi blocca il workflow tipico.

## Description

**Cosa succede**:
- Utente apre TUI (doppio click `labnexus.app` o lancio diretto binario)
- Alla domanda "Cartella di input", trascina la cartella da Finder dentro Terminal
- Terminal incolla il path circondato da single quote + uno spazio finale: `'/Users/denis/Test1' `
- `validateExistingDir(s)` chiama `isExistingDir(s)` che fa `os.Stat(s)` sulla stringa raw
- `os.Stat("'/Users/denis/Test1' ")` ritorna ENOENT (la cartella non si chiama esattamente così)
- TUI mostra: "la cartella '...' non esiste" e fa rieseguire la prompt

**Cosa dovrebbe succedere**: il TrimSpace + strip quote viene fatto PRIMA della validate, così l'input riconosce il path effettivo.

## Steps to Reproduce

1. Lancia binario in TTY (terminale o `.app`)
2. Alla prompt "Cartella di input", digita un path con uno spazio finale (oppure usa drag&drop su macOS)
3. Osserva: validate fallisce con "la cartella non esiste" anche se la cartella esiste

## Environment

- macOS (drag&drop da Finder a Terminal genera questo pattern di default)
- Stesso pattern possibile su Linux con file manager + paste su xterm
- Indipendente dal modello LLM e dal flusso di esecuzione

## Analysis

**Causa diretta (alta confidenza)**:

In `internal/tui/tui.go`, `Run()` applicava `strings.TrimSpace` SOLO dopo `form.Run()`:

```go
if err := form.Run(); err != nil { return nil, err }
sel.InputDir = strings.TrimSpace(sel.InputDir)  // ← troppo tardi
sel.OutputDir = strings.TrimSpace(sel.OutputDir)
```

Ma le validate (`validateExistingDir`, `validateNonEmpty`) vengono chiamate DURANTE la digitazione/conferma del campo, prima che `form.Run()` ritorni. Quindi un input `"/path/ "` veniva validato sulla stringa raw e falliva.

**Fix**:

1. Nuova funzione `cleanPath(s)`: `TrimSpace → Trim quote singole+doppie → TrimSpace`. Idempotente, simmetrica.
2. `validateExistingDir(s)` e `validateNonEmpty(s)` chiamano `cleanPath(s)` come prima riga.
3. Il post-form `strings.TrimSpace` diventa `cleanPath` per coerenza (catch-all di sicurezza).

**Affected files**:
- `internal/tui/tui.go` (cleanPath nuova + 2 validator + post-form)
- `internal/tui/tui_test.go` (3 nuovi test, totale ~15 test)

**Risk of fix**: **molto basso**. Modifica isolata a 2 funzioni di validazione + 1 helper nuovo. Nessun impatto sul flow di esecuzione, sui provider, sul parsing. Backward compatible: input puliti continuano a funzionare invariati.

## Note di scope (prove-out)

Bug fix budget-free. Stima fix: ~15 min (effettivi).

**Pattern emergente — 6° bug consecutivo** post-deploy Fetta 1, di cui:
- 5 trovati al primo smoke test reale
- 1 (questo) trovato durante il riavvio post-fix-5

Tutti riconducibili a UN unico gap: nessuno scenario `@manual` BDD eseguito CTO pre-deploy. Da formalizzare con FORZA al `/v-compound` finale di Fetta 1 (già abbozzato in `.pipeline/proposed-updates/2026-05-19-manual-scenarios-deploy-checkpoint.md`, da rinforzare).
