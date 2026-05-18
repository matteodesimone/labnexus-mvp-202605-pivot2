# /v-status — Show current pipeline state

## Procedure

1. Read `.pipeline/state.json`.

2. If the file is empty or doesn't exist:
   ```
   No active pipeline. Start one with:
   - `/v-spec [requirement]` — new feature
   - `/v-bugfix [bug description]` — fix a bug
   - `/v-refactor [target or "auto"]` — improve existing code
   ```

3. If a pipeline is active, display:

```
Pipeline: [type] — [feature name]
Current step: [step name]
Status: [active | blocked (reason)]

Progress:
  [✓] step-1 — completed at [time]
  [✓] step-2 — completed at [time]  
  [→] step-3 — IN PROGRESS
  [ ] step-4
  [ ] step-5

Artifacts:
  - .pipeline/spec.md
  - .pipeline/plan.md
  - [test files]

Next action: [what to do next]
```

4. If blocked, emphasize what the human needs to do to unblock.

5. If `.pipeline/fette.md` exists, display fetta progress:

```
Fetta: [current]/[total] — "[fetta name]"
Completed: [list of completed fette]
Next: [next fetta name, or "all complete"]
```

6. Read `.pipeline/config.sh`. If `PROJECT_TYPE` is `prove-out`, also display budget:

```
Budget (Package S):
  Total: 20 gettoni (15 client + 5 infra)
  Spent: 8 gettoni (5 client + 3 infra)
  Remaining: 12 gettoni
  Current fetta: 2 (6 allocated, 3 spent)
```

Read budget data from `.pipeline/state.json` field `budget`. See `.pipeline/standards/project-types/prove-out.md` for the full budget schema.

7. **Health indicators.** Read tracking fields from `.pipeline/state.json`:

   If `last_refactor_fetta` exists and `.pipeline/fette.md` exists:
   - Calculate fette since last refactor = `current_fetta` - `last_refactor_fetta`
   - If 0-2: don't mention it
   - If 3-4: show `Fette since last refactor: [N] (consider /v-refactor if debt is accumulating)`
   - If 5+: show `⚠ Fette since last refactor: [N] — run /v-refactor or /v-audit to check for tech debt`

   If `last_refactor_fetta` doesn't exist and `current_fetta` >= 3:
   - Show `No refactoring done yet (fetta [N]). Consider /v-refactor if early fette cut corners.`

   These are suggestions, not blockers. The human decides.
