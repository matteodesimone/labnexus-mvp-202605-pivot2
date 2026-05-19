---
severity: high
status: fixed
created: 2026-05-19
fixed: 2026-05-19
source: developer (CTO 2° smoke test reale, Test 1 completo con 3 file)
fix: "Refactor del default HTTP client in OllamaProvider. Nuova funzione OllamaDefaultHTTPClient(p) ritorna `&http.Client{Transport: &http.Transport{ResponseHeaderTimeout: OllamaDefaultTimeout(p)}}` SENZA Client.Timeout. Semantica: TTFB (response headers) ha timeout, body streaming è ILLIMITATO. Il body può durare ore (output lunghi su modelli grandi) senza essere cancellato. La protezione contro server hung resta tramite ResponseHeaderTimeout. doRequest aggiornato per usare la nuova funzione."
test: "internal/provider/ollama_streaming_unlimited_test.go (3 test): TestBodyCanStreamLongerThanHeaderTimeout (body con sleep 600ms × 4 = 2.4s totali, default client completa OK), TestDefaultClientHasNoOverallTimeout (verifica DESIGN: Client.Timeout=0 + Transport.ResponseHeaderTimeout ≥ 30min), TestHeaderTimeoutStillProtectsAgainstDeadServer (server hung 3s di header, client custom con ResponseHeaderTimeout=1s cancella → protezione preservata)."
---

# Bug: `http.Client.Timeout` overall tronca lo streaming LLM lungo (output incompleto)

## Reported

- **Feature/Area**: FR-7 (provider Ollama streaming) — `internal/provider/ollama.go`
- **Trovato da**: CTO 2° smoke test reale (Test 1 completo con `PG_RISK_LAB_Rev_00.docx` + 2 PDF RT-08), 1800015 ms = 30 min esatti
- **Sprint impact**: HIGH. Su hardware più lento di quanto ipotizzato nel bugfix-3, anche il timeout 30 min default non è sufficiente. Aumentare ulteriormente è palliativo: il design del client HTTP è sbagliato per streaming LLM.

## Description

**Cosa succede**:
- Lancio binario su Test 1 completo (3 file, 75150 token in)
- Qwen 3.6 36B inizia a generare
- A 1800.015s (= 30 min esatti, = `Client.Timeout` default post-bugfix-3) il client cancella:
  ```
  warning: stream error: ollama: scanner: context deadline exceeded
  (Client.Timeout or context cancellation while reading body)
  ```
- Output file scritto MA effettivamente troncato a metà:
  ```
  > [!ISPETTORE] La transizione da rev. 03 a rev. 05 non è un aggiornamento
  terminologico. È un riaccomplamento strutturale completo verso ISO/IEC 1
  ```
  ← troncato a metà parola, ISO/IEC 17025 mai completato

**Cosa dovrebbe succedere**: lo streaming può durare quanto serve (1 ora, 2 ore, finché Ollama emette token). Il client HTTP NON deve cancellare durante body reading. Solo se il **primo header** non arriva (= Ollama davvero down) deve abort.

## Steps to Reproduce

1. Hardware dove qwen3.6 36B + 75k token input richiede > 30 min totali
2. Lancia binario su Test 1 completo
3. Aspetta 30 min: il client cancella + output troncato

In test: fake server che emette headers immediati + sleep 2s + emette content. Client con `Client.Timeout=1s` cancella mid-body. Riproduce lo stesso pattern del bug reale.

## Environment

- macOS Apple Silicon (utente)
- Ollama + qwen3.6:latest 23 GB 36B MoE
- Prompt 75150 token = Test 1 completo (3 file: documento SGQ + 2 norme RT-08)
- Output atteso: ~15-30 KB di markdown (bozza di revisione + frontmatter + 13+ sezioni)
- Tempo realistico sul Mac utente: > 30 min (non ho stima precisa, ma evidentemente > default attuale)

Setup problematico (post-bugfix-3):
```go
client = &http.Client{Timeout: 30 * time.Minute}
// Client.Timeout copre TUTTO il request lifecycle, inclusi i body reads.
// Per streaming LLM dove il body può durare arbitrariamente, è sbagliato.
```

## Analysis

**Likely cause (alta confidenza)**:

`http.Client.Timeout` in Go copre l'**intero ciclo di vita della request**:
- TCP connect
- TLS handshake (se HTTPS)
- Write request body
- Read response headers (TTFB)
- **Read response body** ← ANCHE QUESTO

Per LLM streaming, il body reading è l'attività principale e può durare arbitrariamente (token-by-token). Mettere un timeout overall significa cancellare lo streaming a metà, che è quello che è successo.

**Pattern Go idiomatic per streaming HTTP**:

```go
// Transport con ResponseHeaderTimeout SOLO sul TTFB
transport := http.DefaultTransport.(*http.Transport).Clone()
transport.ResponseHeaderTimeout = 30 * time.Minute  // TTFB max
client := &http.Client{
    Transport: transport,
    Timeout:   0,  // NESSUN timeout overall — body può streaming all'infinito
}
```

Semantic:
- Se Ollama davvero down → header non arriva entro 30 min → abort rapido ✓ (mantiene la protezione bug-app-bundle-doppio-click)
- Se Ollama OK e sta generando → header arriva subito (Ollama scrive 200 OK appena riceve la request), body può durare quanto vuole ✓
- Se Ollama hangs nel mezzo (rete cade) → il caso "stream chiuso senza done marker" (bugfix-4) lo intercetta tramite `io.ErrUnexpectedEOF` ✓

**Affected files**:
- `internal/provider/ollama.go` (refactor client builder)
- `internal/provider/eurouter.go` (per coerenza — anche EUrouter è streaming SSE)
- `internal/provider/ollama_timeout_test.go` (test esistenti restano validi: `Client.Timeout` continua a funzionare quando l'utente lo setta esplicitamente come HTTPClient)
- NEW test: `internal/provider/ollama_streaming_unlimited_test.go` per il pattern "header rapidi + body lungo arbitrario"

**Risk of fix**: **basso**. Refactor isolato al builder del default client. Override custom (`p.HTTPClient`) resta supportato per test.

## Test plan

Test cruciale (red phase): fake Ollama emette headers in 100ms, poi streaming lento (5s totali). Con `Client.Timeout=1s`, fallisce (bug attuale). Con `Transport.ResponseHeaderTimeout=1s` + `Client.Timeout=0`, passa. Atteso: green dopo fix.

## Note di scope (prove-out)

Bug fix budget-free. Stima fix: ~30-45 min.

**Pattern emergente — 5° bug consecutivo** post-deploy Fetta 1, tutti scoperti dai smoke test reali del CTO. Cumulativo ~3.5h CTO. **Nessuno** di questi bug sarebbe stato catturato da un test E2E unico su Mac CTO con Ollama qwen3.6 reale e prompt Test 1 reale. Da formalizzare con forza al compound finale.
