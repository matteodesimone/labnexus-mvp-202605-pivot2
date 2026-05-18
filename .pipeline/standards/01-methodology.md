# Development Methodology

These rules apply to ALL projects regardless of stack. They are the direct consequences of the core principles applied to software development.

## TDD + BDD (non-negotiable)

Development is driven by tests at multiple levels. The process combines TDD (test-driven development) for implementation correctness with BDD (behaviour-driven development) for feature correctness.

### The BDD→TDD flow

1. **Spec produces Gherkin feature files.** Every functional requirement becomes a `Given/When/Then` scenario in plain language. These are proposed during `/v-spec` and reviewed during the spec gate. The client can read and validate them — they ARE the shared quality metric.
2. **Test-scaffold writes step implementations + unit tests.** The BDD steps drive integration/acceptance tests. Unit tests cover pure logic.
3. **Implement makes everything green.** BDD scenarios pass = the user can do what the spec says. Unit tests pass = the logic is correct.
4. **Refactor while ALL tests stay green.** Breaking any test — BDD or unit — is never acceptable.

A task is NOT done until the ENTIRE test suite is green: unit, integration, AND acceptance.

### Test layers (all required)

**Acceptance tests (BDD — Gherkin feature files)**
The most valuable tests. They describe what the USER can do, in language the CLIENT understands. Every functional requirement from the spec has at least one scenario. These tests run the real application — CLI gets invoked as a subprocess, web app gets hit via HTTP, frontend gets driven by a browser.

**Integration tests**
Test components working together with real infrastructure. Real database (sqlite in-memory for speed, or test database). Real HTTP calls to the test server. Real filesystem operations in temp directories. These catch the wiring bugs that unit tests miss.

**Unit tests**
Test pure logic in isolation: algorithms, validations, transformations, calculations. The smallest scope — one function, one behavior.

**No test is optional.** A feature without acceptance tests is not tested. A function without a unit test is not tested. "It works when I try it manually" is not a test.

### Mock rules (the three zones)

Mocking is a tool, not a default. Most agent-generated tests over-mock, which means they test the mocks, not the system. Follow these zones strictly:

**NEVER mock (zone 1):**
- The database. Use sqlite in-memory or a test database. It's fast, it's real, it catches real bugs.
- The filesystem. Use temp directories. They're free and they're real.
- Internal modules. If module A calls module B, test them together. Mocking B means you're not testing the actual behavior.
- The HTTP layer in integration tests. Use a real test client against the real app.

**Mock only at the boundary (zone 2):**
- External APIs you don't control (Stripe, SendGrid, OAuth providers, third-party webhooks). Mock the HTTP boundary, not the internal client code. Use contract tests or recorded responses (VCR pattern) to ensure your mock matches reality.
- Time-dependent operations. Mock the clock, not the function that uses it.
- Randomness. Seed the RNG, don't mock the function.

**Test with the real application (zone 3 — acceptance):**
- Start the actual application process.
- Feed it real test data (from `.pipeline/test-data/`).
- Interact with it as a user would: CLI commands, HTTP requests, browser actions.
- Assert on observable outcomes: output, database state, files created, HTTP responses.

**The litmus test for any mock:** if you removed the mock and used the real thing, would the test still pass? If yes, the mock is unnecessary — remove it. If no, you need the mock — but document WHY.

### Test data

Test data is curated, not random. It lives in `.pipeline/test-data/` and is designed to exercise specific behaviors.

```
.pipeline/test-data/
├── README.md           # Describes datasets, their purpose, and how to use them
├── fixtures/           # Structured seed data for database/state
│   ├── users.yaml
│   ├── projects.yaml
│   └── ...
├── scenarios/          # Data organized by BDD scenario
│   ├── happy-path/
│   ├── edge-cases/
│   └── error-cases/
└── golden/             # Expected outputs for golden file testing
    └── ...
```

**Fixture creation flow:**
1. The developer designs the core test data based on the spec requirements and edge cases.
2. An LLM can assist by proposing an initial dataset — asking clarifying questions about data shapes, relationships, boundary values, and realistic scenarios. The LLM's proposal is a starting point, never the final word.
3. The developer reviews, adjusts, and approves the fixtures. They are committed to the repo as curated artifacts.
4. Tests reference fixtures by name, not by generating data inline. This makes tests deterministic and readable.

**Fixture rules:**
- Fixtures represent realistic data, not trivial placeholders. A test user is "Maria Rossi" with a real-looking email, not "test_user_1" with "test@test.com".
- Edge case fixtures are explicit: an empty name, a 256-character email, a unicode string with emoji, a date at midnight UTC. Each exists for a reason documented in README.md.
- Fixtures never contain real personal data. They are realistic but fictional. GDPR applies even to test data if it's based on real people.

## The Right Way

Code must be idiomatic for its language. Python code follows the Python way. Ruby code follows the Ruby way. JavaScript follows the JavaScript way. Go code follows the Go way. GDScript follows the Godot way.

This means: use the language's conventions, standard library idioms, community patterns, and established tools. Don't write Java in Python. Don't write Python in Go. If you don't know the idiomatic way, learn it before writing the code.

## Style Enforcement

Style compliance is automated via the appropriate linter for the stack. Code must pass the linter before being considered final. No exceptions, no "we'll fix it later."

The linter is not optional tooling — it's part of the build. If it doesn't pass the linter, it doesn't exist.

## Modularity (balanced)

Code is organized logically without falling into the trap of hyper-fragmentation. Micro-components that increase cognitive complexity are as harmful as monolithic files.

**The test**: can a developer understand what a module does by reading its name and first 10 lines? If yes, the granularity is right. If a module's purpose requires reading 5 other modules to understand, it's too fragmented. If a module does 4 unrelated things, it's too monolithic.

## Separation of logic and presentation

Business logic NEVER depends on presentation, I/O, or infrastructure. This is the most important architectural rule across all stacks, all project types, all domains.

**Logic** is pure computation: calculations, validations, state transitions, business rules, domain decisions. It imports nothing related to HTTP, CLI, GUI, rendering, databases, or file systems. It is testable without starting any server, any engine, or any framework.

**Presentation** is everything that connects logic to the outside world: HTTP handlers, CLI commands, UI components, game scene scripts, database queries, file I/O. Presentation CALLS logic, never the other way around. Presentation translates between the external world and the logic layer.

**The connection**: signals, events, callbacks, or simple function calls. The logic emits events or returns values. The presentation listens and reacts. The logic never knows or cares what the presentation does with the results.

**The test**: take any logic module. Can you test it with ZERO imports from your framework, your ORM, your game engine, your HTTP library? If yes, the separation is correct. If you need to import `fastapi`, `http`, `Node2D`, `ActiveRecord`, or `net/http` to test pure business logic, the boundary is in the wrong place.

This applies identically to:
- **Web**: service modules (logic) vs route handlers (presentation)
- **CLI**: internal packages (logic) vs cmd/ entry points (presentation)
- **Games**: logic classes/resources (logic) vs scene scripts (presentation)
- **Frontend**: stores/utils (logic) vs components (presentation)

## Design patterns

The pipeline uses a curated set of design patterns to guide architectural decisions. These are documented in `.pipeline/skills/framework/design-patterns.md`. The `/v-plan` command consults this skill when designing solutions, and `/v-review` verifies correct application.

Patterns are tools, not obligations. A pattern is used ONLY when its trigger condition is met. Forcing a pattern where it doesn't fit is over-engineering — the opposite of "sottrarre è moltiplicare." Every pattern applied in a plan must have a one-line justification.

## Intuitive interfaces

Every interface — GUI, web, CLI, API output, error message — must be intuitive. The definition is objective, not subjective (following Everett McKay's "Intuitive Design"):

**An interface is intuitive when the target user can use it without resorting to reasoning, memorizing, experimenting, seeking help, or training.**

If the user has to guess what a button does, read a manual to understand a command, experiment to discover a feature, or memorize a sequence — the interface has failed. This applies equally to a web dashboard, a CLI tool, a game menu, and an API error response.

### The eight attributes (McKay)

Every interface is evaluated against these eight attributes. They form a checklist for `/v-review` and a design guide for `/v-plan`. See `.pipeline/skills/framework/intuitive-design.md` for the full skill with per-interface-type examples.

1. **Discoverability** — can the user find the feature without searching? Is the starting point obvious?
2. **Affordance** — does the element visually/textually communicate what it does? Does a button look clickable? Does a command name describe its action?
3. **Comprehensibility** — does the user understand what will happen before acting? Are labels clear? Is jargon absent?
4. **Responsive feedback** — does the interface confirm what happened after the user acts? Is success visible? Is failure explained?
5. **Predictability** — does the interface behave consistently? Same action, same result? No surprises?
6. **Efficiency** — can the user accomplish the task with minimum steps? No unnecessary screens, confirmations, or clicks?
7. **Forgiveness** — can the user recover from mistakes? Is undo available? Are destructive actions confirmed? Are errors recoverable?
8. **Explorability** — can the user try things safely without fear of breaking something? Is the interface safe to explore?

### Main instruction

Every screen, page, or command has a **main instruction** — one sentence that tells the user what to do. If you can't write the main instruction for a screen in one clear sentence, the screen is doing too much. The main instruction is documented in the spec for every feature that has a user interface.

### The intuitive test

For any interface element, ask: "would a target user who has never seen this before know what to do?" If the answer is "probably not," the design needs work. This test is applied during `/v-review` alongside code quality checks.

## Dependencies

Dependencies fall into two categories with different rules.

**Foundational dependencies (frameworks, ORMs, core libraries)**: these SIMPLIFY your work by absorbing complexity. A good framework (FastAPI, Rails, SvelteKit) is a force multiplier — it handles routing, security, validation, and a hundred edge cases you'd otherwise code yourself. An ORM handles SQL injection prevention, query building, migrations, and schema management. These are not "added complexity" — they are removed complexity.

Foundational dependencies are chosen once per ecosystem and used consistently. Don't change ORM every project. Don't mix frameworks. Don't use an ORM for half the queries and raw SQL for the other half — if the ORM provides a way to do it, use the ORM's way.

**Utility dependencies (small libraries for specific tasks)**: these must be evaluated carefully. Before adding one, answer:
1. Can I achieve this with the framework or standard library in under 50 lines?
2. Is this dependency actively maintained?
3. Does it pull in a tree of transitive dependencies?
4. What happens to my project if this dependency is abandoned?

If the answer to #1 is yes, don't add the dependency.

## Ecosystem coherence

Pick ONE ecosystem per project domain and go deep. Consistency across projects reduces cognitive load, reduces bugs, and makes you faster over time.

Concrete ecosystem choices:
- **JavaScript, not TypeScript**. Fewer tools, less config, fewer failure modes, less build complexity. TypeScript adds a compilation step, type gymnastics for third-party libraries, and a configuration surface that contradicts "sottrarre è moltiplicare." JavaScript with JSDoc type hints where needed gives 80% of the safety with 20% of the complexity.
- **One ORM per language, used consistently**. Pick the one that prioritizes readability for the developer. Use it everywhere. Don't mix with raw SQL unless the ORM genuinely cannot express the query.
- **One framework per domain**. FastAPI for Python APIs. Rails for Ruby full-stack. SvelteKit for frontend SPAs. Don't switch frameworks between projects without a strong reason.
- **One test framework per language**. pytest for Python. minitest or rspec for Ruby (pick one, never both). vitest for JavaScript.

The cost of switching tools between projects is invisible but real: different patterns to remember, different gotchas, different debugging approaches. Consistency compounds.

## Simplicity and compactness

Technical choices are always guided by the search for the simplest, most readable solution. Over-engineering is strictly forbidden.

Examples of this principle in action:
- sqlite over PostgreSQL unless the project specifically needs concurrent writes, full-text search at scale, or specific Postgres features
- A single file over a multi-file module if the total is under 200 lines
- A function over a class if there's no state to manage
- A shell script over a Python script if it's under 30 lines
- Environment variables over a config management system if there are fewer than 10 settings

## Appropriate scale

The application is designed for its actual scale — not an imagined future scale. If the app serves dozens of concurrent users, the infrastructure and code choices must reflect that reality without adding premature complexity.

**Premature scaling is over-engineering.** A message queue for 10 requests/second is over-engineering. A microservices architecture for a 5-endpoint API is over-engineering. A distributed cache for data that fits in memory is over-engineering.

Scale up when you observe the need (decidi, fai, osserva, impara), not when you imagine it.

## Security by design

Security is a primary guiding principle in every implementation and design choice. It's not a layer added at the end — it's a constraint from the start.

This means:
- Input validation at every boundary
- Authentication and authorization checked on every protected operation
- Secrets never in code, never in logs, never in error messages
- Principle of least privilege everywhere
- Dependencies audited for known vulnerabilities
- HTTPS/TLS for all network communication
- Use the framework's and ORM's built-in security features (CSRF protection, parameterized queries, XSS escaping) — this is why we use frameworks

## Database migrations (non-negotiable for projects with a DB)

Every project with a database uses a migration system. No exceptions. No "I'll just edit the schema manually." Migrations are version-controlled schema changes that can be applied forward and rolled back. This is a safety net as fundamental as git.

### Migration principles

**Every schema change is a migration.** Adding a column, creating a table, adding an index, changing a type — all migrations. Never modify the database schema directly. The migration history IS the schema documentation.

**Every migration must be reversible.** Every `up` has a `down`. If you add a column, the down removes it. If you create a table, the down drops it. If a migration genuinely can't be reversed (dropping data), document WHY in the migration file and mark it as irreversible explicitly.

**Migrations are sequential and immutable.** Once a migration is committed and shared (pushed to development or main), it is NEVER modified. If a migration was wrong, create a new migration that fixes it. Editing a deployed migration breaks other environments.

**Migrations run in CI/test.** The test suite applies all migrations from scratch before running tests. This catches migration bugs before they hit production.

**Seed data is separate from migrations.** Migrations change structure, seeds populate data. Seeds are re-runnable and idempotent. Test fixtures (in `.pipeline/test-data/`) are not seeds — they're for testing only.

### Migration workflow in the pipeline

During `/v-plan`: the plan identifies which schema changes are needed and includes them in the implementation steps.

During `/v-test-scaffold`: tests are written assuming the schema changes exist (they'll fail until implemented).

During `/v-implement`: migrations are created FIRST, applied, then the code that uses the new schema is written. This order matters — code can't use a column that doesn't exist yet.

During `/v-review`: migrations are reviewed for reversibility, safety (no data-destructive changes without explicit justification), and correctness.

During `/v-deploy`: migrations are applied as part of the deploy process. If they fail, the deploy stops. Rollback plan is always: `migrate down` to the previous version.

### Migration tool per stack

- **Python**: Alembic (integrates with SQLAlchemy). Always. `alembic revision --autogenerate`, `alembic upgrade`, `alembic downgrade`.
- **Go**: goose. SQL-based migrations. `goose create`, `goose up`, `goose down`. Integrates naturally with sqlc.
- **Rails**: ActiveRecord migrations. Built-in. `rails generate migration`, `rails db:migrate`, `rails db:rollback`.
- **Godot** (if using sqlite): plain SQL migration scripts executed by a helper function. Numbered files: `001_create_save_table.sql`, `002_add_score_column.sql`. Up and down in the same file separated by `-- +goose Up` / `-- +goose Down` style markers, or in separate `up/` and `down/` directories.

## Privacy by design (GDPR)

Defined in Layer 0 (Principles). The implementation consequences for every project:

- **Data model**: design with deletion in mind. Can you delete a user and all their data cleanly? If not, the schema is wrong.
- **Consent tracking**: if the app collects personal data, track what was consented to, when, and make it revocable.
- **Logging**: audit your logs. No names, emails, phone numbers, IPs in application logs. Use opaque identifiers.
- **Third-party services**: every external service that receives personal data needs a data processing agreement. Evaluate before integrating, not after.
- **Data retention**: define how long each category of personal data is kept. Implement automatic purging.
- **Cookie banner / consent**: if the frontend uses cookies beyond strictly necessary ones, implement proper consent (not a "we use cookies" dismissable banner — real opt-in).

These are not "nice to have" — they are legal obligations in the EU and good practice everywhere.

## Git strategy

See `.pipeline/standards/git-strategy.md` (core). Git rules are universal across all domains.

## Sviluppo per fette orizzontali (horizontal slices)

Defined in Layer 0 ("Sempre funzionante"). The implementation consequences for every project:

**The principle**: after every deploy, the product is usable end-to-end. The user can go from start to finish. It may be reduced — a scooter, not a car — but it works.

**How to slice**: in a User Story Map, the user flow runs left to right. A horizontal slice cuts across the entire flow, selecting one thin layer of features that together produce a complete journey. The first slice (fetta 1) picks the minimum version of each essential step. Subsequent slices enrich the journey.

**In the pipeline**: `/v-plan` proposes the slicing before detailing the implementation. Each fetta is one pipeline cycle (spec→plan→test→implement→review→deploy). After deploying a fetta, the product is demonstrable — to the client, to testers, to yourself.

**The wrong approach** (vertical slicing): building one area in depth before connecting it to the rest. "Backend is done, now let's do frontend" — this leaves the product non-functional until the last piece connects. Build the thinnest horizontal layer that works end-to-end, then go deeper.

**The test**: after deploying this fetta, can the user do the primary task from start to finish? If not, the slice is vertical, not horizontal. Rethink the cut.

## Review iteration (review → fix → re-review)

A single review pass catches most issues, but corrections can introduce new inconsistencies. The review loop addresses this: review, fix, re-review until clean or until a configured limit.

**Why it matters**: in a multi-reviewer system (Claude + external LLMs), each reviewer finds different issues from different perspectives. Fixing issue A from reviewer 1 may create issue B that reviewer 2 catches on the next pass. A single review is necessary. Multiple reviews are sometimes necessary. Infinite reviews are never necessary.

**Configured via `REVIEW` in `config.sh`:**

- `"manual"` — single pass. If it fails, the human fixes and re-runs `/v-review`. Default — simplest, most control.
- `"high"` — until-clean loop at HIGH+ severity. Auto-fix and re-review up to `REVIEW_MAX_LOOPS` times. Escalates to human if still failing.
- `"medium"` — until-clean loop at MEDIUM+. More thorough.
- `"all"` — fix everything. Most autonomous, but uses the most LLM calls.

`REVIEW_MAX_LOOPS` (default 3) caps the iterations. CRITICAL findings always escalate to human regardless of the configured threshold.

**Regression detection**: if a fix iteration produces MORE findings than the previous iteration, it means fixes are introducing problems faster than resolving them. This always escalates to human, regardless of mode — it indicates a structural issue that auto-fix can't solve.

**Auto-fix uses the same fix-and-verify pattern as `/v-bugfix`**: write a failing test/scenario that captures the issue → apply the minimal fix → verify it passes → run the full verification suite. This pattern applies regardless of what's being fixed — code, docs, or config. The domain defines what "test" and "verification suite" mean (in coding: unit tests + BDD; in other domains: their verification mechanism).

For trivial fixes (typo, formatting): apply directly — the re-review loop verifies.

### BDD coherence tests for framework files

The framework's documentation and configuration form a system with internal contracts. BDD scenarios express these contracts in human-readable Gherkin, with step implementations as lightweight shell assertions.

Coherence test feature files live in `features/framework/` and cover:
- **Naming consistency**: all toggle commands use the same pattern, all stack files have the same sections
- **Cross-reference completeness**: every command in `.claude/commands/` appears in README, every principle in Layer 0 is checked in `/v-review`, every review check in Claude is also in REVIEWER.md
- **Numbering integrity**: FR, EC, NFR, RD numbers are sequential with no gaps
- **Terminology consistency**: "horizontal slice" used correctly everywhere, no dead references to old names
- **Config completeness**: every config variable referenced by a command exists in `config.sh`

Example feature file (`features/framework/coherence.feature`):

```gherkin
Feature: Framework internal coherence
  The framework's documentation, configuration, and commands
  must be internally consistent.

  Scenario: All principles are verified in review
    Given the principles file "00-principles.md" defines N principles
    When I check the review command's Layer 0 section
    Then it should list exactly N principle checks

  Scenario: All pipeline commands are documented in README
    Given there are command files in ".claude/commands/"
    When I read the README command tables
    Then every command file should appear in the README

  Scenario: Toggle commands use consistent syntax
    Given the spec defines toggle commands
    Then all binary toggles should use "on|off" format
    And all gate mode switches should use "set" keyword

  Scenario: Every stack has a logic/presentation section
    Given there are stack files in ".pipeline/standards/stacks/"
    Then each stack file should contain "Logic/Presentation Separation"

  Scenario: FR numbering is sequential
    Given the spec defines functional requirements
    Then FR numbers should run from 1 to N with no gaps or duplicates

  Scenario: External reviewer covers all review checks
    Given the review command checks for "Intuitiveness"
    And the review command checks for "End-to-end completeness"
    And the review command checks for "Migrations"
    Then REVIEWER.md should reference each of these

  Scenario: MANUAL-OPS covers all vibbly commands
    Given vibbly has commands "init", "doctor", "status", "config", "reviewer", "gate", "framework update", "propose", "hook review", "git sync"
    Then MANUAL-OPS.md should document a manual equivalent for each
```

Step implementations are bash one-liners or short scripts — the same checks we'd do manually during review, but automated and repeatable. They run as part of `/v-review` and can be added to the pre-commit hook.

This closes the TDD loop for the entire system: **unit/integration tests for code, BDD coherence tests for framework files, re-review for everything.**

**CRITICAL findings are never auto-fixed.** They always require human approval, even in auto-iterate modes. CRITICAL means "will cause bugs, data loss, or security vulnerability" — that's not a decision to delegate.

## Multi-reviewer dialectic

When multiple LLM reviewers participate, their independent analyses sometimes conflict. The framework resolves this with a dialectic process.

### Contextual lens

Each reviewer receives not just the code and standards, but also the project's **quality metric** (from `02-project.md`) and a **contextual lens** specific to their check type. A security reviewer is told "think adversarially about data protection." A quality reviewer is told "evaluate against the agreed quality metric." This prevents generic reviews and focuses each reviewer on what matters for THIS project.

### Independent analysis first

Reviewers analyze independently — they don't see each other's output. This prevents groupthink and ensures genuine diversity of perspective. Four reviewers that parrot the first one's findings are worse than one good reviewer.

### Dialectic on disagreements

When reviewers disagree (finding classified as DISAGREE in cross-reference), the conflict is sent back to the involved reviewers with both positions visible. Each reviewer can: maintain their position with additional reasoning, concede, or propose a synthesis. This is not a vote — it's a structured discussion that produces either a resolution or a well-articulated disagreement for human decision.

### Why this works

Different models have genuinely different biases: one may be stronger on security patterns, another on code structure, another on edge cases. Independent analysis captures these different perspectives. Dialectic on disagreements ensures conflicts are resolved through reasoning, not arbitration by a single model. The result is higher-quality findings than any single reviewer could produce.
