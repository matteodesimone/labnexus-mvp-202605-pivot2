# Gitignore exception for a tracked subdir under an ignored parent

Git re-inclusion (`!pattern`) only works against **patterns**, not against a
fully-ignored directory. If the parent dir is listed bare (`foo/`), you
**cannot** re-include anything inside it — git never even recurses to evaluate
the exception. The fix is to ignore the contents of the parent via glob
(`foo/*`) and then re-include the subdir.

## Trigger

You want one specific subdir tracked in git while everything else under the
same parent stays ignored. Typical case: an archive parent where most content
is ephemeral but one subtree is audit trail.

## Anti-trigger

- If the parent dir contains ONLY the thing you want tracked → just don't
  ignore it at all.
- If you want everything under the parent tracked → remove the ignore rule.
- If the "tracked subdir" is actually one file → use a file-level exception
  (`!foo/specific-file`) and keep the parent-dir ignore as-is.

## Rule shape

```gitignore
# WRONG — re-include is silently ignored:
parent/
!parent/tracked-subdir/

# RIGHT — parent ignored via glob, subdir re-included:
parent/*
!parent/tracked-subdir/
!parent/tracked-subdir/**
```

The `/*` on the parent is load-bearing. Without it, git treats the parent as
fully excluded and never evaluates the `!` line.

Both `tracked-subdir/` and `tracked-subdir/**` are needed:

- `!subdir/` re-includes the directory itself (so `git add subdir/` works).
- `!subdir/**` re-includes all nested content (files created later).

Without `**`, a file created later inside the re-included dir is still ignored.

## Verification (always do this)

After adding the rules, commit a sentinel file (e.g. `.gitkeep`) in the target
subdir and run `git check-ignore -v path/to/file`. If the file is reported as
ignored, the rule shape is wrong.

```bash
touch parent/tracked-subdir/.gitkeep
git check-ignore -v parent/tracked-subdir/.gitkeep
# → empty output = tracked. Any output naming the rule = still ignored.
git status --short parent/tracked-subdir/.gitkeep
# → should show `??` (untracked, ready to add), NOT nothing (ignored).
```

## Gotcha

- Rules are evaluated top-to-bottom with later rules winning on ties. Exceptions
  should sit below the broad rule they're negating.
- `.gitkeep` is convention, not special. Any committable file works. Empty dirs
  cannot be tracked — git tracks files, not directories.
- Pre-existing ignored files won't magically appear as untracked after changing
  rules. Use `git add -f` once or delete/recreate.
- Framework gitignore templates (e.g. `gitignore-base` distributed via
  `vibbly init`) must mirror the same rule shape as the live `.gitignore`,
  otherwise new projects won't inherit the exception.

## Precedent

Applied in vibbly (2026-04-19) to track `.pipeline/archive/applied-proposals/`
under an otherwise-ignored `.pipeline/archive/` parent. Initial attempt used
`.pipeline/archive/` (bare) and the sentinel stayed ignored; fixed by switching
to `.pipeline/archive/*`.
