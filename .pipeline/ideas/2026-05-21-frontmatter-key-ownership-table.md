---
slug: frontmatter-key-ownership-table
created: 2026-05-21
source: compound learning Fetta 2 (semantic collision `stato`)
status: in-progress (spec started 2026-05-21, see .pipeline/spec.md)
priority: P3
pulled_into_spec: 2026-05-21 — assorbita in Sprint 1.5.B (config master TOML + profile-override). La tabella key-ownership diventerà esplicita nello schema TOML.
---

# Idea — Tabella ownership chiavi frontmatter (engine vs profile default)

## Contesto

Durante il fix di `frontmatter-default-non-applicato` (bug #006, Fetta 2 review loop 1), ho scoperto che il campo `stato` nei profili YAML aveva semantica diversa dal `stato` runtime impostato dall'engine:

- **Engine (runtime)**: `stato: completato | timeout | parziale_no_done` — descrive l'esito dello stream LLM.
- **Profilo default (intent)**: `stato: bozza_da_validare_qm` — descrive il workflow QM (output ancora da validare da Denis).

Stessa chiave, due semantiche. Fix immediato: rinominare nel profilo a `stato_qm` (5 profili).

Il problema sottostante: **non c'è una tabella documentata di chi possiede ogni chiave del frontmatter**. Un futuro profile default potrebbe ri-collidere con `modello`, `provider`, `data_esecuzione` se l'autore non sa cosa l'engine genera.

## Proposta

Aggiungere a `.pipeline/standards/02-project.md` (o `02-project.md` Project-Specific Conventions) una **tabella esplicita**:

| Chiave frontmatter | Owner | Valori tipici | Note |
|---|---|---|---|
| `profilo` | engine | nome del profilo | sovrascrive sempre il default |
| `modello` | engine | stringa qwen3.6 / mistralai/... | runtime, dal provider |
| `provider` | engine | ollama / eurouter | runtime |
| `data_esecuzione` | engine | ISO 8601 | runtime |
| `durata_secondi` | engine | float | runtime |
| `token_stimati` | engine | int | runtime |
| `file_input` | engine | array di nomi file | runtime |
| `stato` | engine | completato / timeout / parziale_no_done | runtime execution status |
| `valutazione_denis` | esterno (Denis manuale) | validata / validata_con_riserva / non_validata | popolato post-L2 |
| `tipo` | profile default | management_review_pack / capa_pack / equipment_alert / ... | semantic capability type |
| `stato_qm` | profile default | bozza_da_validare_qm | workflow QM state |
| `profilo_labnexus` | profile default | review-pack / audit-checklist / ... | redundancy per filtering |
| `locale` | profile default | true | declarazione esecuzione on-device |

E una regola: **se aggiungi una chiave a `output.frontmatter_default` di un profilo, verifica che NON collida con nessuna delle chiavi "engine" sopra**. Se collide, scegli una chiave alternativa.

## Stima

- 1 gettone infra (~30 min CTO) per:
  - scrivere la tabella in `02-project.md`
  - aggiungere un test unit in `internal/output/output_test.go` che assert: per ognuno dei 5 profili shippati, le chiavi di `frontmatter_default` NON sono in {profilo, modello, provider, data_esecuzione, durata_secondi, token_stimati, file_input, stato}
  - documentare nel template di `feat-meta` (Fetta 3) la regola per il prompt meta-generator

## Quando

Sprint 2 hardening cluster, o subito quando si aggiunge `feat-meta` (Fetta 3): il prompt meta deve sapere quali chiavi NON usare come default per evitare collisioni future.

## Riferimento incrociato

- Solution: `.pipeline/solutions/2026-05-21-fetta2-3-capability-declarative-piu-hardening-security.md`
- Bug origine collision: `.pipeline/bugs/frontmatter-default-non-applicato.md` (fixed Fetta 2)
- Skill correlata: in attesa di promozione (potenziale `frontmatter-key-namespace-management.md`)
