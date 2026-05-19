package runlog_test

import (
	"bytes"
	"strings"
	"testing"

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
