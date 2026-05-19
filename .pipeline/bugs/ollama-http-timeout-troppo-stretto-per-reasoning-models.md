---
severity: high
status: fixed
created: 2026-05-19
fixed: 2026-05-19
source: developer (CTO smoke test post-bugfix-2 path-resolution)
fix: "internal/provider/ollama.go: nuovo OllamaDefaultTimeout(p) che legge env LABNEXUS_HTTP_TIMEOUT (in secondi) con fallback robusto a 1800s (30 min). Usato in doRequest quando p.HTTPClient è nil. Coerente con: qwen3.6 36B MoE warmup ~30s + context fill 84k token ~minuti + thinking + decode 8192 token = totale realistico 15-30 min. docs/USER-GUIDE.md aggiornato con troubleshooting (curl /api/tags + env override). EUrouter timeout invariato (120s, sufficiente per provider cloud)."
test: "internal/provider/ollama_timeout_test.go (5 test): TestOllamaProvider_TolerateSlowResponseHeaders (fake server con delay 2s, client con timeout 1s fail / 10s pass — verifica behavior), TestOllamaProvider_DefaultTimeoutSupportsLongPrompts (default ≥ 30min), TestOllamaProvider_HonorsTimeoutEnv (env LABNEXUS_HTTP_TIMEOUT=60 → 60s), TestOllamaProvider_HonorsTimeoutEnv_Invalid (env invalido → fallback default), TestOllamaProvider_TimeoutErrorMentionsCorrectCause."
---

# Bug: timeout HTTP 90s è troppo stretto per modelli reasoning grandi (qwen3.6 36B MoE)

## Reported

- **Feature/Area**: FR-7 (provider Ollama streaming) — `internal/provider/ollama.go`
- **Trovato da**: CTO al primo smoke test reale post-bugfix path-resolution (2026-05-19)
- **Sprint impact**: HIGH. Denis non può completare uno shakedown del Test 1: il modello qwen3.6 36B MoE su prompt da 84k token impiega 15-30 min totali, il client abort dopo 90s con "context deadline exceeded".

## Description

**Cosa succede** (verificato sul Mac CTO con Ollama attivo + qwen3.6:latest 23 GB):
- Il binario carica profilo + KB + parsa input + compone prompt (84749 token stimati) — OK
- Apre POST a `http://localhost:11434/api/chat`
- Aspetta 90s per la response → `Client.Timeout exceeded while awaiting headers`
- Exit code 1 con messaggio "Ollama non raggiungibile su localhost:11434 — verifica che ollama serve sia attivo"

**Cosa dovrebbe succedere**:
- Aspettare il tempo necessario (anche 30 min) perché il modello completi context fill + thinking + decode
- Lo streaming Ollama produce token via NDJSON-streaming: il primo token può arrivare dopo molti minuti su prompt large, ma una volta partito lo stream produce token incrementali → il timeout HTTP totale (request_timeout) è troppo restrittivo

## Steps to Reproduce

```bash
# Sul Mac dell'utente con Ollama attivo + qwen3.6:latest installato:
'/path/labnexus.app/Contents/MacOS/labnexus-bin'

# Il binario apre la TUI, l'utente sceglie:
#   - revisione (Capability A)
#   - input: cartella Test 1 (PG_RISK_LAB_Rev_00.docx + 2 PDF RT-08)
#   - output: cartella qualsiasi
# 
# Risultato osservato:
#   [HH:MM:SS] chiamata provider ollama ...
#   [HH:MM+1:30] chiamata provider ollama ok (90026 ms)
#   Error: Ollama non raggiungibile su localhost:11434 — verifica che ollama serve sia attivo:
#     Post "http://localhost:11434/api/chat": context deadline exceeded
#     (Client.Timeout exceeded while awaiting headers).
```

Diagnosi ambiente (eseguita dal CTO):
- `curl http://localhost:11434/api/tags` → Ollama RISPONDE, lista modelli OK
- `ollama list | grep qwen` → `qwen3.6:latest 23 GB`
- `time ollama run qwen3.6 "ciao"` → **32 secondi** per UNA parola (modello reasoning con `<think>` block)

Quindi Ollama è OK; il modello fa thinking pesante; il timeout client HTTP è semplicemente troppo basso per la combinazione "modello 36B + prompt 84k token + reasoning model".

## Environment

- macOS Apple Silicon
- Ollama in esecuzione su `localhost:11434`
- Modello: `qwen3.6:latest` (Qwen 3.5 MoE 36B params, Q4_K_M, 23 GB su disco)
- Profilo: `revisione` (context_window 128000, max_tokens 8192, temperature 0.9)
- Input: Test 1 reale = PG_RISK_LAB_Rev_00 + 2 PDF RT-08 = ~84749 token

Setup del client problematico:
```go
// internal/provider/ollama.go (post-bugfix-2):
client = &http.Client{Timeout: 90 * time.Second}
// 90s è il timeout TOTALE della request, incluso awaiting headers.
```

## Analysis

**Likely cause (alta confidenza)**:

`http.Client.Timeout` di 90s era stato introdotto nel bugfix `app-bundle-doppio-click-no-output` (loop unmasked failure) per evitare hang lunghi sul caso "Ollama down" via macOS port 1. Allora 90s sembrava generoso. Per modelli reasoning grandi con prompt grandi è invece STRETTO.

Considerazioni:
1. **Il timeout HTTP misura il TIME-TO-FIRST-BYTE**, non il tempo totale dello stream. Per LLM streaming, il TTFB può legittimamente essere minuti (context fill).
2. **Qwen 3.6 36B MoE** sul Mac Apple Silicon: ~30s warmup + ~2-5 min context fill (84k token) + thinking interno + decode → primo HTTP-header response può venire dopo 5-10 min.
3. **Mentre il timeout 90s previene hang infiniti** su endpoint sbagliati (port 1), il caso d'uso "vero Ollama, modello vero, prompt vero" richiede minuti.

**Soluzione (idiomatic Go HTTP)**: il `http.Client.Timeout` è il timeout totale; quello che vogliamo è:
- `Timeout` molto generoso (es. 30 min) — protezione di sicurezza
- OPPURE separare `Transport.ResponseHeaderTimeout` (TTFB, es. 10 min) + nessun timeout totale (lo streaming deve poter durare)
- Override via env `LABNEXUS_HTTP_TIMEOUT` per casi estremi (modelli ancora più grandi, prompt enormi)

**Affected files**:
- `internal/provider/ollama.go` (`OllamaProvider.doRequest` → client creation)
- `internal/provider/eurouter.go` (idem, anche se EUrouter è cloud e tipicamente più rapido)
- (Documentation) `docs/USER-GUIDE.md` aggiungere troubleshooting "primo lancio impiega più di X minuti"

**Risk of fix**: **basso**. Cambio della costante di timeout + lettura env var opzionale.

## Note di scope (prove-out)

Bug fix budget-free. Stima fix: ~30-45 min (cambio timeout + env override + 1 test integration con fake server che ritarda > 90s + smoke reale CTO + nota troubleshooting in USER-GUIDE).

## Domanda separata emersa dal report CTO

> "il sistema legge i pdf sul serio o in realtà devo io convertire le cose in md o meno"

**Risposta tecnica** (NON un bug, è feature documentata):

Sì, il sistema legge i PDF veri via `ledongthuc/pdf` (pure-Go). Estrae il testo embedded:
- ✅ PDF generati digitalmente (Word→PDF, LaTeX→PDF, browser print) — funziona
- ⚠ PDF complessi (multi-colonna fitto, formule LaTeX, tabelle annidate) — estrazione imperfetta
- ❌ PDF scansionati senza OCR — il binario NON fa OCR, estrae solo testo embedded vuoto

Per i tuoi Test 1 PDF (`PG_RISK_LAB`, `RT-08`): il dry-run estrae ~84749 token, quindi il PDF è stato letto correttamente. Non serve convertire manualmente.

Il warning ".DS_Store" che vedi nel log è un metadata file macOS, normale skipparlo come "formato non supportato".

**Quando convertire manualmente in `.md`?** Solo se:
- Vedi nei log "PDF non parsabile (magic-bytes mancanti)" su un file specifico
- L'output del modello sembra confuso vs il contenuto del PDF originale (= testo estratto male)

→ Da aggiungere come voce troubleshooting in `docs/USER-GUIDE.md` durante questo bugfix.
