package provider_test

import (
	"testing"

	"github.com/labnexus/labnexus/internal/provider"
)

func TestSelect_FlagWinsOverEnvAndProfile(t *testing.T) {
	p, err := provider.Select("ollama", "eurouter", "eurouter")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "ollama" {
		t.Errorf("flag should win: got %s, want ollama", p.Name())
	}
}

func TestSelect_EnvWinsOverProfile(t *testing.T) {
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
