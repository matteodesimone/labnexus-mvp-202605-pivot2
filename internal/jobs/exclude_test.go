package jobs_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/jobs"
)

func TestLoadMetadata_ParsesExclude(t *testing.T) {
	dir := t.TempDir()
	body := `profile = "revisione"
exclude = ["LEGGIMI.txt", "bozza.docx"]`
	if err := os.WriteFile(filepath.Join(dir, "_labnexus.toml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := jobs.LoadMetadata(filepath.Join(dir, "_labnexus.toml"))
	if err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}
	if len(m.Exclude) != 2 || m.Exclude[0] != "LEGGIMI.txt" || m.Exclude[1] != "bozza.docx" {
		t.Errorf("Exclude = %v, want [LEGGIMI.txt bozza.docx]", m.Exclude)
	}
}

func TestLoadMetadata_NoExcludeIsNil(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "_labnexus.toml"), []byte(`profile = "revisione"`), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := jobs.LoadMetadata(filepath.Join(dir, "_labnexus.toml"))
	if err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}
	if len(m.Exclude) != 0 {
		t.Errorf("Exclude deve essere vuoto, got %v", m.Exclude)
	}
}

// Il toml generato dal wizard contiene un commento-suggerimento `# exclude`:
// deve restare valido (TOML ignora i commenti) e non popolare Exclude.
func TestWriteMetadata_ExcludeHintIsCommentedNotActive(t *testing.T) {
	dir := t.TempDir()
	if err := jobs.WriteMetadata(dir, "revisione", ""); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(filepath.Join(dir, "_labnexus.toml"))
	if !strings.Contains(string(raw), "# exclude") {
		t.Error("il toml generato deve contenere il suggerimento commentato # exclude")
	}
	m, err := jobs.LoadMetadata(filepath.Join(dir, "_labnexus.toml"))
	if err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}
	if len(m.Exclude) != 0 {
		t.Errorf("il suggerimento è commentato: Exclude deve essere vuoto, got %v", m.Exclude)
	}
}
