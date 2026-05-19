---
date: 2026-05-19
pipeline: new-feature
feature: Sprint 1 LabNexus — Fetta 1 (motore CLI + revisione + rilievi)
stack: go-cli (Go 1.22 + cobra + godog + charmbracelet/huh + ledongthuc/pdf)
tags: cli, go, llm, ollama, eurouter, prove-out, dual-provider, bdd, tui, macos-bundle, sprint-1
---

# Fetta 1 — Scooter (motore + 2 profili shakedown)

## Problem

Sprint 1 di LabNexus / AICertus: serve un eseguibile macOS/Linux che esegua una capability ispettiva alla volta del modello AICertus, chiamando Qwen 3 in locale via Ollama. Deve essere consegnabile a Denis (QM) come scooter end-to-end — usabile via CLI o doppio click su `.app`, con 2 profili shakedown (`revisione` e `rilievi`) che replicano i Test 1 e 2 manuali già approvati. Locale mandatory per il deliverable; provider EUrouter come canale di sviluppo quando il CTO è su macOS dove Ollama non gira.

## Solution

Architettura Pipeline + Strategy:
- **CLI cobra**: 5 sottocomandi (`run`/`list`/`describe`/`check`/`validate`) + TUI cross-platform via `charmbracelet/huh` quando lanciato senza argomenti.
- **LLMProvider interface** (Strategy): `OllamaProvider` (NDJSON streaming) ed `EurouterProvider` (SSE OpenAI-compatible). Selezione runtime via `provider.Select(flag>env>profile)`. Env override `LABNEXUS_OLLAMA_ENDPOINT` / `LABNEXUS_EUROUTER_ENDPOINT` per test BDD senza Ollama reale.
- **Pipeline stages** in `runner.Run`: load profile → load KB → parse input → compose prompt → estimate tokens → stream LLM → write output. Ogni stage è una funzione `≤ 22 righe`.
- **Profile schema** (FR-3): YAML con path-traversal guard sul kb_dir (NFR-6).
- **Validazione L1 + L2** (NFR-8): pre-check tecnico CTO automatizzato via BDD (frontmatter completo, struttura coerente col template, no allucinazioni macroscopiche) + giudizio qualitativo Denis manuale (fuori CI).
- **Bundle macOS .app non firmato** + binario Linux/amd64 standalone via cross-compile Go. Zip di consegna con KB-ispettore inclusa (symlink seguito).

## What Worked

- **Slicing in 3 fette equilibrate** (6+3+3) vs proposta granulare (8×1+): meno overhead di gate pipeline, checkpoint cliente naturali, "scooter ricco" che chiude lo shakedown del motore in un solo ciclo.
- **Pattern L1+L2 emerso al gate plan** ed esplicitato in `spec.md` NFR-8 + `fette.md`: ha protetto la pipeline da test BDD che pretendono di "validare Qwen" (impossibile senza Denis).
- **Strategy pattern per provider** con env override: ha permesso fake server `httptest.NewServer` nei BDD step **senza modificare il codice** del provider — boundary mock zone 2 perfettamente isolato.
- **Magic-bytes check su PDF in `extractPDF`**: ha catturato il fatto che `ledongthuc/pdf` succede silenziosamente su input non-PDF (0 pagine, no error). Senza, EC-1 sarebbe stato un falso-negativo in produzione.
- **Funzione `classifyError` in `cmd/labnexus/main.go`**: mappare semantica exit code NFR-5 (1=runtime, 2=misuse) in un solo posto. Pattern riusabile.
- **Cross-compile Go nativo** (`GOOS=darwin GOARCH=arm64`) da Linux/WSL2 a darwin/arm64: ha funzionato out-of-the-box grazie a `CGO_ENABLED=0` + librerie pure-Go (`ledongthuc/pdf`, `charmbracelet/huh`).

## What Didn't Work

- **Bug aritmetico nei test scaffold** (`tokens_test.go`: ratio dichiarato 50%/70% ma textLen e context_window non tornavano): scoperto solo durante `/v-implement` quando i test fallivano per ragione sbagliata. Lezione: in `/v-test-scaffold` verificare l'aritmetica dei boundary case, non solo la struttura del test.
- **BDD step regex senza prefisso `(?:che )?`**: il Gherkin italiano `Dato che X` produce step text `che X`, non `X`. Scoperto solo dopo la prima esecuzione godog. Costo: 1 iterazione di test+fix dei pattern.
- **Single-reviewer mode forzato dall'assenza di `vobbly`**: la review multi-LLM dialettica era una delle feature del framework che speravamo di sfruttare. Setup framework lo deve risolvere prima della prossima fetta.
- **macOS port 1 risponde HTTP 404** invece di "connection refused": il test "Ollama non disponibile" si aspettava un errore di rete ma riceveva una risposta HTTP. Fix: in `OllamaProvider.doRequest`, HTTP 4xx/5xx → `ErrOllamaUnreachable` (server raggiunto ma non è un vero Ollama). Pattern accettato come "endpoint sembra OK ma non è Ollama".
- **Working tree non committato durante l'intera fetta**: il `/v-implement` non ha fatto commit incrementali — un singolo grande commit al `/v-deploy`. Il git log della fetta è poco granulare. Decidere in Fetta 2: commit per sub-step o per gate? Probabilmente per gate (spec gate / plan gate / implement gate / deploy gate).
- **Symlink `KB-ispettore` nello zip**: `zip -ry` preservava il link → broken al primo extract di Denis. Fix tardivo (post-build). Lezione: testare lo zip su un mount diverso prima di considerarlo "consegnabile".

## Reusable Pattern

**"Dual-provider con env override per BDD subprocess testing"** (emerso da FR-12):

Quando un CLI Go ha un client di rete configurabile (LLM, DB, API), aggiungere:
1. Un env var override (`LABNEXUS_<COMPONENT>_ENDPOINT`) letto al `provider.Select`-level
2. Precedenza: flag CLI > env var > campo nel config file
3. Nel BDD harness, montare `httptest.NewServer` e settare l'env var prima di `exec.Command(binary, args...)`

Beneficio: il subprocess CLI è testabile end-to-end senza dipendenze esterne reali, mantenendo il codice di produzione invariato (nessun "test mode" flag o iniezione di mock). Più documentato in `.pipeline/skills/project/dual-provider-env-override.md` se promuoviamo.

## System Updates Applied

Nessuno applicato in questa fetta. Proposte (da approvare):
- **Skill di progetto** `validazione-due-livelli-cto-cliente.md`: codifica il pattern L1+L2 emerso in NFR-8.
- **Skill di progetto** `dual-provider-env-override.md`: codifica il pattern di testabilità via env var override sopra.

(Entrambe potrebbero essere candidate per framework promotion in `proposed-updates/` se l'utente concorda.)

## Tech Debt

Identificati 3 trade-off accettati (LOW severity, documentati in `review.md` loop 2):

1. **`features/steps_test.go` 1142 righe in singolo file**: navigabile via `grep -n` ma scomodo; splittare in 7 file (uno per gruppo) richiede tempo non giustificato finché il file non cresce > 1500 righe o entrano altri developer.
2. **`vobbly` missing**: `.pipeline/bin/vobbly` non installato → review multi-LLM impossibile. È setup workstation, non codice. Da risolvere prima del prossimo `/v-review`.
3. **Tag esatto `qwen3.6` nel registry Ollama non verificato runtime**: documentato in `DEV-GUIDE.md` → l'utente al primo `ollama list` verifica e aggiorna il campo `modello:` nei profili se necessario.

**Raccomandazione**: **NON serve `/v-refactor`** prima di Fetta 2. Tutti e tre i punti sono trade-off documentati, non breaking debt. Procedere con `/v-plan` per Fetta 2.

Funzioni 32-34 righe ancora presenti (consumeOllamaStream, consumeEurouterStream, OllamaProvider.doRequest, runner.Run): accettate come "thin orchestration coesa" (Layer 1 balanced modularity). Se Fetta 2 le toccherà, valutare ulteriore split solo se la nuova logica le rende davvero più complesse.
