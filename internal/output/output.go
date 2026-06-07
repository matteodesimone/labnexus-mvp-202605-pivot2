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
//
// ProfileDefaults sono chiavi/valori dichiarati dal profilo via
// `output.frontmatter_default`. Vengono mergeati nel frontmatter renderizzato
// dopo le chiavi engine-generated: in caso di collisione **vince l'engine**
// (l'engine sa il valore effettivo runtime — modello, provider, data, durata).
// Vedi bugfix `.pipeline/bugs/frontmatter-default-non-applicato.md` (todo #006).
//
// # Frontmatter Key Ownership Table (FR-24, Sprint 1.5.B)
//
// Le chiavi del frontmatter MD output appartengono a uno di 3 owner.
// Convenzione critica per evitare collisioni come `stato` (smascherata
// Fetta 2 review):
//
//  ENGINE-OWNED (popolato a runtime dal runner, valori autoritativi):
//    - profilo            (string)  — nome del profile usato
//    - modello            (string)  — modello LLM (post-merge col master)
//    - provider           (string)  — provider effettivo (ollama|eurouter)
//    - data_esecuzione    (RFC3339) — timestamp di inizio run
//    - durata_secondi     (float64) — wall-clock secondi
//    - token_stimati      (int)     — stima char/4 del prompt composto
//    - file_input         ([]string) — file processati nell'input dir
//    - stato              (string)  — "completato" | "interrotto"
//    - log_file           (string, optional) — path relativo al .log accoppiato
//                                              (Sprint 1.5.C, FR-34, NFR-11 audit trail ISO 17025)
//    - valutazione_denis  (string, optional) — popolato dal runner solo se non
//                                              presente già nei ProfileDefaults
//
//  PROFILE-DEFAULT-OWNED (popolato dal profilo via output.frontmatter_default,
//                         merged DOPO i campi engine; engine vince in collisione):
//    - tipo               — convenzione output del profilo (es. "bozza_revisione",
//                           "capa_pack", "audit_checklist", ecc.)
//    - stato_qm           — convenzione SGQ (es. "bozza_da_validare_qm")
//    - profilo_labnexus   — duplicato del campo engine, mantenuto per backward
//                           compat con script di Denis
//    - (altri campi custom del profilo specifico)
//
//  EXTERNAL-OWNED (popolato da Denis manualmente in fase L2 validation post-run):
//    - valutazione_denis  — "validata" | "validata_con_riserva" | "non_validata"
//                           (Denis edita il frontmatter del file output dopo aver
//                           rivisto il contenuto; non viene MAI sovrascritto dal
//                           runner se popolato dal profilo o dall'umano)
//    - (note manuali, esiti review, ecc.)
//
// Regola di precedenza in caso di collisione: ENGINE > PROFILE-DEFAULT > EXTERNAL.
// L'EXTERNAL viene preservato solo se non è popolato da nessuno dei due.
type Frontmatter struct {
	Profilo          string            `yaml:"profilo"`
	Modello          string            `yaml:"modello"`
	Provider         string            `yaml:"provider"`
	DataEsecuzione   string            `yaml:"data_esecuzione"`
	DurataSecondi    float64           `yaml:"durata_secondi"`
	TokenStimati     int               `yaml:"token_stimati"`
	FileInput        []string          `yaml:"file_input"`
	Stato            string            `yaml:"stato,omitempty"`
	ValutazioneDenis string            `yaml:"valutazione_denis,omitempty"`
	LogFile          string            `yaml:"log_file,omitempty"`
	ProfileDefaults  map[string]string `yaml:"-"`
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
	base := fmt.Sprintf("%s_%s", ts, Slug(fm.Profilo))
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
	// Bug #006 fix: merge dei default del profilo. Per ogni chiave in
	// ProfileDefaults, se NON è già presente nelle chiavi top-level del
	// frontmatter engine-generated, aggiunge in coda "key: value\n". Le
	// chiavi engine vincono in caso di collisione (l'engine conosce i
	// valori effettivi runtime).
	if len(fm.ProfileDefaults) > 0 {
		header, err = mergeFrontmatterDefaults(header, fm.ProfileDefaults)
		if err != nil {
			return "", fmt.Errorf("output: merge frontmatter defaults: %w", err)
		}
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

// reTopLevelKey matcha le chiavi top-level di un blocco YAML: linea che
// inizia con una lettera/underscore (non spazio, non #, non -) seguita da
// ":". Usato per scoprire quali chiavi sono già presenti nell'header
// engine-generated.
var reTopLevelKey = regexp.MustCompile(`(?m)^([A-Za-z_][A-Za-z0-9_]*)\s*:`)

// mergeFrontmatterDefaults aggiunge le chiavi di `defaults` in coda al
// blocco YAML `header`, ma SOLO se non già presenti come chiave top-level.
// L'ordine di aggiunta è stabile (sort delle chiavi defaults).
//
// Ritorna nuovo header bytes; non modifica l'input.
func mergeFrontmatterDefaults(header []byte, defaults map[string]string) ([]byte, error) {
	existing := map[string]bool{}
	for _, m := range reTopLevelKey.FindAllSubmatch(header, -1) {
		existing[string(m[1])] = true
	}
	missing := make([]string, 0, len(defaults))
	for k := range defaults {
		if !existing[k] {
			missing = append(missing, k)
		}
	}
	if len(missing) == 0 {
		return header, nil
	}
	// Sort per stabilità output.
	sortStringSlice(missing)
	out := make([]byte, 0, len(header)+len(missing)*40)
	out = append(out, header...)
	if len(out) > 0 && out[len(out)-1] != '\n' {
		out = append(out, '\n')
	}
	for _, k := range missing {
		out = append(out, []byte(fmt.Sprintf("%s: %s\n", k, defaults[k]))...)
	}
	return out, nil
}

func sortStringSlice(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j-1] > s[j]; j-- {
			s[j-1], s[j] = s[j], s[j-1]
		}
	}
}

var nonSlugRe = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

// Slug normalizza s in un componente di filename sicuro: ogni sequenza di
// caratteri non [a-zA-Z0-9_-] diventa "-". Esportata perché usata sia qui (nomi
// .md/.pdf/.docx) sia da runner.outputBaseName per il `.log` accoppiato, così i
// quattro file condividono lo stesso nome base (no spazi/punti/accenti/virgole).
func Slug(s string) string {
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
	return Slug(first)
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
