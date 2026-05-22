package input_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/input"
)

// writeFile è un helper test per scrivere un file con contenuto dato.
func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("writeFile %q: %v", name, err)
	}
}

func TestParseDir_PlainFormatsExtractText(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "doc.md", "# Markdown\nhello md")
	writeFile(t, dir, "doc.txt", "plain text content")
	writeFile(t, dir, "doc.csv", "a,b,c\n1,2,3")

	res, err := input.ParseDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Files) != 3 {
		t.Fatalf("expected 3 parsed files, got %d", len(res.Files))
	}
	if len(res.Skipped) != 0 {
		t.Errorf("expected no skipped, got %d (%v)", len(res.Skipped), res.Skipped)
	}
}

func TestParseDir_EmptyDirReturnsError(t *testing.T) {
	dir := t.TempDir()
	_, err := input.ParseDir(dir)
	if err == nil {
		t.Fatal("expected error for empty dir, got nil")
	}
	if !strings.Contains(err.Error(), "vuota") {
		t.Errorf("error should mention 'vuota', got: %v", err)
	}
}

func TestParseDir_OnlyUnsupportedReturnsError(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "image.jpg", "fake jpg")
	writeFile(t, dir, "video.mov", "fake mov")
	_, err := input.ParseDir(dir)
	if err == nil {
		t.Fatal("expected error when only unsupported files present, got nil")
	}
	if !strings.Contains(err.Error(), ".md") {
		t.Errorf("error should list supported formats, got: %v", err)
	}
}

// TestParseDir_RecursiveIncludesSubdirs (post-Sprint 1.5 fix): il walk è
// ricorsivo. Le sottocartelle vengono attraversate e i file inclusi nel risultato
// con Name = path relativo alla input dir (POSIX), così il prompt preserva il
// significato semantico della struttura cartelle scelta dall'utente.
func TestParseDir_RecursiveIncludesSubdirs(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "top.md", "top level")
	sub := filepath.Join(dir, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	writeFile(t, sub, "deep.md", "deep content")

	res, err := input.ParseDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Files) != 2 {
		t.Fatalf("expected 2 files (walk ricorsivo), got %d: %+v", len(res.Files), res.Files)
	}
	names := []string{res.Files[0].Name, res.Files[1].Name}
	expected := []string{"sub/deep.md", "top.md"}
	for i, want := range expected {
		if names[i] != want {
			t.Errorf("Files[%d].Name = %q, want %q (ordine alfabetico, path POSIX)", i, names[i], want)
		}
	}
}

// TestParseDir_NestedSubdirs verifica depth > 1 e ordine alfabetico path-based.
func TestParseDir_NestedSubdirs(t *testing.T) {
	dir := t.TempDir()
	nested := filepath.Join(dir, "a", "b", "c")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	writeFile(t, nested, "deep.md", "very deep")
	writeFile(t, dir, "root.md", "root")

	res, err := input.ParseDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Files) != 2 {
		t.Fatalf("expected 2 files, got %d: %+v", len(res.Files), res.Files)
	}
	if res.Files[0].Name != "a/b/c/deep.md" {
		t.Errorf("Files[0].Name = %q, want %q", res.Files[0].Name, "a/b/c/deep.md")
	}
	if res.Files[1].Name != "root.md" {
		t.Errorf("Files[1].Name = %q, want %q", res.Files[1].Name, "root.md")
	}
}

// TestParseDir_OutputSubdirSkipped: la cartella `output/` (dove il tool stesso
// scrive i risultati) viene esclusa dal walk per evitare loop di feedback
// (il risultato di un run precedente non deve entrare come input del prossimo).
func TestParseDir_OutputSubdirSkipped(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "input.md", "real input")
	outDir := filepath.Join(dir, "output")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("mkdir output: %v", err)
	}
	writeFile(t, outDir, "previous-run.md", "questo NON deve essere letto")

	res, err := input.ParseDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Files) != 1 {
		t.Fatalf("expected 1 file (output/ skipped), got %d: %+v", len(res.Files), res.Files)
	}
	if res.Files[0].Name != "input.md" {
		t.Errorf("expected input.md, got %q", res.Files[0].Name)
	}
}

// TestParseDir_DotDirsSkipped: cartelle con prefisso `.` (convenzione unix per
// directory di sistema/nascoste, es. `.git`) vengono saltate.
func TestParseDir_DotDirsSkipped(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "real.md", "real content")
	gitDir := filepath.Join(dir, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	writeFile(t, gitDir, "config.md", "git internal — NON leggere")

	res, err := input.ParseDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Files) != 1 {
		t.Fatalf("expected 1 file (.git/ skipped), got %d: %+v", len(res.Files), res.Files)
	}
	if res.Files[0].Name != "real.md" {
		t.Errorf("expected real.md, got %q", res.Files[0].Name)
	}
}

func TestParseDir_UnsupportedExtensionSkippedWithWarning(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "doc.md", "valid")
	writeFile(t, dir, "image.jpg", "unsupported")

	res, err := input.ParseDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Files) != 1 {
		t.Errorf("expected 1 parsed, got %d", len(res.Files))
	}
	if len(res.Skipped) != 1 || res.Skipped[0].Name != "image.jpg" {
		t.Errorf("expected image.jpg skipped, got %v", res.Skipped)
	}
}

func TestParseDir_PDFGarbageSkippedHalfBoundaryProceeds(t *testing.T) {
	// 1 PDF rotto su 2 supportati → 50% = NON trigger errore (soglia strict-majority > 50%).
	dir := t.TempDir()
	writeFile(t, dir, "good.md", "markdown valid")
	writeFile(t, dir, "broken.pdf", "non-pdf-garbage")
	res, err := input.ParseDir(dir)
	if err != nil {
		t.Fatalf("expected no error at exactly 50%% fail rate, got: %v", err)
	}
	if len(res.Files) != 1 {
		t.Errorf("expected 1 file parsed (the .md), got %d", len(res.Files))
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped (the .pdf), got %d", len(res.Skipped))
	}
}

func TestParseDir_FailureRatioAbove50PercentReturnsError(t *testing.T) {
	// 3 PDF rotti + 1 valido = 75% → error.
	dir := t.TempDir()
	writeFile(t, dir, "good.md", "valid")
	writeFile(t, dir, "broken1.pdf", "garbage1")
	writeFile(t, dir, "broken2.pdf", "garbage2")
	writeFile(t, dir, "broken3.pdf", "garbage3")
	_, err := input.ParseDir(dir)
	if err == nil {
		t.Fatal("expected error when >50%% of supported files fail, got nil")
	}
	if !strings.Contains(err.Error(), "oltre il 50") {
		t.Errorf("error should mention 'oltre il 50', got: %v", err)
	}
}

func TestParseDir_PDFMagicBytesCheckRejectsNonPDF(t *testing.T) {
	dir := t.TempDir()
	// Solo PDF rotto: 1/1 = 100% fail → error (e ratio check si attiva)
	writeFile(t, dir, "broken.pdf", "questo non è un PDF")
	_, err := input.ParseDir(dir)
	if err == nil {
		t.Fatal("expected error for single non-PDF file, got nil")
	}
}

// TestParseTriggerPromptFile_HappyPathTxt (fix review 1.5.B mistral MEDIUM):
// FR-17 happy path con file .txt.
func TestParseTriggerPromptFile_HappyPathTxt(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "trigger.txt", "Trigger plain text di test")
	text, err := input.ParseTriggerPromptFile(dir, "trigger.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(text, "Trigger plain text di test") {
		t.Errorf("text estratto inatteso: %q", text)
	}
}

// TestParseTriggerPromptFile_PathTraversalRejected (fix review 1.5.B mistral
// MEDIUM): NFR-6 amendment esteso a trigger_prompt_file. Path con `..`
// che esce da inputDir deve essere rifiutato (lexical check pre-symlink).
func TestParseTriggerPromptFile_PathTraversalRejected(t *testing.T) {
	dir := t.TempDir()
	// File legittimo dentro inputDir (perché l'errore atteso è path-not-allowed,
	// non file-not-found).
	writeFile(t, dir, "ok.txt", "ok")
	_, err := input.ParseTriggerPromptFile(dir, "../../etc/passwd")
	if err == nil {
		t.Fatal("expected error per path traversal, got nil")
	}
	if !strings.Contains(err.Error(), "non consentito") {
		t.Errorf("error dovrebbe menzionare 'non consentito', got: %v", err)
	}
}

// TestParseTriggerPromptFile_SymlinkEscapeRejected (fix review 1.5.B mistral
// MEDIUM): symlink dentro inputDir che punta fuori → EvalSymlinks rileva e
// rifiuta. Privacy violation prevention.
func TestParseTriggerPromptFile_SymlinkEscapeRejected(t *testing.T) {
	dir := t.TempDir()
	// File sensibile fuori da inputDir.
	parentDir := filepath.Dir(dir)
	sensitive := filepath.Join(parentDir, "sensitive-leak.txt")
	if err := os.WriteFile(sensitive, []byte("data riservata"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	defer os.Remove(sensitive)
	// Symlink dentro inputDir che punta al file fuori.
	linkPath := filepath.Join(dir, "trigger.txt")
	if err := os.Symlink(sensitive, linkPath); err != nil {
		t.Skipf("symlink creation not supported on this filesystem: %v", err)
	}
	_, err := input.ParseTriggerPromptFile(dir, "trigger.txt")
	if err == nil {
		t.Fatal("expected error per symlink escape, got nil (privacy violation)")
	}
	if !strings.Contains(err.Error(), "non consentito") && !strings.Contains(err.Error(), "esce") {
		t.Errorf("error dovrebbe menzionare 'non consentito' o 'esce', got: %v", err)
	}
}
