# /v-plan — Transform spec into development plan

## Procedure

1. Read `.pipeline/state.json`. 
   - **Between-fette resume**: if `current_step` is "plan" AND spec is in `steps_completed` AND `current_fetta` > 1 AND `.pipeline/fette.md` exists — this is a resume after completing a fetta. Skip steps 2-6 (spec is already read, fette already sliced). Jump to step 7 to detail the next fetta.
   - If `current_step` is not "spec" or spec is not in `steps_completed`, REFUSE: "Cannot plan without an approved spec. Run `/v-spec` first."
   - If `blocked` is true and reason is about spec approval, REFUSE: "Spec not yet approved. Review `.pipeline/spec.md` and approve first."

2. Read `.pipeline/spec.md` thoroughly.

3. Read `.pipeline/standards/02-project.md` for project type and stack.
   **If project type is `prove-out`**: read `.pipeline/standards/project-types/prove-out.md`. Check if the spec references a breadboard or sketch from the workshop. If present, use it as the starting point for architecture — the breadboard already captures places (screens/states), affordances (user actions), and connections (flow). The plan refines this into implementation detail, it doesn't redesign from scratch. Order fette by decreasing uncertainty from the Discovery Goal's sub-experiments.

4. Analyze the current codebase to understand:
   - Existing architecture and patterns
   - Where new code should live (which files/modules)
   - What existing code can be reused
   - What needs to be modified vs created from scratch

   **Pre-implementation greppage (avoid duplicate-feature plans).** Before writing a plan that ADDS a check, EXTENDS an existing behavior, or otherwise touches functionality the spec implies might already exist, grep the codebase for the feature name + spec keywords. If pre-existing implementation surfaces, plan to integrate or remove the duplicate, not parallel-add. Cost of skipping this check: ~30 min of scaffold/impl back-and-forth when the duplication is discovered later (precedent: vibbly Fetta 4 Decision B "extend doctor to flag gamedev project_type" — already done in Fetta 1; a 30-second grep at plan time would have caught it).

5. **Consult skills.** Read `.pipeline/skills/` — if skill files exist that are relevant to the current task (design patterns, deployment patterns, etc.), use them to inform the plan. Reference any applicable pattern in the Architecture Decision section. If no skill applies, that's fine — don't force one.

6. **Propose horizontal slices (fette).** Before detailing the plan, analyze the spec and propose how to slice it into fette — horizontal slices that each deliver something usable end-to-end. In a User Story Map, the user flow runs left to right; a horizontal slice cuts across the entire flow, giving the user a complete journey at a thin level of depth.

   Rules for slicing:
   - Each fetta must produce a working product. After shipping a fetta, the user can go from start to finish — even if the experience is minimal.
   - Fetta 1 is the "scooter" — the thinnest possible end-to-end path through ALL the essential components. It may be ugly and minimal, but it WORKS.
   - Each subsequent fetta enriches the experience — more features, better UX, more robustness.
   - A fetta should be completable in one pipeline cycle (spec→plan→verify→execute→review→ship).
   - If the spec is small enough to be a single fetta, say so — no forced splitting.

   Present the proposed slicing to the human:

   ```
   Proposed slicing into N fette:

   Fetta 1 — "The Scooter" (estimated: X hours)
     What works after this: [end-to-end description]
     FRs included: FR-1, FR-3, FR-7 (minimal versions)
     
   Fetta 2 — "[name]" (estimated: X hours)
     What works after this: [what's added to the experience]
     FRs included: FR-2, FR-5, FR-8
     
   Fetta 3 — "[name]" (estimated: X hours)
     What works after this: [what's added]
     FRs included: FR-4, FR-6, FR-9, FR-10

   Does this slicing make sense? Should I adjust?
   ```

   Wait for human approval on the slicing before proceeding.
   Save the slicing to `.pipeline/fette.md`.

7. **Detail the plan for the CURRENT fetta.** After slicing is approved, produce `.pipeline/plan.md` for the first (or next) fetta only — not the entire spec.

```markdown
# Plan: [feature name] — Fetta [N]: [fetta name]

## Fetta Overview
- FRs covered: [list]
- After this fetta: [one sentence — what the user can do end-to-end]
- Estimated effort: [hours]

## Architecture Decision
How this fetta fits into the existing codebase.

### Logic/Presentation Separation
Where does the logic live? Where does the presentation live? How do they connect?
(This section is mandatory for every plan.)

### Design Patterns Applied
| Pattern | Where | Justification (one line) |
|---------|-------|-------------------------|
| Repository | db access for users | Multiple modules query users differently — centralize access |
| (none beyond logic/presentation separation) | — | Simple CRUD, no complex patterns needed |

### Embedded Artifact Dependency Check (when applicable)

If the feature involves embedding a compiled artifact, generated file, or asset into a binary via `go:embed`, `include_bytes!` (Rust), `@resource` (JVM), or similar, document the dependency analysis here — BEFORE implementation:

1. Identify the package P that will host the embed directive.
2. Identify the artifact A being embedded.
3. Identify the process that produces A (typically a compile step).
4. **Verify**: does the compilation of A depend, directly or transitively, on package P?
   - If YES: it's a build cycle. P cannot host the embed directive. Options:
     - Create a dedicated package E (standalone, not imported by A) to host the embed directive. See skill `go-embed-extract-cycle-safe` for the Go pattern.
     - Move artifact generation to a separate process (external script, build tool) that doesn't go through the language's own compiler.
   - If NO: P is safe as the embed home.
5. **Verify in Go**: `go list -deps ./cmd/<A-binary> | grep <P-package-path>` should be empty for the safe case.

Document the decision explicitly in this section, even when there's no cycle — so a future reader can see the check was done. This prevents repeat-discoveries at implement-time (precedent: vibbly's vobbly-local-install fetta 2026-04-19, where the cycle was found only after the first `go build` failed).

## Files to Create
| File | Purpose | Layer (logic/presentation) | Dependencies |
|------|---------|---------------------------|-------------|
| path/to/service.py | order processing logic | logic | models |
| path/to/routes.py | HTTP endpoints | presentation | service |

## Files to Modify
| File | Change | Reason |
|------|--------|--------|
| path/to/existing.py | what changes | why |

## Implementation Steps (ordered)
1. **Step name**: what to do, which file, expected outcome
2. ...

## Migration Steps (if DB changes needed)
1. Create migration: `[migration tool command]`
2. What schema changes
3. Reversibility: how to roll back

## Test Strategy
For each implementation step:
- BDD scenarios (which feature file, which scenarios)
- Integration tests (what components working together)
- Unit tests (what pure logic to test)
- What test data fixtures are needed from `.pipeline/test-data/`

## Risk Assessment
- What could go wrong
- What to watch out for
- Performance implications

## Estimated Complexity
Simple (< 2h) | Medium (2-8h) | Complex (> 8h)
```

8. Update `.pipeline/state.json`:
   - Mark spec step as completed
   - Set `current_step` to "plan"
   - Set `blocked: true`, reason: "Plan produced — awaiting human review"
   - Set `current_fetta` to the fetta number

9. Tell the human: "Plan ready for Fetta [N]. Review `.pipeline/plan.md`. Approve or request changes. Then run `/v-verify` to proceed."

## Between Fette

After completing a full pipeline cycle (ship) for a fetta:
- `/v-plan` for the next fetta reads `.pipeline/fette.md`, finds the next fetta, and produces a new `.pipeline/plan.md`.
- The previous plan is archived in `.pipeline/archive/`.
- The fette.md is updated with completion status.
- If the human wants to re-slice (priorities changed, scope changed), they can re-run the slicing step.
