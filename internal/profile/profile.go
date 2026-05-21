// Package profile carica e valida i file TOML di profilo LabNexus (FR-3, Sprint 1.5.B).
package profile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

// MinTriggerLen è la lunghezza minima del trigger_prompt inline (FR-3).
const MinTriggerLen = 50

// AllowedProviders elenca i provider supportati (FR-3, FR-7).
var AllowedProviders = []string{"ollama", "eurouter"}

// Profile è lo schema di un file profili/<nome>.toml.
//
// Sprint 1.5.B (post-pivot-3): Provider, Modello, Temperature, MaxTokens,
// ContextWindow sono OPZIONALI nello schema (zero-value detection); se vuoti,
// vengono ereditati dal config master via `internal/config.Merge`. Validate
// va chiamato sul profilo MERGED, non isolato.
//
// TriggerPrompt e TriggerPromptFile sono mutually exclusive (XOR): uno solo
// dei due deve essere setted (FR-16).
type Profile struct {
	Profilo           string   `toml:"profilo"`
	Descrizione       string   `toml:"descrizione"`
	Provider          string   `toml:"provider"`
	Modello           string   `toml:"modello"`
	KbFiles           []string `toml:"kb_files"`
	TriggerPrompt     string   `toml:"trigger_prompt"`
	TriggerPromptFile string   `toml:"trigger_prompt_file"`
	Temperature       float64  `toml:"temperature"`
	MaxTokens         int      `toml:"max_tokens"`
	ContextWindow     int      `toml:"context_window"`
	Output            Output   `toml:"output"`
}

// Output configura come l'engine scrive l'output del profilo. Bug #006:
// FrontmatterDefault sono chiavi/valori che vanno mergeati nel frontmatter
// dell'output (le chiavi engine-generated vincono in caso di collisione).
type Output struct {
	FrontmatterDefault map[string]string `toml:"frontmatter_default"`
}

// Load legge e parsa un file di profilo dal path indicato.
func Load(path string) (*Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("profile: lettura %q: %w", path, err)
	}
	var p Profile
	if err := toml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("profile: parse TOML %q: %w", path, err)
	}
	return &p, nil
}

// Validate verifica che il profilo rispetti lo schema (FR-3):
//   - tutti i campi obbligatori presenti (post-merge)
//   - provider in AllowedProviders
//   - trigger_prompt XOR trigger_prompt_file (FR-16)
//   - trigger_prompt ≥ MinTriggerLen se inline
//   - tutti i kb_files esistono in kbDir (path relativo a kbDir)
//
// Path traversal protetto su kb_files: ogni file risolto deve restare dentro kbDir (NFR-6).
// Path traversal su trigger_prompt_file è verificato a runtime da input.ParseTriggerPromptFile
// (relativo a input dir, non a kbDir).
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
	if err := validateTrigger(p); err != nil {
		return err
	}
	return validateKbFiles(p.KbFiles, kbDir)
}

// validateTrigger applica la regola XOR su trigger_prompt vs trigger_prompt_file (FR-16):
//   - entrambi setted → errore "mutually exclusive"
//   - entrambi assenti → errore "trigger_prompt obbligatorio"
//   - solo trigger_prompt inline → minLen check
//   - solo trigger_prompt_file → path-safety lexical check (defense-in-depth)
//     + runtime symlink check via input.ParseTriggerPromptFile
//
// Review 1.5.B HIGH-mistral-1 fix: aggiunto static lexical check qui per
// far emergere path traversal a `labnexus validate` invece che a runtime.
func validateTrigger(p *Profile) error {
	hasInline := strings.TrimSpace(p.TriggerPrompt) != ""
	hasFile := strings.TrimSpace(p.TriggerPromptFile) != ""
	if hasInline && hasFile {
		return errors.New("profile: trigger_prompt e trigger_prompt_file mutually exclusive (specifica uno solo dei due)")
	}
	if !hasInline && !hasFile {
		return errors.New("profile: trigger_prompt obbligatorio (inline o file)")
	}
	if hasInline && len(p.TriggerPrompt) < MinTriggerLen {
		return fmt.Errorf("profile: trigger_prompt troppo corto (%d caratteri), almeno %d caratteri richiesti", len(p.TriggerPrompt), MinTriggerLen)
	}
	if hasFile {
		f := p.TriggerPromptFile
		if strings.Contains(f, "..") {
			return fmt.Errorf("profile: trigger_prompt_file %q path non consentito (contiene '..' — defense-in-depth lexical check)", f)
		}
		if filepath.IsAbs(f) {
			return fmt.Errorf("profile: trigger_prompt_file %q path non consentito (path assoluto — deve essere relativo a --input)", f)
		}
	}
	return nil
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
// Non valida i profili — `validate` è un comando separato. Sprint 1.5.B: cerca .toml.
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
		if !strings.HasSuffix(strings.ToLower(name), ".toml") {
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
