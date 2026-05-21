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

// --- Bug frontmatter-default-non-applicato (cross-fetta, P2 todo #006) ---
// I profili dichiarano output.frontmatter_default con chiavi/valori che
// devono comparire nel frontmatter dell'output (es. tipo:
// management_review_pack, stato: bozza_da_validare_qm, profilo_labnexus,
// locale). L'engine deve merge-arli al frontmatter generato, con precedenza
// alle chiavi engine-generated in caso di collisione.

func TestWrite_AppliesProfileDefaults(t *testing.T) {
	dir := t.TempDir()
	fm := &output.Frontmatter{
		Profilo:        "review-pack",
		Modello:        "qwen3.6",
		Provider:       "ollama",
		DataEsecuzione: "2026-05-20T10:00:00+02:00",
		ProfileDefaults: map[string]string{
			"tipo":             "management_review_pack",
			"stato":            "bozza_da_validare_qm",
			"profilo_labnexus": "review-pack",
			"locale":           "true",
		},
	}
	path, err := output.Write(dir, fm, "body")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, _ := os.ReadFile(path)
	s := string(content)
	for _, want := range []string{
		"tipo: management_review_pack",
		"profilo_labnexus: review-pack",
		"locale: true",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("default mancante %q in:\n%s", want, s)
		}
	}
}

func TestWrite_EngineKeysWinOverProfileDefaults(t *testing.T) {
	// Se il profilo dichiarasse "modello: qwen-default" come default, la
	// chiave engine-generated `modello: qwen3.6` deve vincere (l'engine sa
	// che modello ha effettivamente chiamato).
	dir := t.TempDir()
	fm := &output.Frontmatter{
		Profilo:        "review-pack",
		Modello:        "qwen3.6",
		Provider:       "ollama",
		DataEsecuzione: "2026-05-20T10:00:00+02:00",
		ProfileDefaults: map[string]string{
			"modello":          "qwen-default-from-profile", // collision con engine
			"tipo":             "management_review_pack",
			"profilo_labnexus": "review-pack",
		},
	}
	path, err := output.Write(dir, fm, "body")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, _ := os.ReadFile(path)
	s := string(content)
	if !strings.Contains(s, "modello: qwen3.6") {
		t.Errorf("engine modello deve vincere, got:\n%s", s)
	}
	if strings.Contains(s, "qwen-default-from-profile") {
		t.Errorf("profile default modello NON deve sovrascrivere engine, got:\n%s", s)
	}
	// Le altre chiavi non-collidenti restano.
	if !strings.Contains(s, "tipo: management_review_pack") {
		t.Errorf("tipo default deve essere presente, got:\n%s", s)
	}
}

// TestWrite_ValutazioneDenisExternalOwnedPreserved (fix review 1.5.B mistral MEDIUM):
// FR-24 tabella ownership chiavi frontmatter. `valutazione_denis` è EXTERNAL-OWNED
// (popolato da Denis post-run). Se il profile-default lo dichiara (es. Denis lo
// pre-popola in un profilo di workflow specifico), l'engine NON deve sovrascriverlo
// con stringa vuota (ValutazioneDenis ha yaml:"...,omitempty" — quando non setted
// dall'engine, viene omesso dal frontmatter renderizzato, e profile-default vince).
func TestWrite_ValutazioneDenisExternalOwnedPreserved(t *testing.T) {
	dir := t.TempDir()
	fm := &output.Frontmatter{
		Profilo:        "revisione",
		Modello:        "qwen3.6",
		Provider:       "eurouter",
		DataEsecuzione: "2026-05-21T10:00:00+02:00",
		// ValutazioneDenis NON setted dall'engine (caso normale: nessun runner che
		// pre-popola la valutazione, Denis la mette manualmente post-run).
		ProfileDefaults: map[string]string{
			"valutazione_denis": "validata_con_riserva", // pre-popolato dal profilo
		},
	}
	path, err := output.Write(dir, fm, "body")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, _ := os.ReadFile(path)
	s := string(content)
	if !strings.Contains(s, "valutazione_denis: validata_con_riserva") {
		t.Errorf("profile-default valutazione_denis deve vincere quando engine omette il campo, got:\n%s", s)
	}
}

func TestWrite_NoProfileDefaultsLeavesFrontmatterUnchanged(t *testing.T) {
	// Backwards-compat: senza ProfileDefaults, comportamento esistente.
	dir := t.TempDir()
	fm := &output.Frontmatter{
		Profilo:        "revisione",
		Modello:        "qwen3.6",
		Provider:       "ollama",
		DataEsecuzione: "2026-05-20T10:00:00+02:00",
	}
	path, err := output.Write(dir, fm, "body")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, _ := os.ReadFile(path)
	s := string(content)
	if !strings.Contains(s, "profilo: revisione") {
		t.Errorf("frontmatter standard intatto, got:\n%s", s)
	}
}
