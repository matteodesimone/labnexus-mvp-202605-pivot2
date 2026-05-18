# /v-propose — File a framework improvement outside a pipeline cycle

## When to Run

When you notice an issue or improvement that should land in the framework
itself (templates, standards, skills, commands, methodology, or the vibbly
CLI tool) and you are NOT in the middle of a `/v-compound` cycle. Inside a
compound cycle, step 7 already covers this — don't duplicate.

Trigger examples:
- A `vibbly init`-generated template misbehaves on first use. **(bug)**
- A standard contradicts itself or is silent on a real situation you hit. **(improvement)**
- An existing slash command is missing a step or has a footgun. **(improvement or bug)**
- A skill you wish existed would have saved time on a recurring task. **(improvement)**
- `vibbly doctor` exits with a misleading error or wrong banner. **(bug)**

`/v-propose` accepts two **kinds** of feedback:
- `improvement` — additive or modifying suggestion. The proposal carries a
  concrete `proposed_diff` against framework documentation/templates that the
  maintainer can apply mechanically.
- `bug` — a defect report. The proposal documents a reproducer plus a
  suggested fix. The framework maintainer does NOT apply the diff directly;
  instead `/v-dev-inbox` converts the proposal into a bug entry in the
  framework-dev repo (`.pipeline/bugs/`), which then enters `/v-bugfix`.

## Procedure

1. Confirm the improvement is **framework-wide**, not project-specific.
   Project-specific work belongs in `/v-idea`, `/v-bug`, or `/v-spec`.
   If unsure, ask the human one clarifying question and wait.

2. Decide the **kind**:
   - `improvement` (default) — additive/modifying change to a framework doc,
     template, skill, or command. The maintainer will apply your diff.
   - `bug` — a defect in the framework or vibbly CLI. The maintainer will
     convert this proposal into a bug entry and fix it via `/v-bugfix`.
     Use this when something is *broken*, not when something could be *better*.

3. Decide the target layer (where the fix or change is expected to land):
   - `principles` — change to core values (rare, high bar).
   - `methodology` — change to the universal process.
   - `stack` — change to a stack-specific overlay (e.g. agent-go, web-app).
   - `template` — change to a generated template (e.g. .mise.toml).
   - `skill` — new or updated skill in `skills/framework/`.
   - `command` — new or updated slash command.
   - `cli` — code change to the vibbly Go CLI (`cmd/`, `internal/`, `pkg/`).
     Only valid when `kind: bug`. CLI improvements still go through `kind: bug`
     if they fix a defect, or open a regular GitHub issue/discussion if they
     are pure feature requests.

4. Pick a slug: lowercase, hyphenated, descriptive. The filename will be
   `.pipeline/proposed-updates/YYYY-MM-DD-<slug>.md`.

5. Write the proposal file with this frontmatter (shared with `/v-compound`
   step 7 — keep the schemas in sync):

   ```yaml
   ---
   source_project: [project name or path]
   source_date: [ISO date]
   source_fetta: n/a
   source_solution: n/a (or path if a solution file exists)
   kind: [improvement | bug]               # default: improvement
   target_layer: [principles | methodology | stack | template | skill | command | cli]
   title: "[short title]"
   description: |
     [What the improvement is, why it matters, and the concrete trigger
      that surfaced it. Include error messages, file paths, or commands
      verbatim where relevant — the framework maintainer is reading this
      cold. For bugs, include a minimal reproducer.]
   proposed_diff: |
     [A concrete change. For improvements, this is what gets applied. For
      bugs, this is the suggested fix — the maintainer may adapt it.
      Either a unified diff against the current framework file, or full
      new-file contents. Avoid prose-only proposals — make the maintainer's
      job mechanical.]
   ---
   ```

6. Tell the human:
   ```
   Proposal written to .pipeline/proposed-updates/[file].
   Run `vibbly propose` from your terminal to send it to the shared inbox.
   ```

## Rules

- One proposal per file. Bundle only when the changes are tightly coupled
  (e.g. a new command + the README pointer to it). Otherwise split.
- Do NOT modify `state.json` or any pipeline artifact. `/v-propose` is
  side-effect-free with respect to the active pipeline.
- Do NOT auto-run `vibbly propose` — that's a human action, like git push.
- Quote the trigger evidence: error messages, command output, file paths.
  Future-you reading this in the inbox six weeks later will thank you.

## Relationship to /v-compound

`/v-compound` step 7 produces the same kind of file at the end of a pipeline
cycle. Both commands share the schema above; if you change one, change the
other. Suggested follow-up: extract the schema to a shared standards file
and have both commands link to it instead of re-specifying.
