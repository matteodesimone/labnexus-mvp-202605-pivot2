# /v-audit — Multi-perspective stress test of the current state

## When to Use

On-demand. Not part of the standard pipeline loop. Run it when you want to:
- Stress-test the product after completing a batch of fette
- Validate before a major release or client delivery
- Check blind spots that standard `/v-review` misses (it checks rules; this simulates people)
- Get a fresh look after working heads-down for a while

The human calls it explicitly. It never runs automatically.

## How It Differs from /v-review

`/v-review` checks against **rules**: quality standards, security checklist, test coverage, McKay attributes. It asks "does this follow the rules?"

`/v-audit` checks against **people**: simulates real stakeholders using, maintaining, or evaluating the work. It asks "would this person be satisfied, confused, or frustrated?"

They find different bugs. `/v-review` finds "input not validated." `/v-audit` finds "the user has no idea what to do after clicking Submit."

## Procedure

1. Read `.pipeline/state.json`, `.pipeline/standards/02-project.md` (Project Context), and the current codebase/work state.

2. Read `.pipeline/perspectives.yaml` (or `.pipeline/checks.yaml` perspective entries if no separate file). Identify which perspectives are configured.

3. For each perspective, do a **focused review pass**:

   a. **Adopt the perspective fully.** You are now this person. Forget you're an engineer. Forget implementation details. Think only about what THIS person cares about.

   b. **Walk through the work** from this person's point of view:
      - What do they see first?
      - What are they trying to accomplish?
      - What confuses them?
      - What frustrates them?
      - What delights them?
      - What would make them say "this is broken" or "this doesn't solve my problem"?

   c. **Produce findings** for this perspective. Each finding has:
      - Perspective name
      - What was observed
      - Why it matters to this person (not why it matters technically)
      - Severity: CRITICAL (this person cannot accomplish their goal), HIGH (painful friction), MEDIUM (annoyance), LOW (polish opportunity)

4. **Cross-perspective analysis.** After all perspectives are done:
   - Identify findings that appear from MULTIPLE perspectives (highest confidence)
   - Identify conflicts between perspectives (e.g., user wants simplicity, maintainer wants flexibility)
   - Rank all findings by impact

5. Produce `.pipeline/reports/[date]-audit.md`:

```markdown
---
date: [ISO date]
perspectives: [list]
fetta: [current fetta or "all"]
---

# Audit Report

## Summary
[2-3 sentences: overall state from all perspectives combined]

## Perspective: [name]
### This person's goal
[What they're trying to accomplish]

### What works
- [Finding that would satisfy this person]

### What doesn't work
- [SEVERITY] [Finding] — [Why it matters to them]

## Perspective: [name]
...

## Cross-Perspective Analysis

### Confirmed by multiple perspectives
- [Finding] (seen by: [perspectives])

### Perspective conflicts
- [User wants X, Maintainer wants Y — trade-off to resolve]

### Priority ranking
1. [Highest impact finding across all perspectives]
2. ...
```

6. Tell the human:
```
Audit complete. [N] findings across [N] perspectives.

Critical: [N]  High: [N]  Medium: [N]  Low: [N]

Top 3 findings:
1. [finding]
2. [finding]
3. [finding]

Full report: .pipeline/reports/[date]-audit.md

These findings are observations, not pipeline blockers.
Use /v-spec or /v-bugfix to act on any of them.
```

## Rules

- **This is observation, not enforcement.** Audit findings don't block the pipeline. They're input for the human to decide what to act on.
- **Stay in character.** When reviewing as "end-user", don't think about code quality. When reviewing as "maintainer", don't think about UX. The whole point is the constrained perspective.
- **Be concrete.** "The user might be confused" is useless. "After clicking 'Save', nothing visible happens for 3 seconds — the user doesn't know if it worked" is useful.
- **Read before imagining.** If the perspective is "end-user" and there's an actual UI, look at it. If the perspective is "maintainer" and there's code, read it. Don't guess.
- **One pass per perspective.** Don't mix perspectives in a single pass — that defeats the purpose.
- **Findings are not todos.** The human decides which findings become work items. Some findings are accepted trade-offs.
