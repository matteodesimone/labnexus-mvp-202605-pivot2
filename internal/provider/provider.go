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
	ErrOllamaUnreachable = errors.New("Ollama non raggiungibile su localhost:11434 — verifica che ollama serve sia attivo")
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
	Done         bool
	Err          error
	NoDoneMarker bool // true quando Done=true emesso per EOF post-content (no done:true esplicito dal server)
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
