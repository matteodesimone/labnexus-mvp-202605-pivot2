package runner

import (
	"strings"
	"testing"
	"time"

	"github.com/labnexus/labnexus/internal/input"
	"github.com/labnexus/labnexus/internal/output"
	"github.com/labnexus/labnexus/internal/profile"
)

// TestOutputBaseName_StripsPathFromInputName: dopo l'abilitazione del walk
// ricorsivo (input package), ParsedFile.Name può essere un path relativo
// (es. "Documento_da_revisionare/MQL_Rev03.docx"). outputBaseName viene usato
// come basename per il `.log` accoppiato (runlog.OpenLogFile fa
// filepath.Join(outputDir, baseName+".log")) — se contiene "/" tenta di
// scrivere in una subdir inesistente e l'audit trail si disabilita.
func TestOutputBaseName_StripsPathFromInputName(t *testing.T) {
	p := &profile.Profile{Profilo: "revisione"}
	parsed := &input.ParsedResult{
		Files: []input.ParsedFile{
			{Name: "Documento_da_revisionare/MQL_Rev03.docx", Text: "x", Size: 1},
		},
	}
	started := time.Date(2026, 5, 22, 8, 56, 5, 0, time.UTC)

	got := outputBaseName(p, parsed, started)

	if strings.Contains(got, "/") {
		t.Fatalf("outputBaseName non deve contenere '/' (rompe runlog.OpenLogFile): got %q", got)
	}
	if !strings.Contains(got, "MQL_Rev03") {
		t.Errorf("outputBaseName dovrebbe contenere il basename 'MQL_Rev03', got %q", got)
	}
	if !strings.Contains(got, "revisione") {
		t.Errorf("outputBaseName dovrebbe contenere il profilo 'revisione', got %q", got)
	}
}

// TestOutputBaseName_NoInputFilesReturnsTimestampAndProfile: caso bordo
// (parsed.Files vuoto) — deve comunque produrre un baseName valido.
func TestOutputBaseName_NoInputFilesReturnsTimestampAndProfile(t *testing.T) {
	p := &profile.Profile{Profilo: "revisione"}
	parsed := &input.ParsedResult{Files: nil}
	started := time.Date(2026, 5, 22, 8, 56, 5, 0, time.UTC)

	got := outputBaseName(p, parsed, started)

	if strings.Contains(got, "/") {
		t.Errorf("outputBaseName non deve contenere '/': got %q", got)
	}
	if !strings.Contains(got, "revisione") {
		t.Errorf("outputBaseName dovrebbe contenere il profilo, got %q", got)
	}
}

// Bug 2026-06-07: il basename del `.log` conteneva spazi/punti/accenti/virgole
// (dal nome file di input) mentre .md/.pdf/.docx erano sanitizzati con output.Slug
// → nomi base divergenti e .log fragile su kDrive/shell/URL. Ora outputBaseName usa
// lo STESSO slug: i quattro file condividono il nome base.
func TestOutputBaseName_SanitizesLikeOutputWrite(t *testing.T) {
	p := &profile.Profile{Profilo: "review-pack"}
	parsed := &input.ParsedResult{
		Files: []input.ParsedFile{
			{Name: "WP1/M113.02_Non Conformità, Azioni Correttive_2026.xlsx", Text: "x", Size: 1},
		},
	}
	started := time.Date(2026, 6, 7, 11, 32, 1, 0, time.UTC)
	got := outputBaseName(p, parsed, started)

	if strings.ContainsAny(got, " ,") {
		t.Errorf("il basename non deve contenere spazi/virgole: %q", got)
	}
	want := "2026-06-07T113201_" + output.Slug("review-pack") + "_" + output.Slug("M113.02_Non Conformità, Azioni Correttive_2026")
	if got != want {
		t.Errorf("basename = %q, want %q (stesso slug di output.Write)", got, want)
	}
}
