// Package jobs gestisce auto-discovery delle cartelle di lavoro Denis (FR-25..FR-30).
//
// Una cartella sotto `lavori/<X>/` è un "job" se contiene `_labnexus.toml` con
// almeno `profile = "..."` (FR-26). Se manca il metadata, applica fallback
// convention naming via regex `Profilo (.+)$` sul folder name (FR-27).
package jobs

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
)

// Job rappresenta una cartella di lavoro auto-discovered in `lavori/`.
type Job struct {
	Name              string // nome del job (default: nome cartella)
	Path              string // path assoluto della cartella `lavori/<X>/`
	Profile           string // nome profile da `_labnexus.toml` o fallback convention
	TriggerPromptFile string // opzionale, path relativo a Path
	InputDir          string // = Path (la cartella stessa è l'input)
	OutputDir         string // = Path + "/output" (default)
	Source            string // "metadata" | "convention" — come è stato risolto
	Warning           string // se non vuoto, semantic-level warning (es. profile inesistente in profili/)
}

// ValidateProfileExists controlla che ogni job referenzi un profile esistente
// in profiliDir. Popola Job.Warning per i job orfani. Pattern: fix review 1.5.C
// HIGH-Claude-1 (jobs.Discover non valida profile esistente). Helper esterno
// per permettere a Discover di rimanere pure-filesystem; la validazione semantica
// è opzionale e gestita dal caller (es. cmd/jobs subcommand).
func ValidateProfileExists(js []*Job, profiliDir string) {
	for _, j := range js {
		profilePath := filepath.Join(profiliDir, j.Profile+".toml")
		if _, err := os.Stat(profilePath); err != nil {
			j.Warning = "profile '" + j.Profile + "' non trovato in " + profiliDir
		}
	}
}

// metadataFileName è il nome canonico del file metadata per cartella.
const metadataFileName = "_labnexus.toml"

// conventionRegex estrae il nome profilo da folder con pattern "... — Profilo <nome>".
var conventionRegex = regexp.MustCompile(`Profilo\s+(.+?)\s*$`)

// metadataDoc è la rappresentazione TOML del file `_labnexus.toml`.
type metadataDoc struct {
	Profile           string `toml:"profile"`
	TriggerPromptFile string `toml:"trigger_prompt_file"`
}

// Discover scansiona `lavoriDir`, legge `_labnexus.toml` per ogni sotto-cartella,
// applica fallback convention naming se assente (FR-27). Cartelle malformate
// o senza match skippate (FR-26/FR-27 amendment). Ordinato per Name.
//
// Fix review 1.5.C CRITICAL-mistral-2: symlink escape protection. Le cartelle
// in `lavori/` che sono symlink fuori da `lavoriDir` (post-EvalSymlinks) sono
// skippate per defense-in-depth. Riusa pattern profile.validateKbFiles (Sprint 1
// bugfix #002 KB symlink escape).
func Discover(lavoriDir string) ([]*Job, error) {
	entries, err := os.ReadDir(lavoriDir)
	if err != nil {
		return nil, fmt.Errorf("jobs: lettura cartella %q: %w", lavoriDir, err)
	}
	absLavori, _ := filepath.Abs(lavoriDir)
	realLavori, err := filepath.EvalSymlinks(absLavori)
	if err != nil {
		realLavori = absLavori
	}
	var out []*Job
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		job := resolveJob(lavoriDir, e.Name(), realLavori)
		if job != nil {
			out = append(out, job)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// resolveJob analizza una singola cartella e ritorna *Job se è un job valido,
// altrimenti nil (cartella skippata). Skip se symlink esce da realLavoriDir.
func resolveJob(lavoriDir, folderName, realLavoriDir string) *Job {
	jobPath := filepath.Join(lavoriDir, folderName)
	if !isWithin(jobPath, realLavoriDir) {
		return nil // symlink escape → skip (defense-in-depth)
	}
	metaPath := filepath.Join(jobPath, metadataFileName)
	if _, err := os.Stat(metaPath); err == nil {
		profile, triggerFile, err := LoadMetadata(metaPath)
		if err != nil {
			return nil // metadata malformato → skip
		}
		return newJob(folderName, jobPath, profile, triggerFile, "metadata")
	}
	// Fallback convention naming
	matches := conventionRegex.FindStringSubmatch(folderName)
	if matches == nil {
		return nil // né metadata né convention → skip
	}
	return newJob(folderName, jobPath, strings.TrimSpace(matches[1]), "", "convention")
}

// isWithin verifica che il path resolto via EvalSymlinks resti dentro
// realLavoriDir. Defense-in-depth contro symlink escape (NFR-10 amendment).
// Fix sub-fix review 1.5.C loop 2: assolutizza jobPath prima di EvalSymlinks
// (era bug subtle con path relativi che faceva sempre fallire prefix check).
func isWithin(jobPath, realLavoriDir string) bool {
	absJob, err := filepath.Abs(jobPath)
	if err != nil {
		return false
	}
	realJob, err := filepath.EvalSymlinks(absJob)
	if err != nil {
		// EvalSymlinks fail (non-symlink ok, broken link no): fallback lexical check
		return strings.HasPrefix(absJob, realLavoriDir+string(os.PathSeparator)) || absJob == realLavoriDir
	}
	return strings.HasPrefix(realJob, realLavoriDir+string(os.PathSeparator)) || realJob == realLavoriDir
}

// newJob crea un Job con i campi di default popolati.
func newJob(folderName, jobPath, profile, triggerFile, source string) *Job {
	return &Job{
		Name:              folderName,
		Path:              jobPath,
		Profile:           profile,
		TriggerPromptFile: triggerFile,
		InputDir:          jobPath,
		OutputDir:         filepath.Join(jobPath, "output"),
		Source:            source,
	}
}

// Resolve cerca un job per nome. Strategia: case-insensitive substring match
// sul Name (la cartella). Ritorna il primo match o nil.
func Resolve(jobs []*Job, name string) *Job {
	needle := strings.ToLower(strings.TrimSpace(name))
	if needle == "" {
		return nil
	}
	for _, j := range jobs {
		if strings.Contains(strings.ToLower(j.Name), needle) {
			return j
		}
	}
	return nil
}

// LoadMetadata legge un singolo `_labnexus.toml` e ritorna
// (profile, triggerPromptFile, error). Errore se il file non esiste o
// il TOML è malformato.
//
// Fix review 1.5.C CRITICAL-mistral-1: defense-in-depth lexical validation
// di TriggerPromptFile contro path traversal. Pattern coerente con
// profile.validateTrigger (Sprint 1.5.B). Runtime check resta in
// input.ParseTriggerPromptFile come second-line.
func LoadMetadata(path string) (string, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("jobs: read metadata %q: %w", path, err)
	}
	var m metadataDoc
	if err := toml.Unmarshal(data, &m); err != nil {
		return "", "", fmt.Errorf("jobs: parse metadata %q: %w", path, err)
	}
	if strings.TrimSpace(m.Profile) == "" {
		return "", "", fmt.Errorf("jobs: metadata %q manca campo 'profile'", path)
	}
	if m.TriggerPromptFile != "" {
		if strings.Contains(m.TriggerPromptFile, "..") {
			return "", "", fmt.Errorf("jobs: metadata %q ha trigger_prompt_file %q non consentito (contiene '..')", path, m.TriggerPromptFile)
		}
		if filepath.IsAbs(m.TriggerPromptFile) {
			return "", "", fmt.Errorf("jobs: metadata %q ha trigger_prompt_file %q non consentito (path assoluto)", path, m.TriggerPromptFile)
		}
	}
	return m.Profile, m.TriggerPromptFile, nil
}
