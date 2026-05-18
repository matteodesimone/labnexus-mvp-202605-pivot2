# Prove Out — Project Type Definition

**"Mettiamo l'idea alla prova."**

Prove Out is a framework for building MVPs as business experiments. When `PROJECT_TYPE` is `prove-out`, the pipeline integrates with this methodology.

This file defines the method. Project-specific data (the actual Discovery Goal, budget numbers, etc.) lives in `02-project.md`.

---

## What It Is

The client buys a business experiment, not finished software. The MVP is the instrument to collect data that validates a hypothesis. If the experiment says "no, the idea doesn't work," that's a success of the method — the client spent little to find out.

## The Project Pact

These rules are shared with the client at the start. They define how the collaboration works.

- **This is an MVP, not a product.** The deliverable is an experiment instrument, not finished software.
- **Software is R&D.** We learn by doing. Discoveries change decisions. We work in short iterations.
- **We are partners, not client and vendor.** The developer operates as CTO (technical decisions). The client is the product manager (business decisions, experiment design, user acquisition).
- **Time and cost are fixed.** The package doesn't change.
- **The client decides WHAT, the CTO decides HOW.** Which information and data are needed, and what we build to get them — the client decides (as product manager). How we build it, how we procure and deliver the data — the developer decides (as CTO). At every change, they sit together and re-discuss priorities: the client decides what matters, the developer decides how to do it.
- **Transparency is mutual.** The developer shows work status, gettoni spent and remaining, technical choices. The client shares new information — market feedback, strategy changes, user conversations. Withheld information wastes gettoni.

### Conformity criterion vs success criterion

The **conformity criterion** is: the software makes the necessary data available to answer the Discovery Goal and the agreed sub-experiments. This is what the developer delivers and what gets evaluated at sprint closure.

The **success criterion** ("Avrò ragione se...") is about the business — did the experiment validate the hypothesis? That's the client's responsibility. The developer delivers the instrument; the client runs the experiment.

This distinction matters: the developer is accountable for the software conformity, not for the business outcome.

### Pact implications for the pipeline

- The developer has CTO prerogative on technical decisions. `/v-plan` does not ask the client HOW — only WHAT to prioritize.
- At every change, the developer presents the trade-off in gettoni and the client decides priorities.
- Bug fixes are budget-free by default. If a bug requires work beyond the original scope, evaluate at next checkpoint.

---

## Discovery Goal

The business experiment formalized. Compiled with the client during the workshop.

### Structure

| Field | Description | Required |
|---|---|---|
| **Credo che...** | The business hypothesis | Yes |
| **Avrò ragione se...** | Success criteria (the metric) | Yes |
| **C'è un problema nella mia ipotesi se...** | Warning signals | No |
| **Non farò...** | Explicit boundaries | No |
| **Starò attento a...** | Risks to monitor | No |
| **Lo faremo perché...** | Motivation | No |
| **Lo faremo con...** | Resources committed | No |
| **Per rispondere, esplorerò queste domande** | Sub-experiments, ordered by decreasing uncertainty | Yes |
| **Risponderò entro il...** | Sprint duration | Yes |

The Discovery Goal drives the project's quality metric. "Is the software good?" is answered by the conformity criterion: "does the software make the necessary data available to answer the Discovery Goal and the agreed sub-experiments?" This is distinct from the success criterion ("Avrò ragione se"), which is the client's business outcome.

---

## The Workshop (2-3 hours)

The workshop produces shared knowledge. The client is the product manager, the developer is the CTO. Together they produce: the experiment to validate, a shared image of how the product works, and the package choice.

### Phase 1 — Discovery Goal

Compile the Discovery Goal together. The key question to surface the metric: "in a month, if it went well, what concretely happened?"

Sub-experiments are intermediate questions to validate the hypothesis. Ordered by decreasing uncertainty: the riskiest first.

For each sub-experiment, ask: "does this question need software to be answered, or is it a client action?" Sub-experiments that need software become the input for the map. Those that are client actions (acquiring users, marketing, collecting feedback in person) remain the product manager's responsibility.

### Phase 2 — Essential Map

The map is derived from the sub-experiments that need software. For each of those sub-experiments, build the user journey in sequence: what the user does, step by step, from start to finish.

The map doesn't emerge from scratch — it's born from the questions the software must answer. This guarantees every step built serves the experiment. The correlation isn't a constraint to verify afterward — it's how the map is constructed.

**Bridges** emerge naturally during construction: "to reach this question, what does the user need to do first?" A Bridge doesn't answer a sub-experiment but is required infrastructure to get there (login, navigation, mandatory steps).

If the client proposes steps that don't serve any sub-experiment and aren't required infrastructure: "this has value, we keep it for later as a growth plan."

### Phase 3 — Exploration

On the important steps of the map — those that answer sub-experiments — explore with breadboard and rough sketches. Places, actions, connections.

Exploration surfaces details and with details, what NOT to do. For each extra feature: "do we need it for the experiment or keep it for later?"

Details that emerge are documented as **Tiny Experiments** — the form used to document project tasks. A Tiny Experiment can be a question ("can we do X?") or a direct task, depending on what's needed. Each has a truth criterion ("true if") which in the implementation domain is the definition of done.

### Workshop output

- Discovery Goal (compiled)
- Essential map (steps that matter)
- Breadboard and sketches (from exploration)
- List of things we won't do
- Growth plan (things with value that don't fit the MVP)
- Package chosen
- Photos of everything

---

## Workshop Output to Pipeline (Input Bridge)

Workshop materials translate into pipeline artifacts:

| Workshop output | Pipeline artifact | Pipeline step |
|---|---|---|
| Discovery Goal | Conformity criterion in `02-project.md` | Project setup |
| Discovery Goal: "Non farò" | Out of scope in spec | `/v-spec` |
| Discovery Goal: "Starò attento a" | Risk assessment in plan | `/v-plan` |
| Discovery Goal: sub-experiments (software-only) | Fette ordering (riskiest first) | `/v-plan` |
| Sub-experiment software filter | Scope boundary (client actions excluded) | `/v-spec` |
| Essential map steps (derived from sub-experiments) | Functional requirements | `/v-spec` |
| Breadboard/sketches | Architecture starting point | `/v-plan` |
| Tiny Experiments (from exploration) | Feature definitions with truth criteria | `/v-spec` |
| Bridges (from map construction) | Requirement priority: minimal (infrastructure) | `/v-spec` |
| List of things we won't do | Out of scope | `/v-spec` |
| Growth plan | Ideas backlog or Reserve priority | `/v-idea` |
| Step classification: Pillar | Requirement priority: must | `/v-spec` |
| Step classification: Reserve | Out of scope for this fetta | `/v-spec` |

---

## Gettoni and Packages

### What a gettone is

A relative unit of work complexity. Not an hour, not a todo — a relative dimension, comparable with other features in the same project. It makes trade-offs visible when the scenario changes.

The client sees total gettoni and remaining gettoni. The client does NOT allocate or estimate gettoni — that's CTO work. When priorities need to change, the developer presents the cost in gettoni and the client decides.

### Packages

| | **S** | **M** |
|---|---|---|
| Duration | 3 weeks | 5 weeks |
| Total gettoni | 20 | 35 |

About 25% of gettoni go to infrastructure, security, and technical setup — the work under the hood. The client doesn't see this split.

Internally, the sprint includes a 1-week buffer to absorb estimation errors. The client doesn't know this.

### Budget tracking

The pipeline tracks gettoni in `state.json`:

```json
"budget": {
  "package": "S",
  "total": 20,
  "infra": 5,
  "client": 15,
  "spent": 8,
  "spent_per_fetta": {
    "1": { "total": 5, "infra": 2, "client": 3 },
    "2": { "total": 3, "infra": 0, "client": 3 }
  },
  "allocated_per_fetta": {
    "1": 6,
    "2": 8,
    "3": 6
  }
}
```

The `/v-status` command shows: total, spent, remaining, per-fetta breakdown (infra vs client).

At checkpoints, the developer tells the client only: "N gettoni remaining out of T" — no infra/client split.

### Configuration fields

```
PACKAGE_SIZE          # S | M
BUDGET_TOTAL          # 20 (S) or 35 (M)
BUDGET_INFRA          # ~25% of total
BUDGET_CLIENT         # remainder
SPRINT_DURATION_WEEKS # 3 (S) or 5 (M)
```

---

## Tiny Experiments

Each feature is documented as a Tiny Experiment — a question ("can we do X?") or a direct task, with a truth criterion ("true if") and bounded tasks.

### Structure

| Field | Description |
|---|---|
| **Domanda** | The question this feature answers |
| **Vero se** | Definition of done — when is the answer "yes"? |
| **Sotto-esperimento** | Which sub-experiment from the Discovery Goal this serves |
| **Gettoni** | Estimated cost (CTO allocation) |
| **Fetta** | Which horizontal slice this belongs to |
| **Breadboard** | Flow: places, actions, connections |
| **Tasks** | Actions to complete (checkable list) |
| **Non farò** | Explicit boundaries for THIS feature |
| **Starò attento a** | Risks for THIS feature (rabbit holes) |

A Tiny Experiment is complete when its Tasks are done. "Non farò" protects against scope creep at the feature level.

---

## Horizontal Slices (Fette)

> **Terminology note**: the Prove Out method calls these "fette verticali" (vertical slices through the user story map, which runs horizontally). The pipeline framework calls the same concept "horizontal slices" (thin cuts across the full user journey). They mean the same thing: each fetta crosses the ENTIRE user journey. The pipeline uses "horizontal" to align with the principle "sempre funzionante" in `00-principles.md`.

Each fetta crosses the entire user journey and produces something working.

- **First fetta** (the scooter): the entire path in the most reduced way possible
- **Subsequent fette** (bicycle, motorbike, car): each adds richness

After each fetta the client sees something usable, not an isolated piece.

Fette are ordered by decreasing uncertainty (from sub-experiments). The riskiest hypothesis is tested first.

---

## Checkpoints

Every 3-4 days the developer shows the working product. Meetings are scheduled in advance.

### Protocol

1. **Review**: the client sees and tries the product. The developer shows remaining gettoni.
2. **Revalidate**: confirm or update the Discovery Goal and priorities.
3. **If scenario changed**: the developer presents the trade-off in gettoni ("this costs 3, we have 8 left, the remaining work costs 10 — what's more important?") and proposes an adjustment. The client decides priorities, the developer decides how.
4. **At fetta delivery**: the developer proposes the next fetta. The client approves or changes priorities. Max 1 hour — if it takes longer, there's a deep problem to address.

### What the client sees

- The working product growing
- Total and remaining gettoni
- Priorities confirmed or to discuss

### What the client doesn't see

- Internal gettoni allocation (infra vs client split)
- Internal estimates
- Architectural details

---

## Sprint Closure

1. **Final walkthrough**: full tour with the client. The criterion is conformity: does the software make the necessary data available to answer the Discovery Goal and the agreed sub-experiments? The developer delivers the instrument — the business experiment is the client's responsibility.
2. **Formal ok**: the client confirms.
3. **Delivery**: deliverables + delivery support, OR working hosted instance transferred. (Domain-specific delivery details belong in the domain or project layer.)
4. **Growth plan**: Reserve items become the roadmap for a potential second sprint.

---

## How Prove Out Modifies Pipeline Steps

| Step | Modification |
|---|---|
| `/v-spec` | Translate workshop materials: Tiny Experiment fields (domanda, vero_se, tasks, non farò) become FRs and edge cases. Pillar/Bridge/Reserve becomes requirement priority. |
| `/v-plan` | Use breadboard as architecture starting point. Order fette by decreasing uncertainty (from sub-experiments). Note gettoni allocation per fetta. |
| `/v-status` | Show budget: package, total, spent, remaining, per-fetta breakdown. |
| `/v-bug` | Note: "Bug fixes are budget-free. If significant work beyond scope, evaluate at next checkpoint." |
| `/v-guide` | Topic "prove-out": explain the method. Topic "budget/gettoni": explain budget system. |
| `/v-report` (client) | Use client-facing language: gettoni remaining, features delivered, decisions needed. No technical jargon, no infra/client split. |
| `/v-compound` | At fetta completion: generate checkpoint data (gettoni spent, features completed, remaining budget). |

---

## Glossary

| Term | Meaning |
|---|---|
| **Prove Out** | The framework name. "Let's put the idea to the test." |
| **MVP** | Minimum Viable Product — the minimum instrument to collect data and validate a business hypothesis |
| **Discovery Goal** | The formalized business experiment: hypothesis, metric, sub-experiments, boundaries |
| **Sub-experiment** | Intermediate question to explore in order to answer the Discovery Goal |
| **Tiny Experiment** | A single feature documented as a question ("can we do X?") or a direct task, with a truth criterion and bounded tasks |
| **Gettone** | Relative unit of work complexity. Makes trade-offs visible |
| **Package** | Fixed budget in gettoni with associated duration. S = 20 gettoni / 3 weeks, M = 35 gettoni / 5 weeks |
| **Horizontal slice (fetta)** | A cut of the product that crosses the entire user journey and produces something working |
| **Scooter** | The first fetta — the entire path in the most reduced way possible |
| **Pillar** | Map step that gives data for the experiment. Where most investment goes |
| **Bridge** | Map step without which the user gets stuck, but that doesn't give experiment data |
| **Reserve** | Step with value but not needed for the MVP. Growth plan for later |
| **Breadboard** | Exploration format from Shape Up: places, actions, connections. Not visual — keeps away from interface details |
| **Checkpoint** | Meeting every 3-4 days: product demo, remaining gettoni, priorities to confirm or change |
| **Pact** | The set of rules and mutual commitments that found the CTO-client collaboration |
| **Tasks** | Actions to complete a Tiny Experiment |
| **Non farò** | Explicit boundaries of what we don't do in a Tiny Experiment |
| **Starò attento a** | Identified risks that could drag you off course (Rabbit hole in Shape Up) |

---

## When NOT to Use Prove Out

- The client has already validated and wants a finished product
- The project is maintenance or evolution of an existing product
- The client has no hypothesis to test (no Discovery Goal = no success criterion)
- The client wants to control the "how" (CTO prerogative is non-negotiable)
