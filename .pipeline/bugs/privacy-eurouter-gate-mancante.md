---
severity: critical
status: open
created: 2026-05-19
source: review esterna (codex + mistral, CONFIRMED multi-reviewer)
fix: ""
test: ""
---

# Bug: provider EUrouter selezionabile senza gate di approvazione → rischio esfiltrazione dati SGQ

## Reported

- **Feature/Area**: FR-7 / FR-12 (selezione provider) — `internal/provider/provider.go:Select` + override env/CLI
- **Trovato da**: review multi-LLM 2026-05-19 (codex-privacy + mistral-privacy, CONFIRMED)
- **Sprint impact**: HIGH per privacy. Pre-esistente, NON introdotto dai cycle bugfix-5/6/UX-progress.

## Description

`provider.Select` accetta `LABNEXUS_PROVIDER=eurouter` o `--provider eurouter` senza alcun gate di approvazione esplicita. Il vincolo di progetto è che EUrouter NON deve processare dati reali del SGQ di Denis senza approvazione esplicita (cfr. `docs/piano_iniziale/` + `.pipeline/standards/02-project.md` Privacy section). Tecnicamente niente impedisce a Denis (o a chiunque setti l'env var) di inviare dati personali a `api.eurouter.ai`.

## Fix proposto

Tre opzioni (decisione CTO):
1. Compile-flag `release` che disabilita totalmente EUrouter per i binari distribuiti a Denis (cleanest).
2. Flag `--allow-external-provider` + env `LABNEXUS_EUROUTER_APPROVED=1` obbligatorio per attivare EUrouter (esplicito).
3. Warning prominente stderr "questo run invierà dati a api.eurouter.ai. Premere Enter per continuare" (interattivo, fragile su non-TTY).

Raccomandazione: 1+2 (binario release senza EUrouter, sviluppo via env approval).

## Risk of fix

**Medio**. Tocca interfaccia public di `provider.Select` + `cmd/labnexus/main.go`. Necessari test BDD per ENC-6 (EUrouter gate). Stima ~2-3h.

## Note

Deferred to backlog: pre-esistente, fuori scope cycle bugfix-5/6/UX-progress (chiuso 2026-05-19).
