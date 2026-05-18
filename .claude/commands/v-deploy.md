# /v-deploy — Pre-check and deploy

## Pre-flight Checks (ALL must pass)

1. Read `.pipeline/state.json`. Review must be in `steps_completed` with gate_passed: true.
   - If not, REFUSE: "Cannot deploy without passing review. Run `/v-review` first."

2. Run the full test suite. ALL tests must pass.
   - If any fail, REFUSE and report which tests fail.

3. Check for uncommitted changes. Warn if working tree is dirty.

4. Check for `# TODO` or `# FIXME` in files modified during this pipeline.
   - If found, WARN but don't block (these might be pre-existing).

5. If the project has a linter configured (ruff, flake8, eslint), run it.
   - If errors (not warnings), REFUSE and report.

## Deploy Procedure

1. Report pre-flight results:
   ```
   Pre-flight check:
   ✓ Review passed
   ✓ All tests passing (N tests)
   ✓ Working tree clean
   ✓ No new TODOs
   ✓ Linter clean
   
   Ready to deploy. Confirm? (This will deploy to Fly.io)
   ```

2. WAIT for human confirmation. Set state `blocked: true`, reason: "Awaiting deploy confirmation."

3. On confirmation:
   - Run `fly deploy` (or the project's deploy command)
   - Monitor output for errors
   - Report success or failure

4. On success, update `.pipeline/state.json`:
   - Mark deploy completed
   - Set `current_step: "done"`
   - Set `blocked: false`

5. **Fetta tracking.** Read `current_fetta` from state and `.pipeline/fette.md` if they exist.
   - Update `fette.md`: mark current fetta as completed with timestamp.
   - Determine if there are more fette remaining.

6. **Budget tracking (prove-out only).** Read `.pipeline/config.sh`. If `PROJECT_TYPE=prove-out`:
   - Update budget in state: increment gettoni spent for this fetta.
   - Calculate remaining gettoni.

7. Report:
   ```
   Deploy complete.
   Pipeline: [feature name]
   Steps: spec → plan → test-scaffold → implement → review → deploy ✓
   ```
   
   If fette exist:
   ```
   Fetta [N]/[total]: "[fetta name]" completed.
   Next fetta: [N+1] — "[next fetta name]"
   ```
   Or if last fetta: `All fette completed.`

   If prove-out:
   ```
   Budget: [spent]/[total] gettoni spent. [remaining] remaining.
   Run `/v-checkpoint` to prepare the client report.
   ```

   Always:
   ```
   Run `/v-compound` to capture learnings from this cycle.
   ```

8. Update state: mark deploy completed, set `current_step: "compound"`, `blocked: false`.
   Do NOT archive yet — the `/v-compound` step closes the cycle.
   If the human skips `/v-compound`, archiving happens on the next `/v-spec` or `/v-status`.

9. **Push to remote.** The work is tested, reviewed, and shipped — safe to push.
   ```
   git push origin development
   ```
   If push fails (no remote, auth error, network): warn but don't block. The push is a safety measure, not a gate.
