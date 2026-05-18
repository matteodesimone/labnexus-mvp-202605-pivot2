# /v-brainstorm — Explore and clarify an idea before formalizing it

## When to Use

Use this BEFORE `/v-spec` when:
- The requirement is vague ("I want something that helps with X")
- You're exploring a problem space, not implementing a known solution
- You need to clarify scope, users, or approach before committing

Skip this and go straight to `/v-spec` when the requirement is already clear.

## Input
$ARGUMENTS (the idea, problem, or fuzzy requirement)

## Procedure

0. Read `.pipeline/state.json`. If a pipeline is active, note it but don't block — brainstorming is pre-pipeline and doesn't conflict.

1. Read the codebase (if it exists) to understand current context.

   **When the topic involves distribution, install, release, versioning, update, or packaging** — audit existing infrastructure that might already cover the need before exploring the new idea:
   - `Makefile` / `justfile` targets (list them)
   - `.goreleaser.yml`, `install.sh`, packaging scripts
   - `.github/workflows/` — release, publish, deploy
   - Existing CLI subcommands (`<cli> release`, `<cli> update`, etc.)
   - Existing git tags / GitHub Releases

   Present the findings during exploration. Often the answer is "activate / repair / document what exists" rather than "build new" — this naturally steers toward *sottrarre è moltiplicare*.

2. Research the idea. Identify:
   - What problem does this solve?
   - Who has this problem?
   - What exists already that partially addresses it?
   - What are the possible approaches?

3. Ask the human clarifying questions ONE AT A TIME. Don't dump 10 questions — ask one, get the answer, let it inform the next question. Questions to explore:
   - What does success look like?
   - Who is the primary user?
   - What's the simplest version that would be useful?
   - What's explicitly NOT needed?
   - What constraints exist (time, tech, budget)?

4. After 3-5 questions (or when clarity is reached), switch to **spirito critico** mode. Stop being accommodating. Challenge the idea relentlessly:
   - "Why this and not [simpler alternative]?"
   - "What happens if [core assumption] is wrong?"
   - "Who would NOT use this, and why?"
   - "What's the hardest part, and are you underestimating it?"
   - "Is this a real problem or a solution looking for a problem?"
   - Walk each branch of the decision tree. Don't let vague answers pass — follow up until you reach a concrete answer or the human explicitly says "I don't know yet."
   - If something can be answered by reading existing code or docs, read them instead of asking.
   - The goal is NOT to kill the idea. The goal is to make it survive scrutiny so it arrives at `/v-spec` battle-tested. An idea that can't survive five hard questions shouldn't become a spec.

5. When the idea has survived the grill (or been reshaped by it), synthesize into a brainstorm document.

6. Save to `.pipeline/brainstorms/[date]-[topic-slug].md`:

```markdown
# Brainstorm: [topic]
Date: [ISO date]

## The Problem
[One paragraph: what problem are we solving and for whom]

## Key Insights from Discussion
- [Insight 1]
- [Insight 2]

## Possible Approaches
### Approach A: [name]
- How: [brief description]
- Pro: [advantage]
- Con: [disadvantage]
- Complexity: simple | medium | complex

### Approach B: [name]
...

## Recommended Approach
[Which approach and why — reference "sottrarre è moltiplicare" if relevant]

## Stress Test Results
[What was challenged during spirito critico phase, what survived, what changed]
- Assumption tested: [assumption] → [held / revised / dropped]
- ...

## Ready for Spec
[Yes/No. If yes, the key inputs for /v-spec are:]
- Core requirement: [one sentence]
- Primary actor: [who]
- Key constraints: [list]
- Out of scope: [list]
```

7. Check `.pipeline/ideas/` for ideas related to this brainstorm topic. If matches exist, mention them: "Related ideas in backlog: [titles]. Consider merging or referencing them in the spec."

8. Tell the human: "Brainstorm captured. When ready, run `/v-spec [core requirement]` to formalize."

## Rules
- This is a CONVERSATION, not a form. Be curious, not procedural.
- Apply "sottrarre è moltiplicare" — push toward the simplest approach.
- **Spirito critico is mandatory, not optional.** Every brainstorm must include a challenge phase. Being nice is not being helpful — an untested idea that fails at implementation wastes everyone's time.
- Don't grill AND explore at the same time. First explore (steps 1-3), then grill (step 4). The human needs space to think before being challenged.
- If the idea is already clear enough for a spec, say so but STILL do a quick grill before suggesting `/v-spec`.
- No pipeline state tracking for brainstorm — it's pre-pipeline.
