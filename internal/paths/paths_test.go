package paths_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/paths"
)

// TestDeliveryRootForBinary_TableDriven cataloga i 4 pattern principali di
// distribuzione del binario LabNexus e verifica che la delivery root
// (= dove vivono profili/ e KB-ispettore/) sia risolta correttamente.
//
// Nota: i path sono trattati come stringhe — `DeliveryRootForBinary` non
// deve fare I/O sul filesystem (funzione pura, no symlink resolve, no Stat).
func TestDeliveryRootForBinary_TableDriven(t *testing.T) {
	cases := []struct {
		name      string
		binary    string
		wantRoot  string
	}{
		{
			// Bundle macOS .app: il binario reale è in Contents/MacOS/.
			// Profili e KB sono adiacenti al .app, quindi 4 livelli sopra.
			name:     "macOS .app bundle (labnexus-bin)",
			binary:   "/Users/denis/lab/labnexus.app/Contents/MacOS/labnexus-bin",
			wantRoot: "/Users/denis/lab",
		},
		{
			// Bundle in /Applications (path classico macOS installation).
			name:     "macOS .app in /Applications",
			binary:   "/Applications/labnexus.app/Contents/MacOS/labnexus-bin",
			wantRoot: "/Applications",
		},
		{
			// Binario standalone (make build → bin/labnexus): profili e KB
			// sono nella stessa dir del binario (= bin/).
			name:     "standalone binary in bin/ directory",
			binary:   "/Users/cto/projects/labnexus/bin/labnexus",
			wantRoot: "/Users/cto/projects/labnexus/bin",
		},
		{
			// Binario Linux (cross-compile output): profili e KB nella stessa dir.
			name:     "linux/amd64 standalone binary",
			binary:   "/home/cto/labnexus/dist/labnexus-linux-amd64",
			wantRoot: "/home/cto/labnexus/dist",
		},
		{
			// Edge case: bundle dentro un path con .app intermedi (es. Documents/Some.app.bak/).
			// Il pattern .app/Contents/MacOS/ deve matchare solo l'ultimo livello.
			name:     "bundle with .app-like intermediate path",
			binary:   "/Users/test/Documents/old.app.bak/labnexus.app/Contents/MacOS/labnexus-bin",
			wantRoot: "/Users/test/Documents/old.app.bak",
		},
		{
			// Edge case: una dir chiamata 'MacOS' che non è dentro un bundle.
			// Es. /Users/foo/MacOS/labnexus → non è un bundle → root = /Users/foo/MacOS/.
			name:     "directory named MacOS not inside a bundle",
			binary:   "/Users/foo/MacOS/labnexus",
			wantRoot: "/Users/foo/MacOS",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := paths.DeliveryRootForBinary(tc.binary)
			if got != tc.wantRoot {
				t.Errorf("DeliveryRootForBinary(%q)\n  got:  %q\n  want: %q", tc.binary, got, tc.wantRoot)
			}
		})
	}
}

// TestDeliveryRoot_UsesRealExecutable verifica che DeliveryRoot() chiami
// os.Executable() e produca un path PLAUSIBILE (esistente come dir).
// Non verifichiamo l'esatto path (dipende da dove gira il test) ma:
//   - non ritorna "." (= fallback patologico)
//   - ritorna una dir reale e accessibile
func TestDeliveryRoot_UsesRealExecutable(t *testing.T) {
	got := paths.DeliveryRoot()
	if got == "." {
		t.Fatalf("DeliveryRoot() ha ritornato \".\" → os.Executable è fallito o stub non implementato")
	}
	info, err := os.Stat(got)
	if err != nil {
		t.Errorf("DeliveryRoot() = %q non è accessibile: %v", got, err)
		return
	}
	if !info.IsDir() {
		t.Errorf("DeliveryRoot() = %q non è una directory", got)
	}
}

// TestDeliveryRootForBinary_RelativeBinaryPath verifica che la funzione gestisca
// path NON assoluti (per robustezza — non dovrebbe accadere in pratica perché
// os.Executable restituisce sempre absolute, ma vogliamo evitare panic).
func TestDeliveryRootForBinary_RelativeBinaryPath(t *testing.T) {
	got := paths.DeliveryRootForBinary("./bin/labnexus")
	want := filepath.Dir("./bin/labnexus") // "bin"
	if got != want {
		t.Errorf("DeliveryRootForBinary(./bin/labnexus) = %q, want %q", got, want)
	}
}

// TestDeliveryRootForBinary_DoesNotTouchFilesystem verifica che la funzione
// sia pura (non chiama Stat, EvalSymlinks, ecc.) lavorando su un path che
// SICURAMENTE non esiste — se la funzione facesse I/O, fallirebbe.
func TestDeliveryRootForBinary_DoesNotTouchFilesystem(t *testing.T) {
	fake := "/tmp/this-path-does-not-exist-12345/labnexus.app/Contents/MacOS/labnexus-bin"
	got := paths.DeliveryRootForBinary(fake)
	want := "/tmp/this-path-does-not-exist-12345"
	if got != want {
		t.Errorf("DeliveryRootForBinary su path inesistente non deve fare I/O: got %q, want %q", got, want)
	}
	if strings.Contains(got, "..") {
		t.Errorf("la funzione non deve normalizzare con `..`: got %q", got)
	}
}
