# Internal Tool — Project Type Definition

An internal tool built for the developer or team. No external client, no budget, no checkpoints.

This is the default project type. When `PROJECT_TYPE` is `internal-tool`, the pipeline runs without any prove-out-specific behavior.

---

## What It Is

A tool, script, library, or system built for internal use. No client relationship, no business experiment, no gettoni. Quality is self-defined — the developer sets the quality metric in `02-project.md`.

## What Changes in the Pipeline

| Step | Modification |
|---|---|
| `/v-spec` | No workshop translation. Requirements come directly from the developer. |
| `/v-plan` | No breadboard. No gettoni allocation. Fette are optional — single-fetta is fine for small tools. |
| `/v-status` | No budget display. |
| `/v-bug` | No budget-free note. |
| `/v-report` | No client report type available. Summary and full reports only. |

## What Doesn't Apply

- Gettoni and packages
- Discovery Goal
- Checkpoints
- Pact
- Client reports
- Tiny Experiment format (features are just features)

## When to Use

- Building tools for yourself or your team
- Open source projects
- Internal automation
- Framework development (like vibbly itself)
