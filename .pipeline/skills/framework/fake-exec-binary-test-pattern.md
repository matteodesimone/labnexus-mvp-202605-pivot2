# Fake exec binary via t.TempDir() + t.Setenv(PATH)

## When to use

Your production code calls `exec.Command("<tool>", ...)` and you want to
unit-test the wrapper logic (argument building, output parsing, result
classification) without depending on the real tool being installed in
the test environment.

Examples from the field: `mise ls` parser (vibbly's `MiseToolsCheck`,
shipped 2026-04-21), `git` invocations, LLM CLI wrappers (`codex`,
`ollama`, `vibe`), `goreleaser` hooks, migration tools.

## The pattern (5 steps)

1. In the test, create an isolated temp dir:

   ```go
   dir := t.TempDir()
   ```

2. Write a shell script to that dir, named exactly like the tool you
   want to stand in for. Make it executable (mode `0o755`):

   ```go
   script := `#!/bin/sh
if [ "$1" = "ls" ]; then
  printf '%s' 'output-you-want-the-wrapper-to-see'
  exit 0
fi
exit 0
`
   os.WriteFile(filepath.Join(dir, "mise"), []byte(script), 0o755)
   ```

3. Prepend the temp dir to PATH for the duration of the test:

   ```go
   t.Setenv("PATH", dir + string(os.PathListSeparator) + os.Getenv("PATH"))
   ```

   Go 1.17+ auto-restores PATH at the end of the test, including on
   failure. No manual cleanup needed.

4. Call the production function. `exec.LookPath("mise")` inside the
   wrapper resolves to your fake first (because it's earlier in PATH),
   so the wrapper runs against controllable output.

5. Assert on the function's return value as if it had called real mise.

## Why it works without mocking

- `exec.Command(name, args...)` uses argv semantics — no shell parsing.
  The fake receives the subcommand as `$1` literally, not as a string
  that a shell might reinterpret. Zero command-injection surface in
  the test setup.
- Go's test isolation (`t.TempDir()`, `t.Setenv()`) handles cleanup
  automatically. No leaked state between tests.
- No production-code changes are required to accommodate the test.
  The wrapper stays vanilla `exec.Command(...)` — which is what you
  want to test anyway.

## Reusable helper

Extract the fake-binary-write step into a helper in the test file
when you have 2+ tests using the same fake:

```go
func writeFakeMiseBinary(t *testing.T, lsOutput string) {
    t.Helper()
    dir := t.TempDir()
    // Escape for safe single-quoted interpolation in the shell script:
    escaped := strings.ReplaceAll(lsOutput, "'", `'\''`)
    script := "#!/bin/sh\nif [ \"$1\" = \"ls\" ]; then\n  printf '%s' '" + escaped + "'\n  exit 0\nfi\nexit 0\n"
    path := filepath.Join(dir, "mise")
    if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
        t.Fatalf("write fake mise: %v", err)
    }
    t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}
```

Then tests become one-liners:

```go
func TestFoo_MissingPath(t *testing.T) {
    writeFakeMiseBinary(t, "ruff 0.15.11 (missing) ~/.mise.toml latest")
    result := MiseToolsCheck(t.TempDir())
    if result.Passed { t.Error("expected failure on missing tool") }
}
```

## Gotchas

- **Platform**: the shell script requires `/bin/sh` — fine on macOS and
  Linux, not on Windows. For cross-platform tests, compile a tiny Go
  binary into the temp dir instead (`go build -o <tmpdir>/<tool>` of a
  test-only `main.go`). More setup, fully portable.
- **Exit code matters**: if your wrapper branches on exit status, make
  the shell script exit with the right code for each scenario. `exit 1`
  for failure paths, `exit 0` for success.
- **Don't shell-interpolate test input**: if the fake's output contains
  single quotes or shell metachars, escape them (see the helper above)
  or write the output to a file the script `cat`s. Shell string
  interpolation in test setup is the one place the argv-safety argument
  doesn't protect you.
- **Keep the fake minimal**: one subcommand, one output. If you need to
  test multiple subcommands, either write separate fakes per test or
  make the fake branch on `$1` / `$2`. Resist turning the fake into a
  mini-reimplementation of the real tool.

## When NOT to use

- When you're writing an acceptance test that's supposed to exercise
  the real tool on a real project (e.g., `vibbly doctor` E2E). Those
  live in `features/` (BDD) and use the actual tool via the project
  `.mise.toml`.
- When the production code's interaction with the tool is so thin
  (one call, no output parsing, just exit-code check) that a real-tool
  integration test is cheaper than the fake-binary setup.
- When the tool's output format is genuinely unstable across versions —
  the fake captures a snapshot, but won't detect if a major version
  bump changes the format. Pair with a version-pin in `.mise.toml` (or
  equivalent) so the real and fake see the same format.

## Precedent in the framework

First landed in vibbly `internal/core/doctor/doctor_test.go` on
2026-04-21 during the `doctor-reports-passed-despite-missing-tool`
bugfix. See solution `.pipeline/solutions/2026-04-21-doctor-reports-
passed-despite-missing-tool.md` for the original context.
