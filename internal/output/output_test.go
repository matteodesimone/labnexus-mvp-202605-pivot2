package output_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/output"
)

func TestWrite_CreatesFileWithNamingConvention(t *testing.T) {
	dir := t.TempDir()
	fm := &output.Frontmatter{
		Profilo:        "revisione",
		Modello:        "qwen3.6",
		Provider:       "ollama",
		DataEsecuzione: "2026-05-18T14:10:22+02:00",
	}
	path, err := output.Write(dir, fm, "body content")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	base := filepath.Base(path)
	if !strings.Contains(base, "revisione") {
		t.Errorf("filename must contain profile name 'revisione', got %q", base)
	}
	if !strings.HasSuffix(base, ".md") {
		t.Errorf("filename must end with .md, got %q", base)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file should exist at %q: %v", path, err)
	}
}

func TestWrite_CollisionGetsIncrementalSuffix(t *testing.T) {
	dir := t.TempDir()
	fm := &output.Frontmatter{Profilo: "revisione", DataEsecuzione: "2026-05-18T14:10:22+02:00"}
	first, err := output.Write(dir, fm, "body1")
	if err != nil {
		t.Fatalf("first write: %v", err)
	}
	second, err := output.Write(dir, fm, "body2")
	if err != nil {
		t.Fatalf("second write: %v", err)
	}
	if first == second {
		t.Fatal("second write should not overwrite first (EC-15)")
	}
	if !strings.Contains(filepath.Base(second), "_2") {
		t.Errorf("second file should have _2 suffix, got %q", second)
	}
}

func TestWrite_FrontmatterContainsAllFields(t *testing.T) {
	dir := t.TempDir()
	fm := &output.Frontmatter{
		Profilo:        "rilievi",
		Modello:        "qwen3.6",
		Provider:       "ollama",
		DataEsecuzione: "2026-05-18T14:10:22+02:00",
		DurataSecondi:  42.3,
		TokenStimati:   1234,
		FileInput:      []string{"a.csv"},
	}
	path, err := output.Write(dir, fm, "body")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, _ := os.ReadFile(path)
	s := string(content)
	for _, want := range []string{"profilo: rilievi", "modello: qwen3.6", "provider: ollama", "durata_secondi: 42.3", "token_stimati: 1234"} {
		if !strings.Contains(s, want) {
			t.Errorf("frontmatter missing %q in:\n%s", want, s)
		}
	}
}
