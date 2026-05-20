// Package provider definisce l'astrazione LLMProvider (Strategy pattern, FR-7)
// e la funzione Select per scegliere quale impl usare secondo precedenza
// flag > env > profilo (FR-12).
package provider

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

var (
	ErrMissingAPIKey     = errors.New("EUROUTER_API_KEY non impostata")
	ErrOllamaUnreachable = errors.New("Ollama non raggiungibile su localhost:11434 — verifica che ollama serve sia attivo")
	// ErrEurouterGateMissing è restituito quando si tenta di selezionare il
	// provider eurouter senza il gate esplicito LABNEXUS_ALLOW_CLOUD_PROVIDER.
	// Bug #001 (.pipeline/bugs/privacy-eurouter-gate-mancante.md): protegge
	// dati SGQ da esfiltrazione accidentale al gateway cloud (NFR-1).
	ErrEurouterGateMissing = errors.New(`provider eurouter non ammesso senza approval esplicito.
Imposta l'env var LABNEXUS_ALLOW_CLOUD_PROVIDER=approved-for-synthetic-data
per autorizzare il provider cloud. SU DATI REALI DEL SGQ DI DENIS USARE OLLAMA LOCALE,
NON EUROUTER. EUrouter è strumento di debug interno su dati sintetici/anonimizzati.`)
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
		// Bug #003 gate: l'endpoint Ollama deve essere loopback per
		// rispettare la promise "provider locale = on-device". Override a
		// host remoti richiede LABNEXUS_ALLOW_CLOUD_PROVIDER (stesso gate
		// di eurouter, bug #001). .pipeline/bugs/ollama-endpoint-env-override-leak.md
		if err := requireLoopbackOrCloudGate(endpoint); err != nil {
			return nil, err
		}
		return &OllamaProvider{Endpoint: endpoint}, nil
	case "eurouter":
		// Bug #001 gate: eurouter è ammesso solo con approval esplicito.
		if strings.TrimSpace(os.Getenv("LABNEXUS_ALLOW_CLOUD_PROVIDER")) == "" {
			return nil, ErrEurouterGateMissing
		}
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

// requireLoopbackOrCloudGate verifica che l'endpoint URL fornito risolva a
// un host loopback (127.0.0.0/8, ::1, localhost). Se non loopback, è ammesso
// solo se LABNEXUS_ALLOW_CLOUD_PROVIDER è settato (gate cloud esplicito, lo
// stesso usato per eurouter). Bug #003 fix.
func requireLoopbackOrCloudGate(endpoint string) error {
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("provider: endpoint Ollama %q non parsabile: %w", endpoint, err)
	}
	host := u.Hostname()
	if isLoopbackHost(host) {
		return nil
	}
	if strings.TrimSpace(os.Getenv("LABNEXUS_ALLOW_CLOUD_PROVIDER")) != "" {
		return nil
	}
	return fmt.Errorf(`LABNEXUS_OLLAMA_ENDPOINT punta a host non-loopback %q.
Il provider "ollama" deve restare locale (NFR-1, promise post-pivot 2).
Per usare un endpoint Ollama remoto in scenari di dev avanzati, imposta anche
LABNEXUS_ALLOW_CLOUD_PROVIDER=approved-for-synthetic-data. SU DATI REALI DEL
SGQ USARE SOLO OLLAMA LOCALE`, host)
}

// isLoopbackHost ritorna true se host è un indirizzo loopback o "localhost".
func isLoopbackHost(host string) bool {
	if host == "" || host == "localhost" {
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback()
	}
	return false
}
