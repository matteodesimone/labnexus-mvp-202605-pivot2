---
severity: critical
status: fixed
created: 2026-05-19
fixed: 2026-05-19
source: developer (CTO smoke test post-deploy bugfix, lanciato il binario reale)
fix: "Nuovo package `internal/paths/` con due funzioni: `DeliveryRootForBinary(path) string` (pura, pattern-matching su `<X>.app/Contents/MacOS/<bin>` → ritorna parent del `.app`; altrimenti `filepath.Dir(path)`) e `DeliveryRoot() string` (usa `os.Executable()` + `filepath.EvalSymlinks` poi delega alla funzione pura). `cmd/labnexus/main.go` ora calcola `defaultProfiliDir = filepath.Join(paths.DeliveryRoot(), \"profili\")` e analogo per kb-dir, invece di stringhe hardcoded `./profili` (che erano relative al cwd). Makefile `smoke` target aggiornato per passare `--profiles-dir profili --kb-dir KB-ispettore` espliciti (perché `bin/labnexus` ora cerca `bin/profili` per default)."
test: "internal/paths/paths_test.go (8 unit test table-driven con 6 pattern di distribuzione: bundle macOS, /Applications, standalone bin/, linux/amd64, .app.bak intermedio, MacOS dir non-bundle + 3 test su DeliveryRoot reale e edge case relative-path / no-IO). internal/bundle/bundle_test.go::TestBundle_BinaryFindsProfilesAdjacentToApp (integration: estrae bundle in tmpdir + crea profili adiacenti + lancia binario con cwd ARBITRARIA + verifica `list` mostra il profilo)."
---

# Bug: `./profili` e `./KB-ispettore` risolti vs cwd dell'invocazione, non vs binario

## Reported

- **Feature/Area**: FR-2 (path resolution per `--profiles-dir` e `--kb-dir`) — deliverable Fetta 1
- **Trovato da**: CTO mentre testava il binario reale `labnexus-bin` post-bugfix-deploy
- **Sprint impact**: **blocker totale della consegna**. Denis che estrae lo zip ed esegue da Finder (doppio click) o da Terminal otterrebbe sempre questo errore — il binario è inusabile fuori dal repo del CTO. Più grave del bug doppio-click (quello almeno permetteva il flusso CLI dalla repo root).

## Description

**Cosa succede**:
- L'utente estrae `dist/labnexus-sprint1-darwin-arm64.zip` in una cartella (es. `~/labnexus/`)
- Lancia il binario reale da Terminal: `'~/labnexus/labnexus.app/Contents/MacOS/labnexus-bin'`
- Il cwd del processo è la directory corrente del Terminal (tipicamente `~`)
- Il binario cerca `./profili` (relativo al cwd = `~`) → **`open ./profili: no such file or directory`**
- L'app esce con errore

Traccia esatta (CTO):
```
'/Users/.../dist/labnexus-sprint1-darwin-arm64/labnexus.app/Contents/MacOS/labnexus-bin'
Error: tui: lettura profili da "./profili": profile: lettura cartella "./profili": open ./profili: no such file or directory
```

**Cosa dovrebbe succedere** (vedi FR-2 nella spec):
> L'eseguibile risolve i path relativi alla propria posizione: profili in `./profili/<nome>.yml`, KB-ispettore in `./KB-ispettore/`.

La spec dice *relativo alla propria posizione* (= directory del binario / del bundle). L'implementazione attuale risolve *relativo al cwd*, che è qualcosa di completamente diverso. **La spec è chiara, il codice è sbagliato.**

## Steps to Reproduce

Lo zip contiene:
```
labnexus-sprint1-darwin-arm64/
  labnexus.app/Contents/MacOS/labnexus-bin   ← binario reale
  labnexus.app/Contents/MacOS/labnexus       ← wrapper bash
  profili/{revisione,rilievi}.yml            ← profili shippable
  KB-ispettore/...                            ← KB
  README.md
```

Riproduzione:
1. Estrai lo zip in `~/labnexus-test/` (qualunque path)
2. Apri Terminal, cwd = `~`
3. Lancia: `~/labnexus-test/labnexus.app/Contents/MacOS/labnexus-bin`
4. **Osserva**: errore `open ./profili: no such file or directory`

Stesso errore anche con:
- `~/labnexus-test/labnexus.app/Contents/MacOS/labnexus-bin list`
- Doppio click su `~/labnexus-test/labnexus.app` dal Finder (cwd di Launch Services = `/`, il wrapper apre Terminal con shell HOME = `~`, stesso problema)
- `~/labnexus-test/labnexus.app/Contents/MacOS/labnexus-bin --profiles-dir ./profili` (flag esplicito ma sempre relativo a cwd)

L'unico caso in cui funziona è lanciare con cwd = directory che CONTIENE `profili/` e `KB-ispettore/` (= la directory `labnexus-sprint1-darwin-arm64/` estratta), es. `cd ~/labnexus-test/labnexus-sprint1-darwin-arm64/ && ./labnexus.app/Contents/MacOS/labnexus-bin`.

## Environment

- macOS Apple Silicon (`darwin/arm64`)
- Zip: `dist/labnexus-sprint1-darwin-arm64.zip` (Fetta 1 + bugfix-1 ship 2026-05-19T11:02)
- Default flag in `cmd/labnexus/main.go:buildRoot()`:
  ```go
  root.PersistentFlags().String("profiles-dir", "./profili", ...)
  root.PersistentFlags().String("kb-dir", "./KB-ispettore", ...)
  ```

## Analysis

**Likely cause (alta confidenza)**:

In `cmd/labnexus/main.go` i default sono path **relativi** stringhe `"./profili"` e `"./KB-ispettore"`. Go interpreta `.` come cwd del processo (CWD ereditato dal parent shell), non come directory del binario.

La spec FR-2 dice "relativo alla propria posizione" → richiede `os.Executable()` per ottenere il path del binario e calcolare `<dir-del-binario>/../../../profili` (3 livelli su perché il binario è in `<bundle>/Contents/MacOS/`).

Più precisamente: se il binario è `/PATH/labnexus.app/Contents/MacOS/labnexus-bin`, la directory "consegnata" (= dir adiacente a `.app`) è `/PATH/`. I `profili/` e `KB-ispettore/` sono attesi lì. Quindi:
- bundle binary path: `/PATH/labnexus.app/Contents/MacOS/labnexus-bin`
- delivery root: `/PATH/`
- profili: `/PATH/profili/`
- KB: `/PATH/KB-ispettore/`

Per il binario **fuori dal bundle** (es. `bin/labnexus` da `make build`, o `dist/labnexus-linux-amd64`), la "delivery root" è la directory del binario stesso.

**Logic di risoluzione (proposta)**:

```go
func defaultDeliveryRoot() string {
    exe, err := os.Executable()
    if err != nil {
        return "."  // fallback al cwd, ma è il caso patologico
    }
    exe, _ = filepath.EvalSymlinks(exe)  // risolvi symlink (es. /usr/local/bin → cellar)
    dir := filepath.Dir(exe)
    // Se siamo dentro un .app bundle, sali a parent del .app
    // Pattern: .../labnexus.app/Contents/MacOS/labnexus-bin
    if strings.HasSuffix(dir, "/Contents/MacOS") {
        return filepath.Dir(filepath.Dir(filepath.Dir(dir)))  // 3 livelli su
    }
    return dir
}
```

Poi:
```go
deliveryRoot := defaultDeliveryRoot()
root.PersistentFlags().String("profiles-dir", filepath.Join(deliveryRoot, "profili"), ...)
root.PersistentFlags().String("kb-dir", filepath.Join(deliveryRoot, "KB-ispettore"), ...)
```

**Affected files**:
- `cmd/labnexus/main.go` (default flags `--profiles-dir`, `--kb-dir` — calcolo a runtime)
- `internal/runner/runner.go` (fallback in `validateConfig`: `"./profili"` e `"./KB-ispettore"` — anche qui sbagliato)
- `internal/tui/tui.go` (eredita i default via flag — OK se main risolve correttamente)

**Risk of fix**: **basso**. Calcolo isolato in una funzione helper, riusato in 1-2 punti. Niente architectural impact.

**Test plan**:
- Test unitario: `defaultDeliveryRoot()` su 3 input:
  1. `/path/labnexus.app/Contents/MacOS/labnexus-bin` → `/path`
  2. `/path/bin/labnexus` → `/path/bin`
  3. `/path/labnexus-linux-amd64` → `/path`
- Test integration (bundle_test.go): lancia il bundle estratto in tmp dir, verifica che `labnexus list` legga i profili senza --profiles-dir esplicito
- BDD scenario nuovo (Fetta 1 fixup): "lancio binario dal Finder/CLI senza flag esplicito risolve profili adiacenti al .app"

## Note di scope (prove-out)

**Bug fix budget-free** (vedi `prove-out.md`). Stima fix: ~30-60 min (helper + 2 punti di sostituzione + test unitario + test integration + rebuild zip).

**Learning aggiuntivo** (oltre quello già catturato in `compound-bugfix-app-bundle-doppio-click.md`):

Lo scenario `@manual` "bundle .app contiene wrapper bash + binario reale" verifica la **struttura** del bundle, NON il **comportamento operativo** dell'estratto. Mancava un test "lancia il binario reale da una cartella estratta arbitraria e verifica che funzioni" — è quello che avrebbe catturato sia il bug doppio-click sia questo bug path-resolution.

Idea: **estendere lo scenario `@manual` con un sub-step automatizzabile** che fa `unzip dist/...zip -d tmp/; tmp/.../labnexus-bin list` da una cwd diversa da `tmp/`. Verifica che il comando giri OK senza flag. Questo è automatizzabile e cattura entrambe le categorie di bug visti finora.

→ Aggiungere a `.pipeline/ideas/2026-05-19-httptest-bundle-doppio-click.md` come scope estesso? Sì, lo annoto.
