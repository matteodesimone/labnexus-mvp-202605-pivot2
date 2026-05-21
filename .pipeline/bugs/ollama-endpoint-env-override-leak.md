---
severity: high
status: intentional_deviation_post_pivot_3
created: 2026-05-19
fixed_at: 2026-05-20
deviation_at: 2026-05-21
source: review esterna (codex-security)
fix: provider.Select valida l'endpoint Ollama tramite requireLoopbackOrCloudGate. Solo host loopback (127.0.0.0/8, ::1, "localhost") sono ammessi di default. Host remoti richiedono LABNEXUS_ALLOW_CLOUD_PROVIDER esplicito (stesso gate di #001 eurouter). Default "http://localhost:11434" preservato.
deviation_rationale: |
  Sprint 1.5 pivot-3 ha rimosso il gate cloud (vedi privacy-eurouter-gate-mancante.md
  intentional_deviation). La protezione `requireLoopbackOrCloudGate` per LABNEXUS_OLLAMA_ENDPOINT
  era coerente con il gate eurouter: stessa filosofia "decisione cloud/local via env var
  esplicita". Ora la decisione passa esclusivamente a configurazione utente
  (labnexus.config.toml master in 1.5.B). Eurouter è già il default; un Ollama endpoint
  remoto è uso esperto (CTO + futuro server backend Sprint 2), non ostile.
test: internal/provider/select_test.go::TestSelect_OllamaRejectsRemoteEndpoint (+5 varianti) — RIMOSSI in Sub-fetta 1.5.A (decommissione del gate)
---

# Bug: `LABNEXUS_OLLAMA_ENDPOINT` può deviare il "provider locale" a un host remoto silenziosamente

## Reported

- **Feature/Area**: `internal/provider/ollama.go:OllamaEndpoint` + env override
- **Trovato da**: review multi-LLM 2026-05-19 (codex-security HIGH)
- **Sprint impact**: HIGH per privacy. Pre-esistente.

## Description

`OllamaProvider` legge `LABNEXUS_OLLAMA_ENDPOINT` per override dell'endpoint. Originariamente introdotto per i test (BDD scenari con `httptest.Server`). Ma il binario di produzione legge la stessa env var → un attaccante che possa settarla (o un cron malevolo, o l'utente per errore) può deviare i dati SGQ a un endpoint remoto MA il frontmatter dell'output continua a riportare `provider: ollama` (cioè "locale").

Garanzia di Sprint 1 violata: `provider: ollama` deve = on-device.

## Fix proposto

Build-tag o init-check: l'env var `LABNEXUS_OLLAMA_ENDPOINT` è onorata solo se host = `localhost` / `127.0.0.1` / `::1`. Qualunque altro host viene rifiutato con error: `"ollama: endpoint remoto rifiutato per provider 'ollama' (privacy guarantee). Usare provider 'eurouter' se necessario remoto."`

Alternativamente: rimuovere completamente l'env var dal binario di produzione (compile-tag) e usare dependency injection solo nei test.

## Risk of fix

**Basso**. ~30 min. Aggiungere test BDD per "endpoint non-loopback rifiutato".

## Note

Deferred to backlog: pre-esistente. Correlato a `privacy-eurouter-gate-mancante.md` (stessa categoria: rendere espliciti gli usi remoti).
