//go:build darwin

// Package bundle contains integration tests for the macOS .app bundle layout
// produced by scripts/build-mac.sh.
//
// Bug reproducer for `.pipeline/bugs/app-bundle-doppio-click-no-output.md`:
// quando l'utente fa doppio click sul bundle dal Finder, l'eseguibile deve
// essere un wrapper bash che apre Terminal (via osascript) e lancia il binario
// reale dentro la finestra. Senza wrapper il binario gira senza TTY,
// fa exit 2 silenzioso e l'utente vede "niente".
package bundle_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// repoRoot è la radice del repository (la dir parent di `internal/bundle/`).
const repoRoot = "../.."

// bundleDir è il path al bundle .app prodotto da build-mac.sh, relativo al repoRoot.
const bundleDir = "labnexus.app"

// TestMain costruisce il bundle UNA volta prima dei test (idempotente: lo script
// fa `rm -rf labnexus.app` all'inizio).
func TestMain(m *testing.M) {
	cmd := exec.Command("bash", "scripts/build-mac.sh")
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build-mac.sh failed: %v\n%s\n", err, out)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

// macOSBin ritorna il path assoluto a un file dentro labnexus.app/Contents/MacOS/.
func macOSBin(name string) string {
	return filepath.Join(repoRoot, bundleDir, "Contents/MacOS", name)
}

// TestBundle_CFBundleExecutableIsWrapperScript verifica che il CFBundleExecutable
// (= "labnexus") punti a un FILE DI TESTO (script bash con shebang), non al binario.
// Questo è il fix root del bug: il Finder lancia il wrapper, non il binario.
func TestBundle_CFBundleExecutableIsWrapperScript(t *testing.T) {
	wrapper := macOSBin("labnexus")
	info, err := os.Stat(wrapper)
	if err != nil {
		t.Fatalf("wrapper script must exist at %q: %v", wrapper, err)
	}
	// Uno script bash è piccolo (< 4 KB); il binario Go è ~8 MB.
	if info.Size() > 4*1024 {
		t.Errorf("wrapper deve essere uno script piccolo (< 4 KB), got %d bytes — probabilmente è il binario", info.Size())
	}
	data, err := os.ReadFile(wrapper)
	if err != nil {
		t.Fatalf("read wrapper: %v", err)
	}
	if !strings.HasPrefix(string(data), "#!") {
		t.Errorf("wrapper deve iniziare con shebang #!, got prefix: %q", string(data[:min(20, len(data))]))
	}
}

// TestBundle_RealBinaryAtLabnexusBin verifica che il binario Go reale viva
// come `labnexus-bin` (non `labnexus`, che è il wrapper).
func TestBundle_RealBinaryAtLabnexusBin(t *testing.T) {
	bin := macOSBin("labnexus-bin")
	info, err := os.Stat(bin)
	if err != nil {
		t.Fatalf("il binario reale deve esistere a %q come `labnexus-bin`: %v", bin, err)
	}
	// Binario Go è > 1 MB tipicamente.
	if info.Size() < 1024*1024 {
		t.Errorf("binario reale troppo piccolo (%d bytes), atteso Mach-O > 1 MB", info.Size())
	}
	out, err := exec.Command("file", bin).CombinedOutput()
	if err != nil {
		t.Fatalf("file(1) on %q: %v", bin, err)
	}
	if !strings.Contains(string(out), "Mach-O") || !strings.Contains(string(out), "arm64") {
		t.Errorf("atteso Mach-O arm64 per labnexus-bin, got: %s", out)
	}
}

// TestBundle_WrapperLaunchesTerminalViaOsascript verifica che il wrapper
// usi `osascript` per aprire Terminal e lanciare il binario reale.
// Senza questo, il doppio click esegue il binario senza TTY → ErrNonTTY → bug.
func TestBundle_WrapperLaunchesTerminalViaOsascript(t *testing.T) {
	wrapper := macOSBin("labnexus")
	data, err := os.ReadFile(wrapper)
	if err != nil {
		t.Fatalf("read wrapper: %v", err)
	}
	// Se il "wrapper" è in realtà un binario Mach-O (bug attuale), il primo
	// test ha già fallito; qui mostriamo solo dimensione e prefisso per evitare
	// di leakare ~8 MB di binary content nei log dei test.
	if !looksLikeText(data) {
		t.Errorf("wrapper sembra essere un binario, non uno script di testo (size %d bytes, prefix: %q)", len(data), safePrefix(data, 20))
		return
	}
	content := string(data)
	for _, expected := range []string{"osascript", "Terminal", "labnexus-bin"} {
		if !strings.Contains(content, expected) {
			t.Errorf("wrapper deve contenere %q (per aprire Terminal e lanciare il binario reale)", expected)
		}
	}
}

// looksLikeText ritorna true se data sembra un file di testo (no null bytes nei primi 512 byte).
func looksLikeText(data []byte) bool {
	probe := data
	if len(probe) > 512 {
		probe = probe[:512]
	}
	for _, b := range probe {
		if b == 0 {
			return false
		}
	}
	return true
}

// safePrefix ritorna fino a n bytes del prefisso, escapati in modo safe per log.
func safePrefix(data []byte, n int) string {
	if len(data) > n {
		data = data[:n]
	}
	return fmt.Sprintf("%q", data)
}

// TestBundle_WrapperIsExecutable verifica chmod +x sul wrapper (altrimenti
// macOS rifiuterebbe di eseguirlo come CFBundleExecutable).
func TestBundle_WrapperIsExecutable(t *testing.T) {
	wrapper := macOSBin("labnexus")
	info, err := os.Stat(wrapper)
	if err != nil {
		t.Fatalf("stat wrapper: %v", err)
	}
	mode := info.Mode().Perm()
	if mode&0o100 == 0 {
		t.Errorf("wrapper non eseguibile (perm: %o), Launch Services lo rifiuterà", mode)
	}
}

// TestBundle_InfoPlistDeclaresCFBundleExecutable verifica che Info.plist
// dichiari CFBundleExecutable = "labnexus" (= il wrapper, non il binario).
func TestBundle_InfoPlistDeclaresCFBundleExecutable(t *testing.T) {
	plist := filepath.Join(repoRoot, bundleDir, "Contents/Info.plist")
	data, err := os.ReadFile(plist)
	if err != nil {
		t.Fatalf("read Info.plist: %v", err)
	}
	if !strings.Contains(string(data), "<key>CFBundleExecutable</key>") {
		t.Errorf("Info.plist deve contenere CFBundleExecutable")
	}
	if !strings.Contains(string(data), "<string>labnexus</string>") {
		t.Errorf("CFBundleExecutable deve essere 'labnexus' (il wrapper), got plist:\n%s", data)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
