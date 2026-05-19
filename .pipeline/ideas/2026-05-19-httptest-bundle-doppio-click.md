---
date: 2026-05-19
status: planned
source: compound-learning-bugfix-app-bundle
priority: low
tags: bdd, macos, ci, automation, sprint-2-or-later
---

# Test automatizzato "doppio click reale sul .app"

## Idea

Sostituire il scenario `@manual` "doppio click sul .app apre Terminal e avvia la TUI" con un test CI automatizzato che simula Launch Services + osascript headless. Eliminerebbe il checkpoint manuale (vedi proposal framework `manual-scenarios-deploy-checkpoint`) per questa categoria di scenari.

## Razionale

Il bugfix `app-bundle-doppio-click-no-output` ha mostrato che il scenario `@manual` era completamente saltato pre-consegna. Anche con il checkpoint framework proposto (rubber-stamping rischio), un test automatizzato è strettamente migliore. Costo: la maggior parte della simulazione macOS GUI non è documentata pubblicamente.

## Approcci possibili

1. **`open` + `osascript` headless**: macOS supporta `open -a labnexus.app` da CLI. Il test integration potrebbe:
   - `exec.Command("open", "labnexus.app")` per simulare l'apertura
   - Attendere ~2s che Terminal si apra
   - Verificare via AppleScript "tell application System Events" che esista una finestra Terminal con titolo contenente "labnexus-bin"
   - Pulizia: chiudere quella finestra Terminal
   Complessità media; richiede `osascript --version` + permission del Sistema (Accessibilità) → problematico in CI strict.

2. **Stub di `osascript` via PATH manipulation**: nel test, prepend a `PATH` un dir con un `osascript` script-mock che logga gli argomenti ricevuti. Il wrapper bash chiama quel mock. Il test verifica che il mock sia stato chiamato con la `do script` corretta contenente `labnexus-bin`.
   Complessità bassa; non testa il vero `osascript`, ma verifica il wrapper. Non simula davvero il doppio click.

3. **Test sintattico AppleScript**: il wrapper genera la stringa `do script "<cmd>"` → estrarla con `bash -x` mode + `grep`. Verificare via `osascript -e 'on run argv ... end run'` che la sintassi sia parsable. Non lancia Terminal.

4. **Container macOS in CI**: praticamente impossibile (GitHub Actions macOS runner sono VM, no nested macOS containers; le simulazioni Launch Services sono in user space e richiedono UI).

## Raccomandazione

**Approccio 2 (stub osascript) + 3 (sintassi AppleScript)** combinati. Copre:
- Il wrapper chiama osascript con gli args giusti (stub PATH)
- La stringa AppleScript è sintatticamente valida (parser AppleScript)

L'approccio 1 (real `open`) resta `@manual` per la "vera consegna" — è l'unico che testa Launch Services reali. Ma 2+3 catturano il 80% dei bug del wrapper senza GUI.

## Costo stimato

- Approccio 2 (PATH stub): ~1h. File `internal/bundle/wrapper_dispatch_test.go` + helper `setupOsascriptStub(t)`.
- Approccio 3 (parser sintattico): ~30min. Test che invoca `osascript -e "<stringa estratta dal wrapper>"` con un `do script` no-op.

## Quando affrontarla

- **Sprint 2 / Sprint 3 / cleanup**: non blocca. Il checkpoint framework (proposal) copre il problema con coverage manuale. Sostituibile con questo test quando ha senso investire l'ora.
- **Pre-requisito**: framework proposal `manual-scenarios-deploy-checkpoint` accettata (così il checkpoint manuale è già obbligatorio e questo test ne diventa l'upgrade automatizzato).

## Vedi anche

- `.pipeline/solutions/2026-05-19-bugfix-app-bundle-doppio-click.md`
- `.pipeline/proposed-updates/2026-05-19-manual-scenarios-deploy-checkpoint.md`
- `.pipeline/skills/project/macos-cli-in-app-bundle.md`
- `internal/bundle/bundle_test.go` — test attuali (struttura statica, non simulano doppio click)

## Notes

(vuoto)
