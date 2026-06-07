package input

import "testing"

// UTF-8 valido deve passare INVARIATO (nessun falso ri-decoding).
func TestDecodeText_UTF8Passthrough(t *testing.T) {
	in := []byte("già è perché — €uro ok")
	if got := decodeText(in); got != "già è perché — €uro ok" {
		t.Errorf("UTF-8 valido deve passare invariato, got %q", got)
	}
}

// MacRoman (export TextEdit/Mac): 0x88=à, 0x8F=è sono nella fascia C1 → decoder
// Macintosh. È il caso reale dell'input "ACIAA A1 rilievi.txt" (shakedown 07/06).
func TestDecodeText_MacRomanAccents(t *testing.T) {
	in := []byte{'q', 'u', 'a', 'l', 'i', 't', 0x88, ' ', 'p', 'e', 'r', 'c', 'h', 0x8F}
	if got := decodeText(in); got != "qualità perchè" {
		t.Errorf("MacRoman deve decodificare in UTF-8 corretto, got %q", got)
	}
}

// Windows-1252/Latin-1 (export Windows): 0xE0=à, 0xE8=è — NESSUN byte C1 →
// decoder Windows1252 (non MacRoman).
func TestDecodeText_Windows1252Accents(t *testing.T) {
	in := []byte{'c', 'i', 't', 't', 0xE0, ' ', 'p', 'e', 'r', 'c', 'h', 0xE8}
	if got := decodeText(in); got != "città perchè" {
		t.Errorf("Windows-1252 deve decodificare in UTF-8 corretto, got %q", got)
	}
}

// Plain ASCII: invariato.
func TestDecodeText_ASCII(t *testing.T) {
	if got := decodeText([]byte("plain ascii 123")); got != "plain ascii 123" {
		t.Errorf("ASCII invariato, got %q", got)
	}
}
