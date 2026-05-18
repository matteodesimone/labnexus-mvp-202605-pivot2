# /v-idea — Capture, list, and review ideas for future work

## Input
$ARGUMENTS (idea description, OR "list", OR "review")

## Modes

### Capture: `/v-idea <description>`

Quick capture. Don't overthink it — the point is to not lose the thought.

1. Generate a slug from the description (e.g., "offline mode for mobile" → `offline-mode-mobile`).
2. Create `.pipeline/ideas/[date]-[slug].md`. State (date, source, priority, status) lives in the **YAML frontmatter** at the top — consistent with bugs/solutions/todos across the framework. Body sections are narrative only.

```markdown
---
date: [today, ISO date]
source: [where it came from — brainstorm, spec review, user feedback, shower thought, compound learning from X]
priority: unset
status: captured
---

# [Title]

## Description

[The idea, in enough detail to understand it later.]

## Why it's interesting

[What problem does it solve? What opportunity does it open?]

## Open questions

[What would need to be answered before implementing?]
```

Do NOT duplicate `date`/`source`/`priority`/`status` as bullet points in the body — that representation is deprecated. Other commands (`/v-spec`, `/v-compound`, `/v-idea list`) read and update the frontmatter.

3. Report: "Idea captured: `.pipeline/ideas/[filename]`. Use `/v-idea list` to see all ideas, `/v-idea review` before your next brainstorm."

### List: `/v-idea list`

1. Read all `.pipeline/ideas/*.md` files.
2. Show a summary table:

```
Ideas ([count]):
  [priority] [date] [title] — [first line of description]
  ...

Priority: ★★★ = high, ★★ = medium, ★ = low, · = unset
Status: captured | considered | planned | in-progress | shipped | declined | promoted to bug
```

Status transitions:
- `captured` — freshly captured, not yet reviewed
- `considered` — reviewed, keep in backlog as-is
- `planned` — slated for a pipeline cycle
- `in-progress (spec started YYYY-MM-DD, see .pipeline/spec.md)` — set by `/v-spec` when an idea seeds a cycle
- `shipped YYYY-MM-DD (see .pipeline/archive/..., solution .pipeline/solutions/...)` — **terminal success state**, set by `/v-compound` when the seeding cycle completes and archives
- `declined` — reviewed and decided not to pursue
- `promoted to bug (see .pipeline/bugs/...)` — reclassified as a bug

Sort by: priority (high first), then date (newest first). Unset priority sorts last.

### Review: `/v-idea review`

Interactive review of all ideas. For each idea with status "captured" or "considered":

1. Show the full idea.
2. Ask: "What do you want to do with this idea?"
   - **skip** — skip for now (don't change anything, show again in next `/v-idea review`)
   - **keep** — mark as reviewed: set status to "considered" (won't show as unreviewed next time)
   - **promote** — set status to "planned", optionally set priority
   - **spec** — immediately start `/v-spec` with this idea as input
   - **bug** — promote to a bug: create `.pipeline/bugs/<slug>.md` with severity assessment, set idea status to "promoted to bug (see .pipeline/bugs/<slug>.md)"
   - **merge** — combine with another idea (ask which)
   - **decline** — set status to "declined" with reason
   - **priority [high/medium/low]** — set priority without changing status

### Set priority: `/v-idea priority <slug> <high|medium|low>`

Quick priority update without full review.

## Integration with other commands

### /v-brainstorm reads ideas

At the START of `/v-brainstorm`, before exploring the new topic:
1. Check `.pipeline/ideas/` for ideas with status "captured" or "considered".
2. If any exist, mention: "You have [N] captured ideas. Any of these relevant to this brainstorm?"
3. List titles only (not full content) so the human can say "yes, pull in the offline-mode one."

### /v-spec reads ideas

At the START of `/v-spec`, before writing the spec:
1. Check `.pipeline/ideas/` for ideas with status "planned".
2. If the spec topic matches a planned idea, ask: "This matches idea '[title]'. Pull in its description and open questions? (y/n)"
3. If yes, pre-fill the spec context with the idea content.
4. Update the idea status to "in-progress" and link to the spec.

### /v-compound captures ideas

During `/v-compound`, if a learning suggests a future feature or improvement:
1. Ask: "This sounds like a future idea. Capture it? (y/n)"
2. If yes, create the idea with source "compound learning".

## Storage

```
.pipeline/ideas/
├── 2026-03-15-offline-mode-mobile.md
├── 2026-03-20-multi-language-support.md
└── 2026-03-22-plugin-system.md
```

Ideas are committed to git (they're project artifacts, not ephemeral state).

## Rules
- Capture is fast. Don't force structure on a raw idea — description and "why it's interesting" are enough.
- Review is optional. Ideas can sit for months. That's fine.
- Priority is optional. Most ideas won't have one until review.
- An idea that becomes a spec is still kept (status: "in-progress") for traceability.
- A shipped cycle ends with the seeding idea at status `shipped YYYY-MM-DD (see archive, solution)`. `/v-compound` sets this automatically — do NOT leave ideas stuck at `in-progress` after their cycle closes.
- Declined ideas are kept (status: "declined") with reason — you might reconsider later.
