# /v-bugfix — Bugfix pipeline: reproduce → fix → review → deploy

## Input
$ARGUMENTS (bug description, error message, OR path to a bug file from `/v-bug`)

## Step 1: Load Context

1. If $ARGUMENTS is a path to a `.pipeline/bugs/*.md` file, read it for context (description, severity, reproduction steps, analysis).
2. If $ARGUMENTS is a free-text description, treat it as an ad-hoc bug (no bug file).
3. Read `.pipeline/state.json`. If a pipeline is already active and not completed, warn the human.

## Step 2: Reproduce (TDD red phase)

1. Analyze the bug description and any available reproduction steps.
2. Identify the likely location in the codebase.
3. Write a FAILING test that reproduces the bug.
   - The test name should describe the bug: `test_user_cannot_login_with_expired_token`
   - The test must FAIL with the current code — proving the bug exists.
4. Run the test suite to confirm the new test fails and all other tests still pass.

Report:
```
Bug reproduced.
- Test: path/to/test_file.py::test_name
- Failure: [what the test shows]
- Likely cause: [your analysis]
- Severity: [from bug file or your assessment]

Proceed with fix? (y/n)
```

Set state: `pipeline: "bugfix"`, `current_step: "reproduce-test"`, `blocked: true`.
If from a bug file: set `bug_file: "<path>"` in state.

## Step 3: Fix (TDD green phase)

After human approval:

1. Fix the bug with MINIMAL changes. Do not refactor unrelated code.
2. Run the reproduction test — it must now PASS.
3. Run the FULL test suite — nothing else should break.
4. If other tests break, the fix has side effects. STOP and report.

Report:
```
Fix applied.
- Files modified: [list]
- Changes: [brief description]
- Reproduction test: PASSING
- Full suite: PASSING (N tests)

Run `/v-review` to verify fix quality.
```

Update state: mark fix completed, `current_step: "review"`.

## Step 4: Update Bug File

If the bug came from a `.pipeline/bugs/*.md` file, update it:
- Set `Status: fixed`
- Set `Fix: [brief description of what was changed]`
- Set `Test: [path to reproduction test]`

## After Fix

The standard `/v-review` and `/v-deploy` commands take over from here. The pipeline is:
reproduce-test → fix → `/v-review` → `/v-deploy`
