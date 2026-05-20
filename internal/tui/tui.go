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

// cleanPath rimuove spazi, single quote, double quote da inizio e fine.
// Il drag&drop di una cartella dal Finder al campo TUI tipicamente incolla
// il path circondato da single quotes E con uno spazio finale (eredità del
// comportamento Terminal). Senza pulizia, validateExistingDir fallisce con
// "la cartella '/Users/foo/bar ' non esiste".
func cleanPath(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, `'"`)
	s = strings.TrimSpace(s) // eventuali spazi tra quote rimaste
	return s
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
