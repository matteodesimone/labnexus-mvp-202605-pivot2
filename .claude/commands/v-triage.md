# /v-triage — Convert review findings and bugs into trackable work items

## Procedure

1. Collect items to triage from TWO sources:

   **Source A: Review findings** — Read `.pipeline/review.md` if it exists. Parse all findings.
   
   **Source B: Open bugs** — Read `.pipeline/bugs/*.md`. Collect any with `Status: open`.

   If neither source has items, REFUSE: "Nothing to triage. Run `/v-review` or `/v-bug` first."

2. Present each item ONE AT A TIME to the human (review findings first by severity, then bugs by severity):

   ```
   [REVIEW] Finding 1/N: [SEVERITY] [Title]
   File: [path]
   Issue: [description]
   
   Action? (a)ccept / (s)kip / (e)dit priority / (d)efer / (b)ugfix
   ```

   For bugs:
   ```
   [BUG] Bug 1/N: [SEVERITY] [Title]
   File: .pipeline/bugs/[slug].md
   Issue: [description]
   
   Action? (a)ccept as todo / (f)ix now → /v-bugfix / (d)efer / (s)kip
   ```

   - **accept**: add to todos with current severity as priority
   - **skip**: discard — not worth fixing
   - **edit**: change priority (P1/P2/P3) before adding
   - **defer**: add to todos with status "deferred" and a note about why
   - **bugfix** / **fix now**: immediately start `/v-bugfix` pipeline for this item


4. After triaging all findings, create/update `.pipeline/todos/` directory.

   Each accepted todo becomes a file:

   ```markdown
   ---
   id: [NNN]
   priority: P1 | P2 | P3
   status: ready | in-progress | done | deferred
   source: review | bug | assessment
   created: [ISO date]
   file: [path to affected file]
   bug_ref: [path to .pipeline/bugs/*.md if from a bug, empty otherwise]
   ---

   # [Short title]

   ## Issue
   [Description of the problem]

   ## Fix
   [What to do — from the review finding or bug analysis]

   ## Context
   [Any additional context from triage discussion]
   ```

   File naming: `[NNN]-[priority]-[status]-[slug].md` (e.g., `001-p1-ready-fix-sql-injection.md`)

5. Report:
   ```
   Triage complete.
   Accepted: N (P1: X, P2: Y, P3: Z)
   Skipped: N
   Deferred: N
   
   Todos in .pipeline/todos/
   Run /v-todo-export to export as Markdown or CSV.
   ```

## Rules
- Present findings in severity order (critical first).
- Don't argue with skip decisions — the human knows the context.
- If there are more than 15 findings, suggest focusing on P1s first.
