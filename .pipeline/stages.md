# Autonomy Model

This pipeline uses a **per-gate** autonomy model. Every decision point (gate) in the pipeline can independently be set to one of two modes:

- **human**: the pipeline blocks and waits for your approval.
- **tribunal**: two LLMs (Claude + Gemini) evaluate independently, cross-reference, and decide. Disagreement or uncertainty always escalates to human.

## How to Configure

Edit `.pipeline/config.sh`. Each gate has its own toggle:

```bash
GATE_SPEC="human"           # After spec, before plan
GATE_PLAN="human"           # After plan, before test-scaffold
GATE_TEST_SCAFFOLD="human"  # After test-scaffold, before implement (test coverage quality)
GATE_REVIEW_FAIL="human"    # When review finds critical/high issues
GATE_DEPLOY="human"         # Before deploying
GATE_COMPOUND="human"       # Before applying system updates
```

Change any gate to `"tribunal"` when you trust the process for that specific decision. Change it back anytime.

Via vibbly: `vibbly gate spec set tribunal` or `vibbly gate spec set human`.

### Per-gate reviewer selection

Each gate (and each review check) can optionally specify WHICH reviewers to use, overriding the project-level defaults. This lets you assign models to their strengths:

```bash
# Example: Gemini is strong on reasoning → use for spec/plan gates
GATE_SPEC_REVIEWERS="gemini"
GATE_PLAN_REVIEWERS="gemini"

# Mistral is strong on code analysis → use for code-facing gates
GATE_TEST_SCAFFOLD_REVIEWERS="mistral,ollama"
GATE_REVIEW_FAIL_REVIEWERS="mistral,ollama"

# Same for review checks
REVIEW_SECURITY_REVIEWERS="mistral,ollama"
REVIEW_PRIVACY_REVIEWERS="gemini"
```

Empty (default) = all enabled reviewers participate. Via vibbly: `vibbly gate spec reviewers gemini,mistral`.

## Reviewer Execution

Reviewers can run in parallel or sequential:

```bash
REVIEWER_PARALLEL=true   # all reviewers run simultaneously (default, faster)
REVIEWER_PARALLEL=false  # one at a time (easier to debug, less CPU/network)
```

Via vibbly: `vibbly config reviewer-parallel on|off`.

## The Tribunal Pattern

When a gate is in tribunal mode:

1. Claude produces the artifact (spec, plan, etc.)
2. Gemini receives the artifact with a gate-specific evaluation prompt
3. Gemini returns: **Pro** (reasons to approve), **Contro** (reasons to reject), **Suggested action** (PASS / PASS WITH NOTES / FAIL)
4. Claude reads Gemini's analysis, cross-references with its own assessment
5. Claude decides:
   - **PASS** → pipeline advances automatically
   - **PASS WITH NOTES** → pipeline advances, concerns logged
   - **FAIL** → pipeline blocks, escalates to human with full tribunal report

**Key safety properties:**
- A FAIL from either reviewer always escalates to human. The tribunal can never auto-pass through a disagreement.
- If Gemini is unavailable, the gate always escalates to human. Never auto-pass with a missing second opinion.
- Every tribunal decision is logged with reasoning in the pipeline state. You can audit why any gate was auto-passed.

## Progressive Delegation

Start with everything on "human" (the default). As you complete pipeline cycles and build confidence:

**First to delegate**: `GATE_PLAN` — if your specs are thorough, the plan is a mechanical derivation. The tribunal catches incomplete plans well.

**Second**: `GATE_TEST_SCAFFOLD` — the tribunal is particularly strong here because checking "do these tests cover the spec?" is a well-defined evaluation task with clear right/wrong answers.

**Third**: `GATE_SPEC` — once your brainstorm/spec process is mature and you trust the tribunal to catch ambiguities.

**Fourth**: `GATE_COMPOUND` — system updates are low-risk (they're additive knowledge) and the tribunal is good at spotting conflicts with existing standards.

**Last to delegate**: `GATE_DEPLOY` and `GATE_REVIEW_FAIL` — these have the highest real-world impact. Keep human until you have extensive evidence the tribunal catches what matters.

**Never delegate permanently.** Even if all gates are on tribunal, switch back to human for: new stacks you haven't used before, high-stakes client projects, anything involving payment or sensitive data, and the first project with a new client.

## Rules

- **"Obsess over quality" overrides velocity.** If delegating a gate reduces quality, take it back to human. No shame.
- **Quality metric is the arbiter.** If the client's quality metric is being met with tribunal gates, the delegation is working. If not, pull back.
- **Any gate, any time.** You can change any gate's mode at any point — even mid-pipeline. Changed your mind after seeing the plan? Switch `GATE_PLAN` back to human before running `/v-test-scaffold`.
- **The tribunal is not a rubber stamp.** It's two independent reviewers with different biases doing real analysis. If you notice the tribunal auto-passing things it shouldn't, the prompts need tightening — file a `/v-compound` note about it.

## Review loop

The `/v-review` step can loop: review → write test → fix → run tests → re-review. Configured via `REVIEW` in `config.sh`:

```bash
REVIEW="manual"       # Single pass (default). Human fixes and re-runs /v-review.
REVIEW="critical"     # Until-clean loop. Auto-fix CRITICAL only (still escalates to human).
REVIEW="high"         # Until-clean loop. Auto-fix HIGH+ findings.
REVIEW="medium"       # Until-clean loop. Auto-fix MEDIUM+ findings.
REVIEW="all"          # Until-clean loop. Fix everything.
REVIEW_MAX_LOOPS=3    # Max iterations for non-manual modes.
```

Each loop is a FULL review — not a partial re-check. Fixes can introduce new issues.

**Regression detection**: if a loop produces more findings than the previous one, the system escalates to human automatically.

**CRITICAL findings are never auto-fixed** — they always require human approval.

## Prove Out projects

When the project type in `02-project.md` is `prove-out`, the pipeline integrates with the Prove Out methodology. Full method definition in `.pipeline/standards/project-types/prove-out.md`.

### Input bridge

The workshop (3 phases: Discovery Goal, Essential Map, Exploration) produces materials that feed into the pipeline. See `prove-out.md` section "Workshop Output to Pipeline" for the full mapping. Key translations:

- **Discovery Goal** → conformity criterion in `02-project.md`
- **Sub-experiments (software-only)** → fette ordering (riskiest first) in `/v-plan`
- **Essential map steps** (derived from sub-experiments) → functional requirements in `/v-spec`
- **Tiny Experiments** (from exploration, can be questions or direct tasks) → feature definitions with truth criteria in `/v-spec`
- **Bridges** (infrastructure steps) → requirement priority: minimal in `/v-spec`
- **Pillar/Reserve** → requirement priority in `/v-spec`
- **Breadboard** → architecture starting point in `/v-plan`
- **"Non farò" and "things we won't do"** → out of scope in `/v-spec`
- **Growth plan** → ideas backlog via `/v-idea`

### Budget tracking

The pipeline state tracks gettoni: package (S/M), total budget, infra vs client split, allocated per fetta, spent per fetta. The `/v-status` command shows gettoni alongside pipeline stage. See `prove-out.md` section "Budget tracking" for the state.json schema.

### Checkpoints

Every 3-4 days: demo the working product, show remaining gettoni, revalidate Discovery Goal and priorities. At fetta delivery: propose next fetta, client approves or changes priorities. See `prove-out.md` section "Checkpoints" for the full protocol.

### Horizontal slices (fette)

The principle "sempre funzionante" applies especially here. Each fetta is a horizontal slice: it cuts across the entire user journey, so after shipping the user can go from start to finish. The first fetta is the scooter — minimal but complete end-to-end. Fette are ordered by decreasing uncertainty from the Discovery Goal's sub-experiments.
