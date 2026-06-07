// Package runner orchestra l'esecuzione di una capability (Pipeline stages).
//
//	load profile → load KB → parse input → compose prompt → estimate tokens → stream LLM → write output
package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"

	"github.com/labnexus/labnexus/internal/config"
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
	ConfigPath       string // Sprint 1.5.B: path al labnexus.config.toml master (opzionale)
	DryRun           bool
	ShowPrompt       bool
	// PDFEnabledCLI: override esplicito CLI (--pdf / --pdf=false). nil = non setted.
	PDFEnabledCLI *bool
	// PDFEnabledJob: override esplicito per-capability da _labnexus.toml [pdf].
	// nil = non setted (eredita dal master). Risoluzione: cli > job > master > false.
	PDFEnabledJob *bool
	// DocxEnabledCLI / DocxEnabledJob: stessa semantica di PDFEnabled per la
	// generazione DOCX accoppiata all'MD (Sprint 1.5.D extension).
	DocxEnabledCLI *bool
	DocxEnabledJob *bool
	// Capability: display name mostrato nell'header del PDF (es. "CAPABILITY A — Profilo revisione").
	// Tipicamente il nome del job. Vuoto = fallback al nome del profilo.
	Capability string
	// TriggerPromptJob / TriggerPromptFileJob: override del trigger a livello job
	// (da _labnexus.toml). XOR fra loro (validato in jobs.LoadMetadata). Se uno
	// dei due è setted, vince sul default del profilo. Vuoti = usa il profilo.
	TriggerPromptJob     string
	TriggerPromptFileJob string
	// ExcludeInput: file (rel a InputDir) da NON inviare al modello, dal campo
	// `exclude` del _labnexus.toml del job. Si somma al file-trigger de-duplicato.
	ExcludeInput []string
}

// triggerOverride è l'override del trigger proveniente dal job (_labnexus.toml).
// Inline e File sono mutuamente esclusivi (XOR già validato a monte).
type triggerOverride struct {
	Inline string
	File   string
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

	p, master, err := loadAndValidateProfile(cfg, log)
	if err != nil {
		return nil, err
	}
	kbTexts, err := loadKBTexts(cfg, p, log)
	if err != nil {
		return nil, err
	}
	ov := triggerOverride{Inline: cfg.TriggerPromptJob, File: cfg.TriggerPromptFileJob}
	// Esclusioni dall'input: i file `exclude` del job + il file usato come trigger
	// (de-dup: se il trigger viene da file non va inviato anche come dato).
	exclude := append([]string{}, cfg.ExcludeInput...)
	if tf := effectiveTriggerFile(p, ov); tf != "" {
		exclude = append(exclude, tf)
	}
	parsed, err := parseInputDir(cfg, exclude, log)
	if err != nil {
		return nil, err
	}
	composed, triggerSource, err := composePromptFromInputs(p, ov, kbTexts, parsed, cfg.InputDir, log)
	if err != nil {
		return nil, err
	}
	check, err := checkTokensAgainstContext(p, composed, log)
	if err != nil {
		return nil, err
	}
	// Budget di output dinamico + guardia: se il prompt occupa quasi tutto il
	// context, resta troppo poco per un output utile → blocca con messaggio
	// chiaro invece di produrre un output troncato/vuoto o sforare il context.
	effMaxTokens, err := resolveOutputBudget(p.MaxTokens, check.Tokens, p.ContextWindow)
	if err != nil {
		return nil, err
	}
	if cfg.ShowPrompt {
		printPromptWithPIIWarning(composed, os.Stdout, os.Stderr, isStdoutTTY())
	}
	if cfg.DryRun {
		log.Info("trigger risolto da: %s", triggerSource)
		log.Info("dry-run: nessuna chiamata LLM eseguita (%d token stimati, budget output %d)", check.Tokens, effMaxTokens)
		return &Result{ExitCode: 0, TokensUsed: check.Tokens}, nil
	}
	return streamAndWriteOutput(cfg, master, p, composed, check, parsed, kbTexts, exclude, effMaxTokens, triggerSource, log)
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

// loadAndValidateProfile carica il config master (se disponibile), carica
// il profilo, fa merge master→profile, propaga EurouterAPIKey + OllamaEndpoint
// del master a env vars (fallback se non già setted), valida il merged
// (Sprint 1.5.B).
func loadAndValidateProfile(cfg Config, log *runlog.Logger) (*profile.Profile, *config.Master, error) {
	log.BeginStep("caricamento profilo")
	defer log.EndStep()
	master := loadMasterOptional(cfg.ConfigPath, log)
	applyMasterToEnv(master)
	p, err := loadProfile(cfg.ProfiliDir, cfg.ProfileName)
	if err != nil {
		return nil, nil, err
	}
	merged := config.Merge(master, p)
	if err := profile.Validate(merged, cfg.KbDir); err != nil {
		return nil, nil, fmt.Errorf("runner: validate profilo: %w", err)
	}
	// Sprint 1.5.C verbose: log effettivo provider/modello post-merge.
	log.Info("profilo %q: provider=%s modello=%s temperature=%.2f max_tokens=%d context_window=%d",
		merged.Profilo, merged.Provider, merged.Modello, merged.Temperature, merged.MaxTokens, merged.ContextWindow)
	return merged, master, nil
}

// applyMasterToEnv propaga EurouterAPIKey + OllamaEndpoint del master a env
// vars per il provider runtime (FR-11 spec compliance, fix review 1.5.B
// HIGH-config-2). Precedenza env > master: se l'env var è già setted, NON
// viene sovrascritta. Permette a Denis di editare un solo file (master) per
// configurare API key + endpoint, senza export manuale di env vars.
func applyMasterToEnv(master *config.Master) {
	if master == nil {
		return
	}
	if master.EurouterAPIKey != "" && os.Getenv("EUROUTER_API_KEY") == "" {
		_ = os.Setenv("EUROUTER_API_KEY", master.EurouterAPIKey)
	}
	if master.OllamaEndpoint != "" && os.Getenv("LABNEXUS_OLLAMA_ENDPOINT") == "" {
		_ = os.Setenv("LABNEXUS_OLLAMA_ENDPOINT", master.OllamaEndpoint)
	}
}

// loadMasterOptional carica labnexus.config.toml se esiste, altrimenti ritorna
// nil (profile-only mode). Errori di sintassi sono fatal.
// Fallback CWD: se path non esiste, prova `./labnexus.config.toml`.
func loadMasterOptional(path string, log *runlog.Logger) *config.Master {
	if path == "" {
		path = "labnexus.config.toml"
	}
	if _, err := os.Stat(path); err != nil {
		cwdPath := "labnexus.config.toml"
		if _, err := os.Stat(cwdPath); err != nil {
			return nil
		}
		path = cwdPath
	}
	master, err := config.Load(path)
	if err != nil {
		log.Warn("config master %q non valido: %v (procedo senza master)", path, err)
		return nil
	}
	return master
}

// loadKBTexts legge dal filesystem i file della KB-ispettore citati nel profilo.
func loadKBTexts(cfg Config, p *profile.Profile, log *runlog.Logger) ([]string, error) {
	log.BeginStep("caricamento KB")
	defer log.EndStep()
	log.Info("kb dir: %s (%d file da caricare)", cfg.KbDir, len(p.KbFiles))
	texts, err := loadKB(cfg.KbDir, p.KbFiles)
	if err != nil {
		return nil, err
	}
	// Sprint 1.5.C verbose: per file, mostra il path + dim caratteri.
	totalChars := 0
	for i, t := range texts {
		log.Info("  • %s (%d caratteri)", p.KbFiles[i], len(t))
		totalChars += len(t)
	}
	log.Info("KB totale: %d caratteri", totalChars)
	return texts, nil
}

// parseInputDir esegue il walk non-ricorsivo e parsa i file della cartella di input.
func parseInputDir(cfg Config, excludeRel []string, log *runlog.Logger) (*input.ParsedResult, error) {
	log.BeginStep("parsing input")
	defer log.EndStep()
	log.Info("input dir: %s", cfg.InputDir)
	for _, e := range excludeRel {
		if strings.TrimSpace(e) != "" {
			log.Info("escluso dall'input: %s", e)
		}
	}
	parsed, err := input.ParseDir(cfg.InputDir, excludeRel...)
	if err != nil {
		return nil, fmt.Errorf("runner: parse input: %w", err)
	}
	// Sprint 1.5.C verbose: lista file parsati con dim.
	totalInputChars := 0
	for _, f := range parsed.Files {
		log.Info("  • %s (%d caratteri estratti)", f.Name, f.Size)
		totalInputChars += f.Size
	}
	log.Info("input totale: %d file parsati, %d caratteri", len(parsed.Files), totalInputChars)
	for _, sk := range parsed.Skipped {
		log.Warn("%s: %s (suggerimento: converti manualmente con pandoc se è un PDF complesso)", sk.Name, sk.Reason)
	}
	return parsed, nil
}

// composePromptFromInputs costruisce system+user message per il provider.
// Sprint 1.5.B: se profile.TriggerPromptFile è setted, carica il file
// referenziato (relativo a inputDir, path-safe) e lo usa come trigger.
// Ritorna anche la descrizione della sorgente del trigger (audit ISO 17025).
func composePromptFromInputs(p *profile.Profile, ov triggerOverride, kbTexts []string, parsed *input.ParsedResult, inputDir string, log *runlog.Logger) (*prompt.Composed, string, error) {
	log.BeginStep("composizione prompt")
	defer log.EndStep()
	inputMap := make(map[string]string, len(parsed.Files))
	for _, f := range parsed.Files {
		inputMap[f.Name] = f.Text
	}
	trigger, source, err := resolveTrigger(p, ov, inputDir)
	if err != nil {
		return nil, "", err
	}
	composed, err := prompt.Compose(trigger, kbTexts, inputMap)
	if err != nil {
		return nil, "", err
	}
	return composed, source, nil
}

// resolveTrigger risolve il trigger applicando la precedenza a DUE livelli
// (design Sprint 1.5.D, deciso col CTO):
//
//	override job (testo|file da _labnexus.toml) › default profilo (inline|file)
//
// Il trigger è XOR (testo O file) a ciascun livello: l'override del job (se
// presente) sovrascrive sempre il default shipped nel profilo. Ritorna anche
// una descrizione della sorgente risolta per l'audit trail. I path file sono
// risolti via input.ParseTriggerPromptFile (FR-17 + NFR-6 path-safe). La
// validazione XOR + minLen è già garantita a monte (jobs.LoadMetadata e
// profile.Validate).
// effectiveTriggerFile ritorna il path relativo (a InputDir) del file usato come
// trigger, quando il trigger risolto proviene da un file. Stessa precedenza di
// resolveTrigger (job file › job inline › profilo file › profilo inline). Ritorna
// "" se il trigger è inline (nessun file da escludere dall'input).
func effectiveTriggerFile(p *profile.Profile, ov triggerOverride) string {
	if strings.TrimSpace(ov.File) != "" {
		return ov.File
	}
	if strings.TrimSpace(ov.Inline) != "" {
		return "" // override inline del job: trigger non da file
	}
	if p.TriggerPromptFile != "" {
		return p.TriggerPromptFile
	}
	return ""
}

func resolveTrigger(p *profile.Profile, ov triggerOverride, inputDir string) (trigger, source string, err error) {
	// Livello 1: override del job (vince sul profilo). Emptiness via TrimSpace,
	// coerente con la validazione XOR di jobs.LoadMetadata: un valore
	// whitespace-only è "assente" e fa fallback al default del profilo.
	if strings.TrimSpace(ov.File) != "" {
		t, err := input.ParseTriggerPromptFile(inputDir, ov.File)
		return t, fmt.Sprintf("override job (file: %s)", ov.File), err
	}
	if strings.TrimSpace(ov.Inline) != "" {
		return ov.Inline, "override job (inline _labnexus.toml)", nil
	}
	// Livello 2: default shipped nel profilo.
	if p.TriggerPromptFile != "" {
		t, err := input.ParseTriggerPromptFile(inputDir, p.TriggerPromptFile)
		return t, fmt.Sprintf("default profilo (file: %s)", p.TriggerPromptFile), err
	}
	return p.TriggerPrompt, "default profilo (inline)", nil
}

// checkTokensAgainstContext valuta la stima token vs context_window (FR-6).
func checkTokensAgainstContext(p *profile.Profile, composed *prompt.Composed, log *runlog.Logger) (*tokens.CheckResult, error) {
	log.BeginStep("stima context")
	defer log.EndStep()
	ctxWin := p.ContextWindow
	if ctxWin <= 0 {
		ctxWin = 128000
	}
	// Sprint 1.5.C verbose: dimensioni system + user message.
	log.Info("system message: %d caratteri | user message: %d caratteri (totale %d)",
		len(composed.System), len(composed.User), len(composed.System)+len(composed.User))
	check, err := tokens.Check(len(composed.System)+len(composed.User), ctxWin)
	if err != nil {
		return nil, err
	}
	if check.Block {
		return nil, fmt.Errorf("runner: context window superato: %d token, limite %d (%.0f%%) — riduci i file di input o aumenta context_window nel profilo", check.Tokens, ctxWin, check.Ratio*100)
	}
	log.Info("stima token: %d / context window %d (%.0f%% utilization)", check.Tokens, ctxWin, check.Ratio*100)
	if check.Warning {
		log.Warn("context window al %.0f%% — riduci input o aumenta context_window per sicurezza", check.Ratio*100)
	}
	return check, nil
}

// piiWarning è il messaggio di avviso emesso su stderr quando
// `labnexus check --show-prompt` è eseguito in modalità non-TTY (stdout
// pipato a file/altro comando). Bugfix .pipeline/bugs/show-prompt-pii-exposure.md
// — protegge contro copia/paste accidentale di output con PII degli SGQ.
const piiWarning = `⚠  --show-prompt: output contiene contenuti dei kb_files e degli input,
   inclusi potenziali PII / dati personali del SGQ (nomi tecnici, codici
   personali, riferimenti a colleghi). NON condividere su canali non
   controllati (ticket di supporto, Slack, log condivisi). Per debug
   interno solo.`

// printPromptWithPIIWarning stampa il prompt composto su `stdout`. Se
// `stdoutIsTTY` è false (l'utente sta ridirezionando l'output verso un
// file / pipe / cattura), emette PRIMA un warning su `stderr` per
// proteggere contro copia/paste accidentale di PII. Bugfix #005.
func printPromptWithPIIWarning(composed *prompt.Composed, stdout, stderr io.Writer, stdoutIsTTY bool) {
	if !stdoutIsTTY {
		fmt.Fprintln(stderr, piiWarning)
	}
	fmt.Fprintln(stdout, "=== SYSTEM MESSAGE ===")
	fmt.Fprintln(stdout, composed.System)
	fmt.Fprintln(stdout, "\n=== USER MESSAGE ===")
	fmt.Fprintln(stdout, composed.User)
}

// isStdoutTTY ritorna true se os.Stdout è un terminale interattivo.
// Usato per decidere se emettere il PII warning di `--show-prompt`.
func isStdoutTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// streamAndWriteOutput chiama il provider, drena lo stream e scrive il file di output.
// Sprint 1.5.C (NFR-11 audit trail): apre un .log accoppiato all'.md e usa
// io.MultiWriter per scrivere il logger contemporaneamente su stderr + log file.
// Body streaming live (FR-35): i token chunks vanno raw su stderr (TTY) +
// sempre al log file.
func streamAndWriteOutput(cfg Config, master *config.Master, p *profile.Profile, composed *prompt.Composed, check *tokens.CheckResult, parsed *input.ParsedResult, kbTexts, exclude []string, effMaxTokens int, triggerSource string, log *runlog.Logger) (*Result, error) {
	prov, err := provider.Select(cfg.ProviderOverride, os.Getenv("LABNEXUS_PROVIDER"), p.Provider)
	if err != nil {
		return nil, err
	}
	started := time.Now()
	baseName := outputBaseName(p, parsed, started)
	logFile, multiLog, bodyWriter := openAuditTrail(cfg.OutputDir, baseName, log)
	if logFile != nil {
		defer logFile.Close()
	}
	// Audit trail (NFR-11): dichiara nel .log quale sorgente di trigger è stata
	// risolta (override job vs default profilo), così l'output è tracciabile.
	multiLog.Info("trigger risolto da: %s", triggerSource)
	// Audit trail: registra COSA è stato inviato al modello (KB, dati, esclusioni).
	// L'.log si apre solo qui (dopo il parsing), quindi senza questo riepilogo il
	// log non proverebbe quali dati sono stati inviati — gap diagnostico (Denis
	// shakedown 2026-06-02: "dati non inviati?").
	logInputAudit(multiLog, kbTexts, parsed, exclude, check, p)
	multiLog.BeginStep("chiamata provider " + prov.Name())
	// Sprint 1.5.C verbose: stampa URL endpoint + modello + parametri per debug.
	logProviderRequest(multiLog, prov, p, effMaxTokens, check.Tokens)
	body, stato, err := callProviderWithStreaming(prov, composed, p, effMaxTokens, multiLog, bodyWriter, true)
	// Backstop di calibrazione provider-authoritative (LOOP). La prima chiamata usa
	// la stima accurata (char/3,2) e di norma ENTRA. Se un contenuto è più denso e
	// sfora, il provider ci dice quanti input token conta e impostiamo max_tokens =
	// conteggio_reale + buffer (NESSUNA nostra stima determina il risultato). È un
	// loop perché al primo errore il conteggio può essere PER DIFETTO ("at least N");
	// l'esatto arriva dopo. Le chiamate sono "quiet": l'overflow di calibrazione NON
	// è loggato come errore (lo gestiamo qui come passo calmo), così l'audit ISO resta
	// pulito. Stop quando entra, o il reale non lascia il minimo (errore chiaro), o
	// non c'è progresso, o si esauriscono i giri.
	lastRequested := effMaxTokens
	for attempt := 1; err != nil && attempt <= maxProviderRetries; attempt++ {
		realInput, ok := parseContextOverflow(err, lastRequested)
		if !ok {
			multiLog.Warn("errore dal provider (non di context length): %v", err)
			break
		}
		retryMax, fits := retryOutputBudget(realInput, p.ContextWindow)
		if !fits {
			err = fmt.Errorf("runner: input troppo grande — il provider conta %d token di prompt; oltre il limite del modello (%d) non restano i %d token minimi di output. Riduci i file di input", realInput, p.ContextWindow, minOutputBudget)
			break
		}
		if retryMax >= lastRequested {
			multiLog.Warn("calibrazione budget non converge (il provider riporta %d token di input senza progresso): mi fermo.", realInput)
			break
		}
		multiLog.Info("calibrazione budget output (giro %d/%d): il provider conta %d token di prompt → max_tokens=%d (buffer %d), procedo.",
			attempt, maxProviderRetries, realInput, retryMax, outputBuffer)
		lastRequested = retryMax
		body, stato, err = callProviderWithStreaming(prov, composed, p, retryMax, multiLog, bodyWriter, true)
	}
	duration := time.Since(started)
	multiLog.EndStep()
	if err != nil {
		return nil, err
	}
	// Sanitizzazione: alcuni modelli (es. Qwen3.5-122B sul profilo revisione
	// 2026-05-22T11:17) racchiudono l'intero output in ```yaml ... ``` rendendo
	// il MD irreso come code block. Rimuoviamo il fence wrapping se rilevato.
	if cleaned := stripFenceWrapping(body); cleaned != body {
		multiLog.Info("body sanitization: rimosso fence wrapping spurio (modello aveva racchiuso l'output in code block)")
		body = cleaned
	}
	// Risposta vuota (EC): alcuni modelli (es. kimi-k2.6 su eurouter, shakedown
	// 2026-06-02 CAPABILITY F) chiudono lo stream con done:true ma ZERO token.
	// Non è un "completato": lo segnaliamo come 'vuoto', scriviamo un .md con la
	// spiegazione (invece di un file vuoto) e NON generiamo PDF/DOCX vuoti.
	var emptyOutput bool
	body, stato, emptyOutput = finalizeOutput(body, stato, p, prov.Name())
	if emptyOutput {
		multiLog.Warn("il modello non ha restituito output (0 token / body vuoto): run marcato 'vuoto', niente PDF/DOCX. Riprova; se persiste, cambia modello o segnala.")
	}
	logRelPath := baseName + ".log"
	outPath, fm, err := writeOutput(cfg, p, prov, parsed, check, body, stato, duration, started, multiLog, logRelPath)
	if err != nil {
		return nil, err
	}
	// PDF/DOCX accoppiati solo per run con output reale (Sprint 1.5.D): gerarchia
	// CLI > job > master > false. Salta del tutto se la risposta è vuota.
	if !emptyOutput {
		var masterPDF *bool
		if master != nil {
			masterPDF = master.PDF.Enabled
		}
		pdfEnabled := config.ResolvePDFEnabled(cfg.PDFEnabledCLI, cfg.PDFEnabledJob, masterPDF)
		if pdfPath := maybeWritePDF(outPath, body, fm, pdfEnabled, cfg.Capability, multiLog); pdfPath != "" {
			multiLog.Info("pdf accoppiato: %s", pdfPath)
		}
		var masterDocx *bool
		if master != nil {
			masterDocx = master.Docx.Enabled
		}
		docxEnabled := config.ResolveDocxEnabled(cfg.DocxEnabledCLI, cfg.DocxEnabledJob, masterDocx)
		if docxPath := maybeWriteDocx(outPath, body, fm, docxEnabled, cfg.Capability, multiLog); docxPath != "" {
			multiLog.Info("docx accoppiato: %s", docxPath)
		}
	}
	return &Result{
		OutputPath:  outPath,
		ExitCode:    exitCodeFor(stato),
		TokensUsed:  check.Tokens,
		DurationSec: duration.Seconds(),
	}, nil
}

// outputBuffer è il piccolo cuscino (token) lasciato sopra il prompt: il budget di
// output è context − prompt − outputBuffer. Tenerlo piccolo massimizza l'output.
const outputBuffer = 500

// conservativePromptTokens converte la stima FR-6 char/4 (che SOTTOSTIMA del ~25%)
// nella densità REALE verificata di kimi (~3,2 char/token: 733445 char = 229169
// token su CAPABILITY A, 2026-06-06). Fattore 4/3,2 = ×5/4. Usata SOLO per il budget
// di output; la stima FR-6 char/4 resta per il warning a 70/100%.
func conservativePromptTokens(charDiv4 int) int {
	return charDiv4 * 5 / 4
}

// defaultOutputBudget è il fallback usato SOLO quando il context_window non è noto
// (raro: il config lo definisce sempre). 32768 = esempio Moonshot per kimi thinking.
const defaultOutputBudget = 32768

// minOutputBudget è il budget di output minimo sensato su un modello "thinking":
// sotto i 16000 token (raccomandazione Moonshot per kimi) il reasoning_content può
// consumare tutto il budget lasciando content vuoto.
const minOutputBudget = 16000

// maxProviderRetries limita le iterazioni di calibrazione (backstop) su context
// overflow. Loop perché al primo errore il provider può riportare un conteggio "at
// least N" PER DIFETTO; l'esatto arriva dopo — si itera finché la richiesta entra.
const maxProviderRetries = 4

// resolveOutputBudget decide il max_tokens INIZIALE da inviare.
//
// configMax <= 0 (default: max_tokens=0) → budget = context − stima_reale(char/3,2)
// − buffer, con floor a minOutputBudget. Usa la densità VERIFICATA di kimi (≈3,2
// char/tok), così la prima chiamata ENTRA al primo colpo (audit pulito, output
// quasi-massimo) sia nel gateway eurouter (≈256000, che pre-stima ~char/4) sia nel
// modello vLLM (262144, conteggio reale). NON omesso (default provider ~1024 →
// vuoto su thinking). Se un contenuto è più denso di 3,2 e sfora, il backstop in
// streamAndWriteOutput calibra sul conteggio REALE restituito dal provider.
//
// configMax > 0 → tetto esplicito.
//
// Guardia pre-volo LENIENTE: usa la stima ottimistica (char/4) per NON rifiutare
// input che in realtà entrano; il rifiuto "vero" lo dà il provider nel backstop.
func resolveOutputBudget(configMax, promptTokens, contextWindow int) (int, error) {
	if contextWindow > 0 && promptTokens+minOutputBudget > contextWindow {
		return 0, fmt.Errorf("runner: input troppo grande — anche la stima ottimistica è ~%d token di prompt; il modello (context %d) non lascia i %d token minimi per l'output. Riduci i file di input",
			promptTokens, contextWindow, minOutputBudget)
	}
	if configMax > 0 {
		return configMax, nil
	}
	if contextWindow <= 0 {
		return defaultOutputBudget, nil // context ignoto: fallback fisso sano
	}
	budget := contextWindow - conservativePromptTokens(promptTokens) - outputBuffer
	if budget < minOutputBudget {
		budget = minOutputBudget
	}
	return budget, nil
}

// parseContextOverflow riconosce un errore di context overflow del provider in
// QUALSIASI dei 3 formati osservati ed estrae il conteggio di input token del
// prompt, per il backstop di calibrazione. lastRequested = max_tokens appena
// inviato: serve a derivare il prompt dal "total" stimato dal gateway.
//
// Formati (2026-06-06):
//   - vLLM esatto: "...229169 tokens from the input messages..."        → reale (preferito)
//   - gateway 400: "Estimated total tokens (T) exceeds model context window (256000)" → prompt = T − lastRequested
//   - vLLM stima:  "...at least N input tokens... (parameter=input_tokens, value=N)"  → lower bound
func parseContextOverflow(err error, lastRequested int) (int, bool) {
	if err == nil {
		return 0, false
	}
	msg := err.Error()
	if !strings.Contains(msg, "maximum context length") && !strings.Contains(msg, "exceeds model context window") {
		return 0, false
	}
	// Formato ESATTO (preferito): "...229169 tokens from the input messages...".
	if n, ok := intBefore(msg, "tokens from the input"); ok {
		return n, true
	}
	// Gateway: "Estimated total tokens (T)" → prompt = T − max_tokens inviato.
	if t, ok := intAfter(msg, "Estimated total tokens ("); ok && t > lastRequested {
		return t - lastRequested, true
	}
	// Formato STIMA ("at least N" / value=N): lower bound; il loop converge poi
	// sul conteggio esatto.
	if n, ok := intAfter(msg, "value="); ok {
		return n, true
	}
	if n, ok := intAfter(msg, "at least "); ok {
		return n, true
	}
	return 0, false
}

// intBefore ritorna l'intero che precede marker in s (saltando spazi). (0,false)
// se marker assente o non preceduto da cifre.
func intBefore(s, marker string) (int, bool) {
	i := strings.Index(s, marker)
	if i <= 0 {
		return 0, false
	}
	end := i
	for end > 0 && s[end-1] == ' ' {
		end--
	}
	start := end
	for start > 0 && s[start-1] >= '0' && s[start-1] <= '9' {
		start--
	}
	if start == end {
		return 0, false
	}
	n, err := strconv.Atoi(s[start:end])
	if err != nil {
		return 0, false
	}
	return n, true
}

// intAfter ritorna l'intero che segue marker in s (saltando spazi). (0,false) se
// marker assente o non seguito da cifre.
func intAfter(s, marker string) (int, bool) {
	i := strings.Index(s, marker)
	if i < 0 {
		return 0, false
	}
	j := i + len(marker)
	for j < len(s) && s[j] == ' ' {
		j++
	}
	start := j
	for j < len(s) && s[j] >= '0' && s[j] <= '9' {
		j++
	}
	if j == start {
		return 0, false
	}
	n, err := strconv.Atoi(s[start:j])
	if err != nil {
		return 0, false
	}
	return n, true
}

// retryOutputBudget calcola il max_tokens dal conteggio di input token che il
// provider ha restituito, col buffer minimo (output massimo). Ritorna (budget,true)
// se resta spazio per l'output minimo; (0,false) se l'input è genuinamente troppo
// grande (→ il caller blocca con messaggio chiaro e azionabile).
func retryOutputBudget(realInput, contextWindow int) (int, bool) {
	avail := contextWindow - realInput - outputBuffer
	if avail < minOutputBudget {
		return 0, false
	}
	return avail, true
}

// formatUsage formatta la riga "usage" gestendo la peculiarità dei modelli
// "thinking" su eurouter: a volte reasoning_tokens NON è un sottoinsieme di
// completion_tokens (kimi su eurouter, 2026-06-07: reasoning 14335 > completion
// 6341). In quel caso il content è completion e il reasoning è SEPARATO, quindi il
// totale reale li somma entrambi (prompt+completion+reasoning) e NON usiamo "di
// cui" (che sarebbe contraddittorio e farebbe sottostimare il totale).
func formatUsage(u *provider.Usage) string {
	if u.ReasoningTokens > u.CompletionTokens {
		realTotal := u.PromptTokens + u.CompletionTokens + u.ReasoningTokens
		return fmt.Sprintf("usage (token reali dal provider): prompt=%d, content=%d, reasoning=%d (separato dal content), totale=%d",
			u.PromptTokens, u.CompletionTokens, u.ReasoningTokens, realTotal)
	}
	return fmt.Sprintf("usage (token reali dal provider): prompt=%d, completion=%d (di cui reasoning=%d), totale=%d",
		u.PromptTokens, u.CompletionTokens, u.ReasoningTokens, u.TotalTokens)
}

// logProviderRequest stampa info diagnostiche prima della chiamata al provider
// (Sprint 1.5.C verbose). Endpoint risolto via env override o default. Visibile
// anche nel .log accoppiato (audit trail NFR-11).
func logProviderRequest(log *runlog.Logger, prov provider.LLMProvider, p *profile.Profile, effMaxTokens, promptTokens int) {
	var endpoint string
	switch prov.Name() {
	case "eurouter":
		endpoint = os.Getenv("LABNEXUS_EUROUTER_ENDPOINT")
		if endpoint == "" {
			endpoint = "https://api.eurouter.ai/v1/chat/completions"
		}
	case "ollama":
		endpoint = os.Getenv("LABNEXUS_OLLAMA_ENDPOINT")
		if endpoint == "" {
			endpoint = "http://localhost:11434"
		}
	}
	maxTokensNote := fmt.Sprintf("%d", effMaxTokens)
	if p.MaxTokens <= 0 {
		maxTokensNote += " (auto: context − stima_reale − buffer)"
	}
	log.Info("provider: %s | endpoint: %s | modello: %s | temperature: %.2f | max_tokens: %s | prompt ~%d token (stima reale char/3,2; char/4=%d)",
		prov.Name(), endpoint, p.Modello, p.Temperature, maxTokensNote, conservativePromptTokens(promptTokens), promptTokens)
}

// outputBaseName ritorna il base name del file output (senza estensione).
// Pattern Sprint 1: <timestamp>_<profile>_<input_descriptor>.
// Post walk ricorsivo: ParsedFile.Name può essere un path relativo
// (es. "Documento_da_revisionare/MQL_Rev03.docx") — usiamo path.Base per
// estrarre solo il filename, così il `.log` accoppiato (runlog.OpenLogFile)
// scrive direttamente in outputDir invece di tentare una subdir inesistente.
func outputBaseName(p *profile.Profile, parsed *input.ParsedResult, started time.Time) string {
	ts := started.Format("2006-01-02T150405")
	descriptor := ""
	if len(parsed.Files) > 0 {
		first := path.Base(parsed.Files[0].Name)
		if idx := strings.LastIndex(first, "."); idx > 0 {
			first = first[:idx]
		}
		// Stesso slug di output.Write (nomi .md/.pdf/.docx), così il `.log`
		// accoppiato condivide il nome base e non contiene spazi/punti/accenti/
		// virgole (fragili su kDrive/shell/URL). Bug 2026-06-07.
		descriptor = "_" + output.Slug(first)
	}
	return ts + "_" + output.Slug(p.Profilo) + descriptor
}

// openAuditTrail apre il .log accoppiato e crea un MultiWriter logger.
// Sprint 1.5.C FR-32/NFR-11. Se l'apertura fallisce, fallback al logger originale
// (logging continua, audit trail incompleto ma run non si rompe).
// Ritorna (logFile, multiLog, bodyWriter). bodyWriter è dove drainStream scrive
// i token chunks raw: stderr+logFile in TTY, logFile only in non-TTY (FR-35).
//
// Fix review 1.5.C MEDIUM-Claude-2: usa isStderrTTY (stream effettivo del body
// streaming) invece di isStdoutTTY (era bug subtle che funzionava per il caso
// comune ma falliva con `2>/dev/null`).
func openAuditTrail(outputDir, baseName string, fallback *runlog.Logger) (*os.File, *runlog.Logger, io.Writer) {
	if outputDir == "" {
		return nil, fallback, io.Discard
	}
	logFile, err := runlog.OpenLogFile(outputDir, baseName)
	if err != nil {
		fallback.Warn("audit trail: impossibile aprire .log accoppiato (%v) — log file disabilitato per questo run", err)
		return nil, fallback, io.Discard
	}
	multiLog := runlog.NewMulti(os.Stderr, logFile)
	if isStderrTTY() {
		return logFile, multiLog, io.MultiWriter(os.Stderr, logFile)
	}
	return logFile, multiLog, logFile
}

// isStderrTTY ritorna true se os.Stderr è un terminale interattivo.
// Usato per decidere se il body streaming live va anche su stderr (TTY) o
// solo nel log file (non-TTY pipe/redirect). Fix review 1.5.C MEDIUM-Claude-2.
func isStderrTTY() bool {
	return term.IsTerminal(int(os.Stderr.Fd()))
}

// quiet=true: l'eventuale errore di stream NON viene loggato come warning (il
// chiamante lo gestisce — es. l'overflow di calibrazione, che è un passo atteso).
func callProviderWithStreaming(prov provider.LLMProvider, composed *prompt.Composed, p *profile.Profile, maxTokens int, log *runlog.Logger, bodyWriter io.Writer, quiet bool) (string, string, error) {
	ch, err := prov.Stream(context.Background(), composed.System, composed.User, provider.Options{
		Modello:       p.Modello,
		Temperature:   p.Temperature,
		MaxTokens:     maxTokens,
		ContextWindow: p.ContextWindow,
	})
	if err != nil {
		return "", "", err
	}
	body, stato, derr := drainStreamWithBody(ch, log, bodyWriter, quiet)
	return body, stato, derr
}

func writeOutput(cfg Config, p *profile.Profile, prov provider.LLMProvider, parsed *input.ParsedResult, check *tokens.CheckResult, body, stato string, duration time.Duration, started time.Time, log *runlog.Logger, logRelPath string) (string, *output.Frontmatter, error) {
	log.BeginStep("scrittura output")
	defer log.EndStep()
	fileNames := make([]string, len(parsed.Files))
	for i, f := range parsed.Files {
		fileNames[i] = f.Name
	}
	fm := &output.Frontmatter{
		Profilo:         p.Profilo,
		Modello:         p.Modello,
		Provider:        prov.Name(),
		DataEsecuzione:  started.Format(time.RFC3339),
		DurataSecondi:   duration.Seconds(),
		TokenStimati:    check.Tokens,
		FileInput:       fileNames,
		Stato:           stato,
		LogFile:         logRelPath,
		ProfileDefaults: p.Output.FrontmatterDefault,
	}
	outPath, err := output.Write(cfg.OutputDir, fm, body)
	if err != nil {
		return "", nil, err
	}
	log.Info("output: %s", outPath)
	return outPath, fm, nil
}

func exitCodeFor(stato string) int {
	if stato == "completato" {
		return 0
	}
	return 1
}

// finalizeOutput gestisce la risposta vuota del modello. Se il body è
// vuoto/solo-whitespace ritorna un .md con la spiegazione, stato "vuoto" ed
// empty=true; altrimenti lascia tutto invariato.
func finalizeOutput(body, stato string, p *profile.Profile, provName string) (string, string, bool) {
	if strings.TrimSpace(body) != "" {
		return body, stato, false
	}
	return emptyBodyNotice(p.Modello, provName), "vuoto", true
}

// emptyBodyNotice è il contenuto .md scritto quando il modello restituisce 0
// token: spiega cosa è successo invece di lasciare un file vuoto (Denis
// shakedown 2026-06-02 CAPABILITY F).
func emptyBodyNotice(modello, provName string) string {
	return fmt.Sprintf("> ⚠️ Il modello non ha restituito alcun output (0 token).\n>\n"+
		"> La richiesta è stata inviata correttamente — vedi il file `.log` accoppiato\n"+
		"> per il dettaglio di KB, dati di input e prompt effettivamente inviati.\n"+
		"> Il modello %q via %s ha però chiuso la risposta vuota.\n>\n"+
		"> Possibili cause: instabilità o cold-start del modello, sovraccarico del\n"+
		"> provider. **Riprova l'esecuzione.** Se persiste, valuta un modello diverso\n"+
		"> o segnala al referente tecnico.\n", modello, provName)
}

// logInputAudit registra nel .log COSA è stato inviato al modello: KB, dati,
// esclusioni, file NON parsati, stima token. Prova auditabile che i dati siano
// arrivati (o che un file sia stato saltato in parsing) — il .log si apre dopo
// il parsing, quindi senza questo riepilogo non lo registrerebbe.
func logInputAudit(log *runlog.Logger, kbTexts []string, parsed *input.ParsedResult, exclude []string, check *tokens.CheckResult, p *profile.Profile) {
	kbChars := 0
	for _, t := range kbTexts {
		kbChars += len(t)
	}
	log.Info("input inviato al modello:")
	log.Info("  KB system context: %d file, %d caratteri", len(kbTexts), kbChars)
	dataChars := 0
	for _, f := range parsed.Files {
		log.Info("  dato: %s (%d caratteri)", f.Name, f.Size)
		dataChars += f.Size
	}
	log.Info("  dati input: %d file, %d caratteri", len(parsed.Files), dataChars)
	for _, e := range exclude {
		if strings.TrimSpace(e) != "" {
			log.Info("  escluso dall'input (de-dup/exclude): %s", e)
		}
	}
	for _, sk := range parsed.Skipped {
		log.Warn("  NON parsato (saltato): %s — %s", sk.Name, sk.Reason)
	}
	log.Info("  stima totale: %d token / context %d", check.Tokens, p.ContextWindow)
}

func loadProfile(profiliDir, name string) (*profile.Profile, error) {
	// Fix review 1.5.C CRITICAL-mistral-3: defense-in-depth contro
	// path traversal via --profile flag o _labnexus.toml malformato.
	// name DEVE essere identifier semplice, no path components.
	if strings.ContainsAny(name, "/\\") || strings.Contains(name, "..") {
		return nil, fmt.Errorf("runner: nome profilo %q non consentito (path components rifiutati)", name)
	}
	path := filepath.Join(profiliDir, name+".toml")
	if _, err := os.Stat(path); err == nil {
		return profile.Load(path)
	}
	return nil, fmt.Errorf("runner: profilo non trovato: %s/%s.toml", profiliDir, name)
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
func drainStream(ch <-chan provider.StreamEvent, log *runlog.Logger) (string, string, error) {
	return drainStreamWithBody(ch, log, io.Discard, false)
}

// drainStreamWithBody è la variant di drainStream che riceve anche un
// bodyWriter dove scrivere raw i token chunks (FR-35 body streaming live).
// In TTY: bodyWriter = io.MultiWriter(stderr, logFile). Non-TTY: logFile only.
// quiet=true sopprime il warning "stream error" (il chiamante gestisce l'errore,
// es. l'overflow atteso durante la calibrazione del budget).
func drainStreamWithBody(ch <-chan provider.StreamEvent, log *runlog.Logger, bodyWriter io.Writer, quiet bool) (string, string, error) {
	var b strings.Builder
	// Default "interrotto": diventa "completato" solo se riceviamo esplicitamente done:true (EC-8).
	stato := "interrotto"
	started := time.Now()
	var tokenCount int
	var firstToken bool
	// Diagnostica modelli "thinking" (es. kimi): reasoning_content streamato a
	// parte dal contenuto + finish_reason. Spiega TTFT lunghi e risposte vuote.
	var reasoningChars int
	var sawReasoning bool
	var finishReason string
	var usage *provider.Usage
	finalize := func() {
		if finishReason != "" {
			log.Info("finish_reason: %s", finishReason)
		}
		if reasoningChars > 0 {
			log.Info("reasoning totale: %d caratteri di ragionamento (non inclusi nell'output)", reasoningChars)
		}
		if usage != nil {
			log.Info("%s", formatUsage(usage))
		}
	}

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
				if sawReasoning {
					log.Info("...il modello sta ragionando (reasoning: %d caratteri, %s elapsed)",
						reasoningChars, time.Since(started).Round(time.Second))
				} else {
					log.Info("...ancora in attesa del primo token (%s elapsed)",
						time.Since(started).Round(time.Second))
				}
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
				finalize()
				return b.String(), stato, nil
			}
			if ev.Err != nil {
				log.StreamEnd(tokenCount, time.Since(started))
				if !quiet {
					log.Warn("stream error: %v", ev.Err)
				}
				finalize()
				return b.String(), stato, ev.Err
			}
			if ev.Reasoning != "" {
				if !sawReasoning {
					sawReasoning = true
					log.Info("───── ragionamento del modello (reasoning) — mostrato e loggato, NON incluso nell'output .md ─────")
				}
				reasoningChars += len(ev.Reasoning)
				// Mostrato (TTY) + loggato (.log) via bodyWriter; NON va in `b`,
				// così il .md resta la risposta pulita. Niente del modello si butta.
				_, _ = bodyWriter.Write([]byte(ev.Reasoning))
			}
			if ev.FinishReason != "" {
				finishReason = ev.FinishReason
			}
			if ev.Usage != nil {
				usage = ev.Usage
			}
			if ev.Token != "" {
				if !firstToken {
					firstToken = true
					if reasoningChars > 0 {
						_, _ = bodyWriter.Write([]byte("\n\n───── output ─────\n\n"))
						log.Info("reasoning concluso: %d caratteri in %s → inizio generazione output",
							reasoningChars, time.Since(started).Round(time.Second))
					}
					log.Info("primo token ricevuto (TTFT %s), generazione in corso...",
						time.Since(started).Round(time.Millisecond))
				}
				b.WriteString(ev.Token)
				// FR-35 body streaming live: scrivi raw su bodyWriter (TTY: stderr+logFile; non-TTY: logFile only).
				_, _ = bodyWriter.Write([]byte(ev.Token))
				tokenCount++
			}
			if ev.Done {
				stato = "completato"
				if ev.NoDoneMarker {
					log.Warn("stream chiuso senza done marker dal server (output verosimilmente completo, ma il modello non ha emesso il chunk finale done:true — vedi bug eof-post-content-falsamente-interrotto)")
				}
				log.StreamEnd(tokenCount, time.Since(started))
				finalize()
				return b.String(), stato, nil
			}
		}
	}
}
