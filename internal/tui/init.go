package tui

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
)

// ErrInitCancelled è ritornato se l'utente annulla il wizard `init` (es. rifiuta
// la sovrascrittura di un _labnexus.toml esistente).
var ErrInitCancelled = errors.New("init annullato dall'utente")

// InitChoice è l'esito del wizard `labnexus init`.
type InitChoice struct {
	Profile       string // nome profilo scelto
	NewPromptFile bool   // true = crea un prompt custom da file; false = usa il default del profilo
}

// RunInit apre il wizard interattivo per creare un _labnexus.toml in una cartella
// di lavoro: scelta del profilo + scelta della sorgente del prompt (default del
// profilo oppure nuovo file editabile). alreadyExists=true chiede conferma prima
// di sovrascrivere un metadata esistente. Errore ErrNonTTY in contesto non-TTY.
func RunInit(profiliDir string, alreadyExists bool) (*InitChoice, error) {
	if !isInteractiveTTY() {
		return nil, ErrNonTTY
	}
	if alreadyExists {
		confirm := false
		f := huh.NewForm(huh.NewGroup(
			huh.NewConfirm().
				Title("Esiste già un _labnexus.toml in questa cartella. Sovrascriverlo?").
				Affirmative("Sì, sovrascrivi").
				Negative("No, annulla").
				Value(&confirm),
		))
		if err := f.Run(); err != nil {
			return nil, fmt.Errorf("tui: %w", err)
		}
		if !confirm {
			return nil, ErrInitCancelled
		}
	}
	options, err := buildProfileOptions(profiliDir)
	if err != nil {
		return nil, err
	}
	choice := &InitChoice{}
	var promptSource string
	f := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Scegli la capability (profilo)").
			Options(options...).
			Value(&choice.Profile),
		huh.NewSelect[string]().
			Title("Prompt per questo lavoro").
			Description("Il profilo ha già un prompt predefinito. Vuoi usarlo o scriverne uno tuo?").
			Options(
				huh.NewOption("Usa il prompt predefinito del profilo", "default"),
				huh.NewOption("Crea un nuovo prompt da file (lo editi in TextEdit/Word)", "newfile"),
			).
			Value(&promptSource),
	))
	if err := f.Run(); err != nil {
		return nil, fmt.Errorf("tui: %w", err)
	}
	choice.NewPromptFile = promptSource == "newfile"
	return choice, nil
}
