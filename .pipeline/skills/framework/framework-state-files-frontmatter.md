# Framework state files: frontmatter is canonical

## The rule

Every framework file that carries machine-readable state uses YAML
frontmatter (top-of-file `--- ... ---` block) for that state. Body
content is narrative only. No duplication of state fields in body
bullets or sections.

Applies to:

- `.pipeline/bugs/*.md` — severity, status, created, source, fix, test
- `.pipeline/ideas/*.md` — date, source, priority, status
- `.pipeline/solutions/*.md` — date, pipeline, feature, stack, tags
- `.pipeline/todos/*.md` — id, priority, status, source, created, file, bug_ref
- Any future state-bearing file type

Does NOT apply to pure-narrative files (specs, plans, reviews, reports,
brainstorms).

## Example — bug file

Correct:

```markdown
---
severity: medium
status: open
created: 2026-04-19
source: developer
---

# Bug: [title]

## Description
[narrative prose]

## Steps to Reproduce
...
```

Incorrect (dual representation — drift guaranteed):

```markdown
---
severity: medium
status: open
---

# Bug: [title]

## Severity
medium

## Resolution
- Status: open
```

## When a skill updates state

Update the frontmatter field. Never the body.

`/v-bugfix` step 4:
> Set frontmatter `status: fixed`, `fix: "..."`, `test: "..."`.
> Do NOT add or write a body `## Resolution` section with status bullets.

`/v-compound` cycle close:
> Scan `.pipeline/ideas/*.md` for the seeding idea (frontmatter
> `status:` starting with `in-progress`). Transition to `shipped ...`
> in the frontmatter.

`/v-spec` step 2:
> Update idea frontmatter `status:` to `in-progress (...)`.

## When a template creates a file

Emit frontmatter with the state fields pre-populated. No body
duplication.

## Enforcement

Write invariant tests that walk the file type and assert:

- Every file has frontmatter `status:` (or equivalent field)
- No file has the deprecated body representation

See the skill `invariant-tests-for-conventions.md` for the test pattern.
Reference: `features/steps/pipeline_metadata_lifecycle_test.go` in vibbly.

## Migration from dual representation

If a project inherits files with both frontmatter and body state:

1. Decide frontmatter is canonical (never body).
2. Script: for each file, read body `Status:` / `Priority:` / etc.,
   write to frontmatter, remove body representation.
3. Run the invariant test — find and manually patch format variants
   the script missed.
4. Commit migration in a single PR labelled clearly.

## Why frontmatter (not body)

1. Machine-parsable with standard YAML libraries (Go: gopkg.in/yaml.v3,
   Python: PyYAML, etc.). Future tooling like `vibbly bugs list --status=open`
   becomes trivial.
2. Industry standard: Jekyll, Hugo, Obsidian, Zettlr all do this.
   Contributors recognize the pattern immediately.
3. Schema-like: adding new fields doesn't distort prose.
4. Separation of concerns: metadata above, narrative below. Readers see
   the shape at a glance.

## Why NOT body bullets

1. Regex-parsing every field is fragile.
2. Duplicate representations (frontmatter + body) always drift.
3. Body format varies (`- **Status**: ...` vs `## Resolution → Status: ...`)
   → parsing ambiguity.

## Precedent

Vibbly bug `pipeline-commands-metadata-lifecycle-drift` (2026-04-19) —
canonical case study: months of drift between frontmatter and body state
across bugs/ideas, `/v-triage` mis-reporting fixed bugs as open. Fixed
by unifying on frontmatter-canonical.
