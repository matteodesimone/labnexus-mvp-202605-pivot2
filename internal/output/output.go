// Package output scrive il file markdown di output con frontmatter YAML (FR-9, EC-15).
package output

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Frontmatter sono i metadati YAML scritti in cima all'output (FR-9).
type Frontmatter struct {
	Profilo          string   `yaml:"profilo"`
	Modello          string   `yaml:"modello"`
	Provider         string   `yaml:"provider"`
	DataEsecuzione   string   `yaml:"data_esecuzione"`
	DurataSecondi    float64  `yaml:"durata_secondi"`
	TokenStimati     int      `yaml:"token_stimati"`
	FileInput        []string `yaml:"file_input"`
	Stato            string   `yaml:"stato,omitempty"`
	ValutazioneDenis string   `yaml:"valutazione_denis,omitempty"`
}

// Write scrive un file markdown in outputDir con il pattern:
//
//	<timestamp>_<profile>_<input_descriptor>.md
//
// dove timestamp è derivato da fm.DataEsecuzione (o time.Now() se vuoto).
// In caso di collisione (EC-15) aggiunge suffisso _2, _3, ... Mai overwrite.
func Write(outputDir string, fm *Frontmatter, body string) (string, error) {
	if fm == nil {
		return "", fmt.Errorf("output: frontmatter nil")
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", fmt.Errorf("output: mkdir %q: %w", outputDir, err)
	}
	ts := timestampFromFrontmatter(fm)
	descriptor := inputDescriptor(fm.FileInput)
	base := fmt.Sprintf("%s_%s", ts, slug(fm.Profilo))
	if descriptor != "" {
		base += "_" + descriptor
	}
	path := uniquePath(outputDir, base, ".md")

	rendered, err := render(fm, body)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(rendered), 0o644); err != nil {
		return "", fmt.Errorf("output: write %q: %w", path, err)
	}
	return path, nil
}

func render(fm *Frontmatter, body string) (string, error) {
	header, err := yaml.Marshal(fm)
	if err != nil {
		return "", fmt.Errorf("output: marshal frontmatter: %w", err)
	}
	var b strings.Builder
	b.WriteString("---\n")
	b.Write(header)
	b.WriteString("---\n\n")
	b.WriteString(body)
	if !strings.HasSuffix(body, "\n") {
		b.WriteString("\n")
	}
	return b.String(), nil
}

var nonSlugRe = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

func slug(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "unknown"
	}
	return nonSlugRe.ReplaceAllString(s, "-")
}

func timestampFromFrontmatter(fm *Frontmatter) string {
	if fm.DataEsecuzione != "" {
		if t, err := time.Parse(time.RFC3339, fm.DataEsecuzione); err == nil {
			return t.Format("2006-01-02T150405")
		}
	}
	return time.Now().Format("2006-01-02T150405")
}

func inputDescriptor(files []string) string {
	if len(files) == 0 {
		return ""
	}
	// usa il nome del primo file (senza estensione) come descrittore
	first := filepath.Base(files[0])
	ext := filepath.Ext(first)
	if ext != "" {
		first = strings.TrimSuffix(first, ext)
	}
	return slug(first)
}

// uniquePath ritorna il primo path libero: base.ext, base_2.ext, base_3.ext, ...
func uniquePath(dir, base, ext string) string {
	p := filepath.Join(dir, base+ext)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return p
	}
	for i := 2; ; i++ {
		p := filepath.Join(dir, fmt.Sprintf("%s_%d%s", base, i, ext))
		if _, err := os.Stat(p); os.IsNotExist(err) {
			return p
		}
		if i > 9999 {
			return p // safety net
		}
	}
}
