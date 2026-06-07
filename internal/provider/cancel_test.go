package provider

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"
)

// Quando il runner cancella il context (idle-timeout su stream in stallo), la
// goroutine consumer NON deve restare appesa su `ch <-`: deve vedere ctx.Done()
// via emit e terminare, chiudendo canale+body. Senza la ctx-awareness ci sarebbe
// leak goroutine + connessione HTTP mai chiusa (finding CRITICAL code-review 07/06).
func TestConsumeEurouterStream_CancelUnblocksBlockedSend(t *testing.T) {
	// Molti più chunk del buffer del canale: la goroutine riempie il buffer e poi
	// si BLOCCA su `ch <-` perché il test non legge. cancel() deve sbloccarla.
	var sb strings.Builder
	for i := 0; i < 200; i++ {
		sb.WriteString(`data: {"choices":[{"delta":{"content":"x"}}]}` + "\n")
	}
	body := io.NopCloser(strings.NewReader(sb.String()))
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan StreamEvent, 1) // buffer minimo; il test NON legge
	done := make(chan struct{})
	go func() {
		consumeEurouterStream(ctx, body, nil, ch)
		close(done)
	}()
	time.Sleep(20 * time.Millisecond) // lascia che la goroutine si blocchi sul send
	cancel()
	select {
	case <-done:
		// ok: terminata dopo cancel (emit ha visto ctx.Done()).
	case <-time.After(2 * time.Second):
		t.Fatal("consumeEurouterStream non terminata dopo cancel: goroutine appesa su ch <- (leak)")
	}
}
