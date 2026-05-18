# /v-implement — Write implementation to make tests pass (TDD Green Phase)

## Pre-flight Check (MANDATORY)

1. Read `.pipeline/state.json`.
   - If test-scaffold is not in `steps_completed`, REFUSE: "Cannot implement without tests. Run `/v-test-scaffold` first."

2. Run the full test suite. Count failing tests. This is your target.
   - If zero tests fail, REFUSE: "All tests pass — nothing to implement. If you need new functionality, update the plan and scaffold new tests first."

3. Read `.pipeline/plan.md` for the implementation steps and architecture decisions.

4. If `.pipeline/fette.md` exists, read `current_fetta` from state. Note which fetta you're implementing — the plan covers only this fetta. Stay within scope: do not implement features from future fette even if they seem easy to add. The principle "sempre funzionante" means this fetta must produce a working end-to-end journey when deployed.

## Implementation Procedure

Follow TDD green phase strictly:

1. Pick ONE failing test (start with the simplest).
2. Write the MINIMUM code to make that test pass.
3. Run the test suite.
4. If the target test passes, move to the next failing test.
5. If other tests broke, fix before moving on.
6. Repeat until all tests pass.

After all tests pass:

7. **Refactor phase (TDD blue/refactor):**
   - Look for duplication introduced during green phase
   - Extract common logic
   - Improve naming
   - Run tests after each refactor to confirm nothing breaks

## Rules

- **NEVER modify a test to make it pass.** If a test seems wrong, STOP and ask the human. The tests are the spec — the implementation conforms to them, not the other way around.
- **Write minimal code.** Don't add features not covered by tests. No speculative code. YAGNI.
- **Follow the plan's architecture decisions.** If you disagree with the plan, STOP and discuss with the human before deviating.
- **Check complexity as you go.** If a function exceeds 20 lines or complexity > 5, refactor immediately.

## Completion

1. Run the FULL test suite one final time. ALL must pass.

2. **Documentation check (Layer 0: "la documentazione è codice").** For every file created or modified, verify:
   - README updated if the feature changes user-facing behavior, CLI usage, or API
   - Command help text matches actual behavior
   - Inline comments on public functions are accurate
   - If config options were added/changed: config.sh comments are updated
   - **Project docs** (read config.sh DOCS_* flags): for each enabled doc, update it:
     - `DOCS_USER_GUIDE=true` → update `docs/USER-GUIDE.md` if user-facing behavior changed
     - `DOCS_DEV_GUIDE=true` → update `docs/DEV-GUIDE.md` if architecture, tooling, or workflow changed
     - `DOCS_LLM_CONTEXT=true` → update `CONTEXT.md` if components, API, state, or dependencies changed
   - If docs say one thing and code does another, FIX THE DOCS NOW — don't leave it for review.

3. Report:
   ```
   Implementation complete.
   - X tests passing (was: 0 at start)
   - Files created: [list]
   - Files modified: [list]
   ```
   If fette exist: `Fetta [N]: "[name]" — implementation done.`
   ```
   Run `/v-review` for quality check before deploying.
   ```

4. Update `.pipeline/state.json`:
   - Mark implement as completed
   - Set `current_step` to "review"
   - Set `blocked: false`
   - Record created/modified files in artifacts
