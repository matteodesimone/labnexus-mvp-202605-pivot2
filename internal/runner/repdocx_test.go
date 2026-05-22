package runner

import (
	"os"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/output"
	"github.com/labnexus/labnexus/internal/pdf"
)

// TestRedocx_RegenerateDocxFromExistingMD: one-shot helper analogo a
// TestRepdf_RegeneratePDFFromExistingMD ma per DOCX. Si attiva con
// LABNEXUS_REDOCX_FROM=<path-to-md>.
func TestRedocx_RegenerateDocxFromExistingMD(t *testing.T) {
	mdPath := os.Getenv("LABNEXUS_REDOCX_FROM")
	if mdPath == "" {
		t.Skip("set LABNEXUS_REDOCX_FROM=<path-md> per generare DOCX accoppiato")
	}
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("read %s: %v", mdPath, err)
	}
	content := string(data)
	body := content
	if strings.HasPrefix(content, "---\n") {
		if end := strings.Index(content[4:], "\n---\n"); end >= 0 {
			body = content[4+end+5:]
		}
	}
	body = stripFenceWrapping(body)

	fm := &output.Frontmatter{}
	extract := func(key string) string {
		needle := "\n" + key + ":"
		idx := strings.Index(content, needle)
		if idx < 0 {
			return ""
		}
		line := content[idx+len(needle):]
		if nl := strings.IndexByte(line, '\n'); nl >= 0 {
			line = line[:nl]
		}
		return strings.Trim(strings.TrimSpace(line), `"`)
	}
	fm.Profilo = extract("profilo")
	fm.Modello = extract("modello")
	fm.Provider = extract("provider")
	fm.DataEsecuzione = extract("data_esecuzione")

	docxPath := strings.TrimSuffix(mdPath, ".md") + ".docx"
	f, err := os.Create(docxPath)
	if err != nil {
		t.Fatalf("create %s: %v", docxPath, err)
	}
	defer f.Close()
	meta := pdf.Metadata{
		Profilo:        fm.Profilo,
		Modello:        fm.Modello,
		Provider:       fm.Provider,
		DataEsecuzione: fm.DataEsecuzione,
		Capability:     os.Getenv("LABNEXUS_REDOCX_CAPABILITY"),
		Stato:          "completato",
	}
	if err := pdf.RenderDocx([]byte(body), meta, pdf.DefaultBrand(), f); err != nil {
		t.Fatalf("RenderDocx: %v", err)
	}
	t.Logf("DOCX rigenerato in %s — apri con `open %s`", docxPath, docxPath)
}
