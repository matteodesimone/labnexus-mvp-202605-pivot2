package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// OllamaProvider chiama Ollama su /api/chat con streaming NDJSON
// (un oggetto JSON per linea con campo message.content). FR-7, EC-7, EC-13.
type OllamaProvider struct {
	Endpoint   string
	HTTPClient *http.Client
}

func (p *OllamaProvider) Name() string { return "ollama" }

const maxConsecutiveBadChunks = 10

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaRequest struct {
	Model    string                 `json:"model"`
	Messages []ollamaMessage        `json:"messages"`
	Stream   bool                   `json:"stream"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

type ollamaChunk struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

// Stream invia la richiesta a Ollama e ritorna un canale di StreamEvent.
func (p *OllamaProvider) Stream(ctx context.Context, system, user string, opts Options) (<-chan StreamEvent, error) {
	resp, err := p.doRequest(ctx, system, user, opts)
	if err != nil {
		return nil, err
	}
	ch := make(chan StreamEvent, 32)
	go consumeOllamaStream(resp.Body, ch)
	return ch, nil
}

func (p *OllamaProvider) doRequest(ctx context.Context, system, user string, opts Options) (*http.Response, error) {
	endpoint := p.Endpoint
	if endpoint == "" {
		endpoint = "http://localhost:11434"
	}
	body, err := json.Marshal(buildOllamaRequest(system, user, opts))
	if err != nil {
		return nil, fmt.Errorf("ollama: encode request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ollama: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := p.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrOllamaUnreachable, err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		// 404 o 5xx significa: connessione aperta ma il server non è Ollama (o è in errore).
		// Trattato come "Ollama non raggiungibile" per coerenza con EC-7.
		if resp.StatusCode == http.StatusNotFound || resp.StatusCode >= http.StatusInternalServerError {
			return nil, fmt.Errorf("%w: l'endpoint ha risposto HTTP %d (non sembra essere Ollama)", ErrOllamaUnreachable, resp.StatusCode)
		}
		return nil, fmt.Errorf("ollama: HTTP %d", resp.StatusCode)
	}
	return resp, nil
}

func buildOllamaRequest(system, user string, opts Options) ollamaRequest {
	r := ollamaRequest{
		Model: opts.Modello,
		Messages: []ollamaMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
		Stream: true,
	}
	r.Options = buildOllamaOptions(opts)
	return r
}

func buildOllamaOptions(opts Options) map[string]interface{} {
	if opts.Temperature <= 0 && opts.ContextWindow <= 0 && opts.MaxTokens <= 0 {
		return nil
	}
	m := map[string]interface{}{}
	if opts.Temperature > 0 {
		m["temperature"] = opts.Temperature
	}
	if opts.ContextWindow > 0 {
		m["num_ctx"] = opts.ContextWindow
	}
	if opts.MaxTokens > 0 {
		m["num_predict"] = opts.MaxTokens
	}
	return m
}

// consumeOllamaStream legge incrementalmente NDJSON e tollera fino a
// maxConsecutiveBadChunks chunk malformati (EC-13).
func consumeOllamaStream(body interface{ Read(p []byte) (n int, err error); Close() error }, ch chan<- StreamEvent) {
	defer close(ch)
	defer body.Close()
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	bad := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var c ollamaChunk
		if err := json.Unmarshal([]byte(line), &c); err != nil {
			bad++
			if bad >= maxConsecutiveBadChunks {
				ch <- StreamEvent{Err: fmt.Errorf("ollama: streaming Ollama corrotto (%d chunk malformati consecutivi)", bad)}
				return
			}
			continue
		}
		bad = 0
		if c.Message.Content != "" {
			ch <- StreamEvent{Token: c.Message.Content}
		}
		if c.Done {
			ch <- StreamEvent{Done: true}
			return
		}
	}
	if err := scanner.Err(); err != nil {
		ch <- StreamEvent{Err: fmt.Errorf("ollama: scanner: %w", err)}
	}
}
