package input

import (
	"strings"
	"testing"
)

// Sprint 1.5.B test-scaffold (red phase): tutti i test devono fallire finché
// /v-implement non implementa parseXls.

func TestParseXls_HappyPath(t *testing.T) {
	// Fixture: usa uno dei file reali di Denis come oracolo.
	// Al /v-implement, verifica che il file esista; se no, skip.
	fixture := "../../lavori/CAPABILITY G — Profilo pt-analysis/M132.00_Rapporto Proficiency Testing.xls"
	text, err := parseXls(fixture)
	if err != nil {
		t.Fatalf("parseXls: unexpected error %v", err)
	}
	if !strings.Contains(text, "--- SHEET:") {
		t.Errorf("parseXls: testo estratto non contiene separator '--- SHEET:'")
	}
	if len(text) < 100 {
		t.Errorf("parseXls: testo estratto troppo corto (%d caratteri), atteso almeno 100", len(text))
	}
}

func TestParseXls_FileNonEsistente(t *testing.T) {
	_, err := parseXls("/nonexistent/file.xls")
	if err == nil {
		t.Fatal("parseXls: expected error per file mancante, got nil")
	}
}
