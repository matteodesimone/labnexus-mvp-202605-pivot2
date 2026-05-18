# Adopt-merge pattern for extending user-authored tooling files

When framework adoption extends a user-authored tooling file (not overwrite),
use format-aware merge with an idempotency sentinel. Two variants exist;
pick per content type.

## Variant A — Whole-file archive (for opaque/narrative files)

For files the framework can't parse safely (CLAUDE.md, .cursorrules,
copilot-instructions.md, AGENTS.md): move the user's version to a
staging area BEFORE overwrite, then let a later pipeline step
(/v-onboard) distill its content into the framework layers.

Precedent: `archiveAIConfig` in vibbly's `internal/core/project/onboarding.go`.

Rules:

- Archive target: `.pipeline/onboarding/` (or an equivalent staging folder
  in other frameworks).
- Idempotent: if `.pipeline/onboarding/<file>` already exists from a prior
  adopt, do NOT overwrite (preserve the FIRST snapshot — it's the one
  closest to the human's original intent).
- Generate an `ADOPT-REPORT.md` alongside listing the archived files and
  pointing at the distillation command (`/v-onboard` for vibbly).

## Variant B — In-place merge (for structured content)

For files with identifiable key/value or line semantics (.mise.toml,
.gitignore, .editorconfig, .env.example, Dockerfile ENV blocks):
merge the framework template INTO the user's file in place, preserving
the user's content verbatim and appending only template-only additions.

Precedent: `mergeMiseToml` + `mergeGitignore` in vibbly's
`internal/core/project/project.go`.

### Rules

1. **Parse both into a structured representation** matching the format:
   - TOML-ish: section (`[name]`) + lines within. Identify keys via the
     leftmost `=`.
   - Line-list (gitignore-shaped): one pattern per line, blank lines
     and comments preserved as context.
   - Block-structured (.editorconfig, ini, etc.): header + lines, same
     shape as TOML.

2. **User wins on conflicts.** If a key exists in both the user's file
   and the template with different values, keep the user's. Their version
   is deliberate.

3. **Template fills gaps.** Keys in the template but not in the user's
   file are appended within the same section (TOML) or after a delimiter
   (gitignore).

4. **Embed a sentinel.** When appending framework-owned content, prefix
   it with a recognizable comment like:

   ```
   # Appended by <framework> init --adopt
   ```

   On re-merge, detect the sentinel to SKIP re-emitting the header. Without
   this, repeated adopts accumulate duplicate delimiter headers.

5. **Normalize whitespace per section.** Before emitting each section,
   trim trailing blank lines. Otherwise re-merge accumulates one extra
   blank per pass (verified bug observed during vibbly dogfooding 2026-04-20).

6. **Test the fixed point explicitly.** Add a test of the form:

   ```go
   merged1 := Merge(existing, template)
   merged2 := Merge(merged1, template)
   require.Equal(t, merged1, merged2)
   ```

   If the two aren't byte-identical, the merge isn't idempotent — fix the
   whitespace/comment accretion before shipping.

7. **Lenient parse, `([]byte, error)` preferred.** Malformed user content
   should fall back gracefully — either skip merge (template wins) with a
   log line, or keep user content verbatim + append template as-is. Never
   crash or truncate.

## When to pick which variant

| User file type | Variant |
|---|---|
| Narrative / docs / LLM-instructions | A (archive + distill) |
| Structured / config / line-lists | B (merge in place) |
| Mixed (e.g. README with config block) | A — too hard to merge safely |

If in doubt, default to A — archiving is always reversible, merging
isn't if the merge logic has a bug that overwrites user content.

## Verification checklist

Before shipping any new adopt-merge path:

- [ ] Reproduction test: pre-existing file with a user-specific entry
      the template doesn't have → after adopt, that entry MUST survive.
- [ ] Integration test through the full adopt flow (not just unit test
      of the merger).
- [ ] Idempotency test (re-adopt → same result).
- [ ] Dogfood against at least one REAL pre-existing project, not just
      synthetic tempdirs. Synthetic data has too-clean edges.
