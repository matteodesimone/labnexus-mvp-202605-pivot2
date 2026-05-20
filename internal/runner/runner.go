// Package runner orchestra l'esecuzione di una capability (Pipeline stages).
//
//	load profile → load KB → parse input → compose prompt → estimate tokens → stream LLM → write output
package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labnexus/labnexus/internal/input"
	"github.com/labnexus/labnexus/internal/output"
	"github.com/labnexus/labnexus/internal/profile"
	"github.com/labnexus/labnexus/internal/prompt"
	"github.com/labnexus/labnexus/internal/provider"
	"github.com/labnexus/labnexus/internal/runlog"
	"github.com/labnexus/labnexus/internal/tokens"
)

// Config sono i parametri del run.
type Config struct {
	ProfileName      string
	InputDir         string
	OutputDir        string
	ProviderOverride string
	ProfiliDir       string
	KbDir            string
	DryRun           bool
	ShowPrompt       bool
}

// Result è l'esito dell'esecuzione.
type Result struct {
	OutputPath  string
	ExitCode    int
	TokensUsed  int
	DurationSec float64
}

// Run esegue la pipeline end-to-end.
func Run(cfg Config) (*Result, error) {
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}
	log := runlog.New(os.Stderr)

	p, err := loadAndValidateProfile(cfg, log)
	if err != nil {
		return nil, err
	}
	kbTexts, err := loadKBTexts(cfg, p, log)
	if err != nil {
		return nil, err
	}
	parsed, err := parseInputDir(cfg, log)
	if err != nil {
		return nil, err
	}
	composed, err := composePromptFromInputs(p, kbTexts, parsed, log)
	if err != nil {
		return nil, err
	}
	check, err := checkTokensAgainstContext(p, composed, log)
	if err != nil {
		return nil, err
	}
	if cfg.ShowPrompt {
		printPrompt(composed)
	}
	if cfg.DryRun {
		log.Info("dry-run: nessuna chiamata LLM eseguita (%d token stimati)", check.Tokens)
		return &Result{ExitCode: 0, TokensUsed: check.Tokens}, nil
	}
	return streamAndWriteOutput(cfg, p, composed, check, parsed, log)
}

// validateConfig riempie i default e verifica i campi obbligatori.
func validateConfig(cfg *Config) error {
	if cfg.ProfileName == "" {
		return fmt.Errorf("runner: --profile è obbligatorio")
	}
	if cfg.ProfiliDir == "" {
		cfg.ProfiliDir = "./profili"
	}
	if cfg.KbDir == "" {
		cfg.KbDir = "./KB-ispettore"
	}
	return nil
}

// loadAndValidateProfile carica e valida il profilo da disco.
func loadAndValidateProfile(cfg Config, log *runlog.Logger) (*profile.Profile, error) {
	log.BeginStep("caricamento profilo")
	defer log.EndStep()
	p, err := loadProfile(cfg.ProfiliDir, cfg.ProfileName)
	if err != nil {
		return nil, err
	}
	if err := profile.Validate(p, cfg.KbDir); err != nil {
		return nil, fmt.Errorf("runner: validate profilo: %w", err)
	}
	return p, nil
}

// loadKBTexts legge dal filesystem i file della KB-ispettore citati nel profilo.
func loadKBTexts(cfg Config, p *profile.Profile, log *runlog.Logger) ([]string, error) {
	log.BeginStep("caricamento KB")
	defer log.EndStep()
	return loadKB(cfg.KbDir, p.KbFiles)
}

// parseInputDir esegue il walk non-ricorsivo e parsa i file della cartella di input.
func parseInputDir(cfg Config, log *runlog.Logger) (*input.ParsedResult, error) {
	log.BeginStep("parsing input")
	defer log.EndStep()
	parsed, err := input.ParseDir(cfg.InputDir)
	if err != nil {
		return nil, fmt.Errorf("runner: parse input: %w", err)
	}
	for _, sk := range parsed.Skipped {
		log.Warn("%s: %s (suggerimento: converti manualmente con pandoc se è un PDF complesso)", sk.Name, sk.Reason)
	}
	return parsed, nil
}

// composePromptFromInputs costruisce system+user message per il provider.
func composePromptFromInputs(p *profile.Profile, kbTexts []string, parsed *input.ParsedResult, log *runlog.Logger) (*prompt.Composed, error) {
	log.BeginStep("composizione prompt")
	defer log.EndStep()
	inputMap := make(map[string]string, len(parsed.Files))
	for _, f := range parsed.Files {
		inputMap[f.Name] = f.Text
	}
	return prompt.Compose(p.TriggerPrompt, kbTexts, inputMap)
}

// checkTokensAgainstContext valuta la stima token vs context_window (FR-6).
func checkTokensAgainstContext(p *profile.Profile, composed *prompt.Composed, log *runlog.Logger) (*tokens.CheckResult, error) {
	log.BeginStep("stima context")
	defer log.EndStep()
	ctxWin := p.ContextWindow
	if ctxWin <= 0 {
		ctxWin = 128000
	}
	check, err := tokens.Check(len(composed.System)+len(composed.User), ctxWin)
	if err != nil {
		return nil, err
	}
	if check.Block {
		return nil, fmt.Errorf("runner: context window superato: %d token, limite %d (%.0f%%) — riduci i file di input o aumenta context_window nel profilo", check.Tokens, ctxWin, check.Ratio*100)
	}
	if check.Warning {
		log.Warn("context window al %.0f%% (%d token su %d)", check.Ratio*100, check.Tokens, ctxWin)
	}
	return check, nil
}

// printPrompt stampa il prompt composto su stdout (FR-5, flag --show-prompt).
func printPrompt(composed *prompt.Composed) {
	fmt.Println("=== SYSTEM MESSAGE ===")
	fmt.Println(composed.System)
	fmt.Println("\n=== USER MESSAGE ===")
	fmt.Println(composed.User)
}

// streamAndWriteOutput chiama il provider, drena lo stream e scrive il file di output.
func streamAndWriteOutput(cfg Config, p *profile.Profile, composed *prompt.Composed, check *tokens.CheckResult, parsed *input.ParsedResult, log *runlog.Logger) (*Result, error) {
	prov, err := provider.Select(cfg.ProviderOverride, os.Getenv("LABNEXUS_PROVIDER"), p.Provider)
	if err != nil {
		return nil, err
	}
	log.BeginStep("chiamata provider " + prov.Name())
	started := time.Now()
	body, stato, err := callProvider(prov, composed, p, log)
	duration := time.Since(started)
	log.EndStep()
	if err != nil {
		return nil, err
	}
	outPath, err := writeOutput(cfg, p, prov, parsed, check, body, stato, duration, started, log)
	if err != nil {
		return nil, err
	}
	return &Result{
		OutputPath:  outPath,
		ExitCode:    exitCodeFor(stato),
		TokensUsed:  check.Tokens,
		DurationSec: duration.Seconds(),
	}, nil
}

func callProvider(prov provider.LLMProvider, composed *prompt.Composed, p *profile.Profile, log *runlog.Logger) (string, string, error) {
	ch, err := prov.Stream(context.Background(), composed.System, composed.User, provider.Options{
		Modello:       p.Modello,
		Temperature:   p.Temperature,
		MaxTokens:     p.MaxTokens,
		ContextWindow: p.ContextWindow,
	})
	if err != nil {
		return "", "", err
	}
	body, stato := drainStream(ch, log)
	return body, stato, nil
}

func writeOutput(cfg Config, p *profile.Profile, prov provider.LLMProvider, parsed *input.ParsedResult, check *tokens.CheckResult, body, stato string, duration time.Duration, started time.Time, log *runlog.Logger) (string, error) {
	log.BeginStep("scrittura output")
	defer log.EndStep()
	fileNames := make([]string, len(parsed.Files))
	for i, f := range parsed.Files {
		fileNames[i] = f.Name
	}
	fm := &output.Frontmatter{
		Profilo:        p.Profilo,
		Modello:        p.Modello,
		Provider:       prov.Name(),
		DataEsecuzione: started.Format(time.RFC3339),
		DurataSecondi:  duration.Seconds(),
		TokenStimati:   check.Tokens,
		FileInput:      fileNames,
		Stato:          stato,
	}
	outPath, err := output.Write(cfg.OutputDir, fm, body)
	if err != nil {
		return "", err
	}
	log.Info("output: %s", outPath)
	return outPath, nil
}

func exitCodeFor(stato string) int {
	if stato == "completato" {
		return 0
	}
	return 1
}

func loadProfile(profiliDir, name string) (*profile.Profile, error) {
	for _, ext := range []string{".yml", ".yaml"} {
		path := filepath.Join(profiliDir, name+ext)
		if _, err := os.Stat(path); err == nil {
			return profile.Load(path)
		}
	}
	return nil, fmt.Errorf("runner: profilo non trovato: %s/%s.{yml,yaml}", profiliDir, name)
}

func loadKB(kbDir string, files []string) ([]string, error) {
	out := make([]string, 0, len(files))
	for _, rel := range files {
		data, err := os.ReadFile(filepath.Join(kbDir, rel))
		if err != nil {
			return nil, fmt.Errorf("runner: kb_file %q: %w", rel, err)
		}
		out = append(out, string(data))
	}
	return out, nil
}

// drainStream consuma il canale di StreamEvent del provider mostrando progress
// all'utente (FR-8 + UX-Sprint1).
//   - Ticker 10s: mentre non arriva il primo token, stampa "ancora in attesa..."
//     così l'utente sa che il warmup del modello (qwen3.6 36B può impiegare 1-2 min
//     al primo caricamento in RAM) è in corso e non è un freeze.
//   - Ticker 500ms (TTY only): aggiorna in-place "> X token, Ys elapsed (Z tok/s)".
//   - Ticker 30s (non-TTY only): stampa una linea di stato periodica.
//   - Al primo token: stampa il TTFT (time-to-first-token) e segna l'inizio della generazione.
//   - A fine stream: log.StreamEnd con statistica finale.
func drainStream(ch <-chan provider.StreamEvent, log *runlog.Logger) (string, string) {
	var b strings.Builder
	// Default "interrotto": diventa "completato" solo se riceviamo esplicitamente done:true (EC-8).
	stato := "interrotto"
	started := time.Now()
	var tokenCount int
	var firstToken bool

	log.Info("attendo risposta dal modello (warmup può richiedere minuti su modelli grandi)...")

	waitTicker := time.NewTicker(10 * time.Second)
	defer waitTicker.Stop()
	progressTicker := time.NewTicker(500 * time.Millisecond)
	defer progressTicker.Stop()
	nonTTYTicker := time.NewTicker(30 * time.Second)
	defer nonTTYTicker.Stop()

	for {
		select {
		case <-waitTicker.C:
			if !firstToken {
				log.Info("...ancora in attesa del primo token (%s elapsed)",
					time.Since(started).Round(time.Second))
			}
		case <-progressTicker.C:
			if firstToken {
				log.StreamProgress(tokenCount, time.Since(started))
			}
		case <-nonTTYTicker.C:
			if firstToken && !log.IsTTY {
				elapsed := time.Since(started)
				rate := float64(tokenCount) / elapsed.Seconds()
				log.Info("streaming in corso: %d token, %s elapsed (%.1f tok/s)",
					tokenCount, elapsed.Round(time.Second), rate)
			}
		case ev, ok := <-ch:
			if !ok {
				log.StreamEnd(tokenCount, time.Since(started))
				return b.String(), stato
			}
			if ev.Err != nil {
				log.StreamEnd(tokenCount, time.Since(started))
				log.Warn("stream error: %v", ev.Err)
				return b.String(), stato
			}
			if ev.Token != "" {
				if !firstToken {
					firstToken = true
					log.Info("primo token ricevuto (TTFT %s), generazione in corso...",
						time.Since(started).Round(time.Millisecond))
				}
				b.WriteString(ev.Token)
				tokenCount++
			}
			if ev.Done {
				stato = "completato"
				if ev.NoDoneMarker {
					log.Warn("stream chiuso senza done marker dal server (output verosimilmente completo, ma il modello non ha emesso il chunk finale done:true — vedi bug eof-post-content-falsamente-interrotto)")
				}
				log.StreamEnd(tokenCount, time.Since(started))
				return b.String(), stato
			}
		}
	}
}
