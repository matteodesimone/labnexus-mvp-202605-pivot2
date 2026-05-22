// Package pdf genera il PDF accoppiato all'output MD via pipeline
// pandoc + typst (Sprint 1.5.D, pivot da goldmark+fpdf).
//
// Pipeline:
//
//	1. Pre-processing MD: (nessuno per ora; pandoc + lua filter coprono tutto)
//	2. pandoc MD → typst body, con lua filter che converte i callout Obsidian
//	   `> [!TYPE] ...` in raw typst blocks `#callout(kind, title)[body]`
//	3. Concatena template LabNexus + body typst → file `.typ` completo
//	4. typst compile --input <meta>=<val> ... → PDF
//
// Binari pandoc + typst sono scaricati da scripts/download-pdf-tools.sh e
// posizionati in `<dist>/bin/`. Path risolto a runtime tramite os.Executable.
package pdf

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

//go:embed template.typ
var templateTypst string

//go:embed callout-filter.lua
var calloutFilterLua string

// Metadata sono i campi del frontmatter passati al template Typst tramite
// `typst compile --input key=value`. Mapping 1:1 con le variabili `sys.inputs`
// nel template.typ.
type Metadata struct {
	Profilo        string
	Modello        string
	Provider       string
	DataEsecuzione string
	DurataSecondi  float64
	TokenStimati   int
	FileInput      []string
	Stato          string
	Capability     string
}

// BrandConfig è retained per backwards-compat API. Sprint 1.5.D non lo usa
// (il branding vive nel template.typ embedded). Iterazione futura: override
// di palette via questo struct propagato come typst --input.
type BrandConfig struct {
	TitlePrefix string
	FooterText  string
}

// DefaultBrand ritorna defaults LabNexus. No-op effettivo per Sprint 1.5.D.
func DefaultBrand() BrandConfig {
	return BrandConfig{TitlePrefix: "LabNexus"}
}

// ErrToolsMissing è restituito quando pandoc o typst non sono trovati in <bin>/.
var ErrToolsMissing = errors.New("pdf: pandoc e/o typst non trovati nella cartella bin/ accanto al binary labnexus (esegui scripts/download-pdf-tools.sh)")

// Render genera il PDF a partire da `md` (body markdown senza frontmatter
// engine) + `meta`. Pipeline pandoc → typst, niente runtime esterno richiesto.
func Render(md []byte, meta Metadata, _ BrandConfig, out io.Writer) error {
	pandocPath, typstPath, err := resolveBinaries()
	if err != nil {
		return err
	}

	tmpDir, err := os.MkdirTemp("", "labnexus-pdf-*")
	if err != nil {
		return fmt.Errorf("pdf: tempdir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// 1. scrivi lua filter su tempfile
	filterPath := filepath.Join(tmpDir, "callout-filter.lua")
	if err := os.WriteFile(filterPath, []byte(calloutFilterLua), 0o644); err != nil {
		return fmt.Errorf("pdf: write filter: %w", err)
	}

	// 2. invoca pandoc MD → typst body
	bodyTypst, err := runPandocToTypst(pandocPath, filterPath, md)
	if err != nil {
		return fmt.Errorf("pdf: pandoc: %w", err)
	}

	// 3. componi typst file completo (template + body)
	typstSource := templateTypst + "\n\n" + bodyTypst
	typPath := filepath.Join(tmpDir, "main.typ")
	if err := os.WriteFile(typPath, []byte(typstSource), 0o644); err != nil {
		return fmt.Errorf("pdf: write typ: %w", err)
	}

	// 4. invoca typst compile → PDF
	pdfPath := filepath.Join(tmpDir, "out.pdf")
	if err := runTypstCompile(typstPath, typPath, pdfPath, meta); err != nil {
		return fmt.Errorf("pdf: typst: %w", err)
	}

	// 5. copia PDF a out
	data, err := os.ReadFile(pdfPath)
	if err != nil {
		return fmt.Errorf("pdf: read pdf: %w", err)
	}
	if _, err := out.Write(data); err != nil {
		return fmt.Errorf("pdf: write out: %w", err)
	}
	return nil
}

// resolveBinaries cerca pandoc + typst in <dir-del-binary-labnexus>/bin/.
// Fallback: variabili env LABNEXUS_PANDOC_PATH / LABNEXUS_TYPST_PATH (utili
// per test e per smoke da repo root dove il binary labnexus non è in dist/).
func resolveBinaries() (string, string, error) {
	if p := os.Getenv("LABNEXUS_PANDOC_PATH"); p != "" {
		if t := os.Getenv("LABNEXUS_TYPST_PATH"); t != "" {
			return p, t, nil
		}
	}
	exe, err := os.Executable()
	if err != nil {
		return "", "", fmt.Errorf("pdf: os.Executable: %w", err)
	}
	exe, _ = filepath.EvalSymlinks(exe)
	binDir := filepath.Join(filepath.Dir(exe), "bin")
	pandocPath := filepath.Join(binDir, "pandoc")
	typstPath := filepath.Join(binDir, "typst")
	if !isExecutable(pandocPath) || !isExecutable(typstPath) {
		// Fallback: cerca nella dist canonica del repo (utile durante dev/test)
		repoFallback, ok := findRepoDistBin()
		if ok {
			pandocPath = filepath.Join(repoFallback, "pandoc")
			typstPath = filepath.Join(repoFallback, "typst")
			if isExecutable(pandocPath) && isExecutable(typstPath) {
				return pandocPath, typstPath, nil
			}
		}
		return "", "", ErrToolsMissing
	}
	return pandocPath, typstPath, nil
}

// findRepoDistBin cerca dist/labnexus-sprint1-darwin-arm64/bin risalendo dalla
// CWD. Per dev/test invocati da varie directory del repo.
func findRepoDistBin() (string, bool) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	for dir := cwd; dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		candidate := filepath.Join(dir, "dist", "labnexus-sprint1-darwin-arm64", "bin")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}
	}
	return "", false
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0o111 != 0
}

func runPandocToTypst(pandocPath, filterPath string, md []byte) (string, error) {
	cmd := exec.Command(pandocPath,
		"-f", "markdown+raw_html+fenced_divs+pipe_tables+task_lists",
		"-t", "typst",
		"--lua-filter", filterPath,
		"--wrap=none",
	)
	cmd.Stdin = bytes.NewReader(md)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w (stderr: %s)", err, strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), nil
}

func runTypstCompile(typstPath, typPath, pdfPath string, meta Metadata) error {
	args := []string{"compile", typPath, pdfPath}
	// Passa i metadata come sys.inputs (Typst li espone come stringhe)
	add := func(k, v string) {
		if v == "" {
			return
		}
		args = append(args, "--input", k+"="+v)
	}
	add("profilo", meta.Profilo)
	add("modello", meta.Modello)
	add("provider", meta.Provider)
	add("data_esecuzione", meta.DataEsecuzione)
	add("capability", meta.Capability)
	add("stato", meta.Stato)
	if meta.DurataSecondi > 0 {
		add("durata_secondi", strconv.FormatFloat(meta.DurataSecondi, 'f', 1, 64))
	}
	if meta.TokenStimati > 0 {
		add("token_stimati", formatThousands(meta.TokenStimati))
	}
	if len(meta.FileInput) > 0 {
		add("file_input", formatFileInputList(meta.FileInput))
	}

	cmd := exec.Command(typstPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w (stderr: %s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// formatThousands rende un int con separatore "." (convenzione italiana).
func formatThousands(n int) string {
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return s
	}
	var b strings.Builder
	rem := len(s) % 3
	if rem > 0 {
		b.WriteString(s[:rem])
		if len(s) > rem {
			b.WriteByte('.')
		}
	}
	for i := rem; i < len(s); i += 3 {
		b.WriteString(s[i : i+3])
		if i+3 < len(s) {
			b.WriteByte('.')
		}
	}
	return b.String()
}

func formatFileInputList(files []string) string {
	if len(files) == 0 {
		return ""
	}
	if len(files) == 1 {
		return filepath.Base(files[0])
	}
	return fmt.Sprintf("%d file (%s + %d altri)", len(files), filepath.Base(files[0]), len(files)-1)
}
