# /v-guide — Explain the framework, suggest next actions

## Input
$ARGUMENTS (optional — topic or question, e.g., "review", "fette", "how do I fix a bug", "what's a gate")

## Procedure

1. Read `.pipeline/state.json` to understand current state.
2. Read `.pipeline/config.sh` to know the domain and project type.
3. Read `.pipeline/standards/00-principles.md` for principles.
4. If `.pipeline/checks.yaml` exists, read it for domain-specific checks.
5. If `MANUAL-OPS.md` exists, be aware of manual alternatives.

## Behavior

**If no arguments**: give an overview of the framework tailored to the current project, then suggest the next action based on state.

```
Pipeline Framework Guide (domain: [domain name])

You're working with an agent pipeline that structures work into cycles:

  [brainstorm] → spec → plan → [verify] → [execute] → review → [ship] → compound

Each step has a command. I guide you through them. Gates can block
for your approval or let a multi-LLM tribunal decide.

Your project:
  Type: [project type]
  Domain: [domain]
  Fette: [N/total or "no fette defined"]
  Reviewers: [list enabled]
  Review mode: [iterations setting]

Current state: [current step or "no active pipeline"]
Next action: [what to do next]

Ask me about any topic: /v-guide review, /v-guide fette, /v-guide gates, /v-guide principles
```

**If arguments match a topic**: explain that topic in the context of the current domain.

### Topics (core — available in all domains):

- **"principles"**: list all principles (universal + domain-specific) with one-line explanations
- **"pipeline"**: explain the full pipeline loop with domain-specific command names
- **"fette" / "slicing"**: explain horizontal slicing, show current fette progress if available
- **"review"**: explain multi-reviewer system, dialectic, iteration, threshold
- **"gate" / "gates"**: explain human vs tribunal, list current gate settings
- **"compound"**: explain the learning loop, solutions, framework promotion
- **"status"**: same as `/v-status` but with more explanation of what each field means
- **"config"**: explain ESSENTIAL vs ADVANCED config, what to change and when
- **"reviewer" / "reviewers"**: list configured reviewers, their status, how to add/remove
- **"git"**: explain branch model, push strategy, when to merge
- **"budget" / "gettoni"**: explain budget system (prove-out only, say so if not applicable). Read `.pipeline/standards/project-types/prove-out.md` for details.
- **"prove-out"**: explain the Prove Out method — what it is, Discovery Goal, Tiny Experiments, checkpoints, the pact. Read `.pipeline/standards/project-types/prove-out.md`.
- **"ideas" / "backlog"**: explain /v-idea command, how to capture, review, search, and how brainstorm/spec surface related ideas
- **"grill" / "spirito critico"**: explain /v-grill command, how brainstorm integrates critical thinking, the decision tree approach, and when to use /v-grill vs /v-brainstorm vs /v-review
- **"onboard" / "adopt"**: explain /v-onboard command, what it does with archived AI files in .pipeline/onboarding/, when to run it, and how 02-project.md gets populated
- **"contexts"**: explain technology context files in `.pipeline/contexts/`, how they're fetched during init, how agents use them, and how to refresh with `vibbly context update`
- **"skills"**: explain two-level skill system (framework/ vs project/), how skills differ from standards and contexts, how /v-compound creates project skills, and how to promote to framework
- **"intuitiveness" / "intuitive" / "mckay"**: explain the "Intuitive by design" principle, the 8 McKay attributes, how /v-review evaluates them, and where to find the full checklist (`skills/framework/intuitive-design.md`)
- **"audit" / "perspectives"**: explain /v-audit command, how it differs from /v-review (rules vs people), available perspectives in `perspectives.yaml`, when to use it (after fette, before releases), and that findings are observations not blockers
- **"manual"**: explain that the framework works without vibbly CLI, reference MANUAL-OPS.md
- **"custom" / "extensions"**: explain how to add project-specific commands, review checks, and standards via the "Project Extensions" section in `02-project.md`
- **"report" / "reports"**: explain /v-report types (summary, full, client, ideas), output formats (Markdown + PDF), and how domain shapes report content
- **"principles N"** (e.g., "principles 3"): explain principle N in detail

### Topics (domain-specific — loaded from checks.yaml and methodology):

Read `.pipeline/checks.yaml`. Each check name is a valid topic:
- In coding domain: "quality", "security", "privacy", "intuitiveness", "test-coverage", "tdd", "bdd", "migrations", "patterns"
- Explain what the check covers, how the reviewer evaluates it, and what a finding looks like

### Conversational questions:

If the argument is a question rather than a topic keyword, answer it:
- "how do I fix a bug" → explain /v-bug vs review findings, one-pipeline-at-a-time rule
- "what do I do next" → read state, suggest next command
- "why did the review fail" → read .pipeline/review.md if exists, summarize findings
- "how do I add a reviewer" → explain config.sh REVIEWER_*_ENABLED

## Rules
- Always be concrete: reference the current project's state, config, and domain.
- Never be abstract: "you might want to..." → "run `/v-plan` to start fetta 4."
- Use domain-specific terminology: in coding say "test", "deploy", "implement". In other domains use their aliases.
- If you don't know the answer, say so and suggest where to look (which file, which section).
- Keep it concise. This is a quick reference, not a lecture. If the user wants depth, they'll ask follow-up questions.
