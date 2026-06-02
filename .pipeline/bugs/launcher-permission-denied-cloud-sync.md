---
severity: high
status: fixed
created: 2026-06-02
fixed: 2026-06-02
source: Stefano (run reale dal deliverable su kDrive)
fix: "scripts/labnexus.command: auto-heal del bit eseguibile prima di lanciare il binary (chmod +x ./labnexus se !-x), + distinzione esplicita dei casi 'non trovato' / 'non eseguibile/chmod fallito' / 'crash'. README.md: nuova sezione 'Cartelle su cloud-sync' con workaround chmod +x. Verificato end-to-end in /tmp: binary senza +x → launcher lo ripristina e avvia."
test: "Verifica funzionale manuale in temp dir (binary -rw-r--r-- → -rwxr-xr-x dopo il launcher). Nessun test automatizzato shell nel progetto."
---

# Bug: doppio click → `./labnexus: Permission denied` (deliverable su kDrive)

## Reported

Stefano lancia `labnexus.command` dal deliverable in
`/Users/stefanofiorina/kDrive/.../SPRINT_1_KBv2/labnexus-sprint1-darwin-arm64/`:

```
labnexus.command: line 41: ./labnexus: Permission denied

ERRORE: binary 'labnexus' non trovato o crashato accanto a labnexus.command.
```

## Root cause

Il binary `labnexus` ha perso il bit di esecuzione (`+x`). NON è Gatekeeper
(messaggio diverso). Causa: la cartella è sotto **kDrive** (cloud-sync). I servizi
di sincronizzazione cloud (kDrive, Dropbox, OneDrive, Google Drive) **non
preservano i permessi POSIX** e azzerano `+x` sui file sincronizzati dopo
l'estrazione dello zip. `build-zip.sh` fa correttamente `chmod +x` a build-time
(righe 50-52), quindi lo zip è valido: il bit viene strippato a valle, in kDrive.

Problema secondario: il messaggio d'errore del launcher confondeva "non trovato"
/ "crashato" / "permission denied" in un'unica stringa fuorviante.

## Fix applicato

1. `scripts/labnexus.command`: prima di `./labnexus`, se il file esiste ma non è
   eseguibile, `chmod +x ./labnexus` (auto-heal). Messaggi distinti per i tre
   casi. Resta il chicken-egg sul `.command` stesso (se perde +x il doppio click
   non parte affatto) → documentato.
2. `README.md`: sezione "Cartelle su cloud-sync" con `chmod +x labnexus labnexus.command`.

## Workaround immediato (dato a Stefano)

```bash
cd "/Users/stefanofiorina/kDrive/Common documents/R&D/Cartella CONDIVISA/SPRINT_1_KBv2/labnexus-sprint1-darwin-arm64"
chmod +x labnexus labnexus.command
```

## Follow-up

- Il deliverable zip va **ricostruito** (build-zip.sh) e ri-consegnato con il
  launcher aggiornato perché Denis/Stefano abbiano l'auto-heal.
- Valutare se conviene sconsigliare l'esecuzione diretta da cartella cloud-sync
  (copiare in locale prima di lanciare) — nota operativa per Stefano.
