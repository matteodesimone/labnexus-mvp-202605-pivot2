---
id: 002
priority: P1
status: done
source: bug
created: 2026-05-20
file: internal/profile/profile.go
bug_ref: .pipeline/bugs/kb-symlink-escape.md
---

# KB-ispettore: symlink possono uscire dalla cartella ed esfiltrare file privati

## Issue

`internal/profile/validateKbFiles` controlla solo il lexical containment (`filepath.Clean` + prefix check). `loadKB` poi legge i file con `os.ReadFile`, che **segue i symlink**. Una entry KB come `kb_files: [evil-link]` dove `evil-link → /Users/denis/.ssh/id_rsa` o `/private/SGQ-altra-cartella/` farebbe iniettare il contenuto referenziato nel prompt, esfiltrando potenzialmente segreti o documenti di altri laboratori se il provider scelto è eurouter.

## Fix

In `internal/profile/profile.go` (validateKbFiles e loadKB):
1. Usare `filepath.EvalSymlinks` sia per la root KB sia per ogni candidate path PRIMA del check di containment.
2. Confrontare la root REALE col path REALE: se il path reale non è dentro la root reale → rifiuta con errore.
3. Alternativa più semplice per Sprint 1: rifiutare in toto entries che siano symlink (`os.Lstat().Mode() & os.ModeSymlink != 0` → errore).

Test:
- Unit test in `internal/profile/profile_test.go` con symlink test (creato con `os.Symlink`).
- BDD scenario: profilo con `kb_files: [out-of-tree-symlink]` → exit 2 + stderr "symlink non consentito".

## Context

- Confermato da Codex L2 review Fetta 2.
- **Blocker per consegna a Denis**: lo zip distribuito a Denis include la KB-ispettore via `zip -r` (resolve symlink). Sul SUO sistema, eventuali manipolazioni accidentali della KB potrebbero introdurre symlink involontari.
- Convertito in todo P1 da `/v-triage` 2026-05-20.
- Stima fix: 1-2 gettoni infra (~30min-1h CTO + test).
