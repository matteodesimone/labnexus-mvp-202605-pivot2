package features

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// stateKey è la chiave del context per scenarioState.
type stateKey struct{}

// getState estrae lo state corrente dal context.
func getState(c context.Context) *scenarioState {
	if s, ok := c.Value(stateKey{}).(*scenarioState); ok {
		return s
	}
	panic("scenarioState non in context — hook Before non registrato")
}

// scenarioState mantiene lo stato di un singolo scenario.
type scenarioState struct {
	tmpDir          string
	profiliDir      string
	kbDir           string
	outputDir       string
	env             map[string]string
	exitCode        int
	stdout          string
	stderr          string
	lastOutput      string
	fakeOllama      *httptest.Server
	fakeEurouter    *httptest.Server
	fixtureDirs     map[string]string
	inputDirs       map[string]string
	specialProfiles []string
	defaultInputDir string // input dir "preferita" registrata da step che hanno preparato dati ad hoc
	autoFill        bool   // true se lo scenario contiene placeholder ("...", "<X>") e vogliamo auto-completare i flag
}

// defaultProfile ritorna l'ultimo profilo speciale registrato o "revisione" come fallback.
func (s *scenarioState) defaultProfile() string {
	if n := len(s.specialProfiles); n > 0 {
		return s.specialProfiles[n-1]
	}
	return "revisione"
}

// registerSpecialProfile registra un profilo "ad hoc" creato da uno step.
func (s *scenarioState) registerSpecialProfile(name string) {
	s.specialProfiles = append(s.specialProfiles, name)
}

func newScenarioState() (*scenarioState, error) {
	tmp, err := os.MkdirTemp("", "labnexus-scenario-")
	if err != nil {
		return nil, err
	}
	s := &scenarioState{
		tmpDir:      tmp,
		profiliDir:  filepath.Join(tmp, "profili"),
		kbDir:       filepath.Join(tmp, "KB-ispettore"),
		outputDir:   filepath.Join(tmp, "out"),
		env:         map[string]string{},
		fixtureDirs: map[string]string{},
		inputDirs:   map[string]string{},
	}
	if err := os.MkdirAll(s.profiliDir, 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(s.kbDir, 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(s.outputDir, 0o755); err != nil {
		return nil, err
	}
	// Default kb file presente — i profili shippable nel repo lo citano
	if err := os.WriteFile(filepath.Join(s.kbDir, "CLAUDE.md"), []byte("# CLAUDE.md test stub\nIdentità ispettore."), 0o644); err != nil {
		return nil, err
	}
	// Profili default disponibili per ogni scenario (i singoli step possono sovrascrivere)
	defaultProfile := func(name string) string {
		return fmt.Sprintf(`profilo: %s
descrizione: profilo di test default per %s
provider: ollama
modello: qwen3.6
kb_files:
  - CLAUDE.md
trigger_prompt: |
  Trigger prompt di test di lunghezza sufficiente per soddisfare la soglia
  minima di 50 caratteri (FR-3). Tono ispettivo.
`, name, name)
	}
	for _, name := range []string{"revisione", "rilievi"} {
		if err := os.WriteFile(filepath.Join(s.profiliDir, name+".yml"), []byte(defaultProfile(name)), 0o644); err != nil {
			return nil, err
		}
	}
	// Default fixtures path mapping per i Test 1/2 reali
	s.fixtureDirs["<test1-input>"] = filepath.Join(repoRoot, "docs", "piano_iniziale", "materiali-dominio", "test-precedenti", "test-1-input")
	s.fixtureDirs["<test2-input>"] = filepath.Join(repoRoot, "docs", "piano_iniziale", "materiali-dominio", "test-precedenti", "test-2-input")
	s.fixtureDirs["<out>"] = s.outputDir
	s.fixtureDirs["<dir>"] = ""  // segnaposto, può essere riempito da step
	return s, nil
}

func (s *scenarioState) teardown() {
	if s.fakeOllama != nil {
		s.fakeOllama.Close()
	}
	if s.fakeEurouter != nil {
		s.fakeEurouter.Close()
	}
	if s.tmpDir != "" {
		_ = os.RemoveAll(s.tmpDir)
	}
}

// resolveDir mappa un path "simbolico" (es. "<test1-input>", "/tmp/X") al path reale.
func (s *scenarioState) resolveDir(symbolic string) string {
	if v, ok := s.fixtureDirs[symbolic]; ok && v != "" {
		return v
	}
	if v, ok := s.inputDirs[symbolic]; ok {
		return v
	}
	// Già un path assoluto reale → usalo
	if filepath.IsAbs(symbolic) {
		return symbolic
	}
	return filepath.Join(s.tmpDir, symbolic)
}

// ensureTmpDir crea (se necessario) una dir tmp per il path simbolico e ritorna il path reale.
func (s *scenarioState) ensureTmpDir(symbolic string) string {
	if v, ok := s.inputDirs[symbolic]; ok {
		return v
	}
	d, err := os.MkdirTemp(s.tmpDir, "in-")
	if err != nil {
		return s.resolveDir(symbolic)
	}
	s.inputDirs[symbolic] = d
	return d
}

// writeProfile crea un file in profiliDir col contenuto YAML dato.
func (s *scenarioState) writeProfile(name, body string) error {
	p := filepath.Join(s.profiliDir, name+".yml")
	return os.WriteFile(p, []byte(body), 0o644)
}

// writeKbFile crea un file in kbDir.
func (s *scenarioState) writeKbFile(name, body string) error {
	full := filepath.Join(s.kbDir, name)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return err
	}
	return os.WriteFile(full, []byte(body), 0o644)
}

// runBinary lancia ./labnexus con args, in ambiente con env override applicato.
// Sostituisce token simbolici nei flag (<test1-input>, <out>, /tmp/X reso reale).
func (s *scenarioState) runBinary(args []string) error {
	resolved := make([]string, len(args))
	for i, a := range args {
		resolved[i] = s.substitute(a)
	}
	// Aggiungi sempre i flag globali per pointing su nostre dir tmp
	resolved = append([]string{"--profiles-dir", s.profiliDir, "--kb-dir", s.kbDir}, resolved...)

	cmd := exec.Command(binaryPath, resolved...)
	cmd.Env = s.buildEnv()
	cmd.Dir = s.tmpDir

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	s.stdout = stdout.String()
	s.stderr = stderr.String()
	if err == nil {
		s.exitCode = 0
		return nil
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		s.exitCode = ee.ExitCode()
		return nil
	}
	s.exitCode = -1
	return fmt.Errorf("runBinary error: %w", err)
}

// substitute risolve i token simbolici dentro un singolo argomento.
func (s *scenarioState) substitute(in string) string {
	// Token comuni
	out := in
	for tok, path := range s.fixtureDirs {
		if path != "" {
			out = strings.ReplaceAll(out, tok, path)
		}
	}
	for tok, path := range s.inputDirs {
		out = strings.ReplaceAll(out, tok, path)
	}
	// `/tmp/Y` simbolico → se non in inputDirs, lascia
	return out
}

func (s *scenarioState) buildEnv() []string {
	base := os.Environ()
	// Rimuovi le env note che potrebbero leak dai test esterni
	keep := []string{}
	skip := map[string]bool{
		"LABNEXUS_PROVIDER":         true,
		"LABNEXUS_OLLAMA_ENDPOINT":  true,
		"LABNEXUS_EUROUTER_ENDPOINT": true,
		"EUROUTER_API_KEY":          true,
	}
	for _, kv := range base {
		eq := strings.IndexByte(kv, '=')
		if eq <= 0 {
			keep = append(keep, kv)
			continue
		}
		k := kv[:eq]
		if skip[k] {
			continue
		}
		keep = append(keep, kv)
	}
	// Aggiungi i nostri override (FR-12, dual-platform dev)
	if s.fakeOllama != nil {
		keep = append(keep, "LABNEXUS_OLLAMA_ENDPOINT="+s.fakeOllama.URL)
	}
	if s.fakeEurouter != nil {
		keep = append(keep, "LABNEXUS_EUROUTER_ENDPOINT="+s.fakeEurouter.URL)
	}
	// Bug #001 gate: nei test BDD il gate è SEMPRE settato perché
	// l'ambiente di test è controllato (httptest server locale o nessuna
	// vera chiamata cloud). Il gate stesso è coperto da unit test in
	// internal/provider/select_test.go.
	keep = append(keep, "LABNEXUS_ALLOW_CLOUD_PROVIDER=approved-for-synthetic-data")
	for k, v := range s.env {
		keep = append(keep, k+"="+v)
	}
	return keep
}

// startFakeOllama monta un httptest server che simula Ollama NDJSON.
// content è il testo che viene streamato chunk-per-chunk.
func (s *scenarioState) startFakeOllama(content string) {
	s.fakeOllama = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/x-ndjson")
		fl, _ := w.(http.Flusher)
		// Streama il content in 3 chunk
		chunks := splitInThree(content)
		for i, ch := range chunks {
			obj := map[string]interface{}{
				"message": map[string]string{"content": ch},
				"done":    i == len(chunks)-1,
			}
			b, _ := json.Marshal(obj)
			w.Write(b)
			w.Write([]byte("\n"))
			if fl != nil {
				fl.Flush()
			}
		}
	}))
}

// startFakeEurouter monta un httptest server che simula EUrouter SSE.
func (s *scenarioState) startFakeEurouter(content string) {
	s.fakeEurouter = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		fl, _ := w.(http.Flusher)
		chunks := splitInThree(content)
		for _, ch := range chunks {
			payload, _ := json.Marshal(map[string]interface{}{
				"choices": []map[string]interface{}{{"delta": map[string]string{"content": ch}}},
			})
			fmt.Fprintf(w, "data: %s\n\n", string(payload))
			if fl != nil {
				fl.Flush()
			}
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
}

func splitInThree(s string) []string {
	if s == "" {
		return []string{""}
	}
	n := len(s)
	a := n / 3
	b := 2 * n / 3
	return []string{s[:a], s[a:b], s[b:]}
}

// fakeOllamaHost ritorna host:port (senza scheme) dello fake server.
func (s *scenarioState) fakeOllamaHost() string {
	if s.fakeOllama == nil {
		return ""
	}
	u, _ := url.Parse(s.fakeOllama.URL)
	return u.Host
}

// ensureFakeOllama crea un fake server se non c'è (con default content).
func (s *scenarioState) ensureFakeOllama(content string) {
	if s.fakeOllama == nil {
		s.startFakeOllama(content)
	}
}

// findOutputFile cerca un file .md nella outputDir; ritorna il path se trovato.
func (s *scenarioState) findOutputFile() (string, error) {
	entries, err := os.ReadDir(s.outputDir)
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			return filepath.Join(s.outputDir, e.Name()), nil
		}
	}
	return "", fmt.Errorf("nessun file .md in outputDir")
}

// readOutputFile legge il contenuto dell'unico file .md in outputDir.
func (s *scenarioState) readOutputFile() (string, error) {
	p, err := s.findOutputFile()
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	s.lastOutput = p
	return string(b), nil
}

