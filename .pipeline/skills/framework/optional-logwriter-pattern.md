# Optional io.Writer on Options struct for best-effort info

## When this applies

You have a core function (no cobra/cli dependency) that needs to tell
the user something non-fatal: progress, counts, warnings that don't
abort the operation. The calling CLI layer wants this to land on
cobra's stdout, but the core package must stay framework-agnostic.

## The pattern

Add an optional `io.Writer` field to your Options struct. Nil means
silent. The core package uses it only if non-nil. The CLI passes
`cmd.OutOrStdout()` (or `cmd.ErrOrStderr()` for warnings).

```go
// core package
type DoThingOpts struct {
    Target    string
    DryRun    bool
    LogWriter io.Writer // optional; nil = silent
}

func DoThing(opts DoThingOpts) error {
    // ... do work ...
    if opts.LogWriter != nil {
        fmt.Fprintf(opts.LogWriter, "✓ processed %d items\n", count)
    }
    return nil
}

// CLI package
RunE: func(cmd *cobra.Command, args []string) error {
    return core.DoThing(core.DoThingOpts{
        Target:    target,
        DryRun:    dryRun,
        LogWriter: cmd.OutOrStdout(),
    })
}

// Test (silent)
core.DoThing(core.DoThingOpts{Target: "x"}) // LogWriter nil

// Test (assert on output)
var buf bytes.Buffer
core.DoThing(core.DoThingOpts{Target: "x", LogWriter: &buf})
require.Contains(t, buf.String(), "processed 3")
```

## Why not alternatives

- **fmt.Fprintln(os.Stdout)** in core: couples core to os.Stdout,
  tests become order-sensitive, core can't be used by a library
  consumer that wants to redirect output.
- **Return InitResult with Warnings []string**: forces every caller
  to handle (or ignore) a new return value, cascades through the
  stack, no streaming — warnings only appear after the call completes.
- **Callback `func(string)`**: verbose at call sites; nil-safe check
  is the same effort as io.Writer but less idiomatic in Go.
- **Structured logger (zap, slog)**: over-engineered for CLI progress;
  introduces a new dependency; output format is constrained by the
  logger.

io.Writer is the right abstraction: same interface `os.Stdout` satisfies,
same interface `bytes.Buffer` satisfies for tests, same interface
`io.Discard` satisfies for silent mode.

## When this does NOT apply

- Errors that abort the operation → use return value, not LogWriter.
- High-volume structured logs (tracing, metrics) → use a real logger.
- Output that's machine-consumable (JSON, table rows) → return it as
  a value, let the CLI serialize.

## Naming convention

Call the field `LogWriter` or `Output` — both are common Go idioms
(`io.Writer`, `http.ResponseWriter`). Avoid `Logger` unless the type
is actually a logger interface.
