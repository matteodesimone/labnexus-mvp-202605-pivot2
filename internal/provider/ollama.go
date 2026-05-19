package provider

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// OllamaDefaultTimeout ritorna il timeout HTTP TTFB (Time-To-First-Byte) per OllamaProvider.
// Semantica post-bugfix-5: questo è il `ResponseHeaderTimeout` del Transport,
// NON un overall timeout. Il body streaming non è limitato.
//
// Precedenza:
//  1. Env var LABNEXUS_HTTP_TIMEOUT (in secondi) — per casi estremi
//  2. ollamaDefaultTimeoutSeconds (30 min) — sufficiente per warmup qwen3.6 36B
func OllamaDefaultTimeout(p *OllamaProvider) time.Duration {
	if v := os.Getenv("LABNEXUS_HTTP_TIMEOUT"); v != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && n > 0 {
			return time.Duration(n) * time.Second
		}
		// Valore non valido: fallback al default.
	}
	return time.Duration(ollamaDefaultTimeoutSeconds) * time.Second
}

// OllamaDefaultHTTPClient costruisce il default *http.Client quando
// OllamaProvider.HTTPClient è nil.
//
// Design (post-bugfix-5 `http-overall-timeout-tronca-streaming-llm`):
//   - `Transport.ResponseHeaderTimeout = OllamaDefaultTimeout(p)` — TTFB
//   - `Client.Timeout = 0` — NESSUN overall timeout, body streaming illimitato
//
// Per LLM streaming il body può legittimamente durare ore (output lunghi su
// modelli grandi). Un overall timeout tronca a metà l'output. La protezione
// contro "server hung" resta tramite ResponseHeaderTimeout (cancella se gli
// headers non arrivano in tempo).
func OllamaDefaultHTTPClient(p *OllamaProvider) *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ResponseHeaderTimeout = OllamaDefaultTimeout(p)
	return &http.Client{Transport: transport} // Timeout: 0 (no overall, body unlimited)
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
		// Default client: Transport.ResponseHeaderTimeout (TTFB) + Client.Timeout=0
		// (no overall, body streaming illimitato). Vedi OllamaDefaultHTTPClient.
		// Fix bug `http-overall-timeout-tronca-streaming-llm`.
		client = OllamaDefaultHTTPClient(p)
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
//
// Edge case "EOF post-content" (bug `eof-post-content-falsamente-interrotto`):
// se il server chiude la connessione DOPO aver emesso content ma SENZA il
// chunk finale `done: true`, emettiamo `Done: true, NoDoneMarker: true` per
// informare il caller che l'output è verosimilmente completo ma manca il
// marker formale (capita con alcuni modelli reasoning su Ollama).
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
			// Caso normale: il server ha emesso done:true. NoDoneMarker resta false.
			ch <- StreamEvent{Done: true}
			return
		}
	}
	// Il loop è terminato senza done:true. Casi:
	//   1) scanner.Err() == nil: EOF clean. Probabile "completato senza done marker".
	//   2) scanner.Err() == io.ErrUnexpectedEOF: NDJSON troncato a metà linea. Stesso caso:
	//      il server ha chiuso, possibly dopo aver completato il modello.
	//   3) altri errori scanner: errore vero (es. token buffer overflow).
	scanErr := scanner.Err()
	if scanErr != nil && !errorsIsUnexpectedEOF(scanErr) {
		ch <- StreamEvent{Err: fmt.Errorf("ollama: scanner: %w", scanErr)}
		return
	}
	// EOF clean o unexpected EOF post-content: trattiamo come completato
	// con NoDoneMarker=true, il caller (drainStream) loggerà un warning.
	ch <- StreamEvent{Done: true, NoDoneMarker: true}
}

func errorsIsUnexpectedEOF(err error) bool {
	return errors.Is(err, io.ErrUnexpectedEOF)
}
