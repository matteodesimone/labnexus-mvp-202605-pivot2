# /v-compound — Capture learnings after a completed pipeline cycle

## When to Run

After `/v-ship` completes (or after a pipeline cycle is done). This is the step that closes the "decidi, fai, osserva, impara" loop. Skip it and you've done traditional engineering. Do it and each cycle makes the next one better.

## Procedure

1. Read `.pipeline/state.json` to understand what was just completed.
2. Read the pipeline artifacts: spec, plan, review report.
3. Read the git log for this feature/fix.

4. Analyze the cycle and answer:

   **What worked?**
   - Which parts of the plan were accurate?
   - Which tests caught real issues?
   - What review findings were most valuable?

   **What didn't work?**
   - Where did the plan miss something?
   - What took longer than expected and why?
   - What did the review catch that should have been caught earlier?

   **Reusable pattern?**
   - Is there a pattern here that applies to future work on this project?
   - Did we solve a problem in a way worth remembering?
   - Is there a new "how we do X" that should become a project skill?

   **System update needed?**
   - Should any standard be added or modified?
   - Should CLAUDE.md or REVIEWER.md be updated?
   - Should a new skill file be created in `.pipeline/skills/project/`?

   **Tech debt check.** Look at the code touched in this cycle with fresh eyes:
   - Are there functions that grew too long or do too much?
   - Copy-paste patterns that should be extracted?
   - Naming that made sense during implementation but is unclear now?
   - Shortcuts taken to make tests pass that should be cleaned up?
   - Dependencies or abstractions that feel wrong now that the feature is complete?
   If you find debt, note it in the solution file. If the debt is significant (3+ items, or anything that will slow down the next fetta), suggest `/v-refactor` as the next action before starting the next fetta.

5. Produce `.pipeline/solutions/[date]-[slug].md`:

```markdown
---
date: [ISO date]
pipeline: [new-feature | bugfix | refactor]
feature: [name]
stack: [stack used]
tags: [relevant tags — auth, api, deploy, database, etc.]
---

# [Feature/fix name]

## Problem
[What we were solving — one paragraph]

## Solution
[How we solved it — the approach, not the code]

## What Worked
- [Concrete thing that went well]

## What Didn't Work
- [Concrete thing that didn't go well]
- [What we'd do differently]

## Reusable Pattern
[If applicable — a pattern future projects can use]

## System Updates Applied
- [What was updated in standards/CLAUDE.md/REVIEWER.md/skills, if anything]
- [Or "none needed"]

## Tech Debt
- [Items found during the tech debt check, or "none"]
- [Recommendation: "suggest /v-refactor before next fetta" or "acceptable, proceed"]
```

6. If system updates are needed, ASK the human before making them. Propose the specific change (e.g., "Add to CLAUDE.md: when working with JWT refresh tokens, always use httponly cookies") and wait for approval.

7. **Propose framework promotion.** For each system update that is general enough to apply to ALL projects (not just this one), ask: "This improvement could benefit all projects using the framework. Promote to framework master?" If the human approves, write a proposal file in `.pipeline/proposed-updates/[date]-[slug].md`:

```yaml
source_project: [project name or path]
source_date: [ISO date]
source_fetta: [fetta number or "n/a"]
source_solution: [path to solution file]
kind: [improvement | bug]               # default: improvement
target_layer: [principles | methodology | stack | template | skill | command | cli]
title: "[short title]"
description: |
  [What the improvement is and why it matters. For bugs, include a reproducer.]
proposed_diff: |
  [Concrete change — new file content, or diff to existing file. For bugs,
   this is the suggested fix; the maintainer may adapt it.]
```

`kind` defaults to `improvement`. Use `kind: bug` only when reporting a defect
in the framework or the vibbly CLI — `target_layer: cli` is reserved for that
case. Improvements proposed during compound are almost always `improvement`;
only emit `bug` if the cycle uncovered a framework defect (rare).

This schema is shared with `/v-propose` (the ad-hoc, outside-cycle entry point).
If you change one, change the other — both must accept the same proposal format.

The human will later run `vibbly propose` to move these to the shared inbox.

8. **Capture future ideas.** During the analysis, if a learning or observation suggests a feature, improvement, or experiment that's NOT part of the current work — ask: "This sounds like a future idea: '[brief title]'. Capture it with `/v-idea`? (y/n)" If yes, create the idea with source "compound learning from [feature name]".

9. Tell the human:
   ```
   Compound complete. Learning captured in .pipeline/solutions/[file].
   [N] system updates proposed. [Applied / Pending approval / None needed].
   [N] framework promotions proposed.
   [N] future ideas captured.
   
   Pipeline cycle complete: spec → plan → verify → execute → review → ship → compound ✓
   ```

   If the tech debt check found significant items, add:
   ```
   ⚠ Tech debt detected ([N] items). Consider running /v-refactor before starting the next fetta.
   See .pipeline/solutions/[file] → Tech Debt section for details.
   ```

   If `.pipeline/fette.md` exists, check for remaining fette:
   ```
   Fetta [N]/[total] complete: "[fetta name]"
   Next: run `/v-plan` to start fetta [N+1]: "[next fetta name]"
   ```
   Or if all fette completed:
   ```
   All fette completed. Project is done.
   ```

10. Archive the pipeline state — **FETTE-AWARE**:

   **If fette remain** (not all fette completed):
   - Archive ONLY the plan and review for this fetta:
     - Move `.pipeline/plan.md` → `.pipeline/archive/[date]-fetta[N]-plan.md`
     - Move `.pipeline/review.md` → `.pipeline/archive/[date]-fetta[N]-review.md`
     - Clear `.pipeline/reviews/` (external reviewer outputs)
   - **Keep** `.pipeline/spec.md` (shared across all fette — DO NOT archive)
   - **Keep** `.pipeline/fette.md` (tracks fetta progress)
   - Reset state to **between-fette** (NOT to null):
     ```json
     {
       "pipeline": "new-feature",
       "feature": "[same feature name]",
       "current_step": "plan",
       "steps_completed": [
         {"step": "spec", "completed_at": "[original timestamp]", "gate_passed": true, "artifacts": [".pipeline/spec.md"]}
       ],
       "blocked": false,
       "current_fetta": [N+1],
       "review_iteration": 0,
       "last_refactor_fetta": [preserve existing value],
       "last_refactor_date": "[preserve existing value]"
     }
     ```

   **If all fette completed** (or no fette):
   - **Close the seeding idea (if any).** Before archiving the spec, scan `.pipeline/ideas/*.md` for any file whose **frontmatter** `status:` starts with `in-progress` and references `.pipeline/spec.md` (the file about to be archived). If found, transition the frontmatter status to the **shipped** terminal state:
     ```yaml
     status: shipped YYYY-MM-DD (see .pipeline/archive/[date]-[slug]-spec.md, solution .pipeline/solutions/[date]-[slug].md)
     ```
     This closes the metadata loop opened by `/v-spec` step 2. Without this step, the idea stays stuck at `in-progress` with a dangling path reference. Same rule for a bugfix pipeline when the state had `bug_file` set — ensure the referenced bug's frontmatter `status: fixed` is set (it should be, from /v-bugfix step 4; if not, set it now).
   - Full archive:
     - Move `.pipeline/spec.md` → `.pipeline/archive/[date]-[slug]-spec.md`
     - Move `.pipeline/plan.md` → `.pipeline/archive/[date]-[slug]-plan.md`
     - Move `.pipeline/review.md` → `.pipeline/archive/[date]-[slug]-review.md`
     - Move `.pipeline/fette.md` → `.pipeline/archive/[date]-[slug]-fette.md` (if exists)
     - Clear `.pipeline/reviews/`
   - Full state reset:
     ```json
     {
       "pipeline": null,
       "current_step": null,
       "steps_completed": [],
       "blocked": false,
       "current_fetta": null,
       "review_iteration": 0,
       "last_refactor_fetta": [preserve existing value],
       "last_refactor_date": "[preserve existing value]"
     }
     ```
   - Keep `.pipeline/solutions/`, `.pipeline/skills/`, `.pipeline/todos/` — permanent knowledge

11. **Push to remote.** Compound adds solutions and may update standards — push to keep remote in sync.
    ```
    git push origin development
    ```
    If push fails: warn but don't block.

## Rules
- Be honest in the analysis. "Everything went perfectly" is almost never true.
- Focus on ACTIONABLE insights, not reflections. "The spec missed edge case X" is actionable. "We should plan better" is not.
- Keep solution files SHORT. If the pattern needs more detail, create a skill file in `.pipeline/skills/project/` instead.
- New skills always go in `skills/project/`. To promote to framework, use the framework promotion mechanism (proposed-updates).
- This step should take 2-3 minutes, not 20. It's a quick retrospective, not a thesis.
