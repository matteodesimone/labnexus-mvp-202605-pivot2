// Package input estrae testo plain dai file della cartella di input dell'utente.
// Walk non ricorsivo (FR-4). Formati supportati: md, txt, csv, pdf, docx, xlsx.
// Su parser fallito → skip con warning (EC-1, EC-2). Su ≥ 50% fallimenti → errore.
package input

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/ledongthuc/pdf"
)

// SupportedExt elenca le estensioni accettate (FR-4).
var SupportedExt = []string{".md", ".txt", ".csv", ".pdf", ".docx", ".xlsx"}

// ParsedFile rappresenta un file di input dopo estrazione testo.
type ParsedFile struct {
	Name string
	Text string
	Size int
}

// ParsedResult è l'esito del walk + parsing.
type ParsedResult struct {
	Files        []ParsedFile
	Skipped      []SkippedFile
	FailureRatio float64
}

// SkippedFile è un file non incluso (formato non supportato o parsing fallito).
type SkippedFile struct {
	Name   string
	Reason string
}

// ParseDir esegue il walk non ricorsivo di dir, estrae testo per ogni file supportato.
func ParseDir(dir string) (*ParsedResult, error) {
	names, err := listFilesDeterministic(dir)
	if err != nil {
		return nil, err
	}
	if len(names) == 0 {
		return &ParsedResult{}, errors.New("input: cartella di input vuota")
	}
	res, supported, parsed := extractAllFiles(dir, names)
	if supported == 0 {
		return res, fmt.Errorf("input: nessun file con estensione supportata (consentite: %s)", strings.Join(SupportedExt, ", "))
	}
	failed := supported - parsed
	res.FailureRatio = float64(failed) / float64(supported)
	// Soglia strict-majority (EC-1): > 50% → errore (1 su 2 = warning, esecuzione prosegue).
	if res.FailureRatio > 0.5 {
		return res, fmt.Errorf("input: oltre il 50%% dei file supportati hanno fallito il parsing (%d/%d)", failed, supported)
	}
	return res, nil
}

// listFilesDeterministic ritorna i nomi file (non dir) ordinati alfabeticamente.
func listFilesDeterministic(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("input: lettura cartella %q: %w", dir, err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue // walk non ricorsivo
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)
	return names, nil
}

// extractAllFiles itera i file e popola ParsedResult; ritorna supported/parsed counts.
func extractAllFiles(dir string, names []string) (*ParsedResult, int, int) {
	res := &ParsedResult{}
	supported, parsed := 0, 0
	for _, name := range names {
		ext := strings.ToLower(filepath.Ext(name))
		if !isSupported(ext) {
			res.Skipped = append(res.Skipped, SkippedFile{Name: name, Reason: "formato non supportato"})
			continue
		}
		supported++
		text, err := extract(filepath.Join(dir, name), ext)
		if err != nil {
			res.Skipped = append(res.Skipped, SkippedFile{Name: name, Reason: err.Error()})
			continue
		}
		parsed++
		res.Files = append(res.Files, ParsedFile{Name: name, Text: text, Size: len(text)})
	}
	return res, supported, parsed
}

func isSupported(ext string) bool {
	for _, s := range SupportedExt {
		if s == ext {
			return true
		}
	}
	return false
}

func extract(path, ext string) (string, error) {
	switch ext {
	case ".md", ".txt", ".csv":
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(data), nil
	case ".pdf":
		return extractPDF(path)
	case ".docx":
		return extractZipXML(path, "word/document.xml")
	case ".xlsx":
		return extractXLSX(path)
	}
	return "", fmt.Errorf("estensione %s non gestita", ext)
}

func extractPDF(path string) (string, error) {
	// Pre-check magic bytes: se non comincia con "%PDF-", non è un PDF valido.
	// pdf.Open può silenziosamente succeedere su contenuti non-PDF lasciando 0 pagine.
	head, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("PDF non leggibile: %w", err)
	}
	if len(head) < 4 || string(head[:4]) != "%PDF" {
		return "", fmt.Errorf("PDF non parsabile — magic bytes mancanti (suggerimento: converti manualmente con pandoc)")
	}
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("PDF non parsabile (suggerimento: converti manualmente con pandoc): %w", err)
	}
	defer f.Close()
	var sb strings.Builder
	totalPage := r.NumPage()
	for i := 1; i <= totalPage; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("PDF pagina %d non estraibile: %w", i, err)
		}
		sb.WriteString(text)
		sb.WriteString("\n")
	}
	return sb.String(), nil
}

func extractZipXML(path, entryName string) (string, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return "", fmt.Errorf("file non leggibile come zip: %w", err)
	}
	defer zr.Close()
	for _, f := range zr.File {
		if f.Name == entryName {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			data, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return "", err
			}
			return stripXMLTags(string(data)), nil
		}
	}
	return "", fmt.Errorf("entry %q non trovata nel file", entryName)
}

// extractXLSX legge tutti i sharedStrings + xl/worksheets/sheet*.xml e ritorna il testo plain.
// Sufficiente per Sprint 1 (FR-4); non interpreta formule, formattazione, ecc.
func extractXLSX(path string) (string, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return "", fmt.Errorf("file non leggibile come zip: %w", err)
	}
	defer zr.Close()
	var sb strings.Builder
	for _, f := range zr.File {
		if !isXLSXEntry(f.Name) {
			continue
		}
		if err := appendZipEntryText(&sb, f); err != nil {
			return "", err
		}
	}
	if sb.Len() == 0 {
		return "", errors.New("XLSX: nessun foglio leggibile")
	}
	return sb.String(), nil
}

func isXLSXEntry(name string) bool {
	if !strings.HasPrefix(name, "xl/") || !strings.HasSuffix(name, ".xml") {
		return false
	}
	return strings.Contains(name, "worksheets/sheet") || strings.HasSuffix(name, "sharedStrings.xml")
}

func appendZipEntryText(sb *strings.Builder, f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	data, err := io.ReadAll(rc)
	rc.Close()
	if err != nil {
		return err
	}
	sb.WriteString(stripXMLTags(string(data)))
	sb.WriteString("\n")
	return nil
}

var xmlTagRe = regexp.MustCompile(`<[^>]+>`)
var multiWS = regexp.MustCompile(`[ \t]+`)
var multiNL = regexp.MustCompile(`\n{3,}`)

func stripXMLTags(s string) string {
	s = xmlTagRe.ReplaceAllString(s, " ")
	s = multiWS.ReplaceAllString(s, " ")
	s = multiNL.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}
