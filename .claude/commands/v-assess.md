# /v-assess — Analyze existing codebase and produce quality report

## Input
$ARGUMENTS (optional: specific area to focus on, e.g., "security", "backend only", "auth module")

## Procedure

1. Read `.pipeline/standards/` (all layers) to know the quality bar.
2. Read `.pipeline/standards/02-project.md` section "Intentional Deviations" — these are deliberate choices that should NOT be flagged as issues. If the project intentionally has no database, no ORM, no HTTP layer, etc., those are not findings.
3. Read `.pipeline/skills/framework/design-patterns.md` and `.pipeline/skills/framework/intuitive-design.md` for pattern and UX evaluation criteria.
4. If $ARGUMENTS specifies a focus area, limit the analysis to that area. Otherwise, analyze the full codebase.

## Analysis Phases

### Phase 1: Structure

- Directory organization: does it follow the stack conventions?
- Separation of logic and presentation: are they properly separated? Does logic import presentation concerns?
- Module boundaries: are responsibilities clear? Any god files (>500 lines doing unrelated things)?
- Dependencies between modules: is the dependency direction correct (presentation → logic, not the reverse)?

### Phase 2: Code Quality

- Function length: flag any over 20 lines (per standards)
- Cyclomatic complexity: flag any over 5 (per standards)
- Nesting depth: flag any over 2 levels (per standards)
- Code duplication: identify patterns that appear 3+ times
- Naming: are names descriptive and consistent?
- Dead code: unused functions, unreachable branches, commented-out blocks
- TODOs and FIXMEs: collect and list them

### Phase 3: Testing

- Test existence: which modules have tests? Which don't?
- Test quality: are tests meaningful (not just `assert True`)? Do they test behavior or implementation?
- BDD coverage: are there feature files? Do they match the user flows?
- Test data: fixtures present? Realistic or trivial?
- Mock usage: are mocks used appropriately (boundary only) or over-used (mocking everything)?

### Phase 4: Security

- Input validation: are all external inputs validated?
- SQL injection: parameterized queries everywhere?
- Authentication/authorization: checked on every protected route?
- Secrets: any hardcoded in source? Any in logs?
- Dependencies: known vulnerabilities? (`pip audit`, `go vet`, `bundle audit`, `npm audit`)

### Phase 5: Privacy (GDPR)

- Personal data in logs?
- Data minimization: are API responses leaking unnecessary fields?
- Deletion path: can a user be fully deleted?
- Consent tracking: if applicable, is it implemented?
- Third-party data flow: who gets personal data?

### Phase 6: Database & Migrations

- Migration system present? (Alembic, goose, Rails migrations)
- Migrations reversible?
- Schema matches code models?
- Indexes on frequently queried columns?

### Phase 7: Interface Intuitiveness (if UI/CLI exists)

Apply the 8 McKay attributes from `.pipeline/skills/framework/intuitive-design.md`:
- Discoverability, Affordance, Comprehensibility, Feedback
- Predictability, Efficiency, Forgiveness, Explorability
Focus on primary user flows. Flag critical failures.

### Phase 8: Design Patterns

- Which patterns are in use (explicitly or implicitly)?
- Are any patterns misapplied (anti-trigger violated)?
- Where would a pattern improve the code?

## Output

Produce `.pipeline/assessment.md`:

```markdown
# Codebase Assessment

## Summary
- Overall health: [good / fair / needs work / critical]
- Highest priority: [the single most important thing to fix]
- Estimated effort: [small: 1-3 days / medium: 1-2 weeks / large: 2+ weeks]

## Strengths
- [what the codebase does well]

## Critical Issues (fix immediately)
| Issue | Location | Risk | Suggested Fix |
|-------|----------|------|---------------|
| ... | ... | ... | ... |

## High Priority (fix soon)
| Issue | Location | Risk | Suggested Fix |
|-------|----------|------|---------------|
| ... | ... | ... | ... |

## Medium Priority (plan for)
| Issue | Location | Risk | Suggested Fix |
|-------|----------|------|---------------|
| ... | ... | ... | ... |

## Low Priority (nice to have)
| Issue | Location | Risk | Suggested Fix |
|-------|----------|------|---------------|
| ... | ... | ... | ... |

## Recommended Action Plan
1. [First intervention — what, why, estimated effort]
2. [Second intervention]
3. ...

Each intervention can become a `/v-spec` → `/v-plan` → ... pipeline cycle.

## Metrics
- Files analyzed: N
- Functions analyzed: N
- Test coverage: estimated X%
- Functions over complexity threshold: N
- Security findings: N
- Privacy findings: N
```

Report to human:
```
Assessment complete: .pipeline/assessment.md

Summary: [overall health]
Critical issues: N
High priority: N

Recommended first action: [one sentence]
Start with: `/v-spec [first intervention description]`
```
