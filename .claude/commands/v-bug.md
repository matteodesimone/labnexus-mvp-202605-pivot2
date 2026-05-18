# /v-bug — Capture and document a bug

## Input
$ARGUMENTS (bug description, error message, steps to reproduce, or user report)

## Procedure

1. Read `.pipeline/state.json`. If a pipeline is active (`current_step` is not null):
   - WARN: "⚠ A pipeline is active (step: [current_step]). Bugs found during /v-review are review findings — fix them in the review cycle, not via /v-bug. Use /v-bug only for bugs found outside the current pipeline cycle."
   - ASK: "This is a bug found outside the current cycle? (y/N)"
   - If no: STOP. Tell human: "Fix the bug as part of the current review cycle. Run `/v-review` after fixing."

2. Read `.pipeline/standards/02-project.md` to check the project type.

3. Analyze the bug description. Extract:
   - **What happens** (the symptom)
   - **What should happen** (the expected behavior)
   - **How to reproduce** (steps, if available)
   - **Context** (which feature, which user flow, when it was found)

4. If reproduction steps are unclear, ASK the human for more detail before proceeding.

5. Assess severity:
   - **critical**: the app crashes, data is lost, or a security vulnerability is exposed
   - **high**: a primary feature doesn't work, but the app doesn't crash
   - **medium**: a secondary feature doesn't work, or a primary feature works but incorrectly in edge cases
   - **low**: cosmetic issue, typo, minor inconvenience

6. Ensure `.pipeline/bugs/` directory exists (create it if not — can happen on `--adopt` projects).

7. Create a bug file in `.pipeline/bugs/<slug>.md`. State lives in the **YAML frontmatter** at the top — never in body bullet lists. Body sections are for narrative (description, reproduction, analysis) only.

```markdown
---
severity: [critical | high | medium | low]
status: open
created: [today, ISO date]
source: [reporter — user, client, developer, test]
fix: ""
test: ""
---

# Bug: [short description]

## Reported
- Feature/Area: [which feature or area of the app]

## Description
[What happens vs what should happen]

## Steps to Reproduce
1. [step]
2. [step]
3. [observe: ...]

## Environment
- [Platform, browser, OS, version if relevant]

## Analysis
- Likely cause: [your initial analysis, or "needs investigation"]
- Affected files: [if you can identify them, or "unknown"]
- Risk of fix: [low: isolated change / medium: touches shared code / high: architectural impact]
```

`/v-bugfix` updates the `status`, `fix`, and `test` fields in the frontmatter when the bug is resolved. Do NOT duplicate these fields as a body `## Resolution` section — that representation is deprecated and causes drift.

8. Report:
```
Bug captured: .pipeline/bugs/<slug>.md
Severity: [severity]

Next steps:
  - `/v-bugfix bugs/<slug>.md` — reproduce and fix
  - `/v-triage` — prioritize alongside other bugs and todos
```

9. If project type is `prove-out`:
   - Note: "Bug fixes are budget-free by default. If this bug requires significant work beyond the original scope, evaluate at next checkpoint whether to allocate gettoni."
