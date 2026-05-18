# Git Strategy (Core)

Git is the safety net that makes everything else possible. Without reliable version control, iterative development is fragile, agent-driven work is dangerous, and rollback is prayer. These rules are non-negotiable and apply to ALL domains.

## Branch model

Two permanent branches, feature branches as needed:

**`main`** — production/published. Always in a shippable state. Only the human pushes to main, never an agent. Merges come from `development` after review and approval. This is what gets shipped/deployed/published.

**`development`** — the working branch. This is where the pipeline operates. Agents work here. Reviews happen here. When a pipeline cycle completes (ship step passed), the human merges development into main.

**Feature branches** — created from `development` when isolation is needed. Naming: `feature/<short-slug>`. Used in two scenarios:

1. **Worktree mode**: when `WORKTREE_ENABLED=true`, each pipeline cycle creates a feature branch in a separate worktree. Merges back to `development` after review.
2. **Manual isolation**: when exploring something risky. If it goes wrong, `git branch -D` and you're clean.

Bugfix branches: `fix/<short-slug>`. Same flow as feature branches.

## Commit rules

**Atomic commits.** Each commit is a coherent unit: one deliverable + its verification, one refactor step, one config change. Not "WIP" or "various fixes".

**Commit messages.** Conventional format: `type: short description`. Types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`.

**Pipeline state is not committed.** `.pipeline/state.json`, `.pipeline/archive/`, `.pipeline/reviews/`, and export files are in `.gitignore`. They're ephemeral. Spec, plan, standards, skills, solutions, features, and test data ARE committed — they're project artifacts.

## Push strategy

**After `/v-ship` completes successfully**: the agent pushes development to remote. The work is tested, reviewed, and shipped — it's safe to push. This protects against local disk failure.

```
git push origin development
```

**After `/v-compound`**: push again (compound adds solutions, may update standards).

**Never push main.** The human merges development → main and pushes main. This is a human-only action.

**If push fails** (no remote, auth error, network): warn but don't block. The pipeline continues. The push is a safety measure, not a gate.

## Rollback as safety net

Git is your undo button. The agent can't break anything permanently because:

- Every change is on a branch, never directly on main.
- `git reset --hard HEAD~N` undoes the last N commits.
- `git stash` saves work-in-progress without committing.
- Feature branches can be deleted without affecting development.

**Rule for agents:** if something goes seriously wrong during `/v-execute` (work that was passing now fails, unexpected behavior), the agent should `git stash` the current work, report the issue, and wait for human decision. Never push forward through a mess.

## Protection rules

- **Never force push** to `main` or `development`.
- **Never rebase** `development` or `main`. Merge only.
- **Feature branches can be rebased** before merging (for clean history), but only if they haven't been shared.
- **Tags for releases.** When main is shipped, tag it: `v1.0.0`, `v1.0.1`. Semantic versioning.
