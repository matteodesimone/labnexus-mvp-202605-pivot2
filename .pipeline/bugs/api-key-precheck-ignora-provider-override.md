---
severity: high
status: open
created: 2026-06-02
source: external reviewer (Codex) durante /v-review bugfix trigger-override-job
area: cmd/labnexus/main.go (preCheckAPIKeyByName, FR-36)
---

# Bug: il pre-check API key (FR-36) ignora gli override di provider CLI/env

## Description

`preCheckAPIKeyByName` (cmd/labnexus/main.go:~133) verifica la presenza della
chiave EUrouter usando **solo** il provider del profilo merged col master. Se il
provider effettivo viene cambiato a runtime con `--provider eurouter` oppure
`LABNEXUS_PROVIDER=eurouter`, il fail-fast FR-36 viene saltato: il binary procede
a caricare profilo + KB + input (potenziali dati SGQ/PII) e fallisce solo più
tardi, dentro `runner.Run`.

`provider.Select` a runtime applica la precedenza `--provider > LABNEXUS_PROVIDER >
profilo/master`, ma il pre-check no → disallineamento.

## Expected

Il pre-check deve risolvere il provider **effettivo** con la stessa precedenza di
`provider.Select` PRIMA di caricare input, e rifiutare all'avvio (exit 2) se la
chiave manca. Nessun input sensibile deve essere parsato prima del check.

## Repro (da scrivere come test)

Profilo con `provider = "ollama"`, nessuna `EUROUTER_API_KEY`, invocare
`run --job <x> --provider eurouter` → oggi parte e fallisce tardi; atteso: exit 2
immediato.

## Note

Trovato fuori dallo scope del bugfix trigger-override (file toccato ma logica
non correlata). Non smascherato dal fix → deferred a pipeline dedicata.
