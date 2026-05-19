---
severity: high
status: fixed
created: 2026-05-19
fixed: 2026-05-19
source: developer (CTO smoke test post-deploy Fetta 1)
fix: "scripts/build-mac.sh ora produce DUE file in labnexus.app/Contents/MacOS/: (1) `labnexus` = wrapper bash che apre Terminal via osascript e lancia il binario reale; (2) `labnexus-bin` = il binario Go Mach-O arm64. CFBundleExecutable resta `labnexus` ma punta al wrapper, non al binario. Drag&drop di una cartella sull'icona .app preservato via `$1` posizionale del wrapper. Fix collaterale: OllamaProvider.doRequest ora ha `http.Client.Timeout=90s` (evita hang lunghi su endpoint sospetti) + scenari BDD `Ollama non in ascolto` usano IP RFC 5737 (192.0.2.1) invece di 127.0.0.1:1 per fail-fast deterministico."
test: "internal/bundle/bundle_test.go (5 test, build tag darwin): TestBundle_CFBundleExecutableIsWrapperScript, TestBundle_RealBinaryAtLabnexusBin, TestBundle_WrapperLaunchesTerminalViaOsascript, TestBundle_WrapperIsExecutable, TestBundle_InfoPlistDeclaresCFBundleExecutable. TestMain costruisce il bundle una volta via scripts/build-mac.sh; i 5 test fanno assertion sulla struttura prodotta."
---

# Bug: doppio click su `labnexus.app` non apre nulla e l'app sembra non partire

## Reported

- **Feature/Area**: FR-10 (modalità interattiva TUI) + FR-11 (bundle `.app` macOS) — deliverable di Fetta 1 (`dist/labnexus-sprint1-darwin-arm64.zip`)
- **Trovato da**: CTO durante smoke test del deliverable post-deploy
- **Sprint impact**: blocker per la consegna a Denis — la modalità d'uso "doppio click" è documentata in `README.md`, `docs/USER-GUIDE.md` e nel piano sprint come il modo primario di utilizzo non-CLI

## Description

**Cosa succede**:
- L'utente fa doppio click su `labnexus.app` dal Finder
- **Niente**: nessuna finestra Terminal, nessuna TUI, nessun messaggio di errore visibile, l'icona del Finder non risponde

**Cosa dovrebbe succedere** (vedi FR-11 nella spec + `motore-modo-interattivo.feature` scenario "doppio click sul .app apre Terminal e avvia la TUI"):
- Si apre una finestra Terminal
- Nella finestra parte la TUI sequenziale di FR-10 (selezione capability → input → output)
- La finestra resta aperta a fine esecuzione per permettere di leggere log/output

## Steps to Reproduce

1. Estrai `dist/labnexus-sprint1-darwin-arm64.zip` in una cartella
2. Vai nel Finder a quella cartella
3. Doppio click su `labnexus.app`
4. **Osserva**: nessuna finestra appare; se controlli `Activity Monitor`, il processo `labnexus` non risulta in esecuzione (è già uscito)

Ambiente verificato: macOS Apple Silicon, Ollama attivo con `qwen3.6` installato (irrilevante per questo bug — il problema accade prima che il binario possa contattare Ollama).

## Environment

- macOS Apple Silicon (`darwin/arm64`)
- Bundle non firmato (Gatekeeper procedure documentata, già accettata dall'utente)
- `labnexus.app/Contents/Info.plist`:
  - `CFBundleExecutable=labnexus` → punta direttamente al binario Go `labnexus.app/Contents/MacOS/labnexus`
- Binario: ~8.7 MB, Mach-O `arm64`

## Analysis

**Likely cause (alta confidenza)**:

Quando macOS Launch Services lancia un `.app` bundle via doppio click, esegue `CFBundleExecutable` **senza terminale collegato** (`os.Stdin` e `os.Stderr` non sono TTY). La sequenza nel binario è:

1. `main()` → `cobra.Execute()` → `rootRunE()` (no args)
2. `rootRunE` → `tui.Run(profiliDir, preInput)`
3. `tui.Run` chiama `isInteractiveTTY()`:
   ```go
   return term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stderr.Fd()))
   ```
4. Lanciato da Finder: entrambi → `false` → ritorna `ErrNonTTY`
5. Il main esce con exit 2 dopo aver scritto su `stderr` "modalità interattiva non disponibile in non-TTY — usare labnexus run con i flag..."

**Ma**: lo `stderr` non è visibile per l'utente perché **non c'è una finestra Terminal aperta**. Il binario è terminato in background in ~50ms. L'utente vede solo "niente".

Il design del bundle è incompleto: `CFBundleExecutable` punta al binario nudo, ma il binario non sa "auto-aprire Terminal" — assume di girare già dentro un TTY. È un classico **chicken-and-egg** dei bundle macOS che lanciano CLI: serve un wrapper.

**Affected files**:
- `scripts/build-mac.sh` (genera `Info.plist` con `CFBundleExecutable=labnexus`)
- `internal/tui/tui.go` (`Run` + `isInteractiveTTY`)
- `cmd/labnexus/main.go` (`rootRunE` non gestisce il caso "lanciato dal Finder")

**Risk of fix**: **basso** (cambio isolato a `scripts/build-mac.sh` + eventualmente piccola logica in main). Pattern standard macOS, ben documentato.

## Fix proposto (verifica con CTO prima di applicare)

**Opzione A — Wrapper script come `CFBundleExecutable`** (idiomatic, Recommended):

1. In `scripts/build-mac.sh`, costruisci due file dentro `labnexus.app/Contents/MacOS/`:
   - `labnexus-bin` → il binario Go reale (rename del current `labnexus`)
   - `labnexus` → wrapper bash:
     ```bash
     #!/usr/bin/env bash
     # CFBundleExecutable wrapper: apre Terminal e lancia il binario reale.
     DIR="$(cd "$(dirname "$0")" && pwd)"
     osascript <<EOF
     tell application "Terminal"
       activate
       do script "'${DIR}/labnexus-bin'"
     end tell
     EOF
     ```
2. `Info.plist` resta con `CFBundleExecutable=labnexus` — ora punta al wrapper
3. Drag&drop: il wrapper deve gestire `$1` (la cartella trascinata) e passarla come argomento posizionale al binario reale

**Opzione B — Auto-relaunch in Terminal dal binario Go**:

Quando `ErrNonTTY` viene rilevato in `rootRunE`, prima di uscire chiama `exec.Command("open", "-a", "Terminal", os.Args[0])`. Più complesso (loop infinito se Terminal non è installato; non gestisce drag&drop pulito).

**Decisione preferita**: **A**, perché è il pattern canonico macOS per CLI-in-app-bundle.

## Test plan

Aggiungere al bundle un test **automatizzabile** (anche se il "doppio click" reale resta `@manual`):

1. Test unit: verifica che `labnexus.app/Contents/MacOS/labnexus` esista come **file di testo** (script) con shebang `#!/usr/bin/env bash`, NON come binario Mach-O
2. Test unit: verifica che `labnexus.app/Contents/MacOS/labnexus-bin` esista come binario Mach-O arm64
3. Test integration: invoca il wrapper con `LABNEXUS_TEST_MODE=1` env var (da implementare) che fa skip dell'`osascript` e lancia direttamente — verifica che riceva gli args correttamente
4. Manuale (resta `@manual`): doppio click reale dal Finder + drag&drop di una cartella

**Scenario BDD da aggiornare** in `features/motore-modo-interattivo.feature`:
- "bundle .app contiene l'eseguibile arm64 e Info.plist valido" → da estendere: contiene **wrapper + binario** + Info.plist
- "doppio click sul .app apre Terminal e avvia la TUI" → resta `@manual` ma diventa testabile come "wrapper script chiama osascript e binario"

## Note di scope (prove-out)

Bug fix budget-free per default (vedi `prove-out.md`). Stima fix: ~30-60 min (script wrapper + Info.plist + un nuovo test unit + ribuildare zip).

**Se questo bug fosse stato scoperto in `/v-review` Fetta 1 sarebbe stato HIGH e bloccante per `/v-deploy`.** È stato missato perché lo scenario `@manual` "doppio click apre Terminal" NON è mai stato eseguito manualmente prima della consegna. **Learning per Fetta 2/3**: i scenari `@manual` devono avere un checkpoint esplicito CTO prima del `/v-deploy`, non saltabile.
