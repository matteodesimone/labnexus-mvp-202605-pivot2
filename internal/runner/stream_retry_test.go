package runner

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labnexus/labnexus/internal/profile"
	"github.com/labnexus/labnexus/internal/prompt"
	"github.com/labnexus/labnexus/internal/provider"
	"github.com/labnexus/labnexus/internal/runlog"
)

func testProfile() *profile.Profile {
	return &profile.Profile{Modello: "kimi-k2.6", Temperature: 1.0, ContextWindow: 256000}
}

// scriptedProvider ritorna un canale diverso a ogni chiamata di Stream, secondo
// la funzione fn (call = numero progressivo della chiamata). Serve a testare i
// retry: 1ª chiamata fallisce/vuota, 2ª riesce.
type scriptedProvider struct {
	name  string
	calls int
	fn    func(call int, ctx context.Context) (<-chan provider.StreamEvent, error)
}

func (s *scriptedProvider) Name() string { return s.name }

func (s *scriptedProvider) Stream(ctx context.Context, _, _ string, _ provider.Options) (<-chan provider.StreamEvent, error) {
	s.calls++
	return s.fn(s.calls, ctx)
}

// stallChan non emette mai eventi (mima un gateway in stallo). Si chiude quando
// ctx è cancellato, replicando il comportamento ctx-aware dei provider reali, così
// il test non lascia goroutine appese dopo l'idle-timeout.
func stallChan(ctx context.Context) <-chan provider.StreamEvent {
	ch := make(chan provider.StreamEvent)
	go func() {
		<-ctx.Done()
		close(ch)
	}()
	return ch
}

func tokenChan(tokens ...string) <-chan provider.StreamEvent {
	ch := make(chan provider.StreamEvent, len(tokens)+1)
	for _, tk := range tokens {
		ch <- provider.StreamEvent{Token: tk}
	}
	ch <- provider.StreamEvent{Done: true}
	close(ch)
	return ch
}

// emptyDoneChan chiude lo stream con done:true ma ZERO token (il caso del modello
// thinking che esaurisce il budget nel reasoning → output vuoto).
func emptyDoneChan() <-chan provider.StreamEvent {
	ch := make(chan provider.StreamEvent, 1)
	ch <- provider.StreamEvent{Done: true}
	close(ch)
	return ch
}

// Idle-timeout: nessun evento entro il timeout → errStreamIdle (anti-hang).
func TestDrainStreamWithBody_IdleTimeoutStall(t *testing.T) {
	ch := make(chan provider.StreamEvent) // mai niente, mai chiuso: stream in stallo
	_, _, err := drainStreamWithBody(ch, runlog.New(io.Discard), io.Discard, io.Discard, true, 30*time.Millisecond)
	if !errors.Is(err, errStreamIdle) {
		t.Fatalf("atteso errStreamIdle su stream in stallo, got %v", err)
	}
}

// Il timer di idle si resetta a ogni evento: un modello che ragiona (reasoning)
// e poi risponde NON deve far scattare il timeout.
func TestDrainStreamWithBody_ResetsIdleOnEvents(t *testing.T) {
	ch := make(chan provider.StreamEvent, 4)
	ch <- provider.StreamEvent{Reasoning: "sto pensando"}
	ch <- provider.StreamEvent{Token: "out"}
	ch <- provider.StreamEvent{Done: true}
	close(ch)
	body, stato, err := drainStreamWithBody(ch, runlog.New(io.Discard), io.Discard, io.Discard, true, 200*time.Millisecond)
	if err != nil || stato != "completato" || body != "out" {
		t.Fatalf("eventi entro l'idle window → nessun timeout; body=%q stato=%q err=%v", body, stato, err)
	}
}

// Stallo al 1° tentativo → idle-timeout → retry → successo.
func TestStreamWithRetries_RetriesOnIdleStall(t *testing.T) {
	old := streamIdleTimeout
	streamIdleTimeout = 30 * time.Millisecond
	defer func() { streamIdleTimeout = old }()

	prov := &scriptedProvider{name: "fake", fn: func(call int, ctx context.Context) (<-chan provider.StreamEvent, error) {
		if call == 1 {
			return stallChan(ctx), nil
		}
		return tokenChan("OK"), nil
	}}
	body, stato, err := streamWithRetries(prov, &prompt.Composed{}, testProfile(), 1000, runlog.New(io.Discard), io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("dopo il retry il run deve riuscire, got err: %v", err)
	}
	if body != "OK" || stato != "completato" {
		t.Errorf("body=%q stato=%q, want OK/completato", body, stato)
	}
	if prov.calls != 2 {
		t.Errorf("calls=%d, want 2 (1 stallo + 1 retry riuscito)", prov.calls)
	}
}

// Output vuoto al 1° tentativo → retry → successo (è la leva quando il prompt non
// si può ridurre).
func TestStreamWithRetries_RetriesOnEmptyOutput(t *testing.T) {
	prov := &scriptedProvider{name: "fake", fn: func(call int, _ context.Context) (<-chan provider.StreamEvent, error) {
		if call == 1 {
			return emptyDoneChan(), nil
		}
		return tokenChan("RISPOSTA"), nil
	}}
	body, stato, err := streamWithRetries(prov, &prompt.Composed{}, testProfile(), 1000, runlog.New(io.Discard), io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if body != "RISPOSTA" || stato != "completato" {
		t.Errorf("body=%q stato=%q, want RISPOSTA/completato", body, stato)
	}
	if prov.calls != 2 {
		t.Errorf("calls=%d, want 2 (1 vuoto + 1 retry riuscito)", prov.calls)
	}
}

// Vuoto persistente: esauriti i retry transitori, ritorna body vuoto SENZA errore
// (finalizeOutput lo marcherà 'vuoto'); il numero di chiamate è limitato.
func TestStreamWithRetries_EmptyAfterAllRetries(t *testing.T) {
	prov := &scriptedProvider{name: "fake", fn: func(call int, _ context.Context) (<-chan provider.StreamEvent, error) {
		return emptyDoneChan(), nil
	}}
	body, _, err := streamWithRetries(prov, &prompt.Composed{}, testProfile(), 1000, runlog.New(io.Discard), io.Discard, io.Discard)
	if err != nil {
		t.Fatalf("vuoto persistente non è un errore (lo gestisce finalizeOutput): %v", err)
	}
	if strings.TrimSpace(body) != "" {
		t.Errorf("body deve restare vuoto, got %q", body)
	}
	if prov.calls != 1+maxTransientRetries {
		t.Errorf("calls=%d, want %d (1 iniziale + %d retry transitori)", prov.calls, 1+maxTransientRetries, maxTransientRetries)
	}
}

// Successo al primo colpo: nessun retry.
func TestStreamWithRetries_SuccessFirstTry(t *testing.T) {
	prov := &scriptedProvider{name: "fake", fn: func(call int, _ context.Context) (<-chan provider.StreamEvent, error) {
		return tokenChan("SUBITO"), nil
	}}
	body, _, err := streamWithRetries(prov, &prompt.Composed{}, testProfile(), 1000, runlog.New(io.Discard), io.Discard, io.Discard)
	if err != nil || body != "SUBITO" {
		t.Fatalf("body=%q err=%v, want SUBITO/nil", body, err)
	}
	if prov.calls != 1 {
		t.Errorf("calls=%d, want 1 (nessun retry su successo)", prov.calls)
	}
}

// Il canale reasoning viene SEMPRE salvato in un .md affiancato quando presente.
func TestWriteReasoningSidecar_WritesFile(t *testing.T) {
	dir := t.TempDir()
	path := writeReasoningSidecar(dir, "2026-06-07_revisione_MQL", "Table columns:\nRows:\n1. RT-23 Scope", runlog.New(io.Discard))
	if path == "" {
		t.Fatal("doveva scrivere il file di ragionamento")
	}
	if !strings.HasSuffix(path, "_ragionamento.md") {
		t.Errorf("nome file inatteso: %s", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if !strings.Contains(s, "RT-23 Scope") || !strings.Contains(s, "Ragionamento del modello") {
		t.Errorf("il sidecar deve contenere reasoning + intestazione, got:\n%s", s)
	}
}

// Niente reasoning (o outputDir vuoto) → nessun file.
func TestWriteReasoningSidecar_SkipsWhenEmpty(t *testing.T) {
	if p := writeReasoningSidecar(t.TempDir(), "base", "   \n  ", runlog.New(io.Discard)); p != "" {
		t.Errorf("reasoning vuoto → niente file, got %q", p)
	}
	if p := writeReasoningSidecar("", "base", "del reasoning c'è", runlog.New(io.Discard)); p != "" {
		t.Errorf("outputDir vuoto → niente file, got %q", p)
	}
}
