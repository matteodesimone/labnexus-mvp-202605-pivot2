---
date: 2026-05-19
pipeline: bugfix
feature: app-bundle-doppio-click-no-output
stack: go-cli (Go 1.22 + bash + osascript)
tags: macos, bundle, gui, tdd, integration-test, applescript, escape
---

# Bugfix: doppio click su labnexus.app non apriva nulla

## Problem

Denis (e il CTO al smoke test post-deploy Fetta 1) faceva doppio click su `labnexus.app` dal Finder, ma niente succedeva. Nessuna finestra Terminal si apriva, nessun feedback all'utente. Causa root: `CFBundleExecutable=labnexus` puntava direttamente al binario Go, che lanciato da Launch Services non aveva un TTY → `tui.Run()` ritornava `ErrNonTTY` → il main usciva exit 2 in background senza che lo stderr fosse visibile. Bug HIGH, blocker della consegna.

## Solution

Pattern idiomatic macOS per CLI-in-`.app`: `CFBundleExecutable` punta a un **wrapper bash** che apre Terminal via `osascript` e lancia il binario reale al suo interno.

Layout finale:
```
labnexus.app/Contents/MacOS/
  labnexus      ← wrapper bash (1.5 KB, CFBundleExecutable)
  labnexus-bin  ← binario Go reale (Mach-O arm64, 8.7 MB)
```

Il wrapper supporta drag&drop di una cartella (`$1` posizionale → input pre-selezionato) e fa doppio escape (shell single-quoting + AppleScript string literal escape per `"` e `\`) per path con caratteri speciali.

## What Worked

- **TDD red phase rigorosa**: 5 test integration in `internal/bundle/` con `TestMain` che esegue `scripts/build-mac.sh` una sola volta + 5 assertion sulla struttura. Build tag `//go:build darwin` per non eseguire su Linux/CI senza Mach-O tools. Tutti i test fallivano col bug, tutti passavano dopo il fix.
- **Diagnosi rapida** (~5 min): ipotesi giusta al primo colpo (CFBundleExecutable + non-TTY → ErrNonTTY silenzioso) confermata dall'osservazione del bundle pre-fix.
- **Pattern macOS standard**, non workaround: gli utenti che mantengono CLI-in-.app (alacritty, ssh-agent helpers, ecc.) usano lo stesso wrapper.
- **Inline fix dell'unmasked failure** (Ollama timeout flaky 87s): regola framework "≤30 min → fix inline" applicata correttamente. Causa: macOS port 1 (tcpmux) + `OllamaProvider` senza HTTP timeout. Fix: `Timeout=90s` esplicito + endpoint test su IP RFC 5737 non-routable (`192.0.2.1`).
- **Review loop 2 ha catturato un edge case importante** (path con `"`): il loop 1 escapava solo `'` per shell, mancava l'escape AppleScript. Il loop 2 ha aggiunto `escape_for_osascript()` con doppio livello + nuovo test `TestBundle_WrapperEscapesApplescriptStringChars`.

## What Didn't Work

- **Lo scenario `@manual` "doppio click sul .app apre Terminal" non è MAI stato eseguito manualmente prima della consegna di Fetta 1**. Il framework non ha un checkpoint pre-`/v-deploy` che imponga la verifica degli scenari `@manual`. Risultato: bundle BROKEN consegnabile a Denis. Salvato solo dal smoke test post-deploy del CTO.
- **Il primo fix (loop 1) aveva escape incompleto**: solo `'` shell, niente `"`/`\` per AppleScript. Path con virgolette doppie avrebbero rotto `do script`. Catturato solo in review loop 2 — sintomo che la "shell injection paranoia" del code-reviewer va applicata anche su AppleScript injection.
- **Mancava HTTP timeout sull'OllamaProvider**: scelta originale "lascia default Go" si è rivelata fragile su sistemi (macOS port 1) che accettano TCP ma non rispondono. Avrebbe dovuto essere un finding HIGH del review Fetta 1, non è stato visto.

## Reusable Pattern

**"CLI in macOS .app bundle"** — pattern documentato in `.pipeline/skills/project/macos-cli-in-app-bundle.md` (vedi). Riassunto:
- `CFBundleExecutable` punta a un wrapper bash
- Wrapper usa `osascript` per aprire Terminal e lanciare il binario reale (`<name>-bin`)
- Doppio escape per path: shell `'` + AppleScript `"` `\`
- Drag&drop via `$1` posizionale
- Test integration con `TestMain` che esegue lo script + 6 assertion sulla struttura

**"Test integration di artefatti di build"** — pattern usato qui: `TestMain` esegue `bash scripts/build-X.sh` una sola volta in tmpdir/repo, i test successivi fanno assertion sul filesystem prodotto. Generalizzabile a qualsiasi build script (zip, .deb, .rpm, container image, ecc.).

## System Updates Applied

Nessuno applicato direttamente in questa cycle. Proposte (vedi richiesta CTO in chat):
- **Skill di progetto** `.pipeline/skills/project/macos-cli-in-app-bundle.md` (codifica pattern wrapper bash + osascript + test integration)
- **Framework promotion** `.pipeline/proposed-updates/2026-05-19-manual-scenarios-deploy-checkpoint.md` (gli scenari BDD `@manual` devono avere checkpoint CTO esplicito pre-`/v-deploy`, non saltabile — questo bug ne dimostra il valore)

## Tech Debt

Identificati 2 trade-off, entrambi accettati come "balanced":

1. **Test integration `internal/bundle/` modifica `labnexus.app/` nel repo root**: ogni `go test ./internal/bundle/` ricostruisce il bundle (rm -rf + build). Race possibile se altri processi/test toccano `labnexus.app/`. Mitigazione: nessun altro test/codice tocca quel path (BDD usano tmp dir per il binario). Trade-off accettato.
2. **HTTP timeout di 90s su OllamaProvider**: generoso per evitare false abort su `review-pack` (output lungo). Ma se Qwen è davvero hang, l'utente aspetta 90s prima del feedback. Acceptable per Sprint 1; valutare riduzione + retry pattern in Sprint 2 se diventa fastidio.

**Raccomandazione**: NON serve `/v-refactor`. Procedere col compound del bugfix e tornare a Fetta 2 plan.
