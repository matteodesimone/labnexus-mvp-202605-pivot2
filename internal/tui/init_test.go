package tui

import (
	"errors"
	"testing"
)

// In contesto non-TTY (CI, pipe) il wizard init deve fallire pulito con ErrNonTTY
// PRIMA di qualsiasi I/O, non bloccarsi su un form interattivo. La parte
// interattiva (scelta profilo/prompt, conferma sovrascrittura) richiede un TTY
// reale ed è verificata manualmente.
func TestRunInit_NonTTYReturnsErrNonTTY(t *testing.T) {
	_, err := RunInit("./profili-inesistente", false)
	if !errors.Is(err, ErrNonTTY) {
		t.Errorf("atteso ErrNonTTY in non-TTY, got %v", err)
	}
	_, err = RunInit("./profili-inesistente", true) // anche con conferma sovrascrittura
	if !errors.Is(err, ErrNonTTY) {
		t.Errorf("atteso ErrNonTTY in non-TTY (alreadyExists), got %v", err)
	}
}
