# /v-approve — Unblock the pipeline (human gate approval)

## Procedure

1. Read `.pipeline/state.json`.

2. If not blocked:
   ```
   Pipeline is not blocked. Nothing to approve.
   Current step: [step]. Run the next command to continue.
   ```

3. If blocked:
   - Show what's being approved (the artifact from the current step)
   - Mark the current step as completed in `steps_completed` with `gate_passed: true`
   - Set `blocked: false`
   - Advance `current_step` to the next step in the pipeline
   - Tell the human what command to run next

4. If $ARGUMENTS contains feedback (e.g., `/v-approve with changes: remove NFR-3`):
   - Apply the requested changes to the relevant artifact first
   - Then approve and advance
