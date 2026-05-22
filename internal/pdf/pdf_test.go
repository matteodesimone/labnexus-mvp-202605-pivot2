package pdf_test

import (
	"bytes"
	"os"
	"path/filepath"
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
