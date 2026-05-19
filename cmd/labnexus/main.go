// labnexus CLI — esegue una capability ispettiva del modello AICertus chiamando
// Qwen 3 in locale via Ollama, oppure EUrouter come canale di sviluppo.
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/labnexus/labnexus/internal/profile"
	"github.com/labnexus/labnexus/internal/runner"
	"github.com/labnexus/labnexus/internal/tui"
	"github.com/spf13/cobra"
)

// classifyError applica NFR-5 (exit code semantics):
//   - exit 1: errori di rete/provider/runtime (Ollama unreachable, streaming interrotto, HTTP)
//   - exit 2: misuse (flag mancante, profilo non valido, context_window superato, API key)
func classifyError(err error) int {
	if err == nil {
		return 0
	}
	msg := err.Error()
	provider1 := []string{
		"Ollama non raggiungibile",
		"ollama: HTTP ", // server raggiunto ma stato non-200 (es. 404, 500) → runtime, non misuse
		"EUrouter: HTTP",
		"eurouter: HTTP ",
		"streaming Ollama corrotto",
		"eurouter: timeout",
		"eurouter: scanner",
		"interrotto",
		"oltre il 50%", // EC-1: data-quality error → runtime, non misuse
	}
	for _, m := range provider1 {
		if strings.Contains(msg, m) {
			return 1
		}
	}
	return 2
}

func main() {
	if err := buildRoot().Execute(); err != nil {
		// I subcommand stampano già su stderr; uniformiamo l'exit code.
		var ee *exitError
		if errors.As(err, &ee) {
			os.Exit(ee.Code)
		}
		// Errori cobra non-typed (es. "required flag(s) ... not set", "unknown command")
		// sono misuse → exit 2 (NFR-5).
		os.Exit(2)
	}
}

// exitError consente ai subcommand di propagare un exit code specifico.
type exitError struct {
	Code int
	Err  error
}

func (e *exitError) Error() string { return e.Err.Error() }

func newExit(code int, format string, args ...interface{}) *exitError {
	return &exitError{Code: code, Err: fmt.Errorf(format, args...)}
}

func buildRoot() *cobra.Command {
	root := &cobra.Command{
		Use:   "labnexus",
		Short: "Eseguibile LabNexus — capability ispettive AICertus via Qwen locale (FR-1).",
		Long:  "labnexus esegue una capability ispettiva alla volta del modello AICertus chiamando Qwen 3 in locale via Ollama (o EUrouter come canale di sviluppo). Senza argomenti apre una TUI sequenziale (profilo → input → output).",
		// Default action: nessun subcommand → TUI o argomento posizionale come input
		RunE: rootRunE,
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	root.PersistentFlags().String("profiles-dir", "./profili", "cartella con i file profilo .yml")
	root.PersistentFlags().String("kb-dir", "./KB-ispettore", "cartella della KB-ispettore di Denis")
	root.PersistentFlags().String("provider", "", "override del provider del profilo: ollama|eurouter (FR-12)")
	root.AddCommand(newRunCmd(), newListCmd(), newDescribeCmd(), newCheckCmd(), newValidateCmd())
	return root
}

func resolveDirs(cmd *cobra.Command) (profiliDir, kbDir, providerOverride string) {
	profiliDir, _ = cmd.Flags().GetString("profiles-dir")
	kbDir, _ = cmd.Flags().GetString("kb-dir")
	providerOverride, _ = cmd.Flags().GetString("provider")
	return
}

// --- root: TUI o input posizionale ---

func rootRunE(cmd *cobra.Command, args []string) error {
	profiliDir, kbDir, providerOverride := resolveDirs(cmd)
	preInput, err := resolvePreInput(args)
	if err != nil {
		return err
	}
	sel, err := runTUI(profiliDir, preInput)
	if err != nil {
		return err
	}
	return executeRunner(runner.Config{
		ProfileName:      sel.ProfileName,
		InputDir:         sel.InputDir,
		OutputDir:        sel.OutputDir,
		ProviderOverride: providerOverride,
		ProfiliDir:       profiliDir,
		KbDir:            kbDir,
	})
}

func resolvePreInput(args []string) (string, error) {
	if len(args) == 0 {
		return "", nil
	}
	fi, err := os.Stat(args[0])
	if err != nil || !fi.IsDir() {
		return "", newExit(2, "argomento %q non è una cartella esistente — usa `labnexus run` con i flag", args[0])
	}
	return args[0], nil
}

func runTUI(profiliDir, preInput string) (*tui.Selection, error) {
	sel, err := tui.Run(profiliDir, preInput)
	if err != nil {
		if errors.Is(err, tui.ErrNonTTY) {
			return nil, newExit(2, "%v", err)
		}
		return nil, newExit(1, "%v", err)
	}
	return sel, nil
}

func executeRunner(cfg runner.Config) error {
	res, err := runner.Run(cfg)
	if err != nil {
		return newExit(classifyError(err), "%v", err)
	}
	if res.ExitCode != 0 {
		return newExit(res.ExitCode, "esecuzione conclusa con stato non OK")
	}
	return nil
}

// --- run ---

func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "esegue una capability su una cartella di input",
		RunE:  runRunE,
	}
	cmd.Flags().String("profile", "", "nome del profilo (es. revisione) — obbligatorio")
	cmd.Flags().String("input", "", "cartella di input — obbligatoria")
	cmd.Flags().String("output", "", "cartella di output — obbligatoria")
	_ = cmd.MarkFlagRequired("profile")
	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("output")
	return cmd
}

func runRunE(cmd *cobra.Command, _ []string) error {
	profiliDir, kbDir, providerOverride := resolveDirs(cmd)
	profileName, _ := cmd.Flags().GetString("profile")
	in, _ := cmd.Flags().GetString("input")
	out, _ := cmd.Flags().GetString("output")
	res, err := runner.Run(runner.Config{
		ProfileName:      profileName,
		InputDir:         in,
		OutputDir:        out,
		ProviderOverride: providerOverride,
		ProfiliDir:       profiliDir,
		KbDir:            kbDir,
	})
	if err != nil {
		return newExit(classifyError(err), "%v", err)
	}
	if res.ExitCode != 0 {
		return newExit(res.ExitCode, "esecuzione conclusa con stato non OK")
	}
	return nil
}

// --- list ---

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "elenca i profili disponibili",
		RunE: func(cmd *cobra.Command, _ []string) error {
			profiliDir, _, _ := resolveDirs(cmd)
			profili, err := profile.List(profiliDir)
			if err != nil {
				return newExit(2, "%v", err)
			}
			if len(profili) == 0 {
				return newExit(2, "nessun profilo trovato in %s", profiliDir)
			}
			for _, p := range profili {
				fmt.Printf("%-20s %s\n", p.Profilo, p.Descrizione)
			}
			return nil
		},
	}
}

// --- describe <profile> ---

func newDescribeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "describe <profile>",
		Short: "mostra i dettagli di un profilo (provider, modello, kb_files, input attesi)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profiliDir, _, _ := resolveDirs(cmd)
			path := filepath.Join(profiliDir, args[0]+".yml")
			if _, err := os.Stat(path); err != nil {
				path = filepath.Join(profiliDir, args[0]+".yaml")
			}
			p, err := profile.Load(path)
			if err != nil {
				return newExit(2, "profilo non trovato: %v", err)
			}
			fmt.Printf("profilo:        %s\n", p.Profilo)
			fmt.Printf("descrizione:    %s\n", p.Descrizione)
			fmt.Printf("provider:       %s\n", p.Provider)
			fmt.Printf("modello:        %s\n", p.Modello)
			fmt.Printf("temperature:    %v\n", p.Temperature)
			fmt.Printf("max_tokens:     %d\n", p.MaxTokens)
			fmt.Printf("context_window: %d\n", p.ContextWindow)
			fmt.Println("kb_files:")
			for _, f := range p.KbFiles {
				fmt.Printf("  - %s\n", f)
			}
			return nil
		},
	}
}

// --- check <profile> --input <dir> [--show-prompt] ---

func newCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check <profile>",
		Short: "dry-run: valida profilo + input + stima token, senza chiamare LLM",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profiliDir, kbDir, providerOverride := resolveDirs(cmd)
			in, _ := cmd.Flags().GetString("input")
			showPrompt, _ := cmd.Flags().GetBool("show-prompt")
			_, err := runner.Run(runner.Config{
				ProfileName:      args[0],
				InputDir:         in,
				OutputDir:        "",
				ProviderOverride: providerOverride,
				ProfiliDir:       profiliDir,
				KbDir:            kbDir,
				DryRun:           true,
				ShowPrompt:       showPrompt,
			})
			if err != nil {
				return newExit(2, "%v", err)
			}
			return nil
		},
	}
	cmd.Flags().String("input", "", "cartella di input — obbligatoria")
	cmd.Flags().Bool("show-prompt", false, "stampa il prompt composto su stdout")
	_ = cmd.MarkFlagRequired("input")
	return cmd
}

// --- validate <profile> ---

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <profile>",
		Short: "valida lo schema YAML di un profilo senza eseguirlo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profiliDir, kbDir, _ := resolveDirs(cmd)
			path := filepath.Join(profiliDir, args[0]+".yml")
			if _, err := os.Stat(path); err != nil {
				path = filepath.Join(profiliDir, args[0]+".yaml")
			}
			p, err := profile.Load(path)
			if err != nil {
				return newExit(2, "%v", err)
			}
			if err := profile.Validate(p, kbDir); err != nil {
				return newExit(2, "%v", err)
			}
			fmt.Println("schema OK")
			return nil
		},
	}
}
