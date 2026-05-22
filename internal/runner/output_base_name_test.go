package runner

import (
	"strings"
	"testing"
	"time"

	"github.com/labnexus/labnexus/internal/input"
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
