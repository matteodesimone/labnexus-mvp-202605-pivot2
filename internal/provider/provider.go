// Package provider definisce l'astrazione LLMProvider (Strategy pattern, FR-7)
// e la funzione Select per scegliere quale impl usare secondo precedenza
// flag > env > profilo (FR-12).
package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrMissingAPIKey     = errors.New("EUROUTER_API_KEY non impostata")
	ErrOllamaUnreachable = errors.New("Ollama non raggiungibile su localhost:11434 — verifica che 'ollama serve' sia attivo (provider opzionale post-Sprint-1.5; il default è eurouter cloud)")
)

// Options sono i parametri di chiamata.
type Options struct {
	Modello       string
	Temperature   float64
	MaxTokens     int
	ContextWindow int
}

// StreamEvent rappresenta un singolo evento dello stream.
type StreamEvent struct {
	Token        string
	Reasoning    string // delta di reasoning_content (modelli "thinking" es. kimi): diagnostico, NON va nell'output finale
	Done         bool
	Err          error
	NoDoneMarker bool   // true quando Done=true emesso per EOF post-content (no done:true esplicito dal server)
	FinishReason string // "stop"|"length"|"content_filter"|... — perché lo stream si è chiuso (diagnostica empty/troncato)
	Usage        *Usage // token reali riportati dal provider (richiede stream_options.include_usage)
}

// Usage sono i conteggi token reali del provider (meta del modello). Per i
// modelli "thinking" ReasoningTokens dice quanto è stato speso a ragionare.
type Usage struct {
	PromptTokens     int
	CompletionTokens int
	ReasoningTokens  int
	TotalTokens      int
}

// LLMProvider è l'interface unica di chiamata.
type LLMProvider interface {
	Stream(ctx context.Context, system string, user string, opts Options) (<-chan StreamEvent, error)
	Name() string
}

// Select risolve quale provider istanziare secondo la precedenza:
//
//	flag --provider > env LABNEXUS_PROVIDER > campo profile.Provider.
//
// EUROUTER_API_KEY è letta da env solo per il provider eurouter.
//
// Sprint 1.5.A (post-pivot-3): nessun gate cloud runtime. La decisione di usare
// eurouter (cloud EU-GDPR) vs ollama (locale) passa esclusivamente per la
// configurazione utente (profilo + futuro config master TOML in 1.5.B).
// Cliente (Stefano Fiorina) formalmente informato e d'accordo al pivot.
// Vedi `.pipeline/bugs/privacy-eurouter-gate-mancante.md` (intentional_deviation_post_pivot_3).
func Select(providerOverride, providerEnv, providerFromProfile string) (LLMProvider, error) {
	chosen := strings.TrimSpace(providerOverride)
	if chosen == "" {
		chosen = strings.TrimSpace(providerEnv)
	}
	if chosen == "" {
		chosen = strings.TrimSpace(providerFromProfile)
	}
	switch chosen {
	case "ollama":
		endpoint := os.Getenv("LABNEXUS_OLLAMA_ENDPOINT")
		if endpoint == "" {
			endpoint = "http://localhost:11434"
		}
		return &OllamaProvider{Endpoint: endpoint}, nil
	case "eurouter":
		endpoint := os.Getenv("LABNEXUS_EUROUTER_ENDPOINT")
		if endpoint == "" {
			endpoint = "https://api.eurouter.ai/v1/chat/completions"
		}
		return &EurouterProvider{
			Endpoint: endpoint,
			APIKey:   os.Getenv("EUROUTER_API_KEY"),
			Timeout:  120,
		}, nil
	case "":
		return nil, errors.New("provider: nessun provider scelto (override/env/profilo tutti vuoti)")
	default:
		return nil, fmt.Errorf("provider: %q non ammesso (consentiti: ollama, eurouter)", chosen)
	}
}
