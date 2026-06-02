package jobs_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/jobs"
)

// Bug (parte 2): il _labnexus.toml supportava SOLO trigger_prompt_file come
// override, non un trigger_prompt inline. Design CTO: l'override del job è
// XOR (testo O file), simmetrico al profilo.
func TestLoadMetadata_TriggerPromptInline(t *testing.T) {
	dir := t.TempDir()
	metaPath := filepath.Join(dir, "_labnexus.toml")

	t.Run("inline_override_parsato", func(t *testing.T) {
		body := `profile = "revisione"
trigger_prompt = "Override inline dal job, sopra la lunghezza minima di 50 caratteri."`
		if err := os.WriteFile(metaPath, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		m, err := jobs.LoadMetadata(metaPath)
		if err != nil {
			t.Fatalf("LoadMetadata: %v", err)
		}
		if m.TriggerPrompt == "" {
			t.Fatal("TriggerPrompt inline non parsato dal _labnexus.toml")
		}
	})

	t.Run("inline_e_file_insieme_rifiutati_XOR", func(t *testing.T) {
		body := `profile = "revisione"
trigger_prompt = "Override inline dal job, sopra la lunghezza minima di 50 caratteri."
trigger_prompt_file = "Prompt.rtf"`
		if err := os.WriteFile(metaPath, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := jobs.LoadMetadata(metaPath)
		if err == nil {
			t.Fatal("atteso errore XOR: trigger_prompt e trigger_prompt_file insieme nel job")
		}
		if !strings.Contains(strings.ToLower(err.Error()), "mutual") && !strings.Contains(strings.ToLower(err.Error()), "uno solo") {
			t.Errorf("errore XOR poco chiaro: %v", err)
		}
	})

	t.Run("inline_troppo_corto_rifiutato_minLen", func(t *testing.T) {
		body := `profile = "revisione"
trigger_prompt = "troppo corto"`
		if err := os.WriteFile(metaPath, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := jobs.LoadMetadata(metaPath)
		if err == nil {
			t.Fatal("atteso errore minLen: trigger_prompt inline sotto i 50 caratteri")
		}
		if !strings.Contains(strings.ToLower(err.Error()), "corto") {
			t.Errorf("errore minLen poco chiaro: %v", err)
		}
	})

	t.Run("nessun_trigger_accettato_eredita_profilo", func(t *testing.T) {
		body := `profile = "revisione"`
		if err := os.WriteFile(metaPath, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		m, err := jobs.LoadMetadata(metaPath)
		if err != nil {
			t.Fatalf("metadata senza trigger deve essere valido (eredita dal profilo): %v", err)
		}
		if m.TriggerPrompt != "" || m.TriggerPromptFile != "" {
			t.Errorf("attesi trigger vuoti, got inline=%q file=%q", m.TriggerPrompt, m.TriggerPromptFile)
		}
	})

	t.Run("file_path_traversal_rifiutato", func(t *testing.T) {
		body := `profile = "revisione"
trigger_prompt_file = "foo/../../bar.rtf"`
		if err := os.WriteFile(metaPath, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := jobs.LoadMetadata(metaPath)
		if err == nil {
			t.Fatal("atteso errore: trigger_prompt_file con '..' deve essere rifiutato")
		}
		if !strings.Contains(err.Error(), "..") {
			t.Errorf("errore path-traversal poco chiaro: %v", err)
		}
	})
}
