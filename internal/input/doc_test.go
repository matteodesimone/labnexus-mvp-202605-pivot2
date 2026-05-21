package input

import (
	"os"
	"path/filepath"
	"testing"
)

// Sprint 1.5.B test-scaffold (red phase): test devono fallire finché
// /v-implement non implementa parseDoc.

func TestParseDoc_HappyPath(t *testing.T) {
	// Fixture: file reale di Denis se disponibile come oracolo.
	fixture := "../../lavori/CAPABILITY F — Profilo competence-gap/M115.00_Scheda personale_Alessandro Bidossi.doc"
	if _, err := os.Stat(fixture); err != nil {
		t.Skipf("fixture %q non disponibile (skip atteso fino al /v-implement che valuterà se usare uno test sintetico al posto): %v", fixture, err)
	}
	text, err := parseDoc(fixture)
	if err != nil {
		t.Fatalf("parseDoc: unexpected error %v", err)
	}
	if len(text) < 50 {
		t.Errorf("parseDoc: testo estratto troppo corto (%d caratteri), atteso almeno 50", len(text))
	}
}

func TestParseDoc_FileCorrotto(t *testing.T) {
	// File con header non OLE: deve ritornare error (il chiamante warning+skip).
	dir := t.TempDir()
	path := filepath.Join(dir, "corrupt.doc")
	if err := os.WriteFile(path, []byte("not a real OLE file"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	_, err := parseDoc(path)
	if err == nil {
		t.Fatal("parseDoc: expected error per file corrotto, got nil")
	}
}
