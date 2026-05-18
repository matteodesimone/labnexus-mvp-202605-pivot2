# /v-review — Multi-LLM quality review

## Pre-requisite

Read `.pipeline/standards/` before starting — all applicable layers per `standards/README.md`. Those are the rules you evaluate against.

## Scope

Review all files created or modified in the current pipeline. Read the artifacts from `.pipeline/state.json` to identify which files to review.

If $ARGUMENTS is provided, review those specific files instead.

## Phase 0: Pre-flight Health Check (FR-7)

Phase 0 runs *before* Phase 1's `vobbly hook review`. It catches obviously-broken reviewer setups in seconds rather than letting the user wait 5+ minutes for a `vobbly hook review` that was predictable to fail.

### Bypass

Check the `VREVIEW_SKIP_PREFLIGHT` environment variable. If set to one of the truthy values `1`, `true`, or `yes` (case-insensitive), **skip Phase 0 entirely**: write the line `[REVIEW] Pre-flight skipped by VREVIEW_SKIP_PREFLIGHT` to stderr and proceed directly to Phase 1. The bypass exists for "I know X is down, proceed anyway in degraded single-reviewer mode".

When the bypass is engaged, the Phase 4 report **must** prefix the "Review Method" section with `Pre-flight: bypassed via VREVIEW_SKIP_PREFLIGHT`. Without this marker, a reader of a week-old review.md cannot tell whether the safety check ran clean, was skipped consciously, or never existed.

### Procedure (when bypass not active)

1. Run: `.pipeline/bin/vobbly hook preflight`
   - **Use the Bash tool with a timeout of 60000 (1 minute). Do NOT set run_in_background.** Wait for the command to complete.
   - It calls `reviewer.SmokeTest` (the same function `vibbly doctor --reviewers` uses — single source of truth per FR-7) against every reviewer enabled in `.pipeline/config.sh`, writes `.pipeline/reviews/_preflight.json`, and emits the line `[REVIEW] PRE-FLIGHT: N/M reviewers ready[ — aborting]` to stderr.
2. Read `.pipeline/reviews/_preflight.json` to get the verdict and per-reviewer status. The schema is `{verdict ("proceed"|"abort"), ratio_threshold (0.5), ratio_failed, total_enabled, total_failed, abort_reason, reviewers, exit_code}`.

**`_preflight.json` is always written**, even on abort (exit 2 or 3), via a `defer` in `RunPreflight` — read it for the structured failure reason. The "Interruption behavior" guarantees that apply to `_run.json` (SIGINT may leave `status: running`, SIGKILL bypasses the deferred write) apply identically here.

### Exit-code handling

The exit code of `vobbly hook preflight` is distinct from Phase 1's `vobbly hook review` exit code. Treat them as separate namespaces:

| Phase 0 exit | Meaning | What you do |
|---|---|---|
| exit 0 | proceed — < 50% of enabled reviewers failed smoke | proceed normally to Phase 1. **Surface** the `verdict` plus any down reviewers in `review.md`'s "Review Method" section (one line per down reviewer with its `status` from the JSON). When `total_failed == 0`, a single positive line "Pre-flight: 3/3 reviewers ready" is enough. This is mandatory — see Bypass section on why optional surfacing is unsafe. |
| exit 2 | abort — ≥ 50% of enabled reviewers failed smoke | **abort the review.** Do NOT invoke `vobbly hook review`. Do NOT write `review.md`. Update `.pipeline/state.json` with `blocked=true` and `blocked_reason` quoting the `abort_reason` field from `_preflight.json`. Tell the user the specific failed reviewers (read `reviewers` map for status), and remind them: set `VREVIEW_SKIP_PREFLIGHT=1` to override. |
| exit 3 | catastrophic setup — zero reviewers enabled in `config.sh` (EC-1 boundary) | same as exit 2: abort with `blocked=true`, do NOT write `review.md`. The `abort_reason` explains the setup error. The bypass does NOT help here — fix the configuration. |

On Phase 0 abort, the review has not happened. Do not write `review.md` from Claude alone — the user invoked `/v-review` expecting a multi-reviewer run; if the pre-flight says the reviewers are down, the right answer is to escalate, not fabricate.

## Phase 1: Launch External Review

**⚠ CRITICAL: DO NOT use `run_in_background` for the review invocation. Run it SEQUENTIALLY. WAIT for it to finish. You CANNOT proceed to Phase 4 without external results.**

1. Identify all files to review (from state artifacts or $ARGUMENTS).
2. Run: `.pipeline/bin/vobbly hook review all <file1> <file2> ...`
   - **Use the Bash tool with a timeout of 300000 (5 minutes). Do NOT set run_in_background.**
   - Wait for the command to complete. Read the last stderr line (pattern `^\[REVIEW\] COMPLETE: …` on success or `^\[REVIEW\] SETUP ERROR: …` on setup failure) and note the exit code.
   - If it times out after 5 minutes: proceed with Claude-only review and note "external reviewers timed out".
3. Results land in `.pipeline/reviews/`:
   - `<reviewer>-<check>.md` — one file per (reviewer, check) pair.
   - `_run.json` — structured run summary (produced on every invocation). Read it to know `exit_code`, per-reviewer status, files written, and whether a `summary.md` was also produced.
   - `summary.md` — human-readable degraded-mode marker. Appears when ≥50% of configured reviewers failed.

### Exit-code handling (FR-4)

`vobbly hook review` returns a semantic exit code matching the `_run.json` `exit_code` field:

| Exit code | Meaning | What you do |
|---|---|---|
| exit 0 | all-green — every enabled reviewer produced ≥1 successful check | **proceed normally** to Phase 2/3/4. The final review.md "Review Method" section lists all reviewers without degraded notes. |
| exit 1 | partial — at least one reviewer ok, at least one failed | proceed to Phase 2/3/4. The final review.md "Review Method" section contains "partial multi-review mode" and names the failed reviewer(s). |
| exit 2 | degraded — every enabled reviewer failed every check | proceed in Claude-only mode. The final review.md "Review Method" section contains "degraded — external reviewers unavailable". |
| exit 3 | catastrophic setup — no reviewers enabled, or enabled reviewers have no invocable binary | **abort the review.** Do NOT write review.md. Update `.pipeline/state.json` with `blocked=true` and `blocked_reason` quoting the setup error message from `_run.json.error`. Notify the human with the specific setup problem. |

On exit 3, the review has not happened. Do not fabricate Phase 2/3/4 output from Claude alone — the user's intent was a multi-reviewer run, and a review.md that hides the fact that external reviewers never ran is misleading. Fix the setup and rerun.

### Interruption behavior (EC-5)

`_run.json` is written in a deferred block, so a panic or early return still produces the file. **SIGINT** (Ctrl-C) and **SIGKILL** are best-effort: SIGINT may leave `_run.json` in `"status": "running"` state, and SIGKILL bypasses the deferred write entirely. If you see a `_run.json` with `"status": "running"`, the previous invocation was interrupted — rerun `vobbly hook review` to get a fresh file.

The external reviewers (configured in `.pipeline/config.sh`) receive `REVIEWER.md` content concatenated into the prompt, along with the layered standards and project context. They see the same quality bar you do.

## Phase 2: Your Independent Review

Do YOUR review now. Do NOT read the external reviewer's output yet — form your own opinion first.

## MANDATORY CHECK before Phase 4

**Before writing the review report, VERIFY:**
1. Check: `ls .pipeline/reviews/*.md` — are there files from external reviewers?
2. If YES: read them ALL. If `summary.md` exists, treat its degraded-mode marker as an input to Phase 4. Proceed to Phase 3 (cross-reference).
3. If NO: you did NOT wait for external results. **STOP. Go back to Phase 1 and run `.pipeline/bin/vobbly hook review` again, this time WITHOUT run_in_background.**
4. A review that ignores external results is INVALID. Never write "external reviewers: running in background" — that means you didn't wait.

Evaluate every file against the layered standards in `.pipeline/standards/`:

**Principles (Layer 0)** — Check ALL 8 principles:
- *Sottrarre è moltiplicare*: is this the simplest solution? No over-engineering? No unnecessary abstractions?
- *C'è un modo*: if unconventional approaches are used, are they justified as last resort?
- *Decidi, fai, osserva, impara*: is the implementation iterative or a big-bang change?
- *Obsess over quality*: does the code meet the agreed quality metric from `02-project.md`?
- *Sicurezza by design*: are security implications considered from the start? Frameworks handling CSRF/XSS/injection by default? No custom auth/crypto when libraries exist? Secrets in env vars, not code? Least-privilege defaults?
- *Tempo fisso, scope variabile*: is anything built beyond what the current spec/fetta requires? Over-building is a scope violation.
- *Sempre funzionante*: does the implementation produce a working end-to-end experience? (Detailed check in "End-to-end completeness" below.)
- *La documentazione è codice*: does every code change that affects behavior have corresponding documentation updates? README, API docs, command help, inline comments — all match the actual behavior? Stale docs are HIGH findings.

**Methodology (Layer 1)** — Process followed correctly? Idiomatic approach? Dependencies justified? Verification passes? Horizontal slicing respected?

**Stack (Layer 2)** — Correct patterns for the framework? Right architecture? ORM used consistently?

**Project (Layer 3)** — Domain conventions? Quality metric met? Intentional deviations respected? Data handling per privacy section? If `02-project.md` has a "Project Extensions" section: apply any Custom Review Checks listed there (in addition to domain checks). Read any Custom Standards files referenced there.

Then the specific checks — in two groups:

**Universal checks (always apply):**

**Correctness** — Cross-check against `.pipeline/spec.md`. All requirements addressed? Edge cases handled? Error paths covered?

**Documentation alignment** — Does every change that affects behavior have corresponding documentation updates? README, help text, config comments match actual behavior? (Principle: la documentazione è codice.)

**End-to-end completeness (sempre funzionante)** — If `.pipeline/fette.md` exists: does this fetta produce a working journey from start to finish? Components built but not connected into a usable flow = HIGH.

**Framework coherence** — If `features/framework/coherence.feature` exists, run all coherence tests. Failing scenarios are findings.

**Domain-specific checks (from `.pipeline/checks.yaml`):**

Read `.pipeline/checks.yaml`. For each enabled check, evaluate the work against its description. The check's `lens` field tells you the focus area. If the file doesn't exist, use the methodology (Layer 1) and stack (Layer 2) standards as the check list.

Additionally, evaluate any domain-specific standards from Layer 1 and Layer 2 not covered by the named checks.

## Phase 3: Cross-Reference

NOW read ALL external reviewer output files in `.pipeline/reviews/` (excluding `tribunal-*` files which are for gate decisions). Files follow the pattern `<reviewer-name>-<check-type>.md` (e.g., `gemini-quality.md`, `mistral-security.md`, `ollama-test-coverage.md`).

For each external finding, classify it:

- **CONFIRMED**: Multiple reviewers (or you + at least one reviewer) found the same issue. Highest confidence.
- **VALID NEW FIND**: One external reviewer caught something others (including you) missed. Upon reflection, it's real. Note which reviewer found it.
- **FALSE POSITIVE**: Flagged something that isn't a problem given full project context. Exclude with brief explanation.
- **DISAGREE**: Reviewers have conflicting assessments. → Send to Phase 3.5 (Dialectic).

When multiple external reviewers participated, note the **agreement level**: a finding flagged by 3/3 reviewers is much more credible than one flagged by 1/3.

## Phase 3.5: Dialectic (only when DISAGREE findings exist)

Skip this phase if there are no DISAGREE findings.

For each DISAGREE finding, send the conflict to the involved reviewers for a **second opinion with context**:

```
Prompt to each reviewer involved in the disagreement:

"A review finding is disputed. Here are the positions:

Finding: [description]
[Reviewer A] says: [their position — quote from their review]
[Reviewer B] says: [their position — quote from their review]

Given this disagreement, reconsider your position. You may:
1. Maintain your position with additional reasoning
2. Concede to the other reviewer's position
3. Propose a synthesis (both are partially right)

Respond with your revised assessment."
```

After collecting responses:
- If a reviewer concedes → reclassify the finding (CONFIRMED or FALSE POSITIVE).
- If both maintain → keep as DISAGREE with enriched reasoning for human decision.
- If synthesis → create a new combined finding that captures both perspectives.

Note: this costs 2-4 extra LLM calls total, only when disagreements exist. In practice, disagreements are rare (typically <10% of findings). If no external reviewers are available for the second round, keep the original DISAGREE classification.

## Phase 4: Unified Report

Write `.pipeline/review.md`:

```markdown
# Review: [feature name]

## Summary
[One paragraph overall assessment]

## Review Method
- Primary: Claude (full project context)
- External reviewers: [list all that responded, e.g., Gemini, Mistral, Ollama(codellama)]
- Standards: .pipeline/standards/ (layered)
- Agreement rate: X/Y findings confirmed by 2+ reviewers

---

## Confirmed Findings (both reviewers)

### [SEVERITY] Title
- **File**: path:line
- **Issue**: what's wrong
- **Fix**: what to do
- **Why**: impact if not fixed

## Claude-Only Findings

### [SEVERITY] Title
- **File**: ...
- **Note**: not flagged by external reviewer

## External-Only Findings (validated)

### [SEVERITY] Title
- **File**: ...
- **Note**: missed by Claude, confirmed valid

## Disputed (before dialectic)

### Title
- **Claude says**: ...
- **External says**: ...
- **Sent to Phase 3.5**: yes/no (no if Phase 3.5 was skipped)

## Dialectic Outcomes (if Phase 3.5 ran)

### Title
- **Original dispute**: [brief summary]
- **Dialectic result**: RESOLVED (conceded by [reviewer]) | SYNTHESIZED | UNRESOLVED
- **Final classification**: [CONFIRMED / FALSE POSITIVE / DISAGREE — human decides]
- **Reasoning**: [key argument that resolved it, or summary of both positions if unresolved]

## Dismissed False Positives
[Brief list with reasons]

---

## Metrics
- Files reviewed: N
- Findings: N total (X critical, Y high, Z medium, W low)
- Confirmed by both: N
- Claude-only: N
- External-only validated: N
- False positives dismissed: N

## Verdict
PASS | PASS WITH NOTES | FAIL
```

## Post-Review — Iteration Logic

Read `REVIEW` and `REVIEW_MAX_LOOPS` from `.pipeline/config.sh`.

Track the current loop in `.pipeline/state.json` field `review_iteration` (starts at 1).

**CRITICAL: ALWAYS update state.json** when the review produces a verdict. This is Golden Rule #4.

**`REVIEW` values:**
- `"manual"` — single pass, human fixes
- `"critical"` / `"high"` / `"medium"` — until-clean loop at that severity threshold
- `"all"` — until-clean, fix everything including LOW. PASS WITH NOTES counts as FAIL in this mode — all findings must be resolved.

**When `REVIEW="manual"`** (default)
- **FAIL**: Set `blocked: true`. Tell human what to fix. Remind: "Write a failing test first, then fix, then verify."
- **PASS**: Mark review completed, advance to next step.

**When `REVIEW` is not "manual"** (until-clean loop)
- The REVIEW value IS the threshold. Auto-fix findings at or above it.
- **FAIL** and loop < `REVIEW_MAX_LOOPS`:
  1. Report: "Review loop [current]/[max] (threshold: [REVIEW]): [count] findings."
  2. For each finding at or above threshold, apply the **fix-and-verify pattern** (same as `/v-bugfix`):
     - Write a test/scenario that captures the issue (MUST FAIL before fix)
     - Apply the minimal fix
     - Verify it passes
     - Run full verification suite — confirm no regressions
  3. Increment `review_iteration`. Re-run from Phase 1.
  CRITICAL always escalates to human.
- **FAIL** and loop = `REVIEW_MAX_LOOPS`: escalate to human.
- **PASS**: Mark review completed.

**Report format:**
```
Review loop [N]/[max] (REVIEW=[value]): [PASS|FAIL]
Findings: [total] (critical: X, high: Y, medium: Z, low: W)
At threshold: [count] findings at [value]+
```

**Regression detection**: more findings than previous loop → escalate to human.

## Rules
- Be genuinely critical. Don't soften findings.
- Don't blindly trust external reviewer. Cross-reference with project context.
- If external reviewer was unavailable, note "single-reviewer mode" in report header.
- NEVER auto-fix CRITICAL findings without human approval, even in auto-iterate mode.
- Each re-review is a FULL review, not a partial re-check. Fixes can introduce new issues.
- **Unmasked pre-existing failures**: if a fix or feature EXPOSES failures that were pre-existing but masked (e.g. a test that previously failed at setup now fails at a merit assertion), the default is **fix inline** when the scope is small (≤30 min of typical work per failure). Do NOT auto-defer to the idea backlog. Rationale: the fix made these failures visible; shipping with a red suite violates "sempre funzionante". Only defer (via `/v-idea` or `/v-bug`) when inline fix would require new spec, architectural change, or significantly exceed the original fix's scope. Rule of thumb: ≤5 unmasked failures AND ≤30 min each → inline; otherwise escalate to the human to choose (fix, stop, or defer). Precedent: vibbly bugfix `bdd-tests-vibbly-not-in-path` (2026-04-19) exposed 2 pre-existing failures, fixed inline in ~20 min, bringing `go test ./...` from 111 failing to 0.
