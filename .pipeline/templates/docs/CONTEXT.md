# CONTEXT.md — LLM Project Context

> Paste this file into any LLM conversation to give it full context about this project.
> This file is auto-generated at project init and kept aligned with the code.
> It is the single source of truth for any AI that works with this codebase.

## Project

- **Name**: [project name]
- **Description**: [one line]
- **Type**: [prove-out / internal-tool / experiment]
- **Stack**: [e.g., Python FastAPI + SQLite, Go CLI, Rails 8, Godot 4]
- **Domain**: [cli-tool / agent / web-app / game]
- **Flavor**: [fullstack / backend-only — web-app only; otherwise omit]

## Architecture

```
[directory structure — key dirs only, with one-line annotations]
```

### Key Components

[List main modules/packages and their responsibility. One line each.]

### Data Flow

[How data moves through the system: input → processing → output.]

## API / Interface

[If the project has an API: list main endpoints with method, path, and purpose.]
[If CLI: list main commands with arguments.]
[If library: list main public functions/classes.]

## Dependencies

[Key external dependencies and why they're used. Not the full list — just the important ones.]

## Configuration

[Environment variables, config files, and what they control.]

## Quality Standards

This project follows the vibbly pipeline framework standards:
- **Principles**: see `.pipeline/standards/00-principles.md`
- **Methodology**: see `.pipeline/standards/01-methodology.md`
- **Stack rules**: see `.pipeline/standards/stacks/[stack].md`
- **Project specifics**: see `.pipeline/standards/02-project.md`

Quality metric: [the agreed metric from 02-project.md]

## Current State

[What's built, what's in progress, what's planned. Update after each fetta.]

### Completed

- [Feature/fetta 1: description]
- [Feature/fetta 2: description]

### In Progress

- [Current work]

### Planned

- [Upcoming work]

## Known Limitations

[What doesn't work yet, known bugs, intentional shortcuts.]

## How to Work With This Project

If you're an AI assistant helping with this project:

1. Read `CLAUDE.md` for build agent instructions
2. Read `.pipeline/standards/` for quality rules (Layer 0 → Layer 3)
3. Read `.pipeline/spec.md` for the current feature spec
4. Read `.pipeline/plan.md` for the implementation plan
5. Follow TDD: tests first, then implementation
6. Never modify `main` branch — work on `development`
