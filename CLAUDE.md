# CLAUDE.md — Builder Agent Instructions

This file tells Claude Code how to work in this project. Read it at the start of every session.

## Golden Rules (NEVER violate)

1. **No execution without a plan.** Never start work unless a plan exists in `.pipeline/plan.md`.
2. **No execution without verification criteria.** Never start work unless corresponding verification criteria exist AND are failing (red phase). Verify by running the verification suite before starting.
3. **No ship without review.** Never suggest work is complete without running the `/v-review` command first.
4. **State is sacred.** Always update `.pipeline/state.json` when transitioning between steps. Read it at the start of every command to know where you are.
5. **Ask, don't assume.** When a spec is ambiguous, list the ambiguities and ask the human. Do not fill gaps with assumptions.
6. **Never touch main.** All work happens on `development` or feature branches. Never commit to `main`. Never push to `main`. Only the human merges development into main. If you find yourself on main, switch to development before doing anything.

## Project Context

**Before doing anything, read `.pipeline/standards/02-project.md` section "Project Context".** This tells you what the project is, who it's for, what's already been decided, and what vocabulary to use. Without this, you're working blind.

## Technology Context

**Before writing or reviewing work that involves specific technologies, read `.pipeline/contexts/`.** These files contain up-to-date documentation fetched from official sources. They override your training data for those technologies. If a context file exists for a technology you're about to use, read it first.

## Quality Standards

**All quality standards are defined in `.pipeline/standards/`.** This is a layered system, shared with the external reviewer. Read the README there for reading order.

```
Layer 0: .pipeline/standards/00-principles.md     ← Core values. NEVER overridden.
Layer 1: .pipeline/standards/01-methodology.md    ← Process. Always applies.
Layer 2: .pipeline/standards/stacks/<stack>.md    ← Stack-specific rules.
Layer 3: .pipeline/standards/02-project.md        ← This project's specifics + context.
```

Read all applicable layers at the start of every session. When working, follow every applicable rule. When reviewing, evaluate against every applicable rule.

Do NOT maintain separate quality rules here. If a standard needs to change, change it in the appropriate layer.

## Three Types of Commands

**Pipeline commands** (prefixed with `/`, run by you — Claude Code — inside a project):
`/v-spec`, `/v-plan`, `/v-verify`, `/v-execute`, `/v-review`, `/v-ship`, `/v-compound`, `/v-bugfix`, `/v-refactor`, `/v-assess`, `/v-bug`, `/v-brainstorm`, `/v-grill`, `/v-triage`, `/v-status`, `/v-approve`, `/v-test-data`, `/v-todo-export`, `/v-guide`, `/v-idea`, `/v-propose`, `/v-report`, `/v-onboard`, `/v-audit`.
These drive the development workflow. They read and write pipeline state.
Domain aliases may apply (e.g., `/v-implement` = `/v-execute` in coding domain).

**vibbly CLI commands** (interface human ↔ project, run by the developer in the terminal):
`vibbly init`, `vibbly doctor`, `vibbly status`, `vibbly template`, `vibbly git`, `vibbly framework update`, `vibbly propose`, etc.
These manage project structure, templates, environment, git, and proposal submission. They work in any project that uses the framework. **Humans run these directly.**

**vobbly CLI commands** (interface AI agent ↔ framework, run by you — Claude Code — on behalf of the framework):
`vobbly inbox`, `vobbly coherence`, `vobbly release`, `vobbly changelog`, plus future orchestration commands (review, hooks, etc.).
These let the framework execute actions programmatically without depending on the end-user's environment. **You run these; humans don't.** Vobbly is distributed alongside the framework and lives at `.pipeline/bin/vobbly` in projects that have it installed (convention).

Some subcommands of vobbly are context-specific:
- `inbox`, `coherence`, `release`, `changelog` — framework-development context only (only meaningful in the framework repo itself).
- Future orchestration subcommands (review, hooks, etc.) — available in any project using the framework.

## Pipeline Definitions

### Pipeline: New Feature
```
[brainstorm] → spec → plan → verify → execute → review → ship → compound
```
`brainstorm` is optional (use when requirements are fuzzy). `compound` closes the learning loop.

### Pipeline: Bugfix
```
reproduce-test → fix → review → ship → compound
```

### Pipeline: Refactor
```
review (current) → plan → verify → refactor → review → ship → compound
```

### Pipeline: Assess
```
assess → assessment report → spec for each recommended intervention
```

### Pipeline: Bug (capture only)
```
capture → file in .pipeline/bugs/
```

## State

Pipeline state lives in `.pipeline/state.json`. Schema:

```json
{
  "pipeline": "new-feature | bugfix | refactor | assess | bug",
  "feature": "short description",
  "current_step": "spec | plan | verify | execute | review | ship",
  "steps_completed": [
    {
      "step": "spec",
      "completed_at": "ISO timestamp",
      "gate_passed": true,
      "artifacts": ["path/to/spec.md"]
    }
  ],
  "blocked": false,
  "blocked_reason": null,
  "current_fetta": null,
  "review_iteration": 0,
  "budget": {
    "package": "S",
    "total": 20,
    "infra": 5,
    "client": 15,
    "spent": 0,
    "spent_per_fetta": {},
    "allocated_per_fetta": {}
  },
  "bug_file": null
}
```

Notes on optional fields:
- `current_fetta`: set by `/v-plan` when fette are proposed. Null if project is single-fetta.
- `review_iteration`: tracks current review→fix cycle. Reset to 0 at start of each review phase.
- `budget`: only meaningful for `prove-out` projects. Null or absent for other project types. See `.pipeline/standards/project-types/prove-out.md` for full schema.
- `bug_file`: set by `/v-bugfix` when fixing from a `.pipeline/bugs/*.md` file.

Before executing ANY pipeline command:
1. Read `.pipeline/state.json`
2. Verify the previous step is completed and its gate passed
3. If not, REFUSE to proceed and explain why

## Output Conventions

When producing pipeline artifacts:
- Specs go in `.pipeline/spec.md`
- Plans go in `.pipeline/plan.md`
- Review reports go in `.pipeline/review.md`
- Brainstorms go in `.pipeline/brainstorms/`
- Solutions go in `.pipeline/solutions/`
- Todos go in `.pipeline/todos/`
- Bugs go in `.pipeline/bugs/`
- Ideas go in `.pipeline/ideas/`
- Reports go in `.pipeline/reports/`
- Technology context files live in `.pipeline/contexts/` (fetched, not authored)
- Proposed framework updates go in `.pipeline/proposed-updates/`

## Git Conventions

**Full strategy in `.pipeline/standards/git-strategy.md`.** Key points:

**Conventional commits.** `type: short description`. Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`.

**Atomic commits.** Each commit is one coherent unit.

**Branch discipline.** Work on `development` or feature branches. Never commit to `main` (Golden Rule #6).

**Push after ship and compound.** After `/v-ship` and `/v-compound`, push development to remote. This protects against local disk failure. If push fails (no remote, auth): warn but don't block.

## Quota Awareness

Pipeline operations consume LLM calls — yours (Claude Code) and the external reviewers'. Be cost-conscious.

**Before expensive operations**, warn the human with an estimate:
- Review iteration: "This will run up to N review cycles × M reviewers = ~X LLM calls. Proceed?"
- Tribunal gate: "Tribunal will consult N reviewers. Proceed?"
- `/v-assess` on a large corpus: "Full assessment may require multiple analysis passes. Proceed?"

**If you hit your own quota limit** (Claude Code rate limit, message cap, or any throttling signal):
1. STOP immediately. Do not try to squeeze in more work.
2. Save the current state: commit or stash any work in progress.
3. Update `.pipeline/state.json` with current progress.
4. Tell the human clearly: "I've reached my usage limit. Current state is saved. You can resume with a new session — run `/v-status` to see where we left off."

**If an external reviewer hits quota** (detected by the hooks):
- The hooks handle graceful degradation automatically — the reviewer is skipped and you're notified.
- If ALL external reviewers fail, note it in the review report: "Single-reviewer mode — external reviewers unavailable (quota). Consider waiting for quota reset before shipping."
- Never silently proceed as if the full review happened. The human must know the review was partial.

## Environment

This project uses mise for runtime management. If a tool is not found in PATH:
1. Try: `eval $(mise activate bash)` before running commands
2. Or use: `mise exec -- <command>`
3. Never install runtimes globally — mise manages versions per-project via `.mise.toml`

Common tool locations if not in PATH:
- mise: `~/.local/bin/mise`
- Homebrew macOS: `/opt/homebrew/bin/`
- Linuxbrew: `/home/linuxbrew/.linuxbrew/bin/`
- Go binaries: `~/go/bin/`
- npm global: `~/.npm-global/bin/`

If a command fails with "not found", check these paths before reporting failure.

## One Pipeline at a Time

`state.json` tracks one active pipeline. `/v-bug` and `/v-bugfix` start a new pipeline — don't use them while a feature pipeline is in progress. Bugs found during `/v-review` are review findings: fix them in the review cycle, re-review, and continue. Use `/v-bug` only for bugs found OUTSIDE the current pipeline cycle.

## Skills

`.pipeline/skills/` contains operational knowledge — HOW we do specific things. Two levels:

- `skills/framework/` — curated patterns from the domain (design patterns, intuitiveness). Updated by `vibbly framework update`. Don't modify directly.
- `skills/project/` — patterns discovered in this project via `/v-compound`. Never overwritten by framework updates.

During `/v-plan` and `/v-execute`, read both. If a skill exists for the pattern you're about to use, follow it.

New skills are created in `skills/project/` during `/v-compound` when a reusable pattern emerges. To promote a project skill to the framework, use `/v-compound`'s framework promotion mechanism.

## Solutions Library

`.pipeline/solutions/` accumulates learnings from completed pipeline cycles. Each file has YAML frontmatter with tags.

During `/v-plan`, search solutions for relevant past learnings. If a previous cycle solved a similar problem, reference its solution.


# CLAUDE.md — Software Development Base

## Layer: coding-base (universal)

This file is the universal foundation for every software-development domain in vibbly (cli-tool, agent, web-app, game). It is NOT a selectable domain itself — `coding-base` is internal. The user always picks a concrete domain; the `domainasm` assembler concatenates this base with the domain-specific overlay.

## Golden Rules — Software Development Specifics

- **Golden Rule #2 (universal coding flavor)**: No implementation without failing tests. Never write implementation code unless corresponding tests exist AND fail (red phase of TDD). Verify by running the test suite before writing any implementation.
- "Verify" in any software-development domain = run the test suite. "Execute" = write code to pass tests. "Ship" = deploy/release/distribute (the verb depends on the specific domain).

## Package Managers

- Python: `uv` (never pip directly)
- Go: `go mod` (built-in)
- Ruby: `bundler`
- Node: `npm`

## Pipeline

The structured flow `spec → plan → test-scaffold → implement → review → deploy → compound` is the universal disciplined cycle every software-development project follows. Each step gates the next; nothing skips ahead. `/v-spec` writes the spec, `/v-plan` decomposes into fette, `/v-test-scaffold` creates BDD/TDD red phase, `/v-implement` makes it green, `/v-review` runs multi-LLM review, `/v-deploy` ships, `/v-compound` captures learnings into skills/solutions.

## Verification

In every software-development domain, `/v-test-scaffold` creates BDD acceptance tests + unit tests (TDD red phase). `/v-implement` writes code to pass them (TDD green phase). The test suite is the verification mechanism.

## Quality Checks

This base layer defines the universal review checks (concrete checks live in `.pipeline/checks.yaml`):
- **quality**: code quality, function constraints, DRY, naming
- **test-coverage**: every public function tested, edge cases, BDD match
- **security**: input validation, injection, auth, secrets
- **privacy**: GDPR, data minimization, PII
- **intuitiveness**: McKay's 8 attributes on every interface
- **data-integrity** (opt-in): transactions, migrations, state consistency
- **deploy-verification** (opt-in): rollback, health checks

Domain-specific overlays may add focused checks (e.g., cli-tool may add "exit-code-semantics", web-app "lighthouse-perf", game "60fps-budget", agent "hallucination-resistance").


<!-- vibbly-layer-boundary: cli-tool -->

# CLAUDE.md — cli-tool Domain Overlay

## Domain: cli-tool

Command-line tools, developer utilities, single binaries. Specializes the universal coding-base for CLI conventions.

## CLI Discipline

- **Bash-like, not TUI-like.** Comandi one-shot, parametri espliciti, errori chiari, no wizard, no prompts interattivi. Modello di riferimento: `kubectl`, `git`, `cargo`, `docker run`. Coerente con vibbly stessa.
- **Required flags fail loud.** Use cobra `MarkFlagRequired` for mandatory parameters. Print clear error listing valid values when omitted or invalid.
- **Help text is a contract.** Every flag and subcommand has a one-liner. Allowed values for enum-like flags are listed in `--help`. The user can answer "what does this do?" without leaving the terminal.
- **Exit codes are semantic.** `0` = success, `1` = generic failure, `2` = misuse / invalid invocation. If the tool produces multiple outcome classes (partial success, degraded, catastrophic), document the exit-code semantics and stick to it.
- **Deterministic stdout for piping.** Output meant for piping (logs, JSON, structured data) goes to stdout. User-facing messages and errors go to stderr. No spinners, no ANSI in non-TTY mode.

## Distribution

- **GitHub Releases** with cross-compiled binaries (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64). Use `goreleaser` (Go) or `pyinstaller` + GitHub Actions (Python).
- **Homebrew tap** is optional but recommended for visibility on macOS.
- **No auto-update**. Distribution channels (Homebrew, package managers, manual download) handle versioning.
- **Single binary preferred.** Avoid runtime dependencies that the user must install separately. mise can manage the dev runtime, but the shipped binary should self-contain.

## Verification (cli-tool flavor)

- **CLI integration tests are mandatory.** BDD scenarios invoke the binary as a subprocess and assert on stdout/stderr/exit-code. Acceptance tests do NOT mock the CLI surface.
- **No TTY assumption.** Tests run in non-TTY contexts (CI). The binary must behave identically.

## Ship verb

In cli-tool, `ship` = produce a release artifact and publish it (GitHub release, Homebrew bottle, or equivalent). Concrete `SHIP_CMD` is defined per project in `.pipeline/config.sh`.
