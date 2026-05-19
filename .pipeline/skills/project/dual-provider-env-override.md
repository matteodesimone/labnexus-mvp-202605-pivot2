# Skill — Dual-provider con env override per BDD subprocess testing

Pattern: rendere il client di rete di un CLI configurabile via env var endpoint, in modo che i BDD acceptance test possano sostituirlo con un `httptest.NewServer` senza modificare il codice di produzione. Emerso da FR-12 LabNexus.

## Quando applicarlo

- CLI Go (o equivalente) che chiama un servizio di rete (LLM, API, DB remoto)
- I BDD acceptance invocano il binario come subprocess (`exec.Command`)
- Vuoi testare l'intero flusso end-to-end senza dipendere dal servizio reale, mantenendo invariato il codice di produzione

## Pattern

### 1. Endpoint configurabile a 3 livelli con precedenza

Nel codice di selezione del provider/client:
```go
func Select(flagOverride, envOverride, configFromProfile string) (Client, error) {
    chosen := strings.TrimSpace(flagOverride)
    if chosen == "" {
        chosen = strings.TrimSpace(envOverride)
    }
    if chosen == "" {
        chosen = strings.TrimSpace(configFromProfile)
    }
    switch chosen {
    case "providerA":
        endpoint := os.Getenv("LABNEXUS_PROVIDERA_ENDPOINT")
        if endpoint == "" {
            endpoint = "https://default.providerA.example/api"
        }
        return &ProviderAClient{Endpoint: endpoint}, nil
    // ...
    }
}
```

Precedenza: **flag > env > config**. La selezione del *provider* è separata dall'endpoint del provider scelto.

### 2. BDD harness con fake server + env injection

In `features/helpers_test.go` (o equivalente):
```go
func (s *scenarioState) startFakeProviderA(content string) {
    s.fakeProviderA = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Risposta predefinita per il test
    }))
}

func (s *scenarioState) buildEnv() []string {
    env := os.Environ()
    if s.fakeProviderA != nil {
        env = append(env, "LABNEXUS_PROVIDERA_ENDPOINT="+s.fakeProviderA.URL)
    }
    return env
}

func (s *scenarioState) runBinary(args []string) error {
    cmd := exec.Command(binaryPath, args...)
    cmd.Env = s.buildEnv()
    // ...
}
```

### 3. Auto-setup nello step `lancio`

```go
func lancio(c context.Context, cmd string) (context.Context, error) {
    s := getState(c)
    if needsLLM(args) && s.fakeProviderA == nil && s.env["LABNEXUS_PROVIDERA_ENDPOINT"] == "" {
        s.startFakeProviderA("default fake response")
    }
    // ...
}
```

Per scenari che testano il fallimento ("provider A non disponibile"):
```go
func providerANonInAscolto(c context.Context) (context.Context, error) {
    s := getState(c)
    // NON avvia fake; punta a un endpoint morto.
    s.env["LABNEXUS_PROVIDERA_ENDPOINT"] = "http://127.0.0.1:1"
    return c, nil
}
```

## Beneficio

- **Codice di produzione invariato**: nessuna "test mode" flag, nessuna iniezione di mock
- **Zone-2 mock boundary corretto** (vedi `01-methodology.md`): si mocka solo l'HTTP boundary, tutto il resto è reale (filesystem, profilo, parsing)
- **BDD end-to-end ripetibili**: nessuna dipendenza da Ollama/EUrouter/qualunque provider reale per girare i test in CI
- **Test del path "provider down"**: facile da scrivere puntando l'env a un endpoint morto

## Tradeoff / gotchas

- Su alcuni sistemi (es. macOS), porte basse (1) possono rispondere con HTTP 404 invece di "connection refused". Soluzione: mappare HTTP 4xx/5xx come `ErrProviderUnreachable` nel client di produzione (vedi `OllamaProvider.doRequest` in LabNexus).
- L'env var override è anche un canale di sicurezza interessante: in produzione **non documentarlo** se non strettamente necessario, per evitare che l'utente la usi per puntare a un endpoint non autorizzato.

## Esempio in LabNexus Sprint 1

- Provider primary `ollama` (env `LABNEXUS_OLLAMA_ENDPOINT`) e secondary `eurouter` (env `LABNEXUS_EUROUTER_ENDPOINT`)
- `provider.Select(providerOverride, providerEnv, providerFromProfile)` implementa la precedenza
- `features/helpers_test.go` ha `startFakeOllama` / `startFakeEurouter` con httptest server
- 40 scenari BDD passing senza Ollama né internet, includono il path "Ollama unreachable" via env a porta morta

## Vedi anche

- `.pipeline/solutions/2026-05-19-fetta1-scooter.md` — emersione del pattern
- FR-12 nella spec originale
- `internal/provider/provider.go` — `Select()`
- `features/helpers_test.go` — `scenarioState.buildEnv()` + fake server
