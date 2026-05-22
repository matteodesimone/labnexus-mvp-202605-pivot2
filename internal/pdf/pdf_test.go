package pdf_test

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	pdfreader "github.com/ledongthuc/pdf"
	"github.com/labnexus/labnexus/internal/pdf"
)

// renderOrSkip esegue pdf.Render e skippa il test se pandoc/typst non sono
// disponibili (es. CI senza dist bin/). Sprint 1.5.D: i tool sono scaricati
// localmente da scripts/download-pdf-tools.sh in dist/.../bin/.
func renderOrSkip(t *testing.T, md []byte, meta pdf.Metadata) []byte {
	t.Helper()
	var buf bytes.Buffer
	err := pdf.Render(md, meta, pdf.DefaultBrand(), &buf)
	if err != nil {
		if err == pdf.ErrToolsMissing || strings.Contains(err.Error(), "pandoc") || strings.Contains(err.Error(), "typst") {
			t.Skipf("pandoc/typst non disponibili: %v (esegui make pdf-tools)", err)
		}
		t.Fatalf("Render: %v", err)
	}
	return buf.Bytes()
}

func readPDFText(t *testing.T, data []byte) string {
	t.Helper()
	tmp := filepath.Join(t.TempDir(), "out.pdf")
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		t.Fatalf("scrittura tempfile: %v", err)
	}
	f, r, err := pdfreader.Open(tmp)
	if err != nil {
		t.Fatalf("pdf.Open: %v", err)
	}
	defer f.Close()
	var sb strings.Builder
	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, _ := page.GetPlainText(nil)
		sb.WriteString(text)
	}
	return sb.String()
}

func sampleMeta() pdf.Metadata {
	return pdf.Metadata{
		Profilo:        "revisione",
		Modello:        "qwen3.5-122b-a10b",
		Provider:       "eurouter",
		DataEsecuzione: "2026-05-22T08:56:05+02:00",
		DurataSecondi:  65.0,
		TokenStimati:   103840,
		FileInput:      []string{"Documento_da_revisionare/MQL_Rev03.docx"},
		Stato:          "completato",
		Capability:     "Revisione SGQ",
	}
}

func TestRender_ProducesValidPDF(t *testing.T) {
	out := renderOrSkip(t, []byte("# Titolo\n\nCorpo del documento.\n"), sampleMeta())
	if len(out) < 4 || string(out[:4]) != "%PDF" {
		t.Fatalf("output non è un PDF valido, primi bytes: %q", out[:min(8, len(out))])
	}
}

func TestRender_HeaderContainsBrand(t *testing.T) {
	out := renderOrSkip(t, []byte("# Test\n\nBody.\n"), sampleMeta())
	text := readPDFText(t, out)
	if !strings.Contains(text, "LabNexus") {
		t.Errorf("header dovrebbe contenere 'LabNexus', text: %q", text)
	}
}

func TestRender_MetadataBoxContainsFrontmatter(t *testing.T) {
	out := renderOrSkip(t, []byte("# Test\n\nBody.\n"), sampleMeta())
	text := readPDFText(t, out)
	for _, want := range []string{"revisione", "qwen3.5-122b-a10b", "eurouter", "Revisione SGQ"} {
		if !strings.Contains(text, want) {
			t.Errorf("metadata %q assente nel PDF, text: %q", want, text)
		}
	}
}

func TestRender_BodyHeadingPresent(t *testing.T) {
	md := []byte("# Bozza revisione MQL_Rev03\n\nDettagli operativi.\n")
	out := renderOrSkip(t, md, sampleMeta())
	text := readPDFText(t, out)
	if !strings.Contains(text, "Bozza revisione") {
		t.Errorf("body heading mancante, text: %q", text)
	}
}

func TestRender_CalloutObsidianRendered(t *testing.T) {
	md := []byte("# Test\n\n> [!INFO] Versione agente 2.0\n> Dettagli sul versionamento.\n")
	out := renderOrSkip(t, md, sampleMeta())
	text := readPDFText(t, out)
	if !strings.Contains(strings.ToUpper(text), "INFO") {
		t.Errorf("callout label 'INFO' mancante, text: %q", text)
	}
	if !strings.Contains(text, "Versione agente") {
		t.Errorf("callout title mancante, text: %q", text)
	}
	if !strings.Contains(text, "Dettagli") {
		t.Errorf("callout body mancante, text: %q", text)
	}
}

func TestRender_CalloutMultilineBodyAllVisible(t *testing.T) {
	md := []byte(`# Test

> [!MODIFICA] Aggiornamento Documenti Vincolanti
> **Testo attuale:** "ACCREDIA RT-24 Prove valutative"
> **Testo proposto:** "ACCREDIA RT-39 Rev.00 Prescrizioni"
> **Motivazione normativa:** RT-39 sostituisce RT-24.
> **Criticità:** Alta.
`)
	out := renderOrSkip(t, md, sampleMeta())
	text := readPDFText(t, out)
	for _, want := range []string{"Testo attuale", "Testo proposto", "Motivazione", "Criticità"} {
		if !strings.Contains(text, want) {
			t.Errorf("callout body line %q assente, text: %q", want, text)
		}
	}
}

func TestRender_ItalianAccentsPreserved(t *testing.T) {
	md := []byte("# Test accenti\n\nLa città è già perché può così però ahimè.\n")
	out := renderOrSkip(t, md, sampleMeta())
	text := readPDFText(t, out)
	for _, want := range []string{"città", "perché", "così"} {
		if !strings.Contains(text, want) {
			t.Errorf("parola accentata %q non trovata, text: %q", want, text)
		}
	}
}

func TestRender_TableRendered(t *testing.T) {
	md := []byte(`# Test

| Sezione | Modifica | Impatto |
| ------- | -------- | ------- |
| 6.2     | x        | Media   |
| 7.1     | y        | Alta    |
`)
	out := renderOrSkip(t, md, sampleMeta())
	text := readPDFText(t, out)
	for _, want := range []string{"Sezione", "Modifica", "Impatto", "Media", "Alta"} {
		if !strings.Contains(text, want) {
			t.Errorf("table cell %q assente, text: %q", want, text)
		}
	}
}

// renderDocxOrSkip esegue pdf.RenderDocx e skippa se pandoc non disponibile.
func renderDocxOrSkip(t *testing.T, md []byte, meta pdf.Metadata) []byte {
	t.Helper()
	var buf bytes.Buffer
	err := pdf.RenderDocx(md, meta, pdf.DefaultBrand(), &buf)
	if err != nil {
		if err == pdf.ErrPandocMissing || strings.Contains(err.Error(), "pandoc") {
			t.Skipf("pandoc non disponibile: %v (esegui make pdf-tools)", err)
		}
		t.Fatalf("RenderDocx: %v", err)
	}
	return buf.Bytes()
}

// readDocxText estrae il testo plain dal .docx unendo body + header + footer.
// Necessario perché il branding LabNexus + capability vive in word/header1.xml,
// la paginazione in word/footer1.xml, e il body content in word/document.xml.
func readDocxText(t *testing.T, data []byte) string {
	t.Helper()
	tmp := filepath.Join(t.TempDir(), "out.docx")
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		t.Fatalf("write tempfile: %v", err)
	}
	zr, err := zip.OpenReader(tmp)
	if err != nil {
		t.Fatalf("zip open: %v", err)
	}
	defer zr.Close()
	var sb strings.Builder
	for _, f := range zr.File {
		switch f.Name {
		case "word/document.xml", "word/header1.xml", "word/footer1.xml":
		default:
			continue
		}
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("open %s: %v", f.Name, err)
		}
		body, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatalf("read %s: %v", f.Name, err)
		}
		s := tagRe.ReplaceAllString(string(body), " ")
		sb.WriteString(s)
		sb.WriteString(" ")
	}
	return wsRe.ReplaceAllString(sb.String(), " ")
}

var (
	tagRe = regexp.MustCompile(`<[^>]+>`)
	wsRe  = regexp.MustCompile(`\s+`)
)

func TestRenderDocx_ProducesValidDOCX(t *testing.T) {
	out := renderDocxOrSkip(t, []byte("# Test\n\nCorpo.\n"), sampleMeta())
	// DOCX = zip → magic bytes "PK"
	if len(out) < 2 || string(out[:2]) != "PK" {
		t.Fatalf("output non è uno zip/docx valido, primi bytes: %q", out[:min(4, len(out))])
	}
}

func TestRenderDocx_HeaderContainsBrand(t *testing.T) {
	out := renderDocxOrSkip(t, []byte("# Test\n\nCorpo.\n"), sampleMeta())
	text := readDocxText(t, out)
	if !strings.Contains(text, "LabNexus") {
		t.Errorf("DOCX dovrebbe contenere brand 'LabNexus' (prependDocxHeader), text: %q", text[:min(500, len(text))])
	}
}

func TestRenderDocx_MetadataBoxContainsFrontmatter(t *testing.T) {
	out := renderDocxOrSkip(t, []byte("# Test\n\nCorpo.\n"), sampleMeta())
	text := readDocxText(t, out)
	for _, want := range []string{"revisione", "qwen3.5-122b-a10b", "eurouter", "Revisione SGQ"} {
		if !strings.Contains(text, want) {
			t.Errorf("metadata %q assente nel DOCX, text: %q", want, text[:min(800, len(text))])
		}
	}
}

func TestRenderDocx_CalloutRenderedWithPrefix(t *testing.T) {
	md := []byte("# Test\n\n> [!ATTENZIONE] Scadenza ACCREDIA\n> Entro il 2026-09-30.\n")
	out := renderDocxOrSkip(t, md, sampleMeta())
	text := readDocxText(t, out)
	if !strings.Contains(text, "[ATTENZIONE]") {
		t.Errorf("callout prefix '[ATTENZIONE]' assente, text: %q", text)
	}
	if !strings.Contains(text, "Scadenza ACCREDIA") {
		t.Errorf("callout title 'Scadenza ACCREDIA' assente, text: %q", text)
	}
	if !strings.Contains(text, "2026-09-30") {
		t.Errorf("callout body assente, text: %q", text)
	}
}

func TestRenderDocx_ItalianAccentsPreserved(t *testing.T) {
	out := renderDocxOrSkip(t, []byte("# Test\n\nLa città è già perché può.\n"), sampleMeta())
	text := readDocxText(t, out)
	for _, want := range []string{"città", "perché", "può"} {
		if !strings.Contains(text, want) {
			t.Errorf("accento %q assente nel DOCX, text: %q", want, text)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
