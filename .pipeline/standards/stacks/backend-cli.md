# Stack: Backend / CLI

For web services, MCP servers, CLI tools, and local multiplatform utilities.

## Language Selection

**Go** is preferred when:
- The deliverable needs to be distributed as a standalone binary (CLI tools, desktop utilities, cross-platform tools)
- Performance under concurrency matters (high-throughput services)
- The target machine may not have Python installed

**Python** is preferred when:
- The service runs on a server or a machine where it's installed once (web services, MCP servers, background workers)
- AI/ML integration is involved (SDK availability)
- Rapid prototyping speed matters more than runtime performance
- The standard library + a few packages cover 90% of the need

**Both** can coexist in the same project when appropriate (e.g., Python service + Go CLI tool that talks to it).

---

## Python Standards

### Runtime & Tooling
- Python 3.12+
- Package management: `uv` (preferred) or `pip` with `requirements.txt`. No Poetry unless the project already uses it.
- Linter: `ruff` (covers formatting + linting in one tool)
- Type hints: use them on public function signatures. Not required on local variables or test code.

### Frameworks (pick one per project, simplest that works)
- **No framework**: for simple scripts, CLI tools, small MCP servers. Just standard library + `click` for CLI args if needed.
- **FastAPI**: for HTTP services that need async, automatic OpenAPI docs, and Pydantic validation.
- **Flask**: for simple HTTP services where async isn't needed and FastAPI feels heavy.

### Database (simplest that works)
- **sqlite** (via `sqlite3` stdlib or via ORM): default choice. Use unless there's a concrete reason not to.
- **PostgreSQL** (via `asyncpg` or `psycopg`, or via ORM): when the project needs concurrent writes from multiple processes, advanced queries (full-text search, JSON operators), or the client explicitly requires it.

### ORM
- **SQLAlchemy is the default ORM for Python.** Used consistently across all Python projects that have a database. It handles query building, schema management, and — critically — SQL injection prevention by design.
- **Do not mix ORM and raw SQL.** If SQLAlchemy provides a way to express a query, use it. Raw SQL only as absolute last resort for queries the ORM genuinely cannot express, and even then wrapped in SQLAlchemy's `text()`.
- For projects without a database, or with trivial data needs (read a JSON file, write a CSV), no ORM is needed — the stdlib is enough.

### Migrations (Python + Alembic)
- **Alembic for all Python projects with a database.** `alembic revision --autogenerate -m "add users table"`, `alembic upgrade head`, `alembic downgrade -1`.
- Every migration must be reversible (`upgrade()` + `downgrade()`).
- Migrations live in `alembic/versions/`, committed to git.
- Tests apply all migrations from scratch before running (sqlite in-memory).
- `alembic.ini` and `env.py` configured at project init.

### Testing

**Unit + Integration**: `pytest` with fixtures. No `unittest.TestCase`.
- `conftest.py` for shared fixtures and test database setup.
- Test naming: `test_<behavior_description>` — reads like a sentence.
- Database: sqlite in-memory via SQLAlchemy for all DB tests. No mocking the database.
- Mocking: only external API boundaries (see mock zones in methodology).

**BDD (acceptance)**: `pytest-bdd` — integrates Gherkin `.feature` files with pytest.
- Feature files in `features/` directory.
- Step implementations in `tests/steps/`.
- For CLI apps: steps invoke the binary via `subprocess` with real arguments and assert on stdout/stderr/exit code.
- For HTTP services: steps use `httpx` test client against the real FastAPI app with test database seeded from fixtures.

**Test data**: loaded from `.pipeline/test-data/fixtures/` via `conftest.py` fixtures. Never generate test data inline in tests.

### Structure (small projects)
```
project/
├── main.py
├── service.py
├── db.py
├── models.py
├── config.py
├── alembic/               # Database migrations
│   ├── env.py
│   └── versions/
│       └── 001_create_users.py
├── features/              # BDD Gherkin feature files
│   ├── create_project.feature
│   └── ...
├── tests/
│   ├── conftest.py        # Shared fixtures, DB setup, fixture loading
│   ├── steps/             # BDD step implementations
│   │   ├── test_create_project.py
│   │   └── ...
│   ├── test_service.py    # Unit tests
│   └── test_db.py         # Integration tests
├── alembic.ini
├── pyproject.toml
└── README.md
```

Don't create more structure than this unless the project genuinely needs it. A 500-line `service.py` is fine. Split when a module exceeds ~300 lines AND has clearly separable responsibilities.

**Logic/Presentation Separation (Python)**:
- **Logic**: `service.py`, `models.py`, pure functions and classes. No imports from `fastapi`, `flask`, `click`, or any I/O framework. Testable with `pytest` alone — no HTTP server needed.
- **Presentation**: `routes.py` / `main.py` (CLI entry). Calls service functions, formats output. Thin — no business logic here.
- The test: can you `import service` and call its functions without starting any server or framework? If yes, the separation is correct.

### The Python Way
- List comprehensions over `map`/`filter` when readable.
- Context managers (`with`) for resource management.
- Generators for lazy sequences.
- f-strings for formatting.
- `pathlib.Path` over `os.path`.
- Explicit is better than implicit. Flat is better than nested.
- No bare `except:`. Catch specific exceptions.

---

## Go Standards

### Runtime & Tooling
- Go 1.22+
- Module management: `go mod`
- Linter: `golangci-lint`
- Formatter: `gofmt` (non-negotiable — Go's formatter is canonical)

### Frameworks (pick one)
- **No framework**: for CLI tools and simple services. `net/http` stdlib is excellent.
- **cobra** + **viper**: for complex CLI tools with subcommands and config.
- **chi** or **echo**: for HTTP services that need routing beyond stdlib.

### Database
- **sqlite** (via `modernc.org/sqlite` — pure Go, no CGO): default for CLI tools and single-process services.
- **PostgreSQL** (via `pgx`): when concurrent multi-process access is needed.
- **`sqlc` is the default data access tool for Go.** It generates type-safe Go code from SQL queries — this is the Go way (explicit, no magic, compile-time safety). It gives the security benefits of parameterized queries with the readability of plain SQL. This is different from the Python ecosystem where SQLAlchemy is preferred, because Go's culture favors code generation over runtime reflection. Each language follows its own idiomatic approach.
- **`goose` for migrations.** SQL-based migrations with up/down. Integrates naturally with sqlc — goose manages the schema, sqlc generates Go code from it. Migration files are plain SQL in `db/migrations/`.
- No ORM (GORM, ent) unless the project has very complex relationships AND the team has strong experience with the tool. The Go community consensus is that ORMs fight the language's explicitness.

### Migrations (Go + goose)
```bash
# Create a new migration
goose -dir db/migrations create add_users_table sql

# Apply all pending migrations
goose -dir db/migrations -database "sqlite3:./app.db" up

# Rollback last migration
goose -dir db/migrations -database "sqlite3:./app.db" down

# Check current version
goose -dir db/migrations -database "sqlite3:./app.db" status
```

Migration file example (`db/migrations/001_create_users.sql`):
```sql
-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;
```

Rules:
- Every migration has both `+goose Up` and `+goose Down`.
- Migrations live in `db/migrations/`, committed to git.
- Tests apply all migrations from scratch before running (in-memory sqlite).
- sqlc's `schema` config points to the migration files so it reads the current schema.

### Testing

**Unit + Integration**: standard `testing` package.
- Table-driven tests (the Go way).
- Test naming: `TestServiceName_BehaviorDescription`.
- `testify` only if assertion readability justifies the dependency.
- Database: sqlite in-memory for all DB tests. No mocking the database.

**BDD (acceptance)**: `godog` (official Cucumber BDD for Go).
- Feature files in `features/` directory (same Gherkin syntax as all stacks).
- Step implementations in `features/steps/` as Go test functions.
- For CLI apps: steps invoke the binary via `exec.Command` with real arguments, assert on stdout/stderr/exit code.
- For HTTP services: steps use `net/http/httptest` against the real server.

**Test data**: loaded from `.pipeline/test-data/fixtures/` via test helpers. YAML parsed with `gopkg.in/yaml.v3`.

### Structure (small projects)
```
project/
├── main.go
├── service.go
├── service_test.go       # Unit tests next to code (Go convention)
├── db.go
├── config.go
├── db/
│   └── migrations/       # goose SQL migrations
│       ├── 001_create_users.sql
│       └── ...
├── features/             # BDD Gherkin feature files
│   ├── init_project.feature
│   └── steps/            # Step implementations
│       └── init_steps_test.go
├── go.mod
├── go.sum
└── README.md
```

For larger projects, use the standard Go project layout with `cmd/`, `internal/`, `pkg/` — but ONLY when there are multiple binaries or genuinely reusable packages.

**Logic/Presentation Separation (Go)**:
- **Logic**: `internal/` packages. Pure Go — no imports from `cobra`, `net/http`, `fmt` for user output. Accept interfaces, return structs.
- **Presentation**: `cmd/` package. Cobra commands parse args, call internal packages, format stdout/stderr. Thin — no business logic here.
- The test: can you `import "project/internal/service"` and call its functions without any CLI or HTTP setup? If yes, the separation is correct.

### The Go Way
- Accept interfaces, return structs.
- Errors are values — handle them explicitly, don't panic.
- No global mutable state. Pass dependencies explicitly.
- Short variable names in small scopes (`r` for reader, `w` for writer). Descriptive names in larger scopes.
- Goroutines + channels for concurrency. No shared mutable state.
- `context.Context` as first param for anything cancellable.
- `defer` for cleanup.
- No `init()` unless absolutely necessary.

---

## Deployment

- **Fly.io**: for services (Python or Go). Single `fly.toml` + `Dockerfile`.
- **Go binaries**: for distributed CLI tools. Cross-compile with `GOOS`/`GOARCH`. Distribute via GitHub releases.
- **Python CLI**: `pipx` for installation on user machines. `uv` for server deployment.

## Privacy & GDPR (reminder)

Privacy by design is defined in Layer 0 (Principles) and Layer 1 (Methodology). Stack-specific reminders:

- **Logging**: never log personal data (names, emails, IPs). Use opaque identifiers.
- **API responses**: return only the fields the consumer needs. No full user objects when an ID suffices.
- **Database**: design schema with deletion in mind. Cascading deletes must work cleanly.
- **Third-party services**: evaluate GDPR compliance before integrating.
- **Test data**: never use real personal data in fixtures. Use realistic but synthetic data.

See `01-methodology.md` → "Privacy by design (GDPR)" for the full rules.
