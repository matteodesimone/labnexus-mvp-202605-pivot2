// Package runlog gestisce log per step + rilevamento TTY (FR-8, NFR-4).
// La progress bar visibile è una funzione opzionale chiamata dal runner durante streaming.
package runlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/term"
)

// Step è la durata di una fase della pipeline.
type Step struct {
	Name     string
	Started  time.Time
	Duration time.Duration
}

// Logger emette log per step su stderr; nessuna ANSI se non TTY (NFR-4).
type Logger struct {
	Out     io.Writer
	IsTTY   bool
	current *Step
	Steps   []Step
}

// New crea un logger per out (di norma os.Stderr). Auto-detect TTY su file fd.
func New(out io.Writer) *Logger {
	tty := false
	if f, ok := out.(*os.File); ok {
		tty = term.IsTerminal(int(f.Fd()))
	}
	return &Logger{Out: out, IsTTY: tty}
}

// BeginStep registra l'inizio di uno step.
func (l *Logger) BeginStep(name string) {
	l.current = &Step{Name: name, Started: time.Now()}
	if l.Out != nil {
		fmt.Fprintf(l.Out, "[%s] %s ...\n", l.current.Started.Format("15:04:05"), name)
	}
}

// EndStep registra la fine dello step corrente e ne stampa la durata in ms.
func (l *Logger) EndStep() {
	if l.current == nil {
		return
	}
	l.current.Duration = time.Since(l.current.Started)
	if l.Out != nil {
		fmt.Fprintf(l.Out, "[%s] %s ok (%d ms)\n",
			time.Now().Format("15:04:05"),
			l.current.Name,
			l.current.Duration.Milliseconds())
	}
	l.Steps = append(l.Steps, *l.current)
	l.current = nil
}

// Info logga una linea informativa generica.
func (l *Logger) Info(format string, args ...interface{}) {
	if l.Out == nil {
		return
	}
	fmt.Fprintf(l.Out, "[%s] %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, args...))
}

// Warn logga un warning, prefisso "warning:".
func (l *Logger) Warn(format string, args ...interface{}) {
	if l.Out == nil {
		return
	}
	fmt.Fprintf(l.Out, "[%s] warning: %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, args...))
}

// Progress aggiorna una linea di progress (no-op se non TTY).
func (l *Logger) Progress(current, total int) {
	if !l.IsTTY || l.Out == nil {
		return
	}
	pct := 0
	if total > 0 {
		pct = current * 100 / total
	}
	bar := buildBar(pct, 30)
	fmt.Fprintf(l.Out, "\r%s %3d%%", bar, pct)
}

// EndProgress finalizza la progress (newline, rimuove la riga).
func (l *Logger) EndProgress() {
	if !l.IsTTY || l.Out == nil {
		return
	}
	fmt.Fprint(l.Out, "\r"+strings.Repeat(" ", 60)+"\r")
}

// StreamProgress aggiorna la riga live di progress streaming LLM (TTY only).
// Formato: "  > 142 token, 23s elapsed (6.1 tok/s)". Su non-TTY è no-op:
// l'aggiornamento periodico per non-TTY passa per Info (vedi drainStream).
func (l *Logger) StreamProgress(tokens int, elapsed time.Duration) {
	if !l.IsTTY || l.Out == nil {
		return
	}
	rate := 0.0
	if elapsed > 0 {
		rate = float64(tokens) / elapsed.Seconds()
	}
	fmt.Fprintf(l.Out, "\r  > %d token, %s elapsed (%.1f tok/s)     ",
		tokens, elapsed.Round(time.Second), rate)
}

// StreamEnd finalizza il progress streaming. Su TTY rimuove la riga in-place,
// poi stampa una linea di riepilogo "[hh:mm:ss] streaming completato: N token in Xs (R tok/s)".
func (l *Logger) StreamEnd(tokens int, elapsed time.Duration) {
	if l.Out == nil {
		return
	}
	if l.IsTTY {
		fmt.Fprint(l.Out, "\r"+strings.Repeat(" ", 70)+"\r")
	}
	rate := 0.0
	if elapsed > 0 {
		rate = float64(tokens) / elapsed.Seconds()
	}
	fmt.Fprintf(l.Out, "[%s] streaming completato: %d token in %s (%.1f tok/s)\n",
		time.Now().Format("15:04:05"), tokens, elapsed.Round(time.Second), rate)
}

func buildBar(pct, width int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	filled := pct * width / 100
	return "[" + strings.Repeat("=", filled) + strings.Repeat(" ", width-filled) + "]"
}

// --- Sprint 1.5.C: audit trail multi-writer (FR-32, NFR-11) ---

// NewMulti crea un Logger che scrive su N writers (es. stderr + log file accoppiato).
// Pattern io.MultiWriter stdlib. Audit trail by-product: lo stesso flusso log
// va in console (live) e in file (post-fatto, ISO 17025 compliance).
func NewMulti(writers ...io.Writer) *Logger {
	if len(writers) == 0 {
		return New(nil)
	}
	if len(writers) == 1 {
		return New(writers[0])
	}
	return New(io.MultiWriter(writers...))
}

// OpenLogFile crea il file di audit log accoppiato all'output MD (FR-32, NFR-11).
// Pattern: nello stesso outputDir, con baseName matching il `.md`, estensione `.log`.
// Caller chiude il *os.File. Crea outputDir se mancante.
//
// Sprint 1.5.C: il `.log` rappresenta l'audit trail ISO 17025 — ogni run produce
// un MD (bozza per Denis) accoppiato a un LOG (trail tecnico per CTO/audit).
func OpenLogFile(outputDir, baseName string) (*os.File, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("runlog: mkdir %q: %w", outputDir, err)
	}
	path := filepath.Join(outputDir, baseName+".log")
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("runlog: create %q: %w", path, err)
	}
	return f, nil
}
