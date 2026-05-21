package runner_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/provider"
	"github.com/labnexus/labnexus/internal/runner"
)

// fakeProvider è una LLMProvider mock per test orchestrator.
// Boundary mock (zone 2): non testiamo il vero HTTP, ma il flusso del runner.
type fakeProvider struct {
	name   string
	chunks []string
	done   bool // se true, invia Done:true alla fine; se false, chiude senza done
	err    error // se non-nil, ritorna error subito da Stream
}

func (f *fakeProvider) Name() string { return f.name }

func (f *fakeProvider) Stream(_ context.Context, _, _ string, _ provider.Options) (<-chan provider.StreamEvent, error) {
	if f.err != nil {
		return nil, f.err
	}
	ch := make(chan provider.StreamEvent, len(f.chunks)+1)
	for _, c := range f.chunks {
		ch <- provider.StreamEvent{Token: c}
	}
	if f.done {
		ch <- provider.StreamEvent{Done: true}
	}
	close(ch)
	return ch, nil
}

// setupRunEnv crea profili/, KB-ispettore/, input/, output/ in una tmp dir
// e ritorna una Config pronta. Usa LABNEXUS_OLLAMA_ENDPOINT per puntare a fake-or-broken.
func setupRunEnv(t *testing.T, profileBody string) runner.Config {
	t.Helper()
	tmp := t.TempDir()
	profili := filepath.Join(tmp, "profili")
	kb := filepath.Join(tmp, "KB-ispettore")
	in := filepath.Join(tmp, "input")
	out := filepath.Join(tmp, "out")
	for _, d := range []string{profili, kb, in, out} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(kb, "CLAUDE.md"), []byte("stub kb"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(profili, "test.toml"), []byte(profileBody), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(in, "doc.md"), []byte("input content"), 0o644); err != nil {
		t.Fatal(err)
	}
	return runner.Config{
		ProfileName: "test",
		InputDir:    in,
		OutputDir:   out,
		ProfiliDir:  profili,
		KbDir:       kb,
	}
}

func defaultTestProfile() string {
	return fmt.Sprintf(`profilo = "test"
descrizione = "profilo di test runner"
provider = "ollama"
modello = "qwen3.6"
kb_files = ["CLAUDE.md"]
trigger_prompt = "%s"
`, strings.Repeat("a", 80))
}

func TestRun_RequiresProfileName(t *testing.T) {
	_, err := runner.Run(runner.Config{})
	if err == nil {
		t.Fatal("expected error when ProfileName empty, got nil")
	}
	if !strings.Contains(err.Error(), "profile") {
		t.Errorf("error should mention 'profile', got: %v", err)
	}
}

func TestRun_DryRunSkipsLLMCall(t *testing.T) {
	cfg := setupRunEnv(t, defaultTestProfile())
	cfg.DryRun = true
	// Senza fake Ollama: se il runner provasse a chiamare LLM, fallirebbe.
	res, err := runner.Run(cfg)
	if err != nil {
		t.Fatalf("dry-run should succeed without LLM, got: %v", err)
	}
	if res.ExitCode != 0 {
		t.Errorf("dry-run exit code: want 0, got %d", res.ExitCode)
	}
	if res.TokensUsed <= 0 {
		t.Errorf("dry-run should still compute token estimate, got %d", res.TokensUsed)
	}
}

func TestRun_ProfileNotFoundReturnsError(t *testing.T) {
	cfg := setupRunEnv(t, defaultTestProfile())
	cfg.ProfileName = "inesistente"
	_, err := runner.Run(cfg)
	if err == nil {
		t.Fatal("expected error for non-existent profile, got nil")
	}
	if !strings.Contains(err.Error(), "non trovato") {
		t.Errorf("error should mention 'non trovato', got: %v", err)
	}
}

func TestRun_ContextWindowExceededBlocks(t *testing.T) {
	cfg := setupRunEnv(t, fmt.Sprintf(`profilo = "test"
descrizione = "ctx piccolo"
provider = "ollama"
modello = "qwen3.6"
context_window = 64
kb_files = ["CLAUDE.md"]
trigger_prompt = "%s"
`, strings.Repeat("a", 80)))
	// Riempi input.md con tantissimi caratteri → supera 64 token (256 chars)
	bigInput := filepath.Join(cfg.InputDir, "big.txt")
	_ = os.WriteFile(bigInput, []byte(strings.Repeat("a", 10000)), 0o644)

	_, err := runner.Run(cfg)
	if err == nil {
		t.Fatal("expected error for context window exceeded, got nil")
	}
	if !strings.Contains(err.Error(), "context window") {
		t.Errorf("error should mention 'context window', got: %v", err)
	}
}

func TestRun_InvalidProfileRejectedByValidate(t *testing.T) {
	cfg := setupRunEnv(t, `profilo = "test"
descrizione = "trigger corto"
provider = "ollama"
modello = "qwen3.6"
kb_files = ["CLAUDE.md"]
trigger_prompt = "corto"
`)
	_, err := runner.Run(cfg)
	if err == nil {
		t.Fatal("expected validation error for short trigger_prompt, got nil")
	}
	if !strings.Contains(err.Error(), "trigger_prompt") {
		t.Errorf("error should mention 'trigger_prompt', got: %v", err)
	}
}

// Note: copertura del path "chiamata provider + scrittura output" è negli
// integration test di provider/{ollama,eurouter}_integration_test.go (che
// validano OllamaProvider e EurouterProvider end-to-end con httptest.NewServer).
// L'integrazione completa runner+provider+writer è validata dagli scenari BDD
// che invocano il binario come subprocess (vedi features/).

// TestRun_MergeMasterAndValidatePostMerge (fix review 1.5.B HIGH-mistral-2):
// verifica end-to-end che il config master fornisca i campi opzionali al profile
// parziale, e che profile.Validate passi sul merged invece che sul profile isolato.
func TestRun_MergeMasterAndValidatePostMerge(t *testing.T) {
	cfg := setupRunEnv(t, `profilo = "test"
descrizione = "profilo parziale (provider/modello ereditati dal master)"
kb_files = ["CLAUDE.md"]
trigger_prompt = "Trigger prompt sufficiente a superare la soglia minima di 50 caratteri richiesta dallo schema."
`)
	// Crea master config nella tmp dir (relativo cfg.ProfiliDir + ../)
	masterPath := filepath.Join(filepath.Dir(cfg.ProfiliDir), "labnexus.config.toml")
	masterBody := `provider = "ollama"
modello = "qwen3.6"
temperature = 0.9
max_tokens = 8192
context_window = 128000
`
	if err := os.WriteFile(masterPath, []byte(masterBody), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg.ConfigPath = masterPath
	cfg.DryRun = true

	res, err := runner.Run(cfg)
	if err != nil {
		t.Fatalf("merge+validate dovrebbe passare con master che fornisce provider/modello: %v", err)
	}
	if res.ExitCode != 0 {
		t.Errorf("exit code: want 0, got %d", res.ExitCode)
	}
}

// TestRun_MergeMasterInvalidProviderFails (fix review 1.5.B HIGH-mistral-2 dual):
// verifica che se il master fornisce un provider invalido e il profile non lo
// override, il merged profile è respinto da Validate.
func TestRun_MergeMasterInvalidProviderFails(t *testing.T) {
	cfg := setupRunEnv(t, `profilo = "test"
descrizione = "profilo parziale"
kb_files = ["CLAUDE.md"]
trigger_prompt = "Trigger prompt sufficiente a superare la soglia minima di 50 caratteri richiesta dallo schema."
`)
	masterPath := filepath.Join(filepath.Dir(cfg.ProfiliDir), "labnexus.config.toml")
	if err := os.WriteFile(masterPath, []byte(`provider = "openai"
modello = "gpt-4"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg.ConfigPath = masterPath
	cfg.DryRun = true

	_, err := runner.Run(cfg)
	if err == nil {
		t.Fatal("validate post-merge dovrebbe fallire per provider invalido nel master, got nil")
	}
	if !strings.Contains(err.Error(), "provider") {
		t.Errorf("error dovrebbe menzionare 'provider', got: %v", err)
	}
}
