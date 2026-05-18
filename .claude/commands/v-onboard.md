# /v-onboard — Extract pre-existing AI context into the vibbly framework

## When to Use

Run ONCE after `vibbly init` on a project that had existing AI configuration files (CLAUDE.md, .cursorrules, copilot-instructions, etc.). `vibbly init` archives those files in `.pipeline/onboarding/` — this command reads them and populates `.pipeline/standards/02-project.md` with the extracted project-specific context.

Skip if `.pipeline/onboarding/` doesn't exist (greenfield project, nothing to onboard).

## Procedure

1. Check if `.pipeline/onboarding/` exists. If not, tell the human: "No pre-existing AI configuration found — nothing to onboard. Fill `.pipeline/standards/02-project.md` manually."

2. Read EVERY file in `.pipeline/onboarding/` (except `ADOPT-REPORT.md`). For each file, extract:

   **Project context** — what the project is, who it's for, what it does, its current stage:
   - Project description, purpose, target users
   - Architecture overview (monolith? microservices? stack components?)
   - Stage of development (greenfield, MVP, production?)

   **Settled decisions** — things the original config told the AI not to change:
   - Technology choices presented as final
   - Architectural patterns explicitly chosen
   - Things marked as "do not modify" or "intentional"

   **Conventions** — naming, formatting, organizational rules:
   - Code style preferences
   - File organization rules
   - API conventions, naming patterns
   - Domain-specific vocabulary

   **Constraints and deviations** — things that look wrong but are deliberate:
   - Missing features that are intentionally absent
   - Unusual patterns with justification

   **Out of scope** — things the project explicitly doesn't do

3. Read the current `.pipeline/standards/02-project.md`. It will be mostly empty (template with comments).

4. Populate `02-project.md` with the extracted information:
   - **Project Context** section: what/why/who/stage, architecture overview, key decisions, vocabulary
   - **Project Identity**: name, client, quality metric, scale
   - **Intentional Deviations**: from the original config's "don't change this" rules
   - **Out of Scope**: from the original config's exclusions
   - **Project-Specific Conventions**: from the original config's style/naming rules

   If the domain is coding, also populate:
   - **Stack Selection**: detected from the codebase or original config
   - **Database Choice**: detected from the codebase or original config
   - **Deploy**: detected from CI/CD files or original config
   - **Privacy & Data Handling**: detected from the original config

5. Check if the original CLAUDE.md (or equivalent) had rules that CONFLICT with vibbly's framework rules. If conflicts exist, list them and ask the human which to keep:
   - "The original config said [X], but vibbly's framework says [Y]. Which should this project follow?"
   - Framework rules go in `.pipeline/standards/`, project overrides go in `02-project.md` Intentional Deviations.

6. If the project has a `README.md`, read it too — it often contains architecture and setup context that belongs in 02-project.md.

7. Report what was onboarded:
```
Onboarding complete.

Populated in 02-project.md:
  ✓ Project Context (from: CLAUDE.md, README.md)
  ✓ Key Decisions (from: CLAUDE.md)
  ✓ Conventions (from: .cursorrules)
  ✓ Intentional Deviations (from: CLAUDE.md)

Conflicts found: N (resolved: N, pending: N)

Originals preserved in: .pipeline/onboarding/

Review 02-project.md and adjust as needed.
```

## Rules

- **Never delete the archived files.** They're the safety net.
- **Don't blindly copy.** Analyze and restructure. The original CLAUDE.md might be a mess — your job is to extract the useful bits and place them in the right vibbly sections.
- **Prefer the framework's rules** when there's a conflict, unless the human explicitly says otherwise. The framework is the methodology; the project file is the context.
- **Don't invent.** If the original config doesn't mention something, leave that section empty in 02-project.md. Don't guess the architecture from file names.
- **Run only once.** After onboarding, the archived files are reference only. Changes go in 02-project.md.
