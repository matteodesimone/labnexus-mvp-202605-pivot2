# Skills

Skills are **operational knowledge** — how we do specific things in our projects.

## Two levels of skills

```
.pipeline/skills/
├── framework/       # Shipped with the domain. Updated by `vibbly framework update`.
│   ├── design-patterns.md
│   └── intuitive-design.md
└── project/         # Created by /compound. NEVER overwritten by framework update.
    └── (grows over time)
```

**Framework skills** are curated reference material that comes with the domain. They describe fundamental patterns every project in this domain should know (design patterns, intuitiveness evaluation). These are few, stable, and updated only when the framework evolves.

**Project skills** are born from experience. When `/compound` identifies a reusable pattern ("how we do JWT auth with Fly.io on this project", "our Stripe webhook pattern"), it creates a skill file in `project/`. These accumulate over the life of the project and are never touched by framework updates.

## Skills vs Standards vs Contexts

**Standards** (`.pipeline/standards/`) = WHAT rules to follow

**Skills** (`.pipeline/skills/`) = HOW we do specific things (internal patterns)

**Contexts** (`.pipeline/contexts/`) = External tech documentation (fetched, not authored)

## When to create a project skill

During `/compound`, when the "Reusable Pattern" section describes something worth remembering. Good candidates: auth patterns, API client patterns, deploy configs, testing patterns, integration patterns.

## When to promote a project skill to framework

During `/compound`, if a project skill proves useful across multiple projects, propose promoting it via `.pipeline/proposed-updates/`. Promoted skills move from `project/` to `framework/` in the domain.

## Format

```markdown
# Skill: [name]

## Context
[When to use this skill]

## The Pattern
[Step-by-step or code example]

## Why This Way
[Brief justification]

## Gotchas
[What goes wrong if you deviate]
```

## Usage by agents

Agents read `.pipeline/skills/` (both `framework/` and `project/`). The separation exists only for update protection — agents treat all skills equally.
