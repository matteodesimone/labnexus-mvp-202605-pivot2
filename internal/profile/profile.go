// Package profile carica e valida i file YAML di profilo LabNexus (FR-3).
package profile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// MinTriggerLen è la lunghezza minima del trigger_prompt (FR-3).
const MinTriggerLen = 50

// AllowedProviders elenca i provider supportati (FR-3, FR-7).
var AllowedProviders = []string{"ollama", "eurouter"}

// Profile è lo schema di un file profili/<nome>.yml.
type Profile struct {
	Profilo       string   `yaml:"profilo"`
	Descrizione   string   `yaml:"descrizione"`
	Provider      string   `yaml:"provider"`
	Modello       string   `yaml:"modello"`
	KbFiles       []string `yaml:"kb_files"`
	TriggerPrompt string   `yaml:"trigger_prompt"`
	Temperature   float64  `yaml:"temperature,omitempty"`
	MaxTokens     int      `yaml:"max_tokens,omitempty"`
	ContextWindow int      `yaml:"context_window,omitempty"`
	Output        Output   `yaml:"output,omitempty"`
}

// Output configura come l'engine scrive l'output del profilo. Bug #006:
// FrontmatterDefault sono chiavi/valori che vanno mergeati nel frontmatter
// dell'output (le chiavi engine-generated vincono in caso di collisione).
type Output struct {
	FrontmatterDefault map[string]string `yaml:"frontmatter_default,omitempty"`
}

// Load legge e parsa un file di profilo dal path indicato.
func Load(path string) (*Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("profile: lettura %q: %w", path, err)
	}
	var p Profile
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("profile: parse YAML %q: %w", path, err)
	}
	return &p, nil
}

// Validate verifica che il profilo rispetti lo schema (FR-3):
//   - tutti i campi obbligatori presenti
//   - provider in AllowedProviders
//   - trigger_prompt ≥ MinTriggerLen
//   - tutti i kb_files esistono in kbDir (path relativo a kbDir).
//
// Path traversal protetto: ogni kb_file risolto deve restare dentro kbDir (NFR-6).
func Validate(p *Profile, kbDir string) error {
	if p == nil {
		return errors.New("profile: nil")
	}
	if strings.TrimSpace(p.Profilo) == "" {
		return errors.New("profile: campo 'profilo' obbligatorio")
	}
	if strings.TrimSpace(p.Descrizione) == "" {
		return errors.New("profile: campo 'descrizione' obbligatorio")
	}
	if strings.TrimSpace(p.Provider) == "" {
		return errors.New("profile: campo 'provider' obbligatorio")
	}
	if !contains(AllowedProviders, p.Provider) {
		return fmt.Errorf("profile: provider %q non ammesso (consentiti: %s)", p.Provider, strings.Join(AllowedProviders, ", "))
	}
	if strings.TrimSpace(p.Modello) == "" {
		return errors.New("profile: campo 'modello' obbligatorio")
	}
	if len(p.KbFiles) == 0 {
		return errors.New("profile: campo 'kb_files' deve avere almeno un file")
	}
	if len(p.TriggerPrompt) < MinTriggerLen {
		return fmt.Errorf("profile: trigger_prompt troppo corto (%d caratteri), almeno %d caratteri richiesti", len(p.TriggerPrompt), MinTriggerLen)
	}
	return validateKbFiles(p.KbFiles, kbDir)
}

func validateKbFiles(kbFiles []string, kbDir string) error {
	absKb, err := filepath.Abs(kbDir)
	if err != nil {
		return fmt.Errorf("profile: kb dir %q non risolvibile: %w", kbDir, err)
	}
	// Bug #002 fix (.pipeline/bugs/kb-symlink-escape.md): risolve kbDir
	// dietro eventuali symlink (la KB-ispettore al root del repo È un
	// symlink legittimo → docs/.../KB-ispettore). Confronto post-resolve.
	realKb, err := filepath.EvalSymlinks(absKb)
	if err != nil {
		// Se kbDir non esiste affatto, fallback su absKb (l'errore di
		// "non esiste" emerge poi sui singoli file).
		realKb = absKb
	}
	for _, rel := range kbFiles {
		if strings.TrimSpace(rel) == "" {
			return errors.New("profile: kb_files contiene voce vuota")
		}
		full := filepath.Join(kbDir, rel)
		absFull, err := filepath.Abs(full)
		if err != nil {
			return fmt.Errorf("profile: kb_file %q non risolvibile: %w", rel, err)
		}
		// NFR-6: path traversal lessicale (defense in depth).
		if !strings.HasPrefix(absFull, absKb+string(os.PathSeparator)) && absFull != absKb {
			return fmt.Errorf("profile: kb_file %q esce dalla cartella KB-ispettore", rel)
		}
		// Esistenza (segue symlink).
		if _, err := os.Stat(absFull); err != nil {
			// Path relativo nel messaggio per uniformità coi feature scenari (NFR-4 user-facing).
			return fmt.Errorf("profile: kb_file %q non esiste in ./KB-ispettore/", rel)
		}
		// Bug #002 fix: risolve symlink sul candidato e verifica che il
		// path REALE rimanga dentro la kbDir reale. Senza questo check, un
		// symlink dentro la KB potrebbe puntare a file privati esterni e
		// iniettarne il contenuto nel prompt LLM (privacy / data leak).
		realFull, err := filepath.EvalSymlinks(absFull)
		if err != nil {
			return fmt.Errorf("profile: kb_file %q non risolvibile via symlink: %w", rel, err)
		}
		if !strings.HasPrefix(realFull, realKb+string(os.PathSeparator)) && realFull != realKb {
			return fmt.Errorf("profile: kb_file %q è un symlink che esce dalla cartella KB-ispettore (real path %q outside %q)", rel, realFull, realKb)
		}
	}
	return nil
}

// List ritorna l'elenco dei profili in profiliDir, ordinati per nome (per `labnexus list`).
// Non valida i profili — `validate` è un comando separato.
func List(profiliDir string) ([]*Profile, error) {
	entries, err := os.ReadDir(profiliDir)
	if err != nil {
		return nil, fmt.Errorf("profile: lettura cartella %q: %w", profiliDir, err)
	}
	var out []*Profile
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".yml") && !strings.HasSuffix(strings.ToLower(name), ".yaml") {
			continue
		}
		p, err := Load(filepath.Join(profiliDir, name))
		if err != nil {
			continue // skip profili non parsabili in list, /validate li rivela
		}
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Profilo < out[j].Profilo })
	return out, nil
}

func contains(s []string, x string) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}
