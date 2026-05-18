# /v-grill — Stress-test a plan, design, or decision through relentless questioning

## When to Use

Use when:
- A plan exists (`.pipeline/plan.md`) and you want it challenged before execution
- A spec exists and you want assumptions tested before investing time
- A design decision feels uncertain and needs adversarial review
- The human says "grill me", "challenge this", "what am I missing", or "play devil's advocate"

## Input
$ARGUMENTS (optional: specific aspect to focus on, e.g., "the auth approach", "the data model", "fetta 2 scope")

## Procedure

1. Read the relevant context:
   - `.pipeline/plan.md` and `.pipeline/spec.md` (if they exist)
   - `.pipeline/fette.md` (if fette exist)
   - The codebase (if implementation has started)
   - `.pipeline/standards/` (to know what the quality bar is)
   - If $ARGUMENTS points to a specific area, focus there. Otherwise, cover everything.

2. Build a mental map of the decision tree — every choice that was made (architecture, scope, approach, tech, priority order) is a branch to walk down.

3. **Interview the human relentlessly.** One question at a time. Each question must:
   - Target a specific decision, assumption, or trade-off
   - Not be answerable with "yes" or "no" — force explanation
   - Follow up on vague answers — don't let "it should be fine" pass
   - If a question can be answered by reading the codebase or docs, read them first and then ask about what you found

   Example questions (adapt to context):
   - "The plan assumes [X]. What happens if [X] is wrong?"
   - "Why [approach A] over [approach B]? What would change your mind?"
   - "This is the most complex part. Walk me through the failure modes."
   - "Who is the user who would hate this design? Why?"
   - "You've scoped [Y] out. What if it turns out to be essential?"
   - "What's the earliest point where this could fail in production?"
   - "If you had half the time, what would you cut?"

4. Track findings as you go. Categorize each into:
   - **Held**: assumption tested, reasoning is solid
   - **Revised**: assumption was weak, human adjusted it during the grill
   - **Open**: couldn't resolve, needs more thought or research
   - **Risk accepted**: human acknowledges the risk and chooses to proceed anyway

5. When every branch of the decision tree has been walked (or the human calls it), produce the grill report.

6. Save to `.pipeline/brainstorms/[date]-grill-[topic-slug].md`:

```markdown
# Grill Report: [topic]
Date: [ISO date]
Source: [what was grilled — plan, spec, design decision, etc.]

## Decisions That Held
- [Decision]: [why it held]

## Revised During Grill
- [Original decision] → [Revised to]: [why]

## Open Questions
- [Question]: [what's needed to resolve it]

## Accepted Risks
- [Risk]: [human's reasoning for accepting it]

## Recommendations
[If the grill revealed issues: what should change before proceeding]
[If everything held: "Plan/spec is solid. Proceed."]
```

7. Tell the human the result:
   - If revisions are needed: "The grill surfaced [N] revisions. Update the plan/spec before proceeding."
   - If clean: "Plan survived the grill. Proceed with confidence."
   - If open questions remain: "There are [N] open questions that need answers before this is safe to build."

## Rules
- **Be relentless but not hostile.** The goal is to strengthen the plan, not to win an argument. If the human gives a solid answer, acknowledge it and move on.
- **One question at a time.** Don't dump a list. Each answer informs the next question.
- **Follow the branches.** If a question reveals a dependency ("well, that depends on X"), follow X before coming back.
- **Don't accept "it should be fine."** That's not an answer. Ask: "What specifically makes you confident it will be fine?"
- **Read before asking.** If the answer is in the codebase, find it yourself. Only ask the human about things that require judgment, not facts.
- **Know when to stop.** When every branch has been walked and answers are concrete, stop. Over-grilling is as wasteful as under-grilling.
