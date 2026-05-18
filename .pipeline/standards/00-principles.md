# Principles

These principles guide every decision. They apply regardless of stack, project, or client.
When in doubt, come back here.

## Sottrarre è moltiplicare

Keep things simple. Simplify complex things. Do fewer things.

Simplicity has THREE beneficiaries, all equally important:
- **Simple for the developer**: the code is readable, the architecture is clear, the tools are consistent. A developer can onboard in hours, not days.
- **Simple for the user**: the interface is intuitive, the behavior is predictable, friction is minimal.
- **Simple for the client**: the system is maintainable, the costs are proportional, the technology choices are understandable.

Simplicity does NOT mean "fewer tools" or "no abstractions." It means choosing the right tools and abstractions that reduce total complexity for all three parties. A high-quality framework (FastAPI, Rails, SvelteKit) is simple because it absorbs complexity that would otherwise leak into your code. An ORM is simple because it handles SQL injection prevention, query building, and schema management — making the code more secure AND more readable. Avoiding an ORM and writing raw SQL "because it's simpler" is a false economy if it produces less secure, less readable code.

**The real test of simplicity**: does this choice reduce the total cognitive load? Not just for writing code today, but for reading it in 3 months, for debugging it at 2am, for explaining it to the client.

**What simplicity IS**:
- Using a proven framework instead of reinventing patterns
- Using an ORM consistently instead of mixing ORM and raw SQL
- Choosing ONE ecosystem and going deep (JavaScript, not TypeScript — less tooling, less config, fewer failure modes)
- Keeping the same ORM, the same framework conventions, the same patterns across projects — consistency IS simplicity
- sqlite over PostgreSQL when the scale doesn't demand otherwise

**What simplicity is NOT**:
- Avoiding all dependencies on principle
- Writing raw SQL "because ORMs are magic"
- Skipping a framework to have "more control"
- Mixing approaches within the same project (ORM here, raw SQL there)
- Changing tools every project because something new looks interesting

Over-engineering is a defect. But so is under-engineering that creates security holes or maintenance nightmares in the name of "keeping it simple."

**Scale appropriately**: design for the actual scale — tens of concurrent users, not thousands. Infrastructure choices reflect reality, not aspirations. Premature complexity is a form of over-engineering. If the project doesn't need PostgreSQL, use SQLite. If it doesn't need microservices, use a monolith. Revisit when the scale actually changes, not when it might.

## C'è un modo

There is always a way to unblock a situation, even unconventional ones. But creative solutions are a LAST RESORT, not a first choice. Before going out of the box, exhaust the standard approaches. An unconventional solution that bypasses quality standards is not a solution — it's debt.

**The test**: if you're reaching for a creative workaround, first prove that the standard approaches genuinely don't work. If they do work but seem "less elegant," use them anyway. Elegance is not a goal — working, simple, quality code is.

## Decidi, fai, osserva, impara

Decide, do, observe, learn. This is the continuous cycle that drives discovery and improvement.

Don't overthink before acting. Don't skip observation after acting. Don't waste observations by not learning. Don't learn without deciding what to do next. The cycle is fast and tight — not waterfall planning followed by months of execution.

**In practice**: ship small increments, measure what happens, adjust. If a design decision turns out wrong after observation, change it immediately. Sunk cost is not a reason to keep a bad decision.

## Obsess over quality

Quality of the result is the guiding star. Quality is measured by a metric that is shared and defined with the client — not by the developer's aesthetic preferences.

Quality means: it works correctly, it's maintainable, it's secure, and it meets the client's definition of "good." If the client defines quality as "fast to market with acceptable trade-offs," that's the quality bar. If the client defines quality as "zero bugs, fully tested," that's the bar. The bar is explicit, not assumed.

**In practice**: define the quality metric with the client BEFORE writing code. Every review, every test, every decision references that metric. "Is this good?" is answered by "does it meet the agreed quality bar?" — not by "does it look good to me?"

## Tempo fisso, scope variabile

Time and cost are fixed. Scope is the variable. This applies at every scale: a sprint, a feature, a single work session.

When time runs out, you stop and evaluate — not push through. What's done? What's missing? Is what's missing worth more time, or can the scope be reduced to fit? The answer is always to simplify scope, never to extend time silently.

This principle prevents the two most common project failures: scope creep (building more than budgeted) and silent overrun (spending more time than planned without acknowledging it). Both destroy trust with the client and undermine quality.

**In practice**: every unit of work has a budget (gettoni, hours, story points — the unit doesn't matter, the constraint does). When the budget is consumed, the work is evaluated against the budget, not against an ideal of completeness. "Done enough" is a valid and honest outcome. "Almost done, just need a bit more time" is a warning sign.

## Sempre funzionante (always working)

After every iteration — every ship, every fetta, every meaningful unit of work — the product is usable end-to-end. It may be reduced (a scooter, not a car) but it WORKS. The user can go from start to finish.

Never leave the product in a state where "this part works but you can't use it yet because the other part isn't done." That's a vertical slice (depth on one area) instead of a horizontal slice (thin but complete journey). Build horizontally: thin slices that cross the entire user journey.

**In practice**: if you're building a feature that requires 3 components to be useful, don't build component A completely, then B, then C (vertical). Build a thin version of A+B+C that works together (horizontal), then enrich. The first ship is a scooter — ugly, minimal, but it gets you from A to B. Each subsequent ship is richer.

## Sicurezza by design

Security is a primary guiding principle in every implementation and design choice — not an afterthought, not a checklist, not a layer added at the end. Every decision, from framework selection to data handling to deployment, considers security implications from the start.

**In practice**: use frameworks and libraries that handle security fundamentals (CSRF, XSS, SQL injection) by default. Never roll your own auth, encryption, or session management when a well-maintained library exists. Validate at system boundaries. Store secrets in environment variables, never in code. Default to the most restrictive permissions and open up only as needed. If a shortcut weakens security, it's not a shortcut — it's debt with compounding interest.

## La documentazione è codice (docs are code)

Documentation that doesn't match the code is worse than no documentation — it actively misleads. Every code change that affects behavior, API, configuration, or user-facing functionality MUST include corresponding documentation updates in the same commit.

This is not a "nice to have" or a "we'll update the docs later." Docs and code are one unit. A feature is not done until its documentation matches. A refactor is not done until affected docs are updated. A bugfix that changes behavior updates the relevant docs.

**In practice**: the `/v-review` checks documentation alignment with the same severity as any other quality finding. If the work says one thing and the docs say another, it's a HIGH finding — not a LOW "docs cleanup." The `/v-execute` step produces docs alongside the work, not after. README, help text, config comments, inline documentation — all are verified against actual behavior. Stale docs are bugs.
