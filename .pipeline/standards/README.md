# Standards

This project uses a layered standards system. Read from top to bottom. Higher layers override lower layers when in conflict.

## Reading Order

```
Layer 0: .pipeline/standards/00-principles.md        ← Always applies. Core values.
Layer 1: .pipeline/standards/01-methodology.md       ← Always applies. Dev process.
Layer 2: .pipeline/standards/stacks/<stack>.md        ← Applies per stack choice.
Layer 3: .pipeline/standards/02-project.md            ← Applies to this specific project.
```

## Which Stack File?

Read the `Stack Selection` field in `02-project.md` to determine which stack file applies:

| Selection | File |
|-----------|------|
| `backend-cli` | `stacks/backend-cli.md` |
| `webapp (FastAPI + Svelte)` | `stacks/webapp.md` (Option A) |
| `webapp (Rails 8)` | `stacks/webapp.md` (Option B) |
| `godot` | `stacks/godot.md` |

> Paths in this table are relative to the project's
> `.pipeline/standards/` directory after `vibbly init` (the layered
> overlay copies stack files there). In the framework source repo,
> stack files live at
> `framework/domains/<domain>/.pipeline/standards/stacks/<stack>.md`.

## Conflict Resolution

If a project-specific rule in Layer 3 contradicts a methodology rule in Layer 1, the project rule wins — but it MUST be documented in the "Intentional Deviations" section with a justification.

Core principles (Layer 0) are never overridden. They are the why. Everything else is the how.

## For Reviewers

When reviewing code:
1. Read all applicable layers before starting
2. Evaluate against Layer 0 principles first (is this the simplest solution?)
3. Then Layer 1 methodology (Process followed? Idiomatic approach?)
4. Then Layer 2 stack rules (correct patterns for the framework?)
5. Then Layer 3 project rules (domain conventions? intentional deviations?)
6. If a finding conflicts with "Intentional Deviations" in Layer 3, it's a false positive — skip it
