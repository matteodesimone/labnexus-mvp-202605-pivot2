package provider_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/provider"
)

// fakeOllamaServer emette NDJSON di chunk Ollama, un oggetto per riga,
// con campo message.content (FR-7).
func fakeOllamaServer(t *testing.T, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		// 3 chunk fittizi seguiti da done:true
		_, _ = w.Write([]byte(`{"message":{"content":"alpha "}}` + "\n"))
		_, _ = w.Write([]byte(`{"message":{"content":"beta "}}` + "\n"))
		_, _ = w.Write([]byte(`{"message":{"content":"gamma"},"done":true}` + "\n"))
	}))
}

func TestOllamaProvider_StreamsTokensIncrementally(t *testing.T) {
	srv := fakeOllamaServer(t, "")
	defer srv.Close()

	p := &provider.OllamaProvider{Endpoint: srv.URL}
	ch, err := p.Stream(context.Background(), "system", "user", provider.Options{Modello: "qwen3.6"})
	if err != nil {
		t.Fatalf("Stream error: %v", err)
	}
	var got []string
	for ev := range ch {
		if ev.Err != nil {
			t.Fatalf("stream event error: %v", ev.Err)
		}
		got = append(got, ev.Token)
		if ev.Done {
			break
		}
	}
	joined := strings.Join(got, "")
	if !strings.Contains(joined, "alpha") || !strings.Contains(joined, "beta") || !strings.Contains(joined, "gamma") {
		t.Fatalf("expected alpha/beta/gamma in stream, got %q", joined)
	}
}

func TestOllamaProvider_UnreachableReturnsTypedError(t *testing.T) {
	p := &provider.OllamaProvider{Endpoint: "http://127.0.0.1:1"} // certo non in ascolto
	_, err := p.Stream(context.Background(), "s", "u", provider.Options{})
	if err == nil {
		t.Fatal("expected error when Ollama unreachable, got nil")
	}
	if !strings.Contains(err.Error(), "Ollama") {
		t.Errorf("error should mention 'Ollama' for diagnosis, got: %v", err)
	}
}
