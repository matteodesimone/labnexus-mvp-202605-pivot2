package profile_test

import (
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/profile"
)

// validProfile è un profilo valido di riferimento per i test.
func validProfile() *profile.Profile {
	return &profile.Profile{
		Profilo:       "revisione",
		Descrizione:   "Revisione documentale",
		Provider:      "ollama",
		Modello:       "qwen3.6",
		KbFiles:       []string{"CLAUDE.md"},
		TriggerPrompt: strings.Repeat("a", 100),
		Temperature:   0.9,
		MaxTokens:     8192,
		ContextWindow: 128000,
	}
}

func TestValidate_AcceptsValid(t *testing.T) {
	p := validProfile()
	// In red phase: Validate ritorna ErrNotImplemented → questo test fallisce, atteso.
	if err := profile.Validate(p, "testdata/kb"); err != nil {
		t.Fatalf("expected nil for valid profile, got %v", err)
	}
}

func TestValidate_RejectsShortTrigger(t *testing.T) {
	p := validProfile()
	p.TriggerPrompt = "ciao" // < MinTriggerLen
	err := profile.Validate(p, "testdata/kb")
	if err == nil {
		t.Fatalf("expected error for trigger_prompt < %d chars, got nil", profile.MinTriggerLen)
	}
}

func TestValidate_RejectsUnknownProvider(t *testing.T) {
	p := validProfile()
	p.Provider = "openai"
	err := profile.Validate(p, "testdata/kb")
	if err == nil {
		t.Fatalf("expected error for unknown provider, got nil")
	}
}

func TestValidate_RejectsMissingKbFile(t *testing.T) {
	p := validProfile()
	p.KbFiles = []string{"non-esiste.md"}
	err := profile.Validate(p, "testdata/kb")
	if err == nil {
		t.Fatalf("expected error for missing kb file, got nil")
	}
}

func TestValidate_TableDriven(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*profile.Profile)
		wantErr bool
	}{
		{"valid", func(p *profile.Profile) {}, false},
		{"empty profilo field", func(p *profile.Profile) { p.Profilo = "" }, true},
		{"empty descrizione", func(p *profile.Profile) { p.Descrizione = "" }, true},
		{"empty provider", func(p *profile.Profile) { p.Provider = "" }, true},
		{"empty modello", func(p *profile.Profile) { p.Modello = "" }, true},
		{"empty kb_files", func(p *profile.Profile) { p.KbFiles = []string{} }, true},
		{"trigger exactly 49 chars", func(p *profile.Profile) { p.TriggerPrompt = strings.Repeat("a", 49) }, true},
		{"trigger exactly 50 chars", func(p *profile.Profile) { p.TriggerPrompt = strings.Repeat("a", 50) }, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := validProfile()
			tc.mutate(p)
			err := profile.Validate(p, "testdata/kb")
			gotErr := err != nil
			if gotErr != tc.wantErr {
				t.Fatalf("case %q: wantErr=%v gotErr=%v (err=%v)", tc.name, tc.wantErr, gotErr, err)
			}
		})
	}
}
