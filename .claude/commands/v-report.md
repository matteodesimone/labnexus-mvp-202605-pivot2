# /v-report — Generate project status report

## Input
$ARGUMENTS (optional: "full", "summary", "client", "ideas", OR custom section name)

## Output Formats

Every report produces:
1. `.pipeline/reports/[date]-[type].md` — Markdown source
2. `.pipeline/reports/[date]-[type].pdf` — PDF rendered from Markdown (print-ready)

## Report Types

### Summary (default): `/v-report` or `/v-report summary`

One-page executive summary. Quick read for yourself or stakeholders.

### Full: `/v-report full`

Complete project status. Everything you need to know.

### Client: `/v-report client`

External-facing report for clients. Professional tone, no internal jargon, no pipeline mechanics. Only what the client cares about: what was delivered, what's next, blockers, timeline.

### Ideas: `/v-report ideas`

Overview of all captured ideas — backlog planning document.

## Procedure

1. Read all project state:
   - `.pipeline/state.json` — current pipeline state
   - `.pipeline/spec.md` — current spec (if exists)
   - `.pipeline/plan.md` — current plan (if exists)
   - `.pipeline/fette.md` — fette progress (if exists)
   - `.pipeline/review.md` — last review report (if exists)
   - `.pipeline/solutions/` — accumulated learnings
   - `.pipeline/ideas/` — captured ideas
   - `.pipeline/bugs/` — open bugs
   - `.pipeline/todos/` — open work items
   - `.pipeline/standards/02-project.md` — project identity, type, quality metric
   - `.pipeline/config.sh` — configuration
   - `git log --oneline -20` — recent commits

2. Read domain-specific state:
   - Check `.pipeline/checks.yaml` for domain check names
   - Read any domain-specific status sources (defined by domain)

3. Generate report based on type (see sections below).

4. Save as Markdown: `.pipeline/reports/[date]-[type].md`

5. **Convert Markdown to PDF** using the first available method. **Never install packages** — use only what's already on the system.

   Check in order:
   ```bash
   command -v pandoc          # Option 1: best quality
   python3 -c "import weasyprint"  # Option 2: good quality
   python3 -c "import xhtml2pdf"   # Option 3: lightweight
   ```

   **Option 1: pandoc** (if installed)
   ```bash
   pandoc report.md -o report.pdf -V geometry:margin=2.5cm -V fontsize=11pt
   ```

   **Option 2: Python weasyprint** (if importable)
   ```python
   import re, markdown, weasyprint
   raw = open('report.md').read()
   # SECURITY: strip raw HTML before Markdown conversion to prevent
   # file:// exfiltration and outbound requests via injected tags
   sanitized = re.sub(r'<[^>]+>', '', raw)
   html = markdown.markdown(sanitized, extensions=['tables', 'fenced_code'])
   styled = f'<html><head><style>body{{font-family:sans-serif;font-size:11pt;margin:2.5cm}}...</style></head><body>{html}</body></html>'
   # Disable external resource loading (no file://, no http://)
   weasyprint.HTML(string=styled, base_url='about:blank').write_pdf('report.pdf')
   ```

   **Option 3: Python xhtml2pdf** (if importable)
   ```python
   import re, markdown
   from xhtml2pdf import pisa
   raw = open('report.md').read()
   sanitized = re.sub(r'<[^>]+>', '', raw)  # strip raw HTML
   html = markdown.markdown(sanitized, extensions=['tables', 'fenced_code'])
   # same styled approach, pisa.CreatePDF(styled, dest=f, link_callback=lambda *a: None)
   ```

   **Option 4: nothing available**
   Save `.md` only. Warn clearly:
   ```
   ⚠ PDF not generated — no PDF tool found on this system.
   The Markdown report is ready: .pipeline/reports/[file].md

   To enable PDF reports, install ONE of:
     brew install pandoc          (macOS, recommended)
     sudo apt install pandoc      (Linux, recommended)
     pip install weasyprint       (in a separate venv, not your project's)
   ```
   Never suggest installing into the project's environment. Never run pip install.

6. Report: "Report generated: `.pipeline/reports/[date]-[type].md` and `.pdf`"

## Section Templates

### Summary Report

```markdown
# [Project Name] — Status Report
**Date**: [date]  **Pipeline**: [active/idle]  **Domain**: [domain]

## At a Glance
- **Current work**: [feature name or "idle"]
- **Progress**: [step] of pipeline, fetta [N]/[total] or "no fette"
- **Last review**: [PASS/FAIL with date, or "none"]
- **Open items**: [N] bugs, [M] todos, [K] ideas

## What's Done
[Completed fette/features as bullet list with dates]

## What's Next
[Current fetta scope or next planned work]

## Blockers
[Anything blocking progress, or "None"]
```

### Full Report

```markdown
# [Project Name] — Full Status Report
**Date**: [date]  **Version**: [from .pipeline/version]

## Project Overview
[From 02-project.md: name, type, quality metric, stack]

## Pipeline Status
[Current state, step, fetta, review iteration]

## Progress

### Completed
[Each completed fetta/feature with: what was built, review result, ship date]

### In Progress
[Current work: spec summary, plan summary, where in pipeline]

### Planned
[Remaining fette from fette.md, or upcoming specs]

## Quality
[Last review summary: findings count by severity, trends across iterations]
[Review config: mode, threshold, reviewers enabled]

## Learnings
[From .pipeline/solutions/ — key patterns discovered, grouped by theme]

## Ideas Backlog
[From .pipeline/ideas/ — table: priority, title, status, date]

## Open Items
[From .pipeline/bugs/ — open bugs with severity]
[From .pipeline/todos/ — open todos with priority]

## Recent Activity
[Last 20 git commits, formatted]

## Configuration
[Key config: project type, review mode, enabled reviewers, ship command]
```

### Client Report

```markdown
# [Project Name] — Progress Report
**Date**: [date]  **Prepared for**: [from 02-project.md client field]

## Summary
[2-3 sentences: what was accomplished since last report, what's next]

## Delivered
[Each completed feature in user-facing terms — no technical jargon]
[For each: what the user can now do that they couldn't before]

## In Progress
[Current work in user-facing terms]
[Expected completion: based on fette progress]

## Upcoming
[Next planned features in user-facing terms]

## Timeline
[Fetta progress bar or list: done/in-progress/planned]
[If prove-out: gettoni budget status]

## Decisions Needed
[Questions or approvals needed from the client]

## Ideas Under Consideration
[From ideas with status "planned" or priority ★★★ — in user-facing terms]
```

### Ideas Report

```markdown
# [Project Name] — Ideas Backlog
**Date**: [date]  **Total**: [count]

## By Priority

### ★★★ High Priority
[Table: title, date, status, one-line description]

### ★★ Medium Priority
[Table]

### ★ Low Priority
[Table]

### · Unprioritized
[Table]

## By Status
- **Planned** ([N]): ready to become specs
- **Considered** ([N]): reviewed, not yet committed
- **Captured** ([N]): raw ideas, not yet reviewed
- **Declined** ([N]): decided against (with reasons)

## Recently Added
[Last 5 ideas with full description]
```

## PDF Generation

Claude writes Markdown. External tools convert to PDF. Claude checks what's available and uses it.

**Rule: NEVER install packages.** Use only what's already on the system. Installing pip packages could break the user's project dependencies.

**Rule: SANITIZE before rendering.** Report content comes from pipeline artifacts (specs, bugs, ideas, commits) which may contain raw HTML. Before passing to any HTML-based renderer (weasyprint, xhtml2pdf), strip ALL raw HTML tags from the Markdown source. This prevents `<img src="file:///...">` exfiltration and `<link href="https://...">` outbound requests. Pandoc handles this internally — no extra step needed for Option 1.

**Fallback chain** (check only, no install):
1. `pandoc` — system tool, best quality
2. `weasyprint` — Python package, check with `import`
3. `xhtml2pdf` — Python package, check with `import`
4. Markdown only — with clear install instructions for the user

**Markdown best practices for good PDF output:**
- Use `---` horizontal rules where page breaks are desired
- Use proper heading hierarchy (# title, ## sections, ### subsections)
- Keep tables simple — avoid merged cells or very wide tables
- Use fenced code blocks with language tags for syntax highlighting

## Rules
- Reports are DESCRIPTIVE, not prescriptive. They state facts, not recommendations.
- Client reports never mention internal tools (Claude Code, pipeline, reviewers, fette).
- Client reports use the language of 02-project.md (quality metric, feature names as the client knows them).
- Ideas in client reports are filtered: only "planned" or high-priority. Don't overwhelm with raw captures.
- Reports are committed to git (they're project artifacts).
- The .pipeline/reports/ directory is created by `vibbly init` as part of the project scaffold.
