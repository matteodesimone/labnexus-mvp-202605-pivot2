# Invariant tests for cross-file conventions

## When this applies

You've established a rule that spans many files: "every bug file must
have frontmatter status", "every command file must have a project-local
mirror", "every migration has a down script", "every README section
includes a table of contents". Documentation alone won't enforce it —
the next person to add a file will forget.

## The pattern

Write a test (in the project's testing language) that:

1. Lists all files of the target type (glob or walk).
2. For each file, checks the invariant.
3. Fails loudly with a specific error message for each violation.

Go template:

```go
func TestEveryBugHasFrontmatterStatus(t *testing.T) {
    root, _ := findRepoRoot()
    entries, err := os.ReadDir(filepath.Join(root, ".pipeline", "bugs"))
    if err != nil { t.Fatal(err) }

    fmStatusRe := regexp.MustCompile(`(?m)^status:\s*\S`)

    for _, e := range entries {
        if e.IsDir() || filepath.Ext(e.Name()) != ".md" { continue }
        data, err := os.ReadFile(filepath.Join(root, ".pipeline", "bugs", e.Name()))
        if err != nil { t.Error(err); continue }

        // Extract frontmatter (first --- ... --- block)
        var frontmatter []byte
        if bytes.HasPrefix(data, []byte("---\n")) {
            end := bytes.Index(data[4:], []byte("\n---\n"))
            if end >= 0 { frontmatter = data[4 : 4+end] }
        }

        if !fmStatusRe.Match(frontmatter) {
            t.Errorf("bug %s: missing frontmatter status: field", e.Name())
        }
    }
}
```

## Properties

- **Self-documenting**: the test IS the spec. Anyone reading it sees
  the convention explicitly.
- **Regression-proof**: breaks the build the moment the convention is
  violated, not weeks later during a triage session.
- **Cheap**: 30-60 lines typically. No new infrastructure.
- **Error messages tell you what to fix**: include the path, the
  expected pattern, and ideally a suggested correction.

## When to write one

Write an invariant test when:

- You've made a decision about "all files of type X must Y".
- The rule isn't enforced by a parser, compiler, or type system.
- Violating the rule would silently bit-rot (not blow up at runtime).

Don't write one when:

- The rule is enforced by the build (e.g., Go imports, YAML syntax).
- There's only one file in question (just put the rule in that file).
- The rule is aesthetic ("commits should be well-written") — these
  need human review, not tests.

## Co-existence with other enforcement

- **Linters**: invariants complement linters. Linters catch per-file
  issues; invariants catch cross-file shape.
- **Pre-commit hooks**: invariants in the test suite run via
  `go test ./...` / `pytest` / etc. A pre-commit hook running the same
  suite is fine.
- **Documentation**: still write the convention down (in a skill, in a
  README section). The test enforces it; the doc explains why.

## Precedent

First applied in vibbly (2026-04-19):

- `features/steps/mirror_commands_test.go::TestProjectLocalCommandMirrorInSync`
  — every `framework/core/.claude/commands/v-*.md` has an identical
  mirror at `.claude/commands/v-*.md`.
- `features/steps/pipeline_metadata_lifecycle_test.go` — every bug and
  idea file has frontmatter-canonical state (no body duplication).

Both surfaced pre-existing drift on first run; both prevent regression now.
