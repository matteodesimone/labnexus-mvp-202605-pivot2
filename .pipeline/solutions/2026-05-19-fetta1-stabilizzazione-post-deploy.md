---
date: 2026-05-19
pipeline: bugfix (cycle accorpato)
feature: Sprint 1 LabNexus — chiusura definitiva Fetta 1 dopo 6 bug + 1 UX patch
stack: go-cli-tool
tags: [llm-streaming, ux, observability, ollama, reasoning-model, prove-out, manual-scenario, multi-llm-review]
---

# Fetta 1 — stabilizzazione post-deploy (6 bug + UX progress)

## Problem

Dopo aver chiuso `/v-deploy` Fetta 1 al mattino (2026-05-19 ~09:30), il primo smoke test CTO reale (qwen3.6 36B + Ollama + Test 1 con 3 file reali, 75k token input) ha rivelato 6 bug consecutivi e 1 gap UX critico:

1. **app-bundle-doppio-click-no-output**: `.app` doppio click non apriva TUI (CFBundleExecutable puntava al binario → no TTY → ErrNonTTY silent exit). Fix: wrapper bash + osascript.
2. **path-resolution-profili-kb-cwd-relative**: `./profili` risolveva contro cwd Finder (root /), non contro il bundle. Fix: `internal/paths.DeliveryRootForBinary` con detection `.app/Contents/MacOS/`.
3. **ollama-http-timeout-troppo-stretto**: 90s default insufficiente per warmup qwen3.6 36B. Fix temporaneo: default 1800s + env override `LABNEXUS_HTTP_TIMEOUT`.
4. **eof-post-content-falsamente-interrotto**: Ollama emette tutti i token + chiude TCP senza chunk `done:true` finale → drainStream marca "interrotto" su output de facto completo. Fix: `NoDoneMarker` flag → "completato" + warning.
5. **http-overall-timeout-tronca-streaming**: post-fix-3, `Client.Timeout=30min` overall tronca il body streaming a metà output. Fix: refactor a `Transport.ResponseHeaderTimeout` (TTFB only) + `Client.Timeout=0` (body illimitato).
6. **tui-drag-drop-trailing-space**: drag&drop da Finder a Terminal incolla path con `'/path' ` (single-quote + trailing space). `os.Stat` fallisce sulla stringa raw. Fix: `cleanPath()` helper applicato in `validateExistingDir`/`validateNonEmpty`/post-form-Run/preInput.

E:
7. **UX progress indicator** (non bug, gap): senza feedback live, il CTO non distingue "modello in warmup/reasoning" da "freeze". Quasi killato un run valido di 1h alle ~17:30. Fix: `runlog.StreamProgress` + `StreamEnd` + `drainStream` refactor con 3 ticker (wait 10s pre-firstToken / progress 500ms su TTY / nonTTY 30s).

Risultato finale: il binario ha completato Test 1 reale end-to-end in **3597s (59m 57s)**, output integro (`stato: completato`, frontmatter pulito, contenuto non troncato).

## Solution

Cycle accorpato `/v-review` + `/v-deploy` + `/v-compound` unico su bugfix-5 + bugfix-6 + UX-progress (i bug 1-4 erano già chiusi nei cicli precedenti del giorno).

Review multi-LLM (codex + mistral, ollama down) ha rivelato 26 findings di cui:
- 2 in-scope al cycle → fixati inline (cleanPath su preInput M8, TTFB→TTFT L1)
- 5 pre-existing severi → aperti come bug files per backlog
- 3 false positive → dismissed con razionale tecnico
- ~16 generici codebase audit → annotati per `/v-triage` futuro

## What Worked

- **Refactor del client HTTP per LLM streaming**: il pattern `Transport.ResponseHeaderTimeout` + `Client.Timeout=0` è ESATTAMENTE la separazione semantica giusta (TTFB capped, body illimitato). Da memorizzare come pattern Go canonico per qualsiasi client LLM streaming/SSE.
- **TDD red→green strict su ogni bug**: il test che riproduce il bug PRIMA del fix è non negoziabile. Bug-5 (HTTP overall timeout) ha richiesto un fake server con sleep 600ms × 4 chunk per riprodurre il pattern reale del troncamento.
- **Progress indicator a 3 ticker**: il design "wait pre-firstToken + progress live + nonTTY periodico" è elegantemente esclusivo (le 3 fasi non si sovrappongono). Mistral l'ha contestato come interleave possibile ma è false positive.
- **Bug files come tracciabilità anche su fix-diretti**: avere `.pipeline/bugs/*.md` per ogni bug (anche quando il fix è 5 minuti) crea un audit trail eccellente. Il `/v-compound` può fare riferimento a 11 bug files per Fetta 1 senza ricostruire dalla memoria.
- **Multi-LLM review**: codex ha trovato i finding privacy più severi (eurouter gate, endpoint override, symlinks). Mistral ha trovato i finding qualità (TTFT vs TTFB, validation gaps). Coverage complementare ben oltre quello che Claude da solo avrebbe trovato.

## What Didn't Work

- **`@manual` BDD scenarios MAI eseguiti pre-deploy**: tutti i 6 bug + gap UX erano negli `@manual` scenarios (drag&drop, doppio-click, run reale con LLM). Il `/v-test-scaffold` di Fetta 1 li aveva creati ma marcati `@manual` (cioè "fuori from CI"). Nessun checkpoint formale obbliga il CTO a eseguirli pre-deploy → tutti scoperti AL primo smoke test reale, dopo che il binario era già "shipped".
- **Stima warmup qwen3.6**: piano iniziale prevedeva qwen3 14B con TTFB ~10s. Realtà: qwen3.6 36B MoE 23GB → TTFB 30-60s + reasoning lunghissimo (Test 1 = 1h netta). Mai validato in spec quanto sarebbe stato lento.
- **Cycle pipeline accorpato non standard**: ho fatto 3 fix di natura diversa (bug-5 server, bug-6 client UX, UX-progress feature) in un singolo "review/deploy/compound" perché il flow CTO era in flow. Tecnicamente la disciplina v-bugfix prevede 1 ciclo per bug. Per prove-out questa accorpatura è ACCETTABILE (budget-free + velocity > formalismo) ma da nominare esplicitamente.
- **REVIEW=all in cycle accorpato**: l'intent "until-clean" è impraticabile su review multi-LLM cross-file (genera findings infiniti su pre-existing codebase). Triage manuale → fix in-scope + deferred pre-existing è obbligatorio.
- **Diagnosi "blocco" prematura**: alle ~17:30, dopo 50min di silenzio, ho diagnosticato "labnexus bloccato in drainStream". Era falso: Ollama stava ancora completando il reasoning prima di emettere `done:true`. Lezione: con modelli reasoning, idle apparente può durare diversi minuti tra ultimo token visibile e done finale. Senza progress indicator, indistinguibile da bug.

## Reusable Pattern

### Pattern 1: HTTP client per LLM streaming (Go)

```go
// Mai usare Client.Timeout overall per streaming LLM:
// il body può durare ore senza essere "appeso".
// Usa Transport.ResponseHeaderTimeout per il SOLO TTFB.

func StreamingHTTPClient(headerTimeout time.Duration) *http.Client {
    transport := http.DefaultTransport.(*http.Transport).Clone()
    transport.ResponseHeaderTimeout = headerTimeout
    return &http.Client{
        Transport: transport,
        // Timeout: 0 = no overall timeout (body può streaming all'infinito)
    }
}
```

### Pattern 2: Progress indicator multi-ticker per long-running stream

Tre fasi mutuamente esclusive:
- `waitTicker` (10s): pre-firstToken, ping "ancora in attesa..."
- `progressTicker` (500ms): post-firstToken, su TTY, `\r` in-place update
- `nonTTYTicker` (30s): post-firstToken, su non-TTY, linea log periodica

Stesso `select { case <-X: ... case ev := <-streamCh: ... }` per orchestrare.

### Pattern 3: Path cleaning per TUI cross-platform

`cleanPath(s) = TrimSpace → Trim(quote) → TrimSpace`. Idempotente, simmetrica. Va applicata in OGNI validator + sul preInput passato come argomento → drag&drop Finder genera `'/path/' ` (quote + space) e questo pattern lo neutralizza ovunque venga.

## System Updates Applied

- Aggiornato `state.json` con cycle accorpato e backlog_opened_during_review
- Aperti 5 bug files in `.pipeline/bugs/` per pre-existing findings di review
- USER-GUIDE.md: nuova sezione "Cosa vedo durante l'esecuzione" + nota TTFT
- Bundle .app ricompilato in `dist/labnexus-sprint1-darwin-arm64/` (data 19:11)

## Tech Debt

5 bug files aperti dal review esterno per backlog Fetta 2/3:
1. `privacy-eurouter-gate-mancante.md` — CRITICAL (eurouter senza gate approvazione)
2. `privacy-paths-personali-in-log.md` — MEDIUM (PII in stderr)
3. `ollama-endpoint-env-override-leak.md` — HIGH (`LABNEXUS_OLLAMA_ENDPOINT` può deviare locale → remoto)
4. `kb-symlink-escape.md` — HIGH (KB symlinks escape path traversal)
5. `nodonemarker-semantica-completato.md` — MEDIUM (design decision di bug-4, da rivisitare con Denis)

**Raccomandazione**: PRIMA di esporre il binario a un utente esterno (anche solo Denis), chiudere ALMENO i 3 HIGH/CRITICAL privacy/security:
- privacy-eurouter-gate (CRITICAL): rischio reputazione/GDPR
- ollama-endpoint-env-override (HIGH): bypass silenzioso garanzia "locale"
- kb-symlink-escape (HIGH): exfil arbitrari file via KB malevola

Ulteriori findings out-of-scope (codice non toccato dal cycle): `classifyError` fragile, `runner.Config` validation incomplete, error message capitalization → triage in `/v-triage` post-Sprint.

**Consiglio**: prima di iniziare Fetta 2, considerare un `/v-refactor` mirato sui 3 finding HIGH/CRITICAL (~3-4h totali di lavoro). Alternativa: rimandare a Fetta 2 e affrontare insieme alle nuove capability.

## Pattern proposto al framework

Vedi `.pipeline/proposed-updates/2026-05-19-manual-scenarios-deploy-checkpoint.md` per la promotion del pattern "`@manual` scenarios obbligatori pre-deploy CTO" che avrebbe prevenuto **tutti** i 6 bug + UX gap di stamattina.

Aggiungo: `.pipeline/proposed-updates/2026-05-19-review-accorpata-prove-out.md` per il pattern "review accorpata multi-bugfix in prove-out: triage pre-existing → bug files, no infinite loop".

E: `.pipeline/proposed-updates/2026-05-19-progress-indicator-llm-streaming.md` per il pattern "progress indicator obbligatorio per CLI con LLM streaming long-running".
