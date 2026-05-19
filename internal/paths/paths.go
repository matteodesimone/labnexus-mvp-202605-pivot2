// Package paths fornisce la risoluzione dei path di default del deliverable
// LabNexus (profili/, KB-ispettore/) **relativa alla posizione del binario**,
// NON al cwd del processo (spec FR-2).
//
// Bug fix per `.pipeline/bugs/path-resolution-profili-kb-cwd-relative.md`:
// i default `./profili` e `./KB-ispettore` di cobra erano risolti vs cwd
// (= HOME quando l'utente lancia il binario da Terminal), rendendo il
// binario inusabile fuori dalla repo del CTO.
package paths

import (
	"os"
	"path/filepath"
)

// DeliveryRootForBinary calcola la "delivery root" (dove vivono profili/ e
// KB-ispettore/) dato un path al binario. Funzione pura, NO I/O.
//
// Pattern di risoluzione:
//   - Se il binario è in `<X>.app/Contents/MacOS/<bin>` (bundle macOS),
//     ritorna la directory CONTENENTE `<X>.app`.
//   - Altrimenti, ritorna la directory del binario stesso.
//
// Esempi:
//
//	/path/labnexus.app/Contents/MacOS/labnexus-bin → /path
//	/path/bin/labnexus                              → /path/bin
//	/path/labnexus-linux-amd64                      → /path
//	/Applications/labnexus.app/Contents/MacOS/labnexus-bin → /Applications
func DeliveryRootForBinary(binaryPath string) string {
	dir := filepath.Dir(binaryPath)
	// Pattern bundle macOS: il binario è in `<X>.app/Contents/MacOS/<bin>`.
	// Vogliamo la directory parent di `<X>.app`.
	//   dir              = .../<X>.app/Contents/MacOS
	//   parent(dir)      = .../<X>.app/Contents
	//   parent(parent)   = .../<X>.app
	//   parent(parent(parent)) = ...  (delivery root)
	if filepath.Base(dir) == "MacOS" {
		contentsDir := filepath.Dir(dir)
		if filepath.Base(contentsDir) == "Contents" {
			appDir := filepath.Dir(contentsDir)
			if filepath.Ext(appDir) == ".app" {
				return filepath.Dir(appDir)
			}
		}
	}
	return dir
}

// DeliveryRoot ottiene il path del binario in esecuzione (via os.Executable
// + resolve dei symlink) e ne deriva la delivery root tramite DeliveryRootForBinary.
//
// Fallback: se os.Executable fallisce (caso patologico), ritorna ".".
func DeliveryRoot() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	// Risolve symlink (es. /usr/local/bin/labnexus → /usr/local/Cellar/...).
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	return DeliveryRootForBinary(exe)
}
