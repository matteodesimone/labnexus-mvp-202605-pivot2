package input

import (
	"strings"
	"testing"
)

// decodeRtf su RTF reale (cocoa/TextEdit, come i Prompt_INPUT_Rev.00.rtf di
// Denis) deve: (1) saltare font/color table, (2) decodificare gli escape
// esadecimali cp1252 \'XX (À, §, à...), (3) trasformare `\`+newline in a-capo.
// Prima del fix produceva "ATTIVITc0", "a74.2.1", run-on, e "Helvetica;" in testa.
func TestDecodeRtf_CocoaRealWorld(t *testing.T) {
	rtf := "{\\rtf1\\ansi\\ansicpg1252\\cocoartf2869\n" +
		"{\\fonttbl\\f0\\fswiss\\fcharset0 Helvetica;}\n" +
		"{\\colortbl;\\red255\\green255\\blue255;}\n" +
		"{\\*\\expandedcolortbl;;}\n" +
		"\\f0\\fs24 \\cf0 ATTIVIT\\'c0: confronta.\\\n" +
		"\\\n" +
		"Cita \\'a74.2.1 e criticit\\'e0 perch\\'e9 conta.}"

	got := decodeRtf(rtf)

	// (2) accenti decodificati, non hex grezzi
	for _, want := range []string{"ATTIVITÀ:", "§4.2.1", "criticità", "perché"} {
		if !strings.Contains(got, want) {
			t.Errorf("decode deve contenere %q\n--- got ---\n%q", want, got)
		}
	}
	for _, bad := range []string{"ATTIVITc0", "a74.2.1", "criticite0", "perche9"} {
		if strings.Contains(got, bad) {
			t.Errorf("decode NON deve contenere l'hex grezzo %q\n--- got ---\n%q", bad, got)
		}
	}
	// (1) niente spazzatura da font/color table
	for _, junk := range []string{"Helvetica", "expandedcolortbl"} {
		if strings.Contains(got, junk) {
			t.Errorf("decode NON deve contenere spazzatura %q\n--- got ---\n%q", junk, got)
		}
	}
	// (3) a-capo preservati: "confronta." e "Cita" su righe diverse
	if !strings.Contains(got, "confronta.\n") {
		t.Errorf("manca l'a-capo dopo 'confronta.'\n--- got ---\n%q", got)
	}
}
