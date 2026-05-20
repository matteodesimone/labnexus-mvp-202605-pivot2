---
id: 003
priority: P1
status: done
source: bug
created: 2026-05-20
file: internal/provider/provider.go
bug_ref: .pipeline/bugs/ollama-endpoint-env-override-leak.md
---

# `LABNEXUS_OLLAMA_ENDPOINT` può deviare il "provider locale" a host remoto silenziosamente

## Issue

L'env var `LABNEXUS_OLLAMA_ENDPOINT` (utile per i test BDD via fake server `httptest`) è accettata anche in build di produzione senza restringerla a loopback hosts. Un valore come `LABNEXUS_OLLAMA_ENDPOINT=https://ollama-cloud-evil.example.com/api` passerebbe il prompt a un endpoint remoto, ma il frontmatter dell'output continuerebbe a dichiarare `provider: ollama` (= "locale, on-device" nella mente di Denis). Breakaggio della core promise.

## Fix

In `internal/provider/ollama.go` (o equivalente):
1. Validare l'endpoint risolto: SE non è `127.0.0.1` / `::1` / `localhost`, rifiutare con errore esplicito.
2. Eccezione: variabile env DEDICATA per i test (`LABNEXUS_TEST_OLLAMA_ENDPOINT`) usata solo dai test BDD, MAI in produzione (build flag o tag).
3. Aggiornare `features/helpers_test.go` per usare la env dedicata.

Alternative: gate dietro lo stesso meccanismo di `LABNEXUS_ALLOW_CLOUD_PROVIDER` del todo #001 (eurouter).

Test:
- Unit: `provider.parseOllamaEndpoint("https://remote.example.com")` → error.
- Unit: `provider.parseOllamaEndpoint("http://127.0.0.1:11434")` → ok.
- BDD: lancio con env `LABNEXUS_OLLAMA_ENDPOINT=https://remote` → exit 2.

## Context

- Confermato da Codex L2 review Fetta 2.
- I test BDD attualmente sfruttano questa env (helpers_test.go avvia fake server e setta `LABNEXUS_OLLAMA_ENDPOINT`). Il fix deve preservare quel meccanismo via env dedicata, NON via lo stesso `LABNEXUS_OLLAMA_ENDPOINT` permissivo.
- Convertito in todo P1 da `/v-triage` 2026-05-20.
- Stima fix: 2 gettoni infra (~1-2h CTO + refactor test helper).
