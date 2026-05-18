# Removing a cobra subcommand cleanly

## When this applies

You're removing a subcommand from a cobra command group (e.g.,
dropping `foo bar baz` where `bar` is the group). Your BDD / manual
test expects `foo bar baz` to return exit 1 with "unknown command".

## The trap

Deleting the subcommand's `AddCommand(...)` registration is
insufficient. Cobra's default for an unknown positional arg on a
parent command is to call the parent's `RunE(cmd, []string{"baz"})`.
If the parent's `RunE` is `return cmd.Help()` (the typical default),
you get exit 0 and a help dump — the removal "silently succeeds".

## The fix

Add `Args: cobra.NoArgs` to the parent command:

```go
cmd := &cobra.Command{
    Use:   "bar",
    Short: "...",
    Args:  cobra.NoArgs,
    RunE: func(cmd *cobra.Command, args []string) error {
        return cmd.Help()
    },
}
```

Cobra now validates args before calling RunE and returns an error
with message `unknown command "baz" for "foo bar"` + "Did you mean..."
when there's a close match. Exit code 1.

## What still works after this

- `foo bar --help` / `foo bar -h` — flags are routed before Args
  validation.
- `foo bar version`, `foo bar update` (legitimate subcommands) —
  subcommands are dispatched before the parent's RunE.
- `foo bar` (no args) — hits the parent's RunE and prints help.

## BDD assertion template

```gherkin
Scenario: removed subcommand fails loudly
  Given the CLI is installed
  When I run "foo bar <removed-subcommand>"
  Then the exit code is not 0
  And the output contains "unknown command"
```

## Related

- Place `Args: cobra.NoArgs` on the immediate parent of the removed
  subcommand, not on the root — you want legitimate siblings of the
  removed subcommand to still dispatch.
- For groups that take positional args by design (rare for nested
  groups, common for leaf commands), use `cobra.ExactArgs(n)` or a
  custom `Args` function instead.
