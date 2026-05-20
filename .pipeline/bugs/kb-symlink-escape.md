---
severity: high
status: fixed
created: 2026-05-19
fixed_at: 2026-05-20
source: review esterna (codex-security)
fix: validateKbFiles in internal/profile/profile.go ora usa filepath.EvalSymlinks su kbDir + ogni candidato kb_file, e verifica che il path REALE (post symlink resolve) rimanga dentro la kbDir reale. Defense in depth: il check lessicale precedente è mantenuto. Il caso legittimo "kbDir stesso è un symlink (KB-ispettore root)" è esplicitamente coperto da test.
test: internal/profile/profile_test.go::TestValidate_RejectsKbFileSymlink (+ TestValidate_AcceptsKbDirItselfSymlink per il caso legittimo)
---

# Bug: `KB-ispettore/` simlink possono uscire dalla cartella KB ed esfiltrare file

## Reported

- **Feature/Area**: `internal/profile/profile.go:validateKbFiles` + `internal/runner/runner.go:loadKB`
- **Trovato da**: review multi-LLM 2026-05-19 (codex-security HIGH)
- **Sprint impact**: HIGH per security. Pre-esistente.

## Description

`validateKbFiles` controlla containment via path lessicale (`strings.HasPrefix`). Non risolve simlink. Quindi:
```
KB-ispettore/
  leak.md -> ~/.ssh/id_rsa
```
passa la validazione lessicale (path inizia con `KB-ispettore/`), poi `runner.loadKB` fa `os.ReadFile` seguendo il simlink → contenuto privato finisce nel prompt LLM (e con EUrouter, esfiltrato).

Vector di attacco: pacchetto KB malevolo o KB compromesso. Sprint 1 distribuisce KB-ispettore al laboratorio Denis senza checksum → realistico.

## Fix proposto

In `validateKbFiles`: dopo controllo path lessicale, fare `filepath.EvalSymlinks` sia su `kbDir` che sul candidato, verificare di nuovo containment. Reject se simlink target esce.

Alternativamente più stringente: rifiutare qualsiasi simlink dentro KB-ispettore (semplice + sicuro per Sprint 1).

## Risk of fix

**Basso**. Modifica isolata in `profile.go`. Test BDD: "KB con simlink fuori cartella deve essere rifiutato".

## Note

Deferred to backlog: pre-esistente, ma alta severità. Da valutare per Fetta 2.
