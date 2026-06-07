package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/labnexus/labnexus/internal/config"
	"github.com/labnexus/labnexus/internal/profile"
)

func writeTOML(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "x.toml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

// Il PERNO della catena config→profilo→_labnexus: Overlay sovrascrive SOLO le
// chiavi PRESENTI nel file e lascia invariate le assenti.
func TestOverlay_OverridesPresentKeepsAbsent(t *testing.T) {
	active := &profile.Profile{Provider: "eurouter", Modello: "kimi-k2.6", MaxTokens: 32768, Temperature: 1.0, ContextWindow: 256000}
	if err := config.Overlay(active, writeTOML(t, "max_tokens = 24576\n")); err != nil {
		t.Fatal(err)
	}
	if active.MaxTokens != 24576 {
		t.Errorf("chiave presente deve essere sovrascritta: max_tokens=%d, want 24576", active.MaxTokens)
	}
	if active.Modello != "kimi-k2.6" || active.Temperature != 1.0 || active.ContextWindow != 256000 {
		t.Errorf("chiavi ASSENTI devono restare invariate: %+v", active)
	}
}

// Presence-aware: temperature = 0.0 ESPLICITO viene applicato (risolve la
// limitazione "0 vs omesso" di Merge sui campi numerici).
func TestOverlay_ExplicitZeroApplied(t *testing.T) {
	active := &profile.Profile{Temperature: 1.0}
	if err := config.Overlay(active, writeTOML(t, "temperature = 0.0\n")); err != nil {
		t.Fatal(err)
	}
	if active.Temperature != 0.0 {
		t.Errorf("temperature = 0.0 esplicita deve essere applicata, got %v", active.Temperature)
	}
}

// Overlay NON tocca i trigger (XOR FR-16 + sorgente tracciata): un _labnexus con
// trigger_prompt_file non deve riversarsi nel profilo via Overlay (evita la
// violazione XOR e lascia la gestione al path dedicato del runner).
func TestOverlay_IgnoresTrigger(t *testing.T) {
	active := &profile.Profile{TriggerPrompt: "default del profilo"}
	toml := "trigger_prompt_file = \"Prompt.rtf\"\nmax_tokens = 20000\n"
	if err := config.Overlay(active, writeTOML(t, toml)); err != nil {
		t.Fatal(err)
	}
	if active.TriggerPromptFile != "" {
		t.Errorf("Overlay non deve applicare trigger_prompt_file: got %q", active.TriggerPromptFile)
	}
	if active.TriggerPrompt != "default del profilo" {
		t.Errorf("il trigger del profilo deve restare invariato: got %q", active.TriggerPrompt)
	}
	if active.MaxTokens != 20000 {
		t.Errorf("i parametri non-trigger devono essere applicati: max_tokens=%d", active.MaxTokens)
	}
}

// Catena a 3 strati: master(32768) già risolto nel profilo → _labnexus(24576) vince.
func TestOverlay_ChainJobWins(t *testing.T) {
	active := &profile.Profile{MaxTokens: 32768} // stato dopo config.Merge(master, profilo)
	if err := config.Overlay(active, writeTOML(t, "max_tokens = 24576\n")); err != nil {
		t.Fatal(err)
	}
	if active.MaxTokens != 24576 {
		t.Errorf("lo strato _labnexus deve vincere: got %d, want 24576", active.MaxTokens)
	}
}
