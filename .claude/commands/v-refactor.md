# /v-refactor — Refactor pipeline: review → plan → verify → refactor → review → ship

## Input
$ARGUMENTS (what to refactor and why, or "auto" to let the agent identify candidates)

## Step 1: Check State

1. Read `.pipeline/state.json`. If a pipeline is already active and not completed, warn the human: "A pipeline is currently active ([type]: [feature]). Starting a refactor will run in parallel. Continue? (y/n)". The refactor pipeline uses its own state key to avoid conflicts.

## Step 2: Review Current State

1. If $ARGUMENTS is "auto":
   - Scan the codebase for code smells: long functions, high complexity, duplication, poor naming, architecture violations.
   - Rank findings by impact.
   - Present top 5 candidates to human for selection.
   - WAIT for human to choose.

2. If $ARGUMENTS specifies a target:
   - Analyze the target code thoroughly.
   - Identify all issues: complexity, duplication, coupling, naming, pattern violations.
   - Identify all callers/dependents of the code being refactored.

Report:
```
Refactor analysis: [target]
- Current issues: [list]
- Dependents: [files that use this code]
- Risk assessment: [what could break]

Proceed with refactor plan? (y/n)
```

Set state: `pipeline: "refactor"`, `current_step: "review-current"`, `blocked: true`.

## Step 2: Plan Refactor

After human approval:

1. Produce a refactor plan in `.pipeline/plan.md`:
   - What changes to make (specific, file by file)
   - What patterns to apply
   - Order of operations (to minimize intermediate breakage)
   - What tests need updating (if public interfaces change)

2. Set state: `blocked: true`, reason: "Refactor plan awaiting approval."

## Step 3: Verify Tests

After human approves plan:

1. Run the FULL verification suite. Record baseline: all tests must pass.
2. If any tests are missing for the code being refactored, WRITE THEM FIRST.
   - This is the safety net. Refactoring without tests is reckless.
3. Run the suite again to confirm new tests pass with current (pre-refactor) code.

Report: "Safety net in place. X tests covering refactor target. Proceeding."

## Step 4: Refactor

1. Apply changes incrementally:
   - Make ONE change
   - Run tests
   - If tests pass, continue
   - If tests fail, undo and rethink
2. NEVER change behavior. If a test fails, either:
   - Your refactor changed behavior (fix the refactor, not the test)
   - The test was testing implementation details (flag to human, get approval before changing test)

## After Refactor

Standard `/v-review` and `/v-ship` commands take over. Run `/v-review` to verify refactor quality before deploying.

When the refactor cycle is fully complete (after ship), update `.pipeline/state.json` with:
```json
{
  "last_refactor_date": "[ISO date]",
  "last_refactor_fetta": [current fetta number or null]
}
```
These fields persist across pipeline resets — they're tracking data, not pipeline state.
