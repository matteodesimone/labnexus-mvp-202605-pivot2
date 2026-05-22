package runner

import (
	"os"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/output"
	"github.com/labnexus/labnexus/internal/pdf"
)

// TestRepdf_RegeneratePDFFromExistingMD è un helper one-shot per rigenerare
// un PDF dal MD esistente quando un fix di rendering è stato applicato e
// vogliamo evitare di rispendere gettoni eurouter per il run completo.
// Si attiva solo con LABNEXUS_REPDF_FROM=<path-to-md>.
func TestRepdf_RegeneratePDFFromExistingMD(t *testing.T) {
	mdPath := os.Getenv("LABNEXUS_REPDF_FROM")
	if mdPath == "" {
		t.Skip("set LABNEXUS_REPDF_FROM=<path-md> per rigenerare PDF accoppiato")
	}
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("read %s: %v", mdPath, err)
	}
	content := string(data)

	// Strip frontmatter engine (---\n...---\n) per estrarre il body
	body := content
	if strings.HasPrefix(content, "---\n") {
		if end := strings.Index(content[4:], "\n---\n"); end >= 0 {
			body = content[4+end+5:]
		}
	}
	// Applica il sanitizer (rimuove fence wrapping spurio se presente)
	body = stripFenceWrapping(body)

	// Costruisci Metadata best-effort dal frontmatter MD parsing rough
	// (cerchiamo le key engine note)
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

	pdfPath := strings.TrimSuffix(mdPath, ".md") + ".pdf"
	f, err := os.Create(pdfPath)
	if err != nil {
		t.Fatalf("create %s: %v", pdfPath, err)
	}
	defer f.Close()
	meta := pdf.Metadata{
		Profilo:        fm.Profilo,
		Modello:        fm.Modello,
		Provider:       fm.Provider,
		DataEsecuzione: fm.DataEsecuzione,
		Capability:     os.Getenv("LABNEXUS_REPDF_CAPABILITY"),
		Stato:          "completato",
	}
	if err := pdf.Render([]byte(body), meta, pdf.DefaultBrand(), f); err != nil {
		t.Fatalf("Render: %v", err)
	}
	t.Logf("PDF rigenerato in %s — apri con `open %s`", pdfPath, pdfPath)
}
