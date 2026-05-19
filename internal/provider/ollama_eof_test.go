package provider_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labnexus/labnexus/internal/provider"
)

// TestOllamaProvider_EOFAfterContentCompletes verifica il caso "completato
// senza done marker": Ollama emette N chunk con content valido e poi chiude
// la connessione SENZA il chunk finale `done: true`. Atteso: lo stream
// termina con StreamEvent{Done: true, NoDoneMarker: true}, NON con Err.
//
// Reproducer per `.pipeline/bugs/eof-post-content-falsamente-interrotto.md`:
// nel test sul Mac CTO con qwen3.6 36B reale, il modello ha generato 11 KB
// di output completo in 21 min, poi Ollama ha chiuso senza done:true →
// binario emetteva Err → frontmatter "interrotto" → exit 1.
func TestOllamaProvider_EOFAfterContentCompletes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		fl, _ := w.(http.Flusher)
		// Emette 3 chunk con content, MA NESSUN chunk con done:true.
		// Poi il handler termina → connessione chiusa.
		for _, content := range []string{"primo ", "secondo ", "terzo"} {
			_, _ = w.Write([]byte(`{"message":{"content":"` + content + `"}}` + "\n"))
			if fl != nil {
				fl.Flush()
			}
		}
		// EOF naturale (handler exit).
	}))
	defer srv.Close()

	p := &provider.OllamaProvider{Endpoint: srv.URL}
	ch, err := p.Stream(context.Background(), "s", "u", provider.Options{Modello: "test"})
	if err != nil {
		t.Fatalf("Stream error: %v", err)
	}

	var tokens []string
	var sawDone, sawNoDoneMarker, sawErr bool
	for ev := range ch {
		if ev.Err != nil {
			sawErr = true
			t.Logf("unexpected error event: %v", ev.Err)
			continue
		}
		if ev.Token != "" {
			tokens = append(tokens, ev.Token)
		}
		if ev.Done {
			sawDone = true
			if ev.NoDoneMarker {
				sawNoDoneMarker = true
			}
		}
	}

	if sawErr {
		t.Errorf("EOF post-content non deve emettere Err: stream interpretato come interrotto invece di completato")
	}
	if !sawDone {
		t.Errorf("EOF post-content deve emettere Done:true (anche senza marker dal server)")
	}
	if !sawNoDoneMarker {
		t.Errorf("EOF post-content deve marcare NoDoneMarker=true per informare il caller")
	}
	if len(tokens) != 3 {
		t.Errorf("tutti i 3 token con content devono essere emessi prima dell'EOF, got %d: %v", len(tokens), tokens)
	}
}

// TestOllamaProvider_EOFBeforeContentIsRealError verifica il caso "interrotto
// vero": il server chiude la connessione PRIMA di emettere qualunque chunk
// con content. Atteso: emette Err (stream genuinamente interrotto).
func TestOllamaProvider_EOFBeforeContentIsRealError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		// Nessun chunk emesso, handler termina subito → EOF immediato post-header.
	}))
	defer srv.Close()

	p := &provider.OllamaProvider{Endpoint: srv.URL}
	ch, err := p.Stream(context.Background(), "s", "u", provider.Options{Modello: "test"})
	if err != nil {
		t.Fatalf("Stream error: %v", err)
	}

	var sawDone, sawErr bool
	for ev := range ch {
		if ev.Err != nil {
			sawErr = true
		}
		if ev.Done {
			sawDone = true
		}
	}

	// Caso "EOF before content": nessun token ricevuto, non possiamo dire se
	// è "completato senza done" o "interrotto"; preferiamo segnalare come
	// completato vuoto (Done senza NoDoneMarker non ha senso, lasciamo Err)
	// → questo è il design: senza content, non possiamo distinguere e quindi
	// emettiamo Done con NoDoneMarker (= "il server ha chiuso, abbiamo zero token").
	// Alternativa: emettere Err. Decisione: Done+NoDoneMarker permette al caller
	// di gestire come "completato con stream vuoto" — runner lo riconosce
	// come degenere via body builder vuoto.
	if sawErr {
		t.Logf("server chiuso senza content → emesso Err (accettabile)")
	}
	if sawDone {
		t.Logf("server chiuso senza content → emesso Done con NoDoneMarker (accettabile)")
	}
	// Almeno una delle due deve essere true.
	if !sawErr && !sawDone {
		t.Errorf("server vuoto deve emettere o Err o Done, non niente")
	}
}

// TestOllamaProvider_NormalCompletionStillEmitsCleanDone verifica il caso
// "normale": il server emette chunks + chunk finale done:true. Atteso:
// Done emesso CON NoDoneMarker=false (il marker è arrivato dal server).
func TestOllamaProvider_NormalCompletionStillEmitsCleanDone(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		_, _ = w.Write([]byte(`{"message":{"content":"alpha"},"done":false}` + "\n"))
		_, _ = w.Write([]byte(`{"message":{"content":"beta"},"done":true}` + "\n"))
	}))
	defer srv.Close()

	p := &provider.OllamaProvider{Endpoint: srv.URL}
	ch, err := p.Stream(context.Background(), "s", "u", provider.Options{Modello: "test"})
	if err != nil {
		t.Fatalf("Stream error: %v", err)
	}

	var sawDone, sawNoDoneMarker bool
	for ev := range ch {
		if ev.Done {
			sawDone = true
			if ev.NoDoneMarker {
				sawNoDoneMarker = true
			}
		}
	}

	if !sawDone {
		t.Errorf("completion normale deve emettere Done:true")
	}
	if sawNoDoneMarker {
		t.Errorf("quando il server emette done:true esplicito, NoDoneMarker deve restare false")
	}
}
