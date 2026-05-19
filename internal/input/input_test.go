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

func TestParseDir_NonRecursiveSkipsSubdirs(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "top.md", "top level")
	sub := filepath.Join(dir, "sub")
	_ = os.MkdirAll(sub, 0o755)
	writeFile(t, sub, "deep.md", "should be skipped")

	res, err := input.ParseDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Files) != 1 {
		t.Fatalf("expected 1 file (walk non-ricorsivo), got %d", len(res.Files))
	}
	if res.Files[0].Name != "top.md" {
		t.Errorf("expected top.md, got %q", res.Files[0].Name)
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
