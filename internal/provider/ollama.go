package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// OllamaProvider chiama Ollama su /api/chat con streaming NDJSON
// (un oggetto JSON per linea con campo message.content). FR-7, EC-7, EC-13.
type OllamaProvider struct {
	Endpoint   string
	HTTPClient *http.Client
}

func (p *OllamaProvider) Name() string { return "ollama" }

const maxConsecutiveBadChunks = 10

// ollamaDefaultTimeoutSeconds è il valore di default in secondi quando l'env
// `LABNEXUS_HTTP_TIMEOUT` non è settato (o ha valore non valido).
// 30 minuti = 1800s coprono il warmup + context fill + decode per:
//   - Qwen 3.6 36B MoE (23 GB su disco) su Apple Silicon
//   - Prompt fino a ~128k token (limite context window)
//   - max_tokens=8192 di decode con reasoning model (<think> block)
//
// Casi più estremi (modelli ancora più grandi, prompt più grandi) richiedono
// l'override esplicito via env var.
const ollamaDefaultTimeoutSeconds = 1800

// OllamaDefaultTimeout ritorna il timeout HTTP di default per OllamaProvider.
// Precedenza:
//  1. Env var LABNEXUS_HTTP_TIMEOUT (in secondi) — per casi estremi
//  2. ollamaDefaultTimeoutSeconds (30 min) — sufficiente per qwen3.6 36B
//
// Fix per `.pipeline/bugs/ollama-http-timeout-troppo-stretto-per-reasoning-models.md`.
func OllamaDefaultTimeout(p *OllamaProvider) time.Duration {
	if v := os.Getenv("LABNEXUS_HTTP_TIMEOUT"); v != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && n > 0 {
			return time.Duration(n) * time.Second
		}
		// Valore non valido: fallback al default.
	}
	return time.Duration(ollamaDefaultTimeoutSeconds) * time.Second
}

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
		// Timeout esplicito via OllamaDefaultTimeout (30 min default, override
		// con env LABNEXUS_HTTP_TIMEOUT). Generoso perché modelli reasoning
		// grandi (qwen3.6 36B MoE) su prompt large fanno context fill +
		// thinking per minuti prima del primo byte di response.
		client = &http.Client{Timeout: OllamaDefaultTimeout(p)}
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
