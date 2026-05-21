---
title: config.Merge — pointer types per distinguere "campo non setted" da "esplicitamente 0"
status: planned
source: review 1.5.B loop 1 HIGH-config-1 (mistral test-coverage CONFIRMED)
created: 2026-05-21
target: Sprint 1.5.C+ o Sprint 2 (refactor invasivo)
---

# Idea: pointer types in profile.Profile per detection "campo non setted"

## Problema

`internal/config/Merge` usa `== ""` e `== 0` per decidere se un campo del profile è "non setted" e quindi va ereditato dal master. Per i campi numerici (`Temperature float64`, `MaxTokens int`, `ContextWindow int`), questo non distingue:

- **caso A**: profilo TOML non dichiara `temperature` → zero-value → override col master (corretto)
- **caso B**: profilo TOML dichiara `temperature = 0.0` esplicitamente → zero-value → override col master (BUG)

Use case Denis che innesca il bug: profilo `pt-analysis` (PT z-score) potrebbe volere `temperature = 0.0` per deterministico (calcoli numerici delicati, no varianza). Setting nel profile viene silenziosamente override col master 0.9.

## Soluzione proposta

Cambiare i campi numerici di `Profile` a pointer types:

```go
type Profile struct {
    // ...
    Temperature   *float64 `toml:"temperature"`
    MaxTokens     *int     `toml:"max_tokens"`
    ContextWindow *int     `toml:"context_window"`
    // ...
}
```

`toml.Unmarshal` lascia il pointer nil se il campo è omesso, lo popola se è presente (anche con 0). `Merge` testa `cp.Temperature == nil`.

## Scope del refactor

Invasivo: tocca ~5 file:
- `internal/profile/profile.go` — struct + Validate (deref con `*` quando si controlla)
- `internal/config/config.go` — Merge logic
- `internal/runner/runner.go` — usi di `p.Temperature` ecc.
- `internal/provider/provider.go` (forse Options struct se passa direttamente)
- Tutti i test che creano Profile inline (`runner_test.go`, `profile_test.go`, BDD helpers)

Stima: 2-3h CTO + ri-test completo.

## Workaround temporaneo (Sprint 1.5.B)

Documentato come known limitation in `internal/config/config.go` comment di Merge. Denis che vuole `temperature = 0` setta `temperature = 0.001` come pragmatic substitute (la varianza differenza è negligible per output LLM, ma il valore != 0 evita il merge override).

## Decisione differita

Rimandato a 1.5.C (concierge mode) o Sprint 2. La review 1.5.B loop 2 ha applicato workaround documentato per non sforare lo scope sub-fetta.
