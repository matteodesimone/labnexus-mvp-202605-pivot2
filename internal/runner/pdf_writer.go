package runner

import (
	"os"
	"strings"

	"github.com/labnexus/labnexus/internal/output"
	"github.com/labnexus/labnexus/internal/pdf"
	"github.com/labnexus/labnexus/internal/runlog"
)

// maybeWritePDF è l'helper che genera il PDF accoppiato all'output MD.
// Pattern non-bloccante: se pdfEnabled=false o se il rendering fallisce,
// loggiamo un warning ma il run prosegue normale (il .md è già stato scritto
// dal runner principal). Audit-friendly: il MD resta source of truth.
//
// Path PDF derivato dal mdPath sostituendo l'estensione .md con .pdf.
// Ritorna il path PDF scritto, o "" se non scritto (disabled o errore).
func maybeWritePDF(mdPath, body string, fm *output.Frontmatter, pdfEnabled bool, capability string, log *runlog.Logger) string {
	if !pdfEnabled {
		return ""
	}
	pdfPath := strings.TrimSuffix(mdPath, ".md") + ".pdf"
	meta := pdf.Metadata{
		Profilo:        fm.Profilo,
		Modello:        fm.Modello,
		Provider:       fm.Provider,
		DataEsecuzione: fm.DataEsecuzione,
		DurataSecondi:  fm.DurataSecondi,
		TokenStimati:   fm.TokenStimati,
		FileInput:      fm.FileInput,
		Stato:          fm.Stato,
		Capability:     capability,
	}
	f, err := os.Create(pdfPath)
	if err != nil {
		log.Warn("PDF: impossibile creare file %q: %v — PDF disabilitato per questo run, .md disponibile", pdfPath, err)
		return ""
	}
	if err := pdf.Render([]byte(body), meta, pdf.DefaultBrand(), f); err != nil {
		f.Close()
		_ = os.Remove(pdfPath) // cleanup file parziale
		log.Warn("PDF: render fallito (%v) — PDF disabilitato per questo run, .md disponibile", err)
		return ""
	}
	if err := f.Close(); err != nil {
		log.Warn("PDF: chiusura file %q fallita: %v", pdfPath, err)
		return ""
	}
	return pdfPath
}

// maybeWriteDocx genera il DOCX accoppiato all'MD (analogo a maybeWritePDF).
// Pattern non-bloccante: warning su errore, .md sempre disponibile.
// Path derivato da mdPath sostituendo .md → .docx.
func maybeWriteDocx(mdPath, body string, fm *output.Frontmatter, docxEnabled bool, capability string, log *runlog.Logger) string {
	if !docxEnabled {
		return ""
	}
	docxPath := strings.TrimSuffix(mdPath, ".md") + ".docx"
	meta := pdf.Metadata{
		Profilo:        fm.Profilo,
		Modello:        fm.Modello,
		Provider:       fm.Provider,
		DataEsecuzione: fm.DataEsecuzione,
		DurataSecondi:  fm.DurataSecondi,
		TokenStimati:   fm.TokenStimati,
		FileInput:      fm.FileInput,
		Stato:          fm.Stato,
		Capability:     capability,
	}
	f, err := os.Create(docxPath)
	if err != nil {
		log.Warn("DOCX: impossibile creare file %q: %v — DOCX disabilitato per questo run, .md disponibile", docxPath, err)
		return ""
	}
	if err := pdf.RenderDocx([]byte(body), meta, pdf.DefaultBrand(), f); err != nil {
		f.Close()
		_ = os.Remove(docxPath)
		log.Warn("DOCX: render fallito (%v) — DOCX disabilitato per questo run, .md disponibile", err)
		return ""
	}
	if err := f.Close(); err != nil {
		log.Warn("DOCX: chiusura file %q fallita: %v", docxPath, err)
		return ""
	}
	return docxPath
}
