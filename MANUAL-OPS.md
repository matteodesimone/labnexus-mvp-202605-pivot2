# Framework Manual Operations (senza vibbly CLI)

Il framework è file. Tutto ciò che vibbly fa si può fare a mano. Questa guida documenta come.

## Setup progetto (= vibbly init)

```bash
# 1. Crea la directory del progetto
mkdir -p ~/code/mio-progetto && cd ~/code/mio-progetto

# 2. Copia i file framework (da un template o da un progetto esistente)
#    Se hai il framework come zip:
unzip agent-pipeline.zip

#    Se hai un progetto esistente da usare come template:
cp -r ~/code/progetto-template/.pipeline .
cp -r ~/code/progetto-template/.claude .
cp -r ~/code/progetto-template/.hooks .
cp ~/code/progetto-template/CLAUDE.md .
cp ~/code/progetto-template/REVIEWER.md .

#    Preferred path: usa la CLI vibbly (installa framework + vobbly locale):
vibbly init --stack <stack> .

# 3. Git
git init
git add -A
git commit -m "feat: initial project with pipeline framework"
git branch development
git checkout development

# 4. Installa git hook
cp .hooks/pre-commit.sh .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

# 5. Onboard existing AI config (if any)
#    If the project already had CLAUDE.md, .cursorrules, AGENTS.md, etc.,
#    vibbly init archived them in .pipeline/onboarding/.
#    In Claude Code, run /v-onboard to migrate that context into 02-project.md.
#    If fresh project (no pre-existing AI files), skip this step.

# 6. Configura il progetto
#    Edita .pipeline/standards/02-project.md:
#    - Project Context section (MOST IMPORTANT — what/why/who/stage, architecture, decisions)
#    - Stack
#    - Quality metric
#    - Intentional deviations
#    - Out of scope
#
#    Edita .pipeline/config.sh — sezione ESSENTIAL:
#    - PROJECT_TYPE
#    - DEPLOY_CMD (es: "fly deploy")
#    - REVIEWER_*_ENABLED
#    - REVIEW

# 7. Verifica (= vibbly doctor)
#    Controlla manualmente:
command -v git && echo "✓ git" || echo "✗ git"
command -v mise && echo "✓ mise" || echo "✗ mise"
command -v gemini && echo "✓ gemini" || echo "✗ gemini"
test -f CLAUDE.md && echo "✓ CLAUDE.md" || echo "✗ CLAUDE.md"
test -f REVIEWER.md && echo "✓ REVIEWER.md" || echo "✗ REVIEWER.md"
test -f .pipeline/config.sh && echo "✓ config.sh" || echo "✗ config.sh"
test -f .pipeline/standards/00-principles.md && echo "✓ principles" || echo "✗ principles"
test -f .git/hooks/pre-commit && echo "✓ pre-commit hook" || echo "✗ pre-commit hook"
```

## Pipeline commands (= Claude Code)

Non cambia nulla. Apri Claude Code e usa i comandi `/`:

```bash
claude
# Pipeline principale
/v-brainstorm "idea vaga"   # Esplora e chiarisci prima di formalizzare
/v-grill                    # Stress-test del piano — spirito critico
/v-spec "descrizione feature"
/v-plan
/v-implement                # (= /v-execute in core)
/v-review
/v-deploy                   # (= /v-ship in core)
/v-compound

# Utility
/v-status                   # Stato pipeline
/v-guide [topic]            # Spiega qualsiasi concetto del framework
/v-idea "descrizione"       # Cattura idea futura
/v-report [type]            # Genera report (summary/full/client/ideas)
/v-onboard                  # Migra config AI esistente dopo init
/v-audit                    # Stress test multi-prospettiva (on-demand)
```

I comandi `.claude/commands/v-*.md` sono file markdown che Claude Code legge. Non dipendono da vibbly.

## Review multi-LLM (= vobbly hook review)

La review multi-reviewer è implementata in Go e invocata tramite vobbly. Claude Code la lancia automaticamente quando esegui `/v-review`; per lanciarla a mano:

```bash
.pipeline/bin/vobbly hook review all file1.py file2.py
# Output in .pipeline/reviews/:
#   <reviewer>-<check>.md  — un file per (reviewer, check)
#   _run.json              — summary strutturato (sempre prodotto)
#   summary.md             — marker degraded, se ≥50% dei reviewer falliscono
#
# L'ultima riga su stderr è "[REVIEW] COMPLETE: N/M reviewers ok, K quota-skipped, ..."
# oppure "[REVIEW] SETUP ERROR: ..." se la configurazione è invalida.
#
# Exit code (FR-3):
#   0 = all-green         (tutti i reviewer ok almeno su 1 check)
#   1 = partial           (almeno uno ok, almeno uno failed)
#   2 = degraded          (tutti i reviewer falliti)
#   3 = catastrophic      (nessun reviewer abilitato o binario missing)
```

## Gestione reviewer (= vibbly reviewer)

Edita `.pipeline/config.sh` a mano:

```bash
# Abilitare Gemini:
REVIEWER_GEMINI_ENABLED=true

# Abilitare Mistral (vibe):
REVIEWER_MISTRAL_ENABLED=true
REVIEWER_MISTRAL_CMD="vibe"

# Abilitare Ollama:
REVIEWER_OLLAMA_ENABLED=true
REVIEWER_OLLAMA_MODEL="codellama"

# Abilitare Codex:
REVIEWER_CODEX_ENABLED=true
```

## Gestione gate (= vibbly gate)

Edita `.pipeline/config.sh`:

```bash
# Cambiare un gate da human a tribunal:
GATE_PLAN="tribunal"

# Resettare tutti a human:
GATE_SPEC="human"
GATE_PLAN="human"
GATE_TEST_SCAFFOLD="human"
GATE_REVIEW_FAIL="human"
GATE_DEPLOY="human"
GATE_COMPOUND="human"
```

## Status (= vibbly status)

```bash
# Leggi lo stato:
cat .pipeline/state.json | python3 -m json.tool

# Leggi le fette:
cat .pipeline/fette.md 2>/dev/null || echo "Nessuna fetta definita"

# Leggi il budget (prove-out):
python3 -c "import json; s=json.load(open('.pipeline/state.json')); print(f'Budget: {s.get(\"budget\", \"n/a\")}')"
```

## Aggiornare il framework (= vibbly framework update)

```bash
# 1. Confronta la versione corrente con quella disponibile
cat .pipeline/version
cat ~/code/vibbly/framework/core/.pipeline/version  # o la source

# 2. Diff dei file che sono cambiati
diff -r .pipeline/standards/ ~/code/vibbly/framework/core/.pipeline/standards/
diff CLAUDE.md ~/code/vibbly/framework/core/CLAUDE-base.md

# 3. Copia manualmente i file aggiornati
cp ~/code/vibbly/framework/core/.pipeline/standards/00-principles.md .pipeline/standards/
# ... eccetera, file per file

# 4. Aggiorna version
cp ~/code/vibbly/framework/core/.pipeline/version .pipeline/version

# 5. Aggiorna project.lock
sed -i "s/framework_version:.*/framework_version: \"$(cat .pipeline/version)\"/" .pipeline/project.lock
```

## Promuovere un learning (= vibbly propose)

Dopo `/v-compound`, se hai un proposed-update:

```bash
# 1. Il file è in:
ls .pipeline/proposed-updates/

# 2. Copia nella inbox condivisa (manuale):
mkdir -p ~/.config/vibbly/inbox
cp .pipeline/proposed-updates/2026-03-22-xdg-pattern.md ~/.config/vibbly/inbox/

# 3. Nel repo framework, leggilo e applicalo:
cat ~/.config/vibbly/inbox/2026-03-22-xdg-pattern.md
# Applica manualmente le modifiche, poi:
rm ~/.config/vibbly/inbox/2026-03-22-xdg-pattern.md
```

## Git operations (= vibbly git)

```bash
# Status:
git branch -a
git status

# Sync (merge development in main):
git checkout main
git merge development
git push origin main
git checkout development

# Feature branch:
git checkout -b feature/mia-feature development

# Finish feature:
git checkout development
git merge feature/mia-feature
git branch -d feature/mia-feature
```

## Bug tracking (= /v-bug e /v-bugfix, senza pipeline)

```bash
# Cattura un bug (manuale):
mkdir -p .pipeline/bugs
cat > .pipeline/bugs/login-crash.md << 'EOF'
# Bug: Login crash on empty email
## Severity: high
## Description: ...
## Steps to Reproduce: ...
## Resolution
- Status: open
EOF

# Poi in Claude Code:
/v-bugfix bugs/login-crash.md
```

## Ideas backlog (= /v-idea)

```bash
# Capture an idea
mkdir -p .pipeline/ideas
cat > .pipeline/ideas/$(date +%Y-%m-%d)-mia-idea.md << 'EOF'
# Titolo dell'idea

- **Date**: 2026-03-23
- **Source**: standalone
- **Priority**: unset
- **Tags**: topic1, topic2

## Idea

Descrizione dell'idea...

## Notes

(vuoto)
EOF

# List ideas
ls -la .pipeline/ideas/*.md

# Review and set priority (edit the file, change Priority)
# Drop an idea (move to archive)
mkdir -p .pipeline/ideas/archive
mv .pipeline/ideas/idea-to-drop.md .pipeline/ideas/archive/

# Search
grep -rl "termine" .pipeline/ideas/
```

## Coherence test (= vobbly coherence)

```bash
# Esegui i BDD test del framework:
# Con godog (Go):
godog features/framework/

# Con pytest-bdd (Python):
pytest features/framework/

# O manualmente — i check sono grep:
# Tutti i principi sono nel review?
principles=$(grep -c '^## ' .pipeline/standards/00-principles.md)
review_checks=$(grep -A20 'Principles (Layer 0)' .claude/commands/v-review.md | grep -c '^\-')
[ "$principles" -eq "$review_checks" ] && echo "✓ Principles" || echo "✗ Principles: $principles defined, $review_checks checked"

# Tutti i comandi sono nel README?
commands=$(ls .claude/commands/ | wc -l)
readme_cmds=$(grep -c '| `/' .pipeline/README.md)
[ "$commands" -eq "$readme_cmds" ] && echo "✓ Commands" || echo "✗ Commands: $commands actual, $readme_cmds in README"
```

## Riepilogo: cosa fa vibbly vs cosa fai a mano

| vibbly CLI | Manuale |
|-----------|---------|
| `vibbly init` | Copia file framework + archive `.pipeline/onboarding/` + git init |
| `vibbly doctor` | Verifica comandi uno per uno |
| `vibbly status` | `cat .pipeline/state.json` |
| `vibbly config` | Edita `config.sh` e `02-project.md` |
| `vibbly reviewer gemini on` | `REVIEWER_GEMINI_ENABLED=true` in config.sh |
| `vibbly gate plan set tribunal` | `GATE_PLAN="tribunal"` in config.sh |
| `vibbly framework update` | Diff + copia file + aggiorna version |
| `vibbly propose` | Sposta file da `.pipeline/proposed-updates/` a `~/.config/vibbly/inbox/` (source clear post-submit) |
| `.pipeline/bin/vobbly hook review` | Invocazione canonica della review multi-reviewer |
| `vibbly git sync` | `git checkout main && git merge development && git push` |

I pipeline commands (`/v-spec`, `/v-plan`, etc.) funzionano identicamente — sono file markdown che Claude Code legge. Non dipendono da vibbly.
