package provider

import (
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
