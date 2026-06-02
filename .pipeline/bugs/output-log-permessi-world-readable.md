---
severity: medium
status: open
created: 2026-06-02
source: external reviewer (Codex security/privacy) durante /v-review bugfix trigger-override-job
area: internal/output/output.go + internal/runner (openAuditTrail)
---

# Bug: output `.md` e `.log` di audit scritti world-readable (0o644 / 0o755)

## Description

Gli output markdown sono scritti con `0o644` e i `.log` accoppiati creati via
`os.Create` sotto cartelle `0o755`. Questi file contengono contenuti SGQ,
possibili dati personali, output del modello e il testo grezzo della risposta
streamata. Su una macchina macOS/Linux multi-utente, altri utenti locali possono
leggerli.

Distinto da `privacy-paths-personali-in-log.md` (quello = path personali NEL
contenuto del log; questo = permessi del file sul filesystem).

## Expected

Permessi privati per artefatti sensibili: cartelle output `0o700`, file output/log
`0o600` (`os.OpenFile(..., 0o600)` per l'audit log).

## Note

Pre-esistente, non introdotto dal trigger-override (il fix aggiunge una riga al
.log ma non ne cambia i permessi). Violazione privacy-by-design → pipeline
dedicata. Deferred.
