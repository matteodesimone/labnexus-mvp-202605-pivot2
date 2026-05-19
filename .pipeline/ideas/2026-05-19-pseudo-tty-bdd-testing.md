---
date: 2026-05-19
status: planned
source: compound-learning-fetta1
priority: medium
tags: bdd, tui, testing, sprint-2, sprint-3
---

# Pseudo-TTY BDD testing per la TUI interattiva

## Idea

Testare la TUI interattiva di `labnexus` (lanciato senza argomenti → flusso `huh` sequenziale) usando uno pseudo-terminal in BDD acceptance. Sostituirebbe i 5 scenari attualmente marcati `@manual` in `features/motore-modo-interattivo.feature` con scenari completi e automatizzati.

## Razionale

Al momento gli scenari TUI sono `@manual` perché:
- `huh` richiede `term.IsTerminal(stdin.Fd()) == true` per funzionare
- `go test` invocato da CI gira con stdin non-TTY → `huh` esce con error
- Risultato: i 5 scenari sono testati solo manualmente dal CTO/Denis prima della consegna

Con pseudo-TTY (es. `creack/pty` in Go) il test harness può:
1. Allocare un pty
2. Lanciare il binario `labnexus` con stdin/stdout/stderr collegati al pty
3. Scrivere "tasti" simulati nella TUI (frecce, invio, testo)
4. Leggere l'output e fare assertion su righe specifiche del frame TUI

## Esempi di scenari recuperabili

- "Lancio senza argomenti apre la TUI con selezione profilo" — lanci, leggi frame, verifica testo "Scegli capability" presente
- "Dopo aver scelto il profilo, la TUI chiede la cartella di input" — manda invio sulla selezione, leggi frame, verifica testo
- "Un argomento posizionale = cartella esistente è usato come input" — lanci con arg, verifica che la TUI parta dalla selezione profilo (skip step input)
- "TUI funziona identicamente su darwin/arm64 e linux/amd64" — eseguire la stessa BDD su entrambi i target in CI

## Costo stimato

- Aggiungere dep `creack/pty` (pure Go, ~150 KB)
- Helper TUI in `features/helpers_test.go`: `(s *scenarioState) startTUI(args []string)` + `s.sendKeys(seq string)` + `s.readFrame() string`
- Aggiornare i 5 scenari rimuovendo `@manual` e implementando step
- Stima: ~2-3 ore di lavoro CTO

## Quando affrontarla

- **Sprint 2 / Sprint 3 / cleanup**: non blocca Sprint 1; ma se entrano altri developer o se la TUI cresce (Sprint 2 può aggiungere watcher su cartella SGQ), la copertura diventa strategica
- **Pre-requisito**: feature mature della TUI; refactor della TUI in Sprint 2 invaliderebbe i test pseudo-TTY

## Vedi anche

- `features/motore-modo-interattivo.feature` — scenari attualmente `@manual`
- `internal/tui/tui.go` — `Run()` con `isInteractiveTTY()` check
- `.pipeline/solutions/2026-05-19-fetta1-scooter.md` — emersione dell'osservazione

## Notes

(vuoto)
