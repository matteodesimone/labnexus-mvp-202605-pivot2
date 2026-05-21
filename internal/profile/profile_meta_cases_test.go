package profile_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/profile"
)

// TestMetaPromptCases — FR-22: validazione funzionale del meta-prompt
// `docs/meta-prompt-genera-profilo.md` (Fetta 3 sub-fetta 3.C).
//
// Per ogni file `caso-*.yml` in `.pipeline/test-data/scenarios/feat-meta/
// casi-inventati/`, il test:
//  1. carica il YAML (profile.Load)
//  2. lo valida contro la KB-ispettore reale (profile.Validate)
//
// Tutti DEVONO passare. Se uno fallisce → il meta-prompt va aggiornato
// (oppure lo schema profile è cambiato in modo incompatibile e i casi
// devono essere ri-generati).
//
// In red phase (subito dopo /v-test-scaffold di Fetta 3), la cartella
// `casi-inventati/` contiene solo il README e nessun `caso-*.yml`: il
// test ritorna FAIL per "zero casi" (verifica esplicita che almeno 3
// casi devono essere committati prima di /v-implement green).
//
// In green phase (post /v-implement quando Matteo ha lanciato i 3 casi
// in Claude esterno e committato i YAML), il test diventa verde.
//
// Regression-safe: se in Sprint 2 il profile.Validate si arricchisce di
// nuovi campi obbligatori, i 3 casi committati produrranno FAIL — segnale
// che il meta-prompt va aggiornato in coerenza con lo schema.
func TestMetaPromptCases(t *testing.T) {
	// Trova la repo root risalendo dal pacchetto. Il test gira da
	// internal/profile/ quindi salendo di 2 livelli arriviamo a repo root.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	repoRoot := filepath.Dir(filepath.Dir(wd))

	casiDir := filepath.Join(repoRoot, ".pipeline", "test-data",
		"scenarios", "feat-meta", "casi-inventati")
	kbDir := filepath.Join(repoRoot, "docs", "piano_iniziale",
		"materiali-dominio", "KB-ispettore")

	if _, err := os.Stat(casiDir); err != nil {
		t.Fatalf("cartella casi-inventati non trovata: %s — eseguire test-scaffold di Fetta 3", casiDir)
	}
	if _, err := os.Stat(kbDir); err != nil {
		t.Fatalf("KB-ispettore di progetto non trovata: %s", kbDir)
	}

	entries, err := os.ReadDir(casiDir)
	if err != nil {
		t.Fatalf("read casi-inventati: %v", err)
	}

	yamlFiles := []string{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasPrefix(name, "caso-") {
			continue
		}
		// Sprint 1.5.B post-migration: il formato corrente è TOML, ma accettiamo
		// anche YAML per backward compat se Denis ha esempi storici.
		if !strings.HasSuffix(name, ".yml") && !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".toml") {
			continue
		}
		yamlFiles = append(yamlFiles, filepath.Join(casiDir, name))
	}

	// RED phase: nessun caso-*.{toml,yml} committato → fail con messaggio
	// diagnostico esplicito (Matteo deve eseguire i 3 casi manualmente in
	// Claude esterno e committare i file in formato corrente TOML).
	if len(yamlFiles) == 0 {
		t.Fatalf("FR-22 red phase: 0 casi in %s — Matteo deve eseguire manualmente i 3 casi inventati (semplice/medio/ambiguo) in Claude esterno con docs/meta-prompt-genera-profilo.md, salvare gli output come caso-1-semplice.toml, caso-2-medio.toml, caso-3-ambiguo.toml (formato Sprint 1.5.B). Vedi README.md della cartella per istruzioni.", casiDir)
	}

	// Idealmente esattamente 3 casi (FR-22). Più sono OK; meno è un
	// warning di coverage incompleta.
	if len(yamlFiles) < 3 {
		t.Errorf("attesi ≥ 3 casi (FR-22: semplice, medio, ambiguo), trovati %d in %s", len(yamlFiles), casiDir)
	}

	for _, yamlPath := range yamlFiles {
		caso := filepath.Base(yamlPath)
		t.Run(caso, func(t *testing.T) {
			p, err := profile.Load(yamlPath)
			if err != nil {
				t.Fatalf("profile.Load(%q): %v — il meta-prompt ha generato un YAML non parsabile", yamlPath, err)
			}
			if err := profile.Validate(p, kbDir); err != nil {
				t.Fatalf("profile.Validate fallita per %q contro KB reale: %v — il meta-prompt ha generato un YAML che viola lo schema FR-3 (campi obbligatori, trigger ≥ 50 char, kb_files esistono in KB-ispettore, provider in {ollama, eurouter})", yamlPath, err)
			}
		})
	}
}

// TestMetaPromptContent — FR-20 review loop 1 fix (Mistral test-coverage):
// verifica che `docs/meta-prompt-genera-profilo.md` contenga le sezioni
// strutturali richieste dallo spec FR-20. Senza questo test, una modifica
// inavvertita al meta-prompt (rimozione sezione, drift) non sarebbe
// catturata. Insieme a TestMetaPromptCases (validazione funzionale) e
// TestGuidaMetaPromptContent (sotto), copre i deliverable testuali FR-20/21.
func TestMetaPromptContent(t *testing.T) {
	wd, _ := os.Getwd()
	repoRoot := filepath.Dir(filepath.Dir(wd))
	metaPath := filepath.Join(repoRoot, "docs", "meta-prompt-genera-profilo.md")

	b, err := os.ReadFile(metaPath)
	if err != nil {
		t.Fatalf("meta-prompt non trovato: %s — FR-20 deliverable mancante", metaPath)
	}
	content := string(b)

	// FR-20: deve avere le 8 sezioni canoniche del meta-prompt.
	// Heading H2 numerati "## N. <titolo>".
	requiredSections := []struct {
		num     int
		keyword string // case-insensitive substring del titolo
	}{
		{1, "identità"},
		{2, "sistema"},
		{3, "schema"},
		{4, "kb-ispettore"},
		{5, "esempi"},
		{6, "clarification"},
		{7, "output finale"},
		{8, "vincoli"},
	}
	for _, sec := range requiredSections {
		// Match ## <num>. <whatever> con keyword case-insensitive nel
		// resto della riga.
		pattern := regexp.MustCompile(`(?i)(?m)^## ` + regexp.QuoteMeta(itoa(sec.num)) + `\. .*` + regexp.QuoteMeta(sec.keyword))
		if !pattern.MatchString(content) {
			t.Errorf("meta-prompt manca sezione attesa: '## %d. ...%s...'", sec.num, sec.keyword)
		}
	}

	// FR-20 sub-d: lista kb_files con descrizione. Deve menzionare ognuno
	// dei file canonici della KB-ispettore.
	kbFiles := []string{
		"CLAUDE.md",
		"come-pensa-un-ispettore.md",
		"MAPPA_ESPLOSA_REQUISITI_ISO17025.md",
		"sezione-4-requisiti-generali-domande.md",
		"sezione-5-requisiti-strutturali-domande.md",
		"sezione-6-risorse-personale-dotazioni-domande.md",
		"sezione-7-processo-metodi-validazione-domande.md",
		"sezione-8-sgq-domande.md",
		"come-rispondere-NC-ACCREDIA.md",
		"errori-fatali-da-evitare.md",
	}
	for _, kb := range kbFiles {
		if !strings.Contains(content, kb) {
			t.Errorf("meta-prompt non cita kb_file canonico '%s' nella lista sezione 4 (FR-20 sub-d)", kb)
		}
	}

	// FR-20 sub-e: 2-3 esempi few-shot. Deve avere ≥ 2 blocchi yaml
	// fenced con campo "profilo:" diverso.
	reYamlBlock := regexp.MustCompile("(?s)```yaml\\s+([^`]+?)```")
	yamlBlocks := reYamlBlock.FindAllStringSubmatch(content, -1)
	if len(yamlBlocks) < 3 {
		// 1 schema completo + ≥ 2 esempi.
		t.Errorf("meta-prompt deve avere ≥ 3 blocchi ```yaml``` (schema + 2-3 esempi few-shot), trovati %d", len(yamlBlocks))
	}

	// FR-20 sub-g: pattern di output finale deve menzionare il comando
	// `labnexus validate`.
	if !strings.Contains(content, "labnexus validate") {
		t.Error("meta-prompt non istruisce Claude a fornire il comando 'labnexus validate' nell'output (FR-20 sub-g)")
	}
}

// TestGuidaMetaPromptContent — FR-21: la guida d'uso per Denis deve
// esistere ed essere autosufficiente (apri Claude → incolla → descrivi →
// salva YAML → valida).
func TestGuidaMetaPromptContent(t *testing.T) {
	wd, _ := os.Getwd()
	repoRoot := filepath.Dir(filepath.Dir(wd))
	guidaPath := filepath.Join(repoRoot, "docs", "guida-meta-prompt-denis.md")

	b, err := os.ReadFile(guidaPath)
	if err != nil {
		t.Fatalf("guida-meta-prompt-denis non trovata: %s — FR-21 deliverable mancante", guidaPath)
	}
	content := strings.ToLower(string(b))

	// FR-21: deve menzionare i passi chiave del workflow.
	required := []string{
		"claude.ai",
		"incolla",
		"profili/",
		"labnexus validate",
	}
	for _, kw := range required {
		if !strings.Contains(content, strings.ToLower(kw)) {
			t.Errorf("guida-meta-prompt-denis non menziona la keyword operativa '%s' (FR-21)", kw)
		}
	}
}

// itoa helper minimo per int → string usato nel test FR-20.
func itoa(i int) string {
	if i < 10 {
		return string(rune('0' + i))
	}
	return string(rune('0'+i/10)) + string(rune('0'+i%10))
}
