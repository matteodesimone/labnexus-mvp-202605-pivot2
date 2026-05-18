# Anchoring LLM reviewers on current state

## When this applies

You're building a pipeline that sends structured context to an LLM
reviewer: spec (describing the problem or goal), plus the current
artifact (code, config, doc, etc.). The reviewer's job is to evaluate
the current artifact, potentially with the spec as background.

Naive prompt structure:

```
[reviewer instructions]
[spec content]
[file contents]
Evaluate.
```

## The trap

LLMs consistently read the spec text as describing CURRENT state, not
historical state. If the spec says "X needs to be removed", the LLM
flags "X is present" as a finding — even when the file contents show
X is already absent.

Vibbly logged this across 4 consecutive review cycles with Mistral:
~25 total false positives over those cycles, each requiring manual
dismissal by the primary reviewer (Claude).

## The fix: current-state framing preamble

Prepend a short paragraph (≤500 chars) that explicitly distinguishes
spec from current artifact and anchors the LLM on evaluating CURRENT
state. Apply to every prompt uniformly — single constant, one
prepend-site in the prompt builder.

Template:

```
IMPORTANT: The spec below describes the problem this cycle set out to
solve. The artifacts below represent the CURRENT state — after the
implementation step has already applied fixes. Your job is to evaluate
the CURRENT state. If the spec says something should be removed or
changed and the artifacts already show it as removed or changed, that
is success — do not report it as a finding. Report only issues present
in the current artifact.
```

## Properties to preserve

- **Prepend, don't replace**: the check-specific instructions stay.
- **Apply uniformly**: every check type gets the same preamble. Single
  named constant, one insertion site. Never inline per-check — leads
  to drift.
- **Pin with tests**: assert the preamble is in the built prompt for every
  check type. An invariant test prevents silent regression if someone
  refactors the builder later.
- **Stay positive-framed**: "report only issues present" (not "skip if
  unsure"). Under-reporting is also a failure mode; keep the reviewer
  motivated to find real issues.

## Reference implementation

Vibbly `internal/core/reviewer/prompts.go`:

- Const `currentStateFraming` (~400 chars).
- `BuildPrompt` prepends it to the sections list before the check
  template.
- 4 pinning tests in `prompts_test.go`: presence per check type,
  position within first 500 chars, source-level anchor-phrase check,
  exactly-once declaration.

## What this does NOT fix

- LLM backends that systematically refuse tasks ("I cannot perform
  security review") — unrelated problem.
- Minimal-response stubs (the LLM returns "no issues found" without
  meaningful reasoning) — this is a separate fragility of specific
  backends. Not addressed by the preamble.
- Over-reporting of trivial issues — tune per-check instructions, not
  the preamble.
- **Bug-file-driven hallucinations** (reviewer reads a `.pipeline/bugs/*.md`
  file's pre-fix description as current state during a bugfix cycle) —
  the preamble's wording focuses on spec artifacts and doesn't generalize
  to bug files. Captured as a follow-up refinement; current evidence
  suggests either extending the preamble to name bug files explicitly
  OR dropping bug files from the review artifact set entirely during
  bugfix cycles (structurally cleaner).

## Measurement

Qualitative: track false-positive count in compound reports across the
next 3+ review cycles. If FP rate stays near zero, the fix is working.
If FPs resurface, either the preamble is being ignored (consider
rephrasing) or the new pattern is different from the one this fix
addressed.

Programmatic measurement is possible (parse "Dismissed False Positives"
sections of archived review reports) but requires pipeline tooling to
aggregate — captured as future framework idea, not blocker.
