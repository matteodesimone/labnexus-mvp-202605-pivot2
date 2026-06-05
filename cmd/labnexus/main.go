// labnexus CLI — esegue una capability ispettiva del modello AICertus chiamando
// Qwen 3 via EUrouter (cloud EU-GDPR, default deliverable Sprint 1.5+) oppure
// Ollama in locale (provider opzionale per utenti con hardware adeguato).
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"encoding/json"

	"github.com/labnexus/labnexus/internal/config"
	"github.com/labnexus/labnexus/internal/input"
	"github.com/labnexus/labnexus/internal/jobs"
	"github.com/labnexus/labnexus/internal/paths"
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
		"eurouter: nessun token", // provider ha risposto vuoto/errore non-SSE → runtime
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
	// FR-2: default di --profiles-dir / --kb-dir risolti **relativi al binario**
	// (via os.Executable + DeliveryRootForBinary), NON al cwd del processo.
	// Fix bug `.pipeline/bugs/path-resolution-profili-kb-cwd-relative.md`.
	deliveryRoot := paths.DeliveryRoot()
	defaultProfiliDir := filepath.Join(deliveryRoot, "profili")
	defaultKbDir := filepath.Join(deliveryRoot, "KB-ispettore")
	defaultConfigPath := filepath.Join(deliveryRoot, "labnexus.config.toml")
	defaultLavoriDir := filepath.Join(deliveryRoot, "lavori")

	root := &cobra.Command{
		Use:   "labnexus",
		Short: "Eseguibile LabNexus — capability ispettive AICertus via Qwen su EUrouter cloud EU-GDPR (FR-1).",
		Long:  "labnexus esegue una capability ispettiva alla volta del modello AICertus chiamando Qwen 3 via EUrouter (cloud EU-GDPR, default Sprint 1.5+) o, in opzione, Ollama in locale. Senza argomenti apre una TUI sequenziale (profilo → input → output).",
		// Default action: nessun subcommand → TUI o argomento posizionale come input
		RunE:          rootRunE,
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	root.PersistentFlags().String("profiles-dir", defaultProfiliDir, "cartella con i file profilo .toml (default: <delivery-root>/profili)")
	root.PersistentFlags().String("kb-dir", defaultKbDir, "cartella della KB-ispettore di Denis (default: <delivery-root>/KB-ispettore)")
	root.PersistentFlags().String("provider", "", "override del provider del profilo: ollama|eurouter (FR-12)")
	root.PersistentFlags().String("config", defaultConfigPath, "path al labnexus.config.toml master (FR-11; default: <delivery-root>/labnexus.config.toml)")
	root.PersistentFlags().String("lavori-dir", defaultLavoriDir, "cartella con i lavori Denis auto-discovered (FR-25; default: <delivery-root>/lavori)")
	root.AddCommand(newRunCmd(), newListCmd(), newDescribeCmd(), newCheckCmd(), newValidateCmd(), newJobsCmd(), newInitCmd())
	return root
}

// resolveLavoriDir ritorna il path della cartella `lavori/`.
func resolveLavoriDir(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("lavori-dir")
	return v
}

// preCheckAPIKey verifica che la chiave EUROUTER sia disponibile prima del
// run quando il provider effettivo è eurouter (FR-36). Fail-fast: meglio
// rifiutare al pre-flight che a runtime dopo aver caricato KB+input.
// Ritorna nil se ok, *exitError code 2 se manca.
func preCheckAPIKey(cmd *cobra.Command, p *profile.Profile) error {
	if p.Provider != "eurouter" {
		return nil
	}
	if os.Getenv("EUROUTER_API_KEY") != "" {
		return nil
	}
	cfgPath := resolveConfigPath(cmd)
	if _, err := os.Stat(cfgPath); err == nil {
		if master, err := config.Load(cfgPath); err == nil && master.EurouterAPIKey != "" {
			return nil
		}
	}
	return newExit(2, "EUROUTER_API_KEY mancante. Apri 'labnexus.config.toml' e inserisci la tua chiave EUROUTER alla riga 'eurouter_api_key = \"...\"' prima di lanciare. (oppure: 'export EUROUTER_API_KEY=<chiave>')")
}

// preCheckAPIKeyByName carica il profile + master, applica merge, e chiama
// preCheckAPIKey. Helper invocato da runRunE (fix review 1.5.C MEDIUM-Claude-4).
// Se il profile non esiste o ha problemi, salta silenziosamente — l'errore
// verrà riportato dal vero load nel runner.Run successivo.
func preCheckAPIKeyByName(cmd *cobra.Command, profileName, profiliDir string) error {
	profilePath := filepath.Join(profiliDir, profileName+".toml")
	p, err := profile.Load(profilePath)
	if err != nil {
		return nil // load failure → runner.Run riporterà
	}
	// Merge col master per provider effettivo
	cfgPath := resolveConfigPath(cmd)
	if _, statErr := os.Stat(cfgPath); statErr == nil {
		if master, loadErr := config.Load(cfgPath); loadErr == nil {
			p = config.Merge(master, p)
		}
	}
	return preCheckAPIKey(cmd, p)
}

func resolveDirs(cmd *cobra.Command) (profiliDir, kbDir, providerOverride string) {
	profiliDir, _ = cmd.Flags().GetString("profiles-dir")
	kbDir, _ = cmd.Flags().GetString("kb-dir")
	providerOverride, _ = cmd.Flags().GetString("provider")
	return
}

// resolveConfigPath ritorna il path del config master TOML (FR-11).
// Risolto via flag `--config` (override) o default relativo al binario.
//
// Fix review 1.5.C HIGH-mistral-4 (--config path validation): documentazione.
// Il flag --config è uso advanced (CTO/CI). Per Denis tipico, il default
// resolve via paths.DeliveryRoot() copre tutti gli scenari. No validation
// hard sui path (Sprint 1.5 single-user trust-the-user); idea backlog
// Sprint 2 per restrict whitelist o sandbox.
func resolveConfigPath(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("config")
	return v
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
		ConfigPath:       resolveConfigPath(cmd),
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
	cmd.Flags().String("profile", "", "nome del profilo (es. revisione) — obbligatorio se --job non setted")
	cmd.Flags().String("input", "", "cartella di input — obbligatoria se --job non setted")
	cmd.Flags().String("output", "", "cartella di output — opzionale se --job (default: <input>/output/)")
	cmd.Flags().String("job", "", "nome del job auto-discovered in lavori/ (FR-29; alternativa XOR a --profile/--input/--output)")
	cmd.Flags().Bool("pdf", false, "genera anche il PDF accoppiato all'output MD. Override massimo della gerarchia: CLI > _labnexus.toml [pdf] > labnexus.config.toml [pdf] > off (default). Usa --pdf=false per disattivare esplicitamente quando il master ha pdf on.")
	cmd.Flags().Bool("docx", false, "genera anche il DOCX accoppiato all'output MD (Word-editable per L2 validation Denis). Gerarchia identica a --pdf: CLI > _labnexus.toml [docx] > labnexus.config.toml [docx] > off.")
	return cmd
}

func runRunE(cmd *cobra.Command, _ []string) error {
	profiliDir, kbDir, providerOverride := resolveDirs(cmd)
	jobName, _ := cmd.Flags().GetString("job")
	profileName, _ := cmd.Flags().GetString("profile")
	in, _ := cmd.Flags().GetString("input")
	out, _ := cmd.Flags().GetString("output")
	// Fix review 1.5.C MEDIUM-Claude-3: XOR explicit check su --job vs altri flag.
	if jobName != "" && (profileName != "" || in != "") {
		return newExit(2, "--job XOR (--profile/--input): specifica UNO solo (--output può essere comunque setted per override del default <input>/output/)")
	}
	if jobName != "" {
		j, err := resolveJobOrError(cmd, jobName)
		if err != nil {
			return err
		}
		profileName, in = j.Profile, j.InputDir
		if out == "" {
			out = j.OutputDir
		}
	}
	if profileName == "" || in == "" || out == "" {
		return newExit(2, "flag obbligatori mancanti: serve --job <name> OPPURE (--profile + --input + --output)")
	}
	// Fix review 1.5.C MEDIUM-Claude-4: preCheckAPIKey invoked binary-side
	// come safety net oltre allo shell-side in labnexus.command (FR-36).
	if err := preCheckAPIKeyByName(cmd, profileName, profiliDir); err != nil {
		return err
	}
	// Tristate detection del flag --pdf: solo se l'utente l'ha esplicitamente
	// passato consideriamo il valore (override CLI); altrimenti nil = "non setted"
	// e la risoluzione cade su job-level → master-level → default(false).
	var pdfCLI *bool
	if cmd.Flags().Changed("pdf") {
		v, _ := cmd.Flags().GetBool("pdf")
		pdfCLI = &v
	}
	var pdfJob, docxJob *bool
	var triggerJobInline, triggerJobFile string
	var excludeJob []string
	if jobName != "" {
		if j, err := resolveJobOrError(cmd, jobName); err == nil {
			pdfJob = j.PDFEnabled
			docxJob = j.DocxEnabled
			triggerJobInline = j.TriggerPrompt
			triggerJobFile = j.TriggerPromptFile
			excludeJob = j.Exclude
		}
	}
	var docxCLI *bool
	if cmd.Flags().Changed("docx") {
		v, _ := cmd.Flags().GetBool("docx")
		docxCLI = &v
	}
	res, err := runner.Run(runner.Config{
		ProfileName:          profileName,
		InputDir:             in,
		OutputDir:            out,
		ProviderOverride:     providerOverride,
		ProfiliDir:           profiliDir,
		KbDir:                kbDir,
		ConfigPath:           resolveConfigPath(cmd),
		PDFEnabledCLI:        pdfCLI,
		PDFEnabledJob:        pdfJob,
		DocxEnabledCLI:       docxCLI,
		DocxEnabledJob:       docxJob,
		Capability:           jobName,
		TriggerPromptJob:     triggerJobInline,
		TriggerPromptFileJob: triggerJobFile,
		ExcludeInput:         excludeJob,
	})
	if err != nil {
		return newExit(classifyError(err), "%v", err)
	}
	if res.ExitCode != 0 {
		return newExit(res.ExitCode, "esecuzione conclusa con stato non OK")
	}
	return nil
}

// resolveJobOrError fa discovery e resolve di un job per nome. Errore exit 2
// se nessun match.
func resolveJobOrError(cmd *cobra.Command, name string) (*jobs.Job, error) {
	lavoriDir := resolveLavoriDir(cmd)
	js, err := jobs.Discover(lavoriDir)
	if err != nil {
		return nil, newExit(2, "auto-discovery lavori/ fallito: %v", err)
	}
	j := jobs.Resolve(js, name)
	if j == nil {
		return nil, newExit(2, "job %q non trovato in %s (usa 'labnexus jobs' per la lista)", name, lavoriDir)
	}
	return j, nil
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
			path := filepath.Join(profiliDir, args[0]+".toml")
			p, err := profile.Load(path)
			if err != nil {
				return newExit(2, "profilo non trovato: %v", err)
			}
			merged := mergeWithMasterIfAvailable(cmd, p)
			fmt.Printf("profilo:        %s\n", merged.Profilo)
			fmt.Printf("descrizione:    %s\n", merged.Descrizione)
			fmt.Printf("provider:       %s\n", merged.Provider)
			fmt.Printf("modello:        %s\n", merged.Modello)
			fmt.Printf("temperature:    %v\n", merged.Temperature)
			fmt.Printf("max_tokens:     %d\n", merged.MaxTokens)
			fmt.Printf("context_window: %d\n", merged.ContextWindow)
			fmt.Println("kb_files:")
			for _, f := range merged.KbFiles {
				fmt.Printf("  - %s\n", f)
			}
			if merged.TriggerPromptFile != "" {
				fmt.Printf("trigger_prompt_file: %s\n", merged.TriggerPromptFile)
			} else {
				preview := merged.TriggerPrompt
				if len(preview) > 100 {
					preview = preview[:100] + "..."
				}
				preview = strings.ReplaceAll(strings.TrimSpace(preview), "\n", " ")
				fmt.Printf("trigger_prompt:    %s\n", preview)
			}
			return nil
		},
	}
}

// mergeWithMasterIfAvailable carica il master config (se path esiste) e
// fa merge col profile. Se master non disponibile, ritorna profile invariato.
// Usato da describe e validate per mostrare/validare il profilo EFFETTIVO
// (post-merge), non solo il file isolato.
//
// Fallback: se il path risolto dal flag --config non esiste, prova anche
// `./labnexus.config.toml` nella CWD (utile per esecuzione da repo dev
// o se Denis sposta il binario in una cartella diversa dal config).
func mergeWithMasterIfAvailable(cmd *cobra.Command, p *profile.Profile) *profile.Profile {
	cfgPath := resolveConfigPath(cmd)
	if _, err := os.Stat(cfgPath); err != nil {
		// Fallback CWD
		cwdPath := "labnexus.config.toml"
		if _, err := os.Stat(cwdPath); err != nil {
			return p
		}
		cfgPath = cwdPath
	}
	master, err := config.Load(cfgPath)
	if err != nil {
		return p
	}
	return config.Merge(master, p)
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
				ConfigPath:       resolveConfigPath(cmd),
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
	cmd.Flags().Bool("show-prompt", false, "stampa il prompt composto su stdout (⚠ contiene contenuti dei kb_files + input, potenziali PII; in non-TTY emette warning su stderr)")
	_ = cmd.MarkFlagRequired("input")
	return cmd
}

// --- init [cartella]: wizard per creare _labnexus.toml ---

// initExitNeedsEdit (exit code 10) segnala che è stato creato un file-prompt da
// editare PRIMA di eseguire: l'Esegui.command-wizard lo intercetta e si ferma.
const initExitNeedsEdit = 10

// initPromptTemplate è lo starter NON vuoto scritto quando l'utente sceglie di
// creare un nuovo prompt: evita un trigger vuoto e guida la scrittura. ASCII
// (viene avvolto in RTF da input.WrapPlainAsRTF, vedi initPromptFileName).
const initPromptTemplate = `TASK: <descrivi qui il compito di questa capability>

INPUT FORNITI:
- <elenca i file che metti nella cartella>

ISTRUZIONI:
- <cosa deve fare il modello, passo per passo>

OUTPUT ATTESO (markdown):
## 1. <prima sezione>
## 2. <seconda sezione>

VINCOLI:
- Non inventare dati non presenti negli input.
- Cita sempre i riferimenti puntuali (norma, documento, codice).
`

// initPromptFileName è il file-prompt creato dal wizard: .rtf cosi' Denis lo
// edita in Word/Pages (rich text), coerente col pattern Prompt_INPUT dei lavori.
const initPromptFileName = "Prompt_INPUT.rtf"

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [cartella]",
		Short: "crea interattivamente il _labnexus.toml di una cartella di lavoro (scegli profilo + prompt)",
		Long:  "init apre un wizard: scegli la capability (profilo) e se usare il prompt predefinito del profilo o crearne uno nuovo da file. Scrive il _labnexus.toml nella cartella indicata (default: cartella corrente). Pensato per il flusso 'cartella template + doppio click'.",
		Args:  cobra.MaximumNArgs(1),
		// Messaggi gestiti internamente; evita il doppio 'Error:' di cobra (serve
		// per l'exit code 10 'prompt da editare' che non è un vero errore).
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE:          initRunE,
	}
	return cmd
}

func initRunE(cmd *cobra.Command, args []string) error {
	profiliDir, _, _ := resolveDirs(cmd)
	workDir := "."
	if len(args) == 1 {
		workDir = args[0]
	}
	absWork, err := filepath.Abs(strings.TrimSpace(workDir))
	if err != nil {
		return initFail(2, "cartella %q non risolvibile: %v", workDir, err)
	}
	if fi, statErr := os.Stat(absWork); statErr != nil || !fi.IsDir() {
		return initFail(2, "%q non è una cartella esistente", workDir)
	}
	metaPath := filepath.Join(absWork, "_labnexus.toml")
	_, statErr := os.Stat(metaPath)
	exists := statErr == nil

	choice, err := tui.RunInit(profiliDir, exists)
	if err != nil {
		if errors.Is(err, tui.ErrInitCancelled) {
			return initFail(2, "init annullato.")
		}
		if errors.Is(err, tui.ErrNonTTY) {
			return initFail(2, "labnexus init richiede un terminale interattivo: lancialo da Terminal o con doppio click su Esegui.command (cartella template).")
		}
		return initFail(2, "%v", err)
	}

	if choice.NewPromptFile {
		promptPath := filepath.Join(absWork, initPromptFileName)
		if _, e := os.Stat(promptPath); e != nil { // non sovrascrivere un prompt esistente
			rtf := input.WrapPlainAsRTF(initPromptTemplate)
			if werr := os.WriteFile(promptPath, []byte(rtf), 0o644); werr != nil {
				return initFail(1, "impossibile creare %q: %v", initPromptFileName, werr)
			}
		}
		if werr := jobs.WriteMetadata(absWork, choice.Profile, initPromptFileName); werr != nil {
			return initFail(1, "impossibile scrivere _labnexus.toml: %v", werr)
		}
		fmt.Printf("\nCreato _labnexus.toml (profilo %q) e %s.\n", choice.Profile, initPromptFileName)
		fmt.Printf("Apri %s (Word/Pages/TextEdit), scrivi le istruzioni del task, salva, poi rilancia.\n", initPromptFileName)
		return &exitError{Code: initExitNeedsEdit, Err: errors.New("prompt da editare")}
	}

	if werr := jobs.WriteMetadata(absWork, choice.Profile, ""); werr != nil {
		return initFail(1, "impossibile scrivere _labnexus.toml: %v", werr)
	}
	fmt.Printf("\nCreato _labnexus.toml (profilo %q, prompt predefinito del profilo).\n", choice.Profile)
	return nil
}

// initFail stampa il messaggio su stderr (cobra è silenziato sul comando init) e
// ritorna l'exit code.
func initFail(code int, format string, a ...interface{}) error {
	fmt.Fprintf(os.Stderr, "ERRORE: "+format+"\n", a...)
	return newExit(code, format, a...)
}

// --- validate <profile> ---

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <profile>",
		Short: "valida lo schema TOML di un profilo (post-merge master) senza eseguirlo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profiliDir, kbDir, _ := resolveDirs(cmd)
			path := filepath.Join(profiliDir, args[0]+".toml")
			p, err := profile.Load(path)
			if err != nil {
				return newExit(2, "%v", err)
			}
			merged := mergeWithMasterIfAvailable(cmd, p)
			if err := profile.Validate(merged, kbDir); err != nil {
				return newExit(2, "%v", err)
			}
			fmt.Println("schema OK")
			return nil
		},
	}
}

// --- jobs <subcommand> (FR-28) ---

func newJobsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "elenca i lavori auto-discovered in lavori/ (FR-28 concierge mode)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			lavoriDir := resolveLavoriDir(cmd)
			js, err := jobs.Discover(lavoriDir)
			if err != nil {
				return newExit(2, "auto-discovery lavori/ fallito: %v", err)
			}
			// Fix review 1.5.C HIGH-Claude-1: valida che ogni job referenzi
			// un profile esistente in profili/. Popola Job.Warning per orfani.
			profiliDir, _, _ := resolveDirs(cmd)
			jobs.ValidateProfileExists(js, profiliDir)
			jsonOut, _ := cmd.Flags().GetBool("json")
			if jsonOut {
				return emitJobsJSON(js)
			}
			emitJobsHuman(js)
			return nil
		},
	}
	cmd.Flags().Bool("json", false, "output JSON strutturato (backend-friendly)")
	return cmd
}

// emitJobsHuman stampa la lista jobs in formato human-readable (FR-28 default).
func emitJobsHuman(js []*jobs.Job) {
	if len(js) == 0 {
		fmt.Println("Nessun lavoro trovato in lavori/ — aggiungi una cartella con _labnexus.toml o usa convention naming \"... — Profilo <nome>\".")
		return
	}
	fmt.Printf("%d lavoro/i auto-discovered:\n\n", len(js))
	for _, j := range js {
		fmt.Printf("  %s\n    profile:   %s\n    input_dir: %s\n    output:    %s\n    [%s]\n",
			j.Name, j.Profile, j.InputDir, j.OutputDir, j.Source)
		if j.Warning != "" {
			fmt.Fprintf(os.Stderr, "    ⚠ WARNING: %s\n", j.Warning)
		}
		fmt.Println()
	}
}

// emitJobsJSON stampa la lista jobs come JSON (FR-28 --json flag).
func emitJobsJSON(js []*jobs.Job) error {
	type jobJSON struct {
		Name      string `json:"name"`
		Profile   string `json:"profile"`
		InputDir  string `json:"input_dir"`
		OutputDir string `json:"output_dir"`
		Source    string `json:"source"`
		Warning   string `json:"warning,omitempty"`
	}
	out := make([]jobJSON, len(js))
	for i, j := range js {
		out[i] = jobJSON{Name: j.Name, Profile: j.Profile, InputDir: j.InputDir, OutputDir: j.OutputDir, Source: j.Source, Warning: j.Warning}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
