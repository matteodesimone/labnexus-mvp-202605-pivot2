# Skill: Per-Folder Launcher Zero-CLI

**Created**: 2026-05-22 in sub-fetta 1.5.C (compound, post-deploy iteration)
**Trigger**: CLI tool con una nozione di "task" / "job" / "lavoro" mappato a una cartella, utenti non-technical (no Terminal), zero-interaction UX desiderata.

## Il problema

LabNexus Sprint 1.5.A+B aveva un singolo `labnexus.command` al root del deliverable. Doppio click → TUI sequenziale che chiede:
1. Selezione profilo (frecce + invio)
2. Cartella input (drag&drop da Finder)
3. Cartella output (drag&drop)
4. Run

3 step di interazione per ogni esecuzione. Per utenti tecnici è accettabile. Per Denis (QM, non sviluppatore) era ancora troppo. Il CTO ha segnalato l'attesa di "concierge vero" = 1 doppio click su un file pre-configurato dentro la cartella di lavoro.

## Il pattern

**Un launcher script identico in ogni cartella di lavoro**, che deriva il nome del job dal proprio `basename`. Doppio click → run immediato, zero interazioni.

### Struttura

```
<deliverable-root>/
├── labnexus                    (binary)
├── labnexus.command            (launcher generico/TUI, optional fallback)
├── lavori/
│   ├── CAPABILITY A — Profilo revisione/
│   │   ├── Esegui.command      ← per-folder launcher
│   │   ├── _labnexus.toml      ← metadata (profile, opzionali)
│   │   └── <file di input...>
│   ├── CAPABILITY B — Profilo rilievi/
│   │   ├── Esegui.command      ← STESSO FILE, no template vars
│   │   ├── _labnexus.toml
│   │   └── <file di input...>
│   └── ... (altre cartelle)
```

### Lo script (esempio bash macOS)

```bash
#!/usr/bin/env bash
set -e
DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(cd "$DIR/../.." && pwd)"   # risali alla deliverable root
JOB_NAME="$(basename "$DIR")"       # nome del lavoro = nome cartella

cd "$ROOT" || exit 2

# (eventuali pre-check: API key, deps, ecc.)
# ...

# Esegui il job. Il binary risolve il job per substring match sul basename.
./<binary> run --job "$JOB_NAME" || {
    echo ""
    echo "ERRORE: esecuzione fallita."
    read -n 1 -s -r -p "Premi un tasto per chiudere..."
    exit 1
}

echo ""
echo "Completato. Output: $DIR/output/"
read -n 1 -s -r -p "Premi un tasto per chiudere..."
```

### Convenzioni critiche

1. **Stesso file identico per tutte le cartelle**: niente template variables, niente generation script. Il binary RISOLVE il job dal basename. Vantaggio: Denis duplica una cartella per crearne una nuova → l'`.command` funziona automaticamente.
2. **Basename autoderive**: lo script usa `$(basename "$DIR")` come job identifier. Il binary deve avere un meccanismo di resolution (substring match, case-insensitive ideal).
3. **Risale a deliverable root via `../..`**: lo script assume la struttura `<root>/<container>/<job-folder>/launcher.command`. Per due livelli (cartella container "lavori/" + cartella job). Adattare se la profondità è diversa.
4. **Sempre `read` finale**: anche su successo, lascia la finestra Terminal aperta finché l'utente non preme un tasto. Altrimenti il Finder chiude immediato e l'utente non vede l'output.
5. **Path con spazi**: il nome cartella ospitante può avere caratteri Unicode / spazi (es. "CAPABILITY A — Profilo revisione" con em-dash). Lo script deve sempre quotare le variabili (`"$DIR"`, `"$JOB_NAME"`).
6. **Permessi**: `chmod +x` su ogni launcher prima dello zip (`scripts/build-zip.sh` deve preservare il bit X).

### Skip filter parser input

Il launcher script + metadata file (`_labnexus.toml`, `.DS_Store`) sono DENTRO la cartella di input. Il parser dell'input deve skipparli silenziosamente per non emettere "formato non supportato" warning ad ogni esecuzione. Pattern:

```go
var internalFiles = map[string]bool{
    "_labnexus.toml":  true,
    "Esegui.command":  true,
    ".DS_Store":       true,
}

for _, e := range entries {
    if internalFiles[e.Name()] {
        continue // skip silenzioso
    }
    // ... parse
}
```

## Trade-off

- **Pro**: 0 interazioni utente per esecuzione. UX concierge vero. Denis aggiunge un lavoro = duplica cartella + cambia _labnexus.toml. Niente Terminal, niente CLI flags da ricordare.
- **Pro**: scalabile a N lavori senza modifiche al deliverable (no menu hardcoded, no config da rebuildare).
- **Con**: macOS Gatekeeper richiede autorizzazione su uno dei `.command` la prima volta (gli altri identici si autorizzano automatic). Documenta.
- **Con**: `.command` non funziona su Windows. Per cross-platform serve `.bat` (Win) + `.command`/`.sh` (Unix) generati in build (out-of-scope Sprint 1.5).
- **Con**: il binary deve avere resolution job-by-name robusto (es. substring case-insensitive). Se manca, il pattern non funziona.

## Anti-pattern

- **Template script con variabili sostituite**: genera N file diversi al build time, ognuno con il job hardcoded. Funziona ma Denis non può duplicare la cartella senza che la copia funzioni male (hardcoded job points to wrong folder).
- **Single launcher root + menu jobs**: ridiventa la TUI sequenziale Sprint 1.5.A — utente deve scegliere. Sopra il setup minimo, l'attrito è uguale.
- **`labnexus run --job X` da Terminal**: richiede CLI, Denis non lo userà. OK come modalità avanzata, non come primario.

## Precedente concreto

Sub-fetta 1.5.C LabNexus (post-deploy iteration `097f5ac`): 7 file `Esegui.command` identici uno in ogni `lavori/CAPABILITY <X> — Profilo <Y>/`. Il binary `./labnexus` risolve il job via `jobs.Resolve` (case-insensitive substring match sul folder basename). Pre-check API key shell-side coerente con `labnexus.command` root (FR-36).

Vedi `.pipeline/solutions/2026-05-22-sub-fetta-1.5.C-concierge-mode-chiusura-sprint-1.5.md` per il caso completo + `scripts/lavoro-template.command` (sorgente canonica copia-incolla).

## Quando NON applicare

- Tool dove l'input è arbitrario / fuori da una cartella fissa (es. file uploader interattivo)
- Tool con 1-2 task totali: 1 launcher root con flag fissi è più semplice
- Utenti tecnici / sviluppatori: `cli run --job X` da Terminal è più rapido di navigare Finder

## Generalizzabilità (verso framework promotion)

Pattern universal per ogni CLI tool con auto-discovery di task per cartella. Candidato a `.pipeline/proposed-updates/` come `skills/framework/`. Particolarmente utile per progetti vibbly con dominio cli-tool quando lo user persona target è non-developer.
