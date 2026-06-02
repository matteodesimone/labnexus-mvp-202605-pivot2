---
severity: high
status: fixed
created: 2026-06-02
fixed: 2026-06-02
source: Stefano (run reale dal deliverable su Desktop)
fix: "internal/tui/tui.go: cleanPath ora chiama unescapeDragPath, che rimuove i backslash di escape shell (es. 'CAPABILITY\\ D\\ —\\ Profilo\\ x' → 'CAPABILITY D — Profilo x') inseriti dal drag&drop di una cartella nel campo TUI in Terminal. Coerenza: estende il fix tui-drag-drop-trailing-space (quote + spazi). Consistenza launcher: aggiunto anche auto-heal +x ai 7 lavori/*/Esegui.command."
test: "internal/tui/tui_test.go: TestCleanPath_DragDropFinderPatterns (2 nuovi casi backslash) + TestValidateExistingDir_HandlesBackslashEscapedSpaces (dir reale con spazi+em-dash, path escapato — riproduce il caso Stefano)."
---

# Bug: TUI — path da drag&drop con backslash-escape → "la cartella non esiste"

## Reported

Stefano lancia il `labnexus.command` root (TUI), sceglie `audit-checklist`, e nel
campo "Cartella di input" **trascina** la cartella del job. Terminal incolla il
path con gli spazi shell-escapati:

```
/Users/stefanofiorina/Desktop/.../lavori/CAPABILITY\ D\ —\ Profilo\ audit-checklist
```

Errore:
```
* la cartella "/Users/.../CAPABILITY\\ D\\ —\\ Profilo\\ audit-checklist" non esiste
```

(il `\\` è il quoting Go %q di backslash singoli.)

## Root cause

La TUI non è una shell: prende il path **letteralmente**, inclusi i backslash di
escape che il drag&drop in Terminal inserisce davanti agli spazi. `os.Stat` cerca
quindi una cartella con i backslash nel nome → ENOENT. `cleanPath` gestiva quote e
trailing space ma non il backslash-escape.

## Fix applicato

- `internal/tui/tui.go`: `unescapeDragPath` rimuove i backslash davanti ai
  metacaratteri shell; `cleanPath` lo invoca dopo trim+dequote.
- Consistenza: i 7 `lavori/*/Esegui.command` ora hanno l'auto-heal `+x` come il
  `labnexus.command` root (stesso root-cause cloud-sync del bug
  `launcher-permission-denied-cloud-sync`).

## Workaround immediato (dato a Stefano)

- **Raccomandato:** non usare la TUI per i job standard — doppio click su
  `lavori/<CAPABILITY> — Profilo <nome>/Esegui.command`: esegue `run --job`
  automaticamente, niente path da digitare/trascinare.
- In alternativa nella TUI: incolla il path senza backslash, oppure cancella i
  `\` dopo il drag.

## Follow-up

- Rebuild + ri-consegna zip col launcher e la TUI aggiornati.
