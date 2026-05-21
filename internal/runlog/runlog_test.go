package runlog_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/labnexus/labnexus/internal/runlog"
)

func TestLogger_BeginEndStepEmitsTimingLine(t *testing.T) {
	var buf bytes.Buffer
	l := runlog.New(&buf)
	l.BeginStep("step uno")
	l.EndStep()
	out := buf.String()
	if !strings.Contains(out, "step uno") {
		t.Errorf("output should contain step name, got: %q", out)
	}
	if !strings.Contains(out, "ms") {
		t.Errorf("output should contain duration in ms, got: %q", out)
	}
	if !strings.Contains(out, "ok") {
		t.Errorf("output should contain 'ok' completion marker, got: %q", out)
	}
}

func TestLogger_WarnPrefixesWithWarning(t *testing.T) {
	var buf bytes.Buffer
	l := runlog.New(&buf)
	l.Warn("attenzione %s", "qualcosa")
	out := buf.String()
	if !strings.Contains(out, "warning:") {
		t.Errorf("Warn output should start with 'warning:', got: %q", out)
	}
	if !strings.Contains(out, "attenzione qualcosa") {
		t.Errorf("Warn should include formatted message, got: %q", out)
	}
}

func TestLogger_StreamProgress_NonTTYIsNoOp(t *testing.T) {
	var buf bytes.Buffer
	l := &runlog.Logger{Out: &buf, IsTTY: false}
	l.StreamProgress(142, 23*time.Second)
	if buf.Len() > 0 {
		t.Errorf("StreamProgress su non-TTY deve essere no-op, got: %q", buf.String())
	}
}

func TestLogger_StreamProgress_TTYRendersTokenCountAndRate(t *testing.T) {
	var buf bytes.Buffer
	l := &runlog.Logger{Out: &buf, IsTTY: true}
	l.StreamProgress(142, 23*time.Second)
	s := buf.String()
	if !strings.HasPrefix(s, "\r") {
		t.Errorf("output deve iniziare con \\r (update in-place), got: %q", s)
	}
	if !strings.Contains(s, "142 token") {
		t.Errorf("output deve contenere '142 token', got: %q", s)
	}
	if !strings.Contains(s, "23s") {
		t.Errorf("output deve contenere '23s' (durata round), got: %q", s)
	}
	if !strings.Contains(s, "tok/s") {
		t.Errorf("output deve contenere 'tok/s' (rate), got: %q", s)
	}
}

func TestLogger_StreamEnd_AlwaysEmitsFinalStat(t *testing.T) {
	for _, tty := range []bool{true, false} {
		var buf bytes.Buffer
		l := &runlog.Logger{Out: &buf, IsTTY: tty}
		l.StreamEnd(500, 60*time.Second)
		s := buf.String()
		if !strings.Contains(s, "500 token") {
			t.Errorf("tty=%v: output deve contenere '500 token', got: %q", tty, s)
		}
		if !strings.Contains(s, "completato") {
			t.Errorf("tty=%v: output deve contenere 'completato', got: %q", tty, s)
		}
		if !strings.Contains(s, "tok/s") {
			t.Errorf("tty=%v: output deve contenere 'tok/s', got: %q", tty, s)
		}
		if tty && !strings.HasPrefix(s, "\r") {
			t.Errorf("tty=true: output deve iniziare con \\r (clear riga), got: %q", s)
		}
	}
}

func TestLogger_NonTTYProgressIsNoOp(t *testing.T) {
	// bytes.Buffer non è un *os.File → IsTerminal=false → Progress no-op.
	var buf bytes.Buffer
	l := runlog.New(&buf)
	if l.IsTTY {
		t.Fatal("bytes.Buffer should not be detected as TTY")
	}
	l.Progress(50, 100)
	if buf.Len() != 0 {
		t.Errorf("Progress in non-TTY must not write anything, got: %q", buf.String())
	}
}

func TestLogger_EndStepWithoutBeginIsSafe(t *testing.T) {
	// EndStep senza BeginStep precedente non deve panicare.
	var buf bytes.Buffer
	l := runlog.New(&buf)
	l.EndStep() // no-op
	if buf.Len() != 0 {
		t.Errorf("EndStep without BeginStep should be no-op, got: %q", buf.String())
	}
}

func TestLogger_StepsAreRecorded(t *testing.T) {
	var buf bytes.Buffer
	l := runlog.New(&buf)
	l.BeginStep("a")
	l.EndStep()
	l.BeginStep("b")
	l.EndStep()
	if len(l.Steps) != 2 {
		t.Errorf("expected 2 steps recorded, got %d", len(l.Steps))
	}
	names := []string{l.Steps[0].Name, l.Steps[1].Name}
	if names[0] != "a" || names[1] != "b" {
		t.Errorf("expected step names [a,b], got %v", names)
	}
}

// --- Sprint 1.5.C test-scaffold (red phase): test stub per audit trail ---

// TestNewMulti_WritesToAllSinks (red phase): atteso che NewMulti scriva su tutti
// i writer passati. Stub corrente usa solo il primo → test fail su secondo writer.
func TestNewMulti_WritesToAllSinks(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	l := runlog.NewMulti(&buf1, &buf2)
	l.Info("test message")
	if !strings.Contains(buf1.String(), "test message") {
		t.Errorf("buf1 dovrebbe contenere 'test message', got %q", buf1.String())
	}
	if !strings.Contains(buf2.String(), "test message") {
		t.Errorf("buf2 dovrebbe contenere 'test message' (multi-writer), got %q", buf2.String())
	}
}

// TestOpenLogFile_CreatesFileWithExpectedName (red phase): stub ritorna error.
func TestOpenLogFile_CreatesFileWithExpectedName(t *testing.T) {
	dir := t.TempDir()
	baseName := "2026-05-22T100000_revisione_Test1"
	f, err := runlog.OpenLogFile(dir, baseName)
	if err != nil {
		t.Fatalf("OpenLogFile: unexpected error %v", err)
	}
	defer f.Close()
	if f == nil {
		t.Fatal("OpenLogFile ritorna *os.File nil")
	}
	// Verifica nome file
	if !strings.HasSuffix(f.Name(), baseName+".log") {
		t.Errorf("nome file atteso *.log, got %q", f.Name())
	}
	// Scrivi qualcosa e verifica
	if _, err := f.WriteString("test log entry\n"); err != nil {
		t.Fatalf("write log: %v", err)
	}
}
