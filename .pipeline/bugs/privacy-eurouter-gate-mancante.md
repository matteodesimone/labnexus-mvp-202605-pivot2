---
severity: critical
status: intentional_deviation_post_pivot_3
created: 2026-05-19
fixed_at: 2026-05-20
deviation_at: 2026-05-21
source: review esterna (codex + mistral, CONFIRMED multi-reviewer)
fix: provider.Select ritorna ErrEurouterGateMissing se LABNEXUS_ALLOW_CLOUD_PROVIDER non è impostato. Senza gate, eurouter è rifiutato indipendentemente da flag/env/profilo. BDD helpers settano il gate per default (ambiente test controllato). Gate test in internal/provider/select_test.go.
deviation_rationale: |
  Sprint 1.5 pivot-3 ha invertito la decisione strategica: il deliverable per Denis ora gira
  con provider eurouter di default perché Denis non ha hardware Ollama. Cliente (Stefano Fiorina)
  formalmente informato e d'accordo (vedi state.json sprint1.5 approval). Il gate
  `LABNEXUS_ALLOW_CLOUD_PROVIDER` è stato RIMOSSO dal codice in Sub-fetta 1.5.A perché
  il suo scopo era proteggere dati SGQ contro esfiltrazione cloud — ora l'esfiltrazione
  è la NORMA approvata. La decisione cloud/local passa esclusivamente a configurazione
  utente (labnexus.config.toml master in 1.5.B).
deviation_test_removed: internal/provider/select_test.go::TestSelect_EurouterRejectedWithoutApprovalGate (+3 varianti) — RIMOSSI in Sub-fetta 1.5.A (decommissione del gate)
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
