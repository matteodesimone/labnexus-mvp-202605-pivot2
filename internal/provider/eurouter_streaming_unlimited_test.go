package provider_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labnexus/labnexus/internal/provider"
)

// TestEurouterProvider_BodyCanStreamLongerThanHeaderTimeout riproduce lo
// stesso bug già fixato per ollama (ollama_streaming_unlimited_test.go):
// con `Client.Timeout=N` overall, anche se gli headers arrivano subito,
// se il body streaming dura > N secondi il client cancella e il run fallisce
// con "context deadline exceeded".
//
// Smoke reale 2026-05-22 10:51 ha mostrato il taglio: 935 token streammati
// in 1m47s + timeout overall = stream error a 120s. Eurouter ha ResponseHeader
// rapidi (TTFB 38s) ma body lunghi su modelli grossi (qwen3.5-122b @ 8.7 tok/s).
//
// Comportamento atteso post-fix: `Transport.ResponseHeaderTimeout` (solo TTFB) +
// `Client.Timeout=0` (body streaming illimitato), identico al pattern ollama.
func TestEurouterProvider_BodyCanStreamLongerThanHeaderTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Headers SSE immediati.
		w.Header().Set("Content-Type", "text/event-stream")
		fl, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("ResponseWriter non supporta http.Flusher")
		}
		w.WriteHeader(http.StatusOK)
		fl.Flush()
		// Streaming lento: 4 chunk a 600ms = ~2.4s totali di body.
		for _, c := range []string{"alpha ", "beta ", "gamma ", "delta"} {
			time.Sleep(600 * time.Millisecond)
			_, _ = w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"" + c + "\"}}]}\n\n"))
			fl.Flush()
		}
		// Terminatore SSE.
		time.Sleep(100 * time.Millisecond)
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
		fl.Flush()
	}))
	defer srv.Close()

	// Default client (post-fix): scarica tutto il body anche se dura > 1s.
	// Timeout=1 setting esplicito a 1s nello struct: post-fix questo si
	// applica a ResponseHeaderTimeout (TTFB), non al body overall.
	p := &provider.EurouterProvider{Endpoint: srv.URL, APIKey: "test-key", Timeout: 1}
	ch, err := p.Stream(context.Background(), "s", "u", provider.Options{Modello: "test"})
	if err != nil {
		t.Fatalf("Stream error: %v", err)
	}
	var got []string
	deadline := time.After(10 * time.Second)
	for {
		select {
		case <-deadline:
			t.Fatalf("test timeout (10s) — stream non completato. Tokens ricevuti finora: %v", got)
		case ev, ok := <-ch:
			if !ok {
				goto done
			}
			if ev.Err != nil {
				t.Fatalf("stream interrotto da Err (bug del Client.Timeout overall): %v\nTokens ricevuti: %v", ev.Err, got)
			}
			if ev.Token != "" {
				got = append(got, ev.Token)
			}
			if ev.Done {
				goto done
			}
		}
	}
done:
	joined := strings.Join(got, "")
	if !strings.Contains(joined, "alpha") || !strings.Contains(joined, "delta") {
		t.Errorf("body completo richiesto (alpha + delta), got: %q", joined)
	}
}
