package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// EurouterProvider chiama https://api.eurouter.ai/v1/chat/completions con SSE
// OpenAI-compatible (data: {...}, terminatore data: [DONE]). FR-7, FR-12, EC-6, EC-14.
type EurouterProvider struct {
	Endpoint   string
	APIKey     string
	Timeout    int
	HTTPClient *http.Client
}

func (p *EurouterProvider) Name() string { return "eurouter" }

type sseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type sseRequest struct {
	Model    string       `json:"model"`
	Messages []sseMessage `json:"messages"`
	Stream   bool         `json:"stream"`
	// Prima NON venivano inviati: il modello girava sui propri default (su un
	// modello "thinking" un max_tokens di default basso, consumato dal reasoning,
	// produce risposte vuote). omitempty: 0 = lascia il default del modello.
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	// StreamOptions.include_usage: chiede al provider di emettere il chunk finale
	// con i token reali (prompt/completion/reasoning). Meta del modello, loggato.
	StreamOptions *sseStreamOptions `json:"stream_options,omitempty"`
}

type sseStreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

type sseUsage struct {
	PromptTokens            int `json:"prompt_tokens"`
	CompletionTokens        int `json:"completion_tokens"`
	TotalTokens             int `json:"total_tokens"`
	CompletionTokensDetails *struct {
		ReasoningTokens int `json:"reasoning_tokens"`
	} `json:"completion_tokens_details"`
}

type sseChunk struct {
	Choices []struct {
		Delta struct {
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
			Reasoning        string `json:"reasoning"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
	Usage *sseUsage `json:"usage"`
	// Errore "SSE-framed": alcuni gateway OpenAI-compatible (incl. eurouter)
	// rispondono 200, aprono lo stream e poi notificano un errore UPSTREAM come
	// evento `data: {"error":{...}}`. Senza questo campo il chunk veniva
	// deserializzato vuoto e l'errore SCARTATO, lasciando solo "stream vuoto".
	Error json.RawMessage `json:"error"`
}

// Stream invia la richiesta SSE e ritorna un canale di StreamEvent.
// Errore esplicito se APIKey è vuota prima di qualsiasi chiamata di rete (FR-12, EC-6).
func (p *EurouterProvider) Stream(ctx context.Context, system, user string, opts Options) (<-chan StreamEvent, error) {
	if strings.TrimSpace(p.APIKey) == "" {
		return nil, ErrMissingAPIKey
	}
	resp, err := p.doRequest(ctx, system, user, opts)
	if err != nil {
		return nil, err
	}
	ch := make(chan StreamEvent, 32)
	go consumeEurouterStream(resp.Body, resp.Header, ch)
	return ch, nil
}

func (p *EurouterProvider) doRequest(ctx context.Context, system, user string, opts Options) (*http.Response, error) {
	endpoint := p.Endpoint
	if endpoint == "" {
		endpoint = "https://api.eurouter.ai/v1/chat/completions"
	}
	body, err := json.Marshal(buildEurouterRequest(system, user, opts))
	if err != nil {
		return nil, fmt.Errorf("eurouter: encode request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("eurouter: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)
	req.Header.Set("Accept", "text/event-stream")
	client := p.httpClientOrDefault()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("eurouter: request POST %s: %w", endpoint, err)
	}
	if resp.StatusCode != http.StatusOK {
		// Leggi il body di errore (max 4KB) per diagnosi: eurouter tipicamente
		// ritorna JSON tipo `{"error": "model 'X' not found"}`. Senza questo,
		// l'utente vedeva solo "HTTP 404" senza modo di capire il problema.
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		resp.Body.Close()
		bodyStr := strings.TrimSpace(string(errBody))
		if bodyStr == "" {
			return nil, fmt.Errorf("eurouter: HTTP %d POST %s (modello=%q, body risposta vuoto)", resp.StatusCode, endpoint, opts.Modello)
		}
		return nil, fmt.Errorf("eurouter: HTTP %d POST %s (modello=%q): %s", resp.StatusCode, endpoint, opts.Modello, bodyStr)
	}
	return resp, nil
}

func buildEurouterRequest(system, user string, opts Options) sseRequest {
	return sseRequest{
		Model: opts.Modello,
		Messages: []sseMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
		Stream:        true,
		MaxTokens:     opts.MaxTokens,
		Temperature:   opts.Temperature,
		StreamOptions: &sseStreamOptions{IncludeUsage: true},
	}
}

// httpClientOrDefault ritorna il client HTTP per le chiamate streaming.
// Pattern post-bugfix-5 (analogo a ollama, vedi OllamaDefaultHTTPClient):
//   - Transport.ResponseHeaderTimeout = p.Timeout (TTFB only) — protegge
//     contro server hung che non emette mai gli headers
//   - Client.Timeout = 0 — NESSUN limite overall sul body streaming
//
// Modelli grossi su eurouter (qwen3.5-122b in thinking mode su KB grandi)
// possono streammare il body per minuti a 5-10 tok/s; un Client.Timeout
// overall tronca il body e il run fallisce con "context deadline exceeded"
// (smoke 2026-05-22 10:51 ha mostrato il taglio a 1m47s con 935 token).
func (p *EurouterProvider) httpClientOrDefault() *http.Client {
	if p.HTTPClient != nil {
		return p.HTTPClient
	}
	timeout := time.Duration(p.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 120 * time.Second
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ResponseHeaderTimeout = timeout
	return &http.Client{Transport: transport} // Timeout: 0 (no overall, body unlimited)
}

// consumeEurouterStream legge SSE OpenAI-compatible (data: {...}, terminatore [DONE]).
// header sono gli header della risposta: su uno stream vuoto contengono spesso la
// prova del motivo (rate-limit, request-id per il supporto, content-type/length).
func consumeEurouterStream(body io.ReadCloser, header http.Header, ch chan<- StreamEvent) {
	defer close(ch)
	defer body.Close()
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	sawDone := false
	contentSeen := false
	var nonData strings.Builder   // righe NON-SSE: spesso contengono l'errore vero del provider (es. {"error":...} su un 200)
	var unhandled strings.Builder // payload `data:` che non producono eventi: errori SSE-framed o chunk malformati
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			if s := strings.TrimSpace(line); s != "" && nonData.Len() < 4096 {
				nonData.WriteString(s)
				nonData.WriteByte('\n')
			}
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "[DONE]" {
			sawDone = true
			ch <- StreamEvent{Done: true}
			return
		}
		var c sseChunk
		if err := json.Unmarshal([]byte(payload), &c); err != nil {
			// NON scartare: un payload `data:` malformato può essere/contenere
			// l'errore vero del gateway. Lo conserviamo per la diagnostica.
			if unhandled.Len() < 4096 {
				unhandled.WriteString(payload)
				unhandled.WriteByte('\n')
			}
			continue
		}
		// Errore SSE-framed: alcuni gateway rispondono 200, aprono lo stream e poi
		// inviano `data: {"error":{...}}` (es. context length, upstream, rate-limit).
		// Va fatto EMERGERE: senza, il chunk veniva deserializzato vuoto e scartato.
		if len(c.Error) > 0 && string(c.Error) != "null" {
			ch <- StreamEvent{Err: fmt.Errorf("eurouter: errore dal provider (evento SSE su HTTP 200): %s", strings.TrimSpace(string(c.Error)))}
			return
		}
		for _, choice := range c.Choices {
			// Reasoning (modelli "thinking" come kimi): delta.reasoning_content
			// o delta.reasoning. Emesso separato dal contenuto — diagnostico.
			if r := choice.Delta.ReasoningContent; r != "" {
				ch <- StreamEvent{Reasoning: r}
			} else if r := choice.Delta.Reasoning; r != "" {
				ch <- StreamEvent{Reasoning: r}
			}
			if choice.Delta.Content != "" {
				contentSeen = true
				ch <- StreamEvent{Token: choice.Delta.Content}
			}
			if choice.FinishReason != nil && *choice.FinishReason != "" {
				ch <- StreamEvent{FinishReason: *choice.FinishReason}
			}
		}
		// Chunk finale con usage (stream_options.include_usage): token reali.
		if u := c.Usage; u != nil {
			ev := StreamEvent{Usage: &Usage{
				PromptTokens:     u.PromptTokens,
				CompletionTokens: u.CompletionTokens,
				TotalTokens:      u.TotalTokens,
			}}
			if u.CompletionTokensDetails != nil {
				ev.Usage.ReasoningTokens = u.CompletionTokensDetails.ReasoningTokens
			}
			ch <- ev
		}
	}
	if err := scanner.Err(); err != nil {
		ch <- StreamEvent{Err: fmt.Errorf("eurouter: scanner: %w", err)}
		return
	}
	if sawDone {
		return
	}
	// Niente [DONE].
	if contentSeen {
		// EOF dopo contenuto: output verosimilmente completo, server senza [DONE].
		ch <- StreamEvent{Done: true, NoDoneMarker: true}
		return
	}
	// Zero contenuto: il provider non ha streammato nulla. Facciamo EMERGERE TUTTO
	// ciò che è arrivato — payload `data:` non gestiti, body non-SSE, e gli header
	// (request-id, rate-limit, content-type/length) che spiegano un 200 vuoto —
	// invece del generico "stream vuoto" che nascondeva la causa reale.
	var diag []string
	if s := strings.TrimSpace(unhandled.String()); s != "" {
		diag = append(diag, "payload non gestito: "+s)
	}
	if s := strings.TrimSpace(nonData.String()); s != "" {
		diag = append(diag, "body non-SSE: "+s)
	}
	if hs := headerSummary(header); hs != "" {
		diag = append(diag, "header risposta: "+hs)
	}
	if len(diag) > 0 {
		ch <- StreamEvent{Err: fmt.Errorf("eurouter: nessun token dal provider (HTTP 200 senza contenuto) — %s", strings.Join(diag, " | "))}
	} else {
		ch <- StreamEvent{Err: fmt.Errorf("eurouter: nessun token e nessun [DONE] — stream vuoto (body 0 byte, nessun header diagnostico): connessione chiusa dal gateway o capacità upstream esaurita")}
	}
}

// headerSummary estrae dagli header della risposta i campi diagnostici utili a
// spiegare un 200 senza contenuto: content-type/length (SSE vs JSON, body vuoto),
// request-id (per i ticket al provider), retry-after e rate-limit (capacità).
func headerSummary(h http.Header) string {
	if h == nil {
		return ""
	}
	var parts []string
	add := func(label, val string) {
		if v := strings.TrimSpace(val); v != "" {
			parts = append(parts, label+"="+v)
		}
	}
	add("content-type", h.Get("Content-Type"))
	add("content-length", h.Get("Content-Length"))
	add("retry-after", h.Get("Retry-After"))
	for _, k := range []string{"X-Request-Id", "X-Request-ID", "Cf-Ray", "X-Amzn-Requestid"} {
		add("request-id", h.Get(k))
	}
	for name, vals := range h {
		if strings.Contains(strings.ToLower(name), "ratelimit") && len(vals) > 0 {
			add(strings.ToLower(name), vals[0])
		}
	}
	return strings.Join(parts, ", ")
}
