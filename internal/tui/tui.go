// Package tui implementa il flusso interattivo cross-platform (FR-10).
// Usa charmbracelet/huh come libreria di prompt sequenziale.
package tui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/labnexus/labnexus/internal/profile"
	"golang.org/x/term"
)

// ErrNonTTY è ritornato se la TUI è invocata in contesto non-interattivo (EC-16).
var ErrNonTTY = errors.New("modalità interattiva non disponibile in non-TTY — usare labnexus run con i flag --profile --input --output")

// Selection è l'esito del flusso TUI.
type Selection struct {
	ProfileName string
	InputDir    string
	OutputDir   string
}

// Run apre il flusso TUI sequenziale. Se preInput è non-vuoto e una cartella
// esistente, salta la schermata di selezione input.
func Run(profiliDir, preInput string) (*Selection, error) {
	if !isInteractiveTTY() {
		return nil, ErrNonTTY
	}
	options, err := buildProfileOptions(profiliDir)
	if err != nil {
		return nil, err
	}
	// cleanPath anche su preInput: drag&drop da Finder sull'icona .app può fornirlo
	// con quote/spazi (stesso pattern del campo TUI). Senza pulizia, isExistingDir(preInput)
	// fallisce e la TUI riapre la prompt di input — vanificando il drag&drop.
	preInput = cleanPath(preInput)
	sel := &Selection{InputDir: preInput}
	form := huh.NewForm(huh.NewGroup(buildFormFields(options, preInput, sel)...))
	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("tui: %w", err)
	}
	// cleanPath rimuove anche eventuali quote da drag&drop Finder
	sel.InputDir = cleanPath(sel.InputDir)
	sel.OutputDir = cleanPath(sel.OutputDir)
	return sel, nil
}

func isInteractiveTTY() bool {
	return term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stderr.Fd()))
}

func buildProfileOptions(profiliDir string) ([]huh.Option[string], error) {
	profili, err := profile.List(profiliDir)
	if err != nil {
		return nil, fmt.Errorf("tui: lettura profili da %q: %w", profiliDir, err)
	}
	if len(profili) == 0 {
		return nil, fmt.Errorf("tui: nessun profilo trovato in %s", profiliDir)
	}
	options := make([]huh.Option[string], 0, len(profili))
	for _, p := range profili {
		label := p.Profilo
		if p.Descrizione != "" {
			label = fmt.Sprintf("%s — %s", p.Profilo, p.Descrizione)
		}
		options = append(options, huh.NewOption(label, p.Profilo))
	}
	return options, nil
}

func buildFormFields(options []huh.Option[string], preInput string, sel *Selection) []huh.Field {
	fields := []huh.Field{
		huh.NewSelect[string]().
			Title("Scegli capability").
			Options(options...).
			Value(&sel.ProfileName),
	}
	if preInput == "" || !isExistingDir(preInput) {
		fields = append(fields, huh.NewInput().
			Title("Cartella di input").
			Description("Path assoluto della cartella con i file da processare.").
			Validate(validateExistingDir).
			Value(&sel.InputDir))
	}
	fields = append(fields, huh.NewInput().
		Title("Cartella di output").
		Description("Path dove verrà scritto il file markdown (verrà creata se non esiste).").
		Validate(validateNonEmpty).
		Value(&sel.OutputDir))
	return fields
}

func isExistingDir(p string) bool {
	if p == "" {
		return false
	}
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}

// cleanPath normalizza il path inserito nel campo TUI gestendo i pattern di
// drag&drop di Finder/Terminal macOS:
//   - path circondato da single/double quote + spazio finale (Finder)
//   - spazi e metacaratteri shell-escapati con backslash, es.
//     "CAPABILITY\ D\ —\ Profilo\ x" (Terminal). La TUI NON è una shell:
//     senza de-escape, os.Stat cerca una cartella coi backslash letterali e
//     fallisce con "la cartella '...\\ ...' non esiste".
func cleanPath(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, `'"`)
	s = strings.TrimSpace(s) // eventuali spazi tra quote rimaste
	s = unescapeDragPath(s)
	return s
}

// shellEscapedChars sono i caratteri che Terminal/Finder fanno precedere da '\'
// quando trascini una cartella nel campo input. Su macOS/Linux il backslash non
// è mai un separatore di path, quindi rimuovere il backslash davanti a questi
// caratteri è sicuro e ripristina il path reale.
const shellEscapedChars = " \t!\"#$&'()*,:;<=>?[]^`{|}~"

// unescapeDragPath rimuove i backslash usati come escape shell dal drag&drop.
// Itera per byte: agisce solo su '\' (ASCII) seguito da un metacarattere ASCII;
// le sequenze UTF-8 multibyte (es. em-dash) passano invariate.
func unescapeDragPath(s string) string {
	if !strings.Contains(s, `\`) {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) && strings.IndexByte(shellEscapedChars, s[i+1]) >= 0 {
			continue // salta il backslash di escape, mantieni il carattere successivo
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func validateExistingDir(s string) error {
	s = cleanPath(s)
	if !isExistingDir(s) {
		return fmt.Errorf("la cartella %q non esiste", s)
	}
	return nil
}

func validateNonEmpty(s string) error {
	s = cleanPath(s)
	if s == "" {
		return errors.New("path richiesto")
	}
	// crea la dir se non esiste (idempotent)
	abs, err := filepath.Abs(s)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return fmt.Errorf("impossibile creare %q: %w", abs, err)
	}
	return nil
}
