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
}

type sseChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
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
	go consumeEurouterStream(resp.Body, ch)
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
		Stream: true,
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
func consumeEurouterStream(body interface{ Read(p []byte) (n int, err error); Close() error }, ch chan<- StreamEvent) {
	defer close(ch)
	defer body.Close()
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	sawDone := false
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
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
			continue // tolerante: chunk malformato isolato
		}
		for _, choice := range c.Choices {
			if choice.Delta.Content != "" {
				ch <- StreamEvent{Token: choice.Delta.Content}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		ch <- StreamEvent{Err: fmt.Errorf("eurouter: scanner: %w", err)}
		return
	}
	if !sawDone {
		ch <- StreamEvent{Err: fmt.Errorf("eurouter: timeout senza terminatore [DONE]")}
	}
}
