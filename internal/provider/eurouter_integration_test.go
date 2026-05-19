package provider_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/provider"
)

// fakeEurouterServer emette SSE OpenAI-compatible (data: {...}, terminatore [DONE]).
func fakeEurouterServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		fl, _ := w.(http.Flusher)
		writeChunk := func(content string) {
			_, _ = w.Write([]byte(`data: {"choices":[{"delta":{"content":"` + content + `"}}]}` + "\n\n"))
			if fl != nil {
				fl.Flush()
			}
		}
		writeChunk("uno ")
		writeChunk("due ")
		writeChunk("tre")
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
	}))
}

func TestEurouterProvider_StreamsSSEIncrementally(t *testing.T) {
	srv := fakeEurouterServer(t)
	defer srv.Close()

	p := &provider.EurouterProvider{Endpoint: srv.URL, APIKey: "test-key", Timeout: 5}
	ch, err := p.Stream(context.Background(), "s", "u", provider.Options{Modello: "test-model"})
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
	if !strings.Contains(joined, "uno") || !strings.Contains(joined, "due") || !strings.Contains(joined, "tre") {
		t.Fatalf("expected uno/due/tre in stream, got %q", joined)
	}
}

func TestEurouterProvider_MissingAPIKeyTypedError(t *testing.T) {
	p := &provider.EurouterProvider{Endpoint: "http://localhost:1234", APIKey: ""}
	_, err := p.Stream(context.Background(), "s", "u", provider.Options{})
	if err == nil {
		t.Fatal("expected ErrMissingAPIKey when APIKey empty, got nil")
	}
}
