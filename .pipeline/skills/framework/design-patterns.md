# Skill: Design Patterns

Curated patterns for the projects we build. Each pattern has a trigger (when to use it) and an anti-trigger (when it LOOKS like you need it but you don't). Patterns are tools — use them when the situation demands, never to demonstrate sophistication.

Rule: if you apply a pattern in a plan, include a one-line justification. If you can't justify it in one line, you probably don't need it.

---

## Architectural Patterns

### Separation of Logic and Presentation

**What**: Business logic in pure modules with no I/O dependencies. Presentation connects logic to the outside world (HTTP, CLI, GUI, DB).

**Trigger**: always. This is not optional — it's a principle (see methodology).

**Anti-trigger**: none. Even a 20-line script benefits from separating "what to compute" from "how to print it."

**Stack examples**:
- Python/FastAPI: `service.py` (logic) → `routes.py` calls service, returns HTTP response
- Go CLI: `internal/` packages (logic) → `cmd/` calls logic, formats stdout
- Rails: service objects and models (logic) → controllers call them, render views
- Godot: `logic/` classes with signals (logic) → scene scripts connect signals to nodes
- Svelte: `lib/stores.js` + `lib/utils.js` (logic) → components bind to stores and call utils

---

### Repository

**What**: Abstract data access behind a clean interface. The business logic asks "give me user by ID" and doesn't know if it comes from sqlite, Postgres, an API, or a file.

**Trigger**: business logic contains SQL queries, ORM calls, or file reads mixed with domain decisions. Multiple parts of the codebase access the same data source differently.

**Anti-trigger**: the project has one model, three queries, and they're all in one file. Wrapping them in a Repository class adds a layer for no benefit. Just keep them in `db.py` / `db.go` and move on.

**Stack notes**:
- Python: SQLAlchemy session queries wrapped in repository functions or class. Not a separate abstraction layer on top of SQLAlchemy — just organized access.
- Go: functions that take a `*sql.DB` and return domain types. sqlc already generates this — your sqlc queries ARE your repository.
- Rails: ActiveRecord scopes and class methods on models. This IS the Rails repository — don't add another layer.

---

### Service Object

**What**: Extract a complex operation that spans multiple models/modules into a dedicated object or function. One public method, one job.

**Trigger**: a function does 5+ things in sequence involving multiple dependencies. A controller/handler contains business logic beyond "parse request, call something, return response."

**Anti-trigger**: simple CRUD. Creating a `UserCreationService` that just calls `User.create()` is ceremony with no value.

**Pattern**: `ServiceName.call(params) → Result` (Ruby) or `def service_name(params) → result` (Python) or `func ServiceName(deps, params) (result, error)` (Go).

---

## Behavioral Patterns

### State Machine

**What**: Model an entity with discrete states and explicit transitions. Each state defines allowed actions and next states.

**Trigger**: code has 4+ `if/elif/match` branches on a status field, and the branches will grow. An entity goes through a lifecycle (draft → published → archived). Game entities have behavioral modes (idle → chase → attack → flee).

**Anti-trigger**: an entity has 2-3 simple states with obvious transitions (active/inactive, on/off). A boolean or enum is enough — don't build a state machine for a toggle.

**Stack notes**:
- Godot: node-based state machine with child nodes per state, or simple script-based dictionary of states. This is the standard Godot pattern for entity AI.
- Python/Go/Rails: simple class with `current_state` and a `transition(event)` method that validates allowed transitions. No framework needed.

---

### Observer / Signal

**What**: Decouple event producer from event consumer. Producer emits, consumers subscribe. Neither knows about the other.

**Trigger**: module A needs to notify module B about something, but A shouldn't depend on B. Multiple consumers need to react to the same event.

**Anti-trigger**: only one producer and one consumer, and they already know about each other. A direct function call is simpler and more readable than a signal.

**Stack notes**:
- Godot: `signal` is native. Use it. This is THE Godot way for decoupling.
- Svelte: stores with `subscribe`. Components react to store changes.
- Python: simple callback lists, or `blinker` library if you need named signals. Don't build a custom event bus for 3 signals.
- Rails: `ActiveSupport::Notifications` or simple callbacks.
- Go: channels for async, callback functions for sync.

---

### Strategy

**What**: Define a family of interchangeable algorithms. The caller picks which one to use at runtime.

**Trigger**: same operation with genuinely different implementations that are selected based on context. Example: different pricing calculations for different customer tiers. Different export formats (CSV, JSON, PDF).

**Anti-trigger**: you have ONE implementation and you're wrapping it in a Strategy "in case we need another one later." YAGNI. Add the pattern when the second implementation actually exists.

**Stack notes**: in all languages, a function parameter or interface is sufficient. Don't create abstract classes with a factory — just pass the function.

---

## Structural Patterns

### Composition over Inheritance

**What**: Build behavior by combining small, focused pieces rather than inheriting from a deep class hierarchy.

**Trigger**: always prefer this over inheritance. Especially when you're tempted to create `BaseModel → AuthenticatedModel → AdminModel → SuperAdminModel`.

**Anti-trigger**: almost never. Shallow inheritance (one level, for framework requirements like `extends Node2D`) is fine. Deep inheritance is almost always wrong.

**Stack notes**:
- Godot: child scene components (HealthComponent, HitboxComponent, AIComponent) composed into an entity. Not a class tree.
- Python: mixins are acceptable (one level), multiple inheritance chains are not.
- Go: embedding + interfaces. Go doesn't have inheritance by design — composition is the only option and the right one.
- Rails: concerns for shared behavior (max 2-3 per model). If you have 5+ concerns, the model is doing too much.

---

### Dependency Injection

**What**: Pass dependencies from outside rather than creating them internally. The caller decides what to inject.

**Trigger**: a module creates its own database connection, HTTP client, or external service client internally. This makes testing impossible without mocks.

**Anti-trigger**: a pure function that takes values and returns values. No dependencies to inject. Also: don't build a DI container for a project with 3 dependencies.

**Stack notes**:
- FastAPI: `Depends()` — native, first-class. Use it.
- Go: constructor functions that take dependencies as params. `func NewService(db *sql.DB, logger *slog.Logger) *Service`.
- Rails: constructor injection for service objects. `UserService.new(mailer: Mailer.new)`.
- Godot: pass dependencies via `@export` vars or scene composition. Don't use autoloads as a DI container.

---

## Patterns we deliberately EXCLUDE

These are well-known patterns that add complexity without proportional value in MVP development:

- **Singleton** — global mutable state, testing nightmare, hidden coupling. Use dependency injection instead.
- **Abstract Factory** — when you have one concrete type, a factory is ceremony. Add it when the second type exists.
- **Decorator chain** — elegant in textbooks, confusing in practice for small teams. Prefer explicit composition.
- **Visitor** — almost never needed outside compilers and AST processors.
- **Mediator** — adds indirection. Direct communication (function calls, signals) is simpler.

If you genuinely need one of these, you'll know — and you'll add it to this skill file as a compound learning.

---

## How this skill is used in the pipeline

**During `/plan`**: the agent reads this skill, identifies which patterns apply to the feature being planned (based on triggers), and references them in the Architecture Decision section of the plan. Each pattern usage has a one-line justification.

**During `/review`**: the agent checks — are the patterns from the plan applied correctly? Is the logic/presentation separation maintained? Is any pattern forced where it doesn't belong (anti-trigger violated)?

**During `/compound`**: if a pattern worked well or failed, update this skill. Add notes, adjust triggers, add stack-specific gotchas.
