package input

import "strings"

// WrapPlainAsRTF avvolge testo semplice (ASCII) in un RTF minimale valido,
// apribile/editabile in TextEdit/Word/Pages e ri-parsabile da parseRtf. È
// l'inverso di parseRtf per il caso starter: usato dal wizard `labnexus init`
// per creare un Prompt_INPUT.rtf che Denis edita nel suo editor preferito.
//
// Escapa i caratteri speciali RTF (`\` `{` `}`) e converte ogni newline in `\par`
// (niente newline letterali nel body: decodeRtf li tratterebbe come testo).
// Assunzione: input ASCII (lo starter lo è). Caratteri non-ASCII non sono
// \u-encodati; quando Denis salva da Word il file diventa RTF completo gestito
// nativamente.
func WrapPlainAsRTF(text string) string {
	var b strings.Builder
	b.WriteString("{\\rtf1\\ansi\\ansicpg1252\\deff0\n")
	b.WriteString("{\\fonttbl{\\f0\\fswiss Helvetica;}}\n")
	b.WriteString("\\f0\\fs24 ")
	for _, line := range strings.Split(text, "\n") {
		b.WriteString(escapeRTF(line))
		b.WriteString("\\par ")
	}
	b.WriteString("}\n")
	return b.String()
}

// escapeRTF protegge i tre caratteri con significato strutturale in RTF.
func escapeRTF(s string) string {
	return strings.NewReplacer(`\`, `\\`, `{`, `\{`, `}`, `\}`).Replace(s)
}
