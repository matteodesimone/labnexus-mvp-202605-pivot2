# /v-todo-export — Export todos to Markdown table or CSV

## Input
$ARGUMENTS: format (md | csv) and optional filters (--priority P1 --status ready)

Default: `md` with all todos.

## Procedure

1. Read all files in `.pipeline/todos/`.
2. Parse YAML frontmatter from each.
3. Apply filters if specified.
4. Sort by priority (P1 first), then by ID.

### Markdown output (default)

Write to `.pipeline/todos/export.md`:

```markdown
# Todos — [project name] — [date]

| ID | Priority | Status | Title | File | Created |
|----|----------|--------|-------|------|---------|
| 001 | P1 | ready | Fix SQL injection | auth.py:42 | 2025-01-15 |
| 002 | P2 | ready | Add missing test | service.py | 2025-01-15 |
```

### CSV output

Write to `.pipeline/todos/export.csv`:

```
id,priority,status,title,file,created,description
001,P1,ready,Fix SQL injection,auth.py:42,2025-01-15,"Parameterize query in login endpoint"
002,P2,ready,Add missing test,service.py,2025-01-15,"test_create_user_with_duplicate_email missing"
```

5. Report: "Exported N todos to .pipeline/todos/export.[md|csv]"

## Rules
- CSV uses RFC 4180 format (quoted fields with commas).
- Both formats are importable: Markdown renders in any viewer, CSV opens in any spreadsheet or task tool.
- If no todos exist, say so and suggest running `/v-triage` first.
