package runner

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/output"
	"github.com/labnexus/labnexus/internal/runlog"
)

// TestMaybeWritePDF_DisabledNoOp: con pdfEnabled=false non scrive nulla.
func TestMaybeWritePDF_DisabledNoOp(t *testing.T) {
	dir := t.TempDir()
	mdPath := filepath.Join(dir, "out_test.md")
	if err := os.WriteFile(mdPath, []byte("# Test\n"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	fm := &output.Frontmatter{Profilo: "revisione"}
	var buf bytes.Buffer
	log := runlog.New(&buf)

	got := maybeWritePDF(mdPath, "# Test\n", fm, false, "Test", log)
	if got != "" {
		t.Errorf("pdf disabled deve ritornare empty path, got %q", got)
	}
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".pdf") {
			t.Errorf("nessun .pdf doveva essere creato, trovato %q", e.Name())
		}
	}
}

// TestMaybeWritePDF_EnabledCreatesPDFAccoppiato: con pdfEnabled=true scrive
// un file .pdf con lo stesso basename del .md.
// Skippato quando pandoc/typst non sono in dist/.../bin/ (es. dopo `make clean`).
func TestMaybeWritePDF_EnabledCreatesPDFAccoppiato(t *testing.T) {
	dir := t.TempDir()
	mdPath := filepath.Join(dir, "2026-05-22T085605_revisione_MQL_Rev03.md")
	body := "# Bozza revisione MQL_Rev03\n\nCorpo del documento.\n"
	if err := os.WriteFile(mdPath, []byte(body), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	fm := &output.Frontmatter{
		Profilo:        "revisione",
		Modello:        "qwen3.5-122b-a10b",
		Provider:       "eurouter",
		DataEsecuzione: "2026-05-22T08:56:05+02:00",
		FileInput:      []string{"Documento_da_revisionare/MQL_Rev03.docx"},
		Stato:          "completato",
	}
	var buf bytes.Buffer
	log := runlog.New(&buf)

	got := maybeWritePDF(mdPath, body, fm, true, "Revisione SGQ", log)

	if got == "" {
		// pdf.Render fallisce con warning non-bloccante se pandoc/typst mancano.
		// Detection: log contiene "PDF disabilitato" + il messaggio menziona "pandoc"/"typst".
		logStr := buf.String()
		if strings.Contains(logStr, "pandoc") || strings.Contains(logStr, "typst") {
			t.Skipf("pandoc/typst non disponibili in dist/.../bin/ (esegui `make pdf-tools`): %s", logStr)
		}
		t.Fatalf("path PDF vuoto, log: %s", logStr)
	}
	wantPath := filepath.Join(dir, "2026-05-22T085605_revisione_MQL_Rev03.pdf")
	if got != wantPath {
		t.Errorf("path PDF: got %q, want %q", got, wantPath)
	}
	data, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf(".pdf accoppiato non trovato: %v", err)
	}
	if len(data) < 4 || string(data[:4]) != "%PDF" {
		t.Errorf(".pdf magic bytes mancanti, primi 4 bytes: %q", data[:4])
	}
}

// TestMaybeWritePDF_RenderErrorNonBlocking: se il render fallisce (es. body
// inaspettatamente non parsabile, in pratica raro), warning a log e ritorna
// "" senza crash. Pattern audit-friendly: il .md è già scritto, il run
// non si rompe.
//
// Simuliamo questo passando un mdPath in una directory NON scrivibile
// (read-only), così il pdf.Render non può creare il file di output.
func TestMaybeWritePDF_WriteFailureWarningNonBlocking(t *testing.T) {
	dir := t.TempDir()
	// mdPath in una dir esistente
	mdPath := filepath.Join(dir, "out.md")
	if err := os.WriteFile(mdPath, []byte("# x\n"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	// Rendi dir read-only così la scrittura del .pdf fallisce
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Skipf("chmod ro non supportato: %v", err)
	}
	defer os.Chmod(dir, 0o755)

	fm := &output.Frontmatter{Profilo: "revisione"}
	var buf bytes.Buffer
	log := runlog.New(&buf)

	got := maybeWritePDF(mdPath, "# x\n", fm, true, "Test", log)
	if got != "" {
		t.Errorf("su errore di scrittura, path PDF deve essere empty (non-blocking), got %q", got)
	}
	if !strings.Contains(buf.String(), "PDF") {
		t.Errorf("log deve contenere warning sul PDF fallito, log: %q", buf.String())
	}
}
