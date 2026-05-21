package input

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Sprint 1.5.B test-scaffold (red phase): test devono fallire finché
// /v-implement non implementa parseRtf.

func TestParseRtf_HappyPath(t *testing.T) {
	// Fixture sintetica minima: RTF valido con plain text "Ciao mondo".
	dir := t.TempDir()
	path := filepath.Join(dir, "test.rtf")
	body := `{\rtf1\ansi\ansicpg1252\cocoartf2580
{\fonttbl\f0\fswiss\fcharset0 Helvetica;}
\f0\fs24 Ciao mondo}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	text, err := parseRtf(path)
	if err != nil {
		t.Fatalf("parseRtf: unexpected error %v", err)
	}
	if !strings.Contains(text, "Ciao mondo") {
		t.Errorf("parseRtf: testo estratto %q non contiene 'Ciao mondo'", text)
	}
	if strings.Contains(text, `\rtf1`) {
		t.Errorf("parseRtf: testo estratto contiene control word '\\rtf1' (atteso strip)")
	}
	if strings.Contains(text, `\fonttbl`) {
		t.Errorf("parseRtf: testo estratto contiene control word '\\fonttbl' (atteso strip)")
	}
}

func TestParseRtf_FileDenis(t *testing.T) {
	// Fixture: file reale di Denis se disponibile.
	fixture := "../../lavori/CAPABILITY A — Profilo revisione/Prompt_INPUT_Rev.00.rtf"
	if _, err := os.Stat(fixture); err != nil {
		t.Skipf("fixture %q non disponibile: %v", fixture, err)
	}
	text, err := parseRtf(fixture)
	if err != nil {
		t.Fatalf("parseRtf: unexpected error %v", err)
	}
	if len(text) < 20 {
		t.Errorf("parseRtf: testo estratto troppo corto (%d caratteri), atteso almeno 20", len(text))
	}
}
