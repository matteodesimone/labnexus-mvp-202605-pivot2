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

// TestOllamaProvider_BodyCanStreamLongerThanHeaderTimeout riproduce il bug
// `.pipeline/bugs/http-overall-timeout-tronca-streaming-llm.md`:
// con `Client.Timeout=N`, anche se gli headers arrivano immediatamente,
// se il body streaming dura > N secondi il client cancella e l'output è troncato.
//
// Comportamento atteso post-fix: usare `Transport.ResponseHeaderTimeout`
// (limite SOLO sul TTFB) + `Client.Timeout=0` (no overall). Headers rapidi
// passano il check, il body può streaming arbitrariamente.
//
// Test: fake server emette headers in ~10ms, poi sleep 600ms per chunk per 4 chunk
// (totale ~2.4s di streaming). Default client (post-fix) deve scaricare l'intero body.
func TestOllamaProvider_BodyCanStreamLongerThanHeaderTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Headers immediati (Ollama scrive 200 OK appena riceve la request).
		w.Header().Set("Content-Type", "application/x-ndjson")
		fl, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("ResponseWriter non supporta http.Flusher")
		}
		w.WriteHeader(http.StatusOK)
		fl.Flush()
		// Streaming lento del body: 4 chunk con ritardo 600ms ciascuno = ~2.4s totali.
		for _, c := range []string{"alpha ", "beta ", "gamma ", "delta"} {
			time.Sleep(600 * time.Millisecond)
			_, _ = w.Write([]byte(`{"message":{"content":"` + c + `"}}` + "\n"))
			fl.Flush()
		}
		// Chunk done finale
		time.Sleep(100 * time.Millisecond)
		_, _ = w.Write([]byte(`{"done":true}` + "\n"))
		fl.Flush()
	}))
	defer srv.Close()

	// Default client (post-fix): deve scaricare tutto il body anche se dura > 1s.
	p := &provider.OllamaProvider{Endpoint: srv.URL} // HTTPClient nil → default
	ch, err := p.Stream(context.Background(), "s", "u", provider.Options{Modello: "test"})
	if err != nil {
		t.Fatalf("Stream error: %v", err)
	}
	var got []string
	timeout := time.After(10 * time.Second)
	for {
		select {
		case <-timeout:
			t.Fatalf("test timeout (10s) — stream non completato. Tokens ricevuti finora: %v", got)
		case ev, ok := <-ch:
			if !ok {
				goto done
			}
			if ev.Err != nil {
				t.Fatalf("stream interrotto da Err (bug del client.Timeout overall): %v\nTokens ricevuti: %v", ev.Err, got)
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

// TestOllamaProvider_DefaultClientHasNoOverallTimeout verifica il DESIGN
// del default client: niente `Client.Timeout` (= no overall limit, body
// streaming può durare arbitrariamente) ma `Transport.ResponseHeaderTimeout`
// (TTFB protection).
//
// PRE-FIX: il default usa `Client.Timeout=30min`. Questo test fallisce
// perché Timeout != 0.
// POST-FIX: default = `&http.Client{Transport: &http.Transport{ResponseHeaderTimeout: N}}`.
// Test passa.
func TestOllamaProvider_DefaultClientHasNoOverallTimeout(t *testing.T) {
	p := &provider.OllamaProvider{}
	client := provider.OllamaDefaultHTTPClient(p)
	if client.Timeout != 0 {
		t.Errorf("Client.Timeout deve essere 0 (no overall timeout — body streaming LLM può durare ore), got %v", client.Timeout)
	}
	tr, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("default Transport deve essere *http.Transport, got %T", client.Transport)
	}
	if tr.ResponseHeaderTimeout < 30*time.Minute {
		t.Errorf("Transport.ResponseHeaderTimeout (TTFB) deve essere ≥ 30 min per protezione contro server hung, got %v", tr.ResponseHeaderTimeout)
	}
}

// TestOllamaProvider_HeaderTimeoutStillProtectsAgainstDeadServer verifica che
// `Transport.ResponseHeaderTimeout` continui a proteggere contro server che
// NON emettono mai gli headers (caso "Ollama davvero down").
func TestOllamaProvider_HeaderTimeoutStillProtectsAgainstDeadServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Server "hung": ritardo prima di scrivere QUALSIASI header → triggera
		// ResponseHeaderTimeout se < ritardo.
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	// Costruisci un client con ResponseHeaderTimeout breve (1s) per validare
	// che il fix non rimuova la protezione contro server hanged.
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ResponseHeaderTimeout = 1 * time.Second
	client := &http.Client{Transport: transport} // Timeout: 0 (no overall)

	p := &provider.OllamaProvider{Endpoint: srv.URL, HTTPClient: client}
	_, err := p.Stream(context.Background(), "s", "u", provider.Options{Modello: "test"})
	if err == nil {
		t.Fatal("ResponseHeaderTimeout=1s deve cancellare il server hung che ritarda 3s di header")
	}
}
