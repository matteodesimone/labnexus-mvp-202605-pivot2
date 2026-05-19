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
