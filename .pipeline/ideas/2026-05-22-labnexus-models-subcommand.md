---
title: labnexus models subcommand per listare i modelli eurouter disponibili
status: planned
source: compound learning sub-fetta 1.5.C (post-deploy smoke HTTP 404 modello)
created: 2026-05-22
target: Sprint 2 (cli-tool ergonomic improvement)
---

# Idea: `labnexus models` subcommand

## Problema

Sub-fetta 1.5.C deploy ha incontrato HTTP 404 al primo smoke reale perché il modello configurato (`qwen3.6`, eredità Ollama-style Sprint 1) non era nel catalogo eurouter. Per scoprire il nome corretto abbiamo dovuto: (1) leggere il body errore JSON, (2) aprire manualmente il dashboard eurouter, (3) cercare nella lista modelli quale era il match più vicino, (4) provare a indovinare la stringa esatta.

Workflow inefficiente. Per debug + onboarding di nuovi profili (Denis o futuri utenti) servirebbe un comando interno.

## Soluzione proposta

Nuovo subcommand:

```bash
labnexus models                      # lista tutti i modelli disponibili
labnexus models --filter qwen        # solo i modelli che matchano substring
labnexus models --json               # output strutturato per scripting
```

Implementazione (Go, ~80 LoC):

```go
// internal/provider/eurouter.go
func (p *EurouterProvider) ListModels(ctx context.Context) ([]Model, error) {
    endpoint := p.modelsEndpoint() // derivato da Endpoint sostituendo /chat/completions → /models
    req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
    req.Header.Set("Authorization", "Bearer "+p.APIKey)
    // ... HTTP call + JSON decode response
}

type Model struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    ContextWindow int   `json:"context_window"`
    Pricing     struct {
        Input  float64 `json:"input"`
        Output float64 `json:"output"`
    } `json:"pricing"`
    // ... altri campi catalogo eurouter
}
```

```go
// cmd/labnexus/main.go
func newModelsCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "models",
        Short: "lista i modelli disponibili sul provider eurouter",
        RunE: func(cmd *cobra.Command, _ []string) error {
            // load master per API key, instantiate EurouterProvider, call ListModels
            // formatta output human o JSON
        },
    }
    cmd.Flags().String("filter", "", "substring per filtrare i modelli (es. 'qwen')")
    cmd.Flags().Bool("json", false, "output JSON strutturato")
    return cmd
}
```

Output human esempio:

```
ID                              Context  Pricing in/out ($/M)   Features
qwen3.5-122b-a10b              262144   $0.50 / $3.50         Thinking, Tools, JSON
qwen3.5-9b                     262144   $0.10 / $0.15         Thinking, Vision, Tools, JSON
qwen3-32b                      131072   $0.08 / $0.23         Thinking, Tools, JSON
qwen3-235b-a22b-instruct       131072   $0.07 / $0.46         Tools, JSON
... (8 altri)

13 modelli (filtro: 'qwen')
```

## Trade-off

- **Pro**: Debug rapidissimo quando modello configurato non funziona. Onboarding di nuovi profili guidato (Denis vede subito le opzioni col loro context e prezzo).
- **Pro**: Generalizzabile a Ollama (`ollama list` equivalente integrato).
- **Con**: Stima 2-3h CTO. Sprint 2 scope.
- **Con**: Eurouter API potrebbe non avere `/v1/models` standardizzato (OpenAI-compatible ne ha uno, da verificare per eurouter). Se manca, fallback a "scrape dashboard" non è elegante.

## Implementazione minimal pragmatic

Sprint 2 quick win: aggiungere comando shell helper nella `GUIDA.md` per chi sa Terminal:

```bash
# Lista modelli Qwen disponibili
curl -s https://api.eurouter.ai/api/v1/models \
  -H "Authorization: Bearer $EUROUTER_API_KEY" \
  | jq '.data[] | select(.id | test("qwen"; "i")) | {id, context: .context_length, pricing}'
```

Questa è 0 LoC Go ma richiede Denis usi Terminal. Per il CTO è OK.

## Status

Rimandato a Sprint 2 — non blocker per shakedown reale Denis post-credit-topup.
