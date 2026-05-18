# [Project Name] — Developer Guide

> This file is auto-generated at project init and kept aligned with the code.
> Update it when development workflow, architecture, or tooling changes.

## Prerequisites

- [Runtime] (managed by mise — see `.mise.toml`)
- [Package manager]
- git

Run `vibbly doctor` (or check manually) to verify your setup.

## Getting Started

```bash
# Clone
git clone [repo-url]
cd [project]

# Install dependencies
[install command]

# Run tests
[test command]

# Run the project
[run command]
```

## Architecture

[High-level overview: main components, how they connect, where to find things.]

```
[directory structure with annotations]
```

## Development Workflow

This project uses the vibbly pipeline framework. See `.pipeline/README.md` for details.

```
/spec → /plan → /test-scaffold → /implement → /review → /deploy → /compound
```

### Branch Model

- `main` — production, only human pushes
- `development` — working branch, agents work here
- `feature/<slug>` — isolation when needed

### Running Tests

```bash
[test command with examples]
```

### Code Style

[Linter and formatter info. Link to stack standards.]

Enforced automatically — see `.pipeline/standards/stacks/[stack].md`.

## Contributing

1. Create a feature branch from `development`
2. Write tests first (TDD)
3. Implement
4. Run `/review` or submit for review
5. Merge to `development`

## Key Files

| File | Purpose |
|------|---------|
| `CLAUDE.md` | Instructions for Claude Code agent |
| `REVIEWER.md` | Instructions for review agents |
| `.pipeline/config.sh` | Pipeline configuration |
| `.pipeline/standards/` | Quality standards (layered) |
