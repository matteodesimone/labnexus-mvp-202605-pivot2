package jobs_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/jobs"
)

func TestWriteMetadata_Roundtrip(t *testing.T) {
	t.Run("solo_profilo_default", func(t *testing.T) {
		dir := t.TempDir()
		if err := jobs.WriteMetadata(dir, "revisione", ""); err != nil {
			t.Fatalf("WriteMetadata: %v", err)
		}
		m, err := jobs.LoadMetadata(filepath.Join(dir, "_labnexus.toml"))
		if err != nil {
			t.Fatalf("LoadMetadata: %v", err)
		}
		if m.Profile != "revisione" {
			t.Errorf("Profile = %q, want revisione", m.Profile)
		}
		if m.TriggerPromptFile != "" {
			t.Errorf("TriggerPromptFile deve essere vuoto, got %q", m.TriggerPromptFile)
		}
	})

	t.Run("con_trigger_prompt_file", func(t *testing.T) {
		dir := t.TempDir()
		if err := jobs.WriteMetadata(dir, "rilievi", "Prompt_INPUT.txt"); err != nil {
			t.Fatalf("WriteMetadata: %v", err)
		}
		m, err := jobs.LoadMetadata(filepath.Join(dir, "_labnexus.toml"))
		if err != nil {
			t.Fatalf("LoadMetadata: %v", err)
		}
		if m.Profile != "rilievi" || m.TriggerPromptFile != "Prompt_INPUT.txt" {
			t.Errorf("got profile=%q file=%q, want rilievi/Prompt_INPUT.txt", m.Profile, m.TriggerPromptFile)
		}
	})

	t.Run("profilo_vuoto_errore", func(t *testing.T) {
		if err := jobs.WriteMetadata(t.TempDir(), "  ", ""); err == nil {
			t.Fatal("atteso errore con profilo vuoto")
		}
	})

	t.Run("prompt_file_traversal_rifiutato", func(t *testing.T) {
		err := jobs.WriteMetadata(t.TempDir(), "revisione", "../fuori.txt")
		if err == nil || !strings.Contains(err.Error(), "..") {
			t.Fatalf("atteso errore path-traversal, got %v", err)
		}
		_ = os.Remove("x")
	})
}
