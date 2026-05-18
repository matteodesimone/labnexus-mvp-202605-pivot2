# Domain-Specific Principles (Coding)

These principles extend the universal principles in `00-principles.md`.

## Privacy by design

Privacy is a fundamental right, not a feature to bolt on. Every system that handles personal data is designed from the ground up with GDPR principles and privacy protection as architectural constraints.

This is non-negotiable regardless of where the client or users are based. GDPR is the FLOOR, not the ceiling.

**The principles in practice**:
- **Data minimization**: collect only what is strictly necessary. If you don't need a birth date, don't ask for it. If you can work with anonymized data, do it.
- **Purpose limitation**: data collected for one purpose is not repurposed without explicit consent.
- **Consent**: explicit, informed, revocable. No dark patterns, no pre-checked boxes, no "by continuing you agree."
- **Right to access and deletion**: every system that stores personal data must have a clear mechanism for users to see what's stored and request its deletion. Design the data model with this in mind from day one.
- **Data protection by default**: the most privacy-friendly settings are the default. Users opt IN to sharing, not opt OUT.
- **Encryption**: personal data encrypted at rest and in transit. No exceptions.
- **Logging**: logs never contain personal data (names, emails, IPs) unless strictly necessary for security, and even then with retention limits.
- **Third parties**: every third-party service that touches personal data is evaluated for GDPR compliance. No "we'll check later."

**In architecture decisions**: choose solutions that keep data under your control. Self-hosted over SaaS when handling sensitive data. EU-based infrastructure when serving EU users. Data processing agreements with every sub-processor.

## Intuitive by design

Every interface — GUI, CLI, API, error message, game menu — is intuitive or it is broken. This is not a subjective quality assessment: it follows Everett N. McKay's objective definition.

**An interface is intuitive when the target user can use it without resorting to reasoning, memorizing, experimenting, seeking help, or training.**

If the user has to read a manual, memorize a flag, guess what a button does, or experiment to discover a feature, the interface has failed. This applies to web dashboards AND to `vibbly init`. There is no exception for "developer tools" — developers are users too.

**The eight attributes (McKay)**:
1. **Discoverability** — the user can find features without searching
2. **Affordance** — elements communicate what they do through appearance or naming
3. **Comprehensibility** — the user knows what will happen before acting
4. **Responsive feedback** — the interface confirms what happened after every action
5. **Predictability** — same action, same result, no surprises
6. **Efficiency** — minimum steps to accomplish the task
7. **Forgiveness** — the user can recover from mistakes
8. **Explorability** — the interface is safe to explore

**The main instruction**: every screen, command, or page has ONE main instruction — a single sentence telling the user what to do. If you can't write it in one sentence, the screen does too much. Split it.

**In practice**:
- `/v-spec` requires a main instruction for every feature with a user interface (including CLI commands)
- `/v-plan` references McKay attributes in design decisions: "delete uses `--dry-run` by default (forgiveness + explorability)"
- `/v-review` evaluates every interface element against the 8 attributes. Violations are findings, not suggestions.
- Severity: CRITICAL if discoverability or comprehensibility fails on primary tasks. HIGH if forgiveness fails on destructive actions. MEDIUM for efficiency/predictability gaps. LOW for minor affordance issues.

**The full skill with per-interface-type examples is in `.pipeline/skills/framework/intuitive-design.md`.** Read it before designing any interface.
