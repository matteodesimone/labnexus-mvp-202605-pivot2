// doc.go — Parser per Word 97-2003 OLE binary (.doc). FR-19.
//
// Strategia: usa github.com/richardlehane/mscfb per estrarre lo stream
// "WordDocument" dal container OLE, poi estrae sequenze ASCII/Latin1
// stampabili lunghe ≥ 4 caratteri come euristica di plain text. Non è
// un parser FIB completo, ma è sufficiente per il use case Denis (estrarre
// testo grezzo da scheda personale, mansionario, ecc.).
package input

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/richardlehane/mscfb"
)

// parseDoc estrae plain text da un file Word 97-2003 OLE binary. Niente
// marcatori di formattazione (font, colori, immagini). Su file corrotto
// (header non-OLE): error che il chiamante converte in warning+skip.
//
// Limitazione nota: il parser non gestisce file in formato Word 2007+
// (.docx — usa parser docx esistente). I caratteri non-ASCII (UTF-16 LE
// nel content stream) vengono decodificati come Latin1 best-effort.
func parseDoc(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("input: open .doc %q: %w", path, err)
	}
	defer f.Close()
	doc, err := mscfb.New(f)
	if err != nil {
		return "", fmt.Errorf("input: parse .doc %q (non OLE): %w", path, err)
	}
	stream, err := readWordDocumentStream(doc)
	if err != nil {
		return "", err
	}
	return extractPrintableRuns(stream), nil
}

// readWordDocumentStream cerca lo stream "WordDocument" nel container OLE
// e ritorna i bytes raw. Errore se non trovato.
func readWordDocumentStream(doc *mscfb.Reader) ([]byte, error) {
	for entry, err := doc.Next(); err == nil; entry, err = doc.Next() {
		if entry.Name != "WordDocument" {
			continue
		}
		buf, err := io.ReadAll(entry)
		if err != nil {
			return nil, fmt.Errorf("input: read WordDocument stream: %w", err)
		}
		return buf, nil
	}
	return nil, fmt.Errorf("input: stream 'WordDocument' non trovato (non è un .doc Word 97-2003)")
}

// extractPrintableRuns estrae sequenze di caratteri stampabili lunghe ≥ 4
// dal content stream di un .doc. Decodifica sia ASCII che UTF-16 LE (alternanza
// byte/zero tipica di Word 97-2003 testo non-OEM).
func extractPrintableRuns(data []byte) string {
	var out strings.Builder
	runs := append(extractAscii(data), extractUtf16LE(data)...)
	seen := make(map[string]bool)
	for _, r := range runs {
		r = strings.TrimSpace(r)
		if r == "" || seen[r] {
			continue
		}
		seen[r] = true
		out.WriteString(r)
		out.WriteByte('\n')
	}
	return strings.TrimRight(out.String(), "\n")
}

// extractAscii scan run ASCII stampabili (32-126 + tab/newline) lunghi ≥ 4.
func extractAscii(data []byte) []string {
	var runs []string
	var cur bytes.Buffer
	flush := func() {
		if cur.Len() >= 4 {
			runs = append(runs, cur.String())
		}
		cur.Reset()
	}
	for _, b := range data {
		if (b >= 32 && b <= 126) || b == '\t' || b == '\n' {
			cur.WriteByte(b)
		} else {
			flush()
		}
	}
	flush()
	return runs
}

// extractUtf16LE scan run UTF-16 LE (byte stampabile + 0x00 alternato) lunghi ≥ 4.
// Word 97-2003 scrive testo Unicode come UTF-16 little-endian nel content stream.
func extractUtf16LE(data []byte) []string {
	var runs []string
	var cur bytes.Buffer
	flush := func() {
		if cur.Len() >= 4 {
			runs = append(runs, cur.String())
		}
		cur.Reset()
	}
	for i := 0; i+1 < len(data); i += 2 {
		b, hi := data[i], data[i+1]
		if hi == 0 && ((b >= 32 && b <= 126) || b == '\t' || b == '\n') {
			cur.WriteByte(b)
		} else {
			flush()
		}
	}
	flush()
	return runs
}
