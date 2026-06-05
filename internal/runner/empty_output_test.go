package runner

import (
	"bytes"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/input"
	"github.com/labnexus/labnexus/internal/profile"
	"github.com/labnexus/labnexus/internal/runlog"
	"github.com/labnexus/labnexus/internal/tokens"
)

// #3 — risposta vuota del modello (es. kimi-k2.6 0 token): non è "completato",
// va segnalata con stato "vuoto", un .md esplicativo non vuoto, ed exit non-zero.
func TestFinalizeOutput_EmptyBody(t *testing.T) {
	p := &profile.Profile{Modello: "kimi-k2.6"}

	for _, raw := range []string{"", "   ", "\n\t  \n"} {
		body, stato, empty := finalizeOutput(raw, "completato", p, "eurouter")
		if !empty {
			t.Errorf("body %q deve essere rilevato come vuoto", raw)
		}
		if stato != "vuoto" {
			t.Errorf("stato = %q, want vuoto", stato)
		}
		if strings.TrimSpace(body) == "" {
			t.Error("il .md non deve essere vuoto: deve spiegare l'accaduto")
		}
		if !strings.Contains(body, "kimi-k2.6") || !strings.Contains(body, "eurouter") {
			t.Errorf("il messaggio deve citare modello e provider, got: %q", body)
		}
	}
}

func TestFinalizeOutput_NonEmptyUnchanged(t *testing.T) {
	p := &profile.Profile{Modello: "kimi-k2.6"}
	body, stato, empty := finalizeOutput("## 1. Sintesi\nContenuto reale.", "completato", p, "eurouter")
	if empty {
		t.Error("un body con contenuto non deve essere marcato vuoto")
	}
	if stato != "completato" || !strings.Contains(body, "Contenuto reale") {
		t.Errorf("body/stato alterati: stato=%q body=%q", stato, body)
	}
}

func TestExitCodeFor_VuotoIsNonZero(t *testing.T) {
	if exitCodeFor("vuoto") == 0 {
		t.Error("stato 'vuoto' deve dare exit non-zero (run fallita)")
	}
	if exitCodeFor("completato") != 0 {
		t.Error("stato 'completato' deve dare exit 0")
	}
}

// #4 — il .log di audit deve registrare COSA è stato inviato: KB, dati,
// esclusioni, file saltati. Prova auditabile contro "dati non inviati?".
func TestLogInputAudit_RecordsWhatWasSent(t *testing.T) {
	var buf bytes.Buffer
	log := runlog.New(&buf)
	parsed := &input.ParsedResult{
		Files: []input.ParsedFile{
			{Name: "M114.xlsx", Size: 20626},
			{Name: "M115.doc", Size: 731},
		},
		Skipped: []input.SkippedFile{
			{Name: "rotto.pdf", Reason: "PDF illeggibile"},
		},
	}
	check := &tokens.CheckResult{Tokens: 69485}
	p := &profile.Profile{Modello: "kimi-k2.6", ContextWindow: 262000}

	logInputAudit(log, []string{"kb1", "kb2"}, parsed, []string{"Prompt_INPUT.rtf"}, check, p)

	out := buf.String()
	for _, want := range []string{
		"M114.xlsx", "20626",
		"M115.doc",
		"KB system context: 2 file",
		"escluso dall'input", "Prompt_INPUT.rtf",
		"NON parsato", "rotto.pdf",
		"69485",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("il log di audit deve contenere %q\n--- log ---\n%s", want, out)
		}
	}
}
