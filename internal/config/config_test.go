package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/labnexus/labnexus/internal/config"
	"github.com/labnexus/labnexus/internal/profile"
)

// Sprint 1.5.B test-scaffold (red phase): tutti i test devono fallire finché
// /v-implement non implementa config.Load + config.Merge. Pattern table-driven
// dove applicabile.

func TestLoad_HappyPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "labnexus.config.toml")
	body := `provider = "eurouter"
modello = "qwen3.6"
temperature = 0.9
max_tokens = 8192
context_window = 128000
eurouter_api_key = "sk-test"
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	m, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: unexpected error %v (red phase expected)", err)
	}
	if m == nil {
		t.Fatal("Load ritorna *Master nil")
	}
	if m.Provider != "eurouter" {
		t.Errorf("Provider: got %q, want eurouter", m.Provider)
	}
	if m.Modello != "qwen3.6" {
		t.Errorf("Modello: got %q, want qwen3.6", m.Modello)
	}
	if m.EurouterAPIKey != "sk-test" {
		t.Errorf("EurouterAPIKey: got %q, want sk-test", m.EurouterAPIKey)
	}
}

func TestLoad_FileMissing(t *testing.T) {
	_, err := config.Load("/nonexistent/labnexus.config.toml")
	if err == nil {
		t.Fatal("Load: expected error per file mancante, got nil")
	}
}

func TestLoad_SyntaxError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "labnexus.config.toml")
	body := `provider = "eurouter"
modello = "qwen3.6"
[malformed section
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("Load: expected error per sintassi TOML errata, got nil")
	}
}

func TestMerge_ProfileCompletoOverridesMaster(t *testing.T) {
	master := &config.Master{Provider: "eurouter", Modello: "qwen3.6"}
	p := &profile.Profile{Provider: "ollama", Modello: "llama3"}
	merged := config.Merge(master, p)
	if merged.Provider != "ollama" {
		t.Errorf("Provider override: got %q, want ollama", merged.Provider)
	}
	if merged.Modello != "llama3" {
		t.Errorf("Modello override: got %q, want llama3", merged.Modello)
	}
}

func TestMerge_ProfileParzialeEreditaMaster(t *testing.T) {
	master := &config.Master{Provider: "eurouter", Modello: "qwen3.6", Temperature: 0.5}
	p := &profile.Profile{}
	merged := config.Merge(master, p)
	if merged.Provider != "eurouter" {
		t.Errorf("Provider eredità: got %q, want eurouter (dal master)", merged.Provider)
	}
	if merged.Modello != "qwen3.6" {
		t.Errorf("Modello eredità: got %q, want qwen3.6 (dal master)", merged.Modello)
	}
	if merged.Temperature != 0.5 {
		t.Errorf("Temperature eredità: got %v, want 0.5 (dal master)", merged.Temperature)
	}
}

func TestMerge_MasterNilProfileInvariato(t *testing.T) {
	p := &profile.Profile{Provider: "ollama"}
	merged := config.Merge(nil, p)
	if merged.Provider != "ollama" {
		t.Errorf("master nil: profile invariato, got %q", merged.Provider)
	}
}

// TestLoad_PDFSectionEnabledTrue: il master TOML può contenere `[pdf]` con
// `enabled = true|false` come default globale per tutte le capability.
func TestLoad_PDFSectionEnabledTrue(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "labnexus.config.toml")
	body := `provider = "eurouter"

[pdf]
enabled = true
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	m, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: unexpected error %v", err)
	}
	if m.PDF.Enabled == nil {
		t.Fatal("PDF.Enabled = nil, want pointer to true (default globale on)")
	}
	if *m.PDF.Enabled != true {
		t.Errorf("PDF.Enabled = %v, want true", *m.PDF.Enabled)
	}
}

// TestLoad_PDFSectionAbsent_EnabledNil: senza sezione [pdf] nel master,
// PDF.Enabled deve essere nil (non setted), così la risoluzione gerarchica
// può distinguere "non setted" da "esplicitamente false".
func TestLoad_PDFSectionAbsent_EnabledNil(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "labnexus.config.toml")
	body := `provider = "eurouter"
modello = "qwen3.5-122b-a10b"
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	m, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: unexpected error %v", err)
	}
	if m.PDF.Enabled != nil {
		t.Errorf("PDF.Enabled dovrebbe essere nil quando [pdf] omesso, got %v", *m.PDF.Enabled)
	}
}
