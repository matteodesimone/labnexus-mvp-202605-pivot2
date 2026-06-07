// Package config gestisce il file labnexus.config.toml master (FR-11..FR-15).
//
// Master fornisce defaults globali (provider, modello, parametri LLM, API key).
// I 7 profili shippati possono ometterli — vengono ereditati dal master.
// L'override del profile prevale sul master (precedenza profile-over-master).
package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/labnexus/labnexus/internal/profile"
)

// Master sono i defaults globali letti da labnexus.config.toml.
// I campi sono opzionali (zero-value detection a runtime); valori non setted
// non vengono propagati al merge.
type Master struct {
	Provider       string     `toml:"provider"`
	Modello        string     `toml:"modello"`
	Temperature    float64    `toml:"temperature"`
	MaxTokens      int        `toml:"max_tokens"`
	ContextWindow  int        `toml:"context_window"`
	EurouterAPIKey string     `toml:"eurouter_api_key"`
	OllamaEndpoint string     `toml:"ollama_endpoint"`
	PDF            PDFConfig  `toml:"pdf"`
	Docx           DocxConfig `toml:"docx"`

	// Tuning anti-stallo dello streaming (opzionali, 0/assente = default del codice).
	// StreamIdleTimeoutSec: secondi di silenzio totale (nessun token) oltre i quali
	// lo stream è considerato in stallo e la chiamata viene interrotta+ritentata
	// (default 180). StreamMaxRetries: numero di retry su stallo a metà generazione
	// o output vuoto (default 2). Lo stallo PRIMA del primo token (provider giù) ne
	// usa sempre 1 solo, per fallire in fretta.
	StreamIdleTimeoutSec int `toml:"stream_idle_timeout_sec"`
	StreamMaxRetries     int `toml:"stream_max_retries"`
}

// PDFConfig controlla la generazione del PDF accoppiato all'output MD.
// Enabled è *bool per distinguere "non setted" (nil) da "esplicitamente false".
// Vedi ResolvePDFEnabled per la gerarchia di risoluzione.
type PDFConfig struct {
	Enabled *bool `toml:"enabled"`
}

// DocxConfig controlla la generazione del DOCX accoppiato all'output MD.
// Stessa semantica di PDFConfig: *bool per distinguere non-setted da false
// esplicito. Vedi ResolveDocxEnabled per la gerarchia.
type DocxConfig struct {
	Enabled *bool `toml:"enabled"`
}

// Load legge il file TOML master e ritorna un *Master. Errori:
//   - fs.ErrNotExist se il file non esiste
//   - errore TOML con file:line se la sintassi è errata
func Load(path string) (*Master, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	var m Master
	if err := toml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("config: parse %s: %w", path, err)
	}
	return &m, nil
}

// Merge ritorna una copia del profile con i campi opzionali vuoti riempiti
// dal master. Precedenza: profile (se setted) > master.
// Se master è nil, ritorna profile invariato.
//
// Campi del master propagati AL PROFILE: Provider, Modello, Temperature,
// MaxTokens, ContextWindow.
//
// Campi del master propagati VIA ENV VARS al provider runtime (vedi
// internal/runner.applyMasterToEnv): EurouterAPIKey → EUROUTER_API_KEY,
// OllamaEndpoint → LABNEXUS_OLLAMA_ENDPOINT. Precedenza env > master:
// se l'env var è già setted, il valore del master è ignorato.
//
// LIMITAZIONE NOTA (review 1.5.B HIGH-config-1, mistral test-coverage):
// la detection "campo non setted" usa `== ""` (string) e `== 0` (numerico).
// Per i campi numerici Temperature/MaxTokens/ContextWindow, questo non
// distingue "campo omesso nel TOML" da "campo esplicitamente setted a 0".
// Caso d'uso problematico: profile con `temperature = 0.0` per output
// deterministico viene erroneamente override col valore master (es. 0.9).
// Fix richiede pointer types `*float64`/`*int` nel Profile (refactor invasivo,
// rimandato a 1.5.C+); workaround: settare temperature a un valore piccolo
// ma non zero (es. 0.001). Vedi `.pipeline/ideas/` per il backlog item.
func Merge(master *Master, p *profile.Profile) *profile.Profile {
	if master == nil || p == nil {
		return p
	}
	cp := *p
	if cp.Provider == "" {
		cp.Provider = master.Provider
	}
	if cp.Modello == "" {
		cp.Modello = master.Modello
	}
	if cp.Temperature == 0 {
		cp.Temperature = master.Temperature
	}
	if cp.MaxTokens == 0 {
		cp.MaxTokens = master.MaxTokens
	}
	if cp.ContextWindow == 0 {
		cp.ContextWindow = master.ContextWindow
	}
	return &cp
}

// Overlay è la funzione UNICA di override a strati: legge il TOML in path e ne
// sovrappone i parametri di config sopra `active`, applicando SOLO le chiavi
// PRESENTI nel file (md.IsDefined) e lasciando invariate le assenti. Va chiamata
// una volta per strato che deve sovrascrivere (es. il `_labnexus.toml` del job
// sopra il profilo già risolto) — catena config → profilo → _labnexus.
//
// Essendo presence-aware, distingue "campo a 0/” ESPLICITO" da "omesso" e
// risolve così la limitazione di Merge sui campi numerici (vedi commento sopra):
// `temperature = 0.0` nel file viene applicato, omesso viene ereditato.
//
// ESCLUSI di proposito: trigger_prompt/trigger_prompt_file (XOR FR-16 + sorgente
// tracciata nell'audit → path dedicato del runner) e [pdf]/[docx] (risolti con
// *bool da ResolvePDFEnabled/ResolveDocxEnabled). Aggiungere un nuovo parametro
// overridable = una riga qui, e funziona a tutti gli strati.
func Overlay(active *profile.Profile, path string) error {
	if active == nil {
		return fmt.Errorf("config: Overlay: active nil")
	}
	var tmp profile.Profile
	md, err := toml.DecodeFile(path, &tmp)
	if err != nil {
		return fmt.Errorf("config: overlay %s: %w", path, err)
	}
	if md.IsDefined("provider") {
		active.Provider = tmp.Provider
	}
	if md.IsDefined("modello") {
		active.Modello = tmp.Modello
	}
	if md.IsDefined("temperature") {
		active.Temperature = tmp.Temperature
	}
	if md.IsDefined("max_tokens") {
		active.MaxTokens = tmp.MaxTokens
	}
	if md.IsDefined("context_window") {
		active.ContextWindow = tmp.ContextWindow
	}
	if md.IsDefined("kb_files") {
		active.KbFiles = tmp.KbFiles
	}
	if md.IsDefined("descrizione") {
		active.Descrizione = tmp.Descrizione
	}
	if md.IsDefined("output", "frontmatter_default") {
		active.Output.FrontmatterDefault = tmp.Output.FrontmatterDefault
	}
	return nil
}
