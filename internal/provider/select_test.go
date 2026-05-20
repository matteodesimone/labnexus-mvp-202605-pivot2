package provider_test

import (
	"testing"

	"github.com/labnexus/labnexus/internal/provider"
)

// allowEurouterEnv è l'helper di test che imposta il gate cloud per i test
// che esplicitamente vogliono testare la selezione di eurouter. NFR-1 / bug
// #001 todo: eurouter è selezionabile solo con LABNEXUS_ALLOW_CLOUD_PROVIDER
// settato (gate esplicito per prevenire esfiltrazione dati SGQ).
func allowEurouterEnv(t *testing.T) {
	t.Helper()
	t.Setenv("LABNEXUS_ALLOW_CLOUD_PROVIDER", "approved-for-synthetic-data")
}

func TestSelect_FlagWinsOverEnvAndProfile(t *testing.T) {
	allowEurouterEnv(t)
	p, err := provider.Select("ollama", "eurouter", "eurouter")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "ollama" {
		t.Errorf("flag should win: got %s, want ollama", p.Name())
	}
}

func TestSelect_EnvWinsOverProfile(t *testing.T) {
	allowEurouterEnv(t)
	p, err := provider.Select("", "eurouter", "ollama")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "eurouter" {
		t.Errorf("env should win over profile: got %s, want eurouter", p.Name())
	}
}

func TestSelect_FallbackToProfile(t *testing.T) {
	p, err := provider.Select("", "", "ollama")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "ollama" {
		t.Errorf("profile fallback: got %s, want ollama", p.Name())
	}
}

func TestSelect_RejectsUnknownProvider(t *testing.T) {
	_, err := provider.Select("openai", "", "ollama")
	if err == nil {
		t.Fatal("expected error for unknown provider 'openai', got nil")
	}
}

// --- Bug #001: eurouter requires explicit approval gate -------------------
// Senza LABNEXUS_ALLOW_CLOUD_PROVIDER settato, qualunque tentativo di
// selezionare eurouter (via flag, env, profilo) DEVE essere rifiutato per
// prevenire esfiltrazione accidentale di dati SGQ al gateway cloud
// (NFR-1, .pipeline/bugs/privacy-eurouter-gate-mancante.md).

func TestSelect_EurouterRejectedWithoutApprovalGate(t *testing.T) {
	// Nessun setenv: LABNEXUS_ALLOW_CLOUD_PROVIDER non impostato.
	t.Setenv("LABNEXUS_ALLOW_CLOUD_PROVIDER", "")
	_, err := provider.Select("eurouter", "", "")
	if err == nil {
		t.Fatal("atteso errore selezionando eurouter senza LABNEXUS_ALLOW_CLOUD_PROVIDER, got nil — gate mancante")
	}
}

func TestSelect_EurouterRejectedViaEnvWithoutApprovalGate(t *testing.T) {
	t.Setenv("LABNEXUS_ALLOW_CLOUD_PROVIDER", "")
	_, err := provider.Select("", "eurouter", "")
	if err == nil {
		t.Fatal("atteso errore selezionando eurouter via env senza approval gate, got nil")
	}
}

func TestSelect_EurouterRejectedViaProfileWithoutApprovalGate(t *testing.T) {
	t.Setenv("LABNEXUS_ALLOW_CLOUD_PROVIDER", "")
	_, err := provider.Select("", "", "eurouter")
	if err == nil {
		t.Fatal("atteso errore selezionando eurouter via profilo senza approval gate, got nil")
	}
}

func TestSelect_EurouterAcceptedWithApprovalGate(t *testing.T) {
	allowEurouterEnv(t)
	p, err := provider.Select("eurouter", "", "")
	if err != nil {
		t.Fatalf("unexpected error con gate attivo: %v", err)
	}
	if p.Name() != "eurouter" {
		t.Errorf("eurouter atteso con gate, got %s", p.Name())
	}
}

// --- Bug #003: Ollama endpoint must be loopback ---------------------------
// LABNEXUS_OLLAMA_ENDPOINT override può deviare il "provider locale" a host
// remoto silenziosamente (frontmatter dichiarerebbe provider: ollama mentre
// la rete va al cloud). NFR-1 / .pipeline/bugs/ollama-endpoint-env-override-leak.md.

func TestSelect_OllamaRejectsRemoteEndpoint(t *testing.T) {
	t.Setenv("LABNEXUS_ALLOW_CLOUD_PROVIDER", "")
	t.Setenv("LABNEXUS_OLLAMA_ENDPOINT", "https://ollama-cloud.example.com/api")
	_, err := provider.Select("ollama", "", "")
	if err == nil {
		t.Fatal("atteso errore: LABNEXUS_OLLAMA_ENDPOINT non-loopback deve essere rifiutato senza approval gate")
	}
}

func TestSelect_OllamaAcceptsLoopbackEndpoint(t *testing.T) {
	t.Setenv("LABNEXUS_OLLAMA_ENDPOINT", "http://127.0.0.1:54321")
	p, err := provider.Select("ollama", "", "")
	if err != nil {
		t.Fatalf("loopback 127.0.0.1 deve essere ammesso, got: %v", err)
	}
	if p.Name() != "ollama" {
		t.Fatalf("got %s, want ollama", p.Name())
	}
}

func TestSelect_OllamaAcceptsLocalhostEndpoint(t *testing.T) {
	t.Setenv("LABNEXUS_OLLAMA_ENDPOINT", "http://localhost:11434")
	p, err := provider.Select("ollama", "", "")
	if err != nil {
		t.Fatalf("loopback 'localhost' deve essere ammesso, got: %v", err)
	}
	if p.Name() != "ollama" {
		t.Fatalf("got %s, want ollama", p.Name())
	}
}

func TestSelect_OllamaAcceptsIPv6LoopbackEndpoint(t *testing.T) {
	t.Setenv("LABNEXUS_OLLAMA_ENDPOINT", "http://[::1]:11434")
	p, err := provider.Select("ollama", "", "")
	if err != nil {
		t.Fatalf("loopback IPv6 ::1 deve essere ammesso, got: %v", err)
	}
	if p.Name() != "ollama" {
		t.Fatalf("got %s, want ollama", p.Name())
	}
}

func TestSelect_OllamaRemoteAcceptedWithCloudGate(t *testing.T) {
	// Edge case: se l'utente ha esplicitamente autorizzato il cloud provider,
	// può anche puntare LABNEXUS_OLLAMA_ENDPOINT remoto (es. ollama gateway
	// privato in cloud, scenari di dev avanzati). Stesso gate di #001.
	allowEurouterEnv(t)
	t.Setenv("LABNEXUS_OLLAMA_ENDPOINT", "https://ollama-cloud.example.com/api")
	p, err := provider.Select("ollama", "", "")
	if err != nil {
		t.Fatalf("con LABNEXUS_ALLOW_CLOUD_PROVIDER, endpoint remoto deve essere ammesso, got: %v", err)
	}
	if p.Name() != "ollama" {
		t.Fatalf("got %s, want ollama", p.Name())
	}
}

func TestSelect_OllamaDefaultEndpointStillLocalhost(t *testing.T) {
	t.Setenv("LABNEXUS_OLLAMA_ENDPOINT", "")
	p, err := provider.Select("ollama", "", "")
	if err != nil {
		t.Fatalf("default ollama endpoint non deve essere bloccato: %v", err)
	}
	if p.Name() != "ollama" {
		t.Fatalf("got %s, want ollama", p.Name())
	}
}
