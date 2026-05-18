# /v-test-data — Generate and curate test data fixtures

## When to Use

Run after `/v-spec` (when you know what data the feature needs) and before `/v-verify` (so verification has fixtures to load). Can also be run independently to enrich existing data.

## Input
$ARGUMENTS: optional focus area (e.g., "users", "edge cases for auth", "empty states")

## Procedure

1. Read `.pipeline/spec.md` to understand what data the feature needs.
2. Read `.pipeline/proposed-features/` to understand the BDD scenarios and what data each scenario requires.
3. Read existing fixtures in `.pipeline/test-data/fixtures/` (if any) to understand the current data shape.

4. **Ask the human clarifying questions** (2-4 questions, one at a time):
   - What does a typical [entity] look like in the real world? (e.g., "What does a typical project look like? Name, size, how many items?")
   - What are the boundary values? (e.g., "What's the longest project name you'd expect? What if a project has zero items?")
   - What states can the data be in? (e.g., "Can a project be archived? Deleted? Draft?")
   - What relationships exist? (e.g., "Does a project always have an owner? Can it have multiple?")

5. **Generate a proposed fixture set.** Based on the spec, BDD scenarios, and human answers:

   Create files in `.pipeline/test-data/fixtures/` (YAML format):

   ```yaml
   # .pipeline/test-data/fixtures/users.yaml
   # Purpose: test users for auth and permission scenarios
   
   - id: user-maria
     name: "Maria Rossi"
     email: "maria.rossi@example.com"
     role: admin
     created_at: "2025-01-15T09:00:00Z"
     notes: "Primary test user — admin with full access"
   
   - id: user-luca
     name: "Luca Bianchi"
     email: "luca.bianchi@example.com"
     role: viewer
     created_at: "2025-02-20T14:30:00Z"
     notes: "Restricted user — read-only access"
   ```

   Create scenario-specific data in `.pipeline/test-data/scenarios/`:

   ```yaml
   # .pipeline/test-data/scenarios/edge-cases/empty-project.yaml
   # Purpose: test behavior when project has no items
   
   project:
     id: project-empty
     name: "Empty Project"
     items: []
     owner: user-maria
   ```

6. **Generate README.md** for the test data:

   ```markdown
   # Test Data
   
   ## Fixtures
   | File | Purpose | Used by scenarios |
   |------|---------|-------------------|
   | users.yaml | Auth and permission testing | all |
   | projects.yaml | Project CRUD testing | create, update, delete |
   
   ## Scenarios
   | Directory | Purpose |
   |-----------|---------|
   | happy-path/ | Standard successful flows |
   | edge-cases/ | Boundary values, empty states, limits |
   | error-cases/ | Invalid data, permission denied, not found |
   
   ## Conventions
   - IDs are human-readable slugs (user-maria, project-empty)
   - Dates are ISO 8601 UTC
   - All data is fictional — no real personal information
   - Each fixture has a `notes` field explaining why it exists
   ```

7. **Present to human for review:**
   ```
   Test data draft generated:
   - X fixture files in .pipeline/test-data/fixtures/
   - Y scenario data sets in .pipeline/test-data/scenarios/
   - README.md with inventory
   
   Review the data. Adjust values, add missing cases, remove unnecessary ones.
   This is YOUR curated test data — the LLM proposed it, you own it.
   ```

## Rules
- Generated data must be REALISTIC. "Maria Rossi" not "User1". "Project Q3 Dashboard" not "Test Project".
- Generated data must be FICTIONAL. Never use data that looks like it could be a real person's information.
- Every fixture entry has a `notes` field explaining its purpose in the verification suite.
- IDs are human-readable slugs, not UUIDs or auto-incremented numbers. Tests should read like documentation.
- The human ALWAYS reviews and curates the output. This command proposes — it doesn't decide.
- If existing fixtures already cover the need, say so. Don't generate redundant data.
- YAML format for all fixtures (human-readable, easy to edit by hand).
