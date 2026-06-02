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

	"github.com/labnexus/labnexus/internal/profile"
)

// Job rappresenta una cartella di lavoro auto-discovered in `lavori/`.
type Job struct {
	Name              string // nome del job (default: nome cartella)
	Path              string // path assoluto della cartella `lavori/<X>/`
	Profile           string // nome profile da `_labnexus.toml` o fallback convention
	TriggerPrompt     string   // opzionale, override trigger inline (XOR con TriggerPromptFile)
	TriggerPromptFile string   // opzionale, override trigger da file, path relativo a Path
	Exclude           []string // opzionale, file (rel a Path) da NON inviare al modello
	InputDir          string   // = Path (la cartella stessa è l'input)
	OutputDir         string // = Path + "/output" (default)
	Source            string // "metadata" | "convention" — come è stato risolto
	Warning           string // se non vuoto, semantic-level warning (es. profile inesistente in profili/)
	// PDFEnabled override per-capability del default master [pdf].
	// nil = non setted (eredita da master); risolto via config.ResolvePDFEnabled.
	PDFEnabled *bool
	// DocxEnabled override per-capability del default master [docx].
	DocxEnabled *bool
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
	Profile           string          `toml:"profile"`
	TriggerPrompt     string          `toml:"trigger_prompt"`
	TriggerPromptFile string          `toml:"trigger_prompt_file"`
	Exclude           []string        `toml:"exclude"`
	PDF               metadataPDFDoc  `toml:"pdf"`
	Docx              metadataDocxDoc `toml:"docx"`
}

// metadataPDFDoc è la sezione [pdf] opzionale dentro _labnexus.toml.
// Enabled è *bool per distinguere "non setted" (nil) da "esplicitamente false".
type metadataPDFDoc struct {
	Enabled *bool `toml:"enabled"`
}

// metadataDocxDoc è la sezione [docx] opzionale dentro _labnexus.toml.
type metadataDocxDoc struct {
	Enabled *bool `toml:"enabled"`
}

// Metadata è il contenuto parsato di un `_labnexus.toml`, ritornato da
// LoadMetadata. Estensione progressiva: aggiungi campi qui invece di
// allargare la signature di LoadMetadata.
type Metadata struct {
	Profile           string
	TriggerPrompt     string
	TriggerPromptFile string
	Exclude           []string // file (rel a cartella job) da NON inviare al modello
	// PDFEnabled override per-capability del default [pdf] master.
	// nil = non setted; risolto via config.ResolvePDFEnabled.
	PDFEnabled *bool
	// DocxEnabled override per-capability del default [docx] master.
	// nil = non setted; risolto via config.ResolveDocxEnabled.
	DocxEnabled *bool
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
		m, err := LoadMetadata(metaPath)
		if err != nil {
			return nil // metadata malformato → skip
		}
		j := newJob(folderName, jobPath, m.Profile, m.TriggerPrompt, m.TriggerPromptFile, "metadata")
		j.PDFEnabled = m.PDFEnabled
		j.DocxEnabled = m.DocxEnabled
		j.Exclude = m.Exclude
		return j
	}
	// Fallback convention naming
	matches := conventionRegex.FindStringSubmatch(folderName)
	if matches == nil {
		return nil // né metadata né convention → skip
	}
	return newJob(folderName, jobPath, strings.TrimSpace(matches[1]), "", "", "convention")
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
func newJob(folderName, jobPath, profileName, triggerInline, triggerFile, source string) *Job {
	return &Job{
		Name:              folderName,
		Path:              jobPath,
		Profile:           profileName,
		TriggerPrompt:     triggerInline,
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

// WriteMetadata scrive un `_labnexus.toml` in dir (usato dal wizard `labnexus
// init`). Se promptFile != "" aggiunge `trigger_prompt_file` (override del
// prompt di default del profilo); altrimenti il job eredita il prompt del
// profilo. Path-safety su promptFile coerente con LoadMetadata.
func WriteMetadata(dir, profileName, promptFile string) error {
	if strings.TrimSpace(profileName) == "" {
		return fmt.Errorf("jobs: WriteMetadata profile vuoto")
	}
	if pf := strings.TrimSpace(promptFile); pf != "" {
		if strings.Contains(pf, "..") || filepath.IsAbs(pf) {
			return fmt.Errorf("jobs: trigger_prompt_file %q non consentito (path relativo senza '..')", pf)
		}
	}
	var b strings.Builder
	b.WriteString("# Metadata Sprint 1.5.C — auto-discovery via labnexus jobs.\n")
	b.WriteString("# Generato da `labnexus init`. Per cambiare profilo, edita la riga 'profile'.\n")
	fmt.Fprintf(&b, "profile = %q\n", profileName)
	if pf := strings.TrimSpace(promptFile); pf != "" {
		fmt.Fprintf(&b, "trigger_prompt_file = %q\n", pf)
	}
	b.WriteString("# Per NON inviare al modello certi file della cartella (note, bozze,\n")
	b.WriteString("# il LEGGIMI...), elenca i loro nomi qui sotto:\n")
	b.WriteString("# exclude = [\"LEGGIMI.txt\"]\n")
	return os.WriteFile(filepath.Join(dir, metadataFileName), []byte(b.String()), 0o644)
}

// LoadMetadata legge un singolo `_labnexus.toml` e ritorna *Metadata.
// Errore se il file non esiste o il TOML è malformato.
//
// Fix review 1.5.C CRITICAL-mistral-1: defense-in-depth lexical validation
// di TriggerPromptFile contro path traversal. Pattern coerente con
// profile.validateTrigger (Sprint 1.5.B). Runtime check resta in
// input.ParseTriggerPromptFile come second-line.
func LoadMetadata(path string) (*Metadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("jobs: read metadata %q: %w", path, err)
	}
	var doc metadataDoc
	if err := toml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("jobs: parse metadata %q: %w", path, err)
	}
	if strings.TrimSpace(doc.Profile) == "" {
		return nil, fmt.Errorf("jobs: metadata %q manca campo 'profile'", path)
	}
	hasInline := strings.TrimSpace(doc.TriggerPrompt) != ""
	hasFile := strings.TrimSpace(doc.TriggerPromptFile) != ""
	// XOR override del job: trigger_prompt e trigger_prompt_file mutuamente
	// esclusivi, simmetrico al profilo (profile.validateTrigger). Entrambi
	// assenti è valido: il job eredita il default del profilo.
	if hasInline && hasFile {
		return nil, fmt.Errorf("jobs: metadata %q ha trigger_prompt e trigger_prompt_file insieme — mutually exclusive (specifica uno solo dei due)", path)
	}
	if hasInline && len(doc.TriggerPrompt) < profile.MinTriggerLen {
		return nil, fmt.Errorf("jobs: metadata %q ha trigger_prompt troppo corto (%d caratteri), almeno %d richiesti", path, len(doc.TriggerPrompt), profile.MinTriggerLen)
	}
	if hasFile {
		if strings.Contains(doc.TriggerPromptFile, "..") {
			return nil, fmt.Errorf("jobs: metadata %q ha trigger_prompt_file %q non consentito (contiene '..')", path, doc.TriggerPromptFile)
		}
		if filepath.IsAbs(doc.TriggerPromptFile) {
			return nil, fmt.Errorf("jobs: metadata %q ha trigger_prompt_file %q non consentito (path assoluto)", path, doc.TriggerPromptFile)
		}
	}
	return &Metadata{
		Profile:           doc.Profile,
		TriggerPrompt:     doc.TriggerPrompt,
		TriggerPromptFile: doc.TriggerPromptFile,
		Exclude:           doc.Exclude,
		PDFEnabled:        doc.PDF.Enabled,
		DocxEnabled:       doc.Docx.Enabled,
	}, nil
}
