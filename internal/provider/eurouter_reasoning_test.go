package provider

import (
	"encoding/json"
	"io"
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
	consumeEurouterStream(io.NopCloser(strings.NewReader(sse)), ch)

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
	consumeEurouterStream(io.NopCloser(strings.NewReader(sse)), ch)

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
	consumeEurouterStream(io.NopCloser(strings.NewReader(body)), ch)
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
	consumeEurouterStream(io.NopCloser(strings.NewReader(sse)), ch)
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
