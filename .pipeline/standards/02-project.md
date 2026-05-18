<!-- Auto-filled by `vibbly init`. Edit freely; vibbly doctor reads the key:value lines below. -->
domain: cli-tool
Stack: go-cli

# Project Standards

## Project Context

<!-- This section gives agents the big picture. Fill it during `vibbly init` or manually.
     Every command reads this file first. Without context, agents work blind. -->

- **What**: <!-- One sentence: what is this project? e.g., "SaaS for veterinary clinic management" -->
- **Why**: <!-- The problem it solves. e.g., "Vets waste 2h/day on paper scheduling" -->
- **Who**: <!-- Primary users. e.g., "Veterinary clinic staff (receptionist, vet, admin)" -->
- **Stage**: <!-- Where the project is. e.g., "MVP in production, 12 paying users" / "Greenfield, nothing built yet" -->

### Architecture Overview

<!-- High-level description of how the system is structured. 2-5 sentences.
     e.g., "Rails 8 monolith serving an API consumed by a Svelte SPA. 
     Background jobs via Solid Queue. PostgreSQL for data, Redis for caching.
     Deployed on Fly.io with auto-scaling." -->

### Key Decisions Already Made

<!-- Decisions that are settled and should NOT be revisited unless explicitly requested.
     e.g., "- Monolith, not microservices (deliberate, revisit at 10K users)"
           "- No GraphQL — REST only, simpler for our team size" -->

### Domain Vocabulary

<!-- Terms specific to this project that agents should use consistently.
     e.g., "- 'appointment' not 'booking'"
           "- 'clinic' not 'practice' or 'office'" -->

## Project Identity

- **Name**: <!-- project name -->
- **Client**: <!-- who is this for? "internal" / "client name" / "open source" -->
- **Quality metric**: <!-- how do you measure "good enough"? e.g., "all user flows complete in < 3 clicks" -->
- **Scale**: <!-- users, traffic, team size. e.g., "100 concurrent users, 2 developers" -->

## Project Type

<!-- Uncomment ONE: -->
<!-- type: prove-out        # Enables: gettoni budget, checkpoint report, input bridge from workshop -->
<!-- type: internal-tool      # No client tooling. Pipeline only. -->
<!-- type: experiment         # Lightweight. No client tooling. -->
<!-- 'gamedev' is no longer a valid project_type — game is now a DOMAIN, not a project type. Pick 'internal-tool' (or 'prove-out' if client work) and set domain: game. -->

## Domain

domain: <!-- one of: cli-tool, agent, web-app, game (library coming later) -->

## Stack Selection

Stack: <!-- depends on --domain. cli-tool: go-cli, python-cli. web-app: rails, fastapi-svelte, python-fastapi. game: godot. agent: go-claude-cli. -->

## UI Component Library

<!-- For web stacks only. Tailwind CSS is always used (mandatory). -->
<!-- Flowbite is optional — provides consistent, accessible components on top of Tailwind. -->
<!-- When enabled, agents use Flowbite components instead of building from scratch. -->
UI_COMPONENTS: <!-- flowbite | none (default: none) -->

## Database Choice

- **Engine**: <!-- PostgreSQL / SQLite / none / etc. -->
- **ORM**: <!-- ActiveRecord / Prisma / GORM / none -->

## Privacy & Data Handling

- **Personal data collected**: <!-- what PII does the system handle? "none" if CLI/internal tool -->

## Project-Specific Conventions

<!-- Naming, file organization, config format, anything the codebase follows consistently. -->

## Intentional Deviations

<!-- Things that look wrong but are deliberate. Agents and reviewers will NOT flag these. -->

## Out of Scope

<!-- What this project explicitly does NOT do. Prevents agents from building unwanted features. -->

## Deploy

- **Target**: <!-- e.g., "Fly.io", "Vercel", "GitHub releases" -->
- **Environments**: <!-- e.g., "staging + production" / "production only" -->
- **CI/CD**: <!-- e.g., "GitHub Actions" / "none (manual)" -->

## Tooling

- **Runtime management**: <!-- mise / asdf / nvm / none -->
- **Platforms**: <!-- e.g., "macOS (arm64), Ubuntu 24 (WSL2)" -->
- **Health check**: <!-- how to verify everything works -->

## Project Extensions

### Custom Commands
(none)

### Custom Review Checks
(none)

### Custom Standards
(none)
