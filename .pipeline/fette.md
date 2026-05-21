# Fette di Sprint 1.5 — LabNexus refactor pre-handoff Denis (pivot 3 + concierge mode)

**Slicing approvato dal CTO** (allineato al brainstorm `2026-05-21-refactor-sprint1-concierge-mode.md` Approach B): 3 sub-fette ordinate, 0 gettoni cliente (offerti da Matteo), ~23-36h CTO totali.

**Strict ordering vincolante**: 1.5.A → 1.5.B → 1.5.C. Le sub-fette sono naturalmente dipendenti: 1.5.A è prerequisito per qualunque uso reale (eurouter default); 1.5.B abilita 1.5.C (config master + parser nuovi formati sono prerequisito per concierge mode).

| # | Nome | Stima CTO | FRs coperti | Cosa funziona alla fine | Stato |
|---|---|---|---|---|---|
| 1.5.A | Pivot 3 + CLI puro | 5-7h | FR-1 → FR-10 | Il binary nudo (no `.app`) gira con eurouter di default. Doppio click su `labnexus.command` apre Terminal + TUI. Gate cloud rimosso. 7 profili shippati con `provider: eurouter`. Bug file dei gate aggiornati a `intentional_deviation`. BDD `@manual` su `.app` eliminati. README aggiornato per Gatekeeper-binary. Review loop 2 ha fixato 11 doc-drift findings + audit-trail mismatch sui frontmatter `locale:`. | **deployed 2026-05-21** ✓ |
| 1.5.B | Config master TOML + parser nuovi formati | 9-13h | FR-11 → FR-24 | `labnexus.config.toml` master con defaults globali. Profili `.toml` con campi opzionali (eredità da master, override locale). `trigger_prompt_file` indirezione (XOR con `trigger_prompt` inline, path-safe). Parser `.xls`, `.doc`, `.rtf` in `internal/input`. Meta-prompt aggiornato a TOML. Tabella ownership chiavi frontmatter documentata. Review loop 2 ha fixato 11/14 findings (8 HIGH + 2 MEDIUM critici + 1 LOW); 8 MEDIUM/LOW deferred a `.pipeline/ideas/` per 1.5.C+. Smoke verde su `data/CAPABILITY G` con parser .xls reali Denis. | **deployed 2026-05-21** ✓ |
| 1.5.C | Concierge mode + auto-discovery + audit log + streaming live | 9-16h | FR-25 → FR-39 (+ EC-1 → EC-16) | Cartella `lavori/` rinominata da `data/`. Auto-discovery via `_labnexus.toml` per cartella + fallback convention naming. Subcommand `labnexus jobs` + flag `--job <name>`. ~~TUI estesa con menu pre-mappato~~ (skip giustificato, decision plan-time: CLI `labnexus jobs + run --job` copre il use case). Logger multi-writer (stderr + file `.log` per ogni run) come audit trail ISO 17025. Body streaming live in stderr + log mirror. Output dentro `lavori/<X>/output/`. Documenti aggiornati con tono concierge italiano. Review loop 2 ha fixato 14 findings (3 CRITICAL security defense-in-depth + 5 HIGH + 6 MEDIUM/LOW). Smoke estensivo verde su zip estratto in tmpdir. | **deployed 2026-05-22** ✓ |

**Totale**: 5-7h + 9-13h + 9-16h = **23-36h CTO equivalent**.

## Razionale dello slicing

- **1.5.A è prerequisito**: senza eurouter default, Denis non può manco lanciare labnexus (no hardware Ollama). Senza `.command` launcher, non può aprire da Finder. Senza rimozione del gate, l'invocazione eurouter di Denis fallisce silentamente. Quindi 1.5.A è la "scooter" minimale: dopo 1.5.A il binary è eseguibile su Mac Denis con eurouter (anche se l'UX è degradata, no log audit, no auto-discovery, configurazione manuale).
- **1.5.B abilita 1.5.C**: il config master TOML è la fonte di verità per i defaults (eurouter, modello, API key). Senza di esso, ogni profilo deve duplicare la config (gestione ostile). I parser `.xls`/`.doc`/`.rtf` sono prerequisito perché senza, Denis NON può aprire i suoi materiali reali (impedirebbero qualunque test concierge).
- **1.5.C completa la UX concierge**: auto-discovery + log audit + streaming live sono "polish" che rendono labnexus un vero strumento di lavoro quotidiano. Senza, è un toolkit CLI grezzo che funziona ma non è confortevole per Denis.

## Pattern di validazione (vale per tutte le sub-fette)

Ogni sub-fetta passa per il cycle pipeline completo: plan → test-scaffold → implement → review → deploy → compound. Niente shortcut. Review = REVIEW="all" come da config.

## Checkpoint cliente

- **Fine 1.5.A**: handoff intermedio possibile a Denis se serve testare urgente con eurouter "a mano" (config in profile). Probabile decisione: completare 1.5.B+C prima di handoff effettivo.
- **Fine 1.5.B**: il sistema è funzionalmente completo per uso CLI, ma manca la UX concierge + audit log. Possibile handoff intermedio se Denis ha urgenza.
- **Fine 1.5.C**: handoff completo. Sprint 1.5 chiuso. Denis riceve lo zip e parte. Inizio shakedown reale + L2 Denis sulle 7 capability + FR-22 3 casi Claude (gettoni cliente residui Sprint 1).

## Possibilità di ricalibrazione

Se durante 1.5.A emergono complicazioni (es. rimuovere il gate spezza altri test BDD non previsti), il CTO ridiscute lo scope di 1.5.B/C. Stabilità di scope (NFR-9 invariato): niente nuovi rami di sviluppo senza dati raccolti.

Se Denis dà feedback urgente prima della chiusura 1.5.C che richiede modifiche (es. "voglio drag&drop dal Finder dopotutto"), il CTO valuta se aggiungere ALLA SUB-FETTA CORRENTE o rimandare a Sprint 2.

## Spec di riferimento

- **`.pipeline/spec.md`** — Sprint 1.5 spec full (39 FR + 11 NFR amendments + 16 EC).
- **`.pipeline/brainstorms/2026-05-21-refactor-sprint1-concierge-mode.md`** — brainstorm sorgente con tutti i razionali decisionali + stress test.
- **`.pipeline/proposed-features/`** — 3 file BDD scenari (uno per sub-fetta).
- **`.pipeline/proposed-test-data/`** — fixture DRAFT (config master, _labnexus.toml, profilo TOML, README).
