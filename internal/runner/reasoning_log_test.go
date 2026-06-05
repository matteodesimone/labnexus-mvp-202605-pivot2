package runner

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/provider"
	"github.com/labnexus/labnexus/internal/runlog"
)

// Il reasoning del modello va MOSTRATO e LOGGATO (bodyWriter) ma NON deve finire
// nel .md (body), che resta la risposta pulita. Niente del modello si butta.
func TestDrainStream_ReasoningLoggedNotInOutput(t *testing.T) {
	ch := make(chan provider.StreamEvent, 8)
	ch <- provider.StreamEvent{Reasoning: "STO-PENSANDO"}
	ch <- provider.StreamEvent{Reasoning: " ancora"}
	ch <- provider.StreamEvent{Token: "RISPOSTA-FINALE"}
	ch <- provider.StreamEvent{Done: true}
	close(ch)

	var shown bytes.Buffer
	body, stato, err := drainStreamWithBody(ch, runlog.New(io.Discard), &shown)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if stato != "completato" {
		t.Errorf("stato = %q, want completato", stato)
	}
	if body != "RISPOSTA-FINALE" {
		t.Errorf("body (.md) = %q — deve contenere solo il contenuto, NON il reasoning", body)
	}
	s := shown.String()
	for _, want := range []string{"STO-PENSANDO", "ancora", "RISPOSTA-FINALE"} {
		if !strings.Contains(s, want) {
			t.Errorf("mostrato/loggato deve contenere %q, got %q", want, s)
		}
	}
}
