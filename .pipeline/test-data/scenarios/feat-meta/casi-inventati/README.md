# Feat-meta casi inventati — FR-22

Validazione funzionale del meta-prompt `docs/meta-prompt-genera-profilo.md`: tre casi inventati di complessità crescente. Per ciascun caso:

1. Matteo (CTO, FUORI pipeline AI) apre `claude.ai`, incolla l'intero contenuto del meta-prompt.
2. Descrive il caso a Claude (testo "Domanda utente" sotto).
3. Segue le eventuali clarification questions di Claude.
4. Salva il YAML prodotto nel file dichiarato in tabella.
5. Annota numero di clarification questions e note in `note-esecuzione-<N>.md`.
6. Esegue: `labnexus validate <case-N>` (con `--profiles-dir .pipeline/test-data/scenarios/feat-meta/casi-inventati/ --kb-dir docs/piano_iniziale/materiali-dominio/KB-ispettore`). DEVE ritornare exit 0.

Il test Go `internal/profile/profile_meta_cases_test.go` cicla tutti i `caso-*.yml` e li valida automaticamente in regressione.

## I 3 casi (da `.pipeline/spec.md` FR-22)

| # | Complessità | Domanda utente (da incollare in Claude DOPO il meta-prompt) | YAML atteso in | Clarification questions attese |
|---|---|---|---|---|
| 1 | semplice | "Voglio un profilo per la verifica di conformità di un certificato di taratura." | `caso-1-semplice.yml` | 0 (Claude produce subito) |
| 2 | medio | "Voglio un profilo per identificare deviazioni dai limiti di accettabilità in un set di rapporti di prova." | `caso-2-medio.yml` | ≥ 1 (Claude chiede chiarimenti su quali limiti, quale tipo di rapporto, ecc.) |
| 3 | ambiguo | "Voglio un profilo per la qualifica dei fornitori." | `caso-3-ambiguo.yml` | ≥ 2 (Claude riconosce ambiguità: qualifica iniziale vs. continua, fornitori beni vs. servizi vs. subappalti, fornitori critici vs. non critici, ecc.) |

## Vincoli sul YAML prodotto

Tutti e tre i YAML DEVONO:

- Passare `labnexus validate` con KB-ispettore reale (FR-22).
- Avere tutti i campi obbligatori dello schema (FR-3): `profilo`, `descrizione`, `provider`, `modello`, `kb_files`, `trigger_prompt ≥ 50 caratteri`.
- Avere `kb_files` che esistono fisicamente in `docs/piano_iniziale/materiali-dominio/KB-ispettore/` (no path inventati).
- `provider` ∈ {`ollama`, `eurouter`}.
- Idealmente: `output.frontmatter_default` con `tipo`, `stato_qm`, `profilo_labnexus`, `locale: true` (pattern Fetta 2 — il meta-prompt deve istruire Claude su questo).

## Stato (red phase /v-test-scaffold)

Al momento i 3 file `caso-*.yml` **non esistono**. Il test Go `profile_meta_cases_test.go` è in red (skip se zero file, o failure esplicito) finché Matteo non esegue manualmente i 3 casi in Claude esterno + commit dei YAML qui.

In `/v-implement` (Fetta 3 green phase) scriverò `docs/meta-prompt-genera-profilo.md` e la guida; tu lanci i 3 casi in Claude e committi i YAML; il test Go diventa verde.

## Note operative

- Se Claude produce un YAML che non passa `labnexus validate`, NON correggerlo manualmente: itera il prompt in Claude finché il YAML è valido al primo colpo. Quello è il vero test del meta-prompt.
- Se invece Claude chiede meno clarification questions del previsto al caso 2/3, è segnale che il meta-prompt va rinforzato. Iterabile.
- I YAML committati sono **fittizi** (test data): non vanno copiati in `profili/` per uso reale. Sono solo prova funzionale del meta-prompt.
