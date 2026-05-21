---
title: /v-deploy aggiunge smoke real-gateway opzionale per validare modello + auth
status: planned
source: compound learning sub-fetta 1.5.C (post-deploy HTTP 404 modello + HTTP 402 credit)
created: 2026-05-22
target: vibbly framework improvement
---

# Idea: smoke pre-deploy "real-gateway" check

## Problema

Sub-fetta 1.5.C deploy: il smoke estensivo che ho fatto (zip → tmpdir → `labnexus jobs --json`) verificava auto-discovery ma NON la chiamata vera al provider. Il primo smoke "vero" (CTO `Esegui.command` post-deploy) ha trovato:
1. HTTP 404 — modello `qwen3.6` non in catalogo eurouter
2. HTTP 402 — account a saldo €0

Entrambi rilevabili pre-deploy con una request minimale al provider (~$0.0001 di costo).

Pattern emerso più volte in Sprint 1.5 (anche 1.5.A: ollama qwen3-coder tag non trovato → richiedeva manual `ollama pull`). Il deploy step canonico vibbly NON ha un "validation against the real downstream system" — solo build + zip + verifica file in zip.

## Soluzione proposta

Estensione opzionale a `/v-deploy` skill (framework-level): aggiungi un sub-step "real-gateway smoke" che fa:

```bash
# pseudocode
real_gateway_smoke() {
  for each provider in active providers (eurouter, ollama, ecc.):
    if provider.is_configured() (API key setted, endpoint raggiungibile):
      send_minimal_request(provider, model=cfg.modello, prompt="ping", max_tokens=10)
      assert response.status == 200
      assert response.body contains some token
}
```

Costo: ~$0.0001 per provider per deploy (qualche token in/out su un prompt minimale).

Esempio output `/v-deploy` con smoke:

```
Pre-flight check:
✓ Review passed
✓ All tests passing
✓ Working tree clean
✓ Linter clean

Real-gateway smoke check:
✓ eurouter (modello: qwen3.5-122b-a10b) — 200 OK, 1247ms, 12 token
✓ ollama (modello: qwen3-coder:30b) — 200 OK, 8341ms, 8 token

Ready to deploy.
```

Su failure:

```
Real-gateway smoke check:
✗ eurouter (modello: qwen3.6) — HTTP 404
  Response: {"error": "Model 'qwen3.6' not found"}
  Suggerimento: verifica modello con `labnexus models --filter qwen`
✗ ollama (modello: qwen3-coder) — HTTP 404
  Suggerimento: lancia `ollama pull qwen3-coder` o usa un modello già installato

Deploy ABORTED. Fix i provider e rilancia /v-deploy.
```

## Trade-off

- **Pro**: catturerebbe il 100% dei 404/401/402 PRIMA della consegna a Denis. Sprint 1.5 avrebbe risparmiato 4 commit di follow-up post-deploy.
- **Pro**: Sprint 2+ qualsiasi cambio provider/modello/key viene validato automaticamente.
- **Con**: richiede network attivo a deploy time. Per ambienti offline (CI air-gapped) opt-out via flag `--skip-real-gateway-smoke`.
- **Con**: costa pochi centesimi per deploy (deploy frequenti = cumulativo).
- **Con**: scope è vibbly framework (skill `/v-deploy`), non LabNexus per-progetto. Richiede framework promotion.

## Approccio implementativo

Vibbly skill `/v-deploy` step 1 (Pre-flight) estesa con:

```
6. **Real-gateway smoke check** (opzionale, on by default):
   - Se `.pipeline/config.sh` ha `DEPLOY_SMOKE_REAL_GATEWAY=true` (default per cli-tool con providers configurati):
     - Carica provider configuration (env vars + master config se esiste)
     - Per ogni provider attivo:
       - Esegui request minimale (ping/pong)
       - Verifica HTTP 200 + body non vuoto
     - Su failure: ABORT deploy con error diagnostico
   - Skip se `--skip-real-gateway-smoke` flag presente o env `VDEPLOY_SKIP_SMOKE=1`
```

## Status

Idea per framework vibbly. Da promuovere via `.pipeline/proposed-updates/` se accettata. Quick win per LabNexus pre-Sprint 2: smoke shell manuale documentato in DEV-GUIDE.md (`curl` al modello prima del build-zip).
