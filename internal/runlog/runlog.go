// Package runlog gestisce log per step + rilevamento TTY (FR-8, NFR-4).
// La progress bar visibile è una funzione opzionale chiamata dal runner durante streaming.
package runlog

import (
	"fmt"
	"io"
	"os"
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
