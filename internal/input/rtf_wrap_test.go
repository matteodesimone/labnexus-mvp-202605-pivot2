package input_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/input"
)

// Round-trip: il RTF generato da WrapPlainAsRTF deve essere accettato da
// ParseTriggerPromptFile (header {\rtf valido) e de-RTF'are al testo originale.
func TestWrapPlainAsRTF_Roundtrip(t *testing.T) {
	plain := "TASK: fai X\n\nINPUT FORNITI:\n- file1\n- file2\n\nVINCOLI:\n- niente { } speciali \\ qui"
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "p.rtf"), []byte(input.WrapPlainAsRTF(plain)), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := input.ParseTriggerPromptFile(dir, "p.rtf")
	if err != nil {
		t.Fatalf("ParseTriggerPromptFile: %v (l'RTF generato non è valido)", err)
	}
	for _, want := range []string{"TASK: fai X", "INPUT FORNITI:", "- file1", "- file2", "VINCOLI:", "{ } speciali \\ qui"} {
		if !strings.Contains(got, want) {
			t.Errorf("output de-RTF non contiene %q\n--- got ---\n%s", want, got)
		}
	}
}
