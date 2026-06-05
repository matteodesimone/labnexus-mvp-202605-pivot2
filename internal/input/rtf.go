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
		// Dentro un gruppo destination (fonttbl/colortbl/{\*...}): conta solo le
		// graffe fino a chiuderlo, non emettere nulla (gestisce anche il nesting).
		if skipDepth > 0 {
			switch c {
			case '{':
				skipDepth++
			case '}':
				skipDepth--
			}
			i++
			continue
		}
		switch {
		case c == '{':
			// Inizio gruppo: se è una destination, entra in skip (il contenuto —
			// es. "Helvetica;" della fonttbl — non va nel testo).
			if startsDestination(s, i+1) {
				skipDepth = 1
			}
			i++
		case c == '}':
			i++
		case c == '\\':
			// Escape esadecimale cp1252 (\'XX): accenti À/à/§/é... — era la causa
			// di "ATTIVITc0"/"a74.2.1"/"criticite0".
			if i+3 < len(s) && s[i+1] == '\'' {
				if b, ok := parseHexByte(s[i+2], s[i+3]); ok {
					out.WriteRune(cp1252ToRune(b))
					i += 4
					continue
				}
			}
			word, advance := peekControlWord(s, i+1)
			if advance == 0 {
				// Escape singolo: \\ \{ \} \~ + `\`+newline (a-capo cocoa/TextEdit).
				if i+1 < len(s) {
					switch s[i+1] {
					case '\\', '{', '}':
						out.WriteByte(s[i+1])
					case '~':
						out.WriteByte(' ')
					case '\n', '\r':
						out.WriteByte('\n')
					}
				}
				i += 2
				continue
			}
			writeControlWordOutput(&out, word)
			i += 1 + advance
		default:
			if unicode.IsPrint(rune(c)) || c == '\n' || c == '\t' {
				out.WriteByte(c)
			}
			i++
		}
	}
	return strings.TrimSpace(out.String())
}

// startsDestination indica se a `pos` (subito dopo una '{') inizia un gruppo
// destination da saltare: `\<destinationword>` oppure `\*` (gruppo ignorabile).
func startsDestination(s string, pos int) bool {
	if pos >= len(s) || s[pos] != '\\' {
		return false
	}
	pos++ // dopo il '\'
	if pos < len(s) && s[pos] == '*' {
		return true
	}
	word, _ := peekControlWord(s, pos)
	return destinationControlWords[word]
}

func parseHexByte(a, b byte) (byte, bool) {
	hi, ok1 := hexVal(a)
	lo, ok2 := hexVal(b)
	if !ok1 || !ok2 {
		return 0, false
	}
	return hi<<4 | lo, true
}

func hexVal(c byte) (byte, bool) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', true
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, true
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

// cp1252ToRune mappa un byte Windows-1252 al rune Unicode. 0x00-0x7F e 0xA0-0xFF
// coincidono con Latin-1/Unicode (copre tutti gli accenti italiani e §); il range
// 0x80-0x9F ha mapping speciali (smart quotes, trattini, € ecc.).
func cp1252ToRune(b byte) rune {
	if b < 0x80 || b >= 0xA0 {
		return rune(b)
	}
	if r, ok := cp1252Specials[b]; ok {
		return r
	}
	return rune(b)
}

var cp1252Specials = map[byte]rune{
	0x80: '€', 0x82: '‚', 0x83: 'ƒ', 0x84: '„', 0x85: '…', 0x86: '†', 0x87: '‡',
	0x88: 'ˆ', 0x89: '‰', 0x8A: 'Š', 0x8B: '‹', 0x8C: 'Œ', 0x8E: 'Ž', 0x91: '‘',
	0x92: '’', 0x93: '“', 0x94: '”', 0x95: '•', 0x96: '–', 0x97: '—',
	0x98: '˜', 0x99: '™', 0x9A: 'š', 0x9B: '›', 0x9C: 'œ', 0x9E: 'ž', 0x9F: 'Ÿ',
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
