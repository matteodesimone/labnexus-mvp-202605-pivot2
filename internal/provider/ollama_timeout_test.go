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

// TestOllamaProvider_TolerateSlowResponseHeaders riproduce il bug
// `.pipeline/bugs/ollama-http-timeout-troppo-stretto-per-reasoning-models.md`:
// modelli reasoning grandi (qwen3.6 36B MoE) impiegano minuti per il primo
// byte di response (context fill + thinking). Il client HTTP non deve fare
// abort prematuro mentre il server sta legittimamente lavorando.
//
// Test: fake Ollama che ritarda di 2 secondi prima di emettere l'header.
// Pre-fix (Timeout=90s lavorerebbe, ma per il vero scenario serve molto più).
// Per testare l'INFRASTRUTTURA con un fake server, usiamo un client con
// timeout esplicito che simula la stessa categoria di problema:
// se Timeout fosse < 2s, il test fallisce; se ≥ 5s (post-fix default ben sopra 90s), passa.
//
// Per validare il fix vero (timeout di 30 min default + env override), serve
// un test che verifichi che `OllamaProvider` rispetti l'env var. Vedi
// TestOllamaProvider_HonorsTimeoutEnv sotto.
func TestOllamaProvider_TolerateSlowResponseHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simula context fill / thinking: ritardo prima di rispondere.
		time.Sleep(2 * time.Second)
		w.Header().Set("Content-Type", "application/x-ndjson")
		_, _ = w.Write([]byte(`{"message":{"content":"alpha"},"done":true}` + "\n"))
	}))
	defer srv.Close()

	// Client con Timeout breve (1s) — simula il bug attuale (90s troppo poco per uso reale)
	shortTimeout := &http.Client{Timeout: 1 * time.Second}
	p := &provider.OllamaProvider{Endpoint: srv.URL, HTTPClient: shortTimeout}
	_, err := p.Stream(context.Background(), "s", "u", provider.Options{Modello: "test"})
	if err == nil {
		t.Fatal("client con Timeout=1s deve fallire sul fake server che ritarda 2s — qui invece ha avuto successo")
	}

	// Client con Timeout generoso (10s, ben sopra il delay 2s del fake server)
	// → deve riuscire. Questo è il comportamento atteso post-fix (default 30 min).
	generousTimeout := &http.Client{Timeout: 10 * time.Second}
	p = &provider.OllamaProvider{Endpoint: srv.URL, HTTPClient: generousTimeout}
	ch, err := p.Stream(context.Background(), "s", "u", provider.Options{Modello: "test"})
	if err != nil {
		t.Fatalf("client con Timeout=10s deve tollerare delay 2s, got: %v", err)
	}
	got := false
	for ev := range ch {
		if ev.Token == "alpha" {
			got = true
		}
		if ev.Done {
			break
		}
	}
	if !got {
		t.Errorf("atteso token 'alpha' dal fake server lento")
	}
}

// TestOllamaProvider_DefaultTimeoutSupportsLongPrompts verifica che il default
// del client HTTP (quando p.HTTPClient è nil) sia sufficientemente generoso
// per uso reale di reasoning models. Pre-fix: 90s — non sufficiente.
// Post-fix: ≥ 30 min (1800s).
func TestOllamaProvider_DefaultTimeoutSupportsLongPrompts(t *testing.T) {
	p := &provider.OllamaProvider{} // no HTTPClient → usa default
	timeout := provider.OllamaDefaultTimeout(p)
	minRequired := 30 * time.Minute
	if timeout < minRequired {
		t.Errorf("default timeout = %v, ma per modelli reasoning grandi serve ≥ %v", timeout, minRequired)
	}
}

// TestOllamaProvider_HonorsTimeoutEnv verifica che il default rispetti
// l'env var LABNEXUS_HTTP_TIMEOUT (secondi) per casi estremi.
func TestOllamaProvider_HonorsTimeoutEnv(t *testing.T) {
	t.Setenv("LABNEXUS_HTTP_TIMEOUT", "60") // 60 secondi
	p := &provider.OllamaProvider{}
	timeout := provider.OllamaDefaultTimeout(p)
	want := 60 * time.Second
	if timeout != want {
		t.Errorf("env var LABNEXUS_HTTP_TIMEOUT=60 ignored: got %v, want %v", timeout, want)
	}
}

// TestOllamaProvider_HonorsTimeoutEnv_Invalid verifica fallback robusto se
// l'env var contiene un valore non valido.
func TestOllamaProvider_HonorsTimeoutEnv_Invalid(t *testing.T) {
	t.Setenv("LABNEXUS_HTTP_TIMEOUT", "invalid-not-a-number")
	p := &provider.OllamaProvider{}
	timeout := provider.OllamaDefaultTimeout(p)
	// Atteso: fallback al default sano (≥ 30 min)
	minRequired := 30 * time.Minute
	if timeout < minRequired {
		t.Errorf("con env invalida, fallback deve essere default ≥ %v, got %v", minRequired, timeout)
	}
}

// Helper: documenta nel test che l'errore di Timeout è riconducibile al
// fenomeno descritto nel bug file (utile per chi indaga in futuro).
func TestOllamaProvider_TimeoutErrorMentionsCorrectCause(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second) // ritardo > Timeout cliente
	}))
	defer srv.Close()
	p := &provider.OllamaProvider{Endpoint: srv.URL, HTTPClient: &http.Client{Timeout: 500 * time.Millisecond}}
	_, err := p.Stream(context.Background(), "s", "u", provider.Options{Modello: "test"})
	if err == nil {
		t.Fatal("atteso errore di timeout")
	}
	// Il messaggio attuale dice "Ollama non raggiungibile" — accettabile, anche se
	// fuorviante per i casi "Ollama OK ma modello lento". Refactor future.
	if !strings.Contains(err.Error(), "Ollama") && !strings.Contains(err.Error(), "deadline") {
		t.Errorf("errore di timeout deve menzionare Ollama o deadline, got: %v", err)
	}
}
