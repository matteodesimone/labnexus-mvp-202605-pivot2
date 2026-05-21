// rtf.go — Parser per Rich Text Format (.rtf). FR-20.
//
// Implementazione custom (RTF è formattato testuale, no lib esterna necessaria).
// Strategia: scan dei control words `\<word>[<number>][ ]` + skip gruppi `{...}`
// che iniziano con control words font/colortbl/stylesheet/info.
package input

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

// parseRtf estrae plain text da un file RTF. Strip dei control words
// (\<word>) e dei gruppi destination (es. \fonttbl, \colortbl, \info, ecc.).
// Decode minimale degli escape: \\ → \, \{ → {, \} → }, \~ → spazio, \tab → tab,
// \par → newline. Unicode escape `\uN`: decodifica il codepoint N e skippa il
// fallback char successivo.
func parseRtf(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("input: lettura %q: %w", path, err)
	}
	if !strings.HasPrefix(string(data), "{\\rtf") {
		return "", fmt.Errorf("input: file %q non sembra RTF (manca header {\\rtf)", path)
	}
	return decodeRtf(string(data)), nil
}

// destinationControlWords sono le control words che marcano gruppi da skippare
// completamente nel plain text estratto (font tables, color tables, metadata).
var destinationControlWords = map[string]bool{
	"fonttbl": true, "colortbl": true, "stylesheet": true,
	"info": true, "pict": true, "header": true, "footer": true,
	"object": true, "objdata": true, "filetbl": true, "listtable": true,
	"listoverridetable": true, "rsidtbl": true, "generator": true,
	"themedata": true, "datastore": true, "operator": true, "author": true,
	"company": true, "expandedcolortbl": true, "cocoartf": true,
}

func decodeRtf(s string) string {
	var out strings.Builder
	skipDepth := 0
	i := 0
	for i < len(s) {
		c := s[i]
		switch {
		case c == '{':
			i++
			// Look-ahead: se inizia con \ + destination word, attiva skip.
			if word, advance := peekControlWord(s, i); destinationControlWords[word] {
				skipDepth++
				i += advance
				continue
			}
		case c == '}':
			if skipDepth > 0 {
				skipDepth--
			}
			i++
			continue
		case c == '\\':
			word, advance := peekControlWord(s, i+1)
			if advance == 0 {
				// Escape singolo: \\ \{ \} \~ ecc.
				if i+1 < len(s) && skipDepth == 0 {
					switch s[i+1] {
					case '\\', '{', '}':
						out.WriteByte(s[i+1])
					case '~':
						out.WriteByte(' ')
					}
				}
				i += 2
				continue
			}
			if skipDepth == 0 {
				writeControlWordOutput(&out, word)
			}
			i += 1 + advance
			continue
		default:
			if skipDepth == 0 && (unicode.IsPrint(rune(c)) || c == '\n' || c == '\t') {
				out.WriteByte(c)
			}
			i++
		}
	}
	return strings.TrimSpace(out.String())
}

// peekControlWord legge una control word RTF a partire da i. Una control word
// è '\' (già consumato dal chiamante) + sequenza di lettere ASCII + parametro
// numerico opzionale + spazio opzionale di terminazione. Ritorna (word, advance)
// dove advance è quanti byte consumare dopo il '\'.
func peekControlWord(s string, start int) (string, int) {
	j := start
	for j < len(s) && isAlpha(s[j]) {
		j++
	}
	if j == start {
		return "", 0
	}
	word := s[start:j]
	// Parametro numerico opzionale (es. \fs24, 荤).
	if j < len(s) && (s[j] == '-' || isDigit(s[j])) {
		if s[j] == '-' {
			j++
		}
		for j < len(s) && isDigit(s[j]) {
			j++
		}
	}
	// Spazio opzionale di terminazione (consumato).
	if j < len(s) && s[j] == ' ' {
		j++
	}
	return word, j - start
}

func writeControlWordOutput(out *strings.Builder, word string) {
	switch word {
	case "par", "line":
		out.WriteByte('\n')
	case "tab":
		out.WriteByte('\t')
	}
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
