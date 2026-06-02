package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsExistingDir_True(t *testing.T) {
	d := t.TempDir()
	if !isExistingDir(d) {
		t.Errorf("isExistingDir(%q) = false, want true", d)
	}
}

func TestIsExistingDir_FalseForFile(t *testing.T) {
	d := t.TempDir()
	file := filepath.Join(d, "x.txt")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if isExistingDir(file) {
		t.Errorf("isExistingDir(file) should be false")
	}
}

func TestIsExistingDir_FalseForEmpty(t *testing.T) {
	if isExistingDir("") {
		t.Error("isExistingDir(\"\") should be false")
	}
}

func TestIsExistingDir_FalseForNonExistent(t *testing.T) {
	if isExistingDir("/nonexistent/path/xyz/123") {
		t.Error("isExistingDir on nonexistent path should be false")
	}
}

func TestValidateExistingDir_RejectsMissing(t *testing.T) {
	err := validateExistingDir("/this/does/not/exist")
	if err == nil {
		t.Fatal("expected error for missing dir, got nil")
	}
}

func TestValidateExistingDir_AcceptsExisting(t *testing.T) {
	d := t.TempDir()
	if err := validateExistingDir(d); err != nil {
		t.Errorf("expected nil for existing dir, got: %v", err)
	}
}

func TestValidateNonEmpty_RejectsEmpty(t *testing.T) {
	if err := validateNonEmpty(""); err == nil {
		t.Fatal("expected error for empty string, got nil")
	}
	if err := validateNonEmpty("   "); err == nil {
		t.Fatal("expected error for whitespace-only, got nil")
	}
}

func TestValidateNonEmpty_CreatesMissingDirIdempotent(t *testing.T) {
	base := t.TempDir()
	newDir := filepath.Join(base, "subdir-nuovo")
	// Prima chiamata: deve crearla
	if err := validateNonEmpty(newDir); err != nil {
		t.Fatalf("first call: %v", err)
	}
	if fi, err := os.Stat(newDir); err != nil || !fi.IsDir() {
		t.Fatalf("dir not created: err=%v", err)
	}
	// Seconda chiamata: idempotente (la dir esiste già)
	if err := validateNonEmpty(newDir); err != nil {
		t.Errorf("second call should be idempotent, got: %v", err)
	}
}

func TestErrNonTTY_MessageMentionsAlternative(t *testing.T) {
	msg := ErrNonTTY.Error()
	if !strings.Contains(msg, "labnexus run") {
		t.Errorf("ErrNonTTY message should suggest 'labnexus run', got: %q", msg)
	}
	if !strings.Contains(msg, "--profile") {
		t.Errorf("ErrNonTTY message should mention --profile flag, got: %q", msg)
	}
}

// TestCleanPath verifica che il path inserito nella TUI venga normalizzato:
// rimosso whitespace + quote, gestendo i pattern tipici di drag&drop Finder.
func TestCleanPath_DragDropFinderPatterns(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"path con trailing space (caso bug CTO)", "/Users/denis/Test1 ", "/Users/denis/Test1"},
		{"path con single quotes + trailing space (drag&drop Finder)", "'/Users/denis/Test1' ", "/Users/denis/Test1"},
		{"path con leading + trailing space", "  /Users/denis/Test1  ", "/Users/denis/Test1"},
		{"path con double quotes", `"/Users/denis/Test1"`, "/Users/denis/Test1"},
		{"path pulito invariato", "/Users/denis/Test1", "/Users/denis/Test1"},
		{"backslash-escaped spaces (drag&drop Terminal)", `/Users/denis/CAPABILITY\ D\ —\ Profilo\ x`, "/Users/denis/CAPABILITY D — Profilo x"},
		{"backslash-escaped altri metacaratteri", `/Users/denis/a\&b\ \(c\)`, "/Users/denis/a&b (c)"},
		{"empty", "", ""},
		{"only whitespace", "   ", ""},
		{"only quotes", "''", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := cleanPath(tc.in)
			if got != tc.want {
				t.Errorf("cleanPath(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// TestValidateExistingDir_HandlesTrailingSpace verifica che validateExistingDir
// gestisca path con trailing space (caso bug CTO): la dir esiste, il path ha
// spazio in fondo (es. da drag&drop), validate deve passare.
func TestValidateExistingDir_HandlesTrailingSpace(t *testing.T) {
	d := t.TempDir()
	// Path con trailing space — pre-fix avrebbe fallito perché os.Stat su
	// "<d> " (con spazio) ritorna ENOENT.
	if err := validateExistingDir(d + " "); err != nil {
		t.Errorf("validateExistingDir deve gestire trailing space, got: %v", err)
	}
	if err := validateExistingDir("  " + d + "  "); err != nil {
		t.Errorf("validateExistingDir deve gestire leading+trailing space, got: %v", err)
	}
}

// TestValidateExistingDir_HandlesSingleQuotes verifica che validateExistingDir
// gestisca path drag&drop con single quotes (pattern Finder/Terminal macOS).
func TestValidateExistingDir_HandlesSingleQuotes(t *testing.T) {
	d := t.TempDir()
	quoted := "'" + d + "' "
	if err := validateExistingDir(quoted); err != nil {
		t.Errorf("validateExistingDir deve gestire single-quoted path da drag&drop, got: %v", err)
	}
}

// TestValidateExistingDir_HandlesBackslashEscapedSpaces riproduce il caso reale
// (Stefano 2026-06-02): trascinando nel campo TUI una cartella con spazi/em-dash,
// Terminal incolla il path con backslash-escape (`CAPABILITY\ D\ —\ Profilo\ x`).
// La TUI non è una shell: deve de-escapare prima di os.Stat, altrimenti cerca
// una cartella con i backslash letterali nel nome e fallisce con "non esiste".
func TestValidateExistingDir_HandlesBackslashEscapedSpaces(t *testing.T) {
	parent := t.TempDir()
	realName := "CAPABILITY D — Profilo audit-checklist"
	realDir := filepath.Join(parent, realName)
	if err := os.MkdirAll(realDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Path come lo incolla il drag&drop: spazi escapati con backslash.
	escaped := filepath.Join(parent, `CAPABILITY\ D\ —\ Profilo\ audit-checklist`)
	if err := validateExistingDir(escaped); err != nil {
		t.Errorf("validateExistingDir deve de-escapare il path da drag&drop, got: %v", err)
	}
}

// TestRun_PreInputWithTrailingSpaceSkipsInputPrompt verifica che il drag&drop
// di una cartella sull'icona .app (preInput con quote/spazio) sia riconosciuto
// come "cartella esistente" e la TUI salti la prompt di input.
// Pre-fix (M8 review codex): buildFormFields chiamava isExistingDir(preInput) raw,
// quindi un preInput="'/path/'" con quote falliva → la TUI riapriva la prompt
// vanificando il drag&drop.
// Test indiretto: chiamiamo cleanPath sul preInput "quoted" e verifichiamo
// che il risultato sia una dir esistente.
func TestCleanPath_PreInputDragDropOnAppIcon(t *testing.T) {
	d := t.TempDir()
	preInputQuoted := "'" + d + "' "
	cleaned := cleanPath(preInputQuoted)
	if cleaned != d {
		t.Errorf("cleanPath(%q) = %q, want %q", preInputQuoted, cleaned, d)
	}
	if !isExistingDir(cleaned) {
		t.Errorf("dopo cleanPath il preInput deve essere riconosciuto come dir esistente: %q", cleaned)
	}
}
