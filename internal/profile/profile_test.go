package profile_test

import (
	"os"
	"path/filepath"
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

// --- Bug #002: KB symlink escape ----------------------------------------
// Un kb_file che è un symlink (anche se il path lessicale è dentro kbDir)
// può puntare OUTSIDE la cartella KB, esfiltrando file privati nel prompt
// LLM. Vedi .pipeline/bugs/kb-symlink-escape.md.

func TestValidate_RejectsKbFileSymlink(t *testing.T) {
	// Setup: kbDir tmp con CLAUDE.md regular + evil-link symlink → /tmp/secret
	kbDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(kbDir, "CLAUDE.md"), []byte("# identità ispettore"), 0o644); err != nil {
		t.Fatalf("setup CLAUDE.md: %v", err)
	}
	// File "esterno" da cui un symlink potrebbe esfiltrare contenuto.
	secretDir := t.TempDir()
	secretFile := filepath.Join(secretDir, "secret.md")
	if err := os.WriteFile(secretFile, []byte("contenuto privato"), 0o600); err != nil {
		t.Fatalf("setup secret: %v", err)
	}
	// Symlink dentro kbDir che punta fuori.
	evilLink := filepath.Join(kbDir, "evil-link.md")
	if err := os.Symlink(secretFile, evilLink); err != nil {
		t.Fatalf("setup symlink: %v", err)
	}

	p := validProfile()
	p.KbFiles = []string{"CLAUDE.md", "evil-link.md"}
	err := profile.Validate(p, kbDir)
	if err == nil {
		t.Fatal("atteso errore: kb_file symlink che punta fuori da kbDir deve essere rifiutato")
	}
	if !strings.Contains(err.Error(), "symlink") && !strings.Contains(err.Error(), "esce") {
		t.Errorf("messaggio errore deve menzionare symlink/escape, got: %v", err)
	}
}

func TestValidate_AcceptsKbDirItselfSymlink(t *testing.T) {
	// kbDir esso stesso può essere un symlink (caso reale: ./KB-ispettore →
	// docs/.../KB-ispettore). Questo è OK: il check post-resolve deve solo
	// verificare che i kb_files restino dentro la kbDir risolta.
	realKb := t.TempDir()
	if err := os.WriteFile(filepath.Join(realKb, "CLAUDE.md"), []byte("# identità"), 0o644); err != nil {
		t.Fatalf("setup CLAUDE.md: %v", err)
	}
	linkParent := t.TempDir()
	kbLink := filepath.Join(linkParent, "KB-ispettore")
	if err := os.Symlink(realKb, kbLink); err != nil {
		t.Fatalf("setup kb link: %v", err)
	}
	p := validProfile()
	p.KbFiles = []string{"CLAUDE.md"}
	if err := profile.Validate(p, kbLink); err != nil {
		t.Fatalf("kbDir symlink valido deve passare, got: %v", err)
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
