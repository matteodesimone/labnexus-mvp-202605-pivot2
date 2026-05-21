---
title: vobbly preflight pre-warma automatico modello Ollama prima dello smoke test
status: planned
source: compound learning sub-fetta 1.5.A + 1.5.B (Sprint 1.5 LabNexus)
created: 2026-05-21
target: framework promotion vibbly/vobbly
---

# Idea: pre-warmup automatico modello Ollama nel `vobbly hook preflight`

## Problema osservato

In sub-fetta 1.5.A loop 1 e 1.5.B loop 1 (entrambi 2026-05-21), il preflight ha sempre fallito alla prima esecuzione per ollama:

```
"ollama": {
  "status": "timeout",
  "duration_ms": 15094,
  "message": "timeout after 15s"
}
```

Smoke test ha 15s di timeout. Il modello `qwen3-coder:30b` (18 GB) richiede 30-60s per cold-start (caricamento in RAM). Quindi ogni volta:

1. Lancio `vobbly hook preflight` → ollama timeout dopo 15s
2. Pre-warmo manualmente: `echo "ping" | timeout 120 ollama run qwen3-coder:30b`
3. Re-lancio `vobbly hook preflight` → ollama ora ok (modello caricato)

Time wasted: ~2 min per ogni invocation di `/v-review`, e il CTO/agent deve ricordare di pre-warmare. Non è un workflow naturale.

## Soluzione proposta

`vobbly hook preflight` dovrebbe **chiamare implicitamente un pre-warmup** del modello ollama prima di eseguire lo smoke test, con timeout più generoso (es. 90s) per il warmup vs 15s per lo smoke vero.

### Algoritmo proposto

```
per ogni reviewer ollama enabled:
    1. tenta uno smoke leggero (es. POST /api/show {name: modello}) per detection modello caricato
       - se risponde in <2s: modello già loaded, skip warmup, passa allo smoke vero
       - se timeout: passa a step 2
    2. lancia warmup: POST /api/generate {model, prompt: "1", stream: false}
       con timeout esteso (60-90s configurabile via env VOBBLY_OLLAMA_WARMUP_TIMEOUT)
    3. quando warmup ritorna, esegui smoke test vero come oggi (15s)
    4. esito normale (ok / timeout / unreachable)
```

### Vantaggi

- Workflow naturale: utente lancia `/v-review` → preflight prende qualche minuto la prima volta, ma poi tutto funziona.
- Niente intervento manuale (pre-warmup esplicito).
- Allinea il preflight a "production load" — uno smoke a 15s testa solo la liveness, non la capacità di rispondere a un prompt realistico.

### Trade-off

- Tempi più lunghi al primo lancio (60-90s vs 15s). Accettabile: è la trade-off di usare modelli grandi locali.
- Più codice nel preflight. Mitigare: leverage di funzione `ollama.PreWarm(model, timeout)` dedicata.
- Non aiuta progetti senza ollama (mistral/codex/gemini): no overhead.

## Implementazione suggerita (framework side)

Aggiungere a `framework/internal/reviewer/ollama.go`:

```go
// PreWarm carica il modello in RAM prima dello smoke test. Idempotente:
// se il modello è già loaded, no-op.
func (o *OllamaReviewer) PreWarm(ctx context.Context, timeout time.Duration) error {
    // 1. Check `/api/show` per detection
    // 2. Se non loaded, POST `/api/generate` con prompt minimal e wait
    // 3. Ritorna nil se ok, error se timeout
}
```

Integrazione in `RunPreflight`:

```go
if r.IsOllama() {
    warmupCtx, cancel := context.WithTimeout(ctx, 90*time.Second)
    if err := r.PreWarm(warmupCtx); err != nil {
        // PreWarm fallito: il modello probabilmente non esiste o ollama down.
        // Skip smoke test, marca exec error.
    }
    cancel()
}
// poi: smoke test normale 15s
```

## Workaround temporaneo (fino a fix framework)

Documentato nella memoria utente: prima di `/v-review`, pre-warmare manualmente con `echo "ping" | timeout 120 ollama run <modello>`. Lesson learned cross-fetta che non sarebbe necessaria con il fix proposto.

## Relazione con altre proposed-updates

- `.pipeline/proposed-updates/2026-05-21-state-json-current-fetta-int-enforcement.md` (da 1.5.A): bug framework su schema check
- Questa idea: improvement framework su UX del preflight

Entrambe escono dal singolo progetto LabNexus e migliorano vobbly/vibbly per qualsiasi progetto che usa ollama come reviewer. Candidato a proposal vibbly framework quando il CTO ha tempo.
