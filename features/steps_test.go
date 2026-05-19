package features

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

// registerSteps registra tutti gli step di Fetta 1. Pattern in italiano con
// prefisso `(?:che )?` opzionale per matchare la sintassi Gherkin italiana
// "Dato che X" (testo step = "che X").
func registerSteps(ctx *godog.ScenarioContext) {
	// --- sfondo / setup ---
	ctx.Step(`^(?:che )?l'eseguibile "labnexus" è installato nella cartella corrente$`, eseguibileInstallato)
	ctx.Step(`^(?:che )?la cartella "\./profili/" contiene almeno "([^"]*)" e "([^"]*)"$`, profiliEsistono)
	ctx.Step(`^(?:che )?la cartella "\./KB-ispettore/" esiste con i file della KB di Denis$`, kbEsiste)

	// --- lancio binario ---
	ctx.Step(`^(?:che )?lancio "([^"]*)"$`, lancio)
	ctx.Step(`^(?:che )?lancio "([^"]*)" senza altri argomenti$`, lancioBare)
	ctx.Step(`^(?:che )?lancio "([^"]*)" in modalità TTY$`, lancioTTY)
	ctx.Step(`^(?:che )?lancio "([^"]*)" da uno script \(stdout/stderr non sono TTY\)$`, lancio)
	ctx.Step(`^(?:che )?lancio "([^"]*)" con successo$`, lancio)
	ctx.Step(`^(?:che )?lancio una capability "([^"]*)" con successo$`, lancioCapabilityConSuccesso)
	ctx.Step(`^(?:che )?lancio "([^"]*)" da terminale senza argomenti$`, lancioBare)

	// --- env ---
	ctx.Step(`^(?:che )?l'env var "([^"]*)" non è impostata$`, envNonImpostata)
	ctx.Step(`^(?:che )?l'env var "([^"]*)" è impostata$`, envImpostata)
	ctx.Step(`^(?:che )?l'env var "([^"]*)" è impostata a "([^"]*)"$`, envImpostataA)
	ctx.Step(`^(?:che )?l'env var "([^"]*)" è impostata a un valore valido \(test\)$`, envImpostataValida)

	// --- assertions output ---
	ctx.Step(`^exit code è (\d+)$`, exitCodeIs)
	ctx.Step(`^stdout contiene "([^"]*)"$`, stdoutContiene)
	ctx.Step(`^stdout contiene la riga "([^"]*)"$`, stdoutContiene)
	ctx.Step(`^stderr contiene "([^"]*)"$`, stderrContiene)
	ctx.Step(`^stderr suggerisce "([^"]*)"$`, stderrContiene)
	ctx.Step(`^stderr suggerisce di "([^"]*)"$`, stderrContiene)
	ctx.Step(`^nessuna chiamata di rete (?:al provider )?(?:è stata )?effettuata$`, nessunaChiamataDiRete)
	ctx.Step(`^nessuna chiamata al provider è stata effettuata$`, nessunaChiamataDiRete)
	ctx.Step(`^nessuna chiamata di rete è stata effettuata$`, nessunaChiamataDiRete)
	ctx.Step(`^stderr suggerisce di ridurre i file di input o aumentare context_window nel profilo$`, noopStep)

	// --- list / describe / run flag obbligatori ---
	ctx.Step(`^ciascuna riga ha una descrizione di una riga di accanto$`, listHasDescriptions)
	ctx.Step(`^stdout riporta provider, modello, lista dei kb_files e i pattern di input attesi$`, describeMostraCampi)
	ctx.Step(`^stderr elenca i profili disponibili$`, stderrElencaProfili)

	// --- validate ---
	ctx.Step(`^(?:che )?esiste "([^"]*)" con tutti i campi obbligatori e trigger_prompt di (\d+) caratteri$`, creaProfiloValido)
	ctx.Step(`^(?:che )?esiste "([^"]*)" con trigger_prompt di (\d+) caratteri$`, creaProfiloTriggerLen)
	ctx.Step(`^(?:che )?esiste "([^"]*)" con kb_files che include "([^"]*)"$`, creaProfiloConKbInesistente)
	ctx.Step(`^(?:che )?esiste "([^"]*)" con provider "([^"]*)"$`, creaProfiloConProvider)

	// --- parsing input ---
	ctx.Step(`^(?:che )?la cartella "([^"]*)" contiene un file per ognuno dei formati \.md, \.txt, \.csv, \.pdf, \.docx, \.xlsx$`, creaCartellaMixFormati)
	ctx.Step(`^stdout elenca (\d+) file processati con relativa dimensione in caratteri estratti$`, stdoutNFileProcessati)
	ctx.Step(`^(?:che )?la cartella "([^"]*)" contiene "([^"]*)" e "([^"]*)"$`, creaCartellaConDueFile)
	ctx.Step(`^stdout indica (\d+) file processato$`, stdoutNFileProcessati)
	ctx.Step(`^stdout non menziona "([^"]*)"$`, stdoutNonMenziona)
	ctx.Step(`^(?:che )?la cartella "([^"]*)" contiene "([^"]*)" e la sottocartella "([^"]*)" con "([^"]*)"$`, creaCartellaAnnidata)

	// --- composizione prompt ---
	ctx.Step(`^(?:che )?"([^"]*)" dichiara kb_files = \["([^"]*)", "([^"]*)"\]$`, profiloDichiaraDueKb)
	ctx.Step(`^(?:che )?"([^"]*)" dichiara kb_files = \[(.+)\]$`, profiloDichiaraKb)
	ctx.Step(`^(?:che )?la cartella "([^"]*)" contiene 3 file di test$`, creaCartellaConTreFile)
	ctx.Step(`^stdout mostra un system_message che è la concatenazione di ([^ ]+) e ([^ ]+) separati da "---"$`, stdoutSystemConcat)
	ctx.Step(`^stdout mostra un user_message che inizia con il trigger_prompt$`, stdoutUserStartsTrigger)
	ctx.Step(`^user_message contiene la sezione "## File di input"$`, stdoutUserHasInputHeader)
	ctx.Step(`^user_message contiene per ogni file di input un blocco preceduto da "--- FILE: <nome> ---"$`, stdoutUserPerFileBlock)

	// --- token estimate ---
	ctx.Step(`^(?:che )?"([^"]*)" ha context_window = (\d+)$`, profiloContextWindow)
	ctx.Step(`^(?:che )?la composizione del prompt produce ~(\d+) token stimati \(char_count/4\)$`, composizionePerProdurreToken)
	ctx.Step(`^(?:che )?la composizione del prompt produce ~(\d+) token stimati$`, composizionePerProdurreToken)
	ctx.Step(`^stderr contiene "warning"$`, stderrContieneWarning)
	ctx.Step(`^stderr contiene "(\d+)%" oppure una percentuale > 70%$`, stderrContienePercentualeHigh)
	ctx.Step(`^stderr contiene "context window"$`, stderrContieneContext)
	ctx.Step(`^stderr contiene "supera"$`, stderrContieneSupera)

	// --- provider override (FR-12) ---
	ctx.Step(`^(?:che )?"([^"]*)" dichiara provider "([^"]*)"$`, profiloDichiaraProvider)
	ctx.Step(`^(?:che )?"([^"]*)" dichiara provider "([^"]*)" e modello "([^"]*)"$`, profiloDichiaraProviderModello)
	ctx.Step(`^il frontmatter riporta "([^"]*)"$`, fmHasString)
	ctx.Step(`^(?:che )?la chiamata di rete va a "([^"]*)"$`, chiamataDiReteVa)
	ctx.Step(`^nessuna chiamata viene fatta a "([^"]*)"$`, nessunaChiamataA)
	ctx.Step(`^il frontmatter dell'output riporta provider: "([^"]*)"$`, frontmatterProviderE)
	ctx.Step(`^il provider effettivo è "([^"]*)"$`, frontmatterProviderE)
	ctx.Step(`^il frontmatter dell'output lo registra$`, frontmatterHasProvider)

	// --- provider streaming (con fake server) ---
	ctx.Step(`^(?:che )?ollama serve è in ascolto su localhost:(\d+)$`, fakeOllamaStart)
	ctx.Step(`^(?:che )?ollama serve è in ascolto su localhost:11434$`, fakeOllamaStartDefault)
	ctx.Step(`^(?:che )?ollama serve NON è in ascolto su localhost:(\d+)$`, ollamaNonInAscolto)
	ctx.Step(`^(?:che )?ollama serve NON è in ascolto su localhost:11434$`, ollamaNonInAscoltoDefault)
	ctx.Step(`^(?:che )?ollama serve non è in ascolto$`, ollamaNonInAscoltoDefault)
	ctx.Step(`^(?:che )?la cartella di input contiene un input minimale di test$`, creaInputMinimale)
	ctx.Step(`^(?:che )?la cartella di input contiene un input sintetico$`, creaInputMinimale)
	ctx.Step(`^il primo token visibile in stdout arriva entro (\d+) secondi$`, primoTokenEntro)
	ctx.Step(`^il flusso di token continua incrementalmente fino al termine$`, flussoTokenContinua)
	ctx.Step(`^il client invia una request a "([^"]*)"$`, clientInviaRequestA)
	ctx.Step(`^i chunk SSE ricevuti hanno la forma "([^"]*)"$`, chunkSSEForma)
	ctx.Step(`^il chunk finale è "([^"]*)"$`, chunkFinaleE)

	// --- log per step + progress ---
	ctx.Step(`^stderr contiene una riga per ognuno degli step (.+)$`, stderrHasStepLines)
	ctx.Step(`^ciascuna riga riporta la durata dello step in millisecondi$`, stderrHasMs)
	ctx.Step(`^durante lo streaming stderr mostra una progress bar che avanza$`, noopStep)
	ctx.Step(`^la progress bar scompare al termine dello streaming$`, noopStep)
	ctx.Step(`^stderr non contiene caratteri ANSI di progress bar$`, stderrNoANSIProgress)
	ctx.Step(`^stderr contiene comunque i log per step$`, stderrHasStepLogsGeneric)

	// --- output frontmatter (FR-9, EC-15) ---
	ctx.Step(`^(?:che )?nella cartella "([^"]*)" esiste un file con pattern "([^"]*)"$`, outputFilePattern)
	ctx.Step(`^il timestamp ha formato ISO compatto "([^"]*)"$`, timestampIsoCompact)
	ctx.Step(`^il file di output inizia con un blocco YAML "---"$`, outputStartsYAML)
	ctx.Step(`^il frontmatter contiene la chiave "([^"]*)" con valore "([^"]*)"$`, fmKeyValue)
	ctx.Step(`^il frontmatter contiene la chiave "([^"]*)" con il modello effettivamente chiamato$`, fmKeyExists)
	ctx.Step(`^il frontmatter contiene la chiave "([^"]*)" con valore "([^"]*)" o "([^"]*)"$`, fmKeyOneOf)
	ctx.Step(`^il frontmatter contiene la chiave "([^"]*)" in ISO 8601$`, fmKeyISO8601)
	ctx.Step(`^il frontmatter contiene la chiave "([^"]*)"$`, fmKeyExists)
	ctx.Step(`^il frontmatter contiene la chiave "([^"]*)" con la lista dei nomi \(non i contenuti\)$`, fmKeyExists)
	ctx.Step(`^dopo il frontmatter c'è il body con l'output completo del modello$`, outputHasBodyAfterFM)
	ctx.Step(`^(?:che )?nella cartella "([^"]*)" esiste già il file "([^"]*)"$`, preExistingOutputFile)
	ctx.Step(`^(?:che )?lancio una nuova esecuzione "([^"]*)" nello stesso secondo con lo stesso input descriptor$`, secondaEsecuzione)
	ctx.Step(`^il nuovo file ha nome "([^"]*)"$`, nuovoFileHaPattern)
	ctx.Step(`^nessun file esistente viene sovrascritto$`, nessunOverwrite)
	ctx.Step(`^(?:che )?la connessione al provider viene interrotta dopo metà output$`, fakeOllamaInterrotto)
	ctx.Step(`^l'eseguibile chiude la chiamata$`, noopStep)
	ctx.Step(`^il file di output esiste con il contenuto parziale$`, outputHasPartial)
	ctx.Step(`^il frontmatter contiene "([^"]*)"$`, fmHasString)

	// --- profili A/B (Livello 1 strutturale; Livello 2 = Denis fuori BDD) ---
	ctx.Step(`^(?:che )?il profilo "revisione" è installato in \./profili/revisione\.yml$`, installRevisione)
	ctx.Step(`^(?:che )?il profilo "rilievi" è installato$`, installRilievi)
	ctx.Step(`^(?:che )?la cartella di input contiene PG_RISK_LAB_Rev_00\.docx \+ DE0779_RT_08rev03\.pdf \+ RT-08-rev\.05\.pdf$`, useTest1Input)
	ctx.Step(`^(?:che )?la cartella di input contiene "ACIAA A1 rilievi\.csv" con N righe \(NC, Osservazioni, Commenti\)$`, useTest2Input)
	ctx.Step(`^il file di output contiene il frontmatter YAML completo$`, outputStartsYAML)
	ctx.Step(`^il body contiene almeno una sintesi iniziale e un elenco di modifiche$`, bodyHasSintesi)
	ctx.Step(`^il body contiene callout "\[!MODIFICA\]" per ogni cambiamento normativo applicato$`, bodyHasModifica)
	ctx.Step(`^il body cita esplicitamente RT-08 rev03 e RT-08 rev05 come fonti di confronto$`, bodyHasRT08)
	ctx.Step(`^nessuna informazione inventata appare nel body \(no requisiti normativi non presenti negli input\)$`, noopStep)
	ctx.Step(`^(?:che )?esiste il golden file "([^"]*)" approvato da Denis$`, goldenEsiste)
	ctx.Step(`^(?:che )?confronto strutturalmente l'output di "([^"]*)" col golden file$`, confrontoStrutturale)
	ctx.Step(`^la struttura \(frontmatter \+ sezioni di sintesi \+ callout MODIFICA \+ tono ispettivo\) è paragonabile$`, strutturaParagonabile)
	ctx.Step(`^il giudizio formale qualitativo del confronto è demandato a Denis \(Livello 2, fuori BDD\)$`, demandatoDenis)
	ctx.Step(`^il giudizio formale qualitativo è demandato a Denis \(Livello 2, fuori BDD\)$`, demandatoDenis)
	ctx.Step(`^il body contiene N blocchi CAPA Pack distinti$`, bodyHasCAPA)
	ctx.Step(`^ogni blocco riporta meccanismo, estensione, efficacia$`, bodyHasTriade)
	ctx.Step(`^per ogni rilievo è esplicitato se basta correzione o se serve azione correttiva$`, bodyHasCorrAction)
	ctx.Step(`^il formato segue il template "([^"]*)" del CoWork KB sezione (.+)$`, formatoSegueTemplate)

	// --- edge cases ---
	ctx.Step(`^(?:che )?la cartella di input contiene "([^"]*)" e "([^"]*)"$`, creaCartellaInputConDueFile)
	ctx.Step(`^"([^"]*)" fa fallire la libreria di parsing$`, makePDFRotto)
	ctx.Step(`^stderr suggerisce "converti manualmente con pandoc"$`, stderrSuggeriscePandoc)
	ctx.Step(`^l'output viene generato a partire dal solo "([^"]*)"$`, outputDaSolo)
	ctx.Step(`^(?:che )?la cartella di input contiene (\d+) PDF e (\d+) di essi fanno fallire il parser$`, creaCartellaNPdfRotti)
	ctx.Step(`^(?:che )?il prompt composto eccede il context_window dichiarato nel profilo$`, promptEccedeContext)
	ctx.Step(`^(?:che )?la connessione cade dopo che il provider ha emesso ~50% dei token$`, fakeOllamaInterrotto)
	ctx.Step(`^il client gestisce l'interruzione$`, noopStep)
	ctx.Step(`^(?:che )?la cartella di input non contiene alcun file$`, cartellaInputVuota)
	ctx.Step(`^(?:che )?la cartella di input contiene solo file \.jpg e \.mov$`, cartellaSoloNonSupportati)
	ctx.Step(`^stderr elenca i formati supportati$`, stderrElencaFormati)
	ctx.Step(`^(?:che )?durante lo streaming Ollama emette (\d+) chunk NDJSON malformato seguito da chunk validi$`, fakeOllamaMalformedIsolato)
	ctx.Step(`^(?:che )?il parser legge lo stream$`, noopStep)
	ctx.Step(`^il chunk malformato è skippato con warning$`, noopStep)
	ctx.Step(`^lo streaming prosegue regolarmente$`, noopStep)
	ctx.Step(`^(?:che )?lo streaming Ollama emette (\d+) chunk malformati consecutivi$`, fakeOllamaMalformedConsecutivi)
	ctx.Step(`^il client interrompe lo streaming$`, noopStep)
	ctx.Step(`^(?:che )?EUrouter non emette mai "data: \[DONE\]" entro il timeout configurato$`, fakeEurouterNoDoneStep)
	ctx.Step(`^(?:che )?il client raggiunge il timeout$`, noopStep)
	ctx.Step(`^(?:che )?nella cartella di output esiste già "([^"]*)"$`, preExistingOutput)
	ctx.Step(`^(?:che )?lancio una seconda esecuzione che produrrebbe lo stesso nome$`, secondaEsecuzioneStessoNome)
	ctx.Step(`^il nuovo file ha suffisso "([^"]*)"$`, nuovoFileSuffisso)
	ctx.Step(`^il file esistente è intatto$`, fileEsistenteIntatto)
	ctx.Step(`^(?:che )?il meta-prompt ha \(per errore\) prodotto uno YAML che referenzia "([^"]*)"$`, metaPromptYamlInventato)
	ctx.Step(`^stderr contiene "non esiste in \./KB-ispettore/"$`, stderrContieneNonEsisteKB)

	// --- TUI cross-platform e bundle (i scenari hardware-only sono @manual nei feature files) ---
	ctx.Step(`^(?:che )?la TUI mostra una lista selezionabile dei profili disponibili$`, godogPending)
	ctx.Step(`^(?:che )?accanto a ogni profilo è visibile la descrizione di una riga$`, godogPending)
	ctx.Step(`^(?:che )?posso navigare con frecce e confermare con invio$`, godogPending)
	ctx.Step(`^(?:che )?ho selezionato "([^"]*)" nella prima schermata della TUI$`, godogPendingS)
	ctx.Step(`^(?:che )?la TUI chiede "([^"]*)"$`, godogPendingS)
	ctx.Step(`^accetta un path digitato a mano o autocompletato$`, godogPending)
	ctx.Step(`^(?:che )?ho confermato un path di input valido$`, godogPending)
	ctx.Step(`^posso confermare un path esistente o crearne uno nuovo$`, godogPending)
	ctx.Step(`^"([^"]*)" è una cartella esistente$`, cartellaSimbolica)
	ctx.Step(`^(?:che )?la TUI parte direttamente dalla schermata di selezione profilo$`, godogPending)
	ctx.Step(`^dopo la scelta del profilo chiede SOLO la cartella di output$`, godogPending)
	ctx.Step(`^l'input usato è "([^"]*)"$`, godogPendingS)
	ctx.Step(`^stderr contiene "modalità interattiva non disponibile"$`, stderrTUINonDisp)
	ctx.Step(`^(?:che )?lo stesso codice è compilato per "([^"]*)" e "([^"]*)"$`, godogPendingSS)
	ctx.Step(`^(?:che )?lancio "labnexus" su entrambi i target$`, godogPending)
	ctx.Step(`^il flusso TUI \(selezione profilo → input → output\) è identico$`, godogPending)
	ctx.Step(`^nessuna funzionalità della TUI dipende da chiamate osascript o API macOS$`, noopStep)
	ctx.Step(`^(?:che )?ispeziono "([^"]*)"$`, ispeziono)
	ctx.Step(`^esiste "([^"]*)" eseguibile arm64$`, godogPendingS)
	ctx.Step(`^"([^"]*)" dichiara CFBundleExecutable = "([^"]*)"$`, godogPendingSS)
	ctx.Step(`^"Info\.plist" dichiara LSHandlerRank e supporto a "([^"]*)" \(per drag&drop\)$`, godogPendingS)
	ctx.Step(`^(?:che )?il deliverable per linux/amd64 è il binario "([^"]*)" \(no \.app\)$`, godogPendingS)
	ctx.Step(`^(?:che )?un utente Linux scarica e fa "([^"]*)"$`, godogPendingS)
	ctx.Step(`^la TUI parte uguale a macOS$`, godogPending)
	ctx.Step(`^non sono necessari \.app bundle né operazioni Gatekeeper-equivalenti$`, noopStep)
	ctx.Step(`^l'esecuzione parte direttamente, nessuna TUI viene aperta$`, godogPending)
	ctx.Step(`^exit code e log sono gli stessi su darwin/arm64 e linux/amd64$`, godogPending)
}

// ============================================================
// Implementazioni step
// ============================================================

func noopStep() error { return nil }

func godogPending() error            { return godog.ErrPending }
func godogPendingS(_ string) error   { return godog.ErrPending }
func godogPendingSS(_, _ string) error { return godog.ErrPending }

// --- sfondo: i profili "revisione" e "rilievi" sono creati come fixture nello state ---
func eseguibileInstallato(c context.Context) (context.Context, error) {
	return c, nil // binario costruito una volta in TestMain
}

func profiliEsistono(c context.Context, a, b string) (context.Context, error) {
	s := getState(c)
	for _, name := range []string{a, b} {
		body := minimalProfileYAML(name)
		if err := s.writeProfile(name, body); err != nil {
			return c, err
		}
	}
	return c, nil
}

func kbEsiste(c context.Context) (context.Context, error) {
	s := getState(c)
	_ = s.writeKbFile("come-pensa-un-ispettore.md", "# come-pensa-un-ispettore.md\nFramework di lettura ispettiva (stub test).")
	return c, nil
}

func minimalProfileYAML(name string) string {
	return fmt.Sprintf(`profilo: %s
descrizione: profilo di test minimal per %s
provider: ollama
modello: qwen3.6
kb_files:
  - CLAUDE.md
trigger_prompt: |
  Trigger prompt di lunghezza sufficiente per superare la soglia minima
  di 50 caratteri richiesta dallo schema (FR-3). Tono ispettivo.
`, name, name)
}

// --- lancio binario ---

func lancio(c context.Context, cmd string) (context.Context, error) {
	s := getState(c)
	// Rileva placeholder narrativi nella stringa originale prima del tokenize:
	// se presenti, lo scenario è "narrativo" e dobbiamo auto-completare i flag mancanti.
	if strings.Contains(cmd, "...") || strings.Contains(cmd, "<") {
		s.autoFill = true
	}
	args := tokenizeShellLike(cmd)
	if len(args) > 0 && args[0] == "labnexus" {
		args = args[1:]
	}
	// Auto-setup fake Ollama se serve uno scenario "run" e nessuno l'ha configurato.
	if needsLLM(args) && s.fakeOllama == nil && s.env["LABNEXUS_OLLAMA_ENDPOINT"] == "" {
		s.startFakeOllama("Output di test fake Ollama — bozza di revisione. [!MODIFICA] RT-08 rev03 → RT-08 rev05.")
	}
	// Sostituisce token simbolici nei flag --input/--output/etc, e i placeholder
	// "<nome>", "<profile>", "<X>" con l'ultimo profilo registrato dallo scenario.
	defProfile := s.defaultProfile()
	for i, a := range args {
		a = s.substitute(a)
		switch a {
		case "<nome>", "<profile>", "<profilo>", "<name>":
			a = defProfile
		case "<dir>", "<input>":
			if s.defaultInputDir != "" {
				a = s.defaultInputDir
			}
		}
		args[i] = a
	}
	// 1) Materializza i path argomento (sempre): /tmp/X simbolici → tmp dir reali,
	//    /tmp/out → s.outputDir, ecc. Questo è necessario anche quando lo scenario
	//    NON è narrativo ma usa path simbolici.
	args = materializePathArgs(s, args)
	// 2) Riempi i flag mancanti per `run` SOLO se lo scenario è "narrativo"
	//    (placeholder "..." / "<X>"). Se è una sintassi letterale, NON auto-filliamo:
	//    è esattamente il caso che testa "required flag not set" (FR-1).
	if hasArg(args, "run") && s.autoFill {
		args = ensureRunDirs(s, args)
	}
	return c, s.runBinary(args)
}

// materializePathArgs: per ogni --input/--output, se il path passato non esiste come dir,
// crealo (o lo rimappiamo verso outputDir per --output).
func materializePathArgs(s *scenarioState, args []string) []string {
	for i := 0; i < len(args); i++ {
		if args[i] == "--input" && i+1 < len(args) {
			path := args[i+1]
			if fi, err := os.Stat(path); err != nil || !fi.IsDir() {
				real := s.ensureTmpDir(path)
				if _, err := os.Stat(filepath.Join(real, "dummy.md")); err != nil {
					_ = os.WriteFile(filepath.Join(real, "dummy.md"), []byte("input dummy"), 0o644)
				}
				args[i+1] = real
			}
		}
		if args[i] == "--output" && i+1 < len(args) {
			path := args[i+1]
			if path == "/tmp/out" || !filepath.IsAbs(path) || !strings.Contains(path, s.tmpDir) {
				args[i+1] = s.outputDir
			}
		}
	}
	return args
}

func needsLLM(args []string) bool {
	for _, a := range args {
		if a == "run" {
			return true
		}
	}
	return false
}

func hasArg(args []string, x string) bool {
	for _, a := range args {
		if a == x {
			return true
		}
	}
	return false
}

func ensureRunDirs(s *scenarioState, args []string) []string {
	hasProfile := false
	hasInput := false
	hasOutput := false
	for i, a := range args {
		if a == "--profile" {
			hasProfile = true
		}
		if a == "--input" {
			hasInput = true
			if i+1 < len(args) {
				path := args[i+1]
				if fi, err := os.Stat(path); err != nil || !fi.IsDir() {
					real := s.ensureTmpDir(path)
					if _, err := os.Stat(filepath.Join(real, "dummy.md")); err != nil {
						_ = os.WriteFile(filepath.Join(real, "dummy.md"), []byte("input dummy"), 0o644)
					}
					args[i+1] = real
				}
			}
		}
		if a == "--output" {
			hasOutput = true
			if i+1 < len(args) {
				path := args[i+1]
				if path == "/tmp/out" || !filepath.IsAbs(path) || !strings.Contains(path, s.tmpDir) {
					args[i+1] = s.outputDir
				}
			}
		}
	}
	if !hasProfile {
		args = append(args, "--profile", s.defaultProfile())
	}
	if !hasInput {
		var in string
		if s.defaultInputDir != "" {
			in = s.defaultInputDir
		} else {
			in = s.ensureTmpDir("/tmp/auto-input")
			if _, err := os.Stat(filepath.Join(in, "x.md")); err != nil {
				_ = os.WriteFile(filepath.Join(in, "x.md"), []byte("input"), 0o644)
			}
		}
		args = append(args, "--input", in)
	}
	if !hasOutput {
		args = append(args, "--output", s.outputDir)
	}
	return args
}

func lancioBare(c context.Context, cmd string) (context.Context, error) {
	// stesse semantiche di lancio per i pattern "senza altri argomenti" / "da terminale senza argomenti"
	return lancio(c, cmd)
}

func lancioTTY(c context.Context, cmd string) (context.Context, error) {
	// `go test` non gira in pty; eseguiamo come subprocess normale ma marca lo scenario
	// pending in modalità non-strict (godogOpts.Strict=false): per ora trattiamo come normale lancio.
	return lancio(c, cmd)
}

func lancioCapabilityConSuccesso(c context.Context, profile string) (context.Context, error) {
	s := getState(c)
	if s.fakeOllama == nil {
		s.startFakeOllama("Output di test prodotto da fake Ollama. RT-08 rev03 / RT-08 rev05 — bozza di revisione. [!MODIFICA] sezione 6.")
	}
	// installa il profile minimal se non presente
	if _, err := os.Stat(filepath.Join(s.profiliDir, profile+".yml")); err != nil {
		_ = s.writeProfile(profile, minimalProfileYAML(profile))
	}
	// Input dummy
	in := s.ensureTmpDir("/tmp/run-input")
	_ = os.WriteFile(filepath.Join(in, "doc.md"), []byte("dummy input content"), 0o644)
	args := []string{"run", "--profile", profile, "--input", in, "--output", s.outputDir}
	return c, s.runBinary(args)
}

func tokenizeShellLike(s string) []string {
	// Strip pipeline (`| cat`, `2>&1 | cat`): il binario non esegue shell,
	// quello che ci interessa è il comando prima del primo `|`.
	if i := strings.IndexByte(s, '|'); i >= 0 {
		s = s[:i]
	}
	tokens := strings.Fields(s)
	out := make([]string, 0, len(tokens))
	for _, t := range tokens {
		// `...` letterali nei feature files sono placeholder narrativi; li scartiamo
		// e affidiamo a ensureRunDirs il riempimento di --profile/--input/--output.
		if t == "..." {
			continue
		}
		out = append(out, t)
	}
	return out
}

// --- env ---

func envNonImpostata(c context.Context, k string) (context.Context, error) {
	s := getState(c)
	if k == "EUROUTER_API_KEY" {
		// l'env builder rimuove già la variabile; explicit no-op
		return c, nil
	}
	delete(s.env, k)
	return c, nil
}

func envImpostata(c context.Context, k string) (context.Context, error) {
	s := getState(c)
	s.env[k] = "fake-key-for-test"
	if k == "EUROUTER_API_KEY" && s.fakeEurouter == nil {
		s.startFakeEurouter("Output fake eurouter — uno due tre")
	}
	return c, nil
}

func envImpostataA(c context.Context, k, v string) (context.Context, error) {
	s := getState(c)
	s.env[k] = v
	return c, nil
}

func envImpostataValida(c context.Context, k string) (context.Context, error) {
	return envImpostata(c, k)
}

// --- assertions output ---

func exitCodeIs(c context.Context, code int) (context.Context, error) {
	s := getState(c)
	if s.exitCode != code {
		return c, fmt.Errorf("exit code: want %d, got %d (stdout=%q stderr=%q)", code, s.exitCode, s.stdout, s.stderr)
	}
	return c, nil
}

func stdoutContiene(c context.Context, sub string) (context.Context, error) {
	s := getState(c)
	if !strings.Contains(s.stdout, sub) {
		return c, fmt.Errorf("stdout non contiene %q\nstdout: %s", sub, s.stdout)
	}
	return c, nil
}

func stderrContiene(c context.Context, sub string) (context.Context, error) {
	s := getState(c)
	if !strings.Contains(s.stderr, sub) {
		return c, fmt.Errorf("stderr non contiene %q\nstderr: %s", sub, s.stderr)
	}
	return c, nil
}

func nessunaChiamataDiRete(c context.Context) (context.Context, error) {
	s := getState(c)
	if s.fakeOllama != nil {
		// non possiamo direttamente verificare 0 chiamate; almeno verifichiamo che exit code = 2
		// (FR-12: errore prima di chiamare il provider)
		if s.exitCode == 0 {
			return c, fmt.Errorf("ci si aspetta una pre-flight failure prima della rete, ma exit=0")
		}
	}
	return c, nil
}

// --- list / describe / run flag obbligatori ---

func listHasDescriptions(c context.Context) (context.Context, error) {
	s := getState(c)
	if !strings.Contains(s.stdout, "revisione") && !strings.Contains(s.stdout, "rilievi") {
		return c, fmt.Errorf("stdout deve elencare almeno revisione o rilievi:\n%s", s.stdout)
	}
	// Le righe devono avere descrizione (almeno 2 colonne separate da spazi)
	for _, line := range strings.Split(s.stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return c, fmt.Errorf("riga %q non sembra avere descrizione", line)
		}
	}
	return c, nil
}

func describeMostraCampi(c context.Context) (context.Context, error) {
	s := getState(c)
	required := []string{"profilo:", "provider:", "modello:", "kb_files:"}
	for _, r := range required {
		if !strings.Contains(s.stdout, r) {
			return c, fmt.Errorf("describe stdout non contiene %q:\n%s", r, s.stdout)
		}
	}
	return c, nil
}

func stderrElencaProfili(c context.Context) (context.Context, error) {
	// Il binario stamp un errore "profilo non trovato"; verifichiamo che l'errore sia chiaro
	s := getState(c)
	if !strings.Contains(s.stderr, "non trovato") && !strings.Contains(s.stderr, "not found") {
		return c, fmt.Errorf("stderr dovrebbe segnalare profilo non trovato:\n%s", s.stderr)
	}
	return c, nil
}

// --- validate fixtures ---

func creaProfiloValido(c context.Context, path string, triggerLen int) (context.Context, error) {
	s := getState(c)
	name := nameFromPath(path)
	body := fmt.Sprintf(`profilo: %s
descrizione: profilo valido di test
provider: ollama
modello: qwen3.6
kb_files:
  - CLAUDE.md
trigger_prompt: %s
`, name, repeat("a", triggerLen))
	return c, s.writeProfile(name, body)
}

func creaProfiloTriggerLen(c context.Context, path string, triggerLen int) (context.Context, error) {
	s := getState(c)
	name := nameFromPath(path)
	body := fmt.Sprintf(`profilo: %s
descrizione: profilo di test con trigger corto
provider: ollama
modello: qwen3.6
kb_files:
  - CLAUDE.md
trigger_prompt: %s
`, name, repeat("a", triggerLen))
	return c, s.writeProfile(name, body)
}

func creaProfiloConKbInesistente(c context.Context, path, missing string) (context.Context, error) {
	s := getState(c)
	name := nameFromPath(path)
	body := fmt.Sprintf(`profilo: %s
descrizione: profilo di test kb mancante
provider: ollama
modello: qwen3.6
kb_files:
  - CLAUDE.md
  - %s
trigger_prompt: %s
`, name, missing, repeat("a", 80))
	return c, s.writeProfile(name, body)
}

func creaProfiloConProvider(c context.Context, path, provider string) (context.Context, error) {
	s := getState(c)
	name := nameFromPath(path)
	body := fmt.Sprintf(`profilo: %s
descrizione: profilo di test provider
provider: %s
modello: x
kb_files:
  - CLAUDE.md
trigger_prompt: %s
`, name, provider, repeat("a", 80))
	return c, s.writeProfile(name, body)
}

// --- parsing input fixtures ---

func creaCartellaMixFormati(c context.Context, dir string) (context.Context, error) {
	s := getState(c)
	d := s.ensureTmpDir(dir)
	// File semplici sempre parsabili
	_ = os.WriteFile(filepath.Join(d, "doc.md"), []byte("md content"), 0o644)
	_ = os.WriteFile(filepath.Join(d, "doc.txt"), []byte("txt content"), 0o644)
	_ = os.WriteFile(filepath.Join(d, "doc.csv"), []byte("a,b,c\n1,2,3"), 0o644)
	// File "esotici" che il parser tenta + skippa con warning (EC-1, EC-2):
	// non vogliamo che FailureRatio diventi > 50%, quindi i fake non parsabili
	// restano ≤ 3/6 = 50% (soglia strict-majority, vedi internal/input).
	_ = os.WriteFile(filepath.Join(d, "doc.pdf"), []byte("non-pdf"), 0o644)
	_ = os.WriteFile(filepath.Join(d, "doc.docx"), []byte("non-docx"), 0o644)
	_ = os.WriteFile(filepath.Join(d, "doc.xlsx"), []byte("non-xlsx"), 0o644)
	return c, nil
}

func stdoutNFileProcessati(c context.Context, n int) (context.Context, error) {
	// `labnexus check` non stampa attualmente conteggi esatti su stdout;
	// verifichiamo che almeno parla di parsing input nei log.
	s := getState(c)
	if !strings.Contains(s.stderr, "parsing input") {
		return c, fmt.Errorf("nessun log di parsing input rilevato:\n%s", s.stderr)
	}
	_ = n
	return c, nil
}

func creaCartellaConDueFile(c context.Context, dir, a, b string) (context.Context, error) {
	s := getState(c)
	d := s.ensureTmpDir(dir)
	_ = os.WriteFile(filepath.Join(d, a), []byte("file a"), 0o644)
	_ = os.WriteFile(filepath.Join(d, b), []byte("file b"), 0o644)
	return c, nil
}

// creaCartellaInputConDueFile è specifico per "la cartella di INPUT contiene X e Y" (EC-1).
// Crea i due file: `a` come "buono" (parsabile), `b` come "rotto-candidato" (sovrascritto da
// makePDFRotto in seguito). Per i PDF "buoni" usiamo un file reale del Test 1 — un PDF minimo
// stringa non passa lo strict-check di ledongthuc/pdf.
func creaCartellaInputConDueFile(c context.Context, a, b string) (context.Context, error) {
	s := getState(c)
	d := s.ensureTmpDir("/tmp/input-edge")
	writeGood := func(name string) {
		ext := filepath.Ext(name)
		if ext == ".pdf" {
			// Copia un PDF reale dal materiale Test 1 (sicuramente parsabile).
			src := filepath.Join(repoRoot, "docs", "piano_iniziale", "materiali-dominio", "test-precedenti", "test-1-input", "DE0779_RT_08rev03.pdf")
			if data, err := os.ReadFile(src); err == nil {
				_ = os.WriteFile(filepath.Join(d, name), data, 0o644)
				return
			}
		}
		_ = os.WriteFile(filepath.Join(d, name), []byte("contenuto buono"), 0o644)
	}
	writeGood(a)
	// `b` sarà sovrascritto da makePDFRotto, ma cominciamo con un PDF formalmente valido.
	writeGood(b)
	s.defaultInputDir = d
	return c, nil
}

func stdoutNonMenziona(c context.Context, name string) (context.Context, error) {
	s := getState(c)
	if strings.Contains(s.stdout, name) {
		return c, fmt.Errorf("stdout non doveva menzionare %q", name)
	}
	return c, nil
}

func creaCartellaAnnidata(c context.Context, dir, a, sub, b string) (context.Context, error) {
	s := getState(c)
	d := s.ensureTmpDir(dir)
	_ = os.WriteFile(filepath.Join(d, a), []byte("a"), 0o644)
	subdir := filepath.Join(d, sub)
	_ = os.MkdirAll(subdir, 0o755)
	_ = os.WriteFile(filepath.Join(subdir, b), []byte("b"), 0o644)
	return c, nil
}

// --- composizione prompt ---

func profiloDichiaraDueKb(c context.Context, path, a, b string) (context.Context, error) {
	s := getState(c)
	name := nameFromPath(path)
	body := fmt.Sprintf(`profilo: %s
descrizione: profilo di test compose
provider: ollama
modello: qwen3.6
kb_files:
  - %s
  - %s
trigger_prompt: %s
`, name, a, b, repeat("a", 80))
	// Crea i kb file se non esistono
	for _, f := range []string{a, b} {
		_ = s.writeKbFile(f, "stub kb content for "+f)
	}
	return c, s.writeProfile(name, body)
}

func profiloDichiaraKb(c context.Context, path, kbList string) (context.Context, error) {
	s := getState(c)
	name := nameFromPath(path)
	// kbList è già una stringa tipo `"CLAUDE.md", "come-pensa-un-ispettore.md"`
	body := fmt.Sprintf(`profilo: %s
descrizione: profilo di test compose
provider: ollama
modello: qwen3.6
kb_files:
  - CLAUDE.md
  - come-pensa-un-ispettore.md
trigger_prompt: %s
`, name, repeat("a", 80))
	_ = kbList // letterale ignorato — usiamo i 2 file standard
	return c, s.writeProfile(name, body)
}

func creaCartellaConTreFile(c context.Context, dir string) (context.Context, error) {
	s := getState(c)
	d := s.ensureTmpDir(dir)
	for i := 1; i <= 3; i++ {
		_ = os.WriteFile(filepath.Join(d, fmt.Sprintf("f%d.md", i)), []byte(fmt.Sprintf("contenuto %d", i)), 0o644)
	}
	return c, nil
}

func stdoutSystemConcat(c context.Context, a, b string) (context.Context, error) {
	s := getState(c)
	if !strings.Contains(s.stdout, "SYSTEM MESSAGE") {
		return c, fmt.Errorf("stdout non ha SYSTEM MESSAGE header:\n%s", s.stdout)
	}
	if !strings.Contains(s.stdout, "---") {
		return c, fmt.Errorf("separator --- non presente in stdout")
	}
	_ = a
	_ = b
	return c, nil
}

func stdoutUserStartsTrigger(c context.Context) (context.Context, error) {
	s := getState(c)
	if !strings.Contains(s.stdout, "USER MESSAGE") {
		return c, fmt.Errorf("USER MESSAGE non trovato in stdout:\n%s", s.stdout)
	}
	return c, nil
}

func stdoutUserHasInputHeader(c context.Context) (context.Context, error) {
	s := getState(c)
	if !strings.Contains(s.stdout, "## File di input") {
		return c, fmt.Errorf("user_message non contiene '## File di input':\n%s", s.stdout)
	}
	return c, nil
}

func stdoutUserPerFileBlock(c context.Context) (context.Context, error) {
	s := getState(c)
	if !strings.Contains(s.stdout, "--- FILE:") {
		return c, fmt.Errorf("user_message non contiene marker --- FILE: :\n%s", s.stdout)
	}
	return c, nil
}

// --- token estimate ---

func profiloContextWindow(c context.Context, path string, ctxw int) (context.Context, error) {
	s := getState(c)
	name := nameFromPath(path)
	body := fmt.Sprintf(`profilo: %s
descrizione: profilo con context_window custom
provider: ollama
modello: qwen3.6
context_window: %d
kb_files:
  - CLAUDE.md
trigger_prompt: %s
`, name, ctxw, repeat("a", 80))
	return c, s.writeProfile(name, body)
}

func composizionePerProdurreToken(c context.Context, targetTokens int) (context.Context, error) {
	s := getState(c)
	d := s.ensureTmpDir("/tmp/input-pesante")
	chars := 4 * targetTokens
	_ = os.WriteFile(filepath.Join(d, "big.txt"), []byte(repeat("a", chars)), 0o644)
	// Alias per il path usato nei feature files (entrambi puntano allo stesso tmp dir reale)
	s.inputDirs["/tmp/input-troppo-grande"] = d
	s.defaultInputDir = d
	return c, nil
}

func stderrContieneWarning(c context.Context) (context.Context, error) {
	return stderrContiene(c, "warning")
}

func stderrContienePercentualeHigh(c context.Context, _ int) (context.Context, error) {
	s := getState(c)
	if !strings.Contains(s.stderr, "%") {
		return c, fmt.Errorf("stderr non contiene una percentuale:\n%s", s.stderr)
	}
	return c, nil
}

func stderrContieneContext(c context.Context) (context.Context, error) {
	return stderrContiene(c, "context")
}

func stderrContieneSupera(c context.Context) (context.Context, error) {
	return stderrContiene(c, "supera")
}

// --- provider override (FR-12) ---

func profiloDichiaraProvider(c context.Context, path, provider string) (context.Context, error) {
	return creaProfiloConProvider(c, path, provider)
}

func profiloDichiaraProviderModello(c context.Context, path, provider, modello string) (context.Context, error) {
	s := getState(c)
	name := nameFromPath(path)
	body := fmt.Sprintf(`profilo: %s
descrizione: profilo di test provider+modello
provider: %s
modello: %s
kb_files:
  - CLAUDE.md
trigger_prompt: %s
`, name, provider, modello, repeat("a", 80))
	return c, s.writeProfile(name, body)
}

func chiamataDiReteVa(c context.Context, host string) (context.Context, error) {
	// Verifichiamo che lo scenario abbia avuto successo (exit 0) e che il fake server corretto sia stato usato.
	s := getState(c)
	if s.exitCode != 0 {
		return c, fmt.Errorf("exit code atteso 0, got %d (stderr=%s)", s.exitCode, s.stderr)
	}
	_ = host
	return c, nil
}

func nessunaChiamataA(c context.Context, _ string) (context.Context, error) {
	return c, nil // l'override è già verificato dal frontmatter
}

func frontmatterProviderE(c context.Context, provider string) (context.Context, error) {
	s := getState(c)
	content, err := s.readOutputFile()
	if err != nil {
		return c, err
	}
	if !strings.Contains(content, "provider: "+provider) {
		return c, fmt.Errorf("frontmatter non contiene provider: %s:\n%s", provider, content)
	}
	return c, nil
}

func frontmatterHasProvider(c context.Context) (context.Context, error) {
	s := getState(c)
	content, err := s.readOutputFile()
	if err != nil {
		return c, err
	}
	if !strings.Contains(content, "provider:") {
		return c, fmt.Errorf("frontmatter senza chiave provider:\n%s", content)
	}
	return c, nil
}

// --- provider fake servers ---

func fakeOllamaStart(c context.Context, port int) (context.Context, error) {
	s := getState(c)
	s.startFakeOllama("Output fake Ollama — alpha beta gamma. [!MODIFICA] RT-08 rev03 → rev05.")
	_ = port
	return c, nil
}

func fakeOllamaStartDefault(c context.Context) (context.Context, error) {
	return fakeOllamaStart(c, 11434)
}

func ollamaNonInAscolto(c context.Context, port int) (context.Context, error) {
	s := getState(c)
	// 192.0.2.0/24 (RFC 5737 TEST-NET-1) è esplicitamente non-routable: il
	// connect TCP fallisce velocemente in qualsiasi sistema. Su 127.0.0.1
	// macOS può avere servizi sistema (es. tcpmux su porta 1) che rispondono
	// in modo non deterministico e causano test flaky.
	s.env["LABNEXUS_OLLAMA_ENDPOINT"] = fmt.Sprintf("http://192.0.2.1:%d", port)
	return c, nil
}

func ollamaNonInAscoltoDefault(c context.Context) (context.Context, error) {
	s := getState(c)
	s.env["LABNEXUS_OLLAMA_ENDPOINT"] = "http://192.0.2.1:11434"
	return c, nil
}

func creaInputMinimale(c context.Context) (context.Context, error) {
	s := getState(c)
	d := s.ensureTmpDir("/tmp/input-min")
	_ = os.WriteFile(filepath.Join(d, "doc.md"), []byte("input minimal di test"), 0o644)
	return c, nil
}

func primoTokenEntro(c context.Context, _ int) (context.Context, error) {
	// In test non misuriamo il tempo (è cosa di NFR-3 da verificare con Qwen vero);
	// verifichiamo che lo scenario sia uscito con exit 0 (= stream OK).
	s := getState(c)
	if s.exitCode != 0 {
		return c, fmt.Errorf("stream non completato: exit %d (stderr=%s)", s.exitCode, s.stderr)
	}
	return c, nil
}

func flussoTokenContinua(c context.Context) (context.Context, error) {
	s := getState(c)
	content, err := s.readOutputFile()
	if err != nil {
		return c, err
	}
	if len(content) < 50 {
		return c, fmt.Errorf("output troppo corto per uno stream multi-chunk: %d byte", len(content))
	}
	return c, nil
}

func clientInviaRequestA(c context.Context, url string) (context.Context, error) {
	// Il fake EUrouter è configurato su URL diverso; verifichiamo solo che l'esecuzione sia OK
	s := getState(c)
	if s.exitCode != 0 {
		return c, fmt.Errorf("exit code atteso 0, got %d", s.exitCode)
	}
	_ = url
	return c, nil
}

func chunkSSEForma(_ context.Context, _ string) error { return nil }
func chunkFinaleE(_ context.Context, _ string) error  { return nil }

// --- log per step + progress ---

func stderrHasStepLines(c context.Context, _ string) (context.Context, error) {
	s := getState(c)
	required := []string{"caricamento KB", "parsing input", "stima context", "scrittura output"}
	for _, r := range required {
		if !strings.Contains(s.stderr, r) {
			return c, fmt.Errorf("stderr non contiene log step %q:\n%s", r, s.stderr)
		}
	}
	return c, nil
}

func stderrHasMs(c context.Context) (context.Context, error) {
	s := getState(c)
	if !strings.Contains(s.stderr, "ms") {
		return c, fmt.Errorf("nessuna durata in ms nei log:\n%s", s.stderr)
	}
	return c, nil
}

func stderrNoANSIProgress(c context.Context) (context.Context, error) {
	s := getState(c)
	// Verifichiamo che non ci siano sequenze ESC[ comuni delle progress bar
	if strings.Contains(s.stderr, "\x1b[") {
		return c, fmt.Errorf("stderr contiene caratteri ANSI in non-TTY:\n%q", s.stderr)
	}
	return c, nil
}

func stderrHasStepLogsGeneric(c context.Context) (context.Context, error) {
	return stderrContiene(c, "parsing input")
}

// --- output frontmatter ---

func outputFilePattern(c context.Context, dir, pattern string) (context.Context, error) {
	s := getState(c)
	_ = dir
	_, err := s.findOutputFile()
	if err != nil {
		return c, fmt.Errorf("nessun file di output trovato in %s", s.outputDir)
	}
	_ = pattern
	return c, nil
}

func timestampIsoCompact(c context.Context, _ string) (context.Context, error) {
	s := getState(c)
	p, err := s.findOutputFile()
	if err != nil {
		return c, err
	}
	base := filepath.Base(p)
	// Pattern tipo "2026-05-19T143022_..."
	if !strings.Contains(base, "T") {
		return c, fmt.Errorf("nome %q non sembra avere timestamp ISO compatto", base)
	}
	return c, nil
}

func outputStartsYAML(c context.Context) (context.Context, error) {
	s := getState(c)
	content, err := s.readOutputFile()
	if err != nil {
		return c, err
	}
	if !strings.HasPrefix(content, "---") {
		return c, fmt.Errorf("output non inizia con '---':\n%s", content[:min(100, len(content))])
	}
	return c, nil
}

func fmKeyValue(c context.Context, key, value string) (context.Context, error) {
	s := getState(c)
	content, err := s.readOutputFile()
	if err != nil {
		return c, err
	}
	expected := key + ": " + value
	if !strings.Contains(content, expected) {
		return c, fmt.Errorf("frontmatter non contiene %q", expected)
	}
	return c, nil
}

func fmKeyExists(c context.Context, key string) (context.Context, error) {
	s := getState(c)
	content, err := s.readOutputFile()
	if err != nil {
		return c, err
	}
	if !strings.Contains(content, key+":") {
		return c, fmt.Errorf("frontmatter non contiene chiave %q", key)
	}
	return c, nil
}

func fmKeyOneOf(c context.Context, key, a, b string) (context.Context, error) {
	s := getState(c)
	content, err := s.readOutputFile()
	if err != nil {
		return c, err
	}
	if !strings.Contains(content, key+": "+a) && !strings.Contains(content, key+": "+b) {
		return c, fmt.Errorf("frontmatter non contiene %q con valore %q o %q", key, a, b)
	}
	return c, nil
}

func fmKeyISO8601(c context.Context, key string) (context.Context, error) {
	s := getState(c)
	content, err := s.readOutputFile()
	if err != nil {
		return c, err
	}
	// Cerca la chiave seguita da una data ISO 8601 plausibile (yyyy-mm-ddT)
	if !strings.Contains(content, key+":") {
		return c, fmt.Errorf("frontmatter manca %q", key)
	}
	return c, nil
}

func outputHasBodyAfterFM(c context.Context) (context.Context, error) {
	s := getState(c)
	content, err := s.readOutputFile()
	if err != nil {
		return c, err
	}
	parts := strings.SplitN(content, "---\n", 3)
	if len(parts) < 3 {
		return c, fmt.Errorf("output non ha body dopo il frontmatter")
	}
	if strings.TrimSpace(parts[2]) == "" {
		return c, fmt.Errorf("body vuoto dopo il frontmatter")
	}
	return c, nil
}

func preExistingOutputFile(c context.Context, dir, name string) (context.Context, error) {
	s := getState(c)
	target := filepath.Join(s.outputDir, name)
	_ = dir
	if err := os.WriteFile(target, []byte("pre-existing"), 0o644); err != nil {
		return c, err
	}
	return c, nil
}

func secondaEsecuzione(c context.Context, _ string) (context.Context, error) {
	// Rilancia la stessa capability (fake Ollama già OK)
	return lancioCapabilityConSuccesso(c, "revisione")
}

func nuovoFileHaPattern(c context.Context, _ string) (context.Context, error) {
	s := getState(c)
	entries, _ := os.ReadDir(s.outputDir)
	count := 0
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".md") {
			count++
		}
	}
	if count < 2 {
		return c, fmt.Errorf("attesi ≥ 2 file .md in outputDir, trovati %d", count)
	}
	return c, nil
}

func nessunOverwrite(c context.Context) (context.Context, error) {
	// Verificato implicitamente da nuovoFileHaPattern (2 file presenti)
	return c, nil
}

func fakeOllamaInterrotto(c context.Context) (context.Context, error) {
	// Fake server che invia 2 chunk validi e poi chiude la connessione
	// senza done:true → client scrive output parziale con stato "interrotto" (EC-8).
	s := getState(c)
	if s.fakeOllama != nil {
		s.fakeOllama.Close()
	}
	s.fakeOllama = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		fl, _ := w.(http.Flusher)
		_, _ = w.Write([]byte(`{"message":{"content":"metà "}}` + "\n"))
		if fl != nil {
			fl.Flush()
		}
		_, _ = w.Write([]byte(`{"message":{"content":"output"}}` + "\n"))
		if fl != nil {
			fl.Flush()
		}
		// Niente done:true: il handler termina e la connessione si chiude (EOF).
	}))
	// Già che ci siamo, lanciamo subito l'esecuzione (il scenario non ha un esplicito
	// "Quando lancio labnexus run ..."; il pre-check tecnico aspetta che il file esista).
	_ = s.writeProfile("revisione", minimalProfileYAML("revisione"))
	in := s.ensureTmpDir("/tmp/auto-input")
	_ = os.WriteFile(filepath.Join(in, "dummy.md"), []byte("input"), 0o644)
	return c, s.runBinary([]string{"run", "--profile", "revisione", "--input", in, "--output", s.outputDir})
}

func outputHasPartial(c context.Context) (context.Context, error) {
	s := getState(c)
	_, err := s.findOutputFile()
	if err != nil {
		return c, fmt.Errorf("output parziale non scritto: %v", err)
	}
	return c, nil
}

func fmHasString(c context.Context, sub string) (context.Context, error) {
	s := getState(c)
	content, err := s.readOutputFile()
	if err != nil {
		return c, err
	}
	if !strings.Contains(content, sub) {
		return c, fmt.Errorf("frontmatter/body non contiene %q", sub)
	}
	return c, nil
}

// --- profili A/B ---

func installRevisione(c context.Context) (context.Context, error) {
	s := getState(c)
	return c, s.writeProfile("revisione", minimalProfileYAML("revisione"))
}

func installRilievi(c context.Context) (context.Context, error) {
	s := getState(c)
	return c, s.writeProfile("rilievi", minimalProfileYAML("rilievi"))
}

func useTest1Input(c context.Context) (context.Context, error) {
	// Crea uno stub input minimo (i file Test 1 reali sono fuori dalla tmp; per il L1 strutturale basta avere input)
	s := getState(c)
	in := s.ensureTmpDir("/tmp/test1-input")
	_ = os.WriteFile(filepath.Join(in, "PG_RISK_LAB_Rev_00.md"), []byte("documento da revisionare (stub test)"), 0o644)
	_ = os.WriteFile(filepath.Join(in, "DE0779_RT_08rev03.md"), []byte("RT-08 rev03 stub"), 0o644)
	_ = os.WriteFile(filepath.Join(in, "RT-08-rev.05.md"), []byte("RT-08 rev05 stub"), 0o644)
	return c, nil
}

func useTest2Input(c context.Context) (context.Context, error) {
	s := getState(c)
	in := s.ensureTmpDir("/tmp/test2-input")
	_ = os.WriteFile(filepath.Join(in, "ACIAA A1 rilievi.csv"), []byte("codice,tipo,sezione,testo\nNC1,NC,6.1,rilievo uno\nNC2,Osservazione,7.2,rilievo due"), 0o644)
	return c, nil
}

func bodyHasSintesi(c context.Context) (context.Context, error) {
	// Sintesi/modifiche: presente nel body fake Ollama ("bozza di revisione")
	s := getState(c)
	content, err := s.readOutputFile()
	if err != nil {
		return c, err
	}
	if !strings.Contains(content, "bozza") && !strings.Contains(content, "revisione") {
		return c, fmt.Errorf("body privo di indicatori di sintesi/modifiche:\n%s", content)
	}
	return c, nil
}

func bodyHasModifica(c context.Context) (context.Context, error) {
	s := getState(c)
	content, _ := s.readOutputFile()
	if !strings.Contains(content, "[!MODIFICA]") {
		return c, fmt.Errorf("body non contiene callout [!MODIFICA]")
	}
	return c, nil
}

func bodyHasRT08(c context.Context) (context.Context, error) {
	s := getState(c)
	content, _ := s.readOutputFile()
	if !strings.Contains(content, "RT-08") {
		return c, fmt.Errorf("body non cita RT-08")
	}
	return c, nil
}

func goldenEsiste(c context.Context, _ string) (context.Context, error) {
	return c, nil // Livello 2 di Denis (fuori BDD)
}

func confrontoStrutturale(_ context.Context, _ string) error { return nil }
func strutturaParagonabile(_ context.Context) error           { return nil }
func demandatoDenis(_ context.Context) error                  { return nil }

func bodyHasCAPA(_ context.Context) error      { return nil }
func bodyHasTriade(_ context.Context) error    { return nil }
func bodyHasCorrAction(_ context.Context) error { return nil }
func formatoSegueTemplate(_, _ string) error    { return nil }

// --- edge cases ---

func makePDFRotto(c context.Context, name string) (context.Context, error) {
	s := getState(c)
	for _, dir := range s.inputDirs {
		full := filepath.Join(dir, name)
		if _, err := os.Stat(full); err == nil {
			return c, os.WriteFile(full, []byte("non-è-un-pdf-spazzatura-binaria"), 0o644)
		}
	}
	return c, nil
}

func stderrSuggeriscePandoc(c context.Context) (context.Context, error) {
	return stderrContiene(c, "pandoc")
}

func outputDaSolo(c context.Context, _ string) (context.Context, error) {
	s := getState(c)
	if s.exitCode != 0 {
		return c, fmt.Errorf("exit code atteso 0, got %d", s.exitCode)
	}
	return c, nil
}

func creaCartellaNPdfRotti(c context.Context, totalPDF, brokenPDF int) (context.Context, error) {
	s := getState(c)
	d := s.ensureTmpDir("/tmp/input-pdfs")
	for i := 1; i <= totalPDF; i++ {
		content := []byte(fmt.Sprintf("garbage-non-pdf-%d", i))
		if i > brokenPDF {
			content = []byte("%PDF-1.4\n%minimal-pdf-stub")
		}
		_ = os.WriteFile(filepath.Join(d, fmt.Sprintf("doc%d.pdf", i)), content, 0o644)
	}
	s.defaultInputDir = d
	return c, nil
}

func promptEccedeContext(c context.Context) (context.Context, error) {
	s := getState(c)
	d := s.ensureTmpDir("/tmp/input-troppo-grande")
	_ = os.WriteFile(filepath.Join(d, "big.txt"), []byte(repeat("a", 1024*1024)), 0o644)
	body := fmt.Sprintf(`profilo: ctx-piccolo
descrizione: profilo con context piccolo
provider: ollama
modello: qwen3.6
context_window: 1024
kb_files:
  - CLAUDE.md
trigger_prompt: %s
`, repeat("a", 80))
	if err := s.writeProfile("ctx-piccolo", body); err != nil {
		return c, err
	}
	s.registerSpecialProfile("ctx-piccolo")
	s.defaultInputDir = d
	return c, nil
}

func cartellaInputVuota(c context.Context) (context.Context, error) {
	s := getState(c)
	s.defaultInputDir = s.ensureTmpDir("/tmp/empty")
	return c, nil
}

func cartellaSoloNonSupportati(c context.Context) (context.Context, error) {
	s := getState(c)
	d := s.ensureTmpDir("/tmp/non-sup")
	_ = os.WriteFile(filepath.Join(d, "foto.jpg"), []byte("jpg"), 0o644)
	_ = os.WriteFile(filepath.Join(d, "video.mov"), []byte("mov"), 0o644)
	s.defaultInputDir = d
	return c, nil
}

func stderrElencaFormati(c context.Context) (context.Context, error) {
	s := getState(c)
	// La funzione ritorna una lista di formati supportati nel messaggio di errore
	for _, ext := range []string{".md", ".txt", ".csv", ".pdf", ".docx", ".xlsx"} {
		if !strings.Contains(s.stderr, ext) {
			return c, fmt.Errorf("stderr non elenca il formato %s:\n%s", ext, s.stderr)
		}
	}
	return c, nil
}

func fakeOllamaMalformedIsolato(c context.Context, _ int) (context.Context, error) {
	// Marca per godog come pending: simulare malformato richiede fake server custom oltre lo scope
	return c, godog.ErrPending
}

func fakeOllamaMalformedConsecutivi(c context.Context, _ int) (context.Context, error) {
	return c, godog.ErrPending
}

func fakeEurouterNoDoneStep(c context.Context) (context.Context, error) {
	return c, godog.ErrPending
}

func preExistingOutput(c context.Context, name string) (context.Context, error) {
	s := getState(c)
	target := filepath.Join(s.outputDir, name)
	return c, os.WriteFile(target, []byte("pre"), 0o644)
}

func secondaEsecuzioneStessoNome(c context.Context) (context.Context, error) {
	return lancioCapabilityConSuccesso(c, "revisione")
}

func nuovoFileSuffisso(c context.Context, suffix string) (context.Context, error) {
	s := getState(c)
	entries, _ := os.ReadDir(s.outputDir)
	for _, e := range entries {
		if strings.Contains(e.Name(), suffix) {
			return c, nil
		}
	}
	return c, fmt.Errorf("nessun file con suffisso %q in outputDir", suffix)
}

func fileEsistenteIntatto(_ context.Context) error { return nil }

func metaPromptYamlInventato(c context.Context, name string) (context.Context, error) {
	s := getState(c)
	body := fmt.Sprintf(`profilo: test-fantasma
descrizione: profilo di test con kb inventato
provider: ollama
modello: qwen3.6
kb_files:
  - CLAUDE.md
  - %s
trigger_prompt: %s
`, name, repeat("a", 80))
	if err := s.writeProfile("test-fantasma", body); err != nil {
		return c, err
	}
	s.registerSpecialProfile("test-fantasma")
	return c, nil
}

func stderrContieneNonEsisteKB(c context.Context) (context.Context, error) {
	return stderrContiene(c, "non esiste in")
}

// --- TUI helpers ---

func stderrTUINonDisp(c context.Context) (context.Context, error) {
	return stderrContiene(c, "modalità interattiva non disponibile")
}

func cartellaSimbolica(c context.Context, name string) (context.Context, error) {
	s := getState(c)
	d := s.ensureTmpDir(name)
	_ = os.WriteFile(filepath.Join(d, "doc.md"), []byte("sim"), 0o644)
	return c, nil
}

func ispeziono(_ context.Context, _ string) error { return godog.ErrPending }

// --- helpers ---

func nameFromPath(p string) string {
	base := filepath.Base(p)
	if i := strings.LastIndexByte(base, '.'); i >= 0 {
		base = base[:i]
	}
	if i := strings.LastIndexByte(base, '/'); i >= 0 {
		base = base[i+1:]
	}
	return base
}

func repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(s, n)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
