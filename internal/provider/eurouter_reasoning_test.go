package provider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// kimi e altri modelli "thinking" streamano i token di ragionamento in
// delta.reasoning_content (o delta.reasoning), separati dal contenuto, e
// chiudono con finish_reason. Il parser deve catturarli (diagnostica TTFT
// lunghi / risposte vuote) senza inquinare l'output.
func TestConsumeEurouterStream_ReasoningAndFinishReason(t *testing.T) {
	sse := strings.Join([]string{
		`data: {"choices":[{"delta":{"reasoning_content":"sto pensando..."}}]}`,
		`data: {"choices":[{"delta":{"reasoning":" ancora"}}]}`,
		`data: {"choices":[{"delta":{"content":"Risposta finale"}}]}`,
		`data: {"choices":[{"delta":{},"finish_reason":"stop"}]}`,
		`data: [DONE]`,
		"",
	}, "\n")

	ch := make(chan StreamEvent, 32)
	consumeEurouterStream(context.Background(), io.NopCloser(strings.NewReader(sse)), http.Header{}, ch)

	var reasoning, content, finish string
	var done bool
	for ev := range ch {
		reasoning += ev.Reasoning
		content += ev.Token
		if ev.FinishReason != "" {
			finish = ev.FinishReason
		}
		if ev.Done {
			done = true
		}
	}

	if reasoning != "sto pensando... ancora" {
		t.Errorf("reasoning = %q, want catturato da reasoning_content + reasoning", reasoning)
	}
	if content != "Risposta finale" {
		t.Errorf("content = %q (il reasoning NON deve finirci dentro)", content)
	}
	if finish != "stop" {
		t.Errorf("finish_reason = %q, want stop", finish)
	}
	if !done {
		t.Error("manca l'evento Done")
	}
}

// Caso "F vuota": tutto reasoning, zero contenuto, finish_reason length →
// deve emergere dai dati (content vuoto, reasoning presente, finish length).
func TestConsumeEurouterStream_EmptyContentAllReasoning(t *testing.T) {
	sse := strings.Join([]string{
		`data: {"choices":[{"delta":{"reasoning_content":"penso a lungo senza concludere"}}]}`,
		`data: {"choices":[{"delta":{},"finish_reason":"length"}]}`,
		`data: [DONE]`,
		"",
	}, "\n")

	ch := make(chan StreamEvent, 32)
	consumeEurouterStream(context.Background(), io.NopCloser(strings.NewReader(sse)), http.Header{}, ch)

	var reasoning, content, finish string
	for ev := range ch {
		reasoning += ev.Reasoning
		content += ev.Token
		if ev.FinishReason != "" {
			finish = ev.FinishReason
		}
	}
	if content != "" {
		t.Errorf("content deve essere vuoto, got %q", content)
	}
	if reasoning == "" || finish != "length" {
		t.Errorf("attesi reasoning presente + finish length, got reasoning=%q finish=%q", reasoning, finish)
	}
}

// Errore del provider su HTTP 200 con body NON-SSE (es. JSON {"error":...}):
// deve EMERGERE nel messaggio d'errore, non sparire dietro "timeout".
func TestConsumeEurouterStream_SurfacesNonSSEError(t *testing.T) {
	body := `{"error":{"message":"insufficient credits","type":"billing"}}`
	ch := make(chan StreamEvent, 8)
	consumeEurouterStream(context.Background(), io.NopCloser(strings.NewReader(body)), http.Header{}, ch)
	var gotErr error
	for ev := range ch {
		if ev.Err != nil {
			gotErr = ev.Err
		}
	}
	if gotErr == nil || !strings.Contains(gotErr.Error(), "insufficient credits") {
		t.Errorf("l'errore del provider deve emergere, got %v", gotErr)
	}
}

// stream_options.include_usage deve essere nella richiesta (per ricevere i
// token reali), e il chunk usage deve essere parsato e emesso.
func TestEurouter_UsageRequestedAndParsed(t *testing.T) {
	b, _ := json.Marshal(buildEurouterRequest("s", "u", Options{Modello: "kimi-k2.6", MaxTokens: 0, Temperature: 1.0}))
	if !strings.Contains(string(b), `"stream_options":{"include_usage":true}`) {
		t.Errorf("la richiesta deve chiedere stream_options.include_usage, got %s", b)
	}
	sse := strings.Join([]string{
		`data: {"choices":[{"delta":{"content":"ciao"}}]}`,
		`data: {"choices":[],"usage":{"prompt_tokens":1200,"completion_tokens":800,"total_tokens":2000,"completion_tokens_details":{"reasoning_tokens":500}}}`,
		`data: [DONE]`,
		"",
	}, "\n")
	ch := make(chan StreamEvent, 16)
	consumeEurouterStream(context.Background(), io.NopCloser(strings.NewReader(sse)), http.Header{}, ch)
	var u *Usage
	for ev := range ch {
		if ev.Usage != nil {
			u = ev.Usage
		}
	}
	if u == nil {
		t.Fatal("usage non parsato")
	}
	if u.PromptTokens != 1200 || u.CompletionTokens != 800 || u.ReasoningTokens != 500 || u.TotalTokens != 2000 {
		t.Errorf("usage = %+v, want prompt=1200 completion=800 reasoning=500 total=2000", u)
	}
}

// Errore "SSE-framed": il gateway risponde 200, apre lo stream e poi invia
// data: {"error":{...}} (es. context length / upstream / rate-limit). Prima
// veniva deserializzato vuoto e SCARTATO; ora deve emergere nell'errore, non
// nel generico "stream vuoto".
func TestConsumeEurouterStream_SurfacesSSEFramedError(t *testing.T) {
	sse := strings.Join([]string{
		`data: {"error":{"message":"the request exceeds the maximum context length","type":"invalid_request_error","code":"context_length_exceeded"}}`,
		`data: [DONE]`,
		"",
	}, "\n")
	ch := make(chan StreamEvent, 8)
	consumeEurouterStream(context.Background(), io.NopCloser(strings.NewReader(sse)), http.Header{}, ch)
	var gotErr error
	var content string
	for ev := range ch {
		if ev.Err != nil {
			gotErr = ev.Err
		}
		content += ev.Token
	}
	if content != "" {
		t.Errorf("nessun contenuto atteso, got %q", content)
	}
	if gotErr == nil || !strings.Contains(gotErr.Error(), "context length") {
		t.Errorf("l'errore SSE-framed del provider deve emergere, got %v", gotErr)
	}
}

// 200 con body 0 byte: senza payload da mostrare, l'errore deve almeno riportare
// gli header diagnostici (request-id, retry-after) — tutto ciò che il provider
// manda, meta inclusa, va mostrato/loggato.
func TestConsumeEurouterStream_EmptyBodySurfacesHeaders(t *testing.T) {
	h := http.Header{}
	h.Set("X-Request-Id", "req_abc123")
	h.Set("Retry-After", "30")
	ch := make(chan StreamEvent, 8)
	consumeEurouterStream(context.Background(), io.NopCloser(strings.NewReader("")), h, ch)
	var gotErr error
	for ev := range ch {
		if ev.Err != nil {
			gotErr = ev.Err
		}
	}
	if gotErr == nil {
		t.Fatal("atteso un errore per stream vuoto")
	}
	if !strings.Contains(gotErr.Error(), "req_abc123") || !strings.Contains(gotErr.Error(), "30") {
		t.Errorf("l'errore deve riportare gli header diagnostici (request-id, retry-after), got %v", gotErr)
	}
}
