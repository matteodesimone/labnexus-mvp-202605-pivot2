package runner

import (
	"bytes"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/prompt"
)

// Bugfix #005 (.pipeline/bugs/show-prompt-pii-exposure.md): in modalità
// non-TTY, `labnexus check --show-prompt` deve emettere un warning PII su
// stderr PRIMA del prompt composto, per proteggere contro copia/paste
// accidentale di dati personali del SGQ.

func TestPrintPromptWithPIIWarning_EmitsWarningOnNonTTY(t *testing.T) {
	var stdout, stderr bytes.Buffer
	composed := &prompt.Composed{
		System: "system msg",
		User:   "user msg with kb content",
	}
	printPromptWithPIIWarning(composed, &stdout, &stderr, false)

	stderrStr := stderr.String()
	lowerErr := strings.ToLower(stderrStr)
	if !strings.Contains(lowerErr, "pii") && !strings.Contains(lowerErr, "dati personali") {
		t.Fatalf("warning PII atteso su stderr in non-TTY, got stderr:\n%s", stderrStr)
	}
	if !strings.Contains(lowerErr, "condivider") && !strings.Contains(lowerErr, "do not share") {
		t.Fatalf("warning deve sconsigliare la condivisione, got stderr:\n%s", stderrStr)
	}
	// Il prompt vero deve essere su stdout, NON su stderr.
	if !strings.Contains(stdout.String(), "system msg") {
		t.Fatalf("stdout deve contenere il system message, got:\n%s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "user msg with kb content") {
		t.Fatalf("stdout deve contenere il user message, got:\n%s", stdout.String())
	}
	// Il warning NON deve essere su stdout (sennò sporca il prompt pipato).
	if strings.Contains(strings.ToLower(stdout.String()), "pii") ||
		strings.Contains(strings.ToLower(stdout.String()), "dati personali") {
		t.Fatalf("warning NON deve essere su stdout, got stdout:\n%s", stdout.String())
	}
}

func TestPrintPromptWithPIIWarning_NoWarningOnTTY(t *testing.T) {
	var stdout, stderr bytes.Buffer
	composed := &prompt.Composed{
		System: "system msg",
		User:   "user msg",
	}
	printPromptWithPIIWarning(composed, &stdout, &stderr, true)

	// In TTY mode, nessun warning (l'utente sta vedendo a video).
	if stderr.Len() != 0 {
		t.Fatalf("nessun warning atteso in TTY mode, got stderr:\n%s", stderr.String())
	}
	if !strings.Contains(stdout.String(), "system msg") {
		t.Fatalf("stdout deve contenere il prompt, got:\n%s", stdout.String())
	}
}

func TestPrintPromptWithPIIWarning_PromptOrderingNonTTY(t *testing.T) {
	// In non-TTY: il warning su stderr deve precedere logicamente il prompt
	// su stdout (qui verifichiamo presenza, l'ordine è separato per
	// destination — stderr/stdout — quindi è un'evidenza forte ma non
	// strettamente "warning prima del prompt nella stessa linea").
	var stdout, stderr bytes.Buffer
	composed := &prompt.Composed{System: "sys", User: "usr"}
	printPromptWithPIIWarning(composed, &stdout, &stderr, false)
	if stderr.Len() == 0 {
		t.Fatal("stderr non deve essere vuoto in non-TTY")
	}
	if stdout.Len() == 0 {
		t.Fatal("stdout non deve essere vuoto (deve contenere il prompt)")
	}
}
