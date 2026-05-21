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

// SupportedExt elenca le estensioni accettate (FR-4 + Sprint 1.5.B FR-21).
var SupportedExt = []string{".md", ".txt", ".csv", ".pdf", ".docx", ".xlsx", ".xls", ".doc", ".rtf"}

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
// labnexusInternalFiles sono file infrastrutturali della cartella di lavoro
// (Sprint 1.5.C concierge mode) che NON vanno parsati come input domain-specific.
// Skip silenzioso: niente warning "formato non supportato".
var labnexusInternalFiles = map[string]bool{
	"_labnexus.toml":  true, // job metadata (FR-26)
	"Esegui.command":  true, // launcher per-capability (Sprint 1.5.C concierge)
	".DS_Store":       true, // macOS metadata
}

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
		if labnexusInternalFiles[e.Name()] {
			continue // skip file infrastrutturali Sprint 1.5.C
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

// ParseTriggerPromptFile carica il file referenziato da `trigger_prompt_file`
// e ne ritorna il plain text (FR-17). Il path è relativo a inputDir e protetto
// contro escape (NFR-6 amendment): il file resolto via EvalSymlinks deve
// restare dentro inputDir. Formati supportati: .md, .txt, .rtf.
func ParseTriggerPromptFile(inputDir, relPath string) (string, error) {
	if relPath == "" {
		return "", fmt.Errorf("input: trigger_prompt_file vuoto")
	}
	absInput, err := filepath.Abs(inputDir)
	if err != nil {
		return "", fmt.Errorf("input: inputDir %q non risolvibile: %w", inputDir, err)
	}
	realInput, err := filepath.EvalSymlinks(absInput)
	if err != nil {
		realInput = absInput
	}
	full := filepath.Join(absInput, relPath)
	absFull, err := filepath.Abs(full)
	if err != nil {
		return "", fmt.Errorf("input: trigger_prompt_file %q non risolvibile: %w", relPath, err)
	}
	if !strings.HasPrefix(absFull, absInput+string(os.PathSeparator)) && absFull != absInput {
		return "", fmt.Errorf("input: trigger_prompt_file %q path non consentito (esce dalla cartella di input)", relPath)
	}
	if _, err := os.Stat(absFull); err != nil {
		return "", fmt.Errorf("input: trigger_prompt_file %q non esiste in inputDir", relPath)
	}
	realFull, err := filepath.EvalSymlinks(absFull)
	if err != nil {
		return "", fmt.Errorf("input: trigger_prompt_file %q symlink non risolvibile: %w", relPath, err)
	}
	if !strings.HasPrefix(realFull, realInput+string(os.PathSeparator)) && realFull != realInput {
		return "", fmt.Errorf("input: trigger_prompt_file %q path non consentito (symlink esce dalla cartella di input)", relPath)
	}
	ext := strings.ToLower(filepath.Ext(absFull))
	switch ext {
	case ".md", ".txt":
		data, err := os.ReadFile(absFull)
		if err != nil {
			return "", err
		}
		return string(data), nil
	case ".rtf":
		return parseRtf(absFull)
	default:
		return "", fmt.Errorf("input: trigger_prompt_file %q estensione non supportata %q (consentite: .md, .txt, .rtf)", relPath, ext)
	}
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
	case ".xls":
		return parseXls(path)
	case ".doc":
		return parseDoc(path)
	case ".rtf":
		return parseRtf(path)
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
