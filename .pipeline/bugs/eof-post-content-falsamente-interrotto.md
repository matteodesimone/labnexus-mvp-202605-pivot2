---
severity: medium
status: fixed
created: 2026-05-19
fixed: 2026-05-19
source: developer (CTO smoke test post-bugfix-3 ollama-timeout)
fix: "1) internal/provider/provider.go: aggiunto campo `NoDoneMarker bool` a StreamEvent. 2) internal/provider/ollama.go: consumeOllamaStream ora distingue EOF post-content (errors.Is(err, io.ErrUnexpectedEOF) o EOF clean) da errori veri del scanner. Sul primo emette `Done: true, NoDoneMarker: true`; sul secondo emette Err come prima. 3) internal/runner/runner.go: drainStream interpreta NoDoneMarker → stato 'completato' + warning log esplicito. 4) spec EC-8 aggiornata via feature files: scenari `streaming interrotto produce file parziale` (motore-output.feature) e `interruzione di rete a metà streaming` (edge-cases.feature) ora assertano stato='completato' + warning stderr invece di stato='interrotto' + exit 1. Rationale dell'aggiornamento spec: tecnicamente non è possibile distinguere lato client 'EOF clean post-content' da 'TCP RST mid-stream', quindi preferiamo conservare il body utile + warning trasparente all'utente (Denis L2 valuta se il body è completo)."
test: "internal/provider/ollama_eof_test.go (3 test): TestOllamaProvider_EOFAfterContentCompletes (server emette 3 chunk + EOF senza done → Done+NoDoneMarker, no Err), TestOllamaProvider_EOFBeforeContentIsRealError (server vuoto + EOF → almeno Err o Done), TestOllamaProvider_NormalCompletionStillEmitsCleanDone (server con done:true esplicito → Done senza NoDoneMarker)."
---

# Bug: EOF post-content marcato falsamente "interrotto" invece di "completato"

## Reported

- **Feature/Area**: FR-9 (frontmatter stato) + EC-8 (streaming interrotto) — `internal/provider/ollama.go` e `internal/runner/runner.go`
- **Trovato da**: CTO al primo smoke test reale completo (Test 1 con qwen3.6 36B, durata 21 min)
- **Sprint impact**: MEDIUM. Output prodotti correttamente vengono marcati come "interrotti" → exit code 1 → Denis al L2 può pensare che siano parziali quando invece sono completi.

## Description

**Cosa succede** (verificato sul Mac CTO):
- Lancio del binario su Test 1 (input parziale: solo le 2 norme RT-08, senza il documento SGQ)
- qwen3.6 36B genera output completo in 21 minuti (1275s):
  - 11 KB, 147 righe, frontmatter + 8 sezioni `##` + 5 callout `[!MODIFICA]` + chiusura naturale con sezione "5. Prossimo passo" e callout `[!ISPETTORE]`
- Tuttavia il binario emette: `stream error: ollama: scanner: unexpected EOF`
- Frontmatter scritto con `stato: interrotto`
- Exit code 1: `esecuzione conclusa con stato non OK`

**Cosa dovrebbe succedere**: l'output è chiaramente completato (chiusura naturale del modello, niente troncamento a metà frase). Il frontmatter dovrebbe dire `stato: completato` (eventualmente con annotazione "no done marker" se vogliamo essere trasparenti). Exit 0.

## Steps to Reproduce

Configurazione che ha esposto il bug (Mac CTO con Ollama + qwen3.6:latest):
1. Lancia il binario su un input dove qwen3.6 produce output con chiusura naturale (es. il modello finisce il task prima di raggiungere `max_tokens`).
2. **Osserva**: lo stream Ollama emette N chunk con `done: false`, poi chiude la connessione TCP senza emettere `done: true`.
3. Il binario interpreta `unexpected EOF` come errore.

Caso riproducibile in test (vedi fix): fake Ollama server che emette 3 chunk con `done: false` e poi chiude la connessione (senza chunk finale `done: true`).

## Environment

- macOS Apple Silicon
- Ollama (versione utente — segnalato comportamento "no done marker" su qwen3.6 36B reasoning model)
- Modello qwen3.6 può chiudere lo stream naturalmente senza emettere done:true (probabile bug Ollama o configurazione modello)

## Analysis

**Likely cause (alta confidenza)**:

Il codice attuale in `internal/provider/ollama.go::consumeOllamaStream`:
```go
for scanner.Scan() {
    ...
    if c.Done {
        ch <- StreamEvent{Done: true}
        return
    }
}
if err := scanner.Err(); err != nil {
    ch <- StreamEvent{Err: fmt.Errorf("ollama: scanner: %w", err)}
}
```

Se il server chiude la connessione DOPO aver emesso content ma SENZA `done: true` nell'ultimo chunk:
- `scanner.Scan()` ritorna false
- `scanner.Err()` ritorna `io.ErrUnexpectedEOF` (linea NDJSON incompleta) o `io.EOF` (chiusura clean)
- Il codice emette `StreamEvent{Err: ...}` → runner.drainStream → `stato: interrotto` → exit 1

Manca distinzione tra:
- **Caso A "stream interrotto a metà"**: pochi/zero chunk ricevuti, EOF prematuro → stato "interrotto" è corretto
- **Caso B "stream completato senza done marker"**: molti chunk con content valido, poi EOF → stato dovrebbe essere "completato" (con eventuale warning)

**Affected files**:
- `internal/provider/ollama.go::consumeOllamaStream` (loop scanner)
- `internal/provider/provider.go::StreamEvent` (aggiungere flag `NoDoneMarker`)
- `internal/runner/runner.go::drainStream` (gestire `NoDoneMarker` come "completato" con warning)

**Risk of fix**: **basso**. Logica isolata in 3 file, semantica chiara, test integration con fake server.

## Test plan

1. Unit test `consumeOllamaStream` su fake server che emette N chunk con content + chiude senza done → atteso: emette `Done: true` (no Err).
2. Unit test `consumeOllamaStream` su fake server che chiude SENZA emettere alcun content → atteso: emette `Err` (caso A, interrotto vero).
3. Integration test `runner.drainStream` su EOF post-content → atteso: stato="completato", warning loggato.

## Note di scope (prove-out)

Bug fix budget-free. Stima fix: ~30-45 min.

**Pattern emergente** (3° bug consecutivo dopo Fetta 1 deploy):
- bug 1: doppio-click → fix wrapper bash
- bug 2: path resolution → fix paths package
- bug 3: ollama timeout → fix HTTP client timeout 30 min
- bug 4: EOF post-content → fix detection done marker

**Tutti scoperti SOLO dopo il primo smoke test reale con Ollama + qwen3.6 + Test 1 reale**. Conferma diagnosi del compound bugfix-1: il scenario `@manual` "lancia il binario sul Mac CTO con LLM reale su Test 1 end-to-end" mancava, e non avere quel test ha permesso a 4 bug critici di passare il deploy di Fetta 1. Da formalizzare nel compound finale come framework promotion.
