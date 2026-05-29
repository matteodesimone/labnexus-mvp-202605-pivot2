package features

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/config"
	"github.com/labnexus/labnexus/internal/profile"
)

// requiredProfiles è la lista canonica dei 7 profili (capability A–G) che DEVONO
// essere shippati. Iterare questa lista — anziché solo `profile.List` — garantisce
// che un profilo cancellato o con TOML non parsabile faccia fallire il test
// (profile.List salta silenziosamente i file non parsabili). Codex review HIGH.
var requiredProfiles = []string{
	"revisione",       // A
	"rilievi",         // B
	"review-pack",     // C
	"audit-checklist", // D
	"equipment-alert", // E
	"competence-gap",  // F
	"pt-analysis",     // G
}

// canonicalKBDir è la fonte di verità del KB-ispettore nel repo (v2.1+):
// docs/piano_iniziale/materiali-dominio/KB-ispettore.
func canonicalKBDir() string {
	return filepath.Join(repoRoot, "docs", "piano_iniziale", "materiali-dominio", "KB-ispettore")
}

// rootKBDir è il path che il runner (default --kb-dir) e build-zip.sh (`cp -RL`)
// usano DAVVERO: il symlink ./KB-ispettore al root del repo.
func rootKBDir() string {
	return filepath.Join(repoRoot, "KB-ispettore")
}

// loadRequiredProfilesMerged carica i 7 profili canonici direttamente per nome,
// fallendo su read/parse error, e li ritorna già merge-ati col master config.
func loadRequiredProfilesMerged(t *testing.T) []*profile.Profile {
	t.Helper()
	master, err := config.Load(filepath.Join(repoRoot, "labnexus.config.toml"))
	if err != nil {
		t.Fatalf("load master config: %v", err)
	}
	out := make([]*profile.Profile, 0, len(requiredProfiles))
	for _, name := range requiredProfiles {
		path := filepath.Join(repoRoot, "profili", name+".toml")
		p, err := profile.Load(path)
		if err != nil {
			t.Fatalf("profilo richiesto %q non caricabile (%s): %v", name, path, err)
		}
		out = append(out, config.Merge(master, p))
	}
	return out
}

// TestShippedProfiles_ValidateAgainstCanonicalKB verifica che TUTTI i 7 profili
// canonici validino contro il KB-ispettore canonico del repo.
//
// Bug 2026-05-29 (refactor KB v01 → v2.1): la KB è stata riorganizzata da
// struttura flat (CLAUDE.md, sezioni-ISO/sezione-N-domande.md, NC-patterns/)
// a cartelle numerate (01_SYSTEM/, 02_LOGICA_ISPETTIVA/, 03_REQUISITI/sez_*,
// 04_APPENDICI/). I kb_files dei 7 profili puntavano ancora ai path v01 →
// profile.Validate fa hard-fail e `labnexus run` rifiuta di partire per tutte
// e 7 le capability A–G.
func TestShippedProfiles_ValidateAgainstCanonicalKB(t *testing.T) {
	kbDir := canonicalKBDir()
	if _, err := os.Stat(kbDir); err != nil {
		t.Fatalf("KB canonico non trovato in %s: %v", kbDir, err)
	}
	for _, p := range loadRequiredProfilesMerged(t) {
		if err := profile.Validate(p, kbDir); err != nil {
			t.Errorf("profilo %q non valida contro il KB canonico v2.1: %v", p.Profilo, err)
		}
	}
}

// TestShippedProfiles_ValidateAgainstRootKB valida i 7 profili contro il path
// REALMENTE usato a runtime: ./KB-ispettore (il symlink al root). Codex review
// HIGH: un root KB stale/parziale (con 00_INDICE.md ma privo di file richiesti
// dai profili) passerebbe i test sul KB canonico ma farebbe fallire il CLI
// all'avvio. Questo test chiude quel gap.
func TestShippedProfiles_ValidateAgainstRootKB(t *testing.T) {
	kbDir := rootKBDir()
	if _, err := os.Stat(kbDir); err != nil {
		t.Fatalf("./KB-ispettore (path runtime) non risolvibile: %v", err)
	}
	for _, p := range loadRequiredProfilesMerged(t) {
		if err := profile.Validate(p, kbDir); err != nil {
			t.Errorf("profilo %q non valida contro ./KB-ispettore (path runtime): %v", p.Profilo, err)
		}
	}
}

// TestRootKBSymlink_ResolvesToCanonicalKB verifica che ./KB-ispettore — su cui
// si basano sia il default --kb-dir del runner sia build-zip.sh (`cp -RL`) —
// sia un symlink risolvibile a una directory che contiene il KB v2.1 E che
// risolva esattamente alla KB canonica (no divergenza stale).
//
// Bug 2026-05-29: ./KB-ispettore era diventato un MacOS Alias file (binario
// Finder da ~1KB, type-change `T` in git) invece del symlink originale.
func TestRootKBSymlink_ResolvesToCanonicalKB(t *testing.T) {
	root := rootKBDir()
	info, err := os.Stat(root) // segue il symlink
	if err != nil {
		t.Fatalf("./KB-ispettore non risolvibile (symlink rotto o alias Finder): %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("./KB-ispettore non risolve a una directory (mode %s) — atteso symlink → KB canonico v2.1", info.Mode().Type())
	}
	if _, err := os.Stat(filepath.Join(root, "00_INDICE.md")); err != nil {
		t.Fatalf("./KB-ispettore non contiene 00_INDICE.md (non punta al KB v2.1): %v", err)
	}
	// Identità forte: il path runtime e quello canonico devono risolvere allo stesso target.
	rootReal, err1 := filepath.EvalSymlinks(root)
	canonReal, err2 := filepath.EvalSymlinks(canonicalKBDir())
	if err1 != nil || err2 != nil {
		t.Fatalf("EvalSymlinks fallita (root=%v, canonical=%v)", err1, err2)
	}
	if rootReal != canonReal {
		t.Fatalf("./KB-ispettore risolve a %q ma il KB canonico è %q — il path runtime diverge", rootReal, canonReal)
	}
}

// TestMetaPromptGenera_NoStaleV01KBPaths guarda contro la regressione del bug:
// il meta-prompt usato per generare nuovi profili (FR-22) NON deve più contenere
// i path della struttura KB v01, altrimenti ogni profilo generato nascerebbe
// rotto contro il KB v2.1. Claude review HIGH.
func TestMetaPromptGenera_NoStaleV01KBPaths(t *testing.T) {
	path := filepath.Join(repoRoot, "docs", "meta-prompt-genera-profilo.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("meta-prompt non leggibile (%s): %v", path, err)
	}
	content := string(data)
	staleTokens := []string{
		"sezioni-ISO/",
		"NC-patterns/",
		"MAPPA_ESPLOSA_REQUISITI_ISO17025",
		"come-pensa-un-ispettore.md",
	}
	for _, tok := range staleTokens {
		if strings.Contains(content, tok) {
			t.Errorf("meta-prompt-genera-profilo.md contiene ancora il path KB v01 %q — i profili generati nascerebbero rotti contro il KB v2.1", tok)
		}
	}
}
