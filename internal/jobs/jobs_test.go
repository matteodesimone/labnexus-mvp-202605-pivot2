package jobs_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/jobs"
)

// Sprint 1.5.C test-scaffold (red phase): tutti i test devono fallire finché
// /v-implement non implementa Discover, Resolve, LoadMetadata.

func TestDiscover_LavoriEmptyReturnsEmptyList(t *testing.T) {
	dir := t.TempDir()
	js, err := jobs.Discover(dir)
	if err != nil {
		t.Fatalf("Discover su cartella vuota: unexpected error %v", err)
	}
	if len(js) != 0 {
		t.Errorf("Discover su cartella vuota: atteso 0 jobs, got %d", len(js))
	}
}

func TestDiscover_LavoriDirNotExistReturnsError(t *testing.T) {
	_, err := jobs.Discover("/nonexistent/lavori/dir")
	if err == nil {
		t.Fatal("Discover su path inesistente: expected error, got nil")
	}
}

func TestDiscover_HappyPathWithMetadata(t *testing.T) {
	lavoriDir := t.TempDir()
	jobDir := filepath.Join(lavoriDir, "CAPABILITY A — Profilo revisione")
	if err := os.Mkdir(jobDir, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	metadata := `profile = "revisione"
trigger_prompt_file = "Prompt_INPUT.rtf"
`
	if err := os.WriteFile(filepath.Join(jobDir, "_labnexus.toml"), []byte(metadata), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	js, err := jobs.Discover(lavoriDir)
	if err != nil {
		t.Fatalf("Discover: unexpected error %v", err)
	}
	if len(js) != 1 {
		t.Fatalf("atteso 1 job, got %d", len(js))
	}
	if js[0].Profile != "revisione" {
		t.Errorf("Profile: got %q, want revisione", js[0].Profile)
	}
	if js[0].TriggerPromptFile != "Prompt_INPUT.rtf" {
		t.Errorf("TriggerPromptFile: got %q, want Prompt_INPUT.rtf", js[0].TriggerPromptFile)
	}
	if js[0].Source != "metadata" {
		t.Errorf("Source: got %q, want metadata", js[0].Source)
	}
}

func TestDiscover_FallbackConventionNaming(t *testing.T) {
	lavoriDir := t.TempDir()
	// Cartella SENZA _labnexus.toml ma con convention "Profilo X" nel nome.
	jobDir := filepath.Join(lavoriDir, "Nuovo Lavoro — Profilo audit-checklist")
	if err := os.Mkdir(jobDir, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	js, err := jobs.Discover(lavoriDir)
	if err != nil {
		t.Fatalf("Discover: unexpected error %v", err)
	}
	if len(js) != 1 {
		t.Fatalf("atteso 1 job via convention, got %d", len(js))
	}
	if js[0].Profile != "audit-checklist" {
		t.Errorf("Profile (convention): got %q, want audit-checklist", js[0].Profile)
	}
	if js[0].Source != "convention" {
		t.Errorf("Source: got %q, want convention", js[0].Source)
	}
}

func TestDiscover_CartellaSenzaMetadataNeConventionSkippata(t *testing.T) {
	lavoriDir := t.TempDir()
	// Cartella senza _labnexus.toml e SENZA "Profilo X" nel nome.
	jobDir := filepath.Join(lavoriDir, "cartella-random")
	if err := os.Mkdir(jobDir, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	js, _ := jobs.Discover(lavoriDir)
	if len(js) != 0 {
		t.Errorf("cartella senza metadata né convention dovrebbe essere skippata, got %d jobs", len(js))
	}
}

func TestDiscover_MetadataMalformatoSkippataConWarning(t *testing.T) {
	lavoriDir := t.TempDir()
	jobDir := filepath.Join(lavoriDir, "malformed-job")
	if err := os.Mkdir(jobDir, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	// _labnexus.toml con sintassi rotta.
	if err := os.WriteFile(filepath.Join(jobDir, "_labnexus.toml"), []byte("profile = \"X\n[broken section"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	js, _ := jobs.Discover(lavoriDir)
	if len(js) != 0 {
		t.Errorf("cartella con metadata malformato dovrebbe essere skippata, got %d jobs", len(js))
	}
}

func TestResolve_ExactMatch(t *testing.T) {
	js := []*jobs.Job{
		{Name: "revisione", Profile: "revisione"},
		{Name: "rilievi", Profile: "rilievi"},
	}
	got := jobs.Resolve(js, "revisione")
	if got == nil {
		t.Fatal("Resolve('revisione'): atteso match, got nil")
	}
	if got.Profile != "revisione" {
		t.Errorf("got Profile %q, want revisione", got.Profile)
	}
}

func TestResolve_CaseInsensitive(t *testing.T) {
	js := []*jobs.Job{{Name: "Revisione DOC", Profile: "revisione"}}
	got := jobs.Resolve(js, "revisione doc")
	if got == nil {
		t.Fatal("Resolve case-insensitive: atteso match, got nil")
	}
}

func TestResolve_SubstringMatch(t *testing.T) {
	js := []*jobs.Job{{Name: "CAPABILITY A — Profilo revisione", Profile: "revisione"}}
	got := jobs.Resolve(js, "revisione")
	if got == nil {
		t.Fatal("Resolve substring: atteso match su 'revisione' dentro nome lungo, got nil")
	}
}

func TestResolve_NoMatchReturnsNil(t *testing.T) {
	js := []*jobs.Job{{Name: "revisione", Profile: "revisione"}}
	got := jobs.Resolve(js, "inesistente")
	if got != nil {
		t.Errorf("Resolve no-match: atteso nil, got %+v", got)
	}
}

func TestLoadMetadata_HappyPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "_labnexus.toml")
	body := `profile = "revisione"
trigger_prompt_file = "Prompt_INPUT.rtf"
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	profile, triggerFile, err := jobs.LoadMetadata(path)
	if err != nil {
		t.Fatalf("LoadMetadata: unexpected error %v", err)
	}
	if profile != "revisione" {
		t.Errorf("profile: got %q, want revisione", profile)
	}
	if triggerFile != "Prompt_INPUT.rtf" {
		t.Errorf("triggerFile: got %q, want Prompt_INPUT.rtf", triggerFile)
	}
}

func TestLoadMetadata_FileMissing(t *testing.T) {
	_, _, err := jobs.LoadMetadata("/nonexistent/_labnexus.toml")
	if err == nil {
		t.Fatal("LoadMetadata su file mancante: atteso error, got nil")
	}
}

func TestLoadMetadata_OnlyProfile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "_labnexus.toml")
	// Solo profile, trigger_prompt_file omesso (opzionale).
	if err := os.WriteFile(path, []byte(`profile = "rilievi"`), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	profile, triggerFile, err := jobs.LoadMetadata(path)
	if err != nil {
		t.Fatalf("LoadMetadata: unexpected error %v", err)
	}
	if profile != "rilievi" {
		t.Errorf("profile: got %q, want rilievi", profile)
	}
	if triggerFile != "" {
		t.Errorf("triggerFile dovrebbe essere vuoto (opzionale), got %q", triggerFile)
	}
}

func TestDiscover_OutputDirDefaultIsInputSlashOutput(t *testing.T) {
	lavoriDir := t.TempDir()
	jobDir := filepath.Join(lavoriDir, "CAPABILITY X — Profilo revisione")
	if err := os.Mkdir(jobDir, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	js, err := jobs.Discover(lavoriDir)
	if err != nil {
		t.Fatalf("Discover: unexpected error %v", err)
	}
	if len(js) != 1 {
		t.Fatalf("atteso 1 job, got %d", len(js))
	}
	expected := filepath.Join(jobDir, "output")
	if !strings.HasSuffix(js[0].OutputDir, "/output") {
		t.Errorf("OutputDir dovrebbe finire in '/output' (default), got %q (expected~ %q)", js[0].OutputDir, expected)
	}
}
