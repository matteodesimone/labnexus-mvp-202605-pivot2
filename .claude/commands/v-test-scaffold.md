# /v-test-scaffold — Create tests before implementation (TDD + BDD Red Phase)

## Procedure

1. Read `.pipeline/state.json`.
   - If plan is not in `steps_completed`, REFUSE: "Cannot write tests without an approved plan. Run `/v-plan` first."

2. Read `.pipeline/plan.md` (Test Strategy section), `.pipeline/spec.md`, and approved feature files.

3. If `.pipeline/fette.md` exists, read `current_fetta` from state. Scaffold tests ONLY for the current fetta — do not create tests for features in future fette. The plan covers only the current fetta, so the test strategy in the plan already scopes this, but verify: if a test references a feature not in the current fetta's FR list, skip it.

4. **BDD acceptance tests (first priority).** 
   
   a. Copy approved feature files from `.pipeline/proposed-features/` to the project's `features/` directory (or the stack-appropriate location).
   b. Write step implementations for each scenario. Steps must:
      - Launch the REAL application (CLI subprocess, HTTP test server, or browser via Playwright)
      - Use test data from `.pipeline/test-data/fixtures/` — load the approved fixtures, never generate data inline
      - Assert on OBSERVABLE outcomes (stdout, HTTP response, database state, UI state)
      - NO mocks in acceptance tests. These test the real system end-to-end.
   c. **For web stacks (webapp-rails, webapp-svelte, etc.): Playwright is mandatory.** Read `.pipeline/standards/stacks/webapp.md` "Web Application Testing Strategy" section. Every acceptance test opens a real browser, loads test data, and verifies the user flow. Specifically:
      - Set up Playwright with the stack's BDD integration (`playwright-bdd` for Svelte, `capybara-playwright-driver` for Rails)
      - Write a test data seed/setup step that loads `.pipeline/test-data/fixtures/` before each test suite
      - Every Given step that describes data ("Given 3 appointments exist") loads from fixtures, not from inline creation
      - Every Then step asserts on what the USER SEES in the browser (text on screen, element visibility, navigation result) — not on database state
   d. Create empty stubs so steps can import but fail on assertions (red phase).

5. **Integration tests.**
   
   a. For each service/module interaction identified in the plan, write integration tests.
   b. Use real database (sqlite in-memory), real HTTP client (test client), real filesystem (temp dir).
   c. Load test data from fixtures.
   d. NO mocks for internal dependencies. Mock ONLY external API boundaries (zone 2 from methodology).

6. **Unit tests.**
   
   a. For each pure logic function in the plan, write unit tests.
   b. Table-driven tests where applicable (especially for Go).
   c. Cover: happy path, error paths, boundary values, edge cases from spec.

7. **Test data setup.**
   
   a. If `.pipeline/test-data/fixtures/` has approved fixtures, write the test helpers that load them (conftest.py for Python, test helpers for Go, seeds for Rails).
   b. If fixtures don't exist yet, warn: "No test data fixtures found. Run `/v-test-data` to create them, or create `.pipeline/test-data/fixtures/` manually."

8. **CRITICAL VERIFICATION**: Run ALL test suites.
   - ALL tests MUST FAIL. This is the expected state (red phase).
   - BDD steps fail because the application doesn't implement the behavior yet.
   - Integration tests fail because modules/endpoints don't exist yet.
   - Unit tests fail because functions return NotImplementedError or similar.
   - If any test passes, something is wrong.
   - Create minimal empty stubs so tests can import but still fail on assertions.

9. Report:
   ```
   Test scaffold complete.
   
   BDD acceptance tests:
   - X feature files (Y scenarios)
   - Step implementations in [path]
   - Browser driver: Playwright [or: N/A for CLI]
   
   Integration tests:
   - X test files (Y test cases)
   
   Unit tests:
   - X test files (Y test cases)
   
   Test data:
   - Fixtures loaded from .pipeline/test-data/fixtures/ [or: no fixtures yet]
   
   All tests FAILING as expected (red phase).
   ```

10. Update `.pipeline/state.json`: mark plan completed, set `current_step` to "test-scaffold", record all test file paths in artifacts.

11. **Gate check** — read `GATE_TEST_SCAFFOLD` from `.pipeline/config.sh`:
    - If `"human"`: set `blocked: true`, reason: "Test scaffold produced — awaiting human review of test coverage and BDD scenarios."
    - If `"tribunal"`: run tribunal with concatenated test files + feature files. Cross-reference. Decide.

12. Tell the human:
    ```
    Test scaffold ready. [N] tests created (all failing — red phase).
    - BDD acceptance: [N] scenarios
    - Integration: [N] tests
    - Unit: [N] tests

    Run `/v-implement` to make them pass.
    ```

## Rules
- NEVER write implementation logic in this step. Only tests and empty stubs.
- BDD scenarios must match the approved feature files exactly. Don't invent new scenarios — the spec gate already approved them.
- If the plan's test strategy is insufficient (missing edge cases, missing error scenarios), ADD tests beyond what the plan specified and note what you added.
- Follow the mock zones strictly: zone 1 (never mock) for DB/filesystem/internal, zone 2 (mock at boundary) for external APIs, zone 3 (real app) for acceptance.
- **For web stacks: Playwright is non-negotiable for acceptance tests.** If the stack is a webapp, every BDD scenario runs in a real browser. Acceptance tests without a browser are integration tests mislabeled.
- **Test data is loaded, not invented.** Acceptance and integration tests use fixtures from `.pipeline/test-data/fixtures/`. If a test creates data inline (`User.create!(name: "test")`), it's wrong — use the curated fixtures.
- Use the project's existing test framework and conventions. Match the style of existing tests.

## Step phrasing — avoid collisions

BDD step handlers are matched by regex. Two collision classes matter when scaffolding:

- **Class A (ambiguity)**: two registered regexes both match the same Gherkin. The BDD runner should be in strict mode (e.g. `godog.Options{Strict: true}` for Go) so Class A fails the suite loudly with "ambiguous step definition" naming both handlers. If Strict's side-effects (failing on undefined/pending steps too) are too broad for the current state of the suite, implement a narrower custom check that scans registered patterns vs feature-step texts and fails only on multi-match — see `features/steps/ambiguity_check_test.go` for the precedent.
- **Class B (wrong-handler)**: one regex matches, but it's from a different file than the one you're authoring. No tool catches this — only phrasing discipline does.

Rule: new step phrasings MUST start with a domain-specific noun-verb pair (`the XDG inbox has a proposal "X"`, `the migration file "M" declares column "C"`) rather than a generic one (`a file "X" containing "Y"`). Generic phrasings belong in shared-helper step files (e.g., `common_test.go`), not in feature-specific step files.

Audit before registering a new step in a feature-specific file:

```bash
grep -rn "^\s*ctx\.Step" features/steps/
```

Eyeball the regex list. If your new Gherkin line would also match an existing regex from a different file, rename your Gherkin or your regex so the match is unique and semantically correct.
