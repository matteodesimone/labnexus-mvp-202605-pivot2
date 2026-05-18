# /v-spec — Transform requirement into formal specification

## Input
$ARGUMENTS

## Procedure

1. Read `.pipeline/state.json`. If a pipeline is already active and not completed, warn the human.

2. **Check existing ideas.** Read `.pipeline/ideas/` for ideas with frontmatter `status: planned`. If the spec topic matches a planned idea, ask: "This matches idea '[title]'. Pull in its description and open questions? (y/n)" If yes, use the idea content as context for the spec. Update the idea's **frontmatter** `status:` to `in-progress (spec started YYYY-MM-DD, see .pipeline/spec.md)` — `/v-compound` will later transition it to `shipped ...` when the cycle closes.

3. Read `.pipeline/standards/02-project.md` to check the project type.

4. Analyze the requirement provided in $ARGUMENTS. Identify:
   - Core functionality (what must it do?)
   - Actors (who uses it?)
   - Inputs and outputs (data in, data out)
   - Edge cases and error scenarios
   - Ambiguities (things that are unclear or could be interpreted multiple ways)

   **If the requirement is a cleanup (remove X, retire Y, consolidate Z): audit the implementation before writing FRs.**

   Reading `--help`, API signatures, or spec descriptions tells you *intent*, not *current state*. A feature that sounds like it does something useful may be a stub; a feature that sounds dangerous to remove may already be dead.

   For each candidate-for-removal X:
   1. Read the actual implementation (`RunE` for a CLI command, the handler for an HTTP endpoint, the query for a DB lookup, etc.).
   2. Classify X as one of:
      - **Dead stub**: present but does nothing (prints "not available", returns empty, is unreachable). Safe to remove.
      - **Preview/shadow**: duplicates a canonical path elsewhere. Safe to remove if the canonical path covers all users.
      - **Genuine feature**: has real logic users depend on. Not dead — migration required, not removal.
   3. Document the classification in an FR ("audit confirmed: X is a dead stub, Y is a preview of Z, no other candidates") with the concrete evidence that led to it.

   This audit becomes a spec deliverable. The solution file in `/v-compound` should record what was read and what was concluded, so the next cleanup cycle can build on it instead of re-discovering.

   **If project type is `prove-out`**: read `.pipeline/standards/project-types/prove-out.md` section "Workshop Output to Pipeline" for the full translation rules. Key translations:
   - **Tiny Experiments** (domanda or direct task, vero_se, tasks, non farò) → functional requirements (FR) and edge cases
   - **Sub-experiment software filter** → only sub-experiments that need software become FRs; client actions are out of scope
   - **Bridges** (from map construction) → requirement priority: minimal (infrastructure to reach Pillars)
   - **Discovery Goal: "Non farò"** and **"things we won't do"** → out of scope
   - **Discovery Goal: "Starò attento a"** → edge cases (flag as high risk)
   - **Gettoni allocated** → note in spec metadata for budget tracking
   - **Pillar/Reserve classification** → requirement priority (Pillar = must, Reserve = out of scope for this fetta)

5. If ambiguities exist, LIST them clearly and ASK the human before proceeding. Do not guess.

6. Produce `.pipeline/spec.md` with this structure:

```markdown
# Spec: [feature name]

## Summary
One paragraph describing what this feature does and why.

## Actors
- [who interacts with this]

## Functional Requirements
- FR-1: [requirement]
- FR-2: [requirement]
...

## Non-Functional Requirements
- NFR-1: [performance, security, etc.]

## Input/Output
- Input: [what comes in, from where]
- Output: [what goes out, to where]

## Interface & Intuitiveness
(Include this section for any feature with a user-facing interface: GUI, web, CLI, game UI.)

### Main Instructions
For each screen, page, or command, one sentence describing what the user should do:
| Screen/Command | Main Instruction |
|---------------|------------------|
| [screen/command name] | [one clear sentence: what the user does here] |

If a main instruction can't be written in one sentence, the screen/command does too much — split it.

### Intuitiveness Checklist
For the primary interaction in this feature, verify:
- [ ] Discoverable: user can find it without searching
- [ ] Comprehensible: user understands what will happen before acting
- [ ] Forgiving: user can recover from mistakes
- [ ] Feedback: user knows what happened after acting

(Full 8-attribute review happens during /v-review. The spec captures the critical four.)

## Edge Cases
- EC-1: [what happens when...]

## Data & Privacy
- Personal data involved: [list any personal data this feature creates, reads, updates, or deletes — or "none"]
- If personal data is involved:
  - Purpose: [why this data is needed]
  - Retention: [how long it's kept]
  - Deletion: [how it gets deleted when the user requests it]
  - Third parties: [does this data flow to external services?]
- If no personal data: state "This feature does not handle personal data."

## Out of Scope
- [what this feature explicitly does NOT do]

## Open Questions
- [anything still unclear, if any remain]
```

7. **Generate proposed verification scenarios.** For each functional requirement and edge case, write a verification scenario in `.pipeline/proposed-features/`. The format depends on the domain — read `.pipeline/standards/01-methodology.md` for the expected format (e.g., BDD/Gherkin for coding, acceptance criteria for other domains).

   Rules:
   - One file per functional requirement (or per logical group if FRs are closely related).
   - Every FR must have at least one happy path scenario.
   - Every EC must have a corresponding scenario.
   - Use concrete data in examples, not abstract placeholders.
   - Write in language the stakeholder would understand. No internal jargon.
   - These are PROPOSALS — the human reviews and adjusts them in the spec gate.

8. **Propose initial verification data.** Analyze the spec and propose sample data in `.pipeline/proposed-test-data/`:
   - Ask the human 2-3 questions about the data: what does a typical record look like? What are the boundary values? What data shapes exist?
   - Based on the answers and the spec, generate a draft structure with realistic data.
   - Mark it clearly as DRAFT — the human curates and approves.

9. Initialize `.pipeline/state.json`:
```json
{
  "pipeline": "new-feature",
  "feature": "[feature name]",
  "current_step": "spec",
  "steps_completed": [],
  "blocked": true,
  "blocked_reason": "Spec produced — awaiting human review and approval"
}
```

10. Tell the human: "Spec ready for review. Read `.pipeline/spec.md` and approve or request changes. Then run `/v-plan` to proceed."
